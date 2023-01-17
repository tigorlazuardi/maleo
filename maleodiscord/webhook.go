package maleodiscord

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"

	"github.com/bwmarrin/snowflake"

	"github.com/tigorlazuardi/maleo"
	"github.com/tigorlazuardi/maleo/bucket"
)

type WebhookPayload struct {
	Wait            bool             `json:"-"`
	ThreadID        snowflake.ID     `json:"-"`
	Content         string           `json:"content,omitempty"`
	Username        string           `json:"username,omitempty"`
	AvatarURL       string           `json:"avatarURL,omitempty"`
	TTS             bool             `json:"TTS,omitempty"`
	Embeds          []*Embed         `json:"embeds,omitempty"`
	AllowedMentions *AllowedMentions `json:"allowed_mentions,omitempty"`
	Attachments     []*Attachment    `json:"attachments,omitempty"`
}

func (w *WebhookPayload) BuildMultipartPayloadJSON() ([]byte, error) {
	fields := map[string]any{}
	if w.Content != "" {
		fields["content"] = w.Content
	}
	if w.Username != "" {
		fields["username"] = w.Username
	}
	if w.AvatarURL != "" {
		fields["avatar_url"] = w.AvatarURL
	}
	if w.TTS {
		fields["tts"] = w.TTS
	}
	if len(w.Embeds) > 0 {
		fields["embeds"] = w.Embeds
	}
	if len(w.Attachments) > 0 {
		fields["attachments"] = w.Attachments
	}
	if w.AllowedMentions != nil {
		fields["allowed_mentions"] = w.AllowedMentions
	}
	return json.Marshal(fields)
}

type DiscordErrorResponse struct {
	Code       int             `json:"code"`
	Message    string          `json:"Message"`
	StatusCode int             `json:"status_code"`
	Raw        json.RawMessage `json:"raw"`
}

func newDiscordErrorResponse(statusCode int, body []byte) (*DiscordErrorResponse, error) {
	var errResp DiscordErrorResponse
	errResp.StatusCode = statusCode
	if err := json.Unmarshal(body, &errResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal discord error response: %w", err)
	}
	errResp.Raw = body
	return &errResp, nil
}

func (d DiscordErrorResponse) PrintJSON() {
	f, _ := json.Marshal(d)
	fmt.Println(string(f))
}

func (d DiscordErrorResponse) String() string {
	return fmt.Sprintf("discord error: [%d] %s", d.Code, d.Message)
}

func (d DiscordErrorResponse) Error() string {
	return d.String()
}

type WebhookFile struct {
	Name        string
	ContentType string
}

type AllowedMentions struct {
	Parse []string
	Roles []snowflake.ID
	Users []snowflake.ID
}

type Attachment struct {
	ID          int    `json:"id"`
	Filename    string `json:"filename,omitempty"`
	Description string `json:"description,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	Size        int    `json:"size,omitempty"`
	URL         string `json:"url,omitempty"`
	ProxyURL    string `json:"proxy_url,omitempty"`
	Height      int    `json:"height,omitempty"`
	Width       int    `json:"width,omitempty"`
	Ephemeral   bool   `json:"ephemeral,omitempty"`
}

type WebhookContext struct {
	Message maleo.MessageContext
	Files   []bucket.File
	Payload *WebhookPayload
	Extra   *ExtraInformation
	// Populated on PostMessageHook, otherwise Nil.
	//
	// Body is the response body from Discord.
	// Nil if the request failed to receive response from Discord.
	ResponseBody []byte
	// Populated on PostMessageHook if the request received response from Discord, otherwise nil.
	// Body is already consumed. Use ResponseBody instead to read the response body.
	Response *http.Response
	Maleo    *maleo.Maleo
}

func (d *Discord) PostWebhookJSON(ctx context.Context, web *WebhookContext) error {
	ctx = d.hook.PreMessageHook(ctx, web)
	var out bytes.Buffer
	enc := json.NewEncoder(&out)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(web.Payload); err != nil {
		return fmt.Errorf("failed to encode webhook payload: %w", err)
	}
	b := out.Bytes()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.webhook, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %w", err)
	}
	if web.Payload.Wait {
		req.URL.Query().Add("wait", "true")
	}
	req.URL.Query().Add("thread_id", web.Payload.ThreadID.String())
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute webhook: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	web.Response = resp
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		d.hook.PostMessageHook(ctx, web, err)
		return fmt.Errorf("failed to read webhook response body: %w", err)
	}
	web.ResponseBody = body
	if resp.StatusCode >= 400 {
		errResp, err := newDiscordErrorResponse(resp.StatusCode, body)
		if err != nil {
			d.hook.PostMessageHook(ctx, web, err)
			return fmt.Errorf("failed to parse discord error response: %w", err)
		}
		errResp.PrintJSON()
		d.hook.PostMessageHook(ctx, web, errResp)
		return errResp
	}
	d.hook.PostMessageHook(ctx, web, err)

	return nil
}

func (d *Discord) PostWebhookMultipart(ctx context.Context, web *WebhookContext) error {
	ctx = d.hook.PreMessageHook(ctx, web)
	requestBody, contentType, err := d.buildMultipartWebhookBody(web)
	if err != nil {
		return fmt.Errorf("failed to build multipart webhook body: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.webhook, requestBody)
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %w", err)
	}
	if web.Payload.Wait {
		req.URL.Query().Add("wait", "true")
	}
	req.URL.Query().Add("thread_id", web.Payload.ThreadID.String())
	req.Header.Set("Content-Type", contentType)

	resp, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute webhook: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		d.hook.PostMessageHook(ctx, web, err)
		return fmt.Errorf("failed to read webhook response body: %w", err)
	}
	if resp.StatusCode >= 400 {
		errResp, err := newDiscordErrorResponse(resp.StatusCode, body)
		if err != nil {
			d.hook.PostMessageHook(ctx, web, err)
			return fmt.Errorf("failed to parse discord error response: %w", err)
		}
		errResp.PrintJSON()
		d.hook.PostMessageHook(ctx, web, errResp)
		return errResp
	}
	d.hook.PostMessageHook(ctx, web, err)

	return nil
}

func (d *Discord) buildMultipartWebhookBody(web *WebhookContext) (body *bytes.Buffer, contentType string, err error) {
	body = &bytes.Buffer{}
	multipartWriter := multipart.NewWriter(body)
	defer func(multipartWriter *multipart.Writer) {
		_ = multipartWriter.Close()
	}(multipartWriter)
	contentType = multipartWriter.FormDataContentType()
	for i, file := range web.Files {
		err := func(i int, file bucket.File) error {
			defer func(file bucket.File) {
				_ = file.Close()
			}(file)
			fw, _ := multipartWriter.CreatePart(textproto.MIMEHeader{
				"Content-Disposition": {fmt.Sprintf(`form-data; name="files[%d]"; filename="%s"`, i, file.Filename())},
				"Content-Type":        {file.ContentType()},
			})
			_, err := io.Copy(fw, file)
			if err != nil {
				return err
			}
			var height, width int
			if img, ok := file.Data().(ImageSizeHint); ok {
				height, width = img.ImageSize()
			}
			web.Payload.Attachments = append(web.Payload.Attachments, &Attachment{
				ID:          i,
				Filename:    file.Filename(),
				Description: file.Pretext(),
				ContentType: file.ContentType(),
				Size:        file.Size(),
				Height:      height,
				Width:       width,
			})
			return nil
		}(i, file)
		if err != nil {
			return body, contentType, web.Maleo.
				Wrap(err).
				Caller(web.Message.Caller()).
				Message("failed to copy file data to multipart writer").
				Context(maleo.F{"index": i, "filename": file.Filename()}).
				Freeze()
		}
	}

	fw, _ := multipartWriter.CreatePart(textproto.MIMEHeader{
		"Content-Disposition": {"form-data; name=\"payload_json\""},
		"Content-Type":        {"application/json"},
	})
	j, err := web.Payload.BuildMultipartPayloadJSON()
	if err != nil {
		return body, contentType, web.Maleo.
			Wrap(err).
			Caller(web.Message.Caller()).
			Message("failed to copy file data to build payload json writer").
			Freeze()
	}
	_, _ = fw.Write(j)
	return body, contentType, nil
}

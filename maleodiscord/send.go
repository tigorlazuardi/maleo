package maleodiscord

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/tigorlazuardi/maleo"
)

func (d *Discord) send(ctx context.Context, msg maleo.MessageContext) {
	key := d.buildKey(msg)
	ticker := time.NewTicker(time.Millisecond * 300)
	for d.lock.Exist(ctx, d.globalKey) {
		<-ticker.C
	}
	id := d.snowflake.Generate()
	extra := &ExtraInformation{CacheKey: key, ThreadID: id}
	ticker.Stop()
	if err := d.lock.Set(ctx, d.globalKey, []byte("locked"), time.Second*30); err != nil {
		_ = msg.Maleo().Wrap(err).Caller(msg.Caller()).Message("%s: failed to set global lock to lock",
			d.Name()).Log(ctx)
	}
	if msg.ForceSend() {
		extra.CooldownTimeEnds = time.Now().Add(time.Second * 2)
		_ = d.postMessage(ctx, msg, extra)
		d.deleteGlobalCacheKeyAfter2Seconds(ctx)
		return
	}
	if d.lock.Exist(ctx, key) {
		d.lock.Delete(ctx, d.globalKey)
		return
	}
	defer d.deleteGlobalCacheKeyAfter2Seconds(ctx)
	iterKey := key + d.lock.Separator() + "iter"
	iter := d.getAndSetIter(ctx, iterKey)
	cooldown := d.countCooldown(msg, iter)
	extra.Iteration = iter
	extra.CooldownTimeEnds = time.Now().Add(cooldown)
	err := d.postMessage(ctx, msg, extra)
	if err == nil {
		message := msg.Message()
		if msg.Err() != nil {
			message = msg.Err().Error()
		}
		if err := d.lock.Set(ctx, key, []byte(message), d.countCooldown(msg, iter)); err != nil {
			_ = msg.Maleo().
				Wrap(err).
				Message("%s: failed to set Message key to lock", d.Name()).
				Caller(msg.Caller()).
				Context(maleo.F{"key": key, "payload": message}).
				Log(ctx)
		}
	}
}

func (d *Discord) deleteGlobalCacheKeyAfter2Seconds(ctx context.Context) {
	time.Sleep(time.Second * 2)
	d.lock.Delete(ctx, d.globalKey)
}

func buildIntro(service maleo.Service, err error) string {
	var s strings.Builder
	if err != nil {
		s.WriteString("@here an error has occurred")
	} else {
		s.WriteString("@here Message")
	}
	if service.Name != "" {
		s.WriteString(" on service **")
		s.WriteString(service.Name)
		s.WriteString("**")
	}
	if service.Type != "" {
		s.WriteString(" on type **")
		s.WriteString(service.Type)
		s.WriteString("**")
	}
	if service.Environment != "" {
		s.WriteString(" on environment **")
		s.WriteString(service.Environment)
		s.WriteString("**")
	}
	return s.String()
}

func (d *Discord) postMessage(ctx context.Context, msg maleo.MessageContext, extra *ExtraInformation) error {
	service := msg.Service()
	err := msg.Err()
	intro := buildIntro(service, err)
	if extra.ThreadID == 0 {
		extra.ThreadID = d.snowflake.Generate()
	}

	embeds, files := d.builder.BuildEmbed(ctx, msg, extra)
	payload := &WebhookPayload{
		Wait:     true,
		ThreadID: extra.ThreadID,
		Content:  intro,
		Embeds:   embeds,
	}

	webhookContext := &WebhookContext{
		Message: msg,
		Files:   files,
		Payload: payload,
		Extra:   extra,
		Maleo:   msg.Maleo(),
	}

	switch {
	case d.bucket != nil && len(files) > 0:
		payload, errUpload := d.bucketUpload(ctx, webhookContext)
		webhookContext.Payload = payload
		err := d.PostWebhookJSON(ctx, webhookContext)
		switch {
		case err != nil:
			return err
		case errUpload != nil:
			return errUpload
		default:
			return nil
		}
	case len(files) > 0:
		return d.PostWebhookMultipart(ctx, webhookContext)
	}
	return d.PostWebhookJSON(ctx, webhookContext)
}

func (d *Discord) bucketUpload(ctx context.Context, web *WebhookContext) (*WebhookPayload, error) {
	ctx = d.hook.PreBucketUploadHook(ctx, web)
	results := d.bucket.Upload(ctx, web.Files)
	d.hook.PostBucketUploadHook(ctx, web, results)
	payload := web.Payload
	errs := make([]error, 0, len(results))
	for i, result := range results {
		if result.Error != nil {
			errs = append(errs, result.Error)
			continue
		}
		var height, width int
		if imgHint, ok := result.File.Data().(ImageSizeHint); ok {
			height, width = imgHint.ImageSize()
		}
		payload.Attachments = append(payload.Attachments, &Attachment{
			ID:          i,
			Filename:    result.File.Filename(),
			Description: result.File.Pretext(),
			ContentType: result.File.ContentType(),
			Size:        result.File.Size(),
			URL:         result.URL,
			Height:      height,
			Width:       width,
		})
	}
	if len(errs) > 0 {
		return payload, maleo.
			Bail("failed to upload some file(s) to bucket").
			Caller(web.Message.Caller()).
			Context(maleo.F{"errors": errs}).
			Freeze()
	}
	return payload, nil
}

func (d *Discord) buildKey(msg maleo.MessageContext) string {
	builder := strings.Builder{}
	builder.WriteString(d.Name())
	builder.WriteString(d.lock.Separator())
	service := msg.Service()
	builder.WriteString(service.Environment)
	builder.WriteString(d.lock.Separator())
	builder.WriteString(service.Name)
	builder.WriteString(d.lock.Separator())
	builder.WriteString(service.Type)
	builder.WriteString(d.lock.Separator())

	key := msg.Key()
	if key == "" {
		key = msg.Caller().FormatAsKey()
	}
	builder.WriteString(key)
	return builder.String()
}

func (d *Discord) countCooldown(msg maleo.MessageContext, iter int) time.Duration {
	multiplier := (iter * iter) >> 1
	if multiplier < 1 {
		multiplier = 1
	}
	cooldown := msg.Cooldown()
	if cooldown == 0 {
		cooldown = d.cooldown
	}
	cooldown *= time.Duration(multiplier)
	if cooldown > time.Hour*24 {
		cooldown = time.Hour * 24
	}
	return cooldown
}

func (d *Discord) getAndSetIter(ctx context.Context, key string) int {
	var iter int
	iterByte, err := d.lock.Get(ctx, key)
	if err == nil {
		iter, _ = strconv.Atoi(string(iterByte))
	}
	iter += 1
	iterByte = []byte(strconv.Itoa(iter))
	nextCooldown := d.cooldown*time.Duration(iter) + d.cooldown
	_ = d.lock.Set(ctx, key, iterByte, nextCooldown)
	return iter
}

package maleodiscord

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/tigorlazuardi/maleo"
	"github.com/tigorlazuardi/maleo/bucket"
)

// discordLimit is actually 6000, but we have to reserve 1000 characters create attachments.
const discordLimit = 5000

func (d *Discord) defaultEmbedBuilder(
	ctx context.Context,
	msg maleo.MessageContext,
	extra *ExtraInformation,
) ([]*Embed, []bucket.File) {
	var (
		files  = make([]bucket.File, 0, 5)
		embeds = make([]*Embed, 0, 5)
		limit  = discordLimit - 150 // we have to take account for titles and timestamps.
	)
	summary, fileSummary, written := d.buildSummary(msg, 500, extra)
	limit -= written

	metadata, fileMetadata, written := d.buildMetadataEmbed(ctx, msg, extra, 500)
	limit -= written

	errorStackEmbed, fileErrorStack, written := d.buildErrorStackEmbed(msg, 1000, extra)
	limit -= written

	// Data limit is 50% of the remaining limit at max when error is available, otherwise 100% until 4096.
	dataLimit := limit
	if msg.Err() == nil && dataLimit > 4096 {
		dataLimit = 4096
	} else {
		dataLimit /= 2
	}

	dataEmbed, fileData, written := d.buildContextEmbed(msg, dataLimit, extra)
	limit -= written

	if limit > 4096 {
		limit = 4096
	}

	// Error will take the remaining limit.
	errorEmbed, errorData, written := d.buildErrorEmbed(msg, limit, extra)

	embeds = append(embeds, summary)
	if errorEmbed != nil {
		embeds = append(embeds, errorEmbed)
	}
	if dataEmbed != nil {
		embeds = append(embeds, dataEmbed)
	}
	if errorStackEmbed != nil {
		embeds = append(embeds, errorStackEmbed)
	}
	embeds = append(embeds, metadata)

	if fileSummary != nil {
		files = append(files, fileSummary)
	}
	if errorData != nil {
		files = append(files, errorData)
	}
	if fileData != nil {
		files = append(files, fileData)
	}
	if fileErrorStack != nil {
		files = append(files, fileErrorStack)
	}
	if fileMetadata != nil {
		files = append(files, fileMetadata)
	}
	return embeds, files
}

//goland:noinspection GoUnhandledErrorResult
func (d *Discord) buildContextEmbed(
	msg maleo.MessageContext,
	limit int,
	extra *ExtraInformation,
) (*Embed, bucket.File, int) {
	if len(msg.Context()) == 0 {
		return nil, nil, 0
	}
	embed := &Embed{
		Type:  "rich",
		Title: "Context",
		Color: 0x063970, // Dark Blue
	}

	display, data := new(bytes.Buffer), new(bytes.Buffer)
	display.Reset()
	display.Grow(limit)
	data.Reset()
	data.Grow(limit)

	contextData := msg.Context()
	err := d.codeBlockBuilder.Build(display, contextData)
	if err != nil {
		_, _ = display.WriteString("Error building Context: ")
		display.WriteString("```")
		_, _ = display.WriteString(err.Error())
		display.WriteString("```\n")
	}
	if display.Len() > limit {
		var v any = contextData
		if len(msg.Context()) == 1 {
			v = contextData[0]
		}
		err := d.dataEncoder.Encode(data, v)
		if err != nil {
			display.WriteString("Error encoding Context to file: ")
			display.WriteString("```")
			display.WriteString(err.Error())
			display.WriteString("```\n")
		}
	}

	return shouldCreateFile(&createFileContext{
		embed:          embed,
		display:        display,
		data:           data,
		contentType:    d.dataEncoder.ContentType(),
		fileExtension:  d.dataEncoder.FileExtension(),
		suffixFilename: "_context",
		limit:          limit,
		extra:          extra,
	})
}

func (d *Discord) buildErrorEmbed(
	msg maleo.MessageContext,
	limit int, extra *ExtraInformation,
) (*Embed, bucket.File, int) {
	err := msg.Err()
	if err == nil {
		return nil, nil, 0
	}
	embed := &Embed{
		Type:  "rich",
		Title: "Error",
		Color: 0x71010b, // Venetian Red
	}
	display, data := new(bytes.Buffer), new(bytes.Buffer)
	display.Reset()
	display.Grow(limit)
	data.Reset()
	data.Grow(limit)
	if err := d.codeBlockBuilder.BuildError(display, err); err != nil {
		_, _ = display.WriteString("Error building error as display: ")
		_, _ = display.WriteString("```")
		_, _ = display.WriteString(err.Error())
		_, _ = display.WriteString("```\n")
	}
	if display.Len() > limit {
		err := d.dataEncoder.Encode(data, err)
		if err != nil {
			_, _ = display.WriteString("Error encoding error to file: ")
			_, _ = display.WriteString("```")
			_, _ = display.WriteString(err.Error())
			_, _ = display.WriteString("```\n")
		}
	}
	return shouldCreateFile(&createFileContext{
		embed:          embed,
		display:        display,
		data:           data,
		contentType:    d.dataEncoder.ContentType(),
		fileExtension:  d.dataEncoder.FileExtension(),
		suffixFilename: "_error",
		limit:          limit,
		extra:          extra,
	})
}

func (d *Discord) buildErrorStackEmbed(
	msg maleo.MessageContext,
	limit int, extra *ExtraInformation,
) (*Embed, bucket.File, int) {
	err := msg.Err()
	if err == nil {
		return nil, nil, 0
	}
	s := make([]string, 0, 4)
	s = stackAccumulator(s, msg.Err())

	if len(s) == 0 {
		return nil, nil, 0
	}
	reverse(s)
	content := strings.Join(s, "\n---\n")
	display, data := new(bytes.Buffer), new(bytes.Buffer)
	display.Reset()
	display.Grow(limit)
	display.WriteString("```")
	display.WriteString(content)
	display.WriteString("```")
	embed := &Embed{
		Type:  "rich",
		Title: "Error Stack",
		Color: 0x5d0e16, // Cardinal Red Dark
	}
	if display.Len() > limit {
		data.WriteString(content)
	}
	return shouldCreateFile(&createFileContext{
		embed:          embed,
		display:        display,
		data:           data,
		contentType:    "text/plain; charset=utf-8",
		fileExtension:  "txt",
		suffixFilename: "_error_stack",
		limit:          limit,
		extra:          extra,
	})
}

func stackAccumulator(s []string, err error) []string {
	if err == nil {
		return s
	}
	ss := &strings.Builder{}
	chWritten := false
	if ch, ok := err.(maleo.CallerHint); ok {
		chWritten = true
		ss.WriteString(ch.Caller().String())
	}
	if chWritten {
		if mh, ok := err.(maleo.MessageHint); ok {
			ss.WriteString(": ")
			ss.WriteString(mh.Message())
		}
	}
	if ss.Len() > 0 {
		s = append(s, ss.String())
	}
	return stackAccumulator(s, errors.Unwrap(err))
}

func reverse[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func closingTicksTruncated(b *bytes.Buffer, countBack int) bool {
	buf := b.Bytes()
	if len(buf) >= countBack {
		buf = buf[len(buf)-countBack:]
	}
	count := bytes.Count(buf, []byte("```"))
	return count%2 != 0
}

type createFileContext struct {
	embed          *Embed
	display        *bytes.Buffer
	data           *bytes.Buffer
	contentType    string
	fileExtension  string
	suffixFilename string
	limit          int
	extra          *ExtraInformation
}

func shouldCreateFile(ctx *createFileContext) (em *Embed, file bucket.File, written int) {
	display := ctx.display
	if display.Len() > ctx.limit {
		outro := "Content is too long to be displayed fully. See attachment for details"
		if closingTicksTruncated(display, len(outro)+5) {
			outro = "\n```\nContent is too long to be displayed fully. See attachment for details"
		}
		display.Truncate(ctx.limit - len(outro))
		display.WriteString(outro)
		ctx.embed.Description = display.String()

		filename := fmt.Sprintf("%s%s.%s", ctx.extra.ThreadID, ctx.suffixFilename, ctx.fileExtension)
		file = bucket.NewFile(
			ctx.data,
			ctx.contentType,
			bucket.WithFilename(filename),
		)
		return ctx.embed, file, display.Len()
	}
	ctx.embed.Description = display.String()
	return ctx.embed, nil, display.Len()
}

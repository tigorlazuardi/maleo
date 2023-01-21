package maleodiscord

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/tigorlazuardi/maleo"
	"github.com/tigorlazuardi/maleo/bucket"
)

func (d *Discord) buildMetadataEmbed(
	ctx context.Context,
	msg maleo.MessageContext,
	extra *ExtraInformation,
	limit int,
) (*Embed, bucket.File, int) {
	count := 0
	embed := &Embed{
		Type:      "rich",
		Title:     "Metadata",
		Color:     0x645a5b, // Scorpion Grey
		Timestamp: msg.Time().Format(time.RFC3339),
	}
	count = buildMetadataEmbedFields(ctx, msg, extra, d, embed, count)
	display := buildMetadataBodyCaller(limit, msg, extra)

	embed, file, written := shouldCreateFile(&createFileContext{
		embed:          embed,
		display:        display,
		data:           bytes.NewBufferString(display.String()),
		contentType:    "text/markdown; charset=utf-8",
		fileExtension:  "md",
		suffixFilename: "_metadata",
		limit:          limit,
		extra:          extra,
	})
	count += written
	return embed, file, count
}

func buildMetadataBodyCaller(limit int, msg maleo.MessageContext, extra *ExtraInformation) *bytes.Buffer {
	display := new(bytes.Buffer)
	display.Grow(limit)
	display.WriteString(`**Caller Origin**`)
	display.WriteString("\n```\n")
	display.WriteString(msg.Caller().String())
	display.WriteString("\n```\n")
	display.WriteString(`**Caller Function**`)
	display.WriteString("\n```\n")
	display.WriteString(msg.Caller().ShortName())
	display.WriteString("\n```\n")
	display.WriteString(`**Cache Key**`)
	display.WriteString("\n```\n")
	display.WriteString(extra.CacheKey)
	display.WriteString("\n```")

	return display
}

func buildMetadataEmbedFields(ctx context.Context, msg maleo.MessageContext, extra *ExtraInformation, d *Discord, embed *Embed, count int) int {
	count = buildTraceEmbedFields(ctx, d, embed, count)
	service := msg.Service()
	count = buildServiceEmbedFields(service, embed, count)
	const threadIDName = "Thread ID"
	embed.Fields = append(embed.Fields, &EmbedField{
		Name:   threadIDName,
		Value:  extra.ThreadID.String(),
		Inline: true,
	})
	count += len(threadIDName) + len(extra.ThreadID.String())
	var iteration string
	if msg.ForceSend() {
		iteration = "(Force Send)"
	} else {
		iteration = strconv.Itoa(extra.Iteration)
	}
	const messageIteration = "Message Iteration"
	embed.Fields = append(embed.Fields, &EmbedField{
		Name:   messageIteration,
		Value:  iteration,
		Inline: true,
	})
	count += len(messageIteration) + len(iteration)
	ts := extra.CooldownTimeEnds.Unix()
	const nextPossibleEarliestRepeat = "Next Possible Earliest Repeat"
	repeatValue := fmt.Sprintf("<t:%d:F> | <t:%d:R>", ts, ts)
	embed.Fields = append(embed.Fields, &EmbedField{
		Name:   nextPossibleEarliestRepeat,
		Value:  repeatValue,
		Inline: false,
	})
	count += len(messageIteration) + len(iteration)
	if len(embed.Fields) > 25 {
		embed.Fields = embed.Fields[:25]
	}
	return count
}

func buildTraceEmbedFields(ctx context.Context, d *Discord, embed *Embed, count int) int {
	for _, v := range d.trace.CaptureTrace(ctx) {
		embed.Fields = append(embed.Fields, &EmbedField{
			Name:   v.Key,
			Value:  v.Value,
			Inline: true,
		})
		count += len(v.Key) + len(v.Value)
	}
	return count
}

func buildServiceEmbedFields(service maleo.Service, embed *Embed, count int) int {
	if service.Name != "" {
		const name = "Service"
		embed.Fields = append(embed.Fields, &EmbedField{
			Name:   name,
			Value:  service.Name,
			Inline: true,
		})
		count += len(name) + len(service.Name)
	}
	if service.Type != "" {
		const sType = "Type"
		embed.Fields = append(embed.Fields, &EmbedField{
			Name:   sType,
			Value:  service.Type,
			Inline: true,
		})
		count += len(sType) + len(service.Type)
	}
	if service.Environment != "" {
		const env = "Environment"
		embed.Fields = append(embed.Fields, &EmbedField{
			Name:   env,
			Value:  service.Environment,
			Inline: true,
		})
		count += len(env) + len(service.Type)
	}
	return count
}

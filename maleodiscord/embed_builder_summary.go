package maleodiscord

import (
	"bytes"

	"github.com/tigorlazuardi/maleo"
	"github.com/tigorlazuardi/maleo/bucket"
)

func (d *Discord) buildSummary(
	msg maleo.MessageContext,
	limit int,
	extra *ExtraInformation,
) (*Embed, bucket.File, int) {
	embed := &Embed{
		Type:  "rich",
		Title: "Summary",
		Color: 0x188544, // Green Jewel
	}
	display, data := new(bytes.Buffer), new(bytes.Buffer)

	buildSummaryPretext(msg, display)
	buildSummaryError(msg, display)

	buildSummaryContext(msg, display)
	if display.Len() > limit {
		data.Write(display.Bytes())
	}

	return shouldCreateFile(&createFileContext{
		embed:          embed,
		display:        display,
		data:           data,
		contentType:    "text/markdown; charset=utf-8",
		fileExtension:  "md",
		suffixFilename: "_summary",
		limit:          limit,
		extra:          extra,
	})
}

func buildSummaryContext(msg maleo.MessageContext, display *bytes.Buffer) {
	dataContext := msg.Context()
	if len(dataContext) > 0 {
		for _, c := range dataContext {
			switch c := c.(type) {
			case maleo.SummaryWriter:
				display.WriteString("\n\n**Context**:\n")
				display.WriteString("```")
				if _, ok := c.(maleo.Fields); ok {
					display.WriteString("yaml")
				}
				display.WriteString("\n")
				lw := maleo.NewLineWriter(display).LineBreak("\n").Build()
				c.WriteSummary(lw)
				_, _ = display.WriteString("\n```")
			case maleo.Summary:
				display.WriteString("\n\n**Context**:\n")
				display.WriteString("```\n")
				display.WriteString(c.Summary())
				display.WriteString("\n```")
			}
		}
	}
}

func buildSummaryError(msg maleo.MessageContext, display *bytes.Buffer) {
	err := msg.Err()
	if err != nil {
		display.WriteString("\n\n**Error**:\n")
		display.WriteString("```\n")
		switch err := err.(type) {
		case maleo.SummaryWriter:
			lw := maleo.NewLineWriter(display).LineBreak("\n").Build()
			err.WriteSummary(lw)
		case maleo.Summary:
			display.WriteString(err.Summary())
		case maleo.ErrorWriter:
			lw := maleo.NewLineWriter(display).LineBreak("\n.. ").Build()
			err.WriteError(lw)
		default:
			display.WriteString(err.Error())
		}
		display.WriteString("\n```")
	}
}

func buildSummaryPretext(msg maleo.MessageContext, display *bytes.Buffer) {
	display.WriteString("**")
	display.WriteString(msg.Message())
	display.WriteString("**")
}

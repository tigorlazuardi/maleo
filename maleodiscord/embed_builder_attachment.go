package maleodiscord

import (
	"strings"

	"github.com/tigorlazuardi/maleo/bucket"
)

func buildAttachmentEmbed(uploads []bucket.UploadResult) *Embed {
	embed := &Embed{
		Type:  "rich",
		Title: "Attachment",
		Color: 0x188544, // Green Jewel
	}

	for _, upload := range uploads {
		nameField := getCategoryFromFilename(upload.File.Filename())
		if upload.Error != nil {
			embed.Fields = append(embed.Fields, &EmbedField{
				Name:  nameField,
				Value: upload.Error.Error(),
			})
		} else {
			embed.Fields = append(embed.Fields, &EmbedField{
				Name:  nameField,
				Value: upload.URL,
			})
		}
	}
	return embed
}

func getCategoryFromFilename(filename string) string {
	switch {
	case strings.Contains(filename, "_stack"):
		return "Stack"
	case strings.Contains(filename, "_summary"):
		return "Summary"
	case strings.Contains(filename, "_context"):
		return "Context"
	case strings.Contains(filename, "_error"):
		return "Error"
	case strings.HasSuffix(filename, ".png"),
		strings.HasSuffix(filename, ".svg"),
		strings.HasSuffix(filename, ".webp"),
		strings.HasSuffix(filename, ".jpg"),
		strings.HasSuffix(filename, ".jpeg"):
		return "Image"
	default:
		return filename
	}
}

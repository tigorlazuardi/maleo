package maleodiscord

import (
	"context"
	"time"

	"github.com/bwmarrin/snowflake"

	"github.com/tigorlazuardi/maleo"
	"github.com/tigorlazuardi/maleo/bucket"
)

type EmbedBuilder interface {
	BuildEmbed(ctx context.Context, msg maleo.MessageContext, info *ExtraInformation) ([]*Embed, []bucket.File)
}

type ExtraInformation struct {
	Iteration        int
	CooldownTimeEnds time.Time
	CacheKey         string
	ThreadID         snowflake.ID
}

type EmbedBuilderFunc func(ctx context.Context, msg maleo.MessageContext, info *ExtraInformation) ([]*Embed, []bucket.File)

func (e EmbedBuilderFunc) BuildEmbed(ctx context.Context, msg maleo.MessageContext, info *ExtraInformation) ([]*Embed, []bucket.File) {
	return e(ctx, msg, info)
}

type Embed struct {
	// Title of embed.
	Title string `json:"title,omitempty"`
	// Type of embed (always "rich" for webhook embeds).
	Type string `json:"type,omitempty"`
	// Description of embed.
	Description string `json:"description,omitempty"`
	// URL of embed.
	Url string `json:"url,omitempty"`
	// Timestamp of embed content.
	Timestamp string `json:"timestamp,omitempty"`
	// Color code of the embed.
	Color int `json:"color,omitempty"`
	// Footer information.
	Footer *EmbedFooter `json:"footer,omitempty"`
	// Image information.
	Image *EmbedImage `json:"image,omitempty"`
	// Thumbnail information.
	Thumbnail *EmbedThumbnail `json:"thumbnail,omitempty"`
	// Video information.
	Video *EmbedVideo `json:"video,omitempty"`
	// Provider information.
	Provider *EmbedProvider `json:"provider,omitempty"`
	// Author information.
	Author *EmbedAuthor `json:"author,omitempty"`
	// Fields information.
	Fields []*EmbedField `json:"fields,omitempty"`
}

type EmbedAuthor struct {
	// Name of author.
	Name string `json:"name,omitempty"`
	// URL of author.
	Url string `json:"url,omitempty"`
	// URL of author icon (only supports http(s) and attachments).
	IconUrl string `json:"icon_url,omitempty"`
	// A proxied url of author icon.
	ProxyIconUrl string `json:"proxy_icon_url,omitempty"`
}

type EmbedField struct {
	// Name of the field.
	Name string `json:"name"`
	// Value of the field.
	Value string `json:"value"`
	// Inline states Whether this field should display inline.
	Inline bool `json:"inline,omitempty"`
}

type EmbedFooter struct {
	// Footer text.
	Text string `json:"text"`
	// URL of footer icon (only supports http(s) and attachments).
	IconUrl string `json:"icon_url,omitempty"`
	// A proxied url of footer icon.
	ProxyIconUrl string `json:"proxy_icon_url,omitempty"`
}

type EmbedImage struct {
	// Source url of image (only supports http(s) and attachments).
	Url string `json:"url,omitempty"`
	// A proxied url of the image.
	ProxyUrl string `json:"proxy_url,omitempty"`
	// Height of image.
	Height int `json:"height,omitempty"`
	// Width of image.
	Width int `json:"width,omitempty"`
}

type EmbedProvider struct {
	// Name of provider.
	Name string `json:"name,omitempty"`
	// URL of provider.
	Url string `json:"url,omitempty"`
}

type EmbedThumbnail struct {
	// Source url of thumbnail (only supports http(s) and attachments).
	Url string `json:"url,omitempty"`
	// A proxied url of the thumbnail.
	ProxyUrl string `json:"proxy_url,omitempty"`
	// Height of thumbnail.
	Height int `json:"height,omitempty"`
	// Width of thumbnail.
	Width int `json:"width,omitempty"`
}

type EmbedVideo struct {
	// Source url of video.
	Url string `json:"url,omitempty"`
	// A proxied url of the video.
	ProxyUrl string `json:"proxy_url,omitempty"`
	// Height of video.
	Height int `json:"height,omitempty"`
	// Width of video.
	Width int `json:"width,omitempty"`
}

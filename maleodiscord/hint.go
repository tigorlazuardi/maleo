package maleodiscord

type HighlightHint interface {
	DiscordHighlight() string
}

type MimetypeHint interface {
	Mimetype() string
}

type ImageSizeHint interface {
	ImageSize() (height, width int)
}

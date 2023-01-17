package maleodiscord

import (
	"encoding/json"
	"io"
)

type DataEncoder interface {
	ContentType() string
	Encode(w io.Writer, value any) error
	FileExtension() string
}

var _ DataEncoder = (*JSONDataEncoder)(nil)

type JSONDataEncoder struct{}

func (J JSONDataEncoder) FileExtension() string {
	return "json"
}

func (J JSONDataEncoder) ContentType() string {
	return "application/json"
}

func (J JSONDataEncoder) Encode(w io.Writer, value any) error {
	const indent = "   "
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", indent)
	err := enc.Encode(value)
	return err
}

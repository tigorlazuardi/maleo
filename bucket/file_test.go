package bucket

import (
	"bytes"
	"reflect"
	"testing"
)

func TestNewFile(t *testing.T) {
	origin := new(bytes.Buffer)
	origin.WriteString("Hello, world!")
	mimetype := "text/plain; charset=utf-8"

	file := NewFile(origin, mimetype, WithFilename("test.txt"), WithPretext("pretext"), WithFilesize(origin.Len()))

	if file.Filename() != "test.txt" {
		t.Errorf("Expected filename to be 'test.txt', got '%s'", file.Filename())
	}
	if file.ContentType() != mimetype {
		t.Errorf("Expected mimetype to be '%s', got '%s'", mimetype, file.ContentType())
	}

	if file.Pretext() != "pretext" {
		t.Errorf("Expected pretext to be 'pretext', got '%s'", file.Pretext())
	}

	if file.Size() != origin.Len() {
		t.Errorf("Expected size to be %d, got %d", origin.Len(), file.Size())
	}

	if !reflect.DeepEqual(file.Data(), origin) {
		t.Errorf("Expected data to be equal to origin, got %v", file.Data())
	}

	if err := file.Close(); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	b := make([]byte, origin.Len())
	l := origin.Len()
	n, err := file.Read(b)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if n != l {
		t.Errorf("expected n to be 0, got %d", n)
	}
	if !reflect.DeepEqual(string(b), "Hello, world!") {
		t.Errorf("Expected read bytes to be '%s' equal to origin, got '%s'", "Hello, World!", b)
	}
}

package loader

import (
	"os"
	"testing"
)

func TestLoadEnv(t *testing.T) {
	LoadEnv()

	if os.Getenv("TEST_ENV") != "test" {
		t.Error("TEST_ENV should have value of 'test'")
	}

	if os.Getenv("TEST_ENV2") != "test2" {
		t.Error("TEST_ENV2 should have value of 'test2'")
	}

	if os.Getenv("TEST_NOT_EXIST") != "" {
		t.Error("TEST_NOT_EXIST should be empty")
	}

	if os.Getenv("TEST_ENV3") != "test3=123" {
		t.Error("TEST_ENV3 should have value of 'test3=123'")
	}

	if os.Getenv("TEST_EMPTY") != "" {
		t.Error("TEST_EMPTY should be empty")
	}
}

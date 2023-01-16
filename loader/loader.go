// Package loader is a simplified env loader.
//
// It does not support fancy stuffs other library gives. It's only meant for testing purposes.
//
// It also ignores errors from OS related APIs (although printed), like permissions. So do not use for live runtime environment.
package loader

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LoadEnv Search current directory and it's children and search for `.env` files then loads the environment variables inside.
//
// It ignores errors from OS related APIs (although printed), like permissions. So do not use for live runtime environment.
func LoadEnv() {
	_ = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".env" {
			setEnv(path)
		}
		return nil
	})
}

func setEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}
		// Ignore comments.
		if text[0] == '#' {
			continue
		}

		// Discard comments post text and clean up whitspaces.
		if index := strings.Index(text, "#"); index != -1 {
			text = text[:index]
			text = strings.TrimSpace(text)
		}
		parts := strings.SplitN(text, "=", 2)
		// Ignore invalid lines.
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		_ = os.Setenv(key, value)
	}
}

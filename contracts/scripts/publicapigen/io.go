package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func writeJSON(path string, value any) {
	content, err := json.MarshalIndent(value, "", "  ")
	must(err)
	content = append(content, '\n')
	must(os.MkdirAll(filepath.Dir(path), 0o755))
	must(os.WriteFile(path, content, 0o644))
}

func writeText(path, content string) {
	must(os.MkdirAll(filepath.Dir(path), 0o755))
	must(os.WriteFile(path, []byte(strings.TrimSpace(content)+"\n"), 0o644))
}

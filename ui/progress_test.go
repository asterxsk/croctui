package ui

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestFormatSize(t *testing.T) {
	cases := []struct {
		in   int64
		want string
	}{
		{0, "0 B"},
		{12, "12 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1024 * 1024, "1.0 MB"},
	}
	for _, c := range cases {
		if got := formatSize(c.in); got != c.want {
			t.Fatalf("formatSize(%d) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestDownloadsDir(t *testing.T) {
	d, err := downloadsDir()
	if err != nil {
		t.Fatalf("downloadsDir: %v", err)
	}
	if d == "" {
		t.Fatal("downloadsDir returned empty")
	}
	// Path must end in Downloads and be absolute.
	if !strings.HasSuffix(d, "Downloads") || !filepath.IsAbs(d) {
		t.Fatalf("unexpected downloads dir: %q", d)
	}
}

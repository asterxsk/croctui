package dialog

import (
	"strings"
	"testing"
)

func TestScriptEnablesMultiselect(t *testing.T) {
	s := dialogScript()
	if !strings.Contains(s, "$dialog.Multiselect = $true") {
		t.Fatalf("script missing Multiselect=true:\n%s", s)
	}
}

func TestScriptSetsUtf8Encoding(t *testing.T) {
	s := dialogScript()
	if !strings.Contains(s, "OutputEncoding") {
		t.Fatalf("script missing OutputEncoding:\n%s", s)
	}
}

func TestFolderScriptUsesFolderBrowser(t *testing.T) {
	s := folderScript()
	if !strings.Contains(s, "FolderBrowserDialog") {
		t.Fatalf("folder script missing FolderBrowserDialog:\n%s", s)
	}
	if !strings.Contains(s, "SelectedPath") {
		t.Fatalf("folder script missing SelectedPath:\n%s", s)
	}
}

// TestFolderScriptOpensAtCwd pins that the folder dialog opens at the current
// directory (not the default Desktop/MyComputer root).
func TestFolderScriptOpensAtCwd(t *testing.T) {
	s := folderScript()
	if !strings.Contains(s, "GetCurrentDirectory") || !strings.Contains(s, "SelectedPath = $cwd") {
		t.Fatalf("folder script should open at the current directory:\n%s", s)
	}
}

// TestScriptOpensAtCwd pins that the file dialog also opens at the current
// directory instead of a fixed default location.
func TestScriptOpensAtCwd(t *testing.T) {
	s := dialogScript()
	if !strings.Contains(s, "GetCurrentDirectory") || !strings.Contains(s, "InitialDirectory = $cwd") {
		t.Fatalf("file dialog script should open at the current directory:\n%s", s)
	}
}

func TestParsePaths(t *testing.T) {
	got := parsePaths("C:\\a\n\n  C:\\b  \n")
	want := []string{"C:\\a", "C:\\b"}
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("parsePaths(%q) = %v, want %v", "C:\\a\n\n  C:\\b  \n", got, want)
	}
	if len(parsePaths("")) != 0 {
		t.Fatal("parsePaths(empty) should return empty slice")
	}
}

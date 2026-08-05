package ui

import (
	"testing"
)

func TestBuildCrocCmd(t *testing.T) {
	cmd := buildCrocCmd("croc", []string{"--yes", "send", "a.txt"}, `C:\tmp`)
	if cmd == nil {
		t.Fatal("buildCrocCmd returned nil")
	}
	// Path will be resolved to full path by exec.Command
	if cmd.Dir != `C:\tmp` {
		t.Fatalf("cmd.Dir = %q, want %q", cmd.Dir, `C:\tmp`)
	}
	if len(cmd.Args) != 4 {
		t.Fatalf("cmd.Args = %v", cmd.Args)
	}
	// Stdout/Stderr must be wired so output can be streamed.
	if cmd.Stdout == nil || cmd.Stderr == nil {
		t.Fatal("Stdout and Stderr must be non-nil")
	}
	// Args should include the binary name as first arg
	if cmd.Args[0] != "croc" {
		t.Fatalf("cmd.Args[0] = %q, want %q", cmd.Args[0], "croc")
	}
}
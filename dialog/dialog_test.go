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

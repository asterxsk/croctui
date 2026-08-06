// dialog/dialog_manual_test.go  — skipped in CI; run with CROCGUI_MANUAL=1 go test ./...
package dialog

import (
	"os"
	"testing"
)

func TestManualPick(t *testing.T) {
	if os.Getenv("CROCGUI_MANUAL") == "" {
		t.Skip("set CROCGUI_MANUAL=1 to run the manual dialog test")
	}
	files, err := PickFiles()
	t.Logf("files=%v err=%v", files, err)
}

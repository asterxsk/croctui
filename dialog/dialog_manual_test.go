// dialog/dialog_manual_test.go  — skipped in CI; run with -run TestManualPick
package dialog

import "testing"

func TestManualPick(t *testing.T) {
	if testing.Short() {
		t.Skip("manual")
	}
	files, err := PickFiles()
	t.Logf("files=%v err=%v", files, err)
}

package dialog

import (
	"bufio"
	"os/exec"
	"strings"
)

// dialogScript is a PowerShell fragment using WinForms to show a native,
// multi-select file dialog. Output is UTF-8, one absolute path per line.
// Empty output means the user cancelled.
func dialogScript() string {
	return `[Console]::OutputEncoding = [System.Text.Encoding]::UTF8
Add-Type -AssemblyName System.Windows.Forms
$dialog = New-Object System.Windows.Forms.OpenFileDialog
$dialog.Multiselect = $true
$dialog.Title = 'Select files to send'
if ($dialog.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK) {
    $dialog.FileNames | ForEach-Object { $_.ToString() }
}`
}

// PickFiles opens a native Windows multi-select dialog and returns the chosen
// absolute paths. An empty slice means the user cancelled.
func PickFiles() ([]string, error) {
	// -STA is required for WinForms dialogs on Windows.
	cmd := exec.Command("powershell.exe", "-NoProfile", "-STA", "-Command", dialogScript())
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var files []string
	sc := bufio.NewScanner(strings.NewReader(string(out)))
	for sc.Scan() {
		if p := strings.TrimSpace(sc.Text()); p != "" {
			files = append(files, p)
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return files, nil
}

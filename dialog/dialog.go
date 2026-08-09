package dialog

import (
	"bufio"
	"os/exec"
	"strings"
)

// dialogScript is a PowerShell fragment using WinForms to show a native,
// multi-select file dialog. Output is UTF-8, one absolute path per line.
// Empty output means the user cancelled. The dialog opens at the process's
// current directory rather than the default fixed location.
func dialogScript() string {
	return `[Console]::OutputEncoding = [System.Text.Encoding]::UTF8
Add-Type -AssemblyName System.Windows.Forms
$dialog = New-Object System.Windows.Forms.OpenFileDialog
$dialog.Multiselect = $true
$dialog.Title = 'Select files to send'
$cwd = [System.IO.Directory]::GetCurrentDirectory()
if (Test-Path $cwd) { $dialog.InitialDirectory = $cwd }
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
	return parsePaths(string(out)), nil
}

// folderScript is a PowerShell fragment using WinForms to show a native
// folder picker. It returns a single absolute path, or nothing on cancel.
// The dialog opens at the process's current directory (inherited from the Go
// process's CWD) rather than the default Desktop/MyComputer root.
func folderScript() string {
	return `[Console]::OutputEncoding = [System.Text.Encoding]::UTF8
Add-Type -AssemblyName System.Windows.Forms
$dialog = New-Object System.Windows.Forms.FolderBrowserDialog
$dialog.Description = 'Select a folder to send'
$cwd = [System.IO.Directory]::GetCurrentDirectory()
if (Test-Path $cwd) { $dialog.SelectedPath = $cwd }
if ($dialog.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK) {
    $dialog.SelectedPath
}`
}

// PickFolder opens a native Windows folder dialog and returns the chosen
// absolute path. An empty slice means the user cancelled.
func PickFolder() ([]string, error) {
	cmd := exec.Command("powershell.exe", "-NoProfile", "-STA", "-Command", folderScript())
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return parsePaths(string(out)), nil
}

// parsePaths splits the dialog's newline-separated UTF-8 output into trimmed,
// non-empty paths.
func parsePaths(out string) []string {
	var files []string
	sc := bufio.NewScanner(strings.NewReader(out))
	for sc.Scan() {
		if p := strings.TrimSpace(sc.Text()); p != "" {
			files = append(files, p)
		}
	}
	return files
}

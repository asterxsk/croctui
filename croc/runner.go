package croc

import "os/exec"

// Detect returns the path to the croc executable, or an error if not installed.
func Detect() (string, error) {
	return exec.LookPath("croc")
}

// SendArgs builds the argv for a send: croc --yes send <files…>.
// croc generates and prints the code itself.
func SendArgs(files []string) []string {
	return append([]string{"--yes", "send"}, files...)
}

// ReceiveArgs builds the argv for a receive: croc --yes --rename <code> --out <dir>.
func ReceiveArgs(code, outDir string) []string {
	return []string{"--yes", "--rename", code, "--out", outDir}
}

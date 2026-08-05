package croc

import "os/exec"

// Install installs croc via winget (Windows). Output is discarded; the UI
// shows a spinner and re-runs Detect() on completion.
func Install() error {
	cmd := exec.Command("winget", "install", "schollz.croc",
		"--silent",
		"--accept-package-agreements",
		"--accept-source-agreements",
	)
	return cmd.Run()
}
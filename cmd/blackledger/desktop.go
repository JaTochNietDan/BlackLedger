package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Distribution launchers opt in; source checkouts retain their existing save path.
func desktopPaths(executable, configDir string) (web, save string, err error) {
	web = filepath.Join(filepath.Dir(executable), "dist")
	if info, e := os.Stat(filepath.Join(web, "index.html")); e != nil || info.IsDir() {
		return "", "", fmt.Errorf("game files missing beside executable: extract the complete archive before launching")
	}
	return web, filepath.Join(configDir, "BlackLedger", "campaign.sqlite3"), nil
}

func openBrowser(url string) error {
	var command *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		command = exec.Command("open", url)
	case "windows":
		command = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		command = exec.Command("xdg-open", url)
	}
	return command.Run()
}

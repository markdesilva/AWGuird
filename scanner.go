package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/gotk3/gotk3/gdk"
)

var amneziaKeys = []string{"jc", "jmin", "jmax", "s1", "s2", "h1", "h2", "h3", "h4"}

func (app *App) RefreshTunnels() {
	app.TunnelListStore.Clear()
	app.TunnelsMap = make(map[string]*Tunnel)
	
	app.scanDirectory(WireGuardPath)
	app.scanDirectory(AmneziaPath)
}

func (app *App) scanDirectory(dirPath string) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return 
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		ext := filepath.Ext(name)

		if ext == ".conf" || ext == ".awg" {
			tunnelName := strings.TrimSuffix(name, ext)
			fullPath := filepath.Join(dirPath, name)

			isAWG := (ext == ".awg") || isAmneziaConfig(fullPath)

			t := &Tunnel{
				Name:  tunnelName,
				Path:  fullPath,
				IsAWG: isAWG,
			}

			app.TunnelsMap[t.Name] = t

			// FIXED: Select embedded byte array instead of local file paths
			chosenBytes := ShieldInactiveBytes
			if app.ActiveTunnel == t.Name {
				chosenBytes = ShieldActiveBytes
			}

			var pixbuf *gdk.Pixbuf
			loader, err := gdk.PixbufLoaderNew()
			if err == nil {
				_, _ = loader.Write(chosenBytes)
				_ = loader.Close()
				if pb, err := loader.GetPixbuf(); err == nil {
					// Perform high-quality structural scaling entirely in-memory
					pixbuf, _ = pb.ScaleSimple(22, 22, gdk.INTERP_BILINEAR)
				}
			}

			iter := app.TunnelListStore.Append()
			_ = app.TunnelListStore.Set(iter, []int{0, 1}, []interface{}{pixbuf, t.Name})
		}
	}
}

func isAmneziaConfig(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		for _, key := range amneziaKeys {
			if strings.HasPrefix(line, key+" ") || strings.HasPrefix(line, key+"=") {
				return true
			}
		}
	}
	return false
}

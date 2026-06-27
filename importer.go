package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

func (app *App) ImportLooseConfigPicker() {
	dialog, _ := gtk.FileChooserDialogNewWith2Buttons(
		"Select Config File", app.Window, gtk.FILE_CHOOSER_ACTION_OPEN,
		"Cancel", gtk.RESPONSE_CANCEL, "Open", gtk.RESPONSE_ACCEPT,
	)
	defer dialog.Destroy()

	filter, _ := gtk.FileFilterNew()
	filter.SetName("Wireguard Configs (*.conf, *.awg)")
	filter.AddPattern("*.conf")
	filter.AddPattern("*.awg")
	dialog.AddFilter(filter)

	if dialog.Run() == gtk.RESPONSE_ACCEPT {
		if err := app.importLooseFile(dialog.GetFilename()); err != nil {
			app.UpdateLog(fmt.Sprintf("Import Error: %v", err))
		} else {
			app.UpdateLog("Config imported successfully.")
			app.RefreshTunnels()
		}
	}
}

func (app *App) ImportZipPicker() {
	dialog, _ := gtk.FileChooserDialogNewWith2Buttons(
		"Select Zip Archive", app.Window, gtk.FILE_CHOOSER_ACTION_OPEN,
		"Cancel", gtk.RESPONSE_CANCEL, "Open", gtk.RESPONSE_ACCEPT,
	)
	defer dialog.Destroy()

	filter, _ := gtk.FileFilterNew()
	filter.SetName("Zip Archives (*.zip)")
	filter.AddPattern("*.zip")
	dialog.AddFilter(filter)

	if dialog.Run() == gtk.RESPONSE_ACCEPT {
		if err := app.importZip(dialog.GetFilename()); err != nil {
			app.UpdateLog(fmt.Sprintf("Archive Extract Error: %v", err))
		} else {
			app.UpdateLog("Zip Archive items processed successfully.")
			app.RefreshTunnels()
		}
	}
}

func (app *App) DeleteSelectedTunnel() {
	if app.CurrentTunnel == nil {
		return
	}

	t := app.CurrentTunnel
	if app.ActiveTunnel == t.Name {
		app.HandleToggle()
	}

	confirmDialog := gtk.MessageDialogNew(
		app.Window,
		gtk.DIALOG_MODAL,
		gtk.MESSAGE_WARNING,
		gtk.BUTTONS_YES_NO,
		"Are you sure you want to permanently delete configuration: %s?", t.Name,
	)
	defer confirmDialog.Destroy()

	if confirmDialog.Run() == gtk.RESPONSE_YES {
		err := os.Remove(t.Path)
		if err != nil {
			app.UpdateLog(fmt.Sprintf("Failed to delete filesystem target: %v", err))
			return
		}
		app.CurrentTunnel = nil
		app.UpdateLog(fmt.Sprintf("Successfully removed tunnel profile: %s", t.Name))
		app.RefreshTunnels()
	}
}

func (app *App) importLooseFile(srcPath string) error {
	content, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}
	isAWG := strings.HasSuffix(strings.ToLower(srcPath), ".awg") || isAmneziaConfig(srcPath)
	var destDir string
	if isAWG {
		destDir = AmneziaPath
	} else {
		destDir = WireGuardPath
	}
	_ = os.MkdirAll(destDir, 0700)
	destPath := filepath.Join(destDir, strings.TrimSuffix(filepath.Base(srcPath), filepath.Ext(srcPath))+".conf")
	return os.WriteFile(destPath, content, 0600)
}

func (app *App) importZip(zipPath string) error {
	archive, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer archive.Close()
	for _, file := range archive.File {
		if file.FileInfo().IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(file.Name))
		if ext != ".conf" && ext != ".awg" {
			continue
		}
		rc, err := file.Open()
		if err != nil {
			return err
		}
		temp, _ := os.CreateTemp("", "import-*.conf")
		_, _ = io.Copy(temp, rc)
		tempPath := temp.Name()
		temp.Close()
		rc.Close()

		isAWG := (ext == ".awg") || isAmneziaConfig(tempPath)
		var destDir string
		if isAWG {
			destDir = AmneziaPath
		} else {
			destDir = WireGuardPath
		}
		_ = os.MkdirAll(destDir, 0700)
		destPath := filepath.Join(destDir, strings.TrimSuffix(filepath.Base(file.Name), ext)+".conf")
		content, _ := os.ReadFile(tempPath)
		_ = os.WriteFile(destPath, content, 0600)
		os.Remove(tempPath)
	}
	return nil
}

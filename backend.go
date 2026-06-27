package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

func getPublicKeyFromPrivate(privateKey string, isAWG bool) string {
	binary := "wg"
	if isAWG {
		binary = "awg"
	}
	cmd := exec.Command(binary, "pubkey")
	cmd.Stdin = strings.NewReader(privateKey)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err == nil {
		return strings.TrimSpace(out.String())
	}
	return "—"
}

func (app *App) SelectTunnel(t *Tunnel) {
	if t == nil {
		return
	}
	app.CurrentTunnel = t

	file, err := os.Open(t.Path)
	if err != nil {
		return
	}
	defer file.Close()

	var addresses, dnsServers []string
	privateKeyString := ""
	peerPubKey := "—"
	allowedIPs := "—"
	endpoint := "—"

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		key := strings.ToLower(strings.TrimSpace(parts[0]))
		val := strings.TrimSpace(parts[1])

		switch key {
		case "privatekey":
			privateKeyString = val
		case "publickey":
			peerPubKey = val
		case "address":
			addresses = append(addresses, val)
		case "dns":
			dnsServers = append(dnsServers, val)
		case "allowedips":
			allowedIPs = val
		case "endpoint":
			endpoint = val
		}
	}

	computedInterfacePubKey := "—"
	if privateKeyString != "" {
		computedInterfacePubKey = getPublicKeyFromPrivate(privateKeyString, t.IsAWG)
	}

	app.PubKeyLabel.SetText(fmt.Sprintf("Public key: %s", computedInterfacePubKey))
	app.AddressLabel.SetText(fmt.Sprintf("Addresses: %s", strings.Join(addresses, ", ")))
	app.DnsLabel.SetText(fmt.Sprintf("DNS servers: %s", strings.Join(dnsServers, ", ")))
	app.PeerPubKeyLabel.SetText(fmt.Sprintf("Public key: %s", peerPubKey))
	app.AllowedIpsLabel.SetText(fmt.Sprintf("Allowed IPs: %s", allowedIPs))
	app.EndpointLabel.SetText(fmt.Sprintf("Endpoint: %s", endpoint))

	app.UpdateStatusUI()
}

func (app *App) UpdateStatusUI() {
	if app.CurrentTunnel == nil {
		app.EditConfigBtn.SetSensitive(false)
		return
	}
	t := app.CurrentTunnel

	if app.ActiveTunnel == t.Name {
		app.StatusLabel.SetText("Status: Active")
		app.ToggleBtn.SetLabel("Deactivate")
		app.HeaderSubtitle.SetText(fmt.Sprintf("Connected to %s", t.Name))
		
		app.EditConfigBtn.SetSensitive(false)
		app.PollLiveStats()
	} else {
		app.StatusLabel.SetText("Status: Disconnected")
		app.ToggleBtn.SetLabel("Activate")
		app.PortLabel.SetText("Listen port: —")
		app.HandshakeLabel.SetText("Latest handshake: —")
		app.TransferLabel.SetText("Transfer: —")
		app.HeaderSubtitle.SetText("Disconnected")
		
		app.EditConfigBtn.SetSensitive(true)
	}
	app.ToggleBtn.SetSensitive(true)
	app.Window.QueueDraw()
}

func (app *App) HandleToggle() {
	if app.CurrentTunnel == nil {
		return
	}

	t := app.CurrentTunnel
	targetName := t.Name
	app.ToggleBtn.SetSensitive(false)
	app.EditConfigBtn.SetSensitive(false)

	if app.ActiveTunnel != "" && app.ActiveTunnel != targetName {
		oldTunnelName := app.ActiveTunnel
		oldTunnel, exists := app.TunnelsMap[oldTunnelName]
		
		if exists {
			app.StatusLabel.SetText(fmt.Sprintf("Disconnecting from %s first...", oldTunnelName))
			app.Window.QueueDraw()

			disconnectBinary := "wg-quick"
			if oldTunnel.IsAWG {
				disconnectBinary = "awg-quick"
			}

			downCmd := exec.Command(disconnectBinary, "down", oldTunnel.Path)
			_ = downCmd.Run()
		}
		app.ActiveTunnel = ""
	}

	commandBinary := "wg-quick"
	if t.IsAWG {
		commandBinary = "awg-quick"
	}

	action := "up"
	if app.ActiveTunnel == targetName {
		action = "down"
	}

	app.StatusLabel.SetText(fmt.Sprintf("Running %s %s...", commandBinary, action))
	cmd := exec.Command(commandBinary, action, t.Path)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()

	outputLog := outBuf.String() + "\n" + errBuf.String()
	app.UpdateLog(outputLog)

	if err != nil {
		app.StatusLabel.SetText("Status: Error executing action")
		app.ToggleBtn.SetSensitive(true)
		return
	}

	if action == "up" {
		app.ActiveTunnel = targetName
	} else {
		app.ActiveTunnel = ""
	}

	app.RefreshTunnels()

	if restoredTunnel, exists := app.TunnelsMap[targetName]; exists {
		app.CurrentTunnel = restoredTunnel
	} else {
		app.CurrentTunnel = t
	}

	app.SelectActiveTunnelInTreeView()
	
	// FIXED: Force the UI data labels to accurately render properties of the toggled tunnel instantly
	app.SelectTunnel(app.CurrentTunnel) 
}

func (app *App) ShowConfigEditorModal() {
	if app.CurrentTunnel == nil {
		return
	}
	t := app.CurrentTunnel

	if app.ActiveTunnel == t.Name {
		return
	}

	// FIXED: Direct instantiation and initialization fixes gotk3 variadic argument restrictions
	dialog, _ := gtk.DialogNew()
	dialog.SetTitle(fmt.Sprintf("Editing %s", t.Name))
	dialog.SetTransientFor(app.Window)
	dialog.SetModal(true)
	dialog.SetDestroyWithParent(true)
	dialog.AddButton("Cancel", gtk.RESPONSE_CANCEL)
	dialog.AddButton("Save", gtk.RESPONSE_ACCEPT)
	
	defer dialog.Destroy()
	dialog.SetDefaultSize(600, 500)

	contentArea, _ := dialog.GetContentArea()
	scroll, _ := gtk.ScrolledWindowNew(nil, nil)
	scroll.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	scroll.SetMarginStart(10)
	scroll.SetMarginEnd(10)
	scroll.SetMarginTop(10)
	scroll.SetMarginBottom(10)

	textView, _ := gtk.TextViewNew()
	textView.SetMonospace(true)
	
	bytesRead, err := os.ReadFile(t.Path)
	if err == nil {
		buf, _ := textView.GetBuffer()
		buf.SetText(string(bytesRead))
	} else {
		app.UpdateLog(fmt.Sprintf("File read error: %v", err))
		return
	}

	scroll.Add(textView)
	contentArea.PackStart(scroll, true, true, 0)
	contentArea.ShowAll()

	if dialog.Run() == gtk.RESPONSE_ACCEPT {
		buf, _ := textView.GetBuffer()
		start, end := buf.GetStartIter(), buf.GetEndIter()
		updatedText, _ := buf.GetText(start, end, false)

		err := os.WriteFile(t.Path, []byte(updatedText), 0600)
		if err != nil {
			app.UpdateLog(fmt.Sprintf("Failed to save changes: %v", err))
		} else {
			app.UpdateLog(fmt.Sprintf("Configuration updated successfully for %s", t.Name))
			app.SelectTunnel(t)
		}
	}
}

func (app *App) PollLiveStats() {
	if app.ActiveTunnel == "" || app.CurrentTunnel == nil || app.ActiveTunnel != app.CurrentTunnel.Name {
		return
	}

	binaryName := "wg"
	if app.CurrentTunnel.IsAWG {
		binaryName = "awg"
	}

	cmd := exec.Command(binaryName, "show", app.ActiveTunnel)
	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	if err := cmd.Run(); err != nil {
		return
	}

	scanner := bufio.NewScanner(&outBuf)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(parts[0]))
		val := strings.TrimSpace(parts[1])

		switch key {
		case "listening port":
			app.PortLabel.SetText(fmt.Sprintf("Listen port: %s", val))
		case "latest handshake":
			app.HandshakeLabel.SetText(fmt.Sprintf("Latest handshake: %s", val))
		case "transfer":
			app.TransferLabel.SetText(fmt.Sprintf("Transfer: %s", val))
		}
	}
	app.Window.QueueDraw()
}

func (app *App) UpdateLog(message string) {
	buffer, _ := app.LogTextView.GetBuffer()
	endIter := buffer.GetEndIter()
	buffer.Insert(endIter, message+"\n")
	
	endIter = buffer.GetEndIter()
	mark := buffer.CreateMark("", endIter, false)
	app.LogTextView.ScrollToMark(mark, 0.0, false, 0.0, 1.0)
}

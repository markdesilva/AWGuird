package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	WireGuardPath = "/etc/wireguard"
	AmneziaPath   = "/etc/amnezia/amneziawg"
)

type Tunnel struct {
	Name  string
	Path  string
	IsAWG bool
}

type App struct {
	Window          *gtk.Window
	TunnelListStore *gtk.ListStore
	TunnelTreeView  *gtk.TreeView
	
	HeaderSubtitle  *gtk.Label

	StatusLabel     *gtk.Label
	PubKeyLabel     *gtk.Label
	PortLabel       *gtk.Label
	AddressLabel    *gtk.Label
	DnsLabel        *gtk.Label
	ToggleBtn       *gtk.Button
	EditConfigBtn   *gtk.Button

	PeerPubKeyLabel *gtk.Label
	AllowedIpsLabel *gtk.Label
	EndpointLabel   *gtk.Label
	HandshakeLabel  *gtk.Label
	TransferLabel   *gtk.Label

	LogTextView     *gtk.TextView
	
	// FIXED: Using generic glib.Object to safely interface with the system's tray wrapper
	StatusIcon      *glib.Object 
	
	CurrentTunnel   *Tunnel
	ActiveTunnel    string
	TunnelsMap      map[string]*Tunnel
}

func main() {
	if os.Geteuid() != 0 {
		fmt.Println("Error: This client must be run as root to manage network interfaces and pull statistics.")
		os.Exit(1)
	}

	gtk.Init(nil)

	app := &App{
		TunnelsMap: make(map[string]*Tunnel),
	}
	app.BuildUI()
	app.RefreshTunnels()

	// Global loop driving background statistics synchronization
	go func() {
		for {
			time.Sleep(2 * time.Second)
			glib.IdleAdd(func() bool {
				if app.ActiveTunnel != "" {
					app.PollLiveStats()
				}
				return false 
			})
		}
	}()

	gtk.Main()
}

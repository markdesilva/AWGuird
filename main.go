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

	// Check if the --minimized flag was passed via autostart config
	minimizedStartup := false
	for _, arg := range os.Args {
		if arg == "--minimized" {
			minimizedStartup = true
			break
		}
	}

	// Race condition protection: allows window rendering components on panels to initialize
	if minimizedStartup {
		time.Sleep(4 * time.Second)
	}

	gtk.Init(nil)

	app := &App{
		TunnelsMap: make(map[string]*Tunnel),
	}
	app.BuildUI()
	app.RefreshTunnels()

	// Adjust application main window map visibility based on start flags
	if minimizedStartup {
		app.Window.Hide()
	} else {
		app.Window.ShowAll()
		app.Window.Present()
	}

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

package main

/*
#cgo pkg-config: gtk+-3.0
#include <gtk/gtk.h>

// Forward declarations for CGO signal forwarding
extern void go_tray_activate_callback(GObject *object, gpointer user_data);
extern void go_tray_popup_callback(GObject *object, guint button, guint activate_time, gpointer user_data);

static inline GObject* create_status_icon(GdkPixbuf *pixbuf, gpointer user_data) {
    GtkStatusIcon *icon = gtk_status_icon_new_from_pixbuf(pixbuf);
    gtk_status_icon_set_title(icon, "AWGuird");
    gtk_status_icon_set_tooltip_text(icon, "AWGuird VPN Status");
    gtk_status_icon_set_visible(icon, TRUE);
    
    g_signal_connect(G_OBJECT(icon), "activate", G_CALLBACK(go_tray_activate_callback), user_data);
    g_signal_connect(G_OBJECT(icon), "popup-menu", G_CALLBACK(go_tray_popup_callback), user_data);
    
    return G_OBJECT(icon);
}

// Positions the context menu aligned right under the panel status icon geometry bounds
static inline void popup_status_menu(GObject *icon, GObject *menu, guint button, guint activate_time) {
    gtk_menu_popup(GTK_MENU(menu), NULL, NULL, gtk_status_icon_position_menu, GTK_STATUS_ICON(icon), button, activate_time);
}
*/
import "C"

import (
	"path/filepath"
	"unsafe"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var globalWin *gtk.Window
var globalTrayMenu *gtk.Menu

//export go_tray_activate_callback
func go_tray_activate_callback(object *C.GObject, user_data C.gpointer) {
	glib.IdleAdd(func() bool {
		if globalWin != nil {
			if globalWin.GetVisible() {
				globalWin.Hide()
			} else {
				globalWin.ShowAll()
				globalWin.Present()
			}
		}
		return false
	})
}

//export go_tray_popup_callback
func go_tray_popup_callback(object *C.GObject, button C.guint, activate_time C.guint, user_data C.gpointer) {
	glib.IdleAdd(func() bool {
		if globalTrayMenu != nil {
			cMenu := (*C.GObject)(unsafe.Pointer(globalTrayMenu.Native()))
			C.popup_status_menu(object, cMenu, button, activate_time)
		}
		return false
	})
}

func (app *App) BuildUI() {
	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	win.SetTitle("AWGuird")
	win.SetDefaultSize(880, 580)
	app.Window = win
	globalWin = win

	loader, err := gdk.PixbufLoaderNew()
	var mainPixbuf *gdk.Pixbuf
	if err == nil {
		_, _ = loader.Write(AppIconBytes)
		_ = loader.Close()
		if pixbuf, err := loader.GetPixbuf(); err == nil {
			win.SetIcon(pixbuf)
			mainPixbuf = pixbuf
		}
	}

	mainVBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)

	headerCenterBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 2)
	headerCenterBox.SetMarginTop(10)
	headerCenterBox.SetMarginBottom(10)
	
	app.HeaderSubtitle, _ = gtk.LabelNew("Disconnected")
	headerCenterBox.PackStart(app.HeaderSubtitle, false, false, 0)
	mainVBox.PackStart(headerCenterBox, false, false, 0)

	notebook, _ := gtk.NotebookNew()
	notebook.SetScrollable(true)

	tunnelsTabBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 12)
	tunnelsTabBox.SetMarginStart(10)
	tunnelsTabBox.SetMarginEnd(10)
	tunnelsTabBox.SetMarginTop(10)
	tunnelsTabBox.SetMarginBottom(10)

	leftColumnBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 8)
	scrollWin, _ := gtk.ScrolledWindowNew(nil, nil)
	scrollWin.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	scrollWin.SetSizeRequest(240, -1)

	app.TunnelListStore, _ = gtk.ListStoreNew(glib.TYPE_OBJECT, glib.TYPE_STRING)
	app.TunnelTreeView, _ = gtk.TreeViewNewWithModel(app.TunnelListStore)
	
	iconRenderer, _ := gtk.CellRendererPixbufNew()
	iconCol, _ := gtk.TreeViewColumnNew()
	iconCol.PackStart(iconRenderer, false)
	iconCol.AddAttribute(iconRenderer, "pixbuf", 0)
	
	textRenderer, _ := gtk.CellRendererTextNew()
	nameCol, _ := gtk.TreeViewColumnNewWithAttribute("Tunnels", textRenderer, "text", 1)
	
	app.TunnelTreeView.AppendColumn(iconCol)
	app.TunnelTreeView.AppendColumn(nameCol)
	scrollWin.Add(app.TunnelTreeView)
	leftColumnBox.PackStart(scrollWin, true, true, 0)

	bottomButtonGrid, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
	addFileBtn, _ := gtk.ButtonNewWithLabel("➕ Add Tunnel")
	deleteBtn, _ := gtk.ButtonNewWithLabel("❌")
	importZipBtn, _ := gtk.ButtonNewWithLabel("🗄️")
	
	bottomButtonGrid.PackStart(addFileBtn, true, true, 0)
	bottomButtonGrid.PackStart(deleteBtn, false, false, 0)
	bottomButtonGrid.PackStart(importZipBtn, false, false, 0)
	leftColumnBox.PackStart(bottomButtonGrid, false, false, 0)

	rightColumnBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)

	interfaceFrame, _ := gtk.FrameNew("Interface")
	interfaceGrid, _ := gtk.GridNew()
	interfaceGrid.SetOrientation(gtk.ORIENTATION_VERTICAL)
	interfaceGrid.SetRowSpacing(6)
	interfaceGrid.SetMarginStart(15)
	interfaceGrid.SetMarginEnd(15)
	interfaceGrid.SetMarginTop(10)
	interfaceGrid.SetMarginBottom(10)

	app.StatusLabel, _ = gtk.LabelNew("Status: Disconnected")
	app.PubKeyLabel, _ = gtk.LabelNew("Public key: —")
	app.PortLabel, _ = gtk.LabelNew("Listen port: —")
	app.AddressLabel, _ = gtk.LabelNew("Addresses: —")
	app.DnsLabel, _ = gtk.LabelNew("DNS servers: —")
	app.ToggleBtn, _ = gtk.ButtonNewWithLabel("Activate")
	app.ToggleBtn.SetSensitive(false)

	app.StatusLabel.SetXAlign(0)
	app.PubKeyLabel.SetXAlign(0)
	app.PubKeyLabel.SetSelectable(true)
	app.PortLabel.SetXAlign(0)
	app.AddressLabel.SetXAlign(0)
	app.DnsLabel.SetXAlign(0)

	interfaceGrid.Add(app.StatusLabel)
	interfaceGrid.Add(app.PubKeyLabel)
	interfaceGrid.Add(app.PortLabel)
	interfaceGrid.Add(app.AddressLabel)
	interfaceGrid.Add(app.DnsLabel)
	
	btnBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	btnBox.SetMarginTop(6)
	btnBox.PackStart(app.ToggleBtn, false, false, 0)
	interfaceGrid.Add(btnBox)
	interfaceFrame.Add(interfaceGrid)
	rightColumnBox.PackStart(interfaceFrame, true, true, 0)

	peerFrame, _ := gtk.FrameNew("Peer")
	peerGrid, _ := gtk.GridNew()
	peerGrid.SetOrientation(gtk.ORIENTATION_VERTICAL)
	peerGrid.SetRowSpacing(6)
	peerGrid.SetMarginStart(15)
	peerGrid.SetMarginEnd(15)
	peerGrid.SetMarginTop(10)
	peerGrid.SetMarginBottom(10)

	app.PeerPubKeyLabel, _ = gtk.LabelNew("Public key: —")
	app.AllowedIpsLabel, _ = gtk.LabelNew("Allowed IPs: —")
	app.EndpointLabel, _ = gtk.LabelNew("Endpoint: —")
	app.HandshakeLabel, _ = gtk.LabelNew("Latest handshake: —")
	app.TransferLabel, _ = gtk.LabelNew("Transfer: —")

	app.PeerPubKeyLabel.SetXAlign(0)
	app.PeerPubKeyLabel.SetSelectable(true)
	app.AllowedIpsLabel.SetXAlign(0)
	app.EndpointLabel.SetXAlign(0)
	app.HandshakeLabel.SetXAlign(0)
	app.TransferLabel.SetXAlign(0)

	peerGrid.Add(app.PeerPubKeyLabel)
	peerGrid.Add(app.AllowedIpsLabel)
	peerGrid.Add(app.EndpointLabel)
	peerGrid.Add(app.HandshakeLabel)
	peerGrid.Add(app.TransferLabel)
	peerFrame.Add(peerGrid)
	rightColumnBox.PackStart(peerFrame, true, true, 0)

	actionActionBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	app.EditConfigBtn, _ = gtk.ButtonNewWithLabel("📝 Edit")
	app.EditConfigBtn.SetSensitive(false)
	actionActionBox.PackEnd(app.EditConfigBtn, false, false, 0)
	rightColumnBox.PackStart(actionActionBox, false, false, 0)

	tunnelsTabBox.PackStart(leftColumnBox, false, false, 0)
	tunnelsTabBox.PackStart(rightColumnBox, true, true, 0)

	logsTabBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	logsTabBox.SetMarginStart(10)
	logsTabBox.SetMarginEnd(10)
	logsTabBox.SetMarginTop(10)
	logsTabBox.SetMarginBottom(10)

	app.LogTextView, _ = gtk.TextViewNew()
	app.LogTextView.SetEditable(false)
	logScroll, _ := gtk.ScrolledWindowNew(nil, nil)
	logScroll.Add(app.LogTextView)
	logsTabBox.PackStart(logScroll, true, true, 0)

	tunnelsTabLabel, _ := gtk.LabelNew("  Tunnels  ")
	logsTabLabel, _ := gtk.LabelNew("  Logs  ")
	notebook.AppendPage(tunnelsTabBox, tunnelsTabLabel)
	notebook.AppendPage(logsTabBox, logsTabLabel)
	
	mainVBox.PackStart(notebook, true, true, 0)
	win.Add(mainVBox)

	notebook.Connect("switch-page", func(nb *gtk.Notebook, page *gtk.Widget, pageNum uint) {
		if pageNum == 0 && app.ActiveTunnel != "" {
			app.SelectActiveTunnelInTreeView()
		}
	})

	app.TunnelTreeView.Connect("cursor-changed", func() {
		selection, _ := app.TunnelTreeView.GetSelection()
		model, iter, ok := selection.GetSelected()
		if ok {
			val, _ := model.(*gtk.TreeModel).GetValue(iter, 1)
			name, _ := val.GetString()
			if tunnel, exists := app.TunnelsMap[name]; exists {
				app.SelectTunnel(tunnel)
			}
		}
	})

	app.TunnelTreeView.Connect("row-activated", func() {
		selection, _ := app.TunnelTreeView.GetSelection()
		model, iter, ok := selection.GetSelected()
		if ok {
			val, _ := model.(*gtk.TreeModel).GetValue(iter, 1)
			name, _ := val.GetString()
			if tunnel, exists := app.TunnelsMap[name]; exists {
				app.CurrentTunnel = tunnel
				app.HandleToggle()
			}
		}
	})

	app.ToggleBtn.Connect("clicked", func() {
		app.HandleToggle()
	})

	app.EditConfigBtn.Connect("clicked", func() {
		app.ShowConfigEditorModal()
	})

	addFileBtn.Connect("clicked", func() {
		app.ImportLooseConfigPicker()
	})

	importZipBtn.Connect("clicked", func() {
		app.ImportZipPicker()
	})

	deleteBtn.Connect("clicked", func() {
		app.DeleteSelectedTunnel()
	})

	if provider, err := gtk.CssProviderNew(); err == nil {
		provider.LoadFromData(`
			notebook { background-color: #f6f6f6; border: 1px solid #bcbcbc; }
			notebook tab { padding: 8px 24px; font-weight: bold; background-color: #e8e8e8; }
			notebook tab:active, notebook tab:checked { background-color: #ffffff; border-bottom: 2px solid #d35400; }
			frame { border: 1px solid #d0d0d0; border-radius: 4px; background-color: #ffffff; }
		`)
		if screen, err := gdk.ScreenGetDefault(); err == nil {
			gtk.AddProviderForScreen(screen, provider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
		}
	}

	if mainPixbuf != nil {
		trayMenu, _ := gtk.MenuNew()
		globalTrayMenu = trayMenu
		
		showItem, _ := gtk.MenuItemNewWithLabel("Open AWGuird")
		showItem.Connect("activate", func() {
			win.ShowAll()
			win.Present()
		})
		
		quitItem, _ := gtk.MenuItemNewWithLabel("Quit Application")
		quitItem.Connect("activate", func() {
			gtk.MainQuit()
		})
		
		trayMenu.Append(showItem)
		trayMenu.Append(quitItem)
		trayMenu.ShowAll()

		cPixbuf := (*C.GdkPixbuf)(unsafe.Pointer(mainPixbuf.Native()))
		_ = C.create_status_icon(cPixbuf, nil)
	}

	win.Connect("delete-event", func() bool {
		win.Hide()
		return true 
	})

	// REMOVED win.ShowAll() here! Visibility is now controlled by main.go
}

func (app *App) SelectActiveTunnelInTreeView() {
	if app.ActiveTunnel == "" {
		return
	}
	model := app.TunnelListStore
	iter, ok := model.GetIterFirst()
	for ok {
		val, _ := model.GetValue(iter, 1)
		name, _ := val.GetString()
		if name == app.ActiveTunnel {
			selection, _ := app.TunnelTreeView.GetSelection()
			selection.SelectIter(iter)
			
			path, _ := model.GetPath(iter)
			app.TunnelTreeView.ScrollToCell(path, nil, false, 0.0, 0.0)
			break
		}
		ok = model.IterNext(iter)
	}
}

func (app *App) getIconPath(tunnelName string) string {
	if app.ActiveTunnel == tunnelName {
		return filepath.Join(".", "shield_active.png")
	}
	return filepath.Join(".", "shield_inactive.png")
}

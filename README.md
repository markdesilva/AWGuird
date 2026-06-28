# AWGuird
A Linux GTK gui client for AmneziaWG heavily inspired by UnnoTed's Wireguird

<img width="900" height="635" alt="AWGuird01" src="https://github.com/user-attachments/assets/396dd032-5c05-45b1-9956-c3287203ea02" />

### Features
+ Looks the same and does almost the same things as the official Wireguard's Windows gui client but for AmneziaWG
+ Lists tunnels from /etc/amnezia/amneziawg
+ Controls AmneziaWG through awg-quick
+ Runs on startup (will prompt for root password)
+ System tray icon to maximize/minimize window (doesn't affect tunnel)

### Prerequisites
+ AmneziaWG

### Compiling
#### Dependencies
+ libgtk-3-dev 
+ build-essential 
+ golang-go

```
sudo apt install libgtk-3-dev build-essential golang-go
```

#### Build
+ Do the following:
  
```
sudo git clone https://github.com/markdesilva/AWGuird.git
cd AWGuird
sudo go mod tidy
sudo go build -o awg-client .
```

### Installing
+ Create **/usr/share/applications/awg-client.desktop** and add the following to it:

```
[Desktop Entry]
Type=Application
Name=AWGuird
Comment=Manage WireGuard and AmneziaWG Tunnels
Exec=pkexec /usr/local/bin/awg-client
Icon=awg-client
Terminal=false
Categories=Network;Utility;
StartupWMClass=awg-client
```

+ Create **/etc/xdg/autostart/awg-client.desktop** and add the following to it:

```
[Desktop Entry]
Type=Application
Name=AWGuird AutoRun
Comment=Start AWGuird VPN interface monitor on session login
Exec=pkexec /usr/local/bin/awg-client
Terminal=false
NoDisplay=true
X-GNOME-Autostart-enabled=true
```

+ Copy the awg-client binary to /usr/local/bin
+ Copy app_icon.png file to /usr/share/pixmaps/awg-client.png
+ Restart gdm if the application doesn't show in the app drawer
  
```
sudo cp awg-client /usr/local/bin
sudo cp app_icon.png /usr/share/pixmaps/awg-client.png
sudo systemctl restart gdm
```

### Packages
Debian package is available.

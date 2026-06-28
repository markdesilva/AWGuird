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

### Installing (for Debian)
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
+ Debian package is available.
+ You can convert the Debian package to an Arch tarball using **debtab**
   + edit both **/usr/share/applications/awg-client.desktop** and **/etc/xdg/autostart/awg-client.desktop** to change _/usr/local/bin_ to _/usr/bin_ since debtap will move the awg-client to /usr/bin
   + make sure that your AmneziaWG config directory is **/etc/amnezia/amneziawg**
   + both showing up in the app tray and the auto start _should_ work if you're running GNOME or KDE Plasma
   + if you are running other windows managers (i3, Sway, Hyprland, etc) make sure the polkit agent (like polkit-gnome or lxsession) is running in the background
+ RPM package is available but note:
   + you need to install **gnome-extensions-app** and **gnome-shell-extension-appindicator** since Gnome for Fedora does away with system tray icons
   + after installing the above, you need to run **Extensions** and enable **AppIndicator and KStatusNotifierItem Support**
   + log out and log back in
   + install the RPM package with dnf, you should be able to see it in the app drawer
   + autostart on boot does **_NOT_** work natively, you need to install **gnome-tweaks** and add AWGuird as a startup application on log in
  


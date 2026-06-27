# AWGuird
A Linux GTK gui client for AmneziaWG heavily inspired by UnnoTed's Wireguird

<img width="900" height="635" alt="AWGuird01" src="https://github.com/user-attachments/assets/396dd032-5c05-45b1-9956-c3287203ea02" />

### Features
+ Looks the same and does almost the same things as the official Wireguard's Windows gui client but for AmneziaWG
+ Lists tunnels from /etc/amnezia/amneziawg
+ Controls AmneziaWG through awg-quick

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
+ Clone the repository
+ CD to the repository directory and do the following:
  
```
sudo go mod tidy
sudo go build -o awg-client .
```

### Installing
+ Create usr/share/applications/awg-client.desktop and add the following to it:

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

+ Copy the awg-client binary to /usr/local/bin
+ Copy app_icon.png file to /usr/share/pixmaps/awg-client.png
+ Restart gdm if the application doesn't show in the app drawer
  
### Packages
Debian package is available.

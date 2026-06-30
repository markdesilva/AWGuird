# AWGuird
A Linux GTK gui client for AmneziaWG heavily inspired by UnnoTed's Wireguird

<img width="900" height="635" alt="AWGuird01" src="https://github.com/user-attachments/assets/396dd032-5c05-45b1-9956-c3287203ea02" />

### Features
+ Looks the same and does almost the same things as the official Wireguard's Windows gui client but for AmneziaWG
+ Lists tunnels from /etc/amnezia/amneziawg (and /etc/wireguard if present)
+ Controls AmneziaWG through awg-quick
+ Runs on startup
+ System tray icon to maximize/minimize window (doesn't affect state of the tunnel)

### Prerequisites
+ AmneziaWG (Not covered here - AmneziaWG must be FULLY working)
+ FULLY working means, even before installing AWGuird, you must be able to connect to your VPN (using awg-quick) and have internet

### Compiling
#### Dependencies
+ gtk3 
+ go
+ libayatana-appindicator
+ polkit (if not already installed by default)

**Debian**
```
sudo apt install libgtk-3-dev build-essential golang-go libayatana-appindicator3-dev
```

**RPM**
```
sudo dnf install libgtk-3-dev go libayatana-appindicator3-dev
```

**Arch**
```
sudo pacman -S base-devel go gtk3 pkgconf polkit libayatana-appindicator
```

#### Build
+ Do the following:
  
```
sudo git clone https://github.com/markdesilva/AWGuird.git
cd AWGuird
sudo go mod tidy
sudo go get github.com/getlantern/systray
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
Exec=pkexec /usr/local/bin/awg-client --minimized
Terminal=false
NoDisplay=true
X-GNOME-Autostart-enabled=true
```

+ Create **/usr/share/polkit-1/actions/com.awguird.client.policy** and add the following to it:
```
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE policyconfig PUBLIC "-//freedesktop//DTD PolicyKit Policy Configuration 1.0//EN"
 "http://www.freedesktop.org/standards/PolicyKit/1/policyconfig.dtd">
<policyconfig>
  <action id="com.awguird.client.run">
    <description>Run AWGuird Client as root</description>
    <message>Authentication is required to manage secure VPN network configurations.</message>
    <defaults>
      <allow_any>yes</allow_any>
      <allow_inactive>yes</allow_inactive>
      <allow_active>yes</allow_active>
    </defaults>
    <annotate key="org.freedesktop.policykit.exec.path">/usr/local/bin/awg-client</annotate>
    <annotate key="org.freedesktop.policykit.exec.allow_gui">true</annotate>
  </action>
</policyconfig>
```

+ Create **/usr/share/polkit-1/rules.d/50-awg-client.rules** and add the following to it (only for Arch) :
```
polkit.addRule(function(action, subject) {
    if (action.id == "com.awguird.client.run") {
        return polkit.Result.YES;
    }
});
```

+ Copy the awg-client binary to /usr/local/bin for Debian or /usr/bin for Fedora
+ Copy app_icon.png file to /usr/share/pixmaps/awg-client.png
+ Set ownership and permissions
  
```
sudo cp awg-client /usr/local/bin
sudo cp app_icon.png /usr/share/pixmaps/awg-client.png
chmod 644 /usr/share/polkit-1/actions/com.awguird.client.policy /etc/xdg/autostart/awg-client.desktop /usr/share/applications/awg-client.desktop
chown root:root /usr/share/polkit-1/actions/com.awguird.client.policy /etc/xdg/autostart/awg-client.desktop /usr/share/applications/awg-client.desktop
chmod 755 /usr/local/bin/awg-client (Debian) or chmod 755 /usr/bin/awg-client (RPM/Arch)
```

+ Restart gdm if the application doesn't show in the app drawer
  
```
sudo systemctl restart gdm
```

### Packages
+ Debian package is available.
   + install the deb package with apt, you should be able to see it in the app drawer
   + log out and back in, it should autostart and show up in system tray minimized
     
+ Arch package is available
   + if you are running other windows managers (i3, Sway, Hyprland, etc) make sure the polkit is installed and running in the background
     
+ RPM package is available but note:
   + you need to install **gnome-extensions-app** and **gnome-shell-extension-appindicator** since Gnome for Fedora does away with system tray icons by default now
   + after installing the above, you need to run **Extensions** and enable **AppIndicator and KStatusNotifierItem Support**
   + log out and log back in
   + install the RPM package with dnf, you should be able to see it in the app drawer
   + log out and back in, it should autostart and show up in system tray minimized
  


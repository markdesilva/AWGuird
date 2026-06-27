package main

import _ "embed"

//go:embed app_icon.png
var AppIconBytes []byte

//go:embed shield_active.png
var ShieldActiveBytes []byte

//go:embed shield_inactive.png
var ShieldInactiveBytes []byte

module github.com/pltanton/autobrowser/linux

go 1.22.1

require github.com/labi-le/hyprland-ipc-client v1.0.3

require github.com/pltanton/autobrowser/common v0.0.0

require github.com/godbus/dbus/v5 v5.1.0 // indirect

replace github.com/pltanton/autobrowser/common => ../common

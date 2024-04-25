module github.com/pltanton/autobrowser/linux

go 1.22.1

require github.com/labi-le/hyprland-ipc-client v1.0.3

require github.com/pltanton/autobrowser/common v0.0.0

require (
	github.com/godbus/dbus/v5 v5.1.0
	github.com/joshuarubin/go-sway v1.2.0
)

require (
	github.com/joshuarubin/lifecycle v1.0.0 // indirect
	go.uber.org/atomic v1.3.2 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	golang.org/x/sync v0.0.0-20190412183630-56d357773e84 // indirect
)

replace github.com/pltanton/autobrowser/common => ../common

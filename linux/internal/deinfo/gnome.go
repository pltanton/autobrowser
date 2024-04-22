package deinfo

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/godbus/dbus/v5"
)

type gnomeProvider struct {
}

// fetchActiveApp implements deInfoProvider.
func (g *gnomeProvider) fetchActiveApp() (App, error) {
	slog.Debug("Fetch active app from gnome")
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return App{}, fmt.Errorf("failed to connect session bus: %w", err)
	}
	defer conn.Close()

	var response string
	obj := conn.Object("org.gnome.Shell", "/org/gnome/shell/extensions/FocusedWindow")
	if err := obj.Call("org.gnome.shell.extensions.FocusedWindow.Get", 0).Store(&response); err != nil {
		slog.Error("Failed to call dbus method. Did you install Focused Window D-Bus gnome extension?", "err", err)
		return App{}, nil
	}

	resultJson := struct {
		Title string `json:"title"`
		Class string `json:"wm_class"`
	}{}

	if err := json.Unmarshal([]byte(response), &resultJson); err != nil {
		slog.Error("Failed to unmarshal dbus response", "err", err, "response", response)
		return App{}, nil
	}

	return App{
		Class: resultJson.Class,
		Title: resultJson.Title,
	}, nil
}

var _ deInfoProvider = &gnomeProvider{}

func newGnomeProvider() deInfoProvider {
	return &gnomeProvider{}
}

package deinfo

import (
	"fmt"

	"github.com/godbus/dbus/v5"
)

type gnomeProvider struct {
}

// fetchActiveApp implements deInfoProvider.
func (g *gnomeProvider) fetchActiveApp() (App, error) {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return App{}, fmt.Errorf("failed to connect session bus: %w", err)
	}
	defer conn.Close()
	panic("unimplemented")
}

var _ deInfoProvider = &gnomeProvider{}

func newGnomeProvider() deInfoProvider {
	return &gnomeProvider{}
}

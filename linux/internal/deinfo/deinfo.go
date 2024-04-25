package deinfo

import (
	"log/slog"

	"github.com/pltanton/autobrowser/linux/internal/envx"
)

type App struct {
	Class string
	Title string
}

type DeInfoProvider struct {
	provider deInfoProvider

	activeAppSet bool
	activeApp    App
}

type deInfoProvider interface {
	fetchActiveApp() (App, error)
}

type noopProvider struct{}

func (noopProvider) fetchActiveApp() (App, error) {
	return App{}, nil
}

var _ deInfoProvider = noopProvider{}

func New(appMode envx.AppMode) *DeInfoProvider {
	var provider deInfoProvider

	switch appMode {
	case envx.HYPRLAND:
		provider = newHyprlandProvider()
	case envx.SWAY:
		provider = newSwayProvider()
	case envx.GNOME:
		provider = newGnomeProvider()
	case envx.UNKNOWN:
		provider = noopProvider{}
	}

	return &DeInfoProvider{
		provider: provider,
	}
}

func (p *DeInfoProvider) GetActiveApp() App {
	if !p.activeAppSet {
		var err error
		p.activeApp, err = p.provider.fetchActiveApp()
		p.activeAppSet = true

		if err != nil {
			slog.Error("Failed to set active active app", "err", err)
		}
	}

	return p.activeApp
}

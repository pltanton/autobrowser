package deinfo

import (
	"fmt"
	"log/slog"
	"os"

	ipc "github.com/labi-le/hyprland-ipc-client/v3"
)

type hyprlandProvider struct {
	c ipc.IPC
}

func newHyprlandProvider() deInfoProvider {
	return &hyprlandProvider{
		c: ipc.MustClient(os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")),
	}
}

func (h *hyprlandProvider) fetchActiveApp() (App, error) {
	slog.Debug("Fetch active app from hyprland")

	window, err := h.c.ActiveWindow()
	if err != nil {
		return App{}, fmt.Errorf("failed to fetch active window from hyprland: %w", err)
	}

	return App{
		Title: window.Title,
		Class: window.Class,
	}, nil
}

var _ deInfoProvider = &hyprlandProvider{}

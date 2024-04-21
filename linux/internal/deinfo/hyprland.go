package deinfo

import (
	"fmt"
	"os"

	ipc "github.com/labi-le/hyprland-ipc-client"
)

type hyprlandProvider struct {
	c ipc.IPC
}

func newHyprlandProvider() deInfoProvider {
	return &hyprlandProvider{
		c: ipc.NewClient(os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")),
	}
}

func (h *hyprlandProvider) fetchActiveApp() (App, error) {

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

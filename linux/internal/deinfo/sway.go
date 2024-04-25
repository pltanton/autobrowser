package deinfo

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	sway "github.com/joshuarubin/go-sway"
)

type swayProvider struct{}

// fetchActiveApp implements deInfoProvider.
func (s *swayProvider) fetchActiveApp() (App, error) {
	slog.Debug("Fetch active app from sway")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	client, err := sway.New(ctx)
	if err != nil {
		return App{}, fmt.Errorf("failed to create new sway client: %w", err)
	}

	node, err := client.GetTree(context.Background())
	if err != nil {
		return App{}, fmt.Errorf("failed to get sway tree: %w", err)
	}

	focusedNode := node.FocusedNode()

	var class string
	var title string = focusedNode.Name

	if focusedNode.WindowProperties != nil {
		// For xwayland clients
		class = focusedNode.WindowProperties.Class
		title = focusedNode.WindowProperties.Title
	} else if focusedNode.AppID != nil {
		class = *focusedNode.AppID
	}

	return App{
		Title: title,
		Class: class,
	}, nil
}

func newSwayProvider() deInfoProvider {
	return &swayProvider{}
}

var _ deInfoProvider = &swayProvider{}

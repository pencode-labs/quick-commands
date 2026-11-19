package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"os"
	"strings"

	"github.com/pencode-labs/pencode/peb"
)

type command struct {
	Name    string `json:"name"`
	Command string `json:"command"`
}

func main() {
	if err := peb.ServeExtension(os.Stdin, os.Stdout, peb.Contributions{
		Settings: []peb.Setting{{
			ID:          "commands",
			Label:       "Commands",
			Description: `JSON array with "name" and "command" fields, for example [{"name":"Run tests","command":"go test ./..."}].`,
			Default:     json.RawMessage(`"[{\"name\":\"Run tests\",\"command\":\"go test ./...\"}]"`),
		}},
		Actions: []peb.Action{{
			ID:       "quick-commands",
			Tooltip:  "Quick Commands",
			Icon:     "resources/terminal.svg",
			Location: peb.TitleBarRight,
		}},
	}, peb.Handlers{
		Settings: configure,
		ActivateActions: func(_ context.Context, p *peb.ActionPublisher) {
			publisher = p
		},
		Action:  openMenu,
		UIEvent: handleEvent,
	}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var (
	configured []command
	publisher  *peb.ActionPublisher
)

func configure(values peb.SettingValues) error {
	value, err := peb.DecodeSetting[string](values, "commands")
	if err != nil {
		return err
	}
	var commands []command
	if err := json.Unmarshal([]byte(value), &commands); err != nil {
		return fmt.Errorf("Quick Commands configuration is invalid: %w", err)
	}
	for _, item := range commands {
		if strings.TrimSpace(item.Name) == "" || strings.TrimSpace(item.Command) == "" {
			return fmt.Errorf("Quick Command names and commands cannot be empty")
		}
	}
	configured = commands
	return nil
}

func openMenu(c peb.ActionContext) error {
	var body strings.Builder
	body.WriteString(`<div class="quick-commands-menu" role="menu">`)
	for i, item := range configured {
		fmt.Fprintf(&body, `<button data-ref="command-%d" data-pencode-action="run-terminal" data-pencode-terminal-command="%s" type="button" role="menuitem">%s</button>`, i, html.EscapeString(item.Command), html.EscapeString(item.Name))
	}
	body.WriteString(`<button data-ref="settings" data-pencode-action="open-extension-settings" type="button" role="menuitem">Configure Quick Commands</button></div>`)
	return publisher.Show(peb.UI{
		ID:      "menu",
		Title:   "Quick Commands",
		Anchor:  c.ID,
		Kind:    peb.ContextMenu,
		Context: c,
		View: &peb.EditorView{
			HTML: body.String(),
			CSS:  `.quick-commands-menu { display: flex; flex-direction: column; gap: 4px; min-width: 220px; } button { padding: 7px 12px; text-align: left; }`,
		},
	})
}

func handleEvent(event peb.UIEvent) error {
	if event.Event.Kind == "resync" {
		return openMenu(event.Context)
	}
	return nil
}

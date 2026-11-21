package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pencode-labs/pencode/peb"
)

type command struct {
	Name    string `json:"name"`
	Command string `json:"command"`
	Icon    string `json:"icon,omitempty"`
	CWD     string `json:"cwd,omitempty"`
}

const configPath = ".pencode/quick-commands.json"

func main() {
	if err := peb.ServeExtension(os.Stdin, os.Stdout, peb.Contributions{
		Actions: []peb.Action{{
			ID:       "quick-commands",
			Tooltip:  "Quick Commands",
			Icon:     "resources/command.svg",
			Location: peb.TitleBarRight,
		}},
	}, peb.Handlers{
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

var publisher *peb.ActionPublisher

const menuCSS = `.quick-commands-menu {
    display: flex;
    flex-direction: column;
    width: 100%;
    min-width: 0;
    max-height: min(420px, calc(100vh - 32px));
    gap: 2px;
    overflow-x: hidden;
    overflow-y: auto;
    scrollbar-width: thin;
    scrollbar-color: var(--vscode-scrollbarSlider-background) transparent;
}
.quick-commands-menu button {
    display: flex;
    align-items: center;
    gap: 8px;
    width: 100%;
    min-height: 32px;
    padding: 6px 10px;
    border: 0;
    border-radius: 6px;
    color: var(--vscode-editorWidget-foreground);
    text-align: left;
}
.quick-commands-menu img {
    width: 16px;
    height: 16px;
    flex-shrink: 0;
    object-fit: contain;
}
.quick-commands-menu span {
    min-width: 0;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}
.quick-commands-menu p {
    margin: 0;
    padding: 6px 10px;
    color: var(--vscode-descriptionForeground);
    overflow-wrap: anywhere;
}
.quick-commands-menu button:hover,
.quick-commands-menu button:focus-visible {
    background: var(--vscode-list-hoverBackground);
}`

func readCommands(project string) ([]command, error) {
	data, err := os.ReadFile(filepath.Join(project, configPath))
	if err != nil {
		return nil, err
	}
	var commands []command
	if err := json.Unmarshal(data, &commands); err != nil {
		return nil, err
	}
	if commands == nil {
		return nil, fmt.Errorf("expected an array of commands")
	}
	for i, item := range commands {
		if strings.TrimSpace(item.Name) == "" || strings.TrimSpace(item.Command) == "" {
			return nil, fmt.Errorf("entry %d requires a name and command", i+1)
		}
	}
	return commands, nil
}

func iconSource(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if strings.HasPrefix(value, "<") {
		return "data:image/svg+xml;base64," + base64.StdEncoding.EncodeToString([]byte(value))
	}
	if u, err := url.Parse(value); err == nil && u.Hostname() != "" && (u.Scheme == "http" || u.Scheme == "https") {
		return html.EscapeString(value)
	}
	return ""
}

func menuView(project string) *peb.EditorView {
	const start = `<div class="quick-commands-menu" role="menu">`
	view := &peb.EditorView{CSS: menuCSS}
	if !filepath.IsAbs(project) {
		view.HTML = start + `<p>Open a project to use Quick Commands.</p></div>`
		return view
	}
	commands, err := readCommands(project)
	if os.IsNotExist(err) {
		view.HTML = start + `<button data-ref="create" type="button" role="menuitem">Create Quick Commands configuration</button></div>`
		view.Events = []peb.EditorBinding{{Ref: "create", Kind: "click"}}
		return view
	}
	var body strings.Builder
	body.WriteString(start)
	if err != nil {
		body.WriteString(`<p>` + configPath + `: ` + html.EscapeString(err.Error()) + `</p>`)
	}
	for i, item := range commands {
		body.WriteString(`<button data-ref="command-` + strconv.Itoa(i) + `" data-pencode-action="run-terminal" data-pencode-terminal-command="`)
		body.WriteString(html.EscapeString(item.Command))
		if item.CWD != "" {
			body.WriteString(`" data-pencode-terminal-cwd="`)
			body.WriteString(html.EscapeString(item.CWD))
		}
		body.WriteString(`" type="button" role="menuitem">`)
		if src := iconSource(item.Icon); src != "" {
			body.WriteString(`<img alt="" src="` + src + `">`)
		}
		body.WriteString(`<span>` + html.EscapeString(item.Name) + `</span></button>`)
	}
	body.WriteString(`<button data-ref="configure" data-pencode-href="` + configPath + `" type="button" role="menuitem">Edit Quick Commands configuration</button></div>`)
	view.HTML = body.String()
	return view
}

func openMenu(c peb.ActionContext) error {
	return publisher.Show(peb.UI{
		ID:        "menu",
		Title:     "Quick Commands",
		Anchor:    c.ID,
		Kind:      peb.ContextMenu,
		Context:   c,
		Placement: peb.UIPlacement{Width: 280},
		View:      menuView(c.ProjectPath),
	})
}

func handleEvent(event peb.UIEvent) error {
	if event.ID == "menu" && event.Event.Kind == "click" && event.Event.Ref == "create" {
		if err := createConfig(event.Context.ProjectPath); err != nil {
			return err
		}
		return openMenu(event.Context)
	}
	if event.Event.Kind == "resync" {
		return openMenu(event.Context)
	}
	return nil
}

func createConfig(project string) error {
	if !filepath.IsAbs(project) {
		return fmt.Errorf("open a project to configure Quick Commands")
	}
	root, err := os.OpenRoot(project)
	if err != nil {
		return err
	}
	defer root.Close()
	if err := root.MkdirAll(".pencode", 0o755); err != nil {
		return err
	}
	file, err := root.OpenFile(configPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if os.IsExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	_, err = file.WriteString("[]\n")
	closeErr := file.Close()
	if err != nil {
		return err
	}
	return closeErr
}

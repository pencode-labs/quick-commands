package main

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProjectCommandsReloadOnOpen(t *testing.T) {
	first, second := t.TempDir(), t.TempDir()
	if view := menuView(first); len(view.Events) != 1 || view.Events[0].Ref != "create" {
		t.Fatal(view)
	}
	if _, err := os.Stat(filepath.Join(first, configPath)); !os.IsNotExist(err) {
		t.Fatalf("opening the menu created a config: %v", err)
	}
	if err := createConfig(first); err != nil {
		t.Fatal(err)
	}
	const svg = "<svg viewBox='0 0 16 16'><path d='M1 1h14v14z'/></svg>"
	path := filepath.Join(first, configPath)
	data := `[{"name":"A & B","command":"echo \"A\"","icon":"` + svg + `"},{"name":"Web","command":"npm start","cwd":"apps/web & tools","icon":"https://example.org/icon.png?a=1&b=2"},{"name":"Plain","command":"ls"}]`
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := createConfig(first); err != nil {
		t.Fatal(err)
	}
	menu := menuView(first).HTML
	for _, want := range []string{"A &amp; B", `data-pencode-terminal-command="echo &#34;A&#34;"`, `data-pencode-terminal-command="npm start" data-pencode-terminal-cwd="apps/web &amp; tools"`, "data:image/svg+xml;base64," + base64.StdEncoding.EncodeToString([]byte(svg)), `src="https://example.org/icon.png?a=1&amp;b=2"`, `data-pencode-href="` + configPath + `"`} {
		if !strings.Contains(menu, want) {
			t.Fatalf("missing %q in %s", want, menu)
		}
	}
	if strings.Count(menu, "<img") != 2 || strings.Count(menu, `data-pencode-action="run-terminal"`) != 3 {
		t.Fatal(menu)
	}
	if strings.Count(menu, `data-pencode-terminal-cwd=`) != 1 {
		t.Fatal("only the configured command should have a working directory", menu)
	}
	for _, project := range []string{second, ""} {
		if menu := menuView(project).HTML; strings.Contains(menu, "run-terminal") {
			t.Fatal("another project inherited commands", menu)
		}
	}
	for _, invalid := range []string{`[`, `null`, `[{"name":"invalid","command":" " }]`} {
		if err := os.WriteFile(path, []byte(invalid), 0o644); err != nil {
			t.Fatal(err)
		}
		menu := menuView(first).HTML
		if strings.Contains(menu, "run-terminal") || !strings.Contains(menu, `data-pencode-href="`+configPath+`"`) || !strings.Contains(menu, "<p>") {
			t.Fatal("invalid configuration should show an error and its edit link", menu)
		}
	}
}

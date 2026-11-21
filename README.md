# Quick Commands

Quick Commands adds a button to the editor title bar for running the current project's terminal commands.

## Configuration

Store commands in `.pencode/quick-commands.json` under the project root. Open the Quick Commands menu to create the file or edit an existing one. Save it and reopen the menu to load your changes. If you previously saved commands in extension settings, copy the JSON array into the projects that use it.

```json
[
  {
    "name": "Run tests",
    "command": "go test ./...",
    "icon": "<svg viewBox='0 0 16 16'><path fill='currentColor' d='M3 2l10 6-10 6z'/></svg>"
  },
  {
    "name": "Start development server",
    "command": "npm run dev",
    "cwd": "apps/web",
    "icon": "https://cdn.jsdelivr.net/gh/devicons/devicon@v2.16.0/icons/nodejs/nodejs-original.svg"
  }
]
```

Each entry requires a non-empty `name` and `command`. The optional `icon` accepts SVG markup or an HTTP(S) URL to a PNG, JPEG, WebP, or SVG image. Omit it for a text-only menu item.

## Execution behavior

Set the optional `cwd` field to start a command in a specific directory. Relative paths are resolved from the project root, and absolute paths refer to the machine running the workspace. An empty or omitted `cwd` uses the project root.

Selecting a command opens a new terminal tab and sends the configured command to the workspace shell. When the command finishes, the terminal stays open so you can review the output and run more commands. Close the tab or enter `exit` to end the shell session.

## Build and install

The main Pencode repository builds and installs this extension together with the other bundled extensions.

To build it independently from this repository:

```sh
go test ./...
go build -o bin/quick-commands .
```

The extension manifest is in `resources/extension.json`.

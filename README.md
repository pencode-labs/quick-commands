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
    "icon": "https://cdn.jsdelivr.net/gh/devicons/devicon@v2.16.0/icons/nodejs/nodejs-original.svg"
  }
]
```

Each entry requires a non-empty `name` and `command`. The optional `icon` accepts SVG markup or an HTTP(S) URL to a PNG, JPEG, WebP, or SVG image. Omit it for a text-only menu item.

## Execution behavior

Selecting a command opens a new terminal tab and sends the configured command to the workspace shell. The extension then sends `exit` after the command, so the terminal tab closes automatically when the command finishes. Long-running commands keep the terminal tab open until they terminate.

## Build and install

The main Pencode repository builds and installs this extension together with the other bundled extensions.

To build it independently from this repository:

```sh
go test ./...
go build -o bin/quick-commands .
```

The extension manifest is in `resources/extension.json`.

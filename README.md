# Quick Commands

Quick Commands adds a command palette button to the editor title bar. It lets you keep frequently used terminal commands in the extension settings and run them from a drop-down menu.

## Configuration

Open the editor settings, select **Extensions**, and choose **Quick Commands**. Set `Commands` to a JSON string containing an array of command objects.

```json
[
  {
    "name": "Run tests",
    "command": "go test ./..."
  },
  {
    "name": "Start development server",
    "command": "npm run dev"
  }
]
```

Each command requires a non-empty `name` and `command` field. The name appears in the drop-down menu. The command is entered into a new terminal session in the current workspace directory.

## Execution behavior

Selecting a command opens a new terminal tab and sends the configured command to the workspace shell. The extension then sends `exit` after the command, so the terminal tab closes automatically when the command finishes. Long-running commands keep the terminal tab open until they terminate.

The **Configure Quick Commands** menu item opens the extension settings.

## Build and install

The main Pencode repository builds and installs this extension together with the other bundled extensions.

To build it independently from this repository:

```sh
go test ./...
go build -o bin/quick-commands .
```

The extension manifest is in `resources/extension.json`.

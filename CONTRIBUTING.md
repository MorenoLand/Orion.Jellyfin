# Contributing

## Before opening a change

Read the existing source and preserve behavior at the native-host boundary. The main window is a direct remote URL, while the Wails frontend is the local asset shell used by development and fallback content.

Install dependencies and run the checks:

~~~text
go mod download
npm --prefix frontend install
npm --prefix frontend run build
gofmt -w main.go adblock.go icon.go scripts.go zorder_windows.go zorder_other.go
go test ./...
git diff --check
~~~

When changing the frontend, verify the keyboard shortcuts, title/header drag, double-click maximize, theater mode, tray show/hide, close-to-hide behavior, custom op.txt URL/title handling, favicon fallback, and cosmetic filtering.

Do not commit generated frontend dist or bin output, local op.txt files, credentials, private certificates, or runtime state. Do not reintroduce src-tauri or Docker tasks.

## Pull requests

Describe the user-visible behavior changed and the evidence used to verify it. Include platform limitations when a check was static or build-only. Keep commits focused, signed, and easy to review.

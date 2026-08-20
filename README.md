# Jellyfin Desktop Client

Jellyfin is a small desktop client for the Jellyfin instance at https://booty.moreno.land. It keeps the original Tauri behavior while using Wails 3 and Go for the native host.

## Features

- Opens the configured Jellyfin URL directly in the native webview.
- Reads an optional op.txt beside the executable: line 1 overrides the URL and line 2 overrides the window title.
- Uses a tray icon to show or hide the main window, with Show and Quit menu entries.
- Uses the configured site favicon when a custom URL is supplied.
- Hides the main window instead of closing it from the window close action.
- Supports T for theater mode, X for immersive mode, double-click maximize on the Jellyfin header or videoOsdPage, and drag gestures from the same regions.
- Uses the monitor work area for manual maximize and the monitor bounds for theater mode.
- Applies cosmetic rules from the bundled EasyList and EasyPrivacy lists.

## Repository branches

features/tauri is the exact Tauri baseline preserved before the conversion. features/wails3 contains the Wails 3 implementation.

## Build

Requirements:

- Go 1.26 or newer
- Node.js and npm
- Wails 3 beta.11
- WebView2 on Windows

Install the Wails CLI:

~~~text
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0-beta.11
~~~

Build the frontend and application:

~~~text
npm --prefix frontend install
npm --prefix frontend run build
wails3 build
~~~

Windows builds apply UPX compression with `--best --lzma` after linking; set the `COMPRESS` task variable to `false` when an uncompressed diagnostic executable is needed.

For development, use:

~~~text
wails3 dev -config ./build/config.yml
~~~

The generated application is written under bin. The frontend dist directory is generated output and is not committed.

## Runtime configuration

Create op.txt beside the executable when a different site or title is required:

~~~text
https://example.invalid
Example Jellyfin
~~~

The first line is the URL. The second line is optional and becomes the native window and tray title. The default URL and default title are used when the file or line is absent.

The browser is started with certificate-error checks disabled to preserve the original client behavior. Do not use an untrusted URL in op.txt.

## Contributions

See CONTRIBUTING.md for the development and review requirements. Security issues should be reported privately as described in SECURITY.md.

## Third-party data

easylist.txt and easyprivacy.txt are upstream filter lists retained with their source headers. They are maintained by EasyList and EasyPrivacy under the terms published at https://easylist.to/pages/licence.html.

## License

This project is licensed under the MIT License. See LICENSE.

# Security policy

## Supported versions

Only the latest commit on the features/wails3 branch is supported for security fixes. The features/tauri branch is a historical baseline and is not the maintained implementation.

## Reporting a vulnerability

Please report vulnerabilities through a private GitHub Security Advisory for MorenoLand/Moreno.Jellyfin. Include the affected commit, operating system, reproduction steps, impact, and any logs that do not contain credentials or private URLs. Do not disclose an unpatched vulnerability in a public issue.

## Security boundaries

The application intentionally loads a remote web page and exposes the Wails window bridge to that page. A page loaded by the application can request the documented window actions, including dragging, hiding, manual maximize, and theater mode. Treat the configured URL and any op.txt file as trusted input.

The Windows webview is launched with certificate-error checks disabled because that is part of the original client behavior. This weakens transport verification and makes an untrusted or intercepted URL unsafe.

The bundled EasyList and EasyPrivacy data is used for cosmetic filtering only. It is not a malware scanner, network firewall, or substitute for transport security. Custom favicon retrieval follows the configured URL host and falls back to the embedded application icon when retrieval or image decoding fails.

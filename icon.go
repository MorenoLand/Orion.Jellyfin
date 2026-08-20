package main

import (
	"bytes"
	_ "embed"
	"image"
	"image/png"
	"net/http"
	"net/url"
	"time"
)

//go:embed app-icon.png
var embeddedIcon []byte

func loadIcon(targetURL string) []byte {
	if targetURL == jellyfinURL {
		return append([]byte(nil), embeddedIcon...)
	}
	parsed, err := url.Parse(targetURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return append([]byte(nil), embeddedIcon...)
	}
	faviconURL := parsed.Scheme + "://" + parsed.Host + "/favicon.ico"
	client := &http.Client{Timeout: 5 * time.Second}
	response, err := client.Get(faviconURL)
	if err != nil {
		return append([]byte(nil), embeddedIcon...)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return append([]byte(nil), embeddedIcon...)
	}
	icon, _, err := image.Decode(response.Body)
	if err != nil {
		return append([]byte(nil), embeddedIcon...)
	}
	var encoded bytes.Buffer
	if err := png.Encode(&encoded, icon); err != nil {
		return append([]byte(nil), embeddedIcon...)
	}
	return encoded.Bytes()
}

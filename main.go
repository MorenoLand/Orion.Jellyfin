package main

import (
	"embed"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed frontend/dist
var assets embed.FS

type hostMessage struct {
	Action string
	URL    string
	Title  string
}

type appState struct {
	app         *application.App
	main        *application.WebviewWindow
	theater     *application.WebviewWindow
	customURL   bool
	customTitle bool
	quitting    atomic.Bool
	maxMu       sync.Mutex
	maximized   bool
	restored    application.Rect
	theaterMu   sync.Mutex
	theaterOpen bool
}

func main() {
	targetURL, title, customURL := loadOptions()
	state := &appState{customURL: customURL, customTitle: title != ""}
	icon := loadIcon(targetURL)
	app := application.New(application.Options{
		Name:              "Jellyfin",
		Description:       "Jellyfin Desktop Client",
		Assets:            application.AssetOptions{Handler: application.AssetFileServerFS(assets)},
		Windows:           application.WindowsOptions{AdditionalBrowserArgs: []string{"--ignore-certificate-errors"}},
		RawMessageHandler: state.handleMessage,
	})
	state.app = app
	app.SetIcon(icon)

	mainWindow := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:            "main",
		Title:           title,
		URL:             targetURL,
		Width:           1280,
		Height:          720,
		MinWidth:        640,
		MinHeight:       400,
		InitialPosition: application.WindowCentered,
		Frameless:       true,
	})
	state.main = mainWindow
	injectPage := func(_ *application.WindowEvent) { mainWindow.ExecJS(injectionScript(customURL)) }
	mainWindow.RegisterHook(events.Windows.WebViewNavigationCompleted, injectPage)
	mainWindow.RegisterHook(events.Mac.WebViewDidFinishNavigation, injectPage)
	mainWindow.RegisterHook(events.Linux.WindowLoadFinished, injectPage)

	theater := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:          "theater",
		Title:         "Jellyfin",
		HTML:          theaterHTML,
		Width:         1,
		Height:        1,
		Frameless:     true,
		Hidden:        true,
		DisableResize: true,
		Windows:       application.WindowsWindow{HiddenOnTaskbar: true},
	})
	state.theater = theater

	mainWindow.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
		if state.quitting.Load() {
			return
		}
		event.Cancel()
		mainWindow.Hide()
	})

	setupTray(state, title, icon)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func loadOptions() (targetURL, title string, customURL bool) {
	targetURL = jellyfinURL
	title = "Jellyfin"
	executable, err := os.Executable()
	if err != nil {
		return targetURL, title, false
	}
	data, err := os.ReadFile(filepath.Join(filepath.Dir(executable), "op.txt"))
	if err != nil {
		return targetURL, title, false
	}
	lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
	if len(lines) > 0 && strings.TrimSpace(lines[0]) != "" {
		targetURL = strings.TrimSpace(lines[0])
	}
	if len(lines) > 1 && strings.TrimSpace(lines[1]) != "" {
		title = strings.TrimSpace(lines[1])
	}
	return targetURL, title, targetURL != jellyfinURL
}

func setupTray(state *appState, title string, icon []byte) {
	tray := state.app.SystemTray.New()
	tray.SetIcon(icon)
	tray.SetTooltip(title)
	menu := state.app.NewMenu()
	menu.Add("Show " + title).OnClick(func(_ *application.Context) {
		state.showMain()
	})
	menu.Add("Quit").OnClick(func(_ *application.Context) {
		state.quitting.Store(true)
		state.app.Quit()
	})
	tray.SetMenu(menu)
	tray.OnClick(func() {
		state.toggleMain()
	})
	tray.OnRightClick(func() {
		tray.ShowMenu()
	})
}

func (s *appState) handleMessage(window application.Window, message string, _ *application.OriginInfo) {
	var request hostMessage
	if json.Unmarshal([]byte(message), &request) != nil {
		return
	}
	switch request.Action {
	case "hide":
		window.Hide()
	case "toggle_maximize":
		s.toggleMaximize(window)
	case "toggle_theater":
		s.toggleTheater()
	case "set_title":
		if s.customURL && !s.customTitle && request.Title != "" {
			window.SetTitle(request.Title)
		}
	case "page_load":
		if request.URL == "" {
			return
		}
		css, script := cosmeticResources(request.URL)
		if css != "" {
			window.ExecJS(injectCSS(css))
		}
		if script != "" {
			window.ExecJS(script)
		}
		if s.customURL && !s.customTitle && request.Title != "" {
			window.SetTitle(request.Title)
		}
	}
}

func (s *appState) toggleMain() {
	if s.main.IsVisible() {
		s.main.Hide()
		return
	}
	s.main.Show()
	s.main.Focus()
}

func (s *appState) showMain() {
	s.main.Show()
	s.main.Focus()
}

func (s *appState) toggleMaximize(window application.Window) {
	w, ok := window.(*application.WebviewWindow)
	if !ok {
		return
	}
	s.maxMu.Lock()
	if s.maximized {
		restored := s.restored
		s.maximized = false
		s.maxMu.Unlock()
		w.SetPhysicalBounds(restored)
		return
	}
	screen, err := w.GetScreen()
	if err != nil {
		s.maxMu.Unlock()
		return
	}
	s.restored = w.PhysicalBounds()
	s.maximized = true
	s.maxMu.Unlock()
	w.SetPhysicalBounds(screen.PhysicalWorkArea)
}

func (s *appState) toggleTheater() {
	s.theaterMu.Lock()
	if s.theaterOpen {
		s.theaterOpen = false
		s.theaterMu.Unlock()
		s.theater.Hide()
		return
	}
	screen, err := s.main.GetScreen()
	if err != nil {
		s.theaterMu.Unlock()
		return
	}
	s.theaterOpen = true
	s.theaterMu.Unlock()
	s.theater.SetPhysicalBounds(screen.PhysicalBounds)
	s.theater.Show()
	setWindowAlwaysOnBottom(s.theater)
	s.main.Focus()
}

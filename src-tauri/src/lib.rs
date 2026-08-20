mod adblock;

use image::GenericImageView;
use std::path::Path;
use std::sync::Mutex;
use tauri::menu::{Menu, MenuItem};
use tauri::tray::{MouseButton, MouseButtonState, TrayIconBuilder, TrayIconEvent};
use tauri::Manager;

struct TrayHandle(#[allow(dead_code)] tauri::tray::TrayIcon<tauri::Wry>);

struct MaxState {
    maximized: bool,
    pos: Option<(i32, i32)>,
    size: Option<(u32, u32)>,
}

struct TheaterState {
    on: bool,
}

const JELLYFIN_URL: &str = "https://booty.moreno.land";

const THEATER_HTML: &str = r#"data:text/html,<html><body style="margin:0;background:%23000;width:100vw;height:100vh;cursor:none"></body></html>"#;

const GENERIC_INJECT: &str = r#"(function(){
  function inj(){
    if(document.getElementById('gf-style'))return;
    var s=document.createElement('style');s.id='gf-style';
    s.textContent='#gf-toast{position:fixed;bottom:20px;left:50%;transform:translateX(-50%);z-index:2147483647;background:rgba(0,0,0,.8);color:#fff;padding:8px 18px;border-radius:8px;font:13px system-ui;pointer-events:none;opacity:0;transition:opacity .3s}#gf-toast.show{opacity:1}:-webkit-full-screen{width:100vw!important;height:100vh!important}';
    document.head.appendChild(s);
    var toast=document.createElement('div');toast.id='gf-toast';document.body.appendChild(toast);
    var toastTimer;
    function showToast(msg){toast.textContent=msg;toast.classList.add('show');clearTimeout(toastTimer);toastTimer=setTimeout(function(){toast.classList.remove('show')},2000);}
    var immersive=false;
    document.addEventListener('mousedown',function(e){
      if(e.button!==0)return;
      if(e.target.closest('button,a,input,select,textarea,[contenteditable],iframe,video,audio,canvas'))return;
      if(immersive)return;
      e.preventDefault();
      var sx=e.clientX,sy=e.clientY,moved=false;
      function onMove(ev){if(Math.abs(ev.clientX-sx)>3||Math.abs(ev.clientY-sy)>3){moved=true;if(window.__TAURI__&&window.__TAURI__.core)window.__TAURI__.core.invoke('drag_window');document.removeEventListener('mousemove',onMove);document.removeEventListener('mouseup',onUp);}}
      function onUp(){document.removeEventListener('mousemove',onMove);document.removeEventListener('mouseup',onUp);}
      document.addEventListener('mousemove',onMove);document.addEventListener('mouseup',onUp);
    },true);
    document.addEventListener('dblclick',function(e){
      if(e.target.closest('button,a,input,select,textarea,[contenteditable],iframe,video,audio,canvas'))return;
      if(immersive)return;
      try{window.__TAURI_INTERNALS__.invoke('toggle_maximize');}catch(ex){try{window.__TAURI__.core.invoke('toggle_maximize');}catch(ex2){}}
      e.preventDefault();e.stopImmediatePropagation();
    },true);
    document.addEventListener('keydown',function(e){
      if(e.target.closest('input,textarea,[contenteditable],select'))return;
      if(e.key==='x'||e.key==='X'){
        if(!e.ctrlKey&&!e.altKey&&!e.metaKey&&!e.shiftKey){
          immersive=!immersive;
          showToast(immersive?'Immersive mode on':'Immersive mode off');
          e.preventDefault();e.stopImmediatePropagation();return false;
        }
      }
    },true);
  }
  inj();setTimeout(inj,500);setTimeout(inj,1500);setTimeout(inj,3000);
})();"#;

const INJECT: &str = r#"(function(){
  function inj(){
    if(document.getElementById('jf-style'))return;
    var s=document.createElement('style');s.id='jf-style';
    s.textContent='.skinHeader,.headerTop,.headerLeft,.headerRight,.MuiToolbar-root{-webkit-app-region:drag;app-region:drag}.skinHeader button,.skinHeader a,.skinHeader input,.skinHeader select,.skinHeader textarea,.MuiToolbar-root button,.MuiToolbar-root a,.MuiToolbar-root input,.MuiToolbar-root select{-webkit-app-region:no-drag;app-region:no-drag}*{-webkit-user-select:none;user-select:none}input,textarea,[contenteditable]{-webkit-user-select:text;user-select:text}#jf-immersive-overlay{position:fixed;inset:0;z-index:2147483646}#jf-toast{position:fixed;bottom:20px;left:50%;transform:translateX(-50%);z-index:2147483647;background:rgba(0,0,0,.8);color:#fff;padding:8px 18px;border-radius:8px;font:13px system-ui;pointer-events:none;opacity:0;transition:opacity .3s}#jf-toast.show{opacity:1}:-webkit-full-screen{width:100vw!important;height:100vh!important}';
    document.head.appendChild(s);
    var toast=document.createElement('div');toast.id='jf-toast';document.body.appendChild(toast);
    var toastTimer;
    function showToast(msg){toast.textContent=msg;toast.classList.add('show');clearTimeout(toastTimer);toastTimer=setTimeout(function(){toast.classList.remove('show')},2000);}
    var immersive=false,ov;
    function blockAll(e){e.stopImmediatePropagation();e.preventDefault();}
    var blockTypes=['mousemove','mouseover','mouseenter','mousedown','mouseup','click','dblclick','contextmenu','pointerdown','pointerup','pointermove','pointerover','pointerenter'];
    document.addEventListener('keydown',function(e){
      if(e.target.closest('input,textarea,[contenteditable],select'))return;
      if(e.key==='t'||e.key==='T'){
        if(!e.ctrlKey&&!e.altKey&&!e.metaKey&&!e.shiftKey){
          if(window.__TAURI__&&window.__TAURI__.core)window.__TAURI__.core.invoke('toggle_theater');
          e.preventDefault();e.stopImmediatePropagation();return false;
        }
      }
      if(e.key==='x'||e.key==='X'){
        if(!e.ctrlKey&&!e.altKey&&!e.metaKey&&!e.shiftKey){
          immersive=!immersive;
          if(immersive){
            ov=document.createElement('div');ov.id='jf-immersive-overlay';document.body.appendChild(ov);
            var osd=document.querySelector('div#videoOsdPage');if(osd)osd.style.display='none';
            blockTypes.forEach(function(t){document.addEventListener(t,blockAll,true);});
            showToast('Immersive mode on');
          }else{
            if(ov)ov.remove();ov=null;
            var osd=document.querySelector('div#videoOsdPage');if(osd)osd.style.display='';
            blockTypes.forEach(function(t){document.removeEventListener(t,blockAll,true);});
            showToast('Immersive mode off');
          }
          e.preventDefault();e.stopImmediatePropagation();return false;
        }
      }
      if(immersive){e.preventDefault();e.stopImmediatePropagation();}
    },true);
     document.addEventListener('dblclick',function(e){
       if(e.target.closest('button,a,input,select,textarea'))return;
       var hdr=e.target.closest('.skinHeader,.headerTop,.headerLeft,.headerRight,.MuiToolbar-root,div#videoOsdPage');
       if(!hdr)return;
       if(window.__TAURI__&&window.__TAURI__.core)window.__TAURI__.core.invoke('toggle_maximize');
       e.preventDefault();e.stopImmediatePropagation();
     },true);
     document.addEventListener('mousedown',function(e){
       if(e.button!==0)return;
       if(e.target.closest('button,a,input,select,textarea,.sliderContainer,.volumeSlider,.osdVolumeSlider'))return;
       var osd=e.target.closest('div#videoOsdPage');
       if(!osd)return;
       var sx=e.clientX,sy=e.clientY,moved=false;
      function onMove(ev){if(Math.abs(ev.clientX-sx)>3||Math.abs(ev.clientY-sy)>3){moved=true;try{window.__TAURI__.core.invoke('drag_window');}catch(ex){}document.removeEventListener('mousemove',onMove);document.removeEventListener('mouseup',onUp);}}
       function onUp(){document.removeEventListener('mousemove',onMove);document.removeEventListener('mouseup',onUp);}
       document.addEventListener('mousemove',onMove);document.addEventListener('mouseup',onUp);
     },true);
   }
  inj();setTimeout(inj,500);setTimeout(inj,1500);setTimeout(inj,3000);
})();"#;

fn load_icon() -> tauri::image::Image<'static> {
    let bytes = include_bytes!("../icons/icon.png");
    if let Ok(img) = image::load_from_memory(bytes) {
        let (w, h) = img.dimensions();
        let rgba = img.into_rgba8().into_raw();
        return tauri::image::Image::new_owned(rgba, w, h);
    }
    let img = image::RgbaImage::from_pixel(64, 64, image::Rgba([120, 70, 220, 255]));
    tauri::image::Image::new_owned(img.into_raw(), 64, 64)
}

fn fetch_favicon(url: &str) -> Option<tauri::image::Image<'static>> {
    let parsed = url::Url::parse(url).ok()?;
    let scheme = parsed.scheme();
    let host = parsed.host_str()?;
    let favicon_url = format!("{}://{}/favicon.ico", scheme, host);
    let response = minreq::get(&favicon_url).with_timeout(5).send().ok()?;
    if response.status_code != 200 { return None; }
    let bytes = response.as_bytes();
    let img = image::load_from_memory(bytes).ok()?;
    let (w, h) = img.dimensions();
    let rgba = img.into_rgba8().into_raw();
    Some(tauri::image::Image::new_owned(rgba, w, h))
}

fn load_app_icon(custom_url: Option<&str>) -> tauri::image::Image<'static> {
    if let Some(url) = custom_url {
        if let Some(icon) = fetch_favicon(url) {
            return icon;
        }
    }
    load_icon()
}

#[tauri::command]
fn set_title(window: tauri::WebviewWindow, title: String) {
    window.set_title(&title).ok();
}

#[tauri::command]
fn drag_window(window: tauri::WebviewWindow) {
    window.start_dragging().ok();
}

#[tauri::command]
fn toggle_maximize(window: tauri::WebviewWindow, state: tauri::State<'_, Mutex<MaxState>>) -> Result<(), String> {
    let mut s = state.lock().map_err(|_| "state poisoned")?;
    if s.maximized {
        if let (Some((x, y)), Some((w, h))) = (s.pos, s.size) {
            window.set_size(tauri::Size::Physical(tauri::PhysicalSize::new(w, h))).map_err(|e| e.to_string())?;
            window.set_position(tauri::Position::Physical(tauri::PhysicalPosition::new(x, y))).map_err(|e| e.to_string())?;
        }
        s.maximized = false;
    } else {
        let pos = window.outer_position().map_err(|e| e.to_string())?;
        let size = window.outer_size().map_err(|e| e.to_string())?;
        s.pos = Some((pos.x, pos.y));
        s.size = Some((size.width, size.height));
        let monitor = window.current_monitor().map_err(|e| e.to_string())?.ok_or("no monitor")?;
        let wa = monitor.work_area();
        window.set_position(tauri::Position::Physical(tauri::PhysicalPosition::new(wa.position.x, wa.position.y))).map_err(|e| e.to_string())?;
        window.set_size(tauri::Size::Physical(tauri::PhysicalSize::new(wa.size.width, wa.size.height))).map_err(|e| e.to_string())?;
        s.maximized = true;
    }
    Ok(())
}

#[tauri::command]
fn toggle_theater(app: tauri::AppHandle, window: tauri::WebviewWindow, state: tauri::State<'_, Mutex<TheaterState>>) -> Result<(), String> {
    let mut s = state.lock().map_err(|_| "state poisoned")?;
    if let Some(dark) = app.get_webview_window("theater") {
        if s.on {
            dark.hide().map_err(|e| e.to_string())?;
            s.on = false;
        } else {
            let monitor = window.current_monitor().map_err(|e| e.to_string())?.ok_or("no monitor")?;
            let screen = monitor.size();
            let pos = monitor.position();
            dark.set_size(tauri::Size::Physical(tauri::PhysicalSize::new(screen.width, screen.height))).map_err(|e| e.to_string())?;
            dark.set_position(tauri::Position::Physical(tauri::PhysicalPosition::new(pos.x, pos.y))).map_err(|e| e.to_string())?;
            dark.show().map_err(|e| e.to_string())?;
            window.set_focus().map_err(|e| e.to_string())?;
            s.on = true;
        }
    }
    Ok(())
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .manage(Mutex::new(MaxState { maximized: false, pos: None, size: None }))
        .manage(Mutex::new(TheaterState { on: false }))
        .setup(|app| {
            let op_content = std::fs::read_to_string(
                Path::new(std::env::current_exe().unwrap().parent().unwrap()).join("op.txt")
            ).unwrap_or_default();
            let mut op_lines: Vec<&str> = op_content.lines().collect();
            let url_str = op_lines.get(0).map(|s| s.trim()).filter(|s| !s.is_empty()).unwrap_or(JELLYFIN_URL).to_string();
            let custom_title = op_lines.get(1).map(|s| s.trim()).filter(|s| !s.is_empty()).map(|s| s.to_string());
            let has_custom_url = url_str != JELLYFIN_URL;
            let title = custom_title.clone().unwrap_or_else(|| "Jellyfin".to_string());

            let icon = load_app_icon(if has_custom_url { Some(&url_str) } else { None });

            let _window = tauri::WebviewWindowBuilder::new(
                app, "main",
                tauri::WebviewUrl::External(url_str.parse().unwrap()),
            )
                .title(&title)
                .inner_size(1280.0, 720.0)
                .min_inner_size(640.0, 400.0)
                .center()
                .decorations(false)
                .maximizable(false)
                .additional_browser_args("--ignore-certificate-errors")
                .on_page_load(move |webview, payload| {
                    if payload.event() == tauri::webview::PageLoadEvent::Finished {
                        if custom_title.is_none() && has_custom_url {
                            webview.eval("if(window.__TAURI__&&window.__TAURI__.core)window.__TAURI__.core.invoke('set_title',{title:document.title});").ok();
                        }
                        let url = webview.url().map(|u| u.to_string()).unwrap_or_default();
                        let css = adblock::get_cosmetic_css(&url);
                        if !css.is_empty() {
                            let escaped = css.replace('\\', "\\\\").replace('`', "\\`").replace('$', "\\$");
                            let inject_css = format!("var s=document.createElement('style');s.textContent=`{}`;document.head.appendChild(s);", escaped);
                            webview.eval(&inject_css).ok();
                        }
                        let script = adblock::get_injected_script(&url);
                        if !script.is_empty() {
                            webview.eval(&script).ok();
                        }
                        let inject_js = if has_custom_url { GENERIC_INJECT } else { INJECT };
                        webview.eval(inject_js).ok();
                    }
                })
                .build()?;

            let window = app.get_webview_window("main").unwrap();
            window.set_icon(icon.clone()).ok();

            let show_label = format!("Show {}", title);
            let show_item = MenuItem::with_id(app, "show", &show_label, true, None::<&str>)?;
            let quit_item = MenuItem::with_id(app, "quit", "Quit", true, None::<&str>)?;
            let menu = Menu::with_items(app, &[&show_item, &quit_item])?;

            let tray = TrayIconBuilder::new()
                .icon(icon)
                .tooltip(&title)
                .menu(&menu)
                .show_menu_on_left_click(false)
                .on_menu_event(|app, event| match event.id.as_ref() {
                    "show" => {
                        if let Some(w) = app.get_webview_window("main") {
                            w.show().ok();
                            w.set_focus().ok();
                        }
                    }
                    "quit" => app.exit(0),
                    _ => {}
                })
                .on_tray_icon_event(|tray, event| {
                    if let TrayIconEvent::Click { button: MouseButton::Left, button_state: MouseButtonState::Up, .. } = event {
                        let app = tray.app_handle();
                        if let Some(w) = app.get_webview_window("main") {
                            if w.is_visible().unwrap_or(false) { w.hide().ok(); }
                            else { w.show().ok(); w.set_focus().ok(); }
                        }
                    }
                })
                .build(app)?;

            app.manage(TrayHandle(tray));

            let _theater = tauri::WebviewWindowBuilder::new(
                app, "theater",
                tauri::WebviewUrl::External(THEATER_HTML.parse().unwrap()),
            )
                .title("")
                .decorations(false)
                .always_on_bottom(true)
                .skip_taskbar(true)
                .resizable(false)
                .visible(false)
                .build()?;

            let w = window.clone();
            window.on_window_event(move |event| {
                if let tauri::WindowEvent::CloseRequested { api, .. } = event {
                    api.prevent_close();
                    w.hide().ok();
                }
            });

            Ok(())
        })
        .invoke_handler(tauri::generate_handler![toggle_maximize, toggle_theater, set_title, drag_window])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}

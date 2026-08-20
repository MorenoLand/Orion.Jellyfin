package main

import (
	"encoding/json"
	"strconv"
)

const jellyfinURL = "https://booty.moreno.land"

const theaterHTML = `<!doctype html><html><head><meta charset="utf-8"><style>html,body{width:100%;height:100%;margin:0;overflow:hidden;background:#000}</style></head><body></body></html>`

func injectionScript(customURL bool) string {
	if customURL {
		return genericInjection
	}
	return jellyfinInjection
}

func injectHost(action string, fields map[string]string) string {
	payload := map[string]string{"action": action}
	for key, value := range fields {
		payload[key] = value
	}
	data, _ := json.Marshal(payload)
	return "try{window._wails&&window._wails.invoke&&window._wails.invoke(" + strconv.Quote(string(data)) + ")}catch(_e){}"
}

const genericInjection = `(function(){
  const host=(action,data={})=>{
    try{
      if(!window._wails||typeof window._wails.invoke!=='function')return;
      if(action==='drag_window'){window._wails.invoke('wails:drag');return}
      window._wails.invoke(JSON.stringify(Object.assign({action},data)))
    }catch(_e){}
  };
  const style=document.createElement('style');
  style.id='gf-style';
  style.textContent='html,body{width:100%!important;height:100%!important;overflow:hidden!important;background:#000!important}body>*{max-width:none!important}.gf-immersive{cursor:none!important}.gf-toast{position:fixed;left:50%;bottom:8%;z-index:2147483647;transform:translateX(-50%);padding:8px 14px;border-radius:4px;background:rgba(0,0,0,.78);color:#fff;font:13px "Segoe UI",sans-serif;pointer-events:none;opacity:0;transition:opacity .15s}.gf-toast.gf-show{opacity:1}';
  (document.head||document.documentElement).appendChild(style);
  const toast=document.createElement('div');
  toast.className='gf-toast';
  document.documentElement.appendChild(toast);
  let toastTimer;
  const showToast=text=>{
    toast.textContent=text;
    toast.classList.add('gf-show');
    clearTimeout(toastTimer);
    toastTimer=setTimeout(()=>toast.classList.remove('gf-show'),1200)
  };
  let immersive=false;
  const excluded=target=>target&&target.closest('button,a,input,select,textarea,[contenteditable="true"],iframe,video,audio,canvas');
  document.addEventListener('mousedown',event=>{
    if(event.button!==0||excluded(event.target))return;
    const startX=event.clientX,startY=event.clientY;
    let moved=false;
    const move=moveEvent=>{
      if(moved||Math.hypot(moveEvent.clientX-startX,moveEvent.clientY-startY)<=3)return;
      moved=true;
      host('drag_window');
      cleanup()
    };
    const cleanup=()=>{
      document.removeEventListener('mousemove',move,true);
      document.removeEventListener('mouseup',cleanup,true)
    };
    document.addEventListener('mousemove',move,true);
    document.addEventListener('mouseup',cleanup,true)
  },true);
  document.addEventListener('dblclick',event=>{
    if(excluded(event.target))return;
    host('toggle_maximize')
  },true);
  document.addEventListener('keydown',event=>{
    if(event.defaultPrevented||event.ctrlKey||event.altKey||event.metaKey||event.shiftKey)return;
    if(event.key.toLowerCase()!=='x')return;
    immersive=!immersive;
    document.documentElement.classList.toggle('gf-immersive',immersive);
    showToast(immersive?'Immersive mode':'Normal mode')
  },true);
  host('page_load',{url:location.href,title:document.title})
})()`

const jellyfinInjection = `(function(){
  const host=(action,data={})=>{
    try{
      if(!window._wails||typeof window._wails.invoke!=='function')return;
      if(action==='drag_window'){window._wails.invoke('wails:drag');return}
      window._wails.invoke(JSON.stringify(Object.assign({action},data)))
    }catch(_e){}
  };
  const style=document.createElement('style');
  style.id='jf-style';
  style.textContent='.skinHeader,.headerTop,.headerLeft,.headerRight,.MuiToolbar-root{-webkit-app-region:drag;app-region:drag}.skinHeader button,.skinHeader a,.skinHeader input,.skinHeader select,.skinHeader textarea,.headerTop button,.headerTop a,.headerLeft button,.headerLeft a,.headerRight button,.headerRight a,.MuiToolbar-root button,.MuiToolbar-root a{-webkit-app-region:no-drag;app-region:no-drag}body{user-select:none!important;-webkit-user-select:none!important}input,select,textarea,[contenteditable="true"]{user-select:text!important;-webkit-user-select:text!important}#jf-immersive-overlay{position:fixed;inset:0;z-index:2147483646;display:none;background:#000;pointer-events:none}.jf-immersive #jf-immersive-overlay{display:block}.jf-toast{position:fixed;left:50%;bottom:8%;z-index:2147483647;transform:translateX(-50%);padding:8px 14px;border-radius:4px;background:rgba(0,0,0,.78);color:#fff;font:13px "Segoe UI",sans-serif;pointer-events:none;opacity:0;transition:opacity .15s}.jf-toast.jf-show{opacity:1}';
  (document.head||document.documentElement).appendChild(style);
  const overlay=document.createElement('div');
  overlay.id='jf-immersive-overlay';
  document.documentElement.appendChild(overlay);
  const toast=document.createElement('div');
  toast.className='jf-toast';
  document.documentElement.appendChild(toast);
  let toastTimer;
  const showToast=text=>{
    toast.textContent=text;
    toast.classList.add('jf-show');
    clearTimeout(toastTimer);
    toastTimer=setTimeout(()=>toast.classList.remove('jf-show'),1200)
  };
  let immersive=false;
  const excluded=target=>target&&target.closest('button,a,input,select,textarea');
  document.addEventListener('keydown',event=>{
    if(event.defaultPrevented||event.ctrlKey||event.altKey||event.metaKey||event.shiftKey)return;
    const key=event.key.toLowerCase();
    if(key==='t'){
      host('toggle_theater');
      return
    }
    if(key!=='x')return;
    immersive=!immersive;
    document.documentElement.classList.toggle('jf-immersive',immersive);
    const video=document.querySelector('div#videoOsdPage');
    if(video)video.style.visibility=immersive?'hidden':'';
    showToast(immersive?'Immersive mode':'Normal mode')
  },true);
  document.addEventListener('dblclick',event=>{
    if(excluded(event.target))return;
    const target=event.target;
    if(target.closest('.skinHeader,.headerTop,.headerLeft,.headerRight,.MuiToolbar-root,div#videoOsdPage'))host('toggle_maximize')
  },true);
  document.addEventListener('mousedown',event=>{
    if(event.button!==0||excluded(event.target)||!event.target.closest('div#videoOsdPage'))return;
    const startX=event.clientX,startY=event.clientY;
    let moved=false;
    const move=moveEvent=>{
      if(moved||Math.hypot(moveEvent.clientX-startX,moveEvent.clientY-startY)<=3)return;
      moved=true;
      host('drag_window');
      cleanup()
    };
    const cleanup=()=>{
      document.removeEventListener('mousemove',move,true);
      document.removeEventListener('mouseup',cleanup,true)
    };
    document.addEventListener('mousemove',move,true);
    document.addEventListener('mouseup',cleanup,true)
  },true);
  host('page_load',{url:location.href,title:document.title})
})()`

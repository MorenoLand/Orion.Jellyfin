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

const immersiveToggleScript = `(function(){
  const video=document.querySelector('div#videoOsdPage');
  const immersive=!window.__morenoImmersiveMode;
  window.__morenoImmersiveMode=immersive;
  const immersiveHeaderStyle=document.getElementById('moreno-immersive-header-style')||document.createElement('style');
  immersiveHeaderStyle.id='moreno-immersive-header-style';
  immersiveHeaderStyle.textContent='html.moreno-immersive .skinHeader,html.moreno-immersive .headerTop{display:none!important}html.moreno-immersive,html.moreno-immersive *{cursor:none!important}';
  (document.head||document.documentElement).appendChild(immersiveHeaderStyle);
  document.documentElement.classList.toggle('moreno-immersive',immersive);
  const keyboardTypes=['keydown','keypress','keyup'];
  const keyboardBlock=window.__morenoImmersiveKeyBlock||(window.__morenoImmersiveKeyBlock=event=>{
    if(String(event.key||'').toLowerCase()==='x'&&!event.ctrlKey&&!event.altKey&&!event.metaKey&&!event.shiftKey)return;
    event.preventDefault();
    event.stopImmediatePropagation()
  });
  const setKeyboardBlock=active=>keyboardTypes.forEach(type=>[window,document].forEach(target=>active?target.addEventListener(type,keyboardBlock,true):target.removeEventListener(type,keyboardBlock,true)));
  setKeyboardBlock(immersive);
  try{window._wails&&window._wails.invoke&&window._wails.invoke(JSON.stringify({action:'set_immersive',immersive}))}catch(_e){}
  if(window.__morenoResizeLayer)window.__morenoResizeLayer.style.display=immersive?'none':'';
  const blockTypes=['mousemove','mouseover','mouseenter','mouseleave','mousedown','mouseup','click','dblclick','contextmenu','auxclick','wheel','dragstart','drag','dragend','drop','pointerdown','pointerup','pointermove','pointerover','pointerenter','pointerleave','pointercancel','touchstart','touchmove','touchend','touchcancel'];
  const block=window.__morenoImmersiveBlock||(window.__morenoImmersiveBlock=event=>{event.stopImmediatePropagation();event.preventDefault()});
  let overlay=document.getElementById('jf-immersive-overlay');
  if(immersive){
    if(!overlay){
      overlay=document.createElement('div');
      overlay.id='jf-immersive-overlay';
      document.body.appendChild(overlay)
    }
    overlay.style.cssText='position:fixed;inset:0;z-index:2147483647;background:transparent;pointer-events:auto;cursor:none!important;touch-action:none;user-select:none';
    if(video)video.style.display='none';
    blockTypes.forEach(type=>window.addEventListener(type,block,true))
  }else{
    if(overlay)overlay.remove();
    if(video)video.style.display='';
    blockTypes.forEach(type=>window.removeEventListener(type,block,true))
  }
  let toast=document.querySelector('.jf-toast,.gf-toast,#moreno-immersive-toast');
  if(!toast){
    toast=document.createElement('div');
    toast.id='moreno-immersive-toast';
    toast.style.cssText='position:fixed;left:50%;bottom:8%;z-index:2147483647;transform:translateX(-50%);padding:8px 14px;border-radius:4px;background:rgba(0,0,0,.78);color:#fff;font:13px "Segoe UI",sans-serif;pointer-events:none;opacity:0;transition:opacity .15s';
    document.documentElement.appendChild(toast)
  }
  toast.textContent=immersive?'Immersive mode on':'Immersive mode off';
  toast.classList.add('jf-show','gf-show');
  toast.style.opacity='1';
  clearTimeout(window.__morenoImmersiveToastTimer);
  window.__morenoImmersiveToastTimer=setTimeout(()=>{toast.classList.remove('jf-show','gf-show');toast.style.opacity=''},1200)
})()`

const genericInjection = `(function(){
  if(window.__morenoJellyfinInjection||window.__morenoJellyfinInjectionPending)return;
  window.__morenoJellyfinInjectionPending=true;
  const install=()=>{
    if(!window._wails||typeof window._wails.invoke!=='function'){
      setTimeout(install,50);
      return
    }
    if(window.__morenoJellyfinInjection)return;
    window.__morenoJellyfinInjection=true;
    if(typeof window._wails.setResizable==='function')window._wails.setResizable(true);
    const host=(action,data={})=>{
    try{
      if(!window._wails||typeof window._wails.invoke!=='function')return;
      if(action==='drag_window'){window._wails.invoke('wails:drag');return}
      if(action==='resize_window'){window._wails.invoke('wails:resize:'+data.edge);return}
      window._wails.invoke(JSON.stringify(Object.assign({action},data)))
    }catch(_e){}
  };
  const resizeSystem=window._wails&&window._wails.flags&&window._wails.flags.system||{};
  const resizeHandleWidth=Math.max(10,Number(resizeSystem.resizeHandleWidth)||5);
  const resizeHandleHeight=Math.max(10,Number(resizeSystem.resizeHandleHeight)||5);
  const resizeCornerExtra=Math.max(10,Number(window._wails&&window._wails.flags&&window._wails.flags.resizeCornerExtra)||10);
  let resizeEdge='',resizeReady=false,resizeActive=false,resizeCursorStyle;
  const resizeCursor=edge=>edge==='se-resize'||edge==='nw-resize'?'nwse-resize':edge==='sw-resize'||edge==='ne-resize'?'nesw-resize':edge==='w-resize'||edge==='e-resize'?'ew-resize':'ns-resize';
  const setResize=edge=>{
    if(edge){
      if(!resizeCursorStyle){
        resizeCursorStyle=document.createElement('style');
        resizeCursorStyle.id='moreno-resize-cursor';
        (document.head||document.documentElement).appendChild(resizeCursorStyle)
      }
      resizeCursorStyle.textContent='html,body,body *{cursor:'+resizeCursor(edge)+'!important}'
    }else if(resizeCursorStyle)resizeCursorStyle.textContent='';
    resizeEdge=edge||''
  };
  const resizeLayer=document.createElement('div');
  resizeLayer.id='moreno-resize-layer';
  resizeLayer.style.cssText='position:fixed;inset:0;z-index:2147483645;pointer-events:none';
  const resizeZone=(edge,position)=>{
    const zone=document.createElement('div');
    zone.style.cssText='position:fixed;'+position+';z-index:2147483646;pointer-events:auto;cursor:'+resizeCursor(edge);
    zone.addEventListener('mousedown',event=>{
      if(event.button!==0||window.__morenoImmersiveMode)return;
      event.preventDefault();
      event.stopImmediatePropagation();
      setResize(edge);
      host('resize_window',{edge})
    },true);
    resizeLayer.appendChild(zone)
  };
  const cornerWidth=resizeHandleWidth+resizeCornerExtra,cornerHeight=resizeHandleHeight+resizeCornerExtra;
  resizeZone('n-resize','top:0;left:0;right:0;height:'+resizeHandleHeight+'px');
  resizeZone('s-resize','left:0;right:0;bottom:0;height:'+resizeHandleHeight+'px');
  resizeZone('w-resize','top:'+resizeHandleHeight+'px;left:0;bottom:'+resizeHandleHeight+'px;width:'+resizeHandleWidth+'px');
  resizeZone('e-resize','top:'+resizeHandleHeight+'px;right:0;bottom:'+resizeHandleHeight+'px;width:'+resizeHandleWidth+'px');
  resizeZone('nw-resize','top:0;left:0;width:'+cornerWidth+'px;height:'+cornerHeight+'px');
  resizeZone('ne-resize','top:0;right:0;width:'+cornerWidth+'px;height:'+cornerHeight+'px');
  resizeZone('sw-resize','bottom:0;left:0;width:'+cornerWidth+'px;height:'+cornerHeight+'px');
  resizeZone('se-resize','bottom:0;right:0;width:'+cornerWidth+'px;height:'+cornerHeight+'px');
  (document.body||document.documentElement).appendChild(resizeLayer);
  window.__morenoResizeLayer=resizeLayer;
  const updateResize=event=>{
    if(window.__morenoImmersiveMode){setResize();return}
    const rightContentEdge=window.innerWidth-Math.max(0,window.innerWidth-document.documentElement.clientWidth);
    const bottomContentEdge=window.innerHeight-Math.max(0,window.innerHeight-document.documentElement.clientHeight);
    const rightBorder=event.clientX<rightContentEdge&&rightContentEdge-event.clientX<resizeHandleWidth;
    const leftBorder=event.clientX<resizeHandleWidth;
    const topBorder=event.clientY<resizeHandleHeight;
    const bottomBorder=event.clientY<bottomContentEdge&&bottomContentEdge-event.clientY<resizeHandleHeight;
    const rightCorner=event.clientX<rightContentEdge&&rightContentEdge-event.clientX<resizeHandleWidth+resizeCornerExtra;
    const leftCorner=event.clientX<resizeHandleWidth+resizeCornerExtra;
    const topCorner=event.clientY<resizeHandleHeight+resizeCornerExtra;
    const bottomCorner=event.clientY<bottomContentEdge&&bottomContentEdge-event.clientY<resizeHandleHeight+resizeCornerExtra;
    if(!leftCorner&&!topCorner&&!bottomCorner&&!rightCorner){setResize();return}
    if(rightCorner&&bottomCorner)setResize('se-resize');
    else if(leftCorner&&bottomCorner)setResize('sw-resize');
    else if(leftCorner&&topCorner)setResize('nw-resize');
    else if(topCorner&&rightCorner)setResize('ne-resize');
    else if(leftBorder)setResize('w-resize');
    else if(topBorder)setResize('n-resize');
    else if(bottomBorder)setResize('s-resize');
    else if(rightBorder)setResize('e-resize');
    else setResize()
  };
  document.addEventListener('mousemove',event=>{
    if(resizeReady&&resizeEdge){
      resizeActive=true;
      resizeReady=false;
      event.preventDefault();
      event.stopImmediatePropagation();
      host('resize_window',{edge:resizeEdge});
      return
    }
    if(!resizeActive)updateResize(event)
  },true);
  document.addEventListener('mousedown',event=>{
    if(event.button!==0||window.__morenoImmersiveMode)return;
    if(resizeEdge){
      resizeReady=true;
      event.preventDefault();
      event.stopImmediatePropagation()
    }
  },true);
  document.addEventListener('mouseup',event=>{
    if(event.button===0){resizeReady=false;resizeActive=false}
  },true);
  window.addEventListener('resize',()=>{resizeReady=false;resizeActive=false;setResize()});
  const style=document.createElement('style');
  style.id='gf-style';
  style.textContent='html,body{width:100%!important;height:100%!important;overflow:hidden!important;background:#000!important}body{--wails-resize:all}body>*{max-width:none!important}.gf-toast{position:fixed;left:50%;bottom:8%;z-index:2147483647;transform:translateX(-50%);padding:8px 14px;border-radius:4px;background:rgba(0,0,0,.78);color:#fff;font:13px "Segoe UI",sans-serif;pointer-events:none;opacity:0;transition:opacity .15s}.gf-toast.gf-show{opacity:1}';
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
  let immersive=!!window.__morenoImmersiveMode;
  const excluded=target=>target&&target.closest('button,a,input,select,textarea,[contenteditable="true"],iframe,video,audio,canvas');
  document.addEventListener('mousedown',event=>{
    if(event.button!==0||window.__morenoImmersiveMode||excluded(event.target))return;
    const startX=event.clientX,startY=event.clientY;
    let moved=false;
    const move=moveEvent=>{
      if(moved||Math.hypot(moveEvent.clientX-startX,moveEvent.clientY-startY)<=3)return;
      moved=true;
      event.preventDefault();
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
    if(window.__morenoImmersiveMode||excluded(event.target))return;
    host('toggle_maximize')
  },true);
  document.addEventListener('keydown',event=>{
    if(event.target&&event.target.closest('input,textarea,[contenteditable],select'))return;
    if(event.defaultPrevented||event.ctrlKey||event.altKey||event.metaKey||event.shiftKey)return;
  },true);
  host('page_load',{url:location.href,title:document.title})
  };
  install()
})()`

const jellyfinInjection = `(function(){
  if(window.__morenoJellyfinInjection||window.__morenoJellyfinInjectionPending)return;
  window.__morenoJellyfinInjectionPending=true;
  const install=()=>{
    if(!window._wails||typeof window._wails.invoke!=='function'){
      setTimeout(install,50);
      return
    }
    if(window.__morenoJellyfinInjection)return;
    window.__morenoJellyfinInjection=true;
    if(typeof window._wails.setResizable==='function')window._wails.setResizable(true);
    const host=(action,data={})=>{
    try{
      if(!window._wails||typeof window._wails.invoke!=='function')return;
      if(action==='drag_window'){window._wails.invoke('wails:drag');return}
      if(action==='resize_window'){window._wails.invoke('wails:resize:'+data.edge);return}
      window._wails.invoke(JSON.stringify(Object.assign({action},data)))
    }catch(_e){}
  };
  const resizeSystem=window._wails&&window._wails.flags&&window._wails.flags.system||{};
  const resizeHandleWidth=Math.max(10,Number(resizeSystem.resizeHandleWidth)||5);
  const resizeHandleHeight=Math.max(10,Number(resizeSystem.resizeHandleHeight)||5);
  const resizeCornerExtra=Math.max(10,Number(window._wails&&window._wails.flags&&window._wails.flags.resizeCornerExtra)||10);
  let resizeEdge='',resizeReady=false,resizeActive=false,resizeCursorStyle;
  const resizeCursor=edge=>edge==='se-resize'||edge==='nw-resize'?'nwse-resize':edge==='sw-resize'||edge==='ne-resize'?'nesw-resize':edge==='w-resize'||edge==='e-resize'?'ew-resize':'ns-resize';
  const setResize=edge=>{
    if(edge){
      if(!resizeCursorStyle){
        resizeCursorStyle=document.createElement('style');
        resizeCursorStyle.id='moreno-resize-cursor';
        (document.head||document.documentElement).appendChild(resizeCursorStyle)
      }
      resizeCursorStyle.textContent='html,body,body *{cursor:'+resizeCursor(edge)+'!important}'
    }else if(resizeCursorStyle)resizeCursorStyle.textContent='';
    resizeEdge=edge||''
  };
  const resizeLayer=document.createElement('div');
  resizeLayer.id='moreno-resize-layer';
  resizeLayer.style.cssText='position:fixed;inset:0;z-index:2147483645;pointer-events:none';
  const resizeZone=(edge,position)=>{
    const zone=document.createElement('div');
    zone.style.cssText='position:fixed;'+position+';z-index:2147483646;pointer-events:auto;cursor:'+resizeCursor(edge);
    zone.addEventListener('mousedown',event=>{
      if(event.button!==0||window.__morenoImmersiveMode)return;
      event.preventDefault();
      event.stopImmediatePropagation();
      setResize(edge);
      host('resize_window',{edge})
    },true);
    resizeLayer.appendChild(zone)
  };
  const cornerWidth=resizeHandleWidth+resizeCornerExtra,cornerHeight=resizeHandleHeight+resizeCornerExtra;
  resizeZone('n-resize','top:0;left:0;right:0;height:'+resizeHandleHeight+'px');
  resizeZone('s-resize','left:0;right:0;bottom:0;height:'+resizeHandleHeight+'px');
  resizeZone('w-resize','top:'+resizeHandleHeight+'px;left:0;bottom:'+resizeHandleHeight+'px;width:'+resizeHandleWidth+'px');
  resizeZone('e-resize','top:'+resizeHandleHeight+'px;right:0;bottom:'+resizeHandleHeight+'px;width:'+resizeHandleWidth+'px');
  resizeZone('nw-resize','top:0;left:0;width:'+cornerWidth+'px;height:'+cornerHeight+'px');
  resizeZone('ne-resize','top:0;right:0;width:'+cornerWidth+'px;height:'+cornerHeight+'px');
  resizeZone('sw-resize','bottom:0;left:0;width:'+cornerWidth+'px;height:'+cornerHeight+'px');
  resizeZone('se-resize','bottom:0;right:0;width:'+cornerWidth+'px;height:'+cornerHeight+'px');
  (document.body||document.documentElement).appendChild(resizeLayer);
  window.__morenoResizeLayer=resizeLayer;
  const updateResize=event=>{
    if(window.__morenoImmersiveMode){setResize();return}
    const rightContentEdge=window.innerWidth-Math.max(0,window.innerWidth-document.documentElement.clientWidth);
    const bottomContentEdge=window.innerHeight-Math.max(0,window.innerHeight-document.documentElement.clientHeight);
    const rightBorder=event.clientX<rightContentEdge&&rightContentEdge-event.clientX<resizeHandleWidth;
    const leftBorder=event.clientX<resizeHandleWidth;
    const topBorder=event.clientY<resizeHandleHeight;
    const bottomBorder=event.clientY<bottomContentEdge&&bottomContentEdge-event.clientY<resizeHandleHeight;
    const rightCorner=event.clientX<rightContentEdge&&rightContentEdge-event.clientX<resizeHandleWidth+resizeCornerExtra;
    const leftCorner=event.clientX<resizeHandleWidth+resizeCornerExtra;
    const topCorner=event.clientY<resizeHandleHeight+resizeCornerExtra;
    const bottomCorner=event.clientY<bottomContentEdge&&bottomContentEdge-event.clientY<resizeHandleHeight+resizeCornerExtra;
    if(!leftCorner&&!topCorner&&!bottomCorner&&!rightCorner){setResize();return}
    if(rightCorner&&bottomCorner)setResize('se-resize');
    else if(leftCorner&&bottomCorner)setResize('sw-resize');
    else if(leftCorner&&topCorner)setResize('nw-resize');
    else if(topCorner&&rightCorner)setResize('ne-resize');
    else if(leftBorder)setResize('w-resize');
    else if(topBorder)setResize('n-resize');
    else if(bottomBorder)setResize('s-resize');
    else if(rightBorder)setResize('e-resize');
    else setResize()
  };
  document.addEventListener('mousemove',event=>{
    if(resizeReady&&resizeEdge){
      resizeActive=true;
      resizeReady=false;
      event.preventDefault();
      event.stopImmediatePropagation();
      host('resize_window',{edge:resizeEdge});
      return
    }
    if(!resizeActive)updateResize(event)
  },true);
  document.addEventListener('mousedown',event=>{
    if(event.button!==0||window.__morenoImmersiveMode)return;
    if(resizeEdge){
      resizeReady=true;
      event.preventDefault();
      event.stopImmediatePropagation()
    }
  },true);
  document.addEventListener('mouseup',event=>{
    if(event.button===0){resizeReady=false;resizeActive=false}
  },true);
  window.addEventListener('resize',()=>{resizeReady=false;resizeActive=false;setResize()});
  const style=document.createElement('style');
  style.id='jf-style';
  style.textContent='.skinHeader,.headerTop,.headerLeft,.headerRight,.MuiToolbar-root{--wails-draggable:drag;-webkit-app-region:drag;app-region:drag}.skinHeader button,.skinHeader a,.skinHeader input,.skinHeader select,.skinHeader textarea,.headerTop button,.headerTop a,.headerLeft button,.headerLeft a,.headerRight button,.headerRight a,.MuiToolbar-root button,.MuiToolbar-root a,.MuiToolbar-root input,.MuiToolbar-root select{--wails-draggable:no-drag;-webkit-app-region:no-drag;app-region:no-drag}html.moreno-immersive .skinHeader,html.moreno-immersive .headerTop{display:none!important}body{--wails-resize:all;user-select:none!important;-webkit-user-select:none!important}input,select,textarea,[contenteditable="true"]{user-select:text!important;-webkit-user-select:text!important}.jf-toast{position:fixed;left:50%;bottom:8%;z-index:2147483647;transform:translateX(-50%);padding:8px 14px;border-radius:4px;background:rgba(0,0,0,.78);color:#fff;font:13px "Segoe UI",sans-serif;pointer-events:none;opacity:0;transition:opacity .15s}.jf-toast.jf-show{opacity:1}';
  (document.head||document.documentElement).appendChild(style);
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
  let immersive=!!window.__morenoImmersiveMode,overlay=document.getElementById('jf-immersive-overlay');
  document.documentElement.classList.toggle('moreno-immersive',immersive);
  const blockAll=event=>{event.stopImmediatePropagation();event.preventDefault()};
  const blockTypes=['mousemove','mouseover','mouseenter','mousedown','mouseup','click','dblclick','contextmenu','pointerdown','pointerup','pointermove','pointerover','pointerenter'];
  const excluded=target=>target&&target.closest('button,a,input,select,textarea,[contenteditable="true"]');
  document.addEventListener('keydown',event=>{
    if(event.target&&event.target.closest('input,textarea,[contenteditable],select'))return;
    if(event.defaultPrevented||event.ctrlKey||event.altKey||event.metaKey||event.shiftKey)return;
    const key=event.key.toLowerCase();
    if(key==='t'){
      host('toggle_theater');
      event.preventDefault();
      event.stopImmediatePropagation();
      return false
    }
  },true);
  document.addEventListener('dblclick',event=>{
    if(excluded(event.target))return;
    const target=event.target;
    if(!target.closest('.skinHeader,.headerTop,.headerLeft,.headerRight,.MuiToolbar-root,div#videoOsdPage'))return;
    host('toggle_maximize');
    event.preventDefault();
    event.stopImmediatePropagation()
  },true);
  document.addEventListener('mousedown',event=>{
    if(event.button!==0||window.__morenoImmersiveMode||excluded(event.target)||!event.target.closest('.skinHeader,.headerTop,.headerLeft,.headerRight,.MuiToolbar-root'))return;
    const startX=event.clientX,startY=event.clientY;
    const move=moveEvent=>{
      if(Math.hypot(moveEvent.clientX-startX,moveEvent.clientY-startY)<=3)return;
      event.preventDefault();
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
  document.addEventListener('mousedown',event=>{
    if(event.button!==0||window.__morenoImmersiveMode||excluded(event.target)||event.target.closest('.sliderContainer,.volumeSlider,.osdVolumeSlider')||!event.target.closest('div#videoOsdPage'))return;
    const startX=event.clientX,startY=event.clientY;
    let moved=false;
    const move=moveEvent=>{
      if(moved||Math.hypot(moveEvent.clientX-startX,moveEvent.clientY-startY)<=3)return;
      moved=true;
      event.preventDefault();
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
  };
  install()
})()`

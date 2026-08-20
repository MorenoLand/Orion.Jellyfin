import './style.css'

type WailsRuntime = { invoke: (message: string) => unknown }

declare global {
  interface Window {
    _wails?: WailsRuntime
  }
}

const sendHost = (action: string) => {
  window._wails?.invoke(JSON.stringify({ action }))
}

const iframe = document.getElementById('jellyfin') as HTMLIFrameElement | null
const loading = document.getElementById('loading') as HTMLDivElement | null

iframe?.addEventListener('load', () => {
  if (loading) loading.style.display = 'none'
})

iframe?.addEventListener('error', () => {
  if (loading) loading.innerHTML = '<span>Failed to connect to Jellyfin</span>'
})

document.getElementById('win-max')?.addEventListener('click', () => sendHost('toggle_maximize'))
document.getElementById('win-close')?.addEventListener('click', () => sendHost('hide'))

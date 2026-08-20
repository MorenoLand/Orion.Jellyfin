import './styles.css'
import { getCurrentWindow } from '@tauri-apps/api/window'

const appWindow = getCurrentWindow()

const iframe = document.getElementById('jellyfin') as HTMLIFrameElement
const loading = document.getElementById('loading') as HTMLDivElement

iframe?.addEventListener('load', () => { loading.style.display = 'none' })
iframe?.addEventListener('error', () => { loading.innerHTML = '<span>Failed to connect to Jellyfin</span>' })

document.getElementById('win-max')?.addEventListener('click', async () => {
  if (await appWindow.isMaximized()) await appWindow.unmaximize()
  else await appWindow.maximize()
})

document.getElementById('win-close')?.addEventListener('click', () => { appWindow.hide() })

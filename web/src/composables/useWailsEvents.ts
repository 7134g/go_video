import { onUnmounted } from 'vue'

interface ProgressInfo {
  id: number
  name: string
  type: string
  done: number
  total: number
  percent: number
  timespec: number
}

interface BroadcastMessage {
  task_id: number
  message: string
}

export function useWailsEvents(handlers: {
  onProgress?: (data: ProgressInfo) => void
  onProgressSnapshot?: (data: ProgressInfo[]) => void
  onBroadcast?: (data: BroadcastMessage) => void
}) {
  const isWails = typeof window !== 'undefined' && !!(window as any).runtime

  if (isWails) {
    const runtime = (window as any).runtime

    if (handlers.onProgress) {
      runtime.EventsOn('task:progress', handlers.onProgress)
    }
    if (handlers.onBroadcast) {
      runtime.EventsOn('task:broadcast', handlers.onBroadcast)
    }

    onUnmounted(() => {
      if (handlers.onProgress) {
        runtime.EventsOff('task:progress', handlers.onProgress)
      }
      if (handlers.onBroadcast) {
        runtime.EventsOff('task:broadcast', handlers.onBroadcast)
      }
    })
  }

  async function getInitialProgress(): Promise<ProgressInfo[]> {
    if (!isWails) return []
    return (window as any).go.main.App.GetAllProgress()
  }

  return { isWails, getInitialProgress }
}

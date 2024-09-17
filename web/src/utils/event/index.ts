import { reactive } from 'vue'

type EventCallback<T = any> = (payload: T) => void

interface EventBusInterface {
  events: Record<string, EventCallback[]>
  emit<T = any>(event: string, data?: T): void
  on<T = any>(event: string, callback: EventCallback<T>): void
  off<T = any>(event: string, callback: EventCallback<T>): void
}

const EventBus: EventBusInterface = reactive({
  events: {} as Record<string, EventCallback[]>,
  emit<T>(event: string, data?: T): void {
    if (this.events[event]) {
      this.events[event].forEach((callback) => callback(data))
    }
  },
  on<T>(event: string, callback: EventCallback<T> = () => {}): void {
    if (!this.events[event]) {
      this.events[event] = []
    }
    this.events[event].push(callback)
  },
  off<T>(event: string, callback: EventCallback<T>): void {
    if (this.events[event]) {
      const index = this.events[event].indexOf(callback)
      if (index > -1) {
        this.events[event].splice(index, 1)
      }
    }
  }
})

export default EventBus

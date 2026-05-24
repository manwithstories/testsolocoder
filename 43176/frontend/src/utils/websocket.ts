export class WebSocketService {
  private ws: WebSocket | null = null
  private url: string
  private reconnectAttempts = 0
  private maxReconnectAttempts = 5
  private reconnectDelay = 3000

  private onOpenHandlers: (() => void)[] = []
  private onMessageHandlers: ((data: any) => void)[] = []
  private onCloseHandlers: (() => void)[] = []
  private onErrorHandlers: ((error: Event) => void)[] = []

  constructor(orderId: string | number, token: string) {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    this.url = `${protocol}//${window.location.host}/api/v1/ws/${orderId}?token=${token}`
  }

  connect(): void {
    try {
      this.ws = new WebSocket(this.url)

      this.ws.onopen = () => {
        console.log('WebSocket connected')
        this.reconnectAttempts = 0
        this.onOpenHandlers.forEach(handler => handler())
      }

      this.ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          this.onMessageHandlers.forEach(handler => handler(data))
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error)
        }
      }

      this.ws.onclose = () => {
        console.log('WebSocket disconnected')
        this.onCloseHandlers.forEach(handler => handler())
        this.reconnect()
      }

      this.ws.onerror = (error) => {
        console.error('WebSocket error:', error)
        this.onErrorHandlers.forEach(handler => handler(error))
      }
    } catch (error) {
      console.error('Failed to create WebSocket connection:', error)
    }
  }

  private reconnect(): void {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++
      console.log(`Reconnecting... attempt ${this.reconnectAttempts}`)

      setTimeout(() => {
        this.connect()
      }, this.reconnectDelay * this.reconnectAttempts)
    }
  }

  send(data: any): void {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(data))
    }
  }

  onOpen(handler: () => void): void {
    this.onOpenHandlers.push(handler)
  }

  onMessage(handler: (data: any) => void): void {
    this.onMessageHandlers.push(handler)
  }

  onClose(handler: () => void): void {
    this.onCloseHandlers.push(handler)
  }

  onError(handler: (error: Event) => void): void {
    this.onErrorHandlers.push(handler)
  }

  close(): void {
    if (this.ws) {
      this.reconnectAttempts = this.maxReconnectAttempts
      this.ws.close()
    }
  }

  isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN
  }
}

export const useWebSocket = (orderId: string | number, token: string) => {
  const ws = new WebSocketService(orderId, token)
  return ws
}

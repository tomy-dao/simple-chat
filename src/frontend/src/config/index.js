export const config = {
  apiUrl: import.meta.env.VITE_API_BASE_API_URL || 'http://localhost/api/v1',
  socketUrl: import.meta.env.VITE_SOCKET_BASE_URL || 'ws://localhost:8080/chat',
}

console.log(config);
export function connectWebSocket(onMessage) {
    const ws = new WebSocket("ws://localhost:8080/ws");
  
    ws.onopen = () => console.log("WebSocket Connected");
  
    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        onMessage(data); // Calls the callback function passed in `Home.jsx`
      } catch (error) {
        console.error("Invalid WebSocket message format:", error);
      }
    };
  
    ws.onerror = (error) => console.error("WebSocket Error:", error);
    ws.onclose = () => console.log("WebSocket Disconnected");
  
    return ws;
  }
  
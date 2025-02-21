import { create } from "zustand";
import { Howl } from "howler";

// Zustand Store
export const useTrafficStore = create((set, get) => ({
  vehicles: [],
  trafficLights: [],
  emergencies: [],
  congestion: 0,
  avgSpeed: 45,
  isSimulationActive: false,

  // Audio
  sirenSound: new Howl({ src: ["/sounds/siren.mp3"], loop: true }),
  trafficSound: new Howl({ src: ["/sounds/traffic.mp3"], loop: true }),

  // Actions
  updateState: (data) => {
    switch (data.type) {
      case "vehicle_update":
        set((state) => ({
          vehicles: [
            ...state.vehicles.filter((v) => v.id !== data.id),
            {
              id: data.id,
              position: [data.lat, data.lng],
              speed: data.speed,
              type: data.isEmergency ? "emergency" : "normal",
            },
          ],
        }));
        break;

      case "traffic_light":
        set((state) => ({
          trafficLights: state.trafficLights.map((light) =>
            light.id === data.id ? { ...light, state: data.state } : light
          ),
        }));
        break;

      case "emergency_start":
        get().sirenSound.play();
        set((state) => ({
          emergencies: [...state.emergencies, data],
        }));
        break;

      default:
        console.warn("Unhandled event type:", data.type);
    }
  },

  toggleSimulation: () => {
    const isActive = !get().isSimulationActive;
    set({ isSimulationActive: isActive });
    isActive ? get().trafficSound.play() : get().trafficSound.stop();
  },
}));

// WebSocket Connection
export function connectWebSocket() {
  const ws = new WebSocket("ws://localhost:8080/ws");

  ws.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      useTrafficStore.getState().updateState(data);
    } catch (error) {
      console.error("Invalid message format:", error);
    }
  };

  ws.onerror = (error) => console.error("WebSocket Error:", error);
  ws.onclose = () => console.log("WebSocket Disconnected");

  return ws;
}

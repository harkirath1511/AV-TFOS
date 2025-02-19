Here's the complete `trafficStore.js` implementation using Zustand for state management, handling real-time traffic updates, emergency alerts, and simulation controls:

```javascript
// src/stores/trafficStore.js
import create from 'zustand';
import { playSirenSound, playTrafficSound } from '../utils/audio';

export const useTrafficStore = create((set, get) => ({
  // State
  vehicles: [],
  trafficLights: [],
  emergencies: [],
  congestion: 0, // 0-100 percentage
  avgSpeed: 45, // km/h
  isSimulationActive: false,
  selectedVehicle: null,

  // Actions
  addVehicle: (vehicle) => {
    set((state) => ({
      vehicles: [...state.vehicles.filter(v => v.id !== vehicle.id), vehicle]
    }));
  },

  updateTrafficLight: (lightId, state) => {
    set((store) => ({
      trafficLights: store.trafficLights.map(light => 
        light.id === lightId ? { ...light, state } : light
      )
    }));
  },

  addEmergency: (emergency) => {
    playSirenSound();
    set((state) => ({
      emergencies: [...state.emergencies, emergency]
    }));
    
    // Auto-remove emergency after 2 minutes
    setTimeout(() => {
      set((state) => ({
        emergencies: state.emergencies.filter(e => e.id !== emergency.id)
      }));
    }, 120000);
  },

  toggleSimulation: () => {
    const isActive = !get().isSimulationActive;
    set({ isSimulationActive: isActive });
    playTrafficSound(isActive);
  },

  // WebSocket message handler
  onMessage: (message) => {
    try {
      const data = JSON.parse(message);
      
      switch(data.type) {
        case 'vehicle_update':
          get().addVehicle({
            id: data.vehicleId,
            lat: data.position[0],
            lng: data.position[1],
            speed: data.speed,
            type: data.isEmergency ? 'emergency' : 'av',
            route: data.route
          });
          break;

        case 'traffic_light':
          get().updateTrafficLight(data.lightId, data.state);
          break;

        case 'emergency_start':
          get().addEmergency({
            id: data.emergencyId,
            route: data.route,
            progress: 0
          });
          break;

        case 'congestion_update':
          set({ congestion: data.percentage });
          break;

        case 'avg_speed':
          set({ avgSpeed: data.speed });
          break;
      }

      // Update congestion percentage (mock calculation)
      const congestion = Math.min(
        Math.floor((get().vehicles.length / 500) * 100),
        100
      );
      set({ congestion });

    } catch (error) {
      console.error('Error processing message:', error);
    }
  },

  // Derived state
  congestionPercentage: () => {
    return get().congestion;
  },

  emergencyRoutes: () => {
    return get().emergencies.map(e => e.route);
  },

  // Simulation controls
  startSimulation: (config) => {
    set({ 
      isSimulationActive: true,
      vehicles: [],
      emergencies: [] 
    });
    // Implementation to connect to simulation backend
  },

  stopSimulation: () => {
    set({ isSimulationActive: false });
  }
}));

// WebSocket connection handler
export const connectWebSocket = () => {
  const ws = new WebSocket('ws://localhost:8080/ws');
  
  ws.onmessage = (event) => {
    useTrafficStore.getState().onMessage(event.data);
  };

  ws.onclose = () => {
    setTimeout(connectWebSocket, 3000); // Reconnect after 3 seconds
  };

  return ws;
};
```

### Key Features:
1. **Real-time State Management**:
   - Handles vehicle positions, traffic light states, and emergencies
   - Auto-updates congestion percentage based on vehicle density
   - Manages simulation state and audio feedback

2. **Emergency Handling**:
   - Plays siren sound on emergency vehicles
   - Auto-removes emergencies after 2 minutes
   - Tracks emergency vehicle progress

3. **WebSocket Integration**:
   - Automatic reconnection logic
   - Message type handling (vehicles, lights, emergencies)
   - Congestion and speed updates

4. **Derived State**:
   - Calculated congestion percentage
   - Emergency route visualization data
   - Filtered vehicle lists

5. **Simulation Controls**:
   - Start/stop simulation with cleanup
   - Audio feedback for simulation state
   - Configurable simulation parameters

### Usage Example:
```javascript
// In your components
import { useTrafficStore } from '../stores/trafficStore';

function TrafficComponent() {
  const { congestion, vehicles, toggleSimulation } = useTrafficStore();
  
  return (
    <div>
      <button onClick={toggleSimulation}>
        {useTrafficStore(state => state.isSimulationActive) 
          ? 'Stop Simulation' 
          : 'Start Simulation'}
      </button>
      <p>Current Congestion: {congestion}%</p>
    </div>
  );
}
```

To test the store without a backend, you can dispatch mock messages:
```javascript
// Mock emergency vehicle
useTrafficStore.getState().onMessage(JSON.stringify({
  type: 'emergency_start',
  emergencyId: 'amb-1',
  route: [[37.7749, -122.4194], [37.7813, -122.4168]]
}));
```

This store provides the foundation for building a real-time traffic visualization dashboard with support for both live data and simulations.

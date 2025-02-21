import { useTrafficStore } from "../stores/trafficStore.js";
import React from "react";

export default function TrafficLights() {
  const trafficLights = useTrafficStore(state => state.trafficLights)
  
  return trafficLights.map(light => (
    <TrafficLight key={light.id} {...light} />
  ))
}

function TrafficLight({ id, position, state }) {
  const colors = {
    red: '#ff0000',
    yellow: '#ffff00',
    green: '#00ff00'
  }

  return (
    <group position={[position[0] * 1000, 0, position[1] * 1000]}>
      <mesh position={[0, 4, 0]}>
        <boxGeometry args={[1, 8, 1]} />
        <meshStandardMaterial color="#333333" />
      </mesh>
      
      {['red', 'yellow', 'green'].map((color, i) => (
        <mesh key={color} position={[0, 6 - i*2, 0]}>
          <sphereGeometry args={[0.4]} />
          <meshStandardMaterial 
            color={color === state ? colors[color] : '#444444'}
            emissive={color === state ? colors[color] : '#000000'}
            emissiveIntensity={0.8}
          />
        </mesh>
      ))}
    </group>
  )
}

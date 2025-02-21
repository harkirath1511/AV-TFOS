import { useGLTF } from '@react-three/drei'
import { useTrafficStore } from "../stores/trafficStore.js";
import React from 'react';

import { a, useSpring } from '@react-spring/three'

export default function EmergencyVehicles() {
  const emergencies = useTrafficStore(state => state.emergencies)
  
  return emergencies.map(emergency => (
    <EmergencyVehicle key={emergency.id} {...emergency} />
  ))
}

function EmergencyVehicle({ id, route }) {
  const { scene } = useGLTF('/assets/ambulance.glb')
  const [spring] = useSpring(() => ({
    from: { emissive: 0 },
    to: async next => {
      while(true) {
        await next({ emissive: 2 })
        await next({ emissive: 0.5 })
      }
    },
    config: { duration: 500 }
  }), [])

  return (
    <a.group position={[route[0][0] * 1000, 0, route[0][1] * 1000]}>
      <primitive object={scene} scale={0.8} />
      <a.meshStandardMaterial 
        color="#ff2222"
        emissive="#ff4444"
        emissiveIntensity={spring.emissive}
      />
    </a.group>
  )
}

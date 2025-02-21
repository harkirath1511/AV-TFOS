import { Canvas } from '@react-three/fiber'
import { Environment, OrbitControls } from '@react-three/drei'
import Vehicles from './Vehicle.jsx'
import TrafficLights from './TrafficLight'
import EmergencyVehicles from './Emergency'
import Roads from './Roads'
import React from 'react'

export default function CityMap() {
  return (
    <Canvas camera={{ position: [0, 500, 500], fov: 50 }}>
      <ambientLight intensity={0.8} />
      <Environment preset="dawn" />
      
      {/* City Infrastructure */}
      <Roads />
      
      {/* Dynamic Elements */}
      <Vehicles />
      <TrafficLights />
      <EmergencyVehicles />
      
      {/* Camera Controls */}
      <OrbitControls
        enablePan={true}
        enableZoom={true}
        maxPolarAngle={Math.PI/2 - 0.1}
      />
    </Canvas>
  )
}

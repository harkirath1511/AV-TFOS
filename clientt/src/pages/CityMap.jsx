import React, { Suspense } from 'react'
import { Environment } from '@react-three/drei'
import Vehicles from './Vehicle.jsx'
import TrafficLights from './TrafficLight'
import EmergencyVehicles from './Emergency'
import Roads from './Roads'

export default function CityMap() {
  return (
    <group>
      {/* Wrap 3D content in Suspense */}
      <Suspense fallback={null}>
        {/* City Infrastructure */}
        <Roads />
        
        {/* Dynamic Elements */}
        <Vehicles />
        <TrafficLights />
        <EmergencyVehicles />
        
        {/* Environment */}
        <Environment preset="dawn" />
      </Suspense>
    </group>
  )
}

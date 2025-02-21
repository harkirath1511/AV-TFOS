import { useGLTF, useAnimations } from '@react-three/drei'
import { useFrame } from '@react-three/fiber'
import { useTrafficStore } from "../stores/trafficStore.js";
import { a, useSpring } from '@react-spring/three'
import React from 'react';

 function Vehicles() {
  const { nodes, materials, animations } = useGLTF('/assets/vehicle.glb')
  const vehicles = useTrafficStore(state => state.vehicles)

  return vehicles.map(vehicle => (
    <Vehicle key={vehicle.id} {...vehicle} />
  ))
}

function Vehicle({ id, position, speed, type }) {
  const { scene, animations } = useGLTF(type === 'emergency' 
    ? '/assets/ambulance.glb' 
    : '/assets/car.glb'
  )
  const { actions } = useAnimations(animations, scene)
  const [spring] = useSpring(() => ({
    position: [position[0] * 1000, 0, position[1] * 1000],
    rotation: [0, Math.PI, 0],
    config: { mass: 1, tension: 500, friction: 40 }
  }), [position])

  useFrame(() => {
    actions?.Drive?.setEffectiveTimeScale(speed / 50)
  })

  return (
    <a.group {...spring}>
      <primitive 
        object={scene} 
        scale={type === 'emergency' ? 0.8 : 0.5}
      />
      <meshStandardMaterial 
        color={type === 'emergency' ? '#ff2222' : '#44aaff'}
        emissive={type === 'emergency' ? '#ff4444' : '#2266ff'}
        emissiveIntensity={1.2}
      />
    </a.group>
  )
}

export default Vehicles

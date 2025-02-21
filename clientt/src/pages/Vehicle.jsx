import React from 'react'
import { useFrame } from '@react-three/fiber'
import * as THREE from 'three'

function Vehicle() {
  const meshRef = React.useRef()

  useFrame(() => {
    if (meshRef.current) {
      meshRef.current.rotation.y += 0.01
    }
  })

  return (
    <mesh ref={meshRef}>
      {/* Simple car shape using box geometry */}
      <group>
        {/* Car body */}
        <mesh position={[0, 0.5, 0]}>
          <boxGeometry args={[2, 1, 4]} />
          <meshStandardMaterial color="#4287f5" />
        </mesh>
        
        {/* Car roof */}
        <mesh position={[0, 1.25, 0]}>
          <boxGeometry args={[1.5, 0.75, 2]} />
          <meshStandardMaterial color="#2961c4" />
        </mesh>

        {/* Wheels */}
        <mesh position={[-1, 0, 1]}>
          <cylinderGeometry args={[0.4, 0.4, 0.3]} rotation={[Math.PI / 2, 0, 0]} />
          <meshStandardMaterial color="#1a1a1a" />
        </mesh>
        <mesh position={[1, 0, 1]}>
          <cylinderGeometry args={[0.4, 0.4, 0.3]} rotation={[Math.PI / 2, 0, 0]} />
          <meshStandardMaterial color="#1a1a1a" />
        </mesh>
        <mesh position={[-1, 0, -1]}>
          <cylinderGeometry args={[0.4, 0.4, 0.3]} rotation={[Math.PI / 2, 0, 0]} />
          <meshStandardMaterial color="#1a1a1a" />
        </mesh>
        <mesh position={[1, 0, -1]}>
          <cylinderGeometry args={[0.4, 0.4, 0.3]} rotation={[Math.PI / 2, 0, 0]} />
          <meshStandardMaterial color="#1a1a1a" />
        </mesh>
      </group>
    </mesh>
  )
}

export default Vehicle

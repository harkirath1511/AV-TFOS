import { useMemo } from 'react'
import { Line, Plane } from '@react-three/drei'
import React from 'react'

export default function Roads() {
  // Define main road layout (grid pattern)
  const gridConfig = {
    gridSize: [100, 100], // City blocks (100x100 units)
    cellSize: 10,         // Road every 10 units
    roadWidth: 2,         // 2 units wide roads
    roadColor: '#2a2a2a',
    lineColor: '#ffffff'
  }

  // Generate road geometry
  const { horizontal, vertical, lines } = useMemo(() => {
    const roads = {
      horizontal: [],
      vertical: [],
      lines: []
    }

    // Generate grid roads
    for (let i = -gridConfig.gridSize[0]/2; i <= gridConfig.gridSize[0]/2; i += gridConfig.cellSize) {
      roads.horizontal.push(
        <Plane
          key={`h-${i}`}
          args={[gridConfig.gridSize[0], gridConfig.roadWidth]}
          position={[i, 0.1, 0]}
          rotation={[-Math.PI/2, 0, 0]}
        >
          <meshStandardMaterial color={gridConfig.roadColor} />
        </Plane>
      )

      roads.vertical.push(
        <Plane
          key={`v-${i}`}
          args={[gridConfig.roadWidth, gridConfig.gridSize[1]]}
          position={[0, 0.1, i]}
          rotation={[-Math.PI/2, 0, 0]}
        >
          <meshStandardMaterial color={gridConfig.roadColor} />
        </Plane>
      )
    }

    // Generate lane markings
    for (let i = -gridConfig.gridSize[0]/2; i <= gridConfig.gridSize[0]/2; i += gridConfig.cellSize) {
      roads.lines.push(
        <Line
          key={`line-h-${i}`}
          points={[
            [-gridConfig.gridSize[0]/2, 0.15, i],
            [gridConfig.gridSize[0]/2, 0.15, i]
          ]}
          color={gridConfig.lineColor}
          lineWidth={0.25}
          dashed
        />
      )

      roads.lines.push(
        <Line
          key={`line-v-${i}`}
          points={[
            [i, 0.15, -gridConfig.gridSize[1]/2],
            [i, 0.15, gridConfig.gridSize[1]/2]
          ]}
          color={gridConfig.lineColor}
          lineWidth={0.25}
          dashed
        />
      )
    }

    return roads
  }, [])

  return (
    <group>
      {/* Base road grid */}
      <group position={[0, 0, 0]}>
        {horizontal}
        {vertical}
      </group>

      {/* Lane markings */}
      <group>
        {lines}
      </group>

      {/* Main highway (example) */}
      <mesh
        position={[0, 0.1, 0]}
        rotation={[-Math.PI/2, 0, 0]}
      >
        <planeGeometry args={[gridConfig.gridSize[0], 8]} />
        <meshStandardMaterial color="#333333" />
        <Line
          points={[
            [-gridConfig.gridSize[0]/2, 0.15, 0],
            [gridConfig.gridSize[0]/2, 0.15, 0]
          ]}
          color={gridConfig.lineColor}
          lineWidth={0.5}
        />
      </mesh>
    </group>
  )
}

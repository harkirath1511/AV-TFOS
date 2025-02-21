import { useTrafficStore } from '../stores/trafficStore'
import { Play, Pause, AlertTriangle } from "react-feather";
import React from 'react';

 function ControlPanel() {
  const {
    congestion,
    avgSpeed,
    emergencies,
    isSimulationActive,
    toggleSimulation
  } = useTrafficStore()

  return (
    <div className="fixed top-0 left-0 p-6 bg-gray-900/90 backdrop-blur-lg rounded-tr-2xl shadow-xl">
      <h1 className="text-3xl font-bold mb-4 bg-gradient-to-r from-cyan-400 to-blue-500 bg-clip-text text-transparent">
        FlowSync Control Center
      </h1>

      <div className="grid grid-cols-3 gap-4 mb-6">
        <StatCard 
          title="Congestion" 
          value={`${congestion}%`}
          color="orange"
        />
        <StatCard
          title="Avg Speed"
          value={`${avgSpeed} km/h`}
          color="green"
        />
        <StatCard
          title="Emergencies"
          value={emergencies.length}
          color="red"
          icon={<AlertTriangle size={18} />}
        />
      </div>

      <div className="space-y-4">
        <button 
          onClick={toggleSimulation}
          className={`w-full flex items-center justify-center gap-2 p-3 rounded-lg transition-all ${
            isSimulationActive 
              ? 'bg-red-500 hover:bg-red-600' 
              : 'bg-cyan-500 hover:bg-cyan-600'
          }`}
        >
          {isSimulationActive ? <Stop size={20} /> : <Play size={20} />}
          {isSimulationActive ? 'Stop Simulation' : 'Start Simulation'}
        </button>

        <div className="p-4 bg-gray-800 rounded-lg">
          <h3 className="text-sm font-semibold mb-2 text-gray-400">
            Emergency Control
          </h3>
          <button
            className="w-full bg-red-500/20 hover:bg-red-500/30 p-2 rounded-md text-red-400 flex items-center gap-2 justify-center transition-colors"
          >
            <AlertTriangle size={16} />
            Trigger Emergency Route
          </button>
        </div>

        <div className="p-4 bg-gray-800 rounded-lg">
          <h3 className="text-sm font-semibold mb-2 text-gray-400">
            Traffic Light Control
          </h3>
          <div className="grid grid-cols-2 gap-2">
            <button className="traffic-light-button bg-red-500/20 hover:bg-red-500/30">
              Stop All
            </button>
            <button className="traffic-light-button bg-green-500/20 hover:bg-green-500/30">
              Clear All
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}

function StatCard({ title, value, color, icon }) {
  return (
    <div className={`p-4 rounded-lg bg-${color}-500/10 border border-${color}-500/20`}>
      <div className="flex items-center justify-between">
        <span className={`text-sm text-${color}-400`}>{title}</span>
        {icon && React.cloneElement(icon, { className: `text-${color}-400` })}
      </div>
      <div className={`text-2xl font-bold mt-2 text-${color}-400`}>
        {value}
      </div>
    </div>
  )
}

export default ControlPanel

import { useEffect, useState } from "react";
import { Canvas } from "@react-three/fiber";
import { OrbitControls } from "@react-three/drei";
import CityMap from "../pages/CityMap";
import Vehicle from "../pages/Vehicle";
import Emergency from "../pages/Emergency";
import TrafficLight from "../pages/TrafficLight";
import ControlPanel from "../pages/ControlPanel";
import { useTrafficStore } from "../stores/trafficStore.js";
import { connectWebSocket } from "../utils/websocket";
import React from "react";

function Home() {
  const { setTrafficData } = useTrafficStore();
  const [socket, setSocket] = useState(null);

  useEffect(() => {
    const ws = connectWebSocket((data) => setTrafficData(data));
    setSocket(ws);
    return () => ws.close();
  }, [setTrafficData]);

  return (
    <div className="h-screen w-screen flex flex-col">
      <ControlPanel socket={socket} />
      <Canvas className="flex-grow">
        <ambientLight intensity={0.5} />
        <OrbitControls />
        <CityMap />
        <Vehicle />
        <Emergency />
        <TrafficLight />
      </Canvas>
    </div>
  );
}

export default Home;

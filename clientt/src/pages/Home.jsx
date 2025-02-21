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

function Scene() {
  return (
    <>
      <ambientLight intensity={0.5} />
      <OrbitControls />
      <mesh>
        <CityMap />
      </mesh>
      <mesh>
        <Vehicle />
      </mesh>
      <mesh>
        <Emergency />
      </mesh>
      <mesh>
        <TrafficLight />
      </mesh>
    </>
  );
}

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
      <div className="flex-grow">
<Canvas
  camera={{ position: [10, 10, 10], fov: 75 }}
  gl={{ antialias: true }}
>
  <Scene />
</Canvas>
      </div>
    </div>
  );
}

export default Home;

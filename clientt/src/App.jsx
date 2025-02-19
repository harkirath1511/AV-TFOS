import React from "react"
import { BrowserRouter as Router , Routes , Route } from "react-router-dom"
import CityMap from "./components/CityMap"
import Compiler from "./components/Compiler"

function App() {

  return (
    <Router>
      <Routes>
        <Route path="/" element={<Compiler/>}>
          <Route path="" element={<CityMap/>}></Route>
        </Route>
      </Routes>
    </Router>
  )
}

export default App

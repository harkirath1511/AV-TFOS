import React from 'react'
import './App.css'
import {BrowserRouter as Router , Routes ,  Route} from 'react-router-dom'
import Compiler from './components/Compiler'


function App() {


  return (
    <Router>
      <Routes>
        <Route path='/' element={<Compiler/>}>
          <Route path='' ></Route>
        </Route>
      </Routes>
    </Router>
  )
}

export default App

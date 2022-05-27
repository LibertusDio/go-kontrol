import React from 'react';
import {BrowserRouter as Router, Routes, Route} from "react-router-dom"
import {Login} from "./features/login/login";
import {Home} from "./features/home/home";
import './App.css';
import {Header} from "./features/header/header";

function App() {
    return (
        <div className="container">
            <Header/>
            <Router>
                <Routes>
                    <Route element={<Home/>} path="/"/>
                    <Route element={<Login/>} path="/login"/>
                    {/*<Route exact component={Signup} path="/signup"/>*/}
                    {/*<PrivateRoute exact component={Dashboard} path="/"/>*/}
                    {/*<Route path="*" element={<NotFound/>}/>*/}
                </Routes>
            </Router>
        </div>
    );
}

export default App;

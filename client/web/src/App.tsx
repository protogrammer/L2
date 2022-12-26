import React from 'react';
import {Route, Routes} from "react-router-dom"
import './App.css';
import {FeedScreen} from "./FeedScreen"

const App: React.FC = () => {
    console.log("App");
    return (
        <Routes>
            <Route path="/" element={<FeedScreen/>}/>
        </Routes>
    );
}

export default App;

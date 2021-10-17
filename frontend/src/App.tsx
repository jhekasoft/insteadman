import React, { useState } from 'react';
import logo from './logo.svg';
import './App.css';

function App() {
  const [games, setGames] = useState("");
  (window as any).backend.games().then((games: string) => setGames(games))

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Edit <code>src/App.tsx</code> and save to reload.
        </p>
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn React
        </a>
        <p>{games}</p>
      </header>
    </div>
  );
}

export default App;

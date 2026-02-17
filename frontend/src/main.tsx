import React, { useState } from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'
import Temp from './temp'

function Root() {
  const [started, setStarted] = useState(false)

  if (!started) {
    return (
      <div>
        <Temp />
        <button
          onClick={() => setStarted(true)}
          style={{
            position: 'fixed',
            bottom: '20px',
            right: '20px',
            padding: '10px 20px',
            fontSize: '16px',
            cursor: 'pointer'
          }}
        >
          进入应用
        </button>
      </div>
    )
  }

  return <App />
}

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <Root />
  </React.StrictMode>,
)

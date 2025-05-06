// frontend/src/App.tsx
import { useEffect, useState } from "react";

function App() {
  const [message, setMessage] = useState<string>("");

  useEffect(() => {
    fetch("/api/hello")
      .then(res => {
        if (!res.ok) throw new Error(`HTTP error ${res.status}`);
        return res.json();
      })
      .then(data => setMessage(data.text))
      .catch(err => console.error("Fetch error:", err));
  }, []);

  return (
    <div style={{ padding: 20 }}>
      <h1>React + Go Integration</h1>
      <p>Pesan dari backend: <strong>{message}</strong></p>
    </div>
  );
}

export default App;
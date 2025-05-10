import { useState } from 'react';
import './App.css';

function App() {
  const [target, setTarget] = useState('');
  const [result, setResult] = useState<string[] | null>(null);
  const [found, setFound] = useState<boolean | null>(null);
  const [loading, setLoading] = useState(false);
  const [method, setMethod] = useState('');
  const [numberRecipe, setNumberRecipe] = useState('');
  const findCombination = async () => {
    if (!target) return;

    setLoading(true);
    setResult(null);
    setFound(null);

    try {
      const res = await fetch(`http://localhost:8080/find?target=${target}&method=${method}&numberRecipe=${numberRecipe}`);
      const data = await res.json();
      setResult(data.steps);
      setFound(data.found);
    } catch (error) {
      console.error('❌ Error:', error);
      setFound(false);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="App">
      <h1>🔍 Cari Kombinasi Elemen</h1>
      <input
        value={target}
        onChange={(e) => setTarget(e.target.value)}
        placeholder="Contoh: human"
      />
      <input 
        value={method}
        onChange={(e) => setMethod(e.target.value)}
        placeholder="Contoh: bfs"
      />
      <input
        value={numberRecipe}
        onChange={(e) => setNumberRecipe(e.target.value)}
        placeholder="Contoh: 3"
      />
      <button onClick={findCombination}>Cari</button>
      {loading && <p>⏳ Mencari...</p>}
      {found === false && <p>❌ Tidak ditemukan</p>}
      {found && result && (
        <div>
          <h2>✅ Ditemukan!</h2>
          <ol>
            {result.map((step, i) => (
              <li key={i}>🧪 {step}</li>
            ))}
          </ol>
        </div>
      )}
    </div>
  );
}

export default App;

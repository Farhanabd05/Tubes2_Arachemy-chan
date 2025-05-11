import { useState } from 'react';
import './App.css';

interface SingleResult {
  found: boolean;
  steps: string[];
}

type PathObject = {[key: string] : string[]};
type MultipleResult = PathObject[];

function App() {
  const [target, setTarget] = useState('');
  const [result, setResult] = useState<SingleResult | MultipleResult | null>(null);
  const [isMultiple, setIsMultiple] = useState(false);
  const [loading, setLoading] = useState(false);
  const [method, setMethod] = useState('');
  const [numberRecipe, setNumberRecipe] = useState('');
  const findCombination = async () => {
    if (!target) return;

    setLoading(true);
    setResult(null);
    setIsMultiple(false);

    try {
      const res = await fetch(`http://localhost:8080/find?target=${target}&method=${method}&numberRecipe=${numberRecipe}`);
      const data = await res.json();

      if (Array.isArray(data)) {
        setResult(data);
        setIsMultiple(true);
      }else{
        setResult(data);
        setIsMultiple(false);
      }
    } catch (error) {
      console.error('❌ Error:', error);
      setResult({ found: false, steps: [] });
      setIsMultiple(false);
    } finally {
      setLoading(false);
    }
  };

  const isFound = () => {
    if (!result) return false;
    
    if (isMultiple) {
      return (result as MultipleResult).length > 0;
    } else {
      return (result as SingleResult).found;
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
      {!loading && result && !isFound() && <p>❌ Tidak Ditemukan</p>}
      {!loading && result && isFound() && (
        <div>
          <h2>✅ Ditemukan!</h2>
          {!isMultiple && (
            <ul>
              {(result as SingleResult).steps.map((step, i) => (
                <li key={i}>{i}. 🧪 {step}</li>
              ))}
            </ul>
          )}

          {isMultiple && (
            <div>
              {(result as MultipleResult).map((pathObj, i) => {
                const pathName = Object.keys(pathObj)[0];

                const steps = pathObj[pathName];

                return (
                  <div key={i}>
                    <h3>{pathName}</h3>
                    <ul>
                      {steps.map((steps, j) => (
                        <li key={j}>{j}. 🧪 {steps}</li>
                      ))}
                    </ul>
                  </div>
                )
              })}
            </div>
          )}
        </div>
      )}
    </div>
  );
}

export default App;

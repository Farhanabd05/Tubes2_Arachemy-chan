import { useState, useEffect } from 'react';
import axios from 'axios';
import './App.css';

// Define types for the recipe data
interface Recipe {
  Element: string;
  Ingredient1: string;
  Ingredient2: string;
  Type: number;
}

interface ScrapeResponse {
  data: Recipe[];
}

interface SingleResult {
  found: boolean;
  steps: string[];

}

type PathObject = { [key: string]: string[] };
type MultipleResult = PathObject[];

function App() {
  // State for scraping status
  const [scrapingStatus, setScrapingStatus] = useState<'idle' | 'loading' | 'success' | 'error'>('idle');
  const [recipes, setRecipes] = useState<Recipe[]>([]);
  const [target, setTarget] = useState('');
  const [result, setResult] = useState<SingleResult | MultipleResult | null>(null);
  const [isMultiple, setIsMultiple] = useState(false);
  const [loading, setLoading] = useState(false);
  const [method, setMethod] = useState('');
  const [numberRecipe, setNumberRecipe] = useState('');
  const [runtime, setRuntime] = useState('');
  const [nodesVisited, setNodesVisited] = useState<number | null>(null);
  // New state for bidirectional toggle
  const [bidirectional, setBidirectional] = useState(false);

  // Tambahkan useEffect untuk handle reset bidirectional
  useEffect(() => {
    if (numberRecipe !== '1') {
      setBidirectional(false);
    }
  }, [numberRecipe]);
  
  const findCombination = async () => {
    if (!target) return;

    setLoading(true);
    setResult(null);
    setIsMultiple(false);

    try {
      // Buat objek URLSearchParams
      const params = new URLSearchParams({
        target: target,
        method: method,
        numberRecipe: numberRecipe,
      });

      // Tambahkan parameter bidirectional jika memenuhi syarat
      if (numberRecipe === '1' && bidirectional) {
        params.append('bidirectional', 'true');
      }
      
      const res = await fetch(
        `${import.meta.env.VITE_BACKEND_URL}/find?${params.toString()}`
      );
      const data = await res.json();

      if (Array.isArray(data)) {
        setResult(data);
        setIsMultiple(true);
      } else {
        setResult(data);
        setIsMultiple(false);
      }
      if (data.runtime !== undefined) setRuntime(data.runtime);
      if (data.nodesVisited !== undefined) setNodesVisited(data.nodesVisited);
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

  
  // Function to scrape data when component mounts
  useEffect(() => {
    const fetchData = async () => {
      try {
        setScrapingStatus('loading');
        console.log('Scraping data from API...');
        console.log("Backend URL:", import.meta.env.VITE_BACKEND_URL);
        const response = await axios.get<ScrapeResponse>(`${import.meta.env.VITE_BACKEND_URL}/scrape`);
        setRecipes(response.data.data);
        setScrapingStatus('success');
        console.log(`Successfully scrape`)
        // display recipes json
        console.log(response.data.data);
      } catch (error) {
        console.error('Error during scraping:', error);
        setScrapingStatus('error');
      }
    };

    fetchData();
  }, []); // Empty dependency array means this runs once on mount

  return (
    <div className="App">
      <header>
        <h1>Little Alchemy 2 Path Finder</h1>
      </header>
      
      {/* Scraping status indicator */}
      {scrapingStatus === 'loading' && <p>⏳ Mengambil data resep...</p>}
      {scrapingStatus === 'error' && 
        <p>❌ Error mengambil data. Silakan refresh halaman untuk mencoba lagi.</p>
      }
      {scrapingStatus === 'success' && 
        <p>✅ Berhasil mengambil resep!</p>
      }
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
      <div className="form-group">
        {numberRecipe === '1' && (
          <label>
            <input
              type="checkbox"
              checked={bidirectional}
              onChange={(e) => setBidirectional(e.target.checked)}
            />
            Gunakan Bidirectional Search
          </label>
        )}
      </div>
      <button onClick={findCombination} disabled={loading || !target || scrapingStatus !== 'success'}>Cari</button>
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
            <div className="multiple-results">
              {(result as MultipleResult).map((pathObj, i) => {
                const pathName = Object.keys(pathObj).find(key => key.startsWith('Path'));

                var steps
                if (pathName !== undefined) {
                    steps = pathObj[pathName];
                    // ...
                  } else {
                    // Handle the case where pathName is undefined
                  }
                  const runtime = pathObj['Runtime']?.[0] || '';
                  const nodesVisited = pathObj['NodesVisited']?.[0] || '';

                return (
                  <div key={i}>
                    <h3>{pathName}</h3>
                    <ul>
                      {steps?.map((steps, j) => (
                        <li key={j}>{j}. 🧪 {steps}</li>
                      ))}
                    </ul>
                    {runtime && (
                      <p><strong>Runtime:</strong> {runtime}</p>
                    )}
                    {nodesVisited && (
                      <p><strong>Nodes Visited:</strong> {nodesVisited}</p>
                    )}
                  </div>
                )
              })}
            </div>
          )}
          {runtime && (!isMultiple) && (
            <div>
              <strong>Runtime:</strong> {runtime} ns
            </div>
          )}
          {nodesVisited !== null && (!isMultiple) && (
            <div>
              <strong>Nodes Visited:</strong> {nodesVisited}
            </div>
          )}
        </div>
      )}
    </div>
  );
}

export default App;

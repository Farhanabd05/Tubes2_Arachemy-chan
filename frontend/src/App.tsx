import { useState, useEffect } from 'react';
import axios from 'axios';
import './App.css';
import { SearchForm } from './components/SearchForm';
import { SingleResultDisplay } from './components/SingleResultDisplay';
import { MultipleResultDisplay } from './components/MultipleResultDisplay';
import { StatsDisplay } from './components/StatsDisplay';
import TreeComponent from './components/TreeComponent';

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
  runtime?: string; 
  nodesVisited?: number | null;
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

  const findCombination = async () => {
    if (!target) return;
    setLoading(true);
    setResult(null);
    setIsMultiple(false);

    try {
      const res = await fetch(`http://localhost:8080/find?target=${target}&method=${method}&numberRecipe=${numberRecipe}`);
      const data = await res.json();
      setResult(data);
      setIsMultiple(Array.isArray(data));

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
        
        const response = await axios.get<ScrapeResponse>('http://localhost:8080/scrape');
        setRecipes(response.data.data);
        setScrapingStatus('success');
        console.log(`Successfully scraped ${response.data.data.length} recipes`);
      } catch (error) {
        console.error('Error during scraping:', error);
        setScrapingStatus('error');
      }
    };

    fetchData();
  }, []); // Empty dependency array means this runs once on mount

  return (
    <div className="App">
      <div className="vertical-stack">
        <header>
          <h1>Little Alchemy 2 Path Finder</h1>
        </header>
        
        {/* Scraping status indicator */}
        {scrapingStatus === 'loading' && <p>⏳ Mengambil data resep...</p>}
        {scrapingStatus === 'error' && 
          <p>❌ Error mengambil data. Silakan refresh halaman untuk mencoba lagi.</p>
        }
        {scrapingStatus === 'success' && 
          <p>✅ Berhasil mengambil {recipes.length} resep!</p>
        }
        <SearchForm
          target={target}
          setTarget={setTarget}
          method={method}
          setMethod={setMethod}
          numberRecipe={numberRecipe}
          setNumberRecipe={setNumberRecipe}
          onSearch={findCombination}
        />
        {loading && <p>⏳ Mencari...</p>}
        {!loading && result && !isFound() && <p>❌ Tidak Ditemukan</p>}
        {!loading && result && isFound() && (
          <div>
            <h2>✅ Ditemukan!</h2>
            {!isMultiple && <SingleResultDisplay result={result as SingleResult} />}
            {isMultiple && <MultipleResultDisplay results={result as MultipleResult} />}
          </div>
        )}
      </div>
    </div>
  );
}

export default App;
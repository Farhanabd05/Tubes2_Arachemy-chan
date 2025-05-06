// File: src/App.tsx
import React, { useState } from 'react';
import { RecipeInput } from './components/RecipeInput';
import { PathSearch } from './components/PathSearch';
import { GraphVisualization } from './components/GraphVisualization';

function App() {
  const [highlightPath, setHighlightPath] = useState<string[]>([]);
  const [refreshKey, setRefreshKey] = useState<number>(0);

  const handleRecipesUploaded = () => {
    setHighlightPath([]);
    setRefreshKey(prev => prev + 1);
  };

  const handlePathFound = (path: string[]) => {
    setHighlightPath(path);
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <header className="mb-8 text-center">
        <h1 className="text-3xl font-bold mb-2">
          Graph Data Structure & BFS Algorithm
        </h1>
        <p className="text-gray-600">
          Visualize and search connections between elements using graph and BFS
        </p>
      </header>
      
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
        <RecipeInput onRecipesUploaded={handleRecipesUploaded} />
        <PathSearch onPathFound={handlePathFound} />
      </div>
      
      <div key={refreshKey}>
        <GraphVisualization highlightPath={highlightPath} />
      </div>
    </div>
  );
}

export default App;


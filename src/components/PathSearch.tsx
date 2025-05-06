// File: src/components/PathSearch.tsx
import React, { useState } from 'react';
import { apiService } from '../api';
import type { SearchResult } from '../types';

interface PathSearchProps {
  onPathFound: (path: string[]) => void;
}

export const PathSearch: React.FC<PathSearchProps> = ({ onPathFound }) => {
  const [startNode, setStartNode] = useState<string>('');
  const [targetNode, setTargetNode] = useState<string>('');
  const [result, setResult] = useState<SearchResult | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setResult(null);
    
    if (!startNode || !targetNode) {
      setError('Both start and target nodes are required');
      return;
    }
    
    try {
      setLoading(true);
      const searchResult = await apiService.searchPath(startNode, targetNode);
      setResult(searchResult);
      
      if (searchResult.found && searchResult.path) {
        onPathFound(searchResult.path);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error searching for path');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="mb-6 p-4 border rounded-lg bg-gray-50">
      <h2 className="text-xl font-bold mb-4">Search Path (BFS)</h2>
      <form onSubmit={handleSearch}>
        <div className="grid grid-cols-2 gap-4 mb-4">
          <div>
            <label htmlFor="startNode" className="block mb-2 font-medium">
              Start Node:
            </label>
            <input
              id="startNode"
              type="text"
              className="w-full p-2 border rounded"
              value={startNode}
              onChange={(e) => setStartNode(e.target.value)}
            />
          </div>
          <div>
            <label htmlFor="targetNode" className="block mb-2 font-medium">
              Target Node:
            </label>
            <input
              id="targetNode"
              type="text"
              className="w-full p-2 border rounded"
              value={targetNode}
              onChange={(e) => setTargetNode(e.target.value)}
            />
          </div>
        </div>
        <button
          type="submit"
          className="px-4 py-2 bg-blue-600 text-white font-medium rounded hover:bg-blue-700"
          disabled={loading}
        >
          {loading ? 'Searching...' : 'Search Path'}
        </button>
      </form>
      
      {error && <div className="mt-4 text-red-500">{error}</div>}
      
      {result && (
        <div className="mt-4">
          {result.found ? (
            <div>
              <p className="font-medium text-green-600">Path found!</p>
              <p className="mt-2">Path: {result.path?.join(' → ')}</p>
            </div>
          ) : (
            <p className="font-medium text-amber-600">No path found between these nodes.</p>
          )}
        </div>
      )}
    </div>
  );
};

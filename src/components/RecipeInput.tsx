// File: src/components/RecipeInput.tsx
import React, { useState } from 'react';
import type { Recipe } from '../types';
import { apiService } from '../api';

interface RecipeInputProps {
  onRecipesUploaded: () => void;
}

export const RecipeInput: React.FC<RecipeInputProps> = ({ onRecipesUploaded }) => {
  const [jsonInput, setJsonInput] = useState<string>('');
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    
    try {
      const recipes: Recipe[] = JSON.parse(jsonInput);
      
      if (!Array.isArray(recipes)) {
        throw new Error('Input must be an array of recipes');
      }
      
      // Validate recipe format
      recipes.forEach((recipe, index) => {
        if (!recipe.input || !Array.isArray(recipe.input)) {
          throw new Error(`Recipe at index ${index} has invalid input format`);
        }
        if (!recipe.output || typeof recipe.output !== 'string') {
          throw new Error(`Recipe at index ${index} has invalid output format`);
        }
      });
      
      await apiService.uploadRecipes(recipes);
      onRecipesUploaded();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Invalid JSON format');
    }
  };

  return (
    <div className="mb-6 p-4 border rounded-lg bg-gray-50">
      <h2 className="text-xl font-bold mb-4">Upload Recipes</h2>
      <form onSubmit={handleSubmit}>
        <div className="mb-4">
          <label htmlFor="jsonInput" className="block mb-2 font-medium">
            Enter Recipe JSON:
          </label>
          <textarea
            id="jsonInput"
            className="w-full h-48 p-2 border rounded font-mono"
            value={jsonInput}
            onChange={(e) => setJsonInput(e.target.value)}
            placeholder='[{"input": ["A", "B"], "output": "C"}, {"input": ["C", "D"], "output": "E"}]'
          />
        </div>
        {error && <div className="text-red-500 mb-4">{error}</div>}
        <button
          type="submit"
          className="px-4 py-2 bg-blue-600 text-white font-medium rounded hover:bg-blue-700"
        >
          Upload Recipes
        </button>
      </form>
    </div>
  );
};

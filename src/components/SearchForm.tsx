import { useMemo, useState } from 'react';
import Select from 'react-select';
import Toggle from '../components/Toggle';

interface SearchParams {
  algorithm: 'BFS' | 'DFS';
  searchMode: 'shortest' | 'multiple';
  maxRecipes?: number;
  targetElement: string;
}

const elements = [
  'brick',
  'stone',
  'mud',
  'sand',
  'air',
  'water',
  'fire',
  'heat',
  'pressure',
  'wind',
  'small',
  'big',
  'motion',
];

const SearchForm = () => {
  const [algorithm, setAlgorithm] = useState<'BFS' | 'DFS'>('BFS');
  const [searchMode, setSearchMode] = useState<'shortest' | 'multiple'>('shortest');
  const [maxRecipes, setMaxRecipes] = useState(5);
  const [selectedElement, setSelectedElement] = useState<string | null>(null);

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    const params: SearchParams = {
      algorithm,
      searchMode,
      targetElement: selectedElement || '',
      ...(searchMode === 'multiple' && { maxRecipes }),
    };

    try {
      const response = await fetch('/api/search', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(params),
      });
      const data = await response.json();
      // Handle response data untuk visualisasi
      console.log(data);
    } catch (error) {
      console.error('Search failed:', error);
    }
  };

  const formElements = useMemo(
    () => (
      <form onSubmit={handleSearch} className="space-y-4 p-6 bg-white rounded-lg shadow-md">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Target Element
          </label>
          <Select
            options={elements.map((element) => ({ label: element, value: element }))}
            onChange={(selected) => setSelectedElement(selected!.value)}
            placeholder="Pilih elemen..."
            className="react-select-container"
            classNamePrefix="react-select"
          />
        </div>

        <div className="flex space-x-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Algoritma
            </label>
            <Toggle
              options={[
                { label: 'BFS', value: 'BFS' },
                { label: 'DFS', value: 'DFS' },
              ]}
              selected={algorithm}
              onToggle={(value: 'BFS' | 'DFS') => setAlgorithm(value)}
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Mode Pencarian
            </label>
            <Toggle
              options={[
                { label: 'Terpendek', value: 'shortest' },
                { label: 'Multiple', value: 'multiple' },
              ]}
              selected={searchMode}
              onToggle={(value: 'shortest' | 'multiple') => setSearchMode(value)}
            />
          </div>
        </div>

        {searchMode === 'multiple' && (
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Maksimal Recipe
            </label>
            <input
              type="number"
              min="1"
              value={maxRecipes}
              onChange={(e) => setMaxRecipes(Number(e.target.value))}
              className="mt-1 block w-full rounded-md border-gray-300 shadow-sm p-2 border"
            />
          </div>
        )}

        <button
          type="submit"
          className="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 transition-colors"
        >
          Cari Recipe
        </button>
      </form>
    ),
    [algorithm, searchMode, maxRecipes, selectedElement],
  );

  return formElements;
};

export default SearchForm;


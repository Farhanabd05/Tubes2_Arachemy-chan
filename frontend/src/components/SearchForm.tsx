import Select from 'react-select';

interface SearchFormProps {
  target: string;
  method: string;
  numberRecipe: string;
  bidirectional: boolean;
  setTarget: (val: string) => void;
  setMethod: (val: string) => void;
  setNumberRecipe: (val: string) => void;
  setBidirectional: (val: boolean) => void;
  onSearch: () => void;
}

const customSelectStyles = {
  control: (base: any) => ({
    ...base,
    backgroundColor: '#2b2b2b',
    color: 'white',
    borderColor: '#555',
  }),
  singleValue: (base: any) => ({
    ...base,
    color: 'white',
  }),
  menu: (base: any) => ({
    ...base,
    backgroundColor: '#2b2b2b',
    color: 'white',
  }),
  option: (base: any, { isFocused, isSelected }: any) => ({
    ...base,
    backgroundColor: isSelected
      ? '#555'
      : isFocused
      ? '#444'
      : '#2b2b2b',
    color: 'white',
    cursor: 'pointer',
  }),
  input: (base: any) => ({
    ...base,
    color: 'white',
  }),
};

const methodOptions = [
  { value: 'bfs', label: 'BFS' },
  { value: 'dfs', label: 'DFS' },
];

export const SearchForm: React.FC<SearchFormProps> = ({
  target,
  method,
  numberRecipe,
  bidirectional,
  setTarget,
  setMethod,
  setNumberRecipe,
  setBidirectional,
  onSearch,
}) => (
  <div className="search-form">
    <input value={target} onChange={(e) => setTarget(e.target.value)} placeholder="Contoh: human" />
    <Select
      className="react-select-container"
      classNamePrefix="react-select"
      styles={customSelectStyles}
      options={methodOptions}
      value={methodOptions.find((opt) => opt.value === method)}
      onChange={(selected) => {
        if (selected) setMethod(selected.value);
      }}
      placeholder="Pilih metode..."
    />
    <input
      type="number"
      value={numberRecipe}
      onChange={(e) => setNumberRecipe(e.target.value)}
      placeholder="Jumlah"
      min={1}
      step={1}
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
    <button 
      onClick={onSearch} 
      style={{ backgroundColor: '#4CAF50', color: '#F5F5F5' }}
    >
      Cari
    </button>
  </div>
);
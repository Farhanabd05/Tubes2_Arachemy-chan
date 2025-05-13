import Select from 'react-select';

interface SearchFormProps {
  target: string;
  method: string;
  numberRecipe: string;
  setTarget: (val: string) => void;
  setMethod: (val: string) => void;
  setNumberRecipe: (val: string) => void;
  onSearch: () => void;
}

const customSelectStyles = {
  control: (base) => ({
    ...base,
    backgroundColor: '#2b2b2b',
    color: 'white',
    borderColor: '#555',
  }),
  singleValue: (base) => ({
    ...base,
    color: 'white',
  }),
  menu: (base) => ({
    ...base,
    backgroundColor: '#2b2b2b',
    color: 'white',
  }),
  option: (base, { isFocused, isSelected }) => ({
    ...base,
    backgroundColor: isSelected
      ? '#555'
      : isFocused
      ? '#444'
      : '#2b2b2b',
    color: 'white',
    cursor: 'pointer',
  }),
  input: (base) => ({
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
  setTarget,
  setMethod,
  setNumberRecipe,
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
    <button onClick={onSearch}>Cari</button>
  </div>
);
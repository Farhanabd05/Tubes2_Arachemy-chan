interface SearchFormProps {
  target: string;
  method: string;
  numberRecipe: string;
  setTarget: (val: string) => void;
  setMethod: (val: string) => void;
  setNumberRecipe: (val: string) => void;
  onSearch: () => void;
}

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
    <input value={method} onChange={(e) => setMethod(e.target.value)} placeholder="Contoh: bfs" />
    <input value={numberRecipe} onChange={(e) => setNumberRecipe(e.target.value)} placeholder="Contoh: 3" />
    <button onClick={onSearch}>Cari</button>
  </div>
);
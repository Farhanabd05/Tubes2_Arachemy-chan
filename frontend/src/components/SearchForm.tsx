import { useState, type FormEvent } from "react";

const SearchForm: React.FC = () => {
  const [element1, setElement1] = useState<string>("");
  const [element2, setElement2] = useState<string>("");
  const [result, setResult] = useState<string | null>(null);
  const [notFound, setNotFound] = useState<boolean>(false);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    const res = await fetch(`/search?e1=${element1}&e2=${element2}`);
    const data = await res.json();

    if (data.found) {
      setResult(data.result);
      setNotFound(false);
    } else {
      setResult(null);
      setNotFound(true);
    }
  };

  return (
    <div>
      <form onSubmit={handleSubmit} className="flex gap-4">
        <input
          type="text"
          placeholder="Element 1"
          value={element1}
          onChange={(e) => setElement1(e.target.value)}
          className="border p-2"
        />
        <input
          type="text"
          placeholder="Element 2"
          value={element2}
          onChange={(e) => setElement2(e.target.value)}
          className="border p-2"
        />
        <button type="submit" className="bg-blue-500 text-white p-2">Search</button>
      </form>

      <div className="mt-4">
        {result && <p>Result: <strong>{result}</strong></p>}
        {notFound && <p className="text-red-500">Combination not found.</p>}
      </div>
    </div>
  );
}

export default SearchForm;

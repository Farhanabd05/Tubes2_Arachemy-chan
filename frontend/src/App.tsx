import React from 'react';
import SearchForm from './components/SearchForm';

const App: React.FC = () => {
  return (
    <div className="p-8 font-sans">
      <h1 className="text-2xl font-bold mb-4">Little Alchemy Search</h1>
      <SearchForm />
    </div>
  );
}

export default App;

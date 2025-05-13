import React, { useState } from 'react';
import { StatsDisplay } from './StatsDisplay';
import { TreeComponent } from './TreeComponent';

type PathObject = { [key: string]: string[] };
type MultipleResult = PathObject[];
type SingleResult = {
  steps: string[];
  runtime?: string;
  nodesVisited?: number | null;
};


const convertPathObjectToSingleResult = (pathObj: PathObject): SingleResult => {
  const pathName = Object.keys(pathObj).find(key => key.startsWith('Path'));
  const steps = pathName ? pathObj[pathName] : [];

  const runtime = pathObj['Runtime']?.[0] || '';
  const nodesVisited = pathObj['NodesVisited']?.[0]
    ? parseInt(pathObj['NodesVisited'][0])
    : null;

  return {
    steps,
    runtime,
    nodesVisited,
  };
};

export const MultipleResultDisplay: React.FC<{ results: MultipleResult }> = ({ results }) => {
  const [page, setPage] = useState(0);

  if (results.length === 0) return <div>No results found.</div>;

  const currentResult = convertPathObjectToSingleResult(results[page]);

  return (
    <div>
      

      <StatsDisplay runtime={currentResult.runtime} nodesVisited={currentResult.nodesVisited} />
      <TreeComponent steps={currentResult.steps} />
      <div style={{ marginBottom: '1rem' }}>
        <button onClick={() => setPage(p => Math.max(0, p - 1))} disabled={page === 0}>
          Previous
        </button>
        <span style={{ margin: '0 1rem' }}>
          Result {page + 1} of {results.length}
        </span>
        <button
          onClick={() => setPage(p => Math.min(results.length - 1, p + 1))}
          disabled={page === results.length - 1}
        >
          Next
        </button>
      </div>
    </div>
  );
};
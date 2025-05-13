type PathObject = { [key: string]: string[] };
type MultipleResult = PathObject[];

export const MultipleResultDisplay: React.FC<{ results: MultipleResult }> = ({ results }) => (
  <div className="multiple-results">
    {results.map((pathObj, i) => {
      const pathName = Object.keys(pathObj).find(key => key.startsWith('Path'));
      const steps = pathName ? pathObj[pathName] : [];
      const runtime = pathObj['Runtime']?.[0] || '';
      const nodesVisited = pathObj['NodesVisited']?.[0] || '';

      return (
        <div key={i}>
          <h3>{pathName}</h3>
          <ul>
            {steps.map((step, j) => (
              <li key={j}>{j}. 🧪 {step}</li>
            ))}
          </ul>
          {runtime && <p><strong>Runtime:</strong> {runtime}</p>}
          {nodesVisited && <p><strong>Nodes Visited:</strong> {nodesVisited}</p>}
        </div>
      );
    })}
  </div>
);
export const StatsDisplay: React.FC<{ runtime?: string; nodesVisited?: number | null }> = ({
  runtime,
  nodesVisited
}) => (
  <>
    {runtime && <div><strong>Runtime:</strong> {runtime} ns</div>}
    {nodesVisited !== null && <div><strong>Nodes Visited:</strong> {nodesVisited}</div>}
  </>
);
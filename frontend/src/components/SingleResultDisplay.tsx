import TreeComponent from "./TreeComponent";
import { StatsDisplay } from "./StatsDisplay";

interface SingleResult {
  found: boolean;
  steps: string[];
  runtime?: string; 
  nodesVisited?: number | null;
}

export const SingleResultDisplay: React.FC<{ result: SingleResult }> = ({ result }) => (
  <div>
    <StatsDisplay runtime={result.runtime} nodesVisited={result.nodesVisited} />
    <TreeComponent steps = { result.steps }/>
  </div>
);
import TreeComponent from "./TreeComponent";

interface SingleResult {
  found: boolean;
  steps: string[];
}

export const SingleResultDisplay: React.FC<{ result: SingleResult }> = ({ result }) => (
  <div>
    <TreeComponent steps = { result.steps }/>
  </div>
);
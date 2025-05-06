
// File: src/types.ts
export interface Recipe {
    input: string[];
    output: string;
  }
  
export interface GraphNode {
id: string;
label: string;
}

export interface GraphEdge {
from: string;
to: string;
}

export interface GraphData {
nodes: GraphNode[];
edges: GraphEdge[];
}

export interface SearchResult {
status: string;
path?: string[];
found: boolean;
}

export interface EdgeWithId extends GraphEdge {
    id: string;
}
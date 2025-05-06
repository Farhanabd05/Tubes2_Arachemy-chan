// / File: src/api.ts
// / File: src/api.ts
import type { Recipe, GraphData, SearchResult, GraphEdge, GraphNode } from './types';

const API_BASE_URL = 'http://localhost:8080/api';

export const apiService = {
  async uploadRecipes(recipes: Recipe[]): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/recipes`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(recipes),
    });
    
    if (!response.ok) {
      throw new Error('Failed to upload recipes');
    }
  },
  
  async searchPath(start: string, target: string): Promise<SearchResult> {
    const response = await fetch(`${API_BASE_URL}/search?start=${start}&target=${target}`);
    
    if (!response.ok) {
      throw new Error('Failed to search path');
    }
    
    return await response.json();
  },
  
  async getGraphData(): Promise<GraphData> {
    const response = await fetch(`${API_BASE_URL}/graph`);
    
    if (!response.ok) {
      throw new Error('Failed to get graph data');
    }
    
    const data = await response.json();
    
    // Transform the data into the format expected by the visualization library
    const nodes: GraphNode[] = data.nodes.map((node: string) => ({
      id: node,
      label: node,
    }));
    
    const edges: GraphEdge[] = [];
    Object.entries(data.connections).forEach(([source, targets]) => {
      (targets as string[]).forEach((target) => {
        edges.push({
          from: source,
          to: target,
        });
      });
    });
    
    return { nodes, edges };
  },
};

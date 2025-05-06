// File: src/components/GraphVisualization.tsx
import React, { useEffect, useRef, useState } from 'react';
import { Network } from 'vis-network';
import { DataSet } from 'vis-data';
import { apiService } from '../api';
import type { GraphData, GraphNode, EdgeWithId } from '../types';

interface GraphVisualizationProps {
  highlightPath?: string[];
}

export const GraphVisualization: React.FC<GraphVisualizationProps> = ({ highlightPath }) => {
  const containerRef = useRef<HTMLDivElement>(null);
  const networkRef = useRef<Network | null>(null);
  const [graphData, setGraphData] = useState<GraphData | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const fetchGraphData = async () => {
    try {
      setLoading(true);
      const data = await apiService.getGraphData();
      setGraphData(data);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load graph data');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchGraphData();
  }, []);

  useEffect(() => {
    if (!containerRef.current || !graphData) return;

    // Create nodes and edges datasets
    const nodes = new DataSet<GraphNode>(
      graphData.nodes.map(node => ({
        ...node,
        color: highlightPath?.includes(node.id) ? '#ff7e00' : undefined,
        font: {
          color: highlightPath?.includes(node.id) ? '#000000' : undefined,
          bold: highlightPath?.includes(node.id)
        }
      }))
    );

    const edges = new DataSet<EdgeWithId>(
      graphData.edges.map((edge, index) => {
        const isHighlighted = 
          highlightPath && 
          highlightPath.includes(edge.from) && 
          highlightPath.includes(edge.to) &&
          highlightPath.indexOf(edge.to) === highlightPath.indexOf(edge.from) + 1;
        
        return {
            id: index.toString(), // Add an id property
          ...edge,
          color: isHighlighted ? '#ff7e00' : undefined,
          width: isHighlighted ? 3 : 1,
          arrows: 'to'
        };
      })
    );

    // Create the network
    const options = {
      physics: {
        enabled: true,
        solver: 'forceAtlas2Based',
        forceAtlas2Based: {
          gravitationalConstant: -100,
          centralGravity: 0.01,
          springLength: 100,
          springConstant: 0.08
        },
        stabilization: {
          iterations: 100
        }
      },
      nodes: {
        shape: 'circle',
        size: 25,
        font: {
          size: 14
        }
      },
      edges: {
        smooth: {
            enabled: true,
            type: 'continuous',
            roundness: 0.5
        }
      }
    };
    
    networkRef.current = new Network(containerRef.current, { nodes, edges }, options);

    return () => {
      if (networkRef.current) {
        networkRef.current.destroy();
        networkRef.current = null;
      }
    };
  }, [graphData, highlightPath]);

  const handleRefresh = () => {
    fetchGraphData();
  };

  if (loading && !graphData) {
    return <div className="text-center py-10">Loading graph data...</div>;
  }

  if (error) {
    return (
      <div className="text-center py-10">
        <p className="text-red-500 mb-4">{error}</p>
        <button
          onClick={handleRefresh}
          className="px-4 py-2 bg-blue-600 text-white font-medium rounded hover:bg-blue-700"
        >
          Retry
        </button>
      </div>
    );
  }

  return (
    <div className="border rounded-lg bg-white overflow-hidden">
      <div className="p-4 border-b bg-gray-50 flex justify-between items-center">
        <h2 className="text-xl font-bold">Graph Visualization</h2>
        <button
          onClick={handleRefresh}
          className="px-3 py-1 bg-gray-200 text-gray-700 rounded hover:bg-gray-300"
        >
          Refresh
        </button>
      </div>
      <div ref={containerRef} style={{ height: '500px' }} />
    </div>
  );
};
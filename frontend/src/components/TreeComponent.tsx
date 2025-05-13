import React, { useEffect, useRef, useState } from "react";
import Tree from 'react-d3-tree';
import '../App.css';

const testTreeData = [
  {
    name: 'Blade',
    children: [ {imageUrl: 'https://png.pngtree.com/png-vector/20190418/ourmid/pngtree-vector-plus-icon-png-image_956060.jpg', children: [
      {
        name: 'Stone',
        children: [ {imageUrl: 'https://png.pngtree.com/png-vector/20190418/ourmid/pngtree-vector-plus-icon-png-image_956060.jpg', children: [
          {
            name: 'Earth',
          },
          {
            name: 'Pressure',
            children: [ {imageUrl: 'https://png.pngtree.com/png-vector/20190418/ourmid/pngtree-vector-plus-icon-png-image_956060.jpg', children: [
              {
                name: 'Air'
              },
              {
                name: 'Air'
              }
            ]}]
          }
        ]}]
      },
      {
        name: 'Metal',
        imageLocal: '',
        children: [ {imageUrl: 'https://png.pngtree.com/png-vector/20190418/ourmid/pngtree-vector-plus-icon-png-image_956060.jpg', children: [
          {
            name: 'Stone',
            children: [ {imageUrl: 'https://png.pngtree.com/png-vector/20190418/ourmid/pngtree-vector-plus-icon-png-image_956060.jpg', children: [
              {
                name: 'Earth',
              },
              {
                name: 'Pressure',
                children: [ {imageUrl: 'https://png.pngtree.com/png-vector/20190418/ourmid/pngtree-vector-plus-icon-png-image_956060.jpg', children: [
                  {
                    name: 'Air'
                  },
                  {
                    name: 'Air'
                  }
                ]}]
              }
            ]}]
          },
          {
            name: 'Fire'
          }
        ]}]
      }
    ]}]
  }
];

const testStepData = [
  "air + air = pressure",
  "earth + pressure = stone",
  "air + data = pressure",
  "earth + pressure = stone",
  "stone + fire = metal",
  "stone + metal = blade"
];

type Step = { left: string; right: string; result: string };

type TreeNode = {
  name?: string;
  children?: TreeNode[];
  imageUrl?: string;
};

function parseStep(step: string): Step {
  const [leftRight, result] = step.split(" = ");
  const [left, right] = leftRight.split(" + ");
  return { left, right, result };
}

function buildTreeFromSteps(steps: string[]): TreeNode {
  const parsedSteps = steps.map(parseStep);

  const stepToNode = (index: number): { node: TreeNode; indexUsed: number } => {
    const { left, right, result } = parsedSteps[index];

    const rightResult = findNodeFromName(right, index);
    const leftResult = findNodeFromName(left, rightResult.indexUsed);

    const minUsed = Math.min(leftResult.indexUsed, rightResult.indexUsed);

    const node: TreeNode = {
      name: result,
      children: [
        {
          imageUrl: 'https://png.pngtree.com/png-vector/20190418/ourmid/pngtree-vector-plus-icon-png-image_956060.jpg',
          children: [leftResult.node, rightResult.node],
        },
      ],
    };

    return { node, indexUsed: minUsed };
  };

  const findNodeFromName = (name: string, beforeIndex: number): { node: TreeNode; indexUsed: number } => {
    for (let i = beforeIndex - 1; i >= 0; i--) {
      if (parsedSteps[i].result === name) {
        return stepToNode(i);
      }
    }
    return { node: { name }, indexUsed: beforeIndex }; // base/raw element
  };

  return stepToNode(parsedSteps.length - 1).node;
}

const CircleImage = ({ imageUrl }: { imageUrl: string }) => (
  <foreignObject width={60} height={60} x={-30} y={-30}>
    <div
      style={{
        width: '100%',
        height: '100%',
        borderRadius: '50%',
        overflow: 'hidden',
        border: '2px solid #4A90E2',
        backgroundColor: '#000',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
      }}
    >
      {imageUrl && (
        <img
          src={imageUrl}
          alt=""
          style={{
            width: '80%',
            height: '80%',
            objectFit: 'cover',
          }}
        />
      )}
    </div>
  </foreignObject>
);

const renderNode = ({ nodeDatum }: { nodeDatum: any }) => {
  const hasName = Boolean(nodeDatum.name);

  return (
    <g>
      {hasName ? (
        // Box layout (image + text)
        <foreignObject width={150} height={60} x={-75} y={-30}>
          <div
            style={{
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              borderRadius: '12px',
              border: '2px solid #4A90E2',
              background: '#000',
              padding: '8px',
              color: 'white',
              fontSize: '14px',
              fontFamily: 'sans-serif',
              width: '100%',
              height: '100%',
              boxSizing: 'border-box',
            }}
          >
            {nodeDatum.imageUrl && (
              <img
                src={nodeDatum.imageUrl}
                alt=""
                style={{
                  width: '40px',
                  height: '40px',
                  borderRadius: '8px',
                  objectFit: 'cover',
                  marginRight: '12px'
                }}
              />
            )}
            <div style={{ whiteSpace: 'nowrap' }}>{nodeDatum.name}</div>
          </div>
        </foreignObject>
      ) : (
        // Circular image-only layout
        <CircleImage imageUrl={nodeDatum.imageUrl} />
      )}
    </g>
  );
};

const NODE_WIDTH = 170;
const NODE_HEIGHT = 80;

function calculateTreeStats(node: any, depth = 0, levels: number[] = []): { depth: number; maxWidth: number } {
  if (!levels[depth]) levels[depth] = 0;
  levels[depth] += 1;

  if (!node.children || node.children.length === 0) {
    return { depth: levels.length, maxWidth: Math.max(...levels) };
  }

  node.children.forEach((child: any) => {
    calculateTreeStats(child, depth + 1, levels);
  });

  return { depth: levels.length, maxWidth: Math.max(...levels) };
}

useEffect(() => {
  if (treeContainerRef.current) {
    const { offsetWidth, offsetHeight } = treeContainerRef.current;

    // Center the tree by default
    setTranslate({
      x: offsetWidth / 2,
      y: offsetHeight / 10 // you can tweak this
    });
  }
}, [treeData]); // make sure it updates when data changes

export const TreeComponent: React.FC<{ steps: string[] }> = ({ steps }) => {
  const treeContainerRef = useRef<HTMLDivElement>(null);
  const [translate, setTranslate] = useState({ x: 0, y: 0 });
  const [zoom, setZoom] = useState(1);

  const treeData = [buildTreeFromSteps(steps)];

  useEffect(() => {
    const container = treeContainerRef.current;
    if (!container) return;

    const { width, height } = container.getBoundingClientRect();
    setTranslate({ x: width / 2, y: 100 });

    const stats = calculateTreeStats(treeData[0]);
    const zoomX = width / (NODE_WIDTH * stats.maxWidth);
    const zoomY = height / (NODE_HEIGHT * stats.depth);
    const finalZoom = Math.min(zoomX, zoomY) * 0.9;

    setZoom(finalZoom);
  }, [steps]);

  return (
    <div ref={treeContainerRef} className="tree-container" style={{ width: "90%", height: "80vh" }}>
      <Tree
        data={treeData}
        orientation="vertical"
        pathFunc="step"
        collapsible={false}
        enableLegacyTransitions={true}
        renderCustomNodeElement={renderNode}
        nodeSize={{ x: NODE_WIDTH, y: NODE_HEIGHT }}
        separation={{ siblings: 1, nonSiblings: 1.2 }}
        zoom={zoom}
        translate={translate}
      />
    </div>
  );
};

export default TreeComponent;
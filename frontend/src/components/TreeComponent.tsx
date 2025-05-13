import React, { useEffect, useRef, useState } from "react";
import Tree from 'react-d3-tree';
import '../App.css';

type Step = { left: string; right: string; result: string };

type TreeNode = {
  name: string;
  children?: TreeNode[];
  imageUrl?: string;
};

let elementImageMap: Record<string, string> = {};
async function loadElementImageMap() {
  if (Object.keys(elementImageMap).length === 0) {
    const response = await fetch('/mapped_elements.json');
    const data = await response.json();
    for (const entry of data) {
      elementImageMap[entry.Element.toLowerCase()] = `/images/${entry.ElementImage}`;
    }
  }
}

function parseStep(step: string): Step {
  const [leftRight, result] = step.split(" = ");
  const [left, right] = leftRight.split(" + ");
  return { left: left.trim(), right: right.trim(), result: result.trim() };
}

function buildTreeFromSteps(steps: string[]): TreeNode {
  const parsedSteps = steps.map(parseStep);

  const stepToNode = (index: number): { node: TreeNode; indexUsed: number } => {
    const { left, right, result } = parsedSteps[index];

    let leftResult, rightResult;
    if (index === 0 || parsedSteps[index - 1].result === right) {
      rightResult = findNodeFromName(right, index);
      leftResult = findNodeFromName(left, rightResult.indexUsed);
    } else {
      leftResult = findNodeFromName(left, index);
      rightResult = findNodeFromName(right, leftResult.indexUsed);
    }

    const minUsed = Math.min(leftResult.indexUsed, rightResult.indexUsed);

    const imageUrl = elementImageMap[result?.toLowerCase()] ?? "/images/default.svg";

    const node: TreeNode = {
      name: result || "Unknown",
      imageUrl,
      children: [
        {
          name: "plus",
          imageUrl: "/images/plus.png",
          children: [leftResult.node, rightResult.node].filter(Boolean),
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

    const fallbackImage = elementImageMap[name?.toLowerCase()] ?? "/images/default.svg";
    if (!(name?.toLowerCase() in elementImageMap)) {
      console.warn(`⚠️ No image mapping for element: ${name}`);
    }

    const fallbackNode: TreeNode = {
      name: name || "Unknown",
      imageUrl: fallbackImage,
    };
    return { node: fallbackNode, indexUsed: beforeIndex };
  };

  const root = stepToNode(parsedSteps.length - 1).node;
  console.log("✅ Final Tree:", JSON.stringify(root, null, 2));
  return root;
}

const CircleImage = ({ imageUrl }: { imageUrl: string }) => (
  <foreignObject width={50} height={50} x={-25} y={-25}>
    <div
      style={{
        width: '100%',
        height: '100%',
        clipPath: 'circle(50% at 50% 50%)',
        backgroundColor: '#4CAF50',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
      }}
    >
      <img
        src={imageUrl}
        alt=""
        style={{
          width: '60%',
          height: '60%',
          objectFit: 'cover',
        }}
      />
    </div>
  </foreignObject>
);

const renderNode = ({ nodeDatum }: { nodeDatum: any }) => {
  return (
    <g>
      {nodeDatum.name != 'plus' ? (
        <foreignObject width={150} height={60} x={-75} y={-30}>
          <div
            style={{
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              borderRadius: '12px',
              border: '2px solid #212121',
              background: ' #00BCD4',
              padding: '8px',
              color: 'white',
              fontSize: '14px',
              fontFamily: 'sans-serif',
              fontWeight: 'bold',
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

export const TreeComponent: React.FC<{ steps: string[] }> = ({ steps }) => {
  const treeContainerRef = useRef<HTMLDivElement>(null);
  const [translate, setTranslate] = useState({ x: 0, y: 0 });
  const [zoom, setZoom] = useState(1);
  const [treeData, setTreeData] = useState<TreeNode[]>([]);

  useEffect(() => {
    const container = treeContainerRef.current;
    if (!container) return;

    const { width, height } = container.getBoundingClientRect();
    setTranslate({ x: width / 2, y: 100 });

    (async () => {
      await loadElementImageMap();
      const tree = buildTreeFromSteps(steps);
      const stats = calculateTreeStats(tree);
      const zoomX = width / (NODE_WIDTH * stats.maxWidth);
      const zoomY = height / (NODE_HEIGHT * stats.depth);
      const finalZoom = Math.min(zoomX, zoomY) * 0.9;

      setTreeData([tree]);
      setZoom(finalZoom);
    })();
  }, [steps]);

  return (
    <div ref={treeContainerRef} className="tree-container" style={{ width: "90%", height: "80vh" }}>
      {treeData.length > 0 && (
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
      )}
    </div>
  );
};

export default TreeComponent;
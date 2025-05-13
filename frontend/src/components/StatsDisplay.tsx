export const StatsDisplay: React.FC<{ runtime?: string; nodesVisited?: number | null }> = ({
  runtime,
  nodesVisited
}) => (
  <div style={{ width: '100%', display: 'flex', justifyContent: 'center', marginBottom: '1rem' }}>
    <div style={{ display: 'flex', width: '100%', maxWidth: '800px' }}>
      {runtime && (
        <div style={{
          flex: 1,
          padding: '0.5rem 1rem',
          border: '3px solid #032202',
          borderRadius: '8px 0 0 8px',
          color: '#000000',
          backgroundColor: '#f9f9f9',
          textAlign: 'center'
        }}>
          <strong>Runtime:</strong> {runtime}
        </div>
      )}
      {nodesVisited !== null && (
        <div style={{
          flex: 1,
          padding: '0.5rem 1rem',
          border: '3px solid #032202',
          borderRadius: '0 8px 8px 0',
          color: '#000000',
          backgroundColor: '#f9f9f9',
          textAlign: 'center',
          borderLeft: 'none'
        }}>
          <strong>Nodes Visited:</strong> {nodesVisited}
        </div>
      )}
    </div>
  </div>
);
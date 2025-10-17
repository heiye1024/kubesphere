import React from 'react';
import { useParams } from 'react-router-dom';
import { CircularProgress, List, ListItem, ListItemText, Typography } from '@mui/material';
import useSWR from 'swr';
import { listSnapshots } from '../api/client';

const SnapshotListPage: React.FC = () => {
  const { namespace = 'default' } = useParams();
  const { data, error } = useSWR(['snapshots', namespace], () => listSnapshots(namespace));

  if (error) {
    return <Typography color="error">Failed to load snapshots: {error.message}</Typography>;
  }
  if (!data) {
    return <CircularProgress />;
  }

  return (
    <List>
      {data.map((snap: any) => (
        <ListItem key={snap.metadata.name} divider>
          <ListItemText
            primary={snap.metadata.name}
            secondary={`Ready: ${snap.status?.readyToUse ? 'Yes' : 'No'}`}
          />
        </ListItem>
      ))}
    </List>
  );
};

export default SnapshotListPage;

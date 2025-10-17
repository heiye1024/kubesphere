import React from 'react';
import { useParams } from 'react-router-dom';
import { CircularProgress, List, ListItem, ListItemText, Typography } from '@mui/material';
import useSWR from 'swr';
import { listDisks } from '../api/client';

const DiskListPage: React.FC = () => {
  const { namespace = 'default' } = useParams();
  const { data, error } = useSWR(['disks', namespace], () => listDisks(namespace));

  if (error) {
    return <Typography color="error">Failed to load disks: {error.message}</Typography>;
  }
  if (!data) {
    return <CircularProgress />;
  }

  return (
    <List>
      {data.map((disk: any) => (
        <ListItem key={disk.metadata.name} divider>
          <ListItemText
            primary={disk.metadata.name}
            secondary={`${disk.spec.size} • ${disk.spec.storageClass}`}
          />
        </ListItem>
      ))}
    </List>
  );
};

export default DiskListPage;

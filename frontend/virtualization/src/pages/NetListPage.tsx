import React from 'react';
import { useParams } from 'react-router-dom';
import { CircularProgress, List, ListItem, ListItemText, Typography } from '@mui/material';
import useSWR from 'swr';
import { listNets } from '../api/client';

const NetListPage: React.FC = () => {
  const { namespace = 'default' } = useParams();
  const { data, error } = useSWR(['nets', namespace], () => listNets(namespace));

  if (error) {
    return <Typography color="error">Failed to load networks: {error.message}</Typography>;
  }
  if (!data) {
    return <CircularProgress />;
  }

  return (
    <List>
      {data.map((net: any) => (
        <ListItem key={net.metadata.name} divider>
          <ListItemText
            primary={net.metadata.name}
            secondary={`${net.spec.nadTemplate.slice(0, 40)}...`}
          />
        </ListItem>
      ))}
    </List>
  );
};

export default NetListPage;

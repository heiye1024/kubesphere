import React from 'react';
import { useParams, Link } from 'react-router-dom';
import { Card, CardContent, CircularProgress, Grid, Typography, Button, Stack } from '@mui/material';
import useSWR from 'swr';
import { listVMs, VM } from '../api/client';

const VMListPage: React.FC = () => {
  const { namespace = 'default' } = useParams();
  const { data, error } = useSWR(['vms', namespace], () => listVMs(namespace));

  if (error) {
    return <Typography color="error">Failed to load VMs: {error.message}</Typography>;
  }
  if (!data) {
    return <CircularProgress />;
  }

  return (
    <Stack spacing={2}>
      <Stack direction="row" justifyContent="space-between" alignItems="center">
        <Typography variant="h5">Virtual Machines</Typography>
        <Button
          component={Link}
          to={`/virtualization/projects/${namespace}/templates`}
          variant="contained"
          color="primary"
        >
          Create from Template
        </Button>
      </Stack>
      <Grid container spacing={2}>
        {data.map((vm: VM) => (
          <Grid item xs={12} md={4} key={vm.metadata.name}>
            <Card>
              <CardContent>
                <Typography variant="h6">{vm.metadata.name}</Typography>
                <Typography variant="body2">CPU: {vm.spec.cpu}</Typography>
                <Typography variant="body2">Memory: {vm.spec.memory}</Typography>
                <Typography variant="body2">
                  Power: {vm.status?.powerState ?? vm.spec.powerState}
                </Typography>
                <Button
                  component={Link}
                  to={`/virtualization/projects/${namespace}/vms/${vm.metadata.name}`}
                  sx={{ mt: 2 }}
                >
                  View Details
                </Button>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>
    </Stack>
  );
};

export default VMListPage;

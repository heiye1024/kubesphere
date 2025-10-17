import React from 'react';
import { Link, useParams } from 'react-router-dom';
import { Button, CircularProgress, List, ListItem, ListItemText, Stack, Typography } from '@mui/material';
import useSWR from 'swr';
import { listTemplates } from '../api/client';

const TemplateListPage: React.FC = () => {
  const { namespace = 'default' } = useParams();
  const { data, error } = useSWR(['templates', namespace], () => listTemplates(namespace));

  if (error) {
    return <Typography color="error">Failed to load templates: {error.message}</Typography>;
  }
  if (!data) {
    return <CircularProgress />;
  }

  return (
    <Stack spacing={2}>
      <Typography variant="h5">Templates</Typography>
      <List>
        {data.map((tpl: any) => (
          <ListItem key={tpl.metadata.name} divider>
            <ListItemText
              primary={tpl.metadata.name}
              secondary={`${tpl.spec.parameters.cpu} CPU • ${tpl.spec.parameters.memory}`}
            />
            <Button
              component={Link}
              to={`/virtualization/projects/${namespace}/templates/${tpl.metadata.name}/create`}
            >
              Launch Wizard
            </Button>
          </ListItem>
        ))}
      </List>
    </Stack>
  );
};

export default TemplateListPage;

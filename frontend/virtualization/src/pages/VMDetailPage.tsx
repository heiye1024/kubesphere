import React, { useState } from 'react';
import { useParams } from 'react-router-dom';
import {
  Alert,
  Box,
  Button,
  Card,
  CardContent,
  CircularProgress,
  Snackbar,
  Stack,
  Tooltip,
  Typography,
} from '@mui/material';
import useSWR from 'swr';
import axios from 'axios';
import { listVMs, openConsole, powerOffVM, powerOnVM, VM } from '../api/client';

const fetchVM = (namespace: string, name: string) =>
  listVMs(namespace).then(vms => vms.find(vm => vm.metadata.name === name));

const VMDetailPage: React.FC = () => {
  const { namespace = 'default', name = '' } = useParams();
  const { data, error, mutate } = useSWR(['vm', namespace, name], () => fetchVM(namespace, name));
  const [denyReasons, setDenyReasons] = useState<Record<string, string | undefined>>({});
  const [snackbar, setSnackbar] = useState<string | null>(null);
  const [consoleURL, setConsoleURL] = useState<string | null>(null);

  const handlePower = async (action: 'powerOn' | 'powerOff') => {
    const request = action === 'powerOn' ? powerOnVM : powerOffVM;
    try {
      await request(namespace, name);
      setDenyReasons(prev => ({ ...prev, [action]: undefined }));
      setSnackbar(`Requested ${action === 'powerOn' ? 'power on' : 'power off'} operation.`);
      mutate();
    } catch (err) {
      if (axios.isAxiosError(err)) {
        const reason =
          (err.response?.headers['x-deny-reason'] as string | undefined) ||
          (typeof err.response?.data === 'object' && err.response?.data !== null
            ? (err.response?.data as { error?: string }).error
            : undefined) ||
          err.message;
        setDenyReasons(prev => ({ ...prev, [action]: reason }));
        setSnackbar(reason);
      }
    }
  };

  const handleConsole = async () => {
    try {
      const { consoleURL: url } = await openConsole(namespace, name);
      setConsoleURL(url);
      setSnackbar('Console URL retrieved.');
    } catch (err) {
      if (axios.isAxiosError(err)) {
        const reason = err.response?.data?.error || err.message;
        setSnackbar(reason);
      }
    }
  };

  if (error) {
    return <Typography color="error">Failed to load VM: {error.message}</Typography>;
  }
  if (!data) {
    return <CircularProgress />;
  }

  return (
    <Card>
      <CardContent>
        <Typography variant="h4" gutterBottom>
          {data.metadata.name}
        </Typography>
        <Stack spacing={1}>
          <Typography>CPU: {data.spec.cpu}</Typography>
          <Typography>Memory: {data.spec.memory}</Typography>
          <Typography>Power: {data.status?.powerState ?? data.spec.powerState}</Typography>
        </Stack>
        <Box sx={{ mt: 3, display: 'flex', gap: 2 }}>
          <Tooltip title={denyReasons.powerOn || ''} disableHoverListener={!denyReasons.powerOn}>
            <span>
              <Button
                variant="contained"
                onClick={() => handlePower('powerOn')}
                disabled={Boolean(denyReasons.powerOn)}
              >
                Power On
              </Button>
            </span>
          </Tooltip>
          <Tooltip title={denyReasons.powerOff || ''} disableHoverListener={!denyReasons.powerOff}>
            <span>
              <Button
                variant="outlined"
                color="warning"
                onClick={() => handlePower('powerOff')}
                disabled={Boolean(denyReasons.powerOff)}
              >
                Power Off
              </Button>
            </span>
          </Tooltip>
        </Box>
        <Box sx={{ mt: 4 }}>
          <Typography variant="h6">Console</Typography>
          <Typography variant="body2" color="text.secondary">
            Launch embedded VNC/Serial console. Connection retries automatically on network errors.
          </Typography>
          <Button sx={{ mt: 1 }} variant="outlined" onClick={handleConsole}>
            Open Web Console
          </Button>
          {consoleURL && (
            <Alert sx={{ mt: 2 }} severity="info">
              Connect via{' '}
              <a href={consoleURL} target="_blank" rel="noopener noreferrer">
                {consoleURL}
              </a>
            </Alert>
          )}
        </Box>
      </CardContent>
      <Snackbar
        open={Boolean(snackbar)}
        autoHideDuration={4000}
        onClose={() => setSnackbar(null)}
        message={snackbar}
      />
    </Card>
  );
};

export default VMDetailPage;

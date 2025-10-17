import React, { useState } from 'react';
import { useParams } from 'react-router-dom';
import {
  Box,
  Button,
  Step,
  StepLabel,
  Stepper,
  TextField,
  Typography,
  Chip,
  Stack,
} from '@mui/material';

const steps = ['Basics', 'Resources', 'Disks & Networks', 'Review'];

const VMCreateWizard: React.FC = () => {
  const { namespace = 'default', name = '' } = useParams();
  const [activeStep, setActiveStep] = useState(0);
  const [cpu, setCpu] = useState('2');
  const [memory, setMemory] = useState('4Gi');
  const [hugepages, setHugepages] = useState('');
  const [numaPolicy, setNumaPolicy] = useState('none');

  const next = () => setActiveStep(prev => Math.min(prev + 1, steps.length - 1));
  const back = () => setActiveStep(prev => Math.max(prev - 1, 0));

  return (
    <Stack spacing={3}>
      <Typography variant="h5">Create VM from template {name}</Typography>
      <Typography variant="body2" color="text.secondary">
        Target namespace: {namespace}
      </Typography>
      <Stepper activeStep={activeStep} alternativeLabel>
        {steps.map(step => (
          <Step key={step}>
            <StepLabel>{step}</StepLabel>
          </Step>
        ))}
      </Stepper>
      {activeStep === 0 && (
        <Stack spacing={2}>
          <TextField label="CPU" value={cpu} onChange={e => setCpu(e.target.value)} fullWidth />
          <TextField label="Memory" value={memory} onChange={e => setMemory(e.target.value)} fullWidth />
        </Stack>
      )}
      {activeStep === 1 && (
        <Stack spacing={2}>
          <TextField label="Hugepages" value={hugepages} onChange={e => setHugepages(e.target.value)} fullWidth />
          <TextField label="NUMA Policy" value={numaPolicy} onChange={e => setNumaPolicy(e.target.value)} fullWidth />
        </Stack>
      )}
      {activeStep === 2 && (
        <Typography>Attach data disks and networks in the next iteration.</Typography>
      )}
      {activeStep === 3 && (
        <Box>
          <Typography variant="h6">Review</Typography>
          <Stack direction="row" spacing={1} mt={2}>
            <Chip label={`CPU ${cpu}`} />
            <Chip label={`Memory ${memory}`} />
            {hugepages && <Chip label={`Hugepages ${hugepages}`} />}
            <Chip label={`NUMA ${numaPolicy}`} />
          </Stack>
        </Box>
      )}
      <Stack direction="row" justifyContent="space-between">
        <Button disabled={activeStep === 0} onClick={back}>
          Back
        </Button>
        {activeStep < steps.length - 1 ? (
          <Button variant="contained" onClick={next}>
            Next
          </Button>
        ) : (
          <Button variant="contained" color="primary">
            Submit
          </Button>
        )}
      </Stack>
    </Stack>
  );
};

export default VMCreateWizard;

import axios from 'axios';

export const virtualizationAPIBase = '/kapis/virtualization.kubesphere.io/v1beta1';

export interface Envelope<T> {
  data: T;
  total?: number;
  traceID: string;
  auditID: string;
  message?: string;
}

const client = axios.create({
  baseURL: virtualizationAPIBase,
});

export interface VM {
  metadata: { name: string };
  spec: { cpu: string; memory: string; powerState: string };
  status?: { powerState?: string };
}

const unwrap = <T>(value: Envelope<T>): T => value.data;

export const listVMs = (namespace: string, params?: Record<string, string | string[]>) =>
  client
    .get<Envelope<VM[]>>(`/projects/${namespace}/vms`, { params })
    .then(res => unwrap(res.data));

export const powerOnVM = (namespace: string, name: string) =>
  client.post(`/projects/${namespace}/vms/${name}:powerOn`);

export const powerOffVM = (namespace: string, name: string) =>
  client.post(`/projects/${namespace}/vms/${name}:powerOff`);

export const openConsole = (namespace: string, name: string) =>
  client
    .post<Envelope<{ consoleURL: string }>>(`/projects/${namespace}/vms/${name}:console`)
    .then(res => unwrap(res.data));

export const listDisks = (namespace: string, params?: Record<string, string | string[]>) =>
  client
    .get<Envelope<unknown[]>>(`/projects/${namespace}/disks`, { params })
    .then(res => unwrap(res.data));

export const listNets = (namespace: string, params?: Record<string, string | string[]>) =>
  client
    .get<Envelope<unknown[]>>(`/projects/${namespace}/nets`, { params })
    .then(res => unwrap(res.data));

export const listSnapshots = (namespace: string, params?: Record<string, string | string[]>) =>
  client
    .get<Envelope<unknown[]>>(`/projects/${namespace}/snapshots`, { params })
    .then(res => unwrap(res.data));

export const listTemplates = (namespace: string, params?: Record<string, string | string[]>) =>
  client
    .get<Envelope<unknown[]>>(`/projects/${namespace}/templates`, { params })
    .then(res => unwrap(res.data));

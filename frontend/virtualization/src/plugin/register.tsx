import React from 'react';
import DiskListPage from '../pages/DiskListPage';
import NetListPage from '../pages/NetListPage';
import SnapshotListPage from '../pages/SnapshotListPage';
import TemplateListPage from '../pages/TemplateListPage';
import VMCreateWizard from '../pages/VMCreateWizard';
import VMDetailPage from '../pages/VMDetailPage';
import VMListPage from '../pages/VMListPage';
import type { PluginAPI } from './types';

const NS_PARAM = ':namespace';

export function registerVirtualizationPlugin(api: PluginAPI) {
  const namespacePrefix = `/virtualization/projects/${NS_PARAM}`;

  api.registerRoute({ path: `${namespacePrefix}/vms`, element: <VMListPage /> });
  api.registerRoute({ path: `${namespacePrefix}/vms/:name`, element: <VMDetailPage /> });
  api.registerRoute({ path: `${namespacePrefix}/disks`, element: <DiskListPage /> });
  api.registerRoute({ path: `${namespacePrefix}/nets`, element: <NetListPage /> });
  api.registerRoute({ path: `${namespacePrefix}/snapshots`, element: <SnapshotListPage /> });
  api.registerRoute({ path: `${namespacePrefix}/templates`, element: <TemplateListPage /> });
  api.registerRoute({
    path: `${namespacePrefix}/templates/:name/create`,
    element: <VMCreateWizard />,
  });

  api.registerNavigation({
    id: 'virtualization',
    label: 'Virtualization',
    pathTemplate: `${namespacePrefix}/vms`,
    children: [
      {
        id: 'virtualization-vms',
        label: 'Virtual Machines',
        pathTemplate: `${namespacePrefix}/vms`,
      },
      {
        id: 'virtualization-disks',
        label: 'Disks',
        pathTemplate: `${namespacePrefix}/disks`,
      },
      {
        id: 'virtualization-nets',
        label: 'Networks',
        pathTemplate: `${namespacePrefix}/nets`,
      },
      {
        id: 'virtualization-snapshots',
        label: 'Snapshots',
        pathTemplate: `${namespacePrefix}/snapshots`,
      },
      {
        id: 'virtualization-templates',
        label: 'Templates',
        pathTemplate: `${namespacePrefix}/templates`,
      },
    ],
  });
}

export default registerVirtualizationPlugin;

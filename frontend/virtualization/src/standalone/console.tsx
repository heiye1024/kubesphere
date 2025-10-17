import React, { useMemo, useState } from 'react';
import { createRoot } from 'react-dom/client';
import { BrowserRouter, Navigate, Route, Routes, useLocation } from 'react-router-dom';
import Layout from '../components/Layout';
import type { NavItem, PluginAPI, RouteConfig } from '../plugin/types';

interface ResolvedNavItem {
  id: string;
  label: string;
  path: string;
  children?: ResolvedNavItem[];
}

const resolveNavItem = (item: NavItem, namespace: string): ResolvedNavItem => ({
  id: item.id,
  label: item.label,
  path: item.pathTemplate.replace(':namespace', namespace),
  children: item.children?.map(child => resolveNavItem(child, namespace)),
});

interface StandaloneShellProps {
  routes: RouteConfig[];
  navItems: NavItem[];
}

const DefaultRedirect: React.FC<{ namespace: string }> = ({ namespace }) => {
  const location = useLocation();
  if (location.pathname && location.pathname !== '/' && location.pathname !== '') {
    return null;
  }
  return <Navigate to={`/virtualization/projects/${namespace}/vms`} replace />;
};

const StandaloneShell: React.FC<StandaloneShellProps> = ({ routes, navItems }) => {
  const [namespace, setNamespace] = useState('default');
  const resolvedNav = useMemo(
    () => navItems.map(item => resolveNavItem(item, namespace)),
    [navItems, namespace]
  );

  return (
    <BrowserRouter>
      <Layout
        navItems={resolvedNav}
        namespace={namespace}
        onNamespaceChange={setNamespace}
      >
        <Routes>
          <Route path="/" element={<DefaultRedirect namespace={namespace} />} />
          {routes.map(route => (
            <Route key={route.path} path={route.path} element={route.element} />
          ))}
          <Route
            path="*"
            element={<Navigate to={`/virtualization/projects/${namespace}/vms`} replace />}
          />
        </Routes>
      </Layout>
    </BrowserRouter>
  );
};

class StandaloneAPI implements PluginAPI {
  private routes: RouteConfig[] = [];
  private navItems: NavItem[] = [];

  registerRoute(route: RouteConfig): void {
    this.routes.push(route);
  }

  registerNavigation(item: NavItem): void {
    this.navItems.push(item);
  }

  start(container: HTMLElement) {
    const root = createRoot(container);
    root.render(
      <React.StrictMode>
        <StandaloneShell routes={this.routes} navItems={this.navItems} />
      </React.StrictMode>
    );
  }
}

export const createStandaloneConsole = () => new StandaloneAPI();

export type StandaloneConsole = ReturnType<typeof createStandaloneConsole>;

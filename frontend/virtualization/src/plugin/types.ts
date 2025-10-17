import type { ReactElement } from 'react';

export interface RouteConfig {
  path: string;
  element: ReactElement;
}

export interface NavItem {
  id: string;
  label: string;
  pathTemplate: string;
  children?: NavItem[];
}

export interface PluginAPI {
  registerRoute(route: RouteConfig): void;
  registerNavigation(item: NavItem): void;
}

import React from 'react';
import {
  AppBar,
  Box,
  CssBaseline,
  Drawer,
  List,
  ListItemButton,
  ListItemText,
  TextField,
  Toolbar,
  Typography,
} from '@mui/material';
import { NavLink } from 'react-router-dom';

export interface NavigationItem {
  id: string;
  label: string;
  path: string;
  children?: NavigationItem[];
}

interface LayoutProps {
  navItems: NavigationItem[];
  namespace: string;
  onNamespaceChange: (ns: string) => void;
  children: React.ReactNode;
}

const drawerWidth = 240;

const Layout: React.FC<LayoutProps> = ({ navItems, namespace, onNamespaceChange, children }) => (
  <Box sx={{ display: 'flex', minHeight: '100vh' }}>
    <CssBaseline />
    <AppBar position="fixed" sx={{ zIndex: theme => theme.zIndex.drawer + 1 }}>
      <Toolbar>
        <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
          KubeSphere Virtualization
        </Typography>
        <TextField
          size="small"
          label="Namespace"
          value={namespace}
          onChange={event => onNamespaceChange(event.target.value)}
        />
      </Toolbar>
    </AppBar>
    <Drawer
      variant="permanent"
      sx={{
        width: drawerWidth,
        flexShrink: 0,
        [`& .MuiDrawer-paper`]: { width: drawerWidth, boxSizing: 'border-box', pt: 8 },
      }}
    >
      <Toolbar />
      <Box sx={{ overflow: 'auto' }}>
        <List>
          {navItems.map(item => (
            <React.Fragment key={item.id}>
              <ListItemButton
                component={NavLink}
                to={item.path}
                sx={theme => ({ '&.active': { backgroundColor: theme.palette.action.selected } })}
              >
                <ListItemText primary={item.label} />
              </ListItemButton>
              {item.children?.map(child => (
                <ListItemButton
                  key={child.id}
                  component={NavLink}
                  to={child.path}
                  sx={theme => ({
                    pl: 4,
                    '&.active': { backgroundColor: theme.palette.action.selected },
                  })}
                >
                  <ListItemText primary={child.label} />
                </ListItemButton>
              ))}
            </React.Fragment>
          ))}
        </List>
      </Box>
    </Drawer>
    <Box component="main" sx={{ flexGrow: 1, p: 3, mt: 8, backgroundColor: '#f7f9fb' }}>
      {children}
    </Box>
  </Box>
);

export default Layout;

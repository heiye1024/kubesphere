import { registerVirtualizationPlugin } from './plugin/register';
import { createStandaloneConsole } from './standalone/console';

const container = document.getElementById('root');
if (container) {
  const consoleApp = createStandaloneConsole();
  registerVirtualizationPlugin(consoleApp);
  consoleApp.start(container);
}

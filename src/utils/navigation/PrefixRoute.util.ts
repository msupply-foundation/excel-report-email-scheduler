// Prefixes the route with the base URL of the plugin
import { PLUGIN_BASE_URL } from '../../constants';

function prefixRoute(route: string): string {
  return `${PLUGIN_BASE_URL}/${route}`;
}

export { prefixRoute };

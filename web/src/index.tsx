import ReactDOM from 'react-dom/client';
import { Theme } from '@radix-ui/themes';
import '@radix-ui/themes/styles.css';
import "./global.css"
import RouterProvider from './shared/providers/RouterProvider';


const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);

root.render(
  // <React.StrictMode>
    <Theme>
      <RouterProvider />
    </Theme>
  // </React.StrictMode>
);
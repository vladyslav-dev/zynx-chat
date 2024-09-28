import ReactDOM from 'react-dom/client';
import { Theme } from '@radix-ui/themes';
import '@radix-ui/themes/styles.css';
import "./global.css"
import RouterProvider from './shared/providers/RouterProvider';
import { validateEnvVariables } from './shared/lib/envValidation';

console.log("NODE_ENV", process.env.NODE_ENV)

const isEnvValid = validateEnvVariables();

if (!isEnvValid) {
  console.error('Error: Missing or invalid environment variables. Please check your configuration.');
  
  ReactDOM.createRoot(document.getElementById('root')!).render(
      <div style={{ color: 'red', textAlign: 'center', padding: '20px' }}>
        Error: Missing or invalid environment variables. Please contact support.
      </div>
  );
} else {
  ReactDOM.createRoot(document.getElementById('root')!).render(
    <Theme>
      <RouterProvider />
    </Theme>
  );
}
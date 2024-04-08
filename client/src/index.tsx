import React from 'react';
import ReactDOM from 'react-dom/client';
import './app/styles/global.css';
import AuthContextProvider from './providers/AuthProvider';
import RouterProvider from './providers/RouterProvider';
import WebsocketProvider from './providers/WebsocketProvider';


const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);
root.render(
  <React.StrictMode>
    <AuthContextProvider>
      <WebsocketProvider>
        <RouterProvider />
      </WebsocketProvider>
    </AuthContextProvider>
  </React.StrictMode>
);
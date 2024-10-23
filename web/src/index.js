import React from 'react';
import ReactDOM from 'react-dom/client';
import 'semantic-ui-css/semantic.min.css';  
import { UserProvider } from './context/User';
import { ProjectProvider } from './context/ProjectContext';
import { ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import { BrowserRouter } from 'react-router-dom';
import App from './App';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
    <ProjectProvider>
      <UserProvider>
        <BrowserRouter>
          <App />   
          <ToastContainer />
        </BrowserRouter>
      </UserProvider>
    </ProjectProvider>
);

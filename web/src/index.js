import React from 'react';
import ReactDOM from 'react-dom/client';
import 'semantic-ui-css/semantic.min.css';  
import { AuthProvider } from './context/AuthContext';
import { ProjectProvider } from './context/ProjectContext';
import { UserProvider } from './context/UserContext';
import { ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import { BrowserRouter } from 'react-router-dom';
import App from './App';
import react from 'react';

console.log(react.version)

const root = ReactDOM.createRoot(document.getElementById('root'));

root.render(
    <AuthProvider>
      <ProjectProvider>
        <UserProvider>
          <BrowserRouter>
            <App />   
            <ToastContainer />
          </BrowserRouter>
        </UserProvider>
      </ProjectProvider>
    </AuthProvider>
);

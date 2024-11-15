import React, { lazy, Suspense, useEffect } from 'react';
import { Route, Routes, Navigate } from 'react-router-dom';
import Loading from './components/Loading';
import PrivateRoute from './components/PrivateRoute';
import Layout from './components/Layout';
import { API, showError, showNotice } from './helpers';
import { useAuth } from './context/AuthContext';

// Public components
const LoginForm = lazy(() => import('./components/LoginForm'));
const RegisterForm = lazy(() => import('./components/RegisterForm'));
const PasswordResetForm = lazy(() => import('./components/PasswordResetForm'));

// Private components
const Home = lazy(() => import('./pages/Home'));
const About = lazy(() => import('./pages/About'));
const ProjectGroupOrg = lazy(() => import('./pages/Project/ProjectGroupOrg'));
const ProjectManage = lazy(() => import('./pages/Project/ProjectManage'));
const NotebookList = lazy(() => import('./pages/ModelDevelop/NotebookList'));
const ImageManagement = lazy(() => import('./pages/ImageManagement'));
const ModelTraining = lazy(() => import('./pages/Training/ModelTraining'));
const ModelDeployment = lazy(() => import('./pages/ModelDeployment'));
const ServingList = lazy(() => import('./pages/Serving/ServingList'));
const User = lazy(() => import('./pages/User'));
const UserList = lazy(() => import('./pages/User/UserList'));
const Setting = lazy(() => import('./pages/Setting'));
const NotFound = lazy(() => import('./pages/NotFound'));

function App() {
  const { isAuthenticated } = useAuth();

  const loadStatus = async () => {
    try {
      const res = await API.get('/api/status');
      const { success, data } = res.data;
      if (success) {
        console.log(`GitHub 仓库地址：https://github.com/MLcore-Engine/MLcore-Engine`);
        localStorage.setItem('status', JSON.stringify(data));
        if (
          data.version !== process.env.REACT_APP_VERSION &&
          data.version !== 'v0.0.0' &&
          process.env.REACT_APP_VERSION !== ''
        ) {
          showNotice(`新版本可用：${data.version}，请使用快捷键 Shift + F5 刷新页面`);
        }
      } else {
        showError('无法正常连接至服务器！');
      }
    } catch (error) {
      showError('无法正常连接至服务器！');
      console.error(error);
    }
  };

  useEffect(() => {
      loadStatus();
  }, []);

  return (
    <Suspense fallback={<Loading />}>
      <Routes>
        {/* Public routes */}
        <Route path="/login" element={isAuthenticated ? <Navigate to="/" replace /> : <LoginForm />} />
        <Route path="/register" element={isAuthenticated ? <Navigate to="/" replace /> : <RegisterForm />} /> 

        {/* Private routes */}
        <Route element={<Layout />}>
          <Route path="/reset" element={<PrivateRoute><PasswordResetForm /> </PrivateRoute>} />
          <Route path="/" element={<PrivateRoute><Home /></PrivateRoute>} />
          <Route path="/about" element={<PrivateRoute><About /></PrivateRoute>} />
          
          {/* project manage module */}
		      <Route path="/project" element={<Navigate to="/project/project_list" replace />} />
          <Route path="/project/project_list" element={<PrivateRoute><ProjectGroupOrg /></PrivateRoute>} />
          <Route path="/project/project_manage" element={<PrivateRoute><ProjectManage /></PrivateRoute>} />
          <Route path="/project/user_list" element={<PrivateRoute><UserList /></PrivateRoute>} />
          
          {/* online developing module */}
          <Route path="/notebook" element={<Navigate to="/notebook/env-list" replace />} />
          <Route path="/notebook/env-list" element={<PrivateRoute><NotebookList /></PrivateRoute>} />
          
          <Route path="/training" element={<Navigate to="/training/training-list" replace />} />
          <Route path="/training/training-list" element={<PrivateRoute><ModelTraining /></PrivateRoute>} />

          <Route path="/serving" element={<Navigate to="/serving/serving-list" replace />} />
          <Route path="/serving/serving-list" element={<PrivateRoute><ServingList /></PrivateRoute>} />

          <Route path="/image" element={<PrivateRoute><ImageManagement /></PrivateRoute>} />
          <Route path="/train" element={<PrivateRoute><ModelTraining /></PrivateRoute>} />
          <Route path="/deploy" element={<PrivateRoute><ModelDeployment /></PrivateRoute>} />
          <Route path="/user/*" element={<PrivateRoute><User /></PrivateRoute>} />
          
          <Route path="/setting" element={<PrivateRoute><Setting /></PrivateRoute>} />
          <Route path="*" element={<NotFound />} />
        </Route>
      </Routes>
    </Suspense>
  );
}

export default App;
import React, { lazy, Suspense, useEffect } from 'react';
import { Route, Routes, Navigate, useLocation } from 'react-router-dom';
import Loading from './components/Loading';
import PrivateRoute from './components/PrivateRoute';
import Layout from './components/Layout';
import { useUser } from './context/User';
import { API, showError, showNotice } from './helpers';

// Public components
const LoginForm = lazy(() => import('./components/LoginForm'));
const RegisterForm = lazy(() => import('./components/RegisterForm'));
const PasswordResetForm = lazy(() => import('./components/PasswordResetForm'));
const GitHubOAuth = lazy(() => import('./components/GitHubOAuth'));

// Private components
const Home = lazy(() => import('./pages/Home'));
const ProjectGroupOrg = lazy(() => import('./pages/Project/ProjectGroupOrg'));
const ProjectManage = lazy(() => import('./pages/Project/ProjectManage'));
const ImageManagement = lazy(() => import('./pages/ImageManagement'));
const ModelTraining = lazy(() => import('./pages/ModelTraining'));
const ModelDeployment = lazy(() => import('./pages/ModelDeployment'));
const User = lazy(() => import('./pages/User'));
const Setting = lazy(() => import('./pages/Setting'));
const NotFound = lazy(() => import('./pages/NotFound'));

function App() {
  const { isAuthenticated } = useUser();

  const loadStatus = async () => {
    try {
      const res = await API.get('/api/status');
      const { success, data } = res.data;
      if (success) {
        console.log(`GitHub 仓库地址：https://github.com/MLcore-Engine/MLcore-Engine`);
        localStorage.setItem('status', JSON.stringify(data));
        // localStorage.setItem('system_name', data.system_name);
        // localStorage.setItem('footer_html', data.footer_html);
        // localStorage.setItem('home_page_link', data.home_page_link);
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
        {/* <Route path="/oauth/github" element={<GitHubOAuth />} /> */}

        {/* Private routes */}
        <Route element={<Layout />}>
          <Route path="/reset" element={<PrivateRoute><PasswordResetForm /> </PrivateRoute>} />
          <Route path="/" element={<PrivateRoute><Home /></PrivateRoute>} />
		      <Route path="/project" element={<Navigate to="/project/project_list" replace />} />
          <Route path="/project/project_list" element={<PrivateRoute><ProjectGroupOrg /></PrivateRoute>} />
          <Route path="/project/project_manage" element={<PrivateRoute><ProjectManage /></PrivateRoute>} />
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
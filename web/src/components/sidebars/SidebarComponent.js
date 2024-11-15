import React from 'react';
import { useLocation } from 'react-router-dom';
import ProjectManagementSidebar from './ProjectManagementSidebar';
import ModelDevelopmentSidebar from './ModelDevelopSidebar';
import ModelTrainingSidebar from './ModelTrainingSidebar';
import ModelServingSidebar from './ModelServingSidebar';

const SidebarComponent = () => {
  const location = useLocation();

  const renderSidebar = () => {
    if (location.pathname.startsWith('/project/')) {
      return <ProjectManagementSidebar />;
    }
    if (location.pathname.startsWith('/notebook')) {
      return <ModelDevelopmentSidebar />;
    }

    if (location.pathname.startsWith('/training')) {
      return <ModelTrainingSidebar />;
    }

    if (location.pathname.startsWith('/serving')) {
      return <ModelServingSidebar />;
    }
    
    return null;
  };



  const sidebarContent = renderSidebar();

  if (!sidebarContent) {
    return null;
  }

  return (
    <div style={{ width: '250px', height: '100%', overflowY: 'auto' }}>
      {sidebarContent}
    </div>
  );
};

export default SidebarComponent;
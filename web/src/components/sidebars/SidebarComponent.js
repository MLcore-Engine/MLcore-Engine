import React from 'react';
import { useLocation } from 'react-router-dom';
import ProjectManagementSidebar from './ProjectManagementSidebar';
import ModelDevelopmentSidebar from './ModelDevelopmentSidebar';

const SidebarComponent = () => {
  const location = useLocation();

  const renderSidebar = () => {
    if (location.pathname.startsWith('/project/')) {
      return <ProjectManagementSidebar />;
    }
    if (location.pathname.startsWith('/notebook')) {
      return <ModelDevelopmentSidebar />;
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
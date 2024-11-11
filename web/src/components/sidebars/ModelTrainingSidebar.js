import React, { useState } from 'react';
import { Menu, Icon } from 'semantic-ui-react';
import { Link, useLocation } from 'react-router-dom';

const ModelTrainingSidebar = () => {
  const location = useLocation();
  const [openSections, setOpenSections] = useState({
    trainingManage: true,
    registryManage: true,
    ossManage: true,
  });

  const disabledRoutes = {
    '/training/training-list': true,
    '/training/registry-list': true,
    '/training/artifact-list': true,
  };

  const isActive = (path) => location.pathname === path;

  const toggleSection = (sectionKey) => {
    setOpenSections((prevState) => ({
      ...prevState,
      [sectionKey]: !prevState[sectionKey],
    }));
  };

  const MenuItem = ({ to, icon, children }) => {
    const isDisabled = disabledRoutes[to];
    
    if (isDisabled) {
      return (
        <Menu.Item 
          style={{ 
            cursor: 'not-allowed',
            opacity: 0.45
          }}
          title="功能未开放"
        >
          <Icon name={icon} />
          {children}
        </Menu.Item>
      );
    }
    return (
      <Menu.Item 
        as={Link} 
        to={to} 
        active={isActive(to)}
      >
        <Icon name={icon} />
        {children}
      </Menu.Item>
    );
  };


  return (
    <Menu vertical fluid>
      <Menu.Item>
        <Menu.Header onClick={() => toggleSection('trainingManage')} style={{ cursor: 'pointer' }}>
          <Icon name={openSections.notebookManage ? 'angle down' : 'angle right'} />
          TrainingJob
        </Menu.Header>
        {openSections.trainingManage && (
          <Menu.Menu>
            <Menu.Item as={Link} to="/training/training-list" active={isActive('/training/training-list')}>
              <Icon name="server" />
              TrainingList
            </Menu.Item>
          </Menu.Menu>
        )}
      </Menu.Item>

      <Menu.Item>
        <Menu.Header onClick={() => toggleSection('registryManage')} style={{ cursor: 'pointer' }}>
          <Icon name={openSections.registryManage ? 'angle down' : 'angle right'} />
          仓库管理
        </Menu.Header>
        {openSections.registryManage && (
          <Menu.Menu>
            <MenuItem to="/training/registry-list" icon="cube">
              仓库列表
            </MenuItem>
          </Menu.Menu>
        )}
      </Menu.Item>

      <Menu.Item>
        <Menu.Header onClick={() => toggleSection('ossManage')} style={{ cursor: 'pointer' }}>
          <Icon name={openSections.ossManage ? 'angle down' : 'angle right'} />
          制品管理
        </Menu.Header>
        {openSections.ossManage && (
          <Menu.Menu>
            <MenuItem to="/training/artifact-list" icon="plus square">
              制品列表
            </MenuItem>
          </Menu.Menu>
        )}
      </Menu.Item>

    </Menu>
  );
};

export default ModelTrainingSidebar;   


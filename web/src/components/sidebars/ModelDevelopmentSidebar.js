import React, { useState } from 'react';
import { Menu, Icon } from 'semantic-ui-react';
import { Link, useLocation } from 'react-router-dom';

const ModelDevelopmentSidebar = () => {
  const location = useLocation();
  const [openSections, setOpenSections] = useState({
    imageManage: true,
    notebookManage: true,
    modelManage: true,
  });

  const disabledRoutes = {
    "/notebook/workspace": true,
    '/notebook/settings': true,
    '/notebook/model-list': true,
    '/notebook/model-version': true,
    '/notebook/model-metrics': true,
    '/notebook/image-create': true,
    '/notebook/image-list': true,
    '/notebook/image-registry': true,
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
        <Menu.Header onClick={() => toggleSection('notebookManage')} style={{ cursor: 'pointer' }}>
          <Icon name={openSections.notebookManage ? 'angle down' : 'angle right'} />
          Notebook
        </Menu.Header>
        {openSections.notebookManage && (
          <Menu.Menu>
            <Menu.Item as={Link} to="/notebook/env-list" active={isActive('/notebook/env-list')}>
              <Icon name="server" />
              环境列表
            </Menu.Item>
            <MenuItem to="/notebook/workspace" icon="code">
              工作空间
            </MenuItem>
            <MenuItem to="/notebook/settings" icon="setting">
              环境配置
            </MenuItem>
          </Menu.Menu>
        )}
      </Menu.Item>

      <Menu.Item>
        <Menu.Header onClick={() => toggleSection('modelManage')} style={{ cursor: 'pointer' }}>
          <Icon name={openSections.modelManage ? 'angle down' : 'angle right'} />
          模型管理
        </Menu.Header>
        {openSections.modelManage && (
          <Menu.Menu>
            <MenuItem to="/notebook/model-list" icon="cube">
              模型列表
            </MenuItem>
            <MenuItem to="/notebook/model-version" icon="history">
              版本管理
            </MenuItem>
            <MenuItem to="/notebook/model-metrics" icon="chart line">
              模型指标
            </MenuItem>
          </Menu.Menu>
        )}
      </Menu.Item>

      <Menu.Item>
        <Menu.Header onClick={() => toggleSection('imageManage')} style={{ cursor: 'pointer' }}>
          <Icon name={openSections.imageManage ? 'angle down' : 'angle right'} />
          镜像管理
        </Menu.Header>
        {openSections.imageManage && (
          <Menu.Menu>
            <MenuItem to="/notebook/image-create" icon="plus square">
              镜像构建
            </MenuItem>
            <MenuItem to="/notebook/image-list" icon="list">
              镜像列表
            </MenuItem>
            <MenuItem to="/notebook/image-registry" icon="database">
              镜像仓库
            </MenuItem>
          </Menu.Menu>
        )}
      </Menu.Item>

    </Menu>
  );
};

export default ModelDevelopmentSidebar;

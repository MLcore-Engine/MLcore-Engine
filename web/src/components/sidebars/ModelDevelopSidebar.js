import React, { useState } from 'react';
import { Menu, Icon } from 'semantic-ui-react';
import { Link, useLocation } from 'react-router-dom';


const applyStyles = (element) => {
  return React.cloneElement(element, {
    style: {
      marginLeft: '1.2em',
      display: 'flex',
    },
    children: React.Children.map(element.props.children, (child) => {
      if (child.type === Icon) {
        // 这里直接在 Icon 元素上应用样式
        return React.cloneElement(child, {
          style: { marginRight: '8px', ...child.props.style }, 
        });
      }
      // 对非 Icon 元素递归应用样式
      if (child.props) {
        return applyStyles(child);
      }
      return child;
    }),
  });
};


const ModelDevelopmentSidebar = () => {
  const location = useLocation();
  const [openSections, setOpenSections] = useState({
    imageManage: true,
    notebookManage: true,
    modelManage: true,
  });

  const disabledRoutes = {
    '/notebook/workspace': true,
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

    const menuItemContent = (
      <>
        <Icon name={icon} />
        {children}
      </>
    );

    if (isDisabled) {
      return applyStyles(
        <Menu.Item
          style={{ cursor: 'not-allowed', opacity: 0.45 }}
          title='功能未开放'
        >
          {menuItemContent}
        </Menu.Item>
      );
    }

    return applyStyles(
      <Menu.Item as={Link} to={to} active={isActive(to)}>
        {menuItemContent}
      </Menu.Item>
    );
  };

  return (
    <Menu vertical fluid style={{ height: '100%' }}>
      <Menu.Item>
        <Menu.Header
          onClick={() => toggleSection('notebookManage')}
          style={{ cursor: 'pointer' }}
        >
          <Icon
            name={openSections.notebookManage ? 'angle down' : 'angle right'}
          />
          Notebook
        </Menu.Header>
        {openSections.notebookManage && (
          <Menu.Menu>
            {applyStyles(
              // <Menu.Item as={Link} to="/notebook/env-list" active={isActive('/notebook/env-list')} >
              //   <Icon name="server" />
              //   环境列表
              // </Menu.Item>

              <MenuItem to='/notebook/env-list' icon='server'>
                工作空间
              </MenuItem>
            )}
            {applyStyles(
              <MenuItem to='/notebook/workspace' icon='code'>
                工作空间
              </MenuItem>
            )}
            {applyStyles(
              <MenuItem to='/notebook/settings' icon='setting'>
                环境配置
              </MenuItem>
            )}
          </Menu.Menu>
        )}
      </Menu.Item>

      <Menu.Item>
        <Menu.Header
          onClick={() => toggleSection('modelManage')}
          style={{ cursor: 'pointer' }}
        >
          <Icon
            name={openSections.modelManage ? 'angle down' : 'angle right'}
          />
          模型管理
        </Menu.Header>
        {openSections.modelManage && (
          <Menu.Menu>
            {applyStyles(
              <MenuItem to='/notebook/model-list' icon='cube'>
                模型列表
              </MenuItem>
            )}
            {applyStyles(
              <MenuItem to='/notebook/model-version' icon='history'>
                版本管理
              </MenuItem>
            )}
            {applyStyles(
              <MenuItem to='/notebook/model-metrics' icon='chart line'>
                模型指标
              </MenuItem>
            )}
          </Menu.Menu>
        )}
      </Menu.Item>

      <Menu.Item>
        <Menu.Header
          onClick={() => toggleSection('imageManage')}
          style={{ cursor: 'pointer' }}
        >
          <Icon
            name={openSections.imageManage ? 'angle down' : 'angle right'}
          />
          镜像管理
        </Menu.Header>
        {openSections.imageManage && (
          <Menu.Menu>
            {applyStyles(
              <MenuItem to='/notebook/image-create' icon='plus square'>
                镜像构建
              </MenuItem>
            )}
            {applyStyles(
              <MenuItem to='/notebook/image-list' icon='list'>
                镜像列表
              </MenuItem>
            )}
            {applyStyles(
              <MenuItem to='/notebook/image-registry' icon='database'>
                镜像仓库
              </MenuItem>
            )}
          </Menu.Menu>
        )}
      </Menu.Item>
    </Menu>
  );
};

export default ModelDevelopmentSidebar;

import React, { useState } from 'react';
import { Menu, Icon } from 'semantic-ui-react';
import { Link, useLocation } from 'react-router-dom';

const ProjectManagementSidebar = () => {
  const location = useLocation();
  const [openSections, setOpenSections] = useState({
    projectGroup: true,
    userManagement: true,
    reousrceManage: true,
  });

  const isActive = (path) => location.pathname === path;

  const toggleSection = (sectionKey) => {
    setOpenSections((prevState) => ({
      ...prevState,
      [sectionKey]: !prevState[sectionKey],
    }));
  };

  const menuHeaderStyle = {
    cursor: 'pointer',
    padding: '0.8rem 1rem',
    fontWeight: 600,
    color: '#333',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    borderRadius: '0.5rem',
    margin: '0.3rem 0.5rem',
    transition: 'all 0.2s ease'
  };

  const menuItemStyle = {
    borderRadius: '0.375rem',
    margin: '0.2rem 0.8rem',
    padding: '0.6rem 0.8rem',
    fontSize: '0.95rem',
    transition: 'all 0.2s ease',
    display: 'flex',
    alignItems: 'center'
  };

  const activeMenuItemStyle = {
    ...menuItemStyle,
    backgroundColor: 'rgba(19, 62, 135, 0.08)',
    color: '#133e87',
    fontWeight: 500
  };

  return (
    <div style={{ padding: '0.5rem 0' }}>
      <div style={{ paddingBottom: '1rem', borderBottom: '1px solid rgba(0,0,0,0.05)' }}>
        <div
          onClick={() => toggleSection('projectGroup')}
          style={menuHeaderStyle}
        >
          <span>
            <Icon name="folder" style={{ marginRight: '0.5rem' }} />
            项目组
          </span>
          <Icon name={openSections.projectGroup ? 'angle down' : 'angle right'} />
        </div>
        {openSections.projectGroup && (
          <Menu.Menu>
            <Menu.Item
              as={Link}
              to="/project/project_list"
              active={isActive('/project/project_list')}
              style={isActive('/project/project_list') ? activeMenuItemStyle : menuItemStyle}
            >
              <Icon name="folder open" style={{ marginRight: '0.5rem' }} />
              项目列表
            </Menu.Item>
            <Menu.Item
              as={Link}
              to="/project/project_manage"
              active={isActive('/project/project_manage')}
              style={isActive('/project/project_manage') ? activeMenuItemStyle : menuItemStyle}
            >
              <Icon name="sitemap" style={{ marginRight: '0.5rem' }} />
              项目管理
            </Menu.Item>
          </Menu.Menu>
        )}
      </div>

      <div style={{ padding: '1rem 0', borderBottom: '1px solid rgba(0,0,0,0.05)' }}>
        <div
          onClick={() => toggleSection('userManagement')}
          style={menuHeaderStyle}
        >
          <span>
            <Icon name="user" style={{ marginRight: '0.5rem' }} />
            用户管理
          </span>
          <Icon name={openSections.userManagement ? 'angle down' : 'angle right'} />
        </div>
        {openSections.userManagement && (
          <Menu.Menu>
            <Menu.Item
              as={Link}
              to="/project/user_list"
              active={isActive('/project/user_list')}
              style={isActive('/project/user_list') ? activeMenuItemStyle : menuItemStyle}
            >
              <Icon name="users" style={{ marginRight: '0.5rem' }} />
              用户列表
            </Menu.Item>
            <Menu.Item
              as={Link}
              to="/project/user_role"
              active={isActive('/project/user_role')}
              style={isActive('/project/user_role') ? activeMenuItemStyle : menuItemStyle}
            >
              <Icon name="user circle" style={{ marginRight: '0.5rem' }} />
              用户角色
            </Menu.Item>
          </Menu.Menu>
        )}
      </div>

      <div style={{ paddingTop: '1rem' }}>
        <div
          onClick={() => toggleSection('reousrceManage')}
          style={menuHeaderStyle}
        >
          <span>
            <Icon name="cubes" style={{ marginRight: '0.5rem' }} />
            资源管理
          </span>
          <Icon name={openSections.reousrceManage ? 'angle down' : 'angle right'} />
        </div>
        {openSections.reousrceManage && (
          <Menu.Menu>
            <Menu.Item
              as={Link}
              to="/project/resource_gauge"
              active={isActive('/project/resource_gauge')}
              style={isActive('/project/resource_gauge') ? activeMenuItemStyle : menuItemStyle}
            >
              <Icon name="chart bar" style={{ marginRight: '0.5rem' }} />
              资源统计
            </Menu.Item>
            <Menu.Item
              as={Link}
              to="/project/project_alloc"
              active={isActive('/project/project_alloc')}
              style={isActive('/project/project_alloc') ? activeMenuItemStyle : menuItemStyle}
            >
              <Icon name="cubes" style={{ marginRight: '0.5rem' }} />
              资源分配
            </Menu.Item>
          </Menu.Menu>
        )}
      </div>
    </div>
  );
};

export default ProjectManagementSidebar;

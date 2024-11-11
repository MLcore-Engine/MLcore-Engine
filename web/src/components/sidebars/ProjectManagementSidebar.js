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

  return (
    <Menu vertical fluid>
      <Menu.Item>
        <Menu.Header onClick={() => toggleSection('projectGroup')} style={{ cursor: 'pointer' }}>
          <Icon name={openSections.projectGroup ? 'angle down' : 'angle right'} />
          项目组
        </Menu.Header>
        {openSections.projectGroup && (
          <Menu.Menu>
            <Menu.Item as={Link} to="/project/project_list" active={isActive('/project/project_list')}>
              <Icon name="folder open" />
              项目列表
            </Menu.Item>
            <Menu.Item as={Link} to="/project/project_manage" active={isActive('/project/project_manage')}>
              <Icon name="sitemap" />
              项目管理
            </Menu.Item>
          </Menu.Menu>
        )}
      </Menu.Item>

      <Menu.Item>
        <Menu.Header onClick={() => toggleSection('userManagement')} style={{ cursor: 'pointer' }}>
          <Icon name={openSections.userManagement ? 'angle down' : 'angle right'} />
          用户管理
        </Menu.Header>
        {openSections.userManagement && (
          <Menu.Menu>
            <Menu.Item as={Link} to="/project/user_list" active={isActive('/project/user_list')}>
              <Icon name="users" />
              用户列表
            </Menu.Item>
            <Menu.Item as={Link} to="/project/user_role" active={isActive('/project/user_role')}>
              <Icon name="user circle" />
              用户角色
            </Menu.Item>
          </Menu.Menu>
        )}
      </Menu.Item>


      <Menu.Item>
        <Menu.Header onClick={() => toggleSection('reousrceManage')} style={{ cursor: 'pointer' }}>
          <Icon name={openSections.reousrceManage ? 'angle down' : 'angle right'} />
          资源管理
        </Menu.Header>
        {openSections.reousrceManage && (
          <Menu.Menu>
            <Menu.Item as={Link} to="/project/resource_gauge" active={isActive('/project/resource_gauge')}>
              <Icon name="chart bar" />
              资源统计
            </Menu.Item>
            <Menu.Item as={Link} to="/project/project_alloc" active={isActive('/project/project_alloc')}>
              <Icon name="cubes" /> 
              资源分配
            </Menu.Item>
          </Menu.Menu>
        )}
      </Menu.Item>
    </Menu>
  );
};

export default ProjectManagementSidebar;
import React from 'react';
import { Menu, Icon } from 'semantic-ui-react';
import { Link, useLocation } from 'react-router-dom';
import { useProjects } from '../../context/ProjectContext';

const ProjectManagementSidebar = () => {
  const location = useLocation();
  const { projects = [] } = useProjects(); // 提供默认空数组

  const isActive = (path) => location.pathname === path;

  return (
    <Menu vertical fluid className="h-full">
      {/* 项目组部分 */}
      <Menu.Item header>项目组</Menu.Item>
      <Menu.Item
        as={Link}
        to="/project/project_list"
        active={isActive('/project/project_list')}
      >
        <Icon name="folder open" />
        项目列表
      </Menu.Item>
      <Menu.Item
        as={Link}
        to="/project/project_manage"
        active={isActive('/project/project_manage')}
      >
        <Icon name="sitemap" />
        项目管理
      </Menu.Item>

      {/* 动态项目列表 */}
      {projects.length > 0 ? (
        projects.map((project) => (
          <Menu.Item
            key={project.id}
            as={Link}
            to={`/project/${project.id}`}
            active={isActive(`/project/${project.id}`)}
          >
            <Icon name="folder" />
            {project.name}
          </Menu.Item>
        ))
      ) : (
        <Menu.Item disabled>
          <Icon name="info circle" />
          无项目可显示
        </Menu.Item>
      )}

      {/* 用户管理部分 */}
      <Menu.Item header>用户管理</Menu.Item>
      <Menu.Item
        as={Link}
        to="/project/user_list"
        active={isActive('/project/user_list')}
      >
        <Icon name="users" />
        用户列表
      </Menu.Item>
      <Menu.Item
        as={Link}
        to="/project/user_role"
        active={isActive('/project/user_role')}
      >
        <Icon name="user circle" />
        用户角色
      </Menu.Item>
    </Menu>
  );
};

export default ProjectManagementSidebar;

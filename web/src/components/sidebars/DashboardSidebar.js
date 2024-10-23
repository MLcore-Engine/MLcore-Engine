import React from 'react';
import { Menu, Icon } from 'semantic-ui-react';
import { Link } from 'react-router-dom';

const DashboardSidebar = () => (
  <Menu vertical inverted>
    <Menu.Item as={Link} to="/dashboard">
      <Icon name="dashboard" />
      仪表板
    </Menu.Item>
    <Menu.Item as={Link} to="/project">
      <Icon name="folder open" />
      项目
    </Menu.Item>
    <Menu.Item as={Link} to="/notebook">
      <Icon name="code" />
      笔记本
    </Menu.Item>
    {/* 可以添加更多相关的菜单项 */}
  </Menu>
);

export default DashboardSidebar;
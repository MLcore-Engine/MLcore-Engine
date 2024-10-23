import React from 'react';
import { Menu, Icon } from 'semantic-ui-react';
import { Link } from 'react-router-dom';

const ModelDevelopmentSidebar = () => (
  <Menu vertical inverted>
    <Menu.Item as={Link} to="/notebook">
      <Icon name="code" />
      笔记本
    </Menu.Item>
    <Menu.Item as={Link} to="/notebook/create">
      <Icon name="plus" />
      创建模型
    </Menu.Item>
    {/* Add more model development items */}
  </Menu>
);

export default ModelDevelopmentSidebar;

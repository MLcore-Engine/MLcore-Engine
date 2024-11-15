import React from 'react';
import { Menu, Icon } from 'semantic-ui-react';
import { Link } from 'react-router-dom';

const ModelDeploymentSidebar = () => (
  <Menu vertical inverted>
    <Menu.Item as={Link} to="/deploy">
      <Icon name="cloud" />
      部署列表
    </Menu.Item>
    <Menu.Item as={Link} to="/deploy/create">
      <Icon name="plus" />
      创建部署
    </Menu.Item>
    {/* Add more model deployment items */}
  </Menu>
);

export default ModelDeploymentSidebar;

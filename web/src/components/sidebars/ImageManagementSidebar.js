import React from 'react';
import { Menu, Icon } from 'semantic-ui-react';
import { Link } from 'react-router-dom';

const ImageManagementSidebar = () => (
    
  <Menu vertical inverted>
    <Menu.Item as={Link} to="/image">
      <Icon name="image" />
      镜像列表
    </Menu.Item>
    <Menu.Item as={Link} to="/image/create">
      <Icon name="plus" />
      添加镜像
    </Menu.Item>
    {/* Add more image management items */}
  </Menu>
);

export default ImageManagementSidebar;

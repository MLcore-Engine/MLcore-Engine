import React from 'react';
import { Menu, Icon } from 'semantic-ui-react';
import { Link } from 'react-router-dom';

const DefaultSidebar = () => (
  <Menu vertical inverted>
    <Menu.Item as={Link} to="/">
      <Icon name="home" />
      首页
    </Menu.Item>
    {/* Add more default items if necessary */}
  </Menu>
);

export default DefaultSidebar;

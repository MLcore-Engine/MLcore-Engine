import React from 'react';
import { Menu, Icon } from 'semantic-ui-react';
import { Link } from 'react-router-dom';

const ModelTrainingSidebar = () => (
  <Menu vertical inverted>
    <Menu.Item as={Link} to="/train">
      <Icon name="cogs" />
      训练任务
    </Menu.Item>
    <Menu.Item as={Link} to="/train/create">
      <Icon name="plus" />
      创建训练
    </Menu.Item>
    {/* Add more model training items */}
  </Menu>
);

export default ModelTrainingSidebar;

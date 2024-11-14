import React, { useState } from 'react';
import { Menu, Icon } from 'semantic-ui-react';
import { Link, useLocation } from 'react-router-dom';

const ModelServingSidebar = () => {
  const location = useLocation();
  const [openSections, setOpenSections] = useState({
    tritonManage: true,
    tensorflowManage: true,
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
        <Menu.Header onClick={() => toggleSection('tritonManage')} style={{ cursor: 'pointer' }}>
          <Icon name={openSections.tritonManage ? 'angle down' : 'angle right'} />
          Serving-List
        </Menu.Header>
        {openSections.tritonManage && (
          <Menu.Menu>
            <Menu.Item as={Link} to="/serving/serving-list" active={isActive('/serving/serving-list')}>
              <Icon name="server" />
              Serving-List
            </Menu.Item>
          </Menu.Menu>
        )}
      </Menu.Item>

      <Menu.Item>
        <Menu.Header onClick={() => toggleSection('tensorflowManage')} style={{ cursor: 'pointer' }}>
          <Icon name={openSections.tensorflowManage ? 'angle down' : 'angle right'} />
          Model-List
        </Menu.Header>
        {openSections.tensorflowManage && (
          <Menu.Menu>
            <Menu.Item as={Link} to="/serving/model-list" active={isActive('/serving/model-list')}>
              <Icon name="server" />
              Model-List
            </Menu.Item>
          </Menu.Menu>
        )}
      </Menu.Item>
    </Menu>
  );
};

export default ModelServingSidebar;
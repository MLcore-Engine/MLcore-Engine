import React, { useState } from 'react';
import { Menu, Icon } from 'semantic-ui-react';
import { Link, useLocation } from 'react-router-dom';

const MENU_ITEMS = [
  {
    key: 'tritonManage',
    title: 'Serving-List',
    items: [
      {
        path: '/serving/serving-list',
        icon: 'server',
        label: 'Serving-List'
      }
    ]
  },
  {
    key: 'tensorflowManage',
    title: 'Model-List',
    items: [
      {
        path: '/serving/model-list',
        icon: 'server',
        label: 'Model-List'
      }
    ]
  }
];

const ModelServingSidebar = () => {
  const location = useLocation();
  const [openSections, setOpenSections] = useState({
    tritonManage: true,
    tensorflowManage: true,
  });

  const isActive = (path) => location.pathname === path;

  const handleToggleSection = (sectionKey) => {
    setOpenSections((prevState) => ({
      ...prevState,
      [sectionKey]: !prevState[sectionKey],
    }));
  };

  return (
    <Menu vertical fluid>
      {MENU_ITEMS.map(({ key, title, items }) => (
        <Menu.Item key={key}>
          <Menu.Header 
            onClick={() => handleToggleSection(key)} 
            style={{ cursor: 'pointer' }}
          >
            <Icon name={openSections[key] ? 'angle down' : 'angle right'} />
            {title}
          </Menu.Header>
          {openSections[key] && (
            <Menu.Menu>
              {items.map(({ path, icon, label }) => (
                <Menu.Item 
                  key={path}
                  as={Link} 
                  to={path} 
                  active={isActive(path)}
                >
                  <Icon name={icon} />
                  {label}
                </Menu.Item>
              ))}
            </Menu.Menu>
          )}
        </Menu.Item>
      ))}
    </Menu>
  );
};

export default ModelServingSidebar;
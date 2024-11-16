import React, { useState } from 'react';
import { Menu, Icon } from 'semantic-ui-react';
import { Link, useLocation } from 'react-router-dom';

// 修改 applyStyles 函数
const applyStyles = (element) => {
  return React.cloneElement(element, {
    style: {
      marginLeft: '1.2em',
      display: 'flex',
    },
    children: React.Children.map(element.props.children, (child) => {
      if (child.type === Icon) {
        return React.cloneElement(child, {
          style: { marginRight: '8px', ...child.props.style }, // 给 Icon 加上 margin-right
        });
      }
      return child;
    }),
  });
};

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
    <Menu vertical fluid style={{ height: '100%' }}>
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
                applyStyles(
                  <Menu.Item 
                    key={path}
                    as={Link} 
                    to={path} 
                    active={isActive(path)}
                  >
                    <Icon name={icon} />
                    {label}
                  </Menu.Item>
                )
              ))}
            </Menu.Menu>
          )}
        </Menu.Item>
      ))}
    </Menu>
  );
};

export default ModelServingSidebar;

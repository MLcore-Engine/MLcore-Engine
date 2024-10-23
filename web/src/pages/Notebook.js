// web/src/components/Notebook.js
import React, { useState } from 'react';
import { Menu, Sidebar, Segment, Icon } from 'semantic-ui-react';
// import './Notebook.css'; // 自定义样式

function Notebook() {
  const [visible, setVisible] = useState(true);
  const [activeItem, setActiveItem] = useState('list');

  const handleItemClick = (e, { name }) => {
    setActiveItem(name);
  };

  const getIframeSrc = () => {
    switch (activeItem) {
      case 'create':
        return '/notebook/create';
      case 'reset':
        return '/notebook/reset';
      case 'list':
      default:
        return '/notebook/list';
    }
  };

  return (
    <div className="notebook-container">
      <Sidebar.Pushable as={Segment} className="notebook-pushable">
        <Sidebar
          as={Menu}
          animation='overlay'
          icon='labeled'
          inverted
          vertical
          visible={visible}
          width='thin'
        >
          <Menu.Item
            name='create'
            active={activeItem === 'create'}
            onClick={handleItemClick}
          >
            <Icon name='plus' />
            Create
          </Menu.Item>
          <Menu.Item
            name='reset'
            active={activeItem === 'reset'}
            onClick={handleItemClick}
          >
            <Icon name='redo' />
            Reset
          </Menu.Item>
          <Menu.Item
            name='list'
            active={activeItem === 'list'}
            onClick={handleItemClick}
          >
            <Icon name='list' />
            List
          </Menu.Item>
        </Sidebar>

        <Sidebar.Pusher dimmed={false}>
          <Segment basic className="notebook-content">
            <iframe
              src={getIframeSrc()}
              title="Notebook Content"
              width="100%"
              height="800px"
              frameBorder="0"
            ></iframe>
          </Segment>
        </Sidebar.Pusher>
      </Sidebar.Pushable>
    </div>
  );
}

export default Notebook;
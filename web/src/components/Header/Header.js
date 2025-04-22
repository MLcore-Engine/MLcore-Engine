import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { GIT_REPO_URL } from '../../constants/common.constant';
import { Menu, Dropdown, Icon, Container } from 'semantic-ui-react';
import { getSystemName, isAdmin } from '../../helpers';
import 'semantic-ui-css/semantic.min.css';
import '../../styles/header.css';


const headerButtons = [
  {
    name: '项目管理',
    to: '/project',
    icon: 'folder',
    root: true,
  },
  {
    name: '模型开发',
    to: '/notebook',
    icon: 'code',
  },
  {
    name: '镜像管理',
    to: '/image',
    icon: 'image',
  },
  {
    name: '模型训练',
    to: '/training',
    icon: 'cogs',
  },
  {
    name: '模型服务',
    to: '/serving',
    icon: 'cloud',
  },
  {
    name: '关于',
    to: '/about',
    icon: 'info circle',
  },
];

const Header = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const systemName = getSystemName();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <Menu 
      fixed="top" 
      style={{ 
        zIndex: 1000, 
        height: '60px', 
        background: 'white', 
        boxShadow: '0 2px 10px rgba(0, 0, 0, 0.06)',
        borderRadius: 0,
        margin: 0,
        padding: '0 1rem'
      }}
    >
      <Container>
        <Menu.Item 
          as={Link} 
          to="/" 
          header
          style={{
            fontWeight: 700,
            fontSize: '1.2rem',
            padding: '0 1rem'
          }}
        >
          <img 
            src="/favicon.ico" 
            alt="logo" 
            style={{ 
              marginRight: '0.8rem', 
              width: '28px', 
              height: '28px' 
            }} 
          />
          {systemName}
        </Menu.Item>

        {headerButtons.map((button) => {
          if (button.root && !isAdmin()) return null;
          return (
            <Menu.Item 
              as={Link} 
              to={button.to} 
              key={button.name}
              style={{
                fontSize: '0.95rem',
                fontWeight: 500
              }}
            >
              <Icon name={button.icon} style={{ marginRight: '0.4rem' }} />
              {button.name}
            </Menu.Item>
          );
        })}

        <Menu.Menu position="right">
          <Menu.Item
            as="a"
            href={GIT_REPO_URL}
            target="_blank"
            rel="noopener noreferrer"
            style={{ fontWeight: 500 }}
          >
            <Icon name="github" />
            GitHub
          </Menu.Item>

          {user ? (
            <Dropdown 
              item 
              text={user.username}
              style={{ 
                fontWeight: 500,
                borderLeft: '1px solid rgba(0,0,0,0.08)',
                marginLeft: '0.5rem',
                paddingLeft: '1rem'
              }}
            >
              <Dropdown.Menu style={{ borderRadius: '0.5rem', boxShadow: '0 4px 12px rgba(0, 0, 0, 0.1)' }}>
                <Dropdown.Item onClick={handleLogout}>
                  <Icon name="sign-out" />
                  注销
                </Dropdown.Item>
              </Dropdown.Menu>
            </Dropdown>
          ) : (
            <Menu.Item 
              as={Link} 
              to="/login"
              style={{ fontWeight: 500 }}
            >
              <Icon name="sign-in" />
              登录
            </Menu.Item>
          )}
        </Menu.Menu>
      </Container>
    </Menu>
  );
};

export default Header;
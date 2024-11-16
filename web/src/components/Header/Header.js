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
    <Menu fixed="top" style={{ zIndex: 1000 }}>
      <Container>
        <Menu.Item as={Link} to="/" header>
          <img src="/favicon.ico" alt="logo" style={{ marginRight: '1em' }} />
          {systemName}
        </Menu.Item>

        {headerButtons.map((button) => {
          if (button.root && !isAdmin()) return null;
          return (
            <Menu.Item as={Link} to={button.to} key={button.name}>
              <Icon name={button.icon} />
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
          >
            <Icon name="github" />
            GitHub
          </Menu.Item>

          {user ? (
            <Dropdown item text={user.username}>
              <Dropdown.Menu>
                <Dropdown.Item onClick={handleLogout}>
                  <Icon name="sign-out" />
                  注销
                </Dropdown.Item>
              </Dropdown.Menu>
            </Dropdown>
          ) : (
            <Menu.Item as={Link} to="/login">
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
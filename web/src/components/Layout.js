import React from 'react';
import { Outlet, useLocation } from 'react-router-dom';
import { Container } from 'semantic-ui-react';
import Header from './Header/Header';
import Footer from './Footer/Footer';
import SidebarComponent from './sidebars/SidebarComponent';

const Layout = () => {
  const location = useLocation();
  const isHomePage = ['/', '/dashboard', '/about'].includes(location.pathname);
  
  const hasSidebar = (
    location.pathname.startsWith('/project/') ||
    location.pathname.startsWith('/notebook') ||
    location.pathname.startsWith('/training') ||
    location.pathname.startsWith('/serving')
  );

  return (
    <div style={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
      <Header />
      <div style={{ display: 'flex', flex: 1, marginTop: '60px' }}>
        {hasSidebar && <SidebarComponent />}
        <div style={{ flex: 1, overflowX: 'hidden' }}>
          <Container style={{
            paddingTop: '2em',
            paddingBottom: '2em',
            minHeight: isHomePage ? 'calc(100vh - 180px)' : 'calc(100vh - 120px)'
          }}>
            <Outlet />
          </Container>
        </div>
      </div>
      {isHomePage && <Footer />}
    </div>
  );
};

export default Layout;
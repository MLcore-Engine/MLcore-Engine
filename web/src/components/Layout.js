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
    <div style={{ 
      display: 'flex', 
      flexDirection: 'column', 
      minHeight: '100vh',
      backgroundColor: '#f8f9fa'
    }}>
      <Header />
      <div style={{ 
        display: 'flex', 
        flex: 1, 
        marginTop: '60px',
        gap: '1rem',
        padding: '0.5rem'
      }}>
        {hasSidebar && (
          <div style={{
            width: '250px',
            transition: 'all 0.3s ease',
            borderRadius: '0.75rem',
            backgroundColor: '#fff',
            boxShadow: '0 4px 12px rgba(0, 0, 0, 0.05)',
            padding: '1rem 0',
            alignSelf: 'flex-start',
            position: 'sticky',
            top: '70px'
          }}>
            <SidebarComponent />
          </div>
        )}
        <div style={{ 
          flex: 1, 
          overflowX: 'hidden',
          transition: 'all 0.3s ease'
        }}>
          <Container style={{
            padding: '1.5rem',
            backgroundColor: '#fff',
            borderRadius: '0.75rem',
            boxShadow: '0 4px 12px rgba(0, 0, 0, 0.05)',
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
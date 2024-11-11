import React from 'react';
import { Container, Header, Icon, Button } from 'semantic-ui-react';
import { Link } from 'react-router-dom';

const NotFound = () => {
  return (
    <Container text textAlign='center'>
      <Header as='h1' icon>
        <Icon name='search' />
        404 - not found
        <Header.Subheader>
          sorry, the page you requested does not exist.
        </Header.Subheader>
      </Header>
      <Button as={Link} to="/" primary>
        back to home
      </Button>
    </Container>
  );
};

export default NotFound;
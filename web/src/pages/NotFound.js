import React from 'react';
import { Container, Header, Icon, Button } from 'semantic-ui-react';
import { Link } from 'react-router-dom';

const NotFound = () => {
  return (
    <Container text textAlign='center'>
      <Header as='h1' icon>
        <Icon name='search' />
        404 - 页面未找到
        <Header.Subheader>
          抱歉,您请求的页面不存在。
        </Header.Subheader>
      </Header>
      <Button as={Link} to="/" primary>
        返回首页
      </Button>
    </Container>
  );
};

export default NotFound;
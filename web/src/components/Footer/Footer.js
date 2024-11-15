import React from 'react';
import { Container, Segment, Grid, List, Header } from 'semantic-ui-react';
import { getFooterHTML, getSystemName } from '../../helpers';

const Footer = () => {
  const systemName = getSystemName();
  const footer = getFooterHTML();

  return (
    <Segment inverted vertical style={{ padding: '5em 0em' }}>
      <Container>
        <Grid divided inverted stackable>
          <Grid.Row>
            <Grid.Column width={3}>
              <Header inverted as='h4' content='关于' />
              <List link inverted>
                <List.Item as='a' href='https://github.com/MLcore-Engine/MLcore-Engine' target='_blank'>项目主页</List.Item>
                <List.Item as='a' href='https://github.com/MLcore-Engine' target='_blank'>MLcore-Engine</List.Item>
              </List>
            </Grid.Column>
            <Grid.Column width={3}>
              <Header inverted as='h4' content='服务' />
              <List link inverted>
                <List.Item as='a'>文档</List.Item>
                <List.Item as='a'>支持</List.Item>
              </List>
            </Grid.Column>
            <Grid.Column width={7}>
              <Header as='h4' inverted>
                {systemName}
              </Header>
              {footer ? (
                <div dangerouslySetInnerHTML={{ __html: footer }}></div>
              ) : (
                <p>
                  {systemName} {process.env.REACT_APP_VERSION} 由 MLcore-Engine 构建，
                  源代码遵循 <a href='https://opensource.org/licenses/mit-license.php' target='_blank' rel="noopener noreferrer">MIT 协议</a>
                </p>
              )}
            </Grid.Column>
          </Grid.Row>
        </Grid>
      </Container>
    </Segment>
  );
};

export default Footer;
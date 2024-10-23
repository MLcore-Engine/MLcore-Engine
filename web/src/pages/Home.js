import React from 'react';
import { Header as SemanticHeader, Segment, Grid, Icon, Button } from 'semantic-ui-react';
import { Link } from 'react-router-dom';

const Home = () => (
  <div className="min-h-screen bg-gray-100 py-12 px-4 sm:px-6 lg:px-8">
    <div className="max-w-7xl mx-auto">
      <SemanticHeader as="h1" textAlign="center" className="text-4xl font-extrabold text-gray-900 mb-8">
        欢迎使用我们的平台
      </SemanticHeader>

      <Segment>
        <Grid columns={2} stackable textAlign="center">
          <Grid.Row verticalAlign="middle">
            <Grid.Column>
              <SemanticHeader as="h2" icon>
                <Icon name="docker" />
                镜像管理
                <SemanticHeader.Subheader>
                  管理和组织您的 Docker 镜像
                </SemanticHeader.Subheader>
              </SemanticHeader>
              <Button as={Link} to="/image-management" primary>
                查看镜像
              </Button>
            </Grid.Column>
            <Grid.Column>
              <SemanticHeader as="h2" icon>
                <Icon name="tasks" />
                任务调度
                <SemanticHeader.Subheader>
                  创建和监控您的计算任务
                </SemanticHeader.Subheader>
              </SemanticHeader>
              <Button as={Link} to="/task-scheduling" primary>
                管理任务
              </Button>
            </Grid.Column>
          </Grid.Row>
        </Grid>
      </Segment>

      <Segment>
        <SemanticHeader as="h3">平台特性</SemanticHeader>
        <Grid columns={3} divided>
          <Grid.Row>
            <Grid.Column>
              <Icon name="shield" size="large" />
              <p>安全可靠的镜像管理</p>
            </Grid.Column>
            <Grid.Column>
              <Icon name="sync" size="large" />
              <p>高效的任务调度系统</p>
            </Grid.Column>
            <Grid.Column>
              <Icon name="chart line" size="large" />
              <p>实时性能监控</p>
            </Grid.Column>
          </Grid.Row>
        </Grid>
      </Segment>

      <Segment textAlign="center">
        <SemanticHeader as="h3">开始使用</SemanticHeader>
        <p>立即体验我们强大的平台功能</p>
        <Button as={Link} to="/signup" size="large" primary>
          注册账号
        </Button>
      </Segment>
    </div>
  </div>
);

export default Home;
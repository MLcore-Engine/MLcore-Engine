import React from 'react';
import {
  Container,
  Header,
  Segment,
  Grid,
  List,
  Icon,
  Button,
} from 'semantic-ui-react';

const About = () => {
  const primaryColor = '#133E87';  
  const platformFeatures = [
    {
      icon: 'code',
      title: '模型开发 ',
      description:
        '支持在线开发模型，提供 Jupyter Notebook 环境，可以直接编写和调试代码。',
    },
    {
      icon: 'cogs',
      title: '模型训练',
      description:
        '基于 Kubernetes Training Operator，支持分布式训练，自动化训练流程管理。',
    },
    {
      icon: 'save',
      title: '模型存储',
      description:
        '使用 MinIO 对象存储，为每个用户提供独立的存储空间，支持数据集和模型文件管理。',
    },
    {
      icon: 'docker',
      title: '镜像管理',
      description:
        '集成 Docker Registry，支持模型镜像的自动构建、版本管理和快速部署。',
    },
    {
      icon: 'server',
      title: '模型部署',
      description:
        '使用 Triton Inference Server，提供高性能的模型推理服务，支持多种深度学习框架。',
    },
    {
      icon: 'chart line',
      title: '监控管理',
      description:
        '实时监控训练任务和部署服务的状态，提供详细的日志和性能指标。',
    },
  ];
  const cardStyle = {
    backgroundImage: 'linear-gradient(120deg, #fdfbfb 0%, #ebedee 100%)',
    border: 'none',
    boxShadow:
      '0 2px 4px 0 rgba(34, 36, 38, .12), 0 2px 10px 0 rgba(34, 36, 38, .15)',
  };
  return (
    <Container style={{ padding: '2em 0' }}>
      {/* 平台介绍 */}
      <Segment basic textAlign='center' style={{ marginBottom: '3em' }}>
        <Header as='h1' style={{ fontSize: '2.5em', marginBottom: '0.5em' }}>
          <div style={{ color: primaryColor }}>机器学习平台</div>
          <Header.Subheader style={{ marginTop: '1em' }}>
            一个简单、高效的机器学习模型开发与部署平台
          </Header.Subheader>
        </Header>
      </Segment>

      {/* 主要功能 */}
      <Segment basic>
        <Header as='h2' style={{ marginBottom: '1.5em' }}>
          <Icon name='cube' style={{ color: primaryColor }} />
          <Header.Content>
            <div style={{ color: primaryColor }}>平台功能</div>
            <Header.Subheader>全面的机器学习工作流支持</Header.Subheader>
          </Header.Content>
        </Header>

        <Grid stackable columns={3}>
          {platformFeatures.map((feature, index) => (
            <Grid.Column key={index}>
              <Segment raised padded style={cardStyle}>
                <Header as='h3' style={{ color: primaryColor }}>
                  <Icon name={feature.icon} />
                  <Header.Content>{feature.title}</Header.Content>
                </Header>
                <p>{feature.description}</p>
              </Segment>
            </Grid.Column>
          ))}
        </Grid>
      </Segment>

      {/* 技术栈 */}
      <Segment basic style={{ marginTop: '3em' }}>
        <Header as='h2'>
          <Icon name='settings' style={{ color: primaryColor }} />
          <Header.Content>
            <div style={{ color: primaryColor }}>技术栈</div>
            <Header.Subheader>使用现代化的技术栈构建</Header.Subheader>
          </Header.Content>
        </Header>

        <Grid columns={2} stackable>
          <Grid.Column>
            <Segment style={cardStyle}>
              <Header as='h3' style={{ color: primaryColor }}>前端技术</Header>
              <List relaxed>
                <List.Item>
                  <Icon name='react' style={{ color: primaryColor }} />
                  <List.Content>React + React Router</List.Content>
                </List.Item>
                <List.Item>
                  <Icon name='js' style={{ color: primaryColor }} />
                  <List.Content>JavaScript/TypeScript</List.Content>
                </List.Item>
                <List.Item>
                  <Icon name='paint brush' style={{ color: primaryColor }} />
                  <List.Content>Semantic UI React</List.Content>
                </List.Item>
              </List>
            </Segment>
          </Grid.Column>

          <Grid.Column>
            <Segment style={cardStyle}>
              <Header as='h3' style={{ color: primaryColor }}>
                后端技术
              </Header>
              <List relaxed>
                <List.Item>
                  <Icon name='docker' style={{ color: primaryColor }} />
                  <List.Content>Docker + Kubernetes</List.Content>
                </List.Item>
                <List.Item>
                  <Icon name='database' style={{ color: primaryColor }}/>
                  <List.Content>MinIO 对象存储</List.Content>
                </List.Item>
                <List.Item>
                  <Icon name='server' style={{ color: primaryColor }}/>
                  <List.Content>Triton Inference Server</List.Content>
                </List.Item>
              </List>
            </Segment>
          </Grid.Column>
        </Grid>
      </Segment>

      {/* 项目信息 */}
      <Segment basic style={{ marginTop: '3em' }}>
        <Header as='h2' >
          <Icon name='github' />
          <Header.Content>
          <div style={{ color: primaryColor }}>项目信息</div>
            <Header.Subheader>开源项目，欢迎贡献</Header.Subheader>
          </Header.Content>
        </Header>

        <Segment>
          <Grid columns={2} stackable verticalAlign='middle'>
            <Grid.Column>
              <List relaxed>
                <List.Item>
                  <Icon name='github' style={{ color: primaryColor }} />
                  <List.Content>
                    <List.Header>项目地址</List.Header>
                    <List.Description>
                      <a
                        href='https://github.com/yourusername/project'
                        target='_blank'
                        rel='noopener noreferrer'
                      >
                        github.com/yourusername/project
                      </a>
                    </List.Description>
                  </List.Content>
                </List.Item>
                <List.Item>
                  <Icon name='book' style={{ color: primaryColor }}/>
                  <List.Content>
                    <List.Header>文档地址</List.Header>
                    <List.Description>
                      <a
                        href='https://docs.yourproject.com'
                        target='_blank'
                        rel='noopener noreferrer'
                      >
                        docs.yourproject.com
                      </a>
                    </List.Description>
                  </List.Content>
                </List.Item>
              </List>
            </Grid.Column>
            <Grid.Column textAlign='center'>
              <Button
                as='a'
                href='https://github.com/yourusername/project'
                target='_blank'
                rel='noopener noreferrer'
                size='large'
                color='black'
              >
                <Icon name='github' /> 访问项目
              </Button>
            </Grid.Column>
          </Grid>
        </Segment>
      </Segment>

      {/* 联系方式 */}
      <Segment basic  style={{ marginTop: '3em' }}>
        <Header as='h2'>
          <Icon name='mail' />
          <Header.Content>
          <div style={{ color: primaryColor }}>联系我们</div>
            <Header.Subheader>有任何问题或建议，欢迎联系</Header.Subheader>
          </Header.Content>
        </Header>
        <Button.Group size='large' style={{ marginTop: '1em' }}>
          <Button as='a' href='mailto:contact@yourproject.com'>
            <Icon name='mail' /> Email
          </Button>
          <Button
            as='a'
            href='https://github.com/yourusername/project/issues'
            target='_blank'
          >
            <Icon name='github' /> Issues
          </Button>
        </Button.Group>
      </Segment>
    </Container>
  );
};

export default About;

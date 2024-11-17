import React from 'react';
import {
  Header as SemanticHeader,
  Segment,
  Grid,
  Icon,
  Divider,
  Card,
} from 'semantic-ui-react';
import '../styles/home.css';


const cardStyle = {
  border: 'none',
  boxShadow: '0 4px 8px rgba(0, 0, 0, 0.1)', // 添加阴影
  borderRadius: '8px', // 卡片圆角
  minHeight: '120px', // 设置统一的最小高度
  display: 'flex', // 启用 flex 布局
  flexDirection: 'column', // 内容垂直排列
  justifyContent: 'center', // 垂直居中
};
const ProcessStep = ({ icon, title, description, color = '#133e87' }) => (
  <Card fluid style={cardStyle}>
    <Card.Content>
      <Card.Header style={{ color: color }}>
        <Icon name={icon} style={{ marginRight: '3px' }} />
        {title}
      </Card.Header>
      <Card.Description>{description}</Card.Description>
    </Card.Content>
  </Card>
);

const ArrowIcon = () => (
  <div style={{ marginLeft: '20px' }}>
    <Icon name='arrow down' size='large' style={{ color: '#B7B7B7' }} />
  </div>
);
const primaryColor = '#133E87';  
const Home = () => {
  return (
    <div className='min-h-screen  py-8 px-4 sm:px-6 lg:px-8 bg'>
      <div className='max-w-7xl mx-auto'>
        <SemanticHeader
          as='h1'
          textAlign='center'
          className='text-3xl font-bold text-gray-900 mb-8'
        >
          <div style={{ color: primaryColor }}>MLcore-Engine workflow</div>
        </SemanticHeader>

        <Grid container stackable>
          {/* 用户操作区域 */}
          <Grid.Row>
            <Grid.Column width={16}>
              <Segment>
                <SemanticHeader as='h3'>User Operations</SemanticHeader>
                <Grid columns={5} divided>
                  <Grid.Row>
                    <Grid.Column>
                      <ProcessStep
                        icon='folder open'
                        title='Project Creation/Management'
                        description='Create and manage project spaces'
                      />
                    </Grid.Column>
                    <Grid.Column>
                      <ProcessStep
                        icon='upload'
                        title='Data Upload'
                        description='Upload training datasets'
                      />
                    </Grid.Column>
                    <Grid.Column>
                      <ProcessStep
                        icon='code'
                        title='Model Development'
                        description='Develop and debug models'
                      />
                    </Grid.Column>
                    <Grid.Column>
                      <ProcessStep
                        icon='play'
                        title='Training Task Submission'
                        description='Submit and manage training tasks'
                      />
                    </Grid.Column>
                    <Grid.Column>
                      <ProcessStep
                        icon='server'
                        title='Model Deployment Request'
                        description='Deploy models as services'
                      />
                    </Grid.Column>
                  </Grid.Row>
                </Grid>
              </Segment>
            </Grid.Column>
          </Grid.Row>

          <ArrowIcon />

          {/* 数据存储层 */}
          <Grid.Row>
            <Grid.Column width={16}>
              <Segment>
                <SemanticHeader as='h3'>Data Storage Layer</SemanticHeader>
                <Grid columns={2} divided>
                  <Grid.Row>
                    <Grid.Column>
                      <ProcessStep
                        icon='database'
                        title='MinIO Object Storage'
                        description='Store training data and model files'
                        color='#3B1E54'
                      />
                    </Grid.Column>
                    <Grid.Column>
                      <ProcessStep
                        icon='docker'
                        title='Docker Registry'
                        description='Store container images'
                        color='#3B1E54'
                      />
                    </Grid.Column>
                  </Grid.Row>
                </Grid>
              </Segment>
            </Grid.Column>
          </Grid.Row>

          <ArrowIcon />

          {/* Kubernetes Cluster */}
          <Grid.Row>
            <Grid.Column width={16}>
              <Segment>
                <SemanticHeader as='h3'>Kubernetes Cluster</SemanticHeader>
                <Grid columns={4} divided>
                  <Grid.Row>
                    <Grid.Column>
                      <ProcessStep
                        icon='settings'
                        title='Training Operator'
                        description='Manage training tasks'
                        color='#31511E'
                      />
                    </Grid.Column>
                    <Grid.Column>
                      <ProcessStep
                        icon='microchip'
                        title='Training Pod'
                        description='Execute model training'
                        color='#31511E'
                      />
                    </Grid.Column>
                    <Grid.Column>
                      <ProcessStep
                        icon='server'
                        title='Triton Inference Server'
                        description='Provide inference services'
                        color='#31511E'
                      />
                    </Grid.Column>
                    <Grid.Column>
                      <ProcessStep
                        icon='cogs'
                        title='Model Serving Pod'
                        description='Run inference services'
                        color='#31511E'
                      />
                    </Grid.Column>
                  </Grid.Row>
                </Grid>
              </Segment>
            </Grid.Column>
          </Grid.Row>

          <ArrowIcon />

          {/* 底部监控和模型管理 */}
          <Grid.Row>
            <Grid.Column width={16}>
              <Grid columns={2} >
                <Grid.Row>
                  <Grid.Column>
                    <Segment>
                      <SemanticHeader as='h3'>Model Management</SemanticHeader>
                      <Grid columns={3} divided>
                        <Grid.Row>
                          <Grid.Column>
                            <ProcessStep
                              icon='chart bar'
                              title='MLflow'
                              description='Model lifecycle management'
                              color='#1A1A19'
                            />
                          </Grid.Column>
                          <Grid.Column>
                            <ProcessStep
                              icon='history'
                              title='Model Version Control'
                              description='Version management and tracking'
                              color='#1A1A19'
                            />
                          </Grid.Column>
                          <Grid.Column>
                            <ProcessStep
                              icon='line graph'
                              title='Model Metrics Tracking'
                              description='Performance metrics monitoring'
                              color='#1A1A19'
                            />
                          </Grid.Column>
                        </Grid.Row>
                      </Grid>
                    </Segment>
                  </Grid.Column>
                  <Grid.Column>
                    <Segment>
                      <SemanticHeader as='h3'>Monitoring System</SemanticHeader>
                      <Grid columns={3} divided>
                        <Grid.Row>
                          <Grid.Column>
                            <ProcessStep
                              icon='eye'
                              title='Prometheus'
                              description='Metrics collection'
                              color='#2A3663'
                            />
                          </Grid.Column>
                          <Grid.Column>
                            <ProcessStep
                              icon='dashboard'
                              title='Grafana'
                              description='Visualization monitoring'
                              color='#2A3663'
                            />
                          </Grid.Column>
                          <Grid.Column>
                            <ProcessStep
                              icon='file alternate'
                              title='EFK Stack'
                              description='Log management'
                              color='#2A3663'
                            />
                          </Grid.Column>
                        </Grid.Row>
                      </Grid>
                    </Segment>
                  </Grid.Column>
                </Grid.Row>
              </Grid>
            </Grid.Column>
          </Grid.Row>
        </Grid>
      </div>
    </div>
  );
};

export default Home;

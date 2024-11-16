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
          <div style={{ color: primaryColor }}>机器学习平台工作流程</div>
        </SemanticHeader>

        <Grid container stackable>
          {/* 用户操作区域 */}
          <Grid.Row>
            <Grid.Column width={16}>
              <Segment>
                <SemanticHeader as='h3'>用户操作</SemanticHeader>
                <Grid columns={5} divided>
                  <Grid.Row>
                    <Grid.Column>
                      <ProcessStep
                        icon='folder open'
                        title='项目创建/管理'
                        description='创建和管理项目空间'
                      />
                    </Grid.Column>
                    <Grid.Column>
                      <ProcessStep
                        icon='upload'
                        title='数据上传'
                        description='上传训练数据集'
                      />
                    </Grid.Column>
                    <Grid.Column>
                      <ProcessStep
                        icon='code'
                        title='模型开发'
                        description='开发和调试模型'
                      />
                    </Grid.Column>
                    <Grid.Column>
                      <ProcessStep
                        icon='play'
                        title='训练任务提交'
                        description='提交和管理训练任务'
                      />
                    </Grid.Column>
                    <Grid.Column>
                      <ProcessStep
                        icon='server'
                        title='模型部署请求'
                        description='部署模型为服务'
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
                <SemanticHeader as='h3'>数据存储层</SemanticHeader>
                <Grid columns={2} divided>
                  <Grid.Row>
                    <Grid.Column>
                      <ProcessStep
                        icon='database'
                        title='MinIO对象存储'
                        description='存储训练数据和模型文件'
                        color='#3B1E54'
                      />
                    </Grid.Column>
                    <Grid.Column>
                      <ProcessStep
                        icon='docker'
                        title='Docker Registry'
                        description='存储容器镜像'
                        color='#3B1E54'
                      />
                    </Grid.Column>
                  </Grid.Row>
                </Grid>
              </Segment>
            </Grid.Column>
          </Grid.Row>

          <ArrowIcon />

          {/* Kubernetes集群 */}
          <Grid.Row>
            <Grid.Column width={16}>
              <Segment>
                <SemanticHeader as='h3'>Kubernetes集群</SemanticHeader>
                <Grid columns={4} divided>
                  <Grid.Row>
                    <Grid.Column>
                      <ProcessStep
                        icon='settings'
                        title='Training Operator'
                        description='管理训练任务'
                        color='#31511E'
                      />
                    </Grid.Column>
                    <Grid.Column>
                      <ProcessStep
                        icon='microchip'
                        title='训练Pod'
                        description='执行模型训练'
                        color='#31511E'
                      />
                    </Grid.Column>
                    <Grid.Column>
                      <ProcessStep
                        icon='server'
                        title='Triton Inference Server'
                        description='提供推理服务'
                        color='#31511E'
                      />
                    </Grid.Column>
                    <Grid.Column>
                      <ProcessStep
                        icon='cogs'
                        title='模型服务Pod'
                        description='运行推理服务'
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
                      <SemanticHeader as='h3'>模型管理</SemanticHeader>
                      <Grid columns={3} divided>
                        <Grid.Row>
                          <Grid.Column>
                            <ProcessStep
                              icon='chart bar'
                              title='MLflow'
                              description='模型生命周期管理'
                              color='#1A1A19'
                            />
                          </Grid.Column>
                          <Grid.Column>
                            <ProcessStep
                              icon='history'
                              title='模型版本控制'
                              description='版本管理与追踪'
                              color='#1A1A19'
                            />
                          </Grid.Column>
                          <Grid.Column>
                            <ProcessStep
                              icon='line graph'
                              title='模型指标追踪'
                              description='性能指标监控'
                              color='#1A1A19'
                            />
                          </Grid.Column>
                        </Grid.Row>
                      </Grid>
                    </Segment>
                  </Grid.Column>
                  <Grid.Column>
                    <Segment>
                      <SemanticHeader as='h3'>监控系统</SemanticHeader>
                      <Grid columns={3} divided>
                        <Grid.Row>
                          <Grid.Column>
                            <ProcessStep
                              icon='eye'
                              title='Prometheus'
                              description='指标收集'
                              color='#2A3663'
                            />
                          </Grid.Column>
                          <Grid.Column>
                            <ProcessStep
                              icon='dashboard'
                              title='Grafana'
                              description='可视化监控'
                              color='#2A3663'
                            />
                          </Grid.Column>
                          <Grid.Column>
                            <ProcessStep
                              icon='file alternate'
                              title='EFK Stack'
                              description='日志管理'
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

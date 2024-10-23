import React from 'react';
import { Header as SemanticHeader, Segment, List, Button, Icon, Breadcrumb } from 'semantic-ui-react';
import { Link } from 'react-router-dom';

const ModelTraining = () => (
  <div>
    <SemanticHeader as="h2">模型训练</SemanticHeader>
    <Segment>
      <Button as={Link} to="/train/create" primary>
        <Icon name="plus" />
        创建训练任务
      </Button>
      <Breadcrumb>
        <Breadcrumb.Section link>主页</Breadcrumb.Section>
        <Breadcrumb.Divider />
        <Breadcrumb.Section active>模型训练</Breadcrumb.Section>
      </Breadcrumb>
      <List divided relaxed>
        {/* 示例训练任务 */}
        <List.Item>
          <Icon name="truck" size="large" verticalAlign="middle" />
          <List.Content>
            <List.Header>训练任务 1</List.Header>
            <List.Description>正在运行 - 75% 完成</List.Description>
          </List.Content>
          <List.Content floated="right">
            <Button icon>
              <Icon name="pause" />
            </Button>
            <Button icon color="red">
              <Icon name="trash" />
            </Button>
          </List.Content>
        </List.Item>
        {/* 在这里动态渲染训练任务列表 */}
      </List>
    </Segment>
  </div>
);

export default ModelTraining;


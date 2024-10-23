import React from 'react';
import { Container, Header } from 'semantic-ui-react';

const Dashboard = () => {
  return (
    <Container>
      <Header as="h1">仪表板</Header>
      <p>欢迎来到您的仪表板。这里可以显示概览信息、最近的活动等。</p>
      {/* 这里可以添加更多的仪表板内容，如统计数据、图表等 */}
    </Container>
  );
};

export default Dashboard;
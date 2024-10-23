import React from 'react';
import { Header as SemanticHeader, Segment, Grid, Statistic, Icon, Button } from 'semantic-ui-react';
import { Link } from 'react-router-dom';

const ModelDeployment = () => (
  <div>
    <SemanticHeader as="h2">模型部署</SemanticHeader>
    <Segment>
      <Grid>
        <Grid.Row columns={2}>
          <Grid.Column>
            <Statistic>
              <Statistic.Value>42</Statistic.Value>
              <Statistic.Label>已部署模型</Statistic.Label>
            </Statistic>
          </Grid.Column>
          <Grid.Column>
            <Button as={Link} to="/deploy/create" primary>
              <Icon name="plus" />
              创建部署
            </Button>
          </Grid.Column>
        </Grid.Row>
      </Grid>
      {/* 部署模型列表或仪表盘 */}
      <Segment placeholder>
        <Icon name="clipboard outline" size="huge" />
        <p>暂无已部署的模型。</p>
        <Button as={Link} to="/deploy/create" primary>
          <Icon name="plus" />
          创建部署
        </Button>
      </Segment>
    </Segment>
  </div>
);

export default ModelDeployment;
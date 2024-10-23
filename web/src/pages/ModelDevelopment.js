import React from 'react';
import { Header as SemanticHeader, Segment, Grid, Card, Icon } from 'semantic-ui-react';
import { Link } from 'react-router-dom';

const ModelDevelopment = () => (
  <div>
    <SemanticHeader as="h2">模型开发</SemanticHeader>
    <Segment>
      <Grid stackable>
        <Grid.Row columns={3}>
          <Grid.Column>
            <Card as={Link} to="/notebook/create">
              <Card.Content>
                <Icon name="plus circle" size="huge" color="blue" />
                <Card.Header>创建新模型</Card.Header>
                <Card.Description>开始开发新的机器学习模型。</Card.Description>
              </Card.Content>
            </Card>
          </Grid.Column>
          <Grid.Column>
            <Card as={Link} to="/notebook">
              <Card.Content>
                <Icon name="list alternate" size="huge" color="green" />
                <Card.Header>模型列表</Card.Header>
                <Card.Description>查看和管理所有已开发的模型。</Card.Description>
              </Card.Content>
            </Card>
          </Grid.Column>
          <Grid.Column>
            <Card as={Link} to="/notebook/import">
              <Card.Content>
                <Icon name="upload" size="huge" color="orange" />
                <Card.Header>导入模型</Card.Header>
                <Card.Description>从外部源导入预训练的模型。</Card.Description>
              </Card.Content>
            </Card>
          </Grid.Column>
        </Grid.Row>
      </Grid>
    </Segment>
  </div>
);

export default ModelDevelopment;


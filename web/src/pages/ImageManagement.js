import React from 'react';
import { Header as SemanticHeader, Segment, Table, Button, Icon } from 'semantic-ui-react';
import { Link } from 'react-router-dom';

const ImageManagement = () => (
  <div>
    <SemanticHeader as="h2">镜像管理</SemanticHeader>
    <Segment>
      <Button as={Link} to="/image/create" primary>
        <Icon name="plus" />
        添加镜像
      </Button>
      <Table celled striped selectable>
        <Table.Header>
          <Table.Row>
            <Table.HeaderCell>ID</Table.HeaderCell>
            <Table.HeaderCell>镜像名称</Table.HeaderCell>
            <Table.HeaderCell>标签</Table.HeaderCell>
            <Table.HeaderCell>操作</Table.HeaderCell>
          </Table.Row>
        </Table.Header>

        <Table.Body>
          {/* 示例行 */}
          <Table.Row>
            <Table.Cell>1</Table.Cell>
            <Table.Cell>tensorflow/tensorflow:latest</Table.Cell>
            <Table.Cell>深度学习</Table.Cell>
            <Table.Cell>
              <Button icon>
                <Icon name="edit" />
              </Button>
              <Button icon color="red">
                <Icon name="trash" />
              </Button>
            </Table.Cell>
          </Table.Row>
          {/* 在这里动态渲染镜像列表 */}
        </Table.Body>
      </Table>
    </Segment>
  </div>
);

export default ImageManagement;


import React from 'react';
import { Modal, Button, Icon, List, Label } from 'semantic-ui-react';

const ServingDetails = ({ deploy, isOpen, onClose }) => {
  if (!deploy) return null;

  // Parse ports from JSON string
  const parsedPorts = deploy.ports ? JSON.parse(deploy.ports) : [];
  
  // Parse access URLs
  const accessUrls = deploy.access_url ? deploy.access_url.split(',') : [];

  const getStatusColor = (status) => {
    switch (status.toLowerCase()) {
      case 'running':
        return 'green';
      case 'pending':
        return 'yellow';
      case 'failed':
        return 'red';
      case 'creating':
        return 'blue';
      case 'deleted':
        return 'grey';
      default:
        return 'grey';
    }
  };

  return (
    <Modal open={isOpen} onClose={onClose}>
      <Modal.Header>Triton Deployment Details</Modal.Header>
      <Modal.Content>
        <List divided relaxed>
          <List.Item>
            <List.Header>Name</List.Header>
            {deploy.name}
          </List.Item>
          
          <List.Item>
            <List.Header>Status</List.Header>
            <Label color={getStatusColor(deploy.status)}>
              {deploy.status}
            </Label>
          </List.Item>

          <List.Item>
            <List.Header>Namespace</List.Header>
            {deploy.namespace}
          </List.Item>

          <List.Item>
            <List.Header>Image</List.Header>
            {deploy.image}
          </List.Item>

          <List.Item>
            <List.Header>Replicas</List.Header>
            {deploy.replicas}
          </List.Item>

          <List.Item>
            <List.Header>Resources</List.Header>
            <List.Description>
              <div>CPU: {deploy.cpu}</div>
              <div>Memory: {deploy.memory}</div>
              <div>GPU: {deploy.gpu}</div>
            </List.Description>
          </List.Item>

          <List.Item>
            <List.Header>Ports</List.Header>
            <List bulleted>
              {parsedPorts.map((port, index) => (
                <List.Item key={index}>
                  {port.name ? `${port.name}: ` : ''}
                  Port: {port.port} → NodePort: {port.nodePort}
                </List.Item>
              ))}
            </List>
          </List.Item>

          <List.Item>
            <List.Header>Access URLs</List.Header>
            <List bulleted>
              {accessUrls.map((url, index) => (
                <List.Item key={index}>
                  <a href={url} target="_blank" rel="noopener noreferrer">
                    {url}
                  </a>
                </List.Item>
              ))}
            </List>
          </List.Item>

          <List.Item>
            <List.Header>Labels</List.Header>
            <List.Description>
              {deploy.labels && Object.entries(JSON.parse(deploy.labels)).map(([key, value]) => (
                <Label key={key} style={{ margin: '2px' }}>
                  {key}: {value}
                </Label>
              ))}
            </List.Description>
          </List.Item>

          <List.Item>
            <List.Header>Created At</List.Header>
            {new Date(deploy.created_at).toLocaleString()}
          </List.Item>

          <List.Item>
            <List.Header>Updated At</List.Header>
            {new Date(deploy.updated_at).toLocaleString()}
          </List.Item>
        </List>
      </Modal.Content>
      <Modal.Actions>
        <Button onClick={onClose} primary>
          <Icon name="checkmark" /> Close
        </Button>
      </Modal.Actions>
    </Modal>
  );
};

export default ServingDetails;

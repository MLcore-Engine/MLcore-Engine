// src/pages/training/TrainingJobDetails.js

import React from 'react';
import { Modal, Button, Icon, List } from 'semantic-ui-react';

const TrainingJobDetails = ({ job, isOpen, onClose }) => {
  if (!job) return null;

  const parsedArgs = job.args ? JSON.parse(job.args) : [];

  return (
    <Modal open={isOpen} onClose={onClose}>
      <Modal.Header>Training Job Details</Modal.Header>
      <Modal.Content>
        <List>
          <List.Item>
            <List.Header>Name</List.Header>
            {job.name}
          </List.Item>
          <List.Item>
            <List.Header>Status</List.Header>
            {job.status}
          </List.Item>
          <List.Item>
            <List.Header>Image</List.Header>
            {job.image}
          </List.Item>
          <List.Item>
            <List.Header>Image Pull Policy</List.Header>
            {job.image_pull_policy}
          </List.Item>
          <List.Item>
            <List.Header>Restart Policy</List.Header>
            {job.restart_policy}
          </List.Item>
          <List.Item>
            <List.Header>Namespace</List.Header>
            {job.namespace}
          </List.Item>
          <List.Item>
            <List.Header>CPU Limit</List.Header>
            {job.cpu_limit}
          </List.Item>
          <List.Item>
            <List.Header>Memory Limit</List.Header>
            {job.memory_limit}
          </List.Item>
          <List.Item>
            <List.Header>Master Replicas</List.Header>
            {job.master_replicas}
          </List.Item>
          <List.Item>
            <List.Header>Worker Replicas</List.Header>
            {job.worker_replicas}
          </List.Item>
          <List.Item>
            <List.Header>GPUs Per Node</List.Header>
            {job.gpus_per_node}
          </List.Item>
          <List.Item>
            <List.Header>Args</List.Header>
            <List bulleted>
              {parsedArgs.map((arg, index) => (
                <List.Item key={index}>{arg}</List.Item>
              ))}
            </List>
          </List.Item>
          <List.Item>
            <List.Header>Description</List.Header>
            {job.describe}
          </List.Item>
          <List.Item>
            <List.Header>Created At</List.Header>
            {new Date(job.created_at).toLocaleString()}
          </List.Item>
          <List.Item>
            <List.Header>Updated At</List.Header>
            {new Date(job.updated_at).toLocaleString()}
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

export default TrainingJobDetails;
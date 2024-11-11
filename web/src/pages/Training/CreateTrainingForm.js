// src/pages/training/CreateTrainingForm.js

import React, { useState, useEffect } from 'react';
import { Modal, Form, Button, Message, Dropdown, Icon, List } from 'semantic-ui-react';
import { trainingAPI } from '../../api/trainingAPI';
import { toast } from 'react-toastify';

const CreateTrainingForm = ({ isOpen, onClose, onRefresh }) => {
  const userInfo = JSON.parse(localStorage.getItem('user') || '{}');
  const [formData, setFormData] = useState({
    user_id: userInfo.id || '',
    project_id: userInfo.projects?.[0]?.ID || '',
    describe: '',
    namespace: '',
    image: '',
    image_pull_policy: '',
    restart_policy: '',
    args: [],
    master_replicas: 1,
    worker_replicas: 1,
    gpus_per_node: 0,
    cpu_limit: '2',
    memory_limit: '8Gi',
    status: 'Creating',
    expand: '{}',
  });

  const [loading, setLoading] = useState(false);
  const [errors, setErrors] = useState({});
  const [dropdownOptions, setDropdownOptions] = useState({
    images: [],
    imagePullPolicies: [],
    restartPolicies: [],
    namespaces: [],
  });
  const [argInput, setArgInput] = useState('');

  useEffect(() => {
    // Initialize dropdown options
    setDropdownOptions({
      images: [
        { key: 'image1', text: 'kubeflow/pytorch-dist-mnist:latest', value: 'kubeflow/pytorch-dist-mnist:latest' },
        { key: 'image2', text: 'custom/image:tag', value: 'custom/image:tag' },
        // Add more image options as needed
      ],
      imagePullPolicies: [
        { key: 'IfNotPresent', text: 'IfNotPresent', value: 'IfNotPresent' },
        { key: 'Always', text: 'Always', value: 'Always' },
        { key: 'Never', text: 'Never', value: 'Never' },
      ],
      restartPolicies: [
        { key: 'OnFailure', text: 'OnFailure', value: 'OnFailure' },
        { key: 'Never', text: 'Never', value: 'Never' },
      ],
      namespaces: [
        { key: 'train', text: 'train', value: 'train' },
        { key: 'development', text: 'development', value: 'development' },
        { key: 'production', text: 'production', value: 'production' },
        // Add more namespace options as needed
      ],
    });
  }, []);

  const handleChange = (e, { name, value }) => {
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleArgChange = (e) => {
    setArgInput(e.target.value);
  };

  const addArg = () => {
    if (argInput.trim() === '') return;
    setFormData((prev) => ({ ...prev, args: [...prev.args, argInput.trim()] }));
    setArgInput('');
  };

  const removeArg = (index) => {
    setFormData((prev) => {
      const newArgs = [...prev.args];
      newArgs.splice(index, 1);
      return { ...prev, args: newArgs };
    });
  };

  const validateForm = () => {
    const newErrors = {};
    if (!formData.image) newErrors.image = 'Image is required';
    if (!formData.image_pull_policy) newErrors.image_pull_policy = 'Image Pull Policy is required';
    if (!formData.restart_policy) newErrors.restart_policy = 'Restart Policy is required';
    if (!formData.namespace) newErrors.namespace = 'Namespace is required';
    if (!formData.cpu_limit) newErrors.cpu_limit = 'CPU limit is required';
    if (!formData.memory_limit) newErrors.memory_limit = 'Memory limit is required';
    if (formData.master_replicas < 1) newErrors.master_replicas = 'At least one master replica is required';
    if (formData.worker_replicas < 1) newErrors.worker_replicas = 'At least one worker replica is required';
    // Add more validation as needed
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async () => {
    if (!validateForm()) return;
    setLoading(true);
    try {
      const payload = {
        ...formData,
        args: JSON.stringify(formData.args),
      };
      await trainingAPI.createTrainingJob(payload);
      toast.info('Training job created successfully');
      onClose();
      onRefresh();
    } catch (error) {
      toast.error(error.message || 'Failed to create training job');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Modal open={isOpen} onClose={onClose}>
      <Modal.Header>CreateTrainingJob</Modal.Header>
      <Modal.Content>
        <Form loading={loading} error={Object.keys(errors).length > 0}>
          <Form.Field
            control={Dropdown}
            label="Image"
            placeholder="Select Image"
            options={dropdownOptions.images}
            name="image"
            onChange={handleChange}
            value={formData.image}
            error={errors.image ? { content: errors.image } : null}
            selection
          />
          <Form.Field
            control={Dropdown}
            label="Image Pull Policy"
            placeholder="Select Image Pull Policy"
            options={dropdownOptions.imagePullPolicies}
            name="image_pull_policy"
            onChange={handleChange}
            value={formData.image_pull_policy}
            error={errors.image_pull_policy ? { content: errors.image_pull_policy } : null}
            selection
          />
          <Form.Field
            control={Dropdown}
            label="Restart Policy"
            placeholder="Select Restart Policy"
            options={dropdownOptions.restartPolicies}
            name="restart_policy"
            onChange={handleChange}
            value={formData.restart_policy}
            error={errors.restart_policy ? { content: errors.restart_policy } : null}
            selection
          />
          <Form.Field
            control={Dropdown}
            label="Namespace"
            placeholder="Select Namespace"
            options={dropdownOptions.namespaces}
            name="namespace"
            onChange={handleChange}
            value={formData.namespace}
            error={errors.namespace ? { content: errors.namespace } : null}
            selection
          />
          <Form.Field>
            <label>Args</label>
            <div className="d-flex">
              <Form.Input
                placeholder="Add argument"
                value={argInput}
                onChange={handleArgChange}
                onKeyPress={(e) => {
                  if (e.key === 'Enter') {
                    e.preventDefault();
                    addArg();
                  }
                }}
              />
              <Button type="button" onClick={addArg} icon>
                <Icon name="add" />
              </Button>
            </div>
            <List divided relaxed>
              {formData.args.map((arg, index) => (
                <List.Item key={index}>
                  <List.Content floated="right">
                    <Button type="button" icon color="red" onClick={() => removeArg(index)}>
                      <Icon name="trash" />
                    </Button>
                  </List.Content>
                  <List.Content>{arg}</List.Content>
                </List.Item>
              ))}
            </List>
          </Form.Field>
          <Form.Input
            label="CPU Limit"
            name="cpu_limit"
            type="text"
            value={formData.cpu_limit}
            onChange={handleChange}
            error={errors.cpu_limit ? { content: errors.cpu_limit } : null}
          />
          <Form.Input
            label="Memory Limit"
            name="memory_limit"
            type="text"
            value={formData.memory_limit}
            onChange={handleChange}
            error={errors.memory_limit ? { content: errors.memory_limit } : null}
          />
          <Form.Input
            label="Master Replicas"
            name="master_replicas"
            type="number"
            min="1"
            value={formData.master_replicas}
            onChange={handleChange}
            error={errors.master_replicas ? { content: errors.master_replicas } : null}
          />
          <Form.Input
            label="Worker Replicas"
            name="worker_replicas"
            type="number"
            min="1"
            value={formData.worker_replicas}
            onChange={handleChange}
            error={errors.worker_replicas ? { content: errors.worker_replicas } : null}
          />
          <Form.Input
            label="GPUs Per Node"
            name="gpus_per_node"
            type="number"
            min="0"
            value={formData.gpus_per_node}
            onChange={handleChange}
          />
          <Form.TextArea
            label="Description"
            name="describe"
            onChange={handleChange}
            value={formData.describe}
          />
          {Object.keys(errors).length > 0 && (
            <Message
              error
              header="There were some errors with your submission"
              list={Object.values(errors)}
            />
          )}
        </Form>
      </Modal.Content>
      <Modal.Actions>
        <Button onClick={onClose} disabled={loading}>
          Cancel
        </Button>
        <Button primary onClick={handleSubmit} loading={loading} disabled={loading}>
          Create
        </Button>
      </Modal.Actions>
    </Modal>
  );
};

export default CreateTrainingForm;
import React, { useState } from 'react';
import { Modal, Form, Button, Message, Dropdown } from 'semantic-ui-react';
import { tritonAPI } from '../../api/tritonAPI';
import { toast } from 'react-toastify';

const CreateTritonForm = ({ isOpen, onClose, onRefresh }) => {
  const [formData, setFormData] = useState({
    name: '',
    namespace: 'triton-serving',
    image: '192.168.12.121:5005/traning/triton:24.10-py3',
    replicas: 1,
    cpu: 2,
    memory: 4,
    gpu: 0
  });

  const [loading, setLoading] = useState(false);
  const [errors, setErrors] = useState({});

  const handleChange = (e, { name, value }) => {
    setFormData(prev => ({ ...prev, [name]: value }));
  };

  const validateForm = () => {
    const newErrors = {};
    if (!formData.image) newErrors.image = 'Image is required';
    if (!formData.cpu) newErrors.cpu = 'CPU is required';
    if (!formData.memory) newErrors.memory = 'Memory is required';
    if (formData.replicas < 1) newErrors.replicas = 'At least one replica is required';
    
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async () => {
    if (!validateForm()) return;
    
    setLoading(true);
    try {
      await tritonAPI.createTritonDeploy(formData);
      toast.success('Triton deployment created successfully');
      onClose();
      onRefresh();
    } catch (error) {
      toast.error(error.message || 'Failed to create Triton deployment');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Modal open={isOpen} onClose={onClose}>
      <Modal.Header>Create Triton Deployment</Modal.Header>
      <Modal.Content>
        <Form loading={loading} error={Object.keys(errors).length > 0}>
          <Form.Input
            label="Image"
            name="image"
            value={formData.image}
            onChange={handleChange}
            error={errors.image}
          />
          
          <Form.Input
            label="Namespace"
            name="namespace"
            value={formData.namespace}
            onChange={handleChange}
          />

          <Form.Input
            label="Replicas"
            name="replicas"
            type="number"
            min="1"
            value={formData.replicas}
            onChange={handleChange}
            error={errors.replicas}
          />

          <Form.Input
            label="CPU"
            name="cpu"
            value={formData.cpu}
            onChange={handleChange}
            error={errors.cpu}
          />

          <Form.Input
            label="Memory"
            name="memory"
            value={formData.memory}
            onChange={handleChange}
            error={errors.memory}
          />

          <Form.Input
            label="GPU"
            name="gpu"
            type="number"
            min="0"
            value={formData.gpu}
            onChange={handleChange}
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

export default CreateTritonForm;
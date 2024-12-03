import React, { useState, useEffect } from 'react';
import { Modal, Form, Button, Message, Tab } from 'semantic-ui-react';
import { tritonAPI } from '../../api/tritonAPI';
import { toast } from 'react-toastify';

const CreateTritonForm = ({ isOpen, onClose, onRefresh }) => {
  const [formData, setFormData] = useState({
    namespace: '',
    image: '',
    replicas: 1,
    cpu: 'small',
    memory: 'small',
    gpu: 'none',
    modelRepository: '',
    httpPort: '',
    grpcPort: '',
    metricsPort: '',
    backend: '',
    logVerbose: 0,
    logFormat: 'default'
  });

  const [loading, setLoading] = useState(false);
  const [errors, setErrors] = useState({});
  const [configOptions, setConfigOptions] = useState(null);

  useEffect(() => {
    const fetchConfig = async () => {
      try {
        const response = await tritonAPI.getTritonConfig();
        console.log('Raw config response:', response);
        
        if (response.success) {
          const config = response.data;
          console.log('Config data:', config);
          
          setConfigOptions(config);
          
          const defaultValues = {
            namespace: config.namespace,
            image: config.images[0],
            modelRepository: config.model_repository[0],
            httpPort: config.ports.http[0],
            grpcPort: config.ports.grpc[0],
            metricsPort: config.ports.metrics[0],
            backend: config.backend[0],
            logFormat: config.logging.formats[0],
            cpu: Object.keys(config.resources.cpu)[0],
            memory: Object.keys(config.resources.memory)[0],
            gpu: Object.keys(config.resources.gpu)[0]
          };
          
          console.log('Setting default values:', defaultValues);
          setFormData(prev => ({
            ...prev,
            ...defaultValues
          }));
        }
      } catch (error) {
        console.error('Config fetch error:', error);
        toast.error('Failed to load configuration options');
      }
    };
    fetchConfig();
  }, []);

  const getResourceOptions = (resources) => {
    if (!resources || typeof resources !== 'object') return [];
    
    return Object.entries(resources).map(([key, value]) => ({
      key,
      text: `${key} (${value})`,
      value: key
    }));
  };

  const getArrayOptions = (items) => {
    if (!Array.isArray(items)) return [];
    
    return items.map(item => ({
      key: String(item),
      text: String(item),
      value: String(item)
    }));
  };

  const handleChange = (e, { name, value }) => {
    console.log('Form field changed:', { name, value });
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

  console.log('Current formData:', formData);
  console.log('Current configOptions:', configOptions);

  if (!configOptions) {
    return (
      <Modal open={isOpen} onClose={onClose}>
        <Modal.Header>Create Triton Service</Modal.Header>
        <Modal.Content>
          <div>Loading configuration...</div>
        </Modal.Content>
      </Modal>
    );
  }

  const panes = [
    {
      menuItem: 'Basic Configuration',
      render: () => {
        console.log('Rendering Basic Configuration tab');
        console.log('Image options:', getArrayOptions(configOptions.images));
        console.log('CPU options:', getResourceOptions(configOptions.resources?.cpu));
        
        return (
          <Tab.Pane>
            <Form.Group widths="equal">
              <Form.Select
                label="Image"
                name="image"
                options={getArrayOptions(configOptions.images)}
                value={formData.image}
                onChange={handleChange}
                error={errors.image}
                required
                placeholder="Select Image"
              />
              <Form.Select
                label="Namespace"
                name="namespace"
                options={[{ 
                  key: configOptions.namespace,
                  text: configOptions.namespace,
                  value: configOptions.namespace 
                }]}
                value={formData.namespace}
                onChange={handleChange}
                placeholder="Select Namespace"
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
            </Form.Group>

            <Form.Group widths="equal">
              <Form.Select
                label="CPU"
                name="cpu"
                options={getResourceOptions(configOptions.resources?.cpu)}
                value={formData.cpu}
                onChange={handleChange}
                error={errors.cpu}
                required
                placeholder="Select CPU"
              />
              <Form.Select
                label="Memory"
                name="memory"
                options={getResourceOptions(configOptions.resources?.memory)}
                value={formData.memory}
                onChange={handleChange}
                error={errors.memory}
                required
                placeholder="Select Memory"
              />
              <Form.Select
                label="GPU"
                name="gpu"
                options={getResourceOptions(configOptions.resources?.gpu)}
                value={formData.gpu}
                onChange={handleChange}
                placeholder="Select GPU"
              />
            </Form.Group>
          </Tab.Pane>
        );
      }
    },
    {
      menuItem: 'Advanced Configuration',
      render: () => (
        <Tab.Pane>
          <Form.Group widths="equal">
            <Form.Select
              label="Model Repository"
              name="modelRepository"
              options={getArrayOptions(configOptions.model_repository)}
              value={formData.modelRepository}
              onChange={handleChange}
              placeholder="Select Model Repository"
            />
            <Form.Select
              label="Backend"
              name="backend"
              options={getArrayOptions(configOptions.backend)}
              value={formData.backend}
              onChange={handleChange}
              placeholder="Select Backend"
            />
          </Form.Group>

          <Form.Group widths="equal">
            <Form.Select
              label="HTTP Port"
              name="httpPort"
              options={getArrayOptions(configOptions.ports?.http)}
              value={formData.httpPort}
              onChange={handleChange}
              placeholder="Select HTTP Port"
            />
            <Form.Select
              label="gRPC Port"
              name="grpcPort"
              options={getArrayOptions(configOptions.ports?.grpc)}
              value={formData.grpcPort}
              onChange={handleChange}
              placeholder="Select gRPC Port"
            />
            <Form.Select
              label="Metrics Port"
              name="metricsPort"
              options={getArrayOptions(configOptions.ports?.metrics)}
              value={formData.metricsPort}
              onChange={handleChange}
              placeholder="Select Metrics Port"
            />
          </Form.Group>

          <Form.Group widths="equal">
            <Form.Select
              label="Log Verbose Level"
              name="logVerbose"
              options={getArrayOptions(configOptions.logging?.verbose)}
              value={formData.logVerbose}
              onChange={handleChange}
              placeholder="Select Log Level"
            />
            <Form.Select
              label="Log Format"
              name="logFormat"
              options={getArrayOptions(configOptions.logging?.formats)}
              value={formData.logFormat}
              onChange={handleChange}
              placeholder="Select Log Format"
            />
          </Form.Group>
        </Tab.Pane>
      ),
    },
  ];

  return (
    <Modal open={isOpen} onClose={onClose} size="large">
      <Modal.Header>Create Triton Service</Modal.Header>
      <Modal.Content>
        <Form loading={loading} error={Object.keys(errors).length > 0}>
          <Tab panes={panes} />
          {Object.keys(errors).length > 0 && (
            <Message
              error
              header="Validation Errors"
              list={Object.values(errors)}
            />
          )}
        </Form>
      </Modal.Content>
      <Modal.Actions>
        <Button onClick={onClose} disabled={loading}>Cancel</Button>
        <Button primary onClick={handleSubmit} loading={loading} disabled={loading}>
          Create
        </Button>
      </Modal.Actions>
    </Modal>
  );
};

export default CreateTritonForm;
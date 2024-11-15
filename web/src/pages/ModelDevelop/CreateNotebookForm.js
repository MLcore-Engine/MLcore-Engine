// src/pages/notebook/CreateNotebookForm.js

import React, { useState, useEffect } from 'react';
import { Modal, Form, Button, Message, Dropdown } from 'semantic-ui-react';
import { notebookAPI } from '../../api/notebookAPI';
import { toast } from 'react-toastify';

const CreateNotebookForm = ({ isOpen, onClose, onRefresh }) => {

  const userInfo = JSON.parse(localStorage.getItem('user') || '{}');
  const [formData, setFormData] = useState({
    user_id: userInfo.id || '', // Should be fetched from user context or authentication
    project_id: userInfo.projects?.[0]?.ID || '', // Could be a dropdown of available projects
    describe: '',
    namespace: 'jupyter',
    image: '',
    ide_type: '',
    working_dir: '',
    volume_mount: '',
    node_selector: 'notebook=true',
    image_pull_policy: 'IfNotPresent',
    resource_memory: '',
    resource_cpu: '',
    resource_gpu: 0,
    status: 'Creating',
    expand: '{}',
  });

  const [loading, setLoading] = useState(false);
  const [errors, setErrors] = useState({});
  const [dropdownOptions, setDropdownOptions] = useState({
    images: [],
    ideTypes: [],
    resourceMemories: [],
    resourceCPUs: [],
  });

  useEffect(() => {
    // Fetch options for dropdown menus
    // This could be replaced with actual API calls if options are dynamic
    setDropdownOptions({
      images: [
        { key: 'image1', text: 'notebook-cpu', value: 'gaoxin2020/notebook-tgqs:jupyter-ubuntu-cpu-base' },
        // Add more image options
      ],
      ideTypes: [
        { key: 'jupyter', text: 'Jupyter', value: 'jupyter' },
        { key: 'vscode', text: 'VSCode', value: 'vscode' },
        // Add more IDE types
      ],
      resourceMemories: [
        { key: '4G', text: '4G', value: '4G' },
        { key: '8G', text: '8G', value: '8G' },
        // Add more memory options
      ],
      resourceCPUs: [
        { key: '2', text: '2', value: '2' },
        { key: '4', text: '4', value: '4' },
        // Add more CPU options
      ],
    });
  }, []);

  const handleChange = (e, { name, value }) => {
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const validateForm = () => {
    const newErrors = {};
    if (!formData.image) newErrors.image = 'Image is required';
    if (!formData.ide_type) newErrors.ide_type = 'IDE Type is required';
    if (!formData.resource_memory) newErrors.resource_memory = 'Resource Memory is required';
    if (!formData.resource_cpu) newErrors.resource_cpu = 'Resource CPU is required';
    // Add more validation as needed
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async () => {
    if (!validateForm()) return;
    setLoading(true);
    try {
      await notebookAPI.createNotebook(formData);
      toast.success('Notebook created successfully');
      onClose();
      onRefresh();
    } catch (error) {
      toast.error(error.message || 'Failed to create notebook');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Modal open={isOpen} onClose={onClose}>
      <Modal.Header>Create Notebook</Modal.Header>
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
            label="IDE Type"
            placeholder="Select IDE Type"
            options={dropdownOptions.ideTypes}
            name="ide_type"
            onChange={handleChange}
            value={formData.ide_type}
            error={errors.ide_type ? { content: errors.ide_type } : null}
            selection
          />
          <Form.Field
            control={Dropdown}
            label="Resource Memory"
            placeholder="Select Memory"
            options={dropdownOptions.resourceMemories}
            name="resource_memory"
            onChange={handleChange}
            value={formData.resource_memory}
            error={errors.resource_memory ? { content: errors.resource_memory } : null}
            selection
          />
          <Form.Field
            control={Dropdown}
            label="Resource CPU"
            placeholder="Select CPU"
            options={dropdownOptions.resourceCPUs}
            name="resource_cpu"
            onChange={handleChange}
            value={formData.resource_cpu}
            error={errors.resource_cpu ? { content: errors.resource_cpu } : null}
            selection
          />
          <Form.Input
            label="Resource GPU"
            name="resource_gpu"
            type="number"
            min="0"
            onChange={handleChange}
            value={formData.resource_gpu}
          />
          <Form.TextArea
            label="Description"
            name="describe"
            onChange={handleChange}
            value={formData.describe}
          />
          {/* Add more fields as needed */}
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

export default CreateNotebookForm;
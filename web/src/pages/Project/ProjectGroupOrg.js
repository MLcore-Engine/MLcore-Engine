import React, { useState } from 'react';
import { Button, Icon, Modal, Form, Message } from 'semantic-ui-react';
import DataList from '../../components/common/DataList';
import { useProjects } from '../../context/ProjectContext';
import { toast } from 'react-toastify';

const ProjectGroupOrg = () => {
  const {
    projects,
    loading,
    createProject,
    updateProject,
    deleteProject,
  } = useProjects();

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [modalType, setModalType] = useState('');
  const [projectData, setProjectData] = useState({ name: '', description: '' });
  const [modalError, setModalError] = useState('');
  const [selectedProjectId, setSelectedProjectId] = useState(null);
  const [isSubmitting, setIsSubmitting] = useState(false);


  const handleAdd = () => {
    setModalType('add');
    setProjectData({ name: '', description: '' });
    setModalError('');
    setIsModalOpen(true);
  };

  const handleEdit = (project) => {
    setModalType('edit');
    setProjectData({ name: project.name, description: project.description });
    setSelectedProjectId(project.ID);
    setModalError('');
    setIsModalOpen(true);
  };

  const handleDelete = async (item) => {
    try {
      await deleteProject(item.ID);
      toast.info('project deleted');
    } catch (err) {
      console.error(err.message);
      toast.error(err.message || 'delete project failed');
    }
  };

  const handleSubmit = async () => {
    if (!projectData.name.trim()) {
      setModalError('project name is required');  
      return;
    }

    setIsSubmitting(true);
    setModalError('');

    try {
      if (modalType === 'add') {
        await createProject(projectData);
      } else if (modalType === 'edit' && selectedProjectId) {
        await updateProject(selectedProjectId, projectData);
      }
      setIsModalOpen(false);
      setProjectData({ name: '', description: '' });
      setSelectedProjectId(null);
      setModalError('');
    } catch (err) {
      setModalError(err.message || 'operation failed');
    } finally {
      setIsSubmitting(false);
    }
  };

  const columns = ['name', 'description'];
  const columnNames = { name: '项目名称', description: '描述' };

  if (loading) return <div className='p-4'>loading...</div>;

  return (
    <div className='p-4'>
      <DataList
        title="Project List"
        data={projects || []}
        columns={columns}
        columnNames={columnNames}
        onAdd={handleAdd}
        onEdit={handleEdit}
        onDelete={handleDelete}
      />

      <Modal
              open={isModalOpen}
              onClose={() => setIsModalOpen(false)}
              closeOnDimmerClick={!isSubmitting}
              closeOnEscape={!isSubmitting}
              aria-labelledby="modal-header"
      >
        <Modal.Header id="modal-header">{modalType === 'add' ? '添加项目' : '编辑项目'}</Modal.Header>
        <Modal.Content>
          <Form error={!!modalError}>
            <Form.Input
              label="项目名称"
              value={projectData.name}
              onChange={(e, { value }) => setProjectData({ ...projectData, name: value })}
              placeholder="输入项目名称"
              aria-label="project-name"
              required
            />
            <Form.TextArea
              label="描述"
              value={projectData.description}
              onChange={(e, { value }) => setProjectData({ ...projectData, description: value })}
              placeholder="输入项目描述"
              aria-label="project-description"
            />
            {modalError && <Message error content={modalError} />}
          </Form>
        </Modal.Content>
        <Modal.Actions>
          <Button 
            onClick={() => setIsModalOpen(false)}
            disabled={isSubmitting}
            aria-label="cancel-button"
          >
            取消
          </Button>
          <Button 
            primary 
            onClick={handleSubmit}
            loading={isSubmitting}
            disabled={isSubmitting}
            aria-label={modalType === 'add' ? '添加' : '保存'}
          >
            {modalType === 'add' ? '添加' : '保存'}
          </Button>
        </Modal.Actions>
      </Modal>
    </div>
  );
};

export default ProjectGroupOrg;
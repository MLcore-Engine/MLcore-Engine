import React, { useState } from 'react';
import { Button, Icon, Modal, Form, Message } from 'semantic-ui-react';
import DataList from '../../components/common/DataList';
import { useProjects } from '../../context/ProjectContext';
import { toast } from 'react-toastify';

const ProjectGroupOrg = () => {
  const {
    projects,
    loading,
    error,
    createProject,
    updateProject,
    deleteProject,
  } = useProjects();

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [modalType, setModalType] = useState('');
  const [projectData, setProjectData] = useState({ name: '', description: '' });
  const [modalError, setModalError] = useState('');

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

  const handleDelete = async (id) => {
    try {
      await deleteProject(id);
      toast.success('Project deleted successfully.');
    } catch (err) {
      console.error(err.message);
      toast.error(err.message || 'Failed to delete project.');
    }
  };

  const [selectedProjectId, setSelectedProjectId] = useState(null);

  const handleSubmit = async () => {
    if (!projectData.name) {
      setModalError('项目名称不能为空');
      return;
    }
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
      setModalError(err.message);
    }
  };

  const columns = ['name', 'description'];

  if (loading) return <div>加载中...</div>;
  if (error) return <div>错误: {error}</div>;

  return (
    <div>
      <DataList
        title="项目管理"
        data={projects}
        columns={columns}
        onAdd={handleAdd}
        onEdit={handleEdit}
        onDelete={(project) => handleDelete(project.ID)}
      />

      <Modal open={isModalOpen} onClose={() => setIsModalOpen(false)}>
        <Modal.Header>{modalType === 'add' ? '添加项目' : '编辑项目'}</Modal.Header>
        <Modal.Content>
          <Form error={!!modalError}>
            <Form.Input
              label="项目名称"
              value={projectData.name}
              onChange={(e, { value }) => setProjectData({ ...projectData, name: value })}
              placeholder="输入项目名称"
            />
            <Form.TextArea
              label="描述"
              value={projectData.description}
              onChange={(e, { value }) => setProjectData({ ...projectData, description: value })}
              placeholder="输入项目描述"
            />
            {modalError && <Message error content={modalError} />}
          </Form>
        </Modal.Content>
        <Modal.Actions>
          <Button onClick={() => setIsModalOpen(false)}>取消</Button>
          <Button primary onClick={handleSubmit}>
            {modalType === 'add' ? '添加' : '保存'}
          </Button>
        </Modal.Actions>
      </Modal>
    </div>
  );
};

export default ProjectGroupOrg;
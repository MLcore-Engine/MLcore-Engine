import React, { useState } from 'react';
import { Table, Button, Icon, Modal, Form, Message } from 'semantic-ui-react';
import { useProjects } from '../../context/ProjectContext';

const ProjectManage = () => {
  const {
    projects,
    loading,
    error,
    addProjectMember,
    removeProjectMember,
    updateProjectMemberRole, // Assuming you might use this
  } = useProjects();

  const [selectedProject, setSelectedProject] = useState(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [newUser, setNewUser] = useState({ userId: '', role: 0 });
  const [modalError, setModalError] = useState('');

  const handleAddUser = (project) => {
    setSelectedProject(project);
    setModalError('');
    setNewUser({ userId: '', role: 0 });
    setIsModalOpen(true);
  };

  const handleRemoveUser = async (projectId, userId) => {
    try {
      await removeProjectMember(projectId, userId);
    } catch (err) {
      console.error(err.message);
    }
  };

  const handleSubmitNewUser = async () => {
    if (!newUser.userId) {
      setModalError('用户ID不能为空');
      return;
    }
    try {
      await addProjectMember(selectedProject.ID, newUser);
      setIsModalOpen(false);
      setNewUser({ userId: '', role: 0 });
      setModalError('');
    } catch (err) {
      setModalError(err.message);
    }
  };

  if (loading) return <div>加载中...</div>;
  if (error) return <div>错误: {error}</div>;

  return (
    <div>
      <h2>项目成员管理</h2>
      {projects.map((project) => (
        <div key={project.ID} style={{ marginBottom: '2em' }}>
          <h3>{project.name}</h3>
          <Table celled>
            <Table.Header>
              <Table.Row>
                <Table.HeaderCell>用户名</Table.HeaderCell>
                <Table.HeaderCell>角色</Table.HeaderCell>
                <Table.HeaderCell>操作</Table.HeaderCell>
              </Table.Row>
            </Table.Header>
            <Table.Body>
              {project.users && project.users.length > 0 ? (
                project.users.map((user) => (
                  <Table.Row key={user.ID}>
                    <Table.Cell>{user.username}</Table.Cell>
                    <Table.Cell>{getRoleName(user.role)}</Table.Cell>
                    <Table.Cell>
                      <Button
                        icon
                        color="red"
                        onClick={() => handleRemoveUser(project.ID, user.ID)}
                      >
                        <Icon name="trash" />
                      </Button>
                    </Table.Cell>
                  </Table.Row>
                ))
              ) : (
                <Table.Row>
                  <Table.Cell colSpan="3">暂无成员</Table.Cell>
                </Table.Row>
              )}
            </Table.Body>
          </Table>
          <Button primary onClick={() => handleAddUser(project)}>
            <Icon name="add user" /> 添加用户
          </Button>
        </div>
      ))}

      <Modal open={isModalOpen} onClose={() => setIsModalOpen(false)}>
        <Modal.Header>添加用户到 {selectedProject?.name}</Modal.Header>
        <Modal.Content>
          <Form error={!!modalError}>
            <Form.Input
              label="用户ID"
              value={newUser.userId}
              onChange={(e, { value }) => setNewUser({ ...newUser, userId: value })}
              placeholder="输入用户ID"
            />
            <Form.Select
              label="角色"
              options={[
                { key: '0', text: '用户', value: 0 },
                { key: '1', text: '管理员', value: 1 },
                { key: '2', text: '根用户', value: 2 },
              ]}
              value={newUser.role}
              onChange={(e, { value }) => setNewUser({ ...newUser, role: value })}
              placeholder="选择角色"
            />
            {modalError && <Message error content={modalError} />}
          </Form>
        </Modal.Content>
        <Modal.Actions>
          <Button onClick={() => setIsModalOpen(false)}>取消</Button>
          <Button primary onClick={handleSubmitNewUser}>
            添加
          </Button>
        </Modal.Actions>
      </Modal>
    </div>
  );
};

// Helper function to convert role number to role name
const getRoleName = (role) => {
  switch (role) {
    case 0:
      return '用户';
    case 1:
      return '管理员';
    case 2:
      return '根用户';
    default:
      return '未知';
  }
};

export default ProjectManage;
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
  } = useProjects();

  const [selectedProject, setSelectedProject] = useState(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [newUser, setNewUser] = useState({ username: '', role: 'user' });
  const [modalError, setModalError] = useState('');

  const handleAddUser = (project) => {
    setSelectedProject(project);
    setModalError('');
    setNewUser({ username: '', role: 'user' });
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
    if (!newUser.username) {
      setModalError('用户名不能为空');
      return;
    }
    try {
      await addProjectMember(selectedProject.id, newUser);
      setIsModalOpen(false);
      setNewUser({ username: '', role: 'user' });
      setModalError('');
    } catch (err) {
      setModalError(err.message);
    }
  };

  if (loading) return <div>加载中...</div>;
  if (error) return <div>错误: {error}</div>;

  return (
    <div>
      <h2>项目管理</h2>
      {projects.map((project) => (
        <div key={project.id} style={{ marginBottom: '2em' }}>
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
              {project.members && project.members.length > 0 ? (
                project.members.map((user) => (
                  <Table.Row key={user.userId}>
                    <Table.Cell>{user.username}</Table.Cell>
                    <Table.Cell>{user.role}</Table.Cell>
                    <Table.Cell>
                      <Button
                        icon
                        color="red"
                        onClick={() => handleRemoveUser(project.id, user.userId)}
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
              label="用户名"
              value={newUser.username}
              onChange={(e, { value }) => setNewUser({ ...newUser, username: value })}
              placeholder="输入用户名"
            />
            <Form.Select
              label="角色"
              options={[
                { key: 'user', text: '用户', value: 'user' },
                { key: 'admin', text: '管理员', value: 'admin' },
                { key: 'root', text: '根用户', value: 'root' },
              ]}
              value={newUser.role}
              onChange={(e, { value }) => setNewUser({ ...newUser, role: value })}
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

export default ProjectManage;
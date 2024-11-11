import React, { useState, useEffect, useCallback } from 'react';
import { Button, Modal, Form, Message, Dropdown, Pagination, Input } from 'semantic-ui-react';
import { useProjects } from '../../context/ProjectContext';
import { userAPI } from '../../api/userAPI';
import { toast } from 'react-toastify';
import DataList from '../../components/common/DataList'; // 确保导入路径正确

// 辅助函数，将角色编号转换为角色名称
const getRoleName = (role) => {
  switch (role) {
    case 1:
      return '用户';
    case 10:
      return '管理员';
    default:
      return '未知';
  }
};

const ProjectManage = () => {
  const {
    projects,
    loading,
    addProjectMember,
    removeProjectMember,
    updateProjectMemberRole,
  } = useProjects();

  // 状态管理
  const [selectedProject, setSelectedProject] = useState(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [newUser, setNewUser] = useState({ userID: '', role: 1 });
  const [modalError, setModalError] = useState('');
  const [users, setUsers] = useState([]);
  const [loadingUsers, setLoadingUsers] = useState(false);
  const [currentPage, setCurrentPage] = useState(1); // 使用1为第一页
  const [searchTerm, setSearchTerm] = useState('');
  const [isSearching, setIsSearching] = useState(false);
  const [totalPages, setTotalPages] = useState(1); // 动态获取总页数

  // 使用 useCallback 优化搜索处理函数
  const fetchUsers = useCallback(async (searchValue, page) => {
    setLoadingUsers(true);
    try {
      let response;
      if (searchValue.trim() !== '') {
        response = await userAPI.searchUsers(searchValue);
        setIsSearching(true);
        setTotalPages(1); // 搜索结果不考虑分页
      } else {
        response = await userAPI.getUsers(page - 1); // 后端页码从0开始
        setIsSearching(false);
        setTotalPages(response.total || 10); // 默认10页
      }
      const userOptions = response.data.map(user => ({
        key: user.ID,
        text: `${user.username}`,
        value: user.ID
      }));
      setUsers(userOptions);
    } catch (err) {
      console.error('Failed to fetch users:', err);
      toast.error(err.message || '获取用户列表失败');
    } finally {
      setLoadingUsers(false);
    }
  }, []);

  // 修改 useEffect，添加防抖
  useEffect(() => {
    if (!isModalOpen) return; // 仅在模态框打开时获取用户

    const timer = setTimeout(() => {
      fetchUsers(searchTerm, currentPage);
    }, 500); // 始终使用防抖

    return () => clearTimeout(timer);
  }, [searchTerm, currentPage, isModalOpen, fetchUsers]);

  // 优化搜索处理函数
  const handleSearchChange = useCallback((e, { value }) => {
    setSearchTerm(value);
    setCurrentPage(1); // 重置为第一页
  }, []);

  // 优化分页处理函数
  const handlePageChange = useCallback((e, { activePage }) => {
    setCurrentPage(activePage);
  }, []);

  // 优化模态框打开处理
  const handleAddUser = useCallback((project) => {
    setSelectedProject(project);
    setModalError('');
    setNewUser({ userID: '', role: 1 });
    setSearchTerm(''); // 重置搜索条件
    setCurrentPage(1); // 重置页码为1
    setIsModalOpen(true);
  }, []);

  // 优化模态框关闭处理
  const handleModalClose = useCallback(() => {
    setIsModalOpen(false);
    setSearchTerm('');
    setCurrentPage(1);
    setUsers([]);
    setModalError('');
  }, []);

  const handleRemoveUser = async (projectId, userID) => {
      try {
        await removeProjectMember(projectId, userID);
      } catch (err) {
        console.error(err.message);
        toast.error(err.message || '移除成员失败');
      }
  };

  const handleSubmitNewUser = async () => {
    if (!newUser.userID) {
      setModalError('请选择用户');
      return;
    }
    try {
      await addProjectMember(selectedProject.ID, newUser);
      handleModalClose();
    } catch (err) {
      setModalError(err.message || '添加成员失败');
      toast.error(err.message || '添加成员失败');
    }
  };

  if (loading) return <div>加载中...</div>;

  return (
    <div className="p-4">
      <h2>Project Member Management</h2>
      {projects.map((project) => (
        <div key={project.ID} className="project-section p-4 border border-gray-300 rounded-md mb-6">
          <DataList
            title={project.name}
            data={project.users || []}
            columns={['username', 'role']}
            columnNames={{ username: '用户名', role: '角色' }}
            onAdd={() => handleAddUser(project)}
            onDelete={(user) => handleRemoveUser(project.ID, user.ID)}
            renderRow={(user) => ({
              ...user,
              role: getRoleName(user.role),
            })}
          />
        </div>
      ))}

      {/* 添加成员的模态框 */}
      <Modal
        open={isModalOpen}
        onClose={handleModalClose}
        size="small"
        closeOnDimmerClick={!isModalOpen}
        closeOnEscape={!isModalOpen}
        aria-labelledby="modal-header"
      >
        <Modal.Header id="modal-header">添加用户到 {selectedProject?.name}</Modal.Header>
        <Modal.Content>
          <Form error={!!modalError}>
            <Form.Field>
              <label>Search User</label>
              <Input
                icon="search"
                placeholder="Search User..."
                value={searchTerm}
                onChange={handleSearchChange}
                loading={loadingUsers}
                aria-label="search-user"
              />
            </Form.Field>
            <Form.Field>
              <label>Select User</label>
              <Dropdown
                placeholder="Select User"
                fluid
                search
                selection
                loading={loadingUsers}
                options={users}
                value={newUser.userID}
                onChange={(e, { value }) => setNewUser({ ...newUser, userID: value })}
                disabled={loadingUsers}
                noResultsMessage={loadingUsers ? "Loading..." : "No matching users"}
                aria-label="select-user"
              />
            </Form.Field>
            <Form.Select
              label="角色"
              options={[
                { key: '0', text: '用户', value: 1 },
                { key: '1', text: '管理员', value: 10 },
              ]}
              value={newUser.role}
              onChange={(e, { value }) => setNewUser({ ...newUser, role: value })}
              placeholder="选择角色"
              aria-label="选择角色"
            />
            {modalError && <Message error content={modalError} />}
          </Form>
          {!isSearching && users.length > 0 && (
            <div className="mt-4 text-center">
              <Pagination
                activePage={currentPage}
                onPageChange={handlePageChange}
                totalPages={totalPages}
                disabled={loadingUsers}
                aria-label="pagination"
              />
            </div>
          )}
        </Modal.Content>
        <Modal.Actions>
          <Button onClick={handleModalClose} aria-label="cancel-button" disabled={loadingUsers}>
            取消
          </Button>
          <Button
            primary
            onClick={handleSubmitNewUser}
            disabled={loadingUsers || !newUser.userID}
            aria-label="add-button"
          >
            添加
          </Button>
        </Modal.Actions>
      </Modal>

      <style jsx>{`
        .project-section {
          margin-bottom: 1.5em;
          border: 1px solid #e0e0e0;
          border-radius: 4px;
        }
      `}</style>
    </div>
  );
};


export default ProjectManage;
import React, { useState, useEffect, useCallback } from 'react';
import { Button, Modal, Form, Message, Dropdown, Pagination, Input } from 'semantic-ui-react';
import { useProjects } from '../../context/ProjectContext';
import { userAPI } from '../../api/userAPI';
import { toast } from 'react-toastify';
import DataList from '../../components/common/DataList';
import { getProjectRoleName, ROLE } from '../../utils/roleUtils';
import RoleSelector from '../../components/common/RoleSelector';

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
  const [newUser, setNewUser] = useState({ userId: '', role: ROLE.COMMON });
  const [modalError, setModalError] = useState('');
  const [users, setUsers] = useState([]);
  const [loadingUsers, setLoadingUsers] = useState(false);
  const [currentPage, setCurrentPage] = useState(1); // 使用1为第一页
  const [searchTerm, setSearchTerm] = useState('');
  const [isSearching, setIsSearching] = useState(false);
  const [totalPages, setTotalPages] = useState(1); // 动态获取总页数
  const [isEditRoleModalOpen, setIsEditRoleModalOpen] = useState(false);
  const [editingMember, setEditingMember] = useState(null);
  const [newRole, setNewRole] = useState(null);

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
        const totalRecords = response.total || 0;
        const pageSize = response.limit || 10;
        const calculatedTotalPages = Math.ceil(totalRecords / pageSize);
        setTotalPages(calculatedTotalPages || 1); // 确保至少有1页
      }
      const userOptions = response.data.map(user => ({
        key: user.id,
        text: `${user.username}`,
        value: user.id
      }));
      setUsers(userOptions);
    } catch (err) {
      console.error('Failed to fetch users:', err);
      toast.error(err.message || '获取用户列表失败');
    } finally {
      setLoadingUsers(false);
    }
  }, []);

  const handleEditRole = (member, projectId) => {
    setEditingMember(member);
    setNewRole(member.role);
    setIsEditRoleModalOpen(true);
  };

  const handleSaveRole = async () => {
    try {
      await updateProjectMemberRole(
        editingMember.projectId,
        editingMember.userId,
        newRole
      );
      setIsEditRoleModalOpen(false);
    } catch (err) {
      toast.error(err.message || '角色更新失败');
    }
  };

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
    setNewUser({ userId: '', role: ROLE.COMMON });
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

  const handleRemoveUser = async (projectId, userId) => {
    try {
      console.log('移除用户，参数检查:', {
        projectId,
        userId,
        projectIdType: typeof projectId,
        userIdType: typeof userId
      });

      if (!userId) {
        throw new Error('用户ID不能为空');
      }

      // 确保ID为字符串类型
      const strProjectId = String(projectId);
      const strUserId = String(userId);

      console.log('最终移除用户参数:', {
        projectId: strProjectId,
        userId: strUserId
      });

      await removeProjectMember(strProjectId, strUserId);
    } catch (err) {
      // console.error('移除项目成员失败:', err);
      toast.info(err.message || '移除成员失败');
    }
  };

  const handleSubmitNewUser = async () => {
    if (!newUser.userId) {
      setModalError('请选择用户');
      return;
    }
    try {
      await addProjectMember(selectedProject.id, newUser);
      handleModalClose();
    } catch (err) {
      setModalError(err.message || '添加成员失败');
      toast.error(err.message || '添加成员失败');
    }
  };

  // 渲染项目成员行数据
  const renderMemberRow = (member) => {
    // 调试日志
    // console.log('渲染成员数据原始值:', member);

    // 获取用户ID，尝试多种可能的属性
    const userId = member.userId || member.userID || member.id || member.ID;

    // 构建增强后的成员数据
    const enhancedMember = {
      ...member,
      // 确保关键字段存在
      username: member.username || '未知用户',
      role: getProjectRoleName(member.role !== undefined ? member.role : 0),
      // 确保各种ID字段都存在
      id: userId,
      ID: userId,
      userId: userId,
    };

    // console.log('增强后的成员数据:', enhancedMember);
    return enhancedMember;
  };

  if (loading) return <div>加载中...</div>;

  return (
    <div className="p-4">
      <h2>项目成员</h2>
      {projects.map((project) => (
        <div key={project.id} className="project-section p-4 border border-gray-300 rounded-md mb-6">
          <DataList
            title={project.name}
            data={project.users || []}
            columns={['username', 'role']}
            columnNames={{ username: '用户名', role: '角色' }}
            onAdd={() => handleAddUser(project)}
            onEdit={(user) => handleEditRole(user, project.id)}
            onDelete={(user) => {
              // 详细打印用户对象，包括所有属性
              console.log('用户对象详情:', user);
              console.log('用户对象键名:', Object.keys(user));

              // 尝试不同方式获取用户ID
              // 1. 直接访问常见属性名
              let userId = user.userId || user.userID || user.id || user.ID;

              // 2. 如果是数字类型，转为字符串
              if (typeof user.userId === 'number') userId = String(user.userId);
              if (typeof user.userID === 'number') userId = String(user.userID);
              if (typeof user.id === 'number') userId = String(user.id);
              if (typeof user.ID === 'number') userId = String(user.ID);

              // 3. 遍历所有属性查找可能的ID字段
              if (!userId) {
                for (const key of Object.keys(user)) {
                  // 检查含有"id"的键名(不区分大小写)
                  if (key.toLowerCase().includes('id') && user[key]) {
                    console.log('找到可能的ID字段:', key, user[key]);
                    userId = user[key];
                    break;
                  }
                }
              }

              // 输出最终选择的用户ID
              // console.log('最终选择的用户ID:', userId);

              if (userId) {
                handleRemoveUser(project.id, userId);
              } else {
                // console.error('找不到用户ID:', user);
                toast.info('无法删除：找不到用户ID');
              }
            }}
            renderRow={renderMemberRow}
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
                value={newUser.userId}
                onChange={(e, { value }) => setNewUser({ ...newUser, userId: value })}
                disabled={loadingUsers}
                noResultsMessage={loadingUsers ? "Loading..." : "No matching users"}
                aria-label="select-user"
              />
            </Form.Field>
            <Form.Field>
              <label>角色</label>
              <RoleSelector
                value={newUser.role}
                onChange={(e, { value }) => setNewUser(prev => ({ ...prev, role: value }))}
                isProjectRole={true}
              />
            </Form.Field>
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
            disabled={loadingUsers || !newUser.userId}
            aria-label="add-button"
          >
            添加
          </Button>
        </Modal.Actions>
      </Modal>


      {/* 编辑角色模态框 */}
      <Modal
        open={isEditRoleModalOpen}
        onClose={() => setIsEditRoleModalOpen(false)}
        size="small"
      >
        <Modal.Header>编辑用户角色</Modal.Header>
        <Modal.Content>
          <Form>
            <Form.Field>
              <label>用户名</label>
              <Input value={editingMember?.username || ''} disabled />
            </Form.Field>
            <Form.Field>
              <label>角色</label>
              <RoleSelector
                value={newRole}
                onChange={(e, { value }) => setNewRole(value)}
                isProjectRole={true}
              />
            </Form.Field>
          </Form>
        </Modal.Content>
        <Modal.Actions>
          <Button onClick={() => setIsEditRoleModalOpen(false)}>
            取消
          </Button>
          <Button primary onClick={handleSaveRole}>
            保存
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
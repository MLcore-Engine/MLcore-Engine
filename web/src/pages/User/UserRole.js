import React, { useState, useEffect } from 'react';
import { Button, Modal, Message } from 'semantic-ui-react';
import { useUsers } from '../../context/UserContext';
import { toast } from 'react-toastify';
import UserTable from '../../components/common/UserTable';
import RoleSelector from '../../components/common/RoleSelector';

// 辅助函数：将角色编号转换为角色名称
const getRoleName = (role) => {
  switch (role) {
    case 1:
      return '普通用户';
    case 10:
      return '管理员';
    case 100:
      return '超级管理员';
    default:
      return '未知';
  }
};

// 角色选项
const roleOptions = [
  { key: '1', text: '普通用户', value: 1 },
  { key: '10', text: '管理员', value: 10 },
  { key: '100', text: '超级管理员', value: 100 },
];

const UserRole = () => {
  const {
    users,
    loading,
    fetchUsers,
    updateUserRole
  } = useUsers();

  // 状态管理
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedUser, setSelectedUser] = useState(null);
  const [newRole, setNewRole] = useState(1);
  const [modalError, setModalError] = useState('');

  // 首次加载时获取用户列表
  useEffect(() => {
    fetchUsers();
  }, [fetchUsers]);

  // 打开修改角色模态框
  const handleEditRole = (user) => {
    setSelectedUser(user);
    setNewRole(user.role);
    setIsModalOpen(true);
    setModalError('');
  };

  // 关闭模态框
  const handleModalClose = () => {
    setIsModalOpen(false);
    setModalError('');
  };

  // 提交角色更新
  const handleUpdateRole = async () => {
    if (!selectedUser) return;

    try {
      await updateUserRole(selectedUser.id, newRole);
      toast.success('用户角色已更新');
      handleModalClose();
    } catch (err) {
      setModalError(err.message || '更新角色失败');
    }
  };

  // 表格列配置
  const columns = ['username', 'display_name', 'email', 'role'];
  const columnNames = {
    username: '用户名',
    display_name: '显示名称',
    email: '邮箱',
    role: '当前角色'
  };

  // 自定义操作按钮
  const customActions = [
    {
      icon: 'edit',
      color: 'blue',
      onClick: handleEditRole
    }
  ];

  return (
    <div className="p-4">
      <h2>用户角色管理</h2>
      <UserTable
        title="角色分配"
        users={users}
        loading={loading}
        columns={columns}
        columnNames={columnNames}
        customActions={customActions}
      />

      {/* 角色编辑模态框 */}
      <Modal
        open={isModalOpen}
        onClose={handleModalClose}
        size="mini"
        closeOnDimmerClick={false}
        aria-labelledby="role-modal-header"
      >
        <Modal.Header id="role-modal-header">更新用户角色</Modal.Header>
        <Modal.Content>
          {selectedUser && (
            <>
              <p>
                <strong>用户:</strong> {selectedUser.username}
              </p>
              <p>
                <strong>当前角色:</strong> {selectedUser.role_name || '未知'}
              </p>
              <RoleSelector
                value={newRole}
                onChange={(e, { value }) => setNewRole(value)}
                placeholder="选择新角色"
              />
              {modalError && <Message error content={modalError} style={{ marginTop: '10px' }} />}
            </>
          )}
        </Modal.Content>
        <Modal.Actions>
          <Button onClick={handleModalClose}>取消</Button>
          <Button primary onClick={handleUpdateRole}>更新</Button>
        </Modal.Actions>
      </Modal>
    </div>
  );
};

export default UserRole; 
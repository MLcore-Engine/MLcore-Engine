import React, { useState, useEffect } from 'react';
import { Button, Modal, Form, Message, Input } from 'semantic-ui-react';
import { useUsers } from '../../context/UserContext';
import { toast } from 'react-toastify';
import UserTable from '../../components/common/UserTable';
import RoleSelector from '../../components/common/RoleSelector';
import { ROLE } from '../../utils/roleUtils';

const UserList = () => {
    const {
        users,
        loading,
        fetchUsers,
        addUser,
        deleteUser,
        toggleUserStatus
    } = useUsers();

    // 状态管理
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [newUserData, setNewUserData] = useState({
        username: '',
        display_name: '',
        email: '',
        password: '',
        role: ROLE.COMMON,
        status: 1
    });
    const [modalError, setModalError] = useState('');
    const [userToDelete, setUserToDelete] = useState(null);

    // 首次加载时获取用户列表
    useEffect(() => {
        fetchUsers();
    }, [fetchUsers]);

    // 打开添加用户模态框
    const handleAddUser = () => {
        setIsModalOpen(true);
        setModalError('');
        setNewUserData({
            username: '',
            display_name: '',
            email: '',
            password: '',
            role: ROLE.COMMON,
            status: 1
        });
    };

    // 关闭模态框
    const handleModalClose = () => {
        setIsModalOpen(false);
        setModalError('');
    };

    // 处理表单字段变化
    const handleInputChange = (e, { name, value }) => {
        setNewUserData(prev => ({ ...prev, [name]: value }));
    };

    // 提交新用户表单
    const handleSubmitNewUser = async () => {
        // 表单验证
        if (!newUserData.username || !newUserData.password || !newUserData.email) {
            setModalError('用户名、密码和邮箱为必填项');
            return;
        }

        try {
            await addUser(newUserData);
            handleModalClose();
            toast.success('用户添加成功');
        } catch (err) {
            setModalError(err.message || '添加用户失败');
        }
    };

    // 处理删除用户点击
    const handleDeleteClick = (user) => {
        setUserToDelete(user);
        // 这里应该显示确认对话框，简化起见直接调用删除
        handleConfirmDelete(user);
    };

    // 执行删除用户
    const handleConfirmDelete = async (user) => {
        if (user) {
            try {
                await deleteUser(user.id);
                toast.success('用户已删除');
            } catch (err) {
                toast.error(err.message || '删除用户失败');
            }
        }
    };

    // 处理状态切换
    const handleToggleStatus = async (user) => {
        try {
            const action = user.status === 1 ? 'disable' : 'enable';
            await toggleUserStatus(user.id, action);
            toast.success(`用户已${user.status === 1 ? '禁用' : '启用'}`);
        } catch (err) {
            toast.error(err.message || `${user.status === 1 ? '禁用' : '启用'}用户失败`);
        }
    };

    // 表格列配置
    const columns = ['username', 'display_name', 'email', 'role', 'status'];
    const columnNames = {
        username: '用户名',
        display_name: '显示名称',
        email: '邮箱',
        role: '角色',
        status: '状态'
    };

    // 自定义操作按钮
    const customActions = [
        {
            icon: 'power',
            color: 'orange',
            onClick: handleToggleStatus,
            confirm: {
                content: (user) => `确定要${user.status === 1 ? '禁用' : '启用'}用户 "${user.username}" 吗?`,
                confirmButton: '确认',
                cancelButton: '取消'
            }
        }
    ];

    return (
        <div className="p-4">
            <h2>用户列表</h2>
            <Button primary onClick={handleAddUser} style={{ marginBottom: '1rem' }}>添加用户</Button>
            
            <UserTable
                title="用户管理"
                users={users}
                loading={loading}
                columns={columns}
                columnNames={columnNames}
                onDelete={handleDeleteClick}
                customActions={customActions}
            />

            {/* 添加用户的模态框 */}
            <Modal
                open={isModalOpen}
                onClose={handleModalClose}
                size="small"
                closeOnDimmerClick={false}
                aria-labelledby="modal-header"
            >
                <Modal.Header id="modal-header">添加新用户</Modal.Header>
                <Modal.Content>
                    <Form error={!!modalError}>
                        <Form.Field>
                            <label>用户名</label>
                            <Input
                                placeholder="输入用户名..."
                                name="username"
                                value={newUserData.username}
                                onChange={handleInputChange}
                            />
                        </Form.Field>
                        <Form.Field>
                            <label>显示名称</label>
                            <Input
                                placeholder="输入显示名称..."
                                name="display_name"
                                value={newUserData.display_name}
                                onChange={handleInputChange}
                            />
                        </Form.Field>
                        <Form.Field>
                            <label>电子邮箱</label>
                            <Input
                                placeholder="输入电子邮箱..."
                                name="email"
                                value={newUserData.email}
                                onChange={handleInputChange}
                            />
                        </Form.Field>
                        <Form.Field>
                            <label>密码</label>
                            <Input
                                type="password"
                                placeholder="输入密码..."
                                name="password"
                                value={newUserData.password}
                                onChange={handleInputChange}
                            />
                        </Form.Field>
                        <Form.Field>
                            <label>角色</label>
                            <RoleSelector
                                value={newUserData.role}
                                onChange={(e, { value }) => setNewUserData(prev => ({ ...prev, role: value }))}
                            />
                        </Form.Field>
                        {modalError && <Message error content={modalError} />}
                    </Form>
                </Modal.Content>
                <Modal.Actions>
                    <Button onClick={handleModalClose}>取消</Button>
                    <Button primary onClick={handleSubmitNewUser}>添加</Button>
                </Modal.Actions>
            </Modal>
        </div>
    );
};

export default UserList;

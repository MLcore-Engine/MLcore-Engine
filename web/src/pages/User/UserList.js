import React, { useState, useEffect, useCallback } from 'react';
import { Dimmer, Loader } from 'semantic-ui-react';
import { userAPI } from '../../api/userAPI';
import { toast } from 'react-toastify';
import DataList from '../../components/common/DataList';

const UserList = () => {
    // 状态管理
    const [users, setUsers] = useState([]);
    const [loading, setLoading] = useState(false);
    const [currentPage, setCurrentPage] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [totalItems, setTotalItems] = useState(0);

    // 获取用户列表数据
    const fetchUsers = useCallback(async (page = currentPage) => {
        setLoading(true);
        try {
            const response = await userAPI.getUsers(page, pageSize);
            if (response.success) {
                setUsers(response.data);
                setTotalItems(response.total);
                setCurrentPage(response.page);
                setPageSize(response.limit);
            } else {
                toast.error(response.message || '获取用户列表失败');
            }
        } catch (error) {
            console.error('获取用户列表错误:', error);
            toast.error(error.message || '获取用户列表失败');
        } finally {
            setLoading(false);
        }
    }, [currentPage, pageSize]);

    // 首次加载和页码变化时获取数据
    useEffect(() => {
        fetchUsers();
    }, [fetchUsers]);

    // 表格列配置
    const columns = ['username', 'email', 'role', 'status'];
    const columnNames = {
        username: '用户名',
        email: '邮箱',
        role: '角色',
        status: '状态'
    };

    // 行数据转换
    const renderRow = (user) => ({
        ...user,
        role: getRoleName(user.role),
        status: user.status === 1 ? '启用' : '禁用'
    });

    // 角色名称转换
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

    // 自定义操作按钮
    const customActions = [
        {
            icon: 'lock',
            color: 'orange',
            onClick: handleToggleStatus
        },
        {
            icon: 'edit',
            color: 'blue',
            onClick: handleEdit
        }
    ];

    // 切换用户状态
    async function handleToggleStatus(user) {
        try {
            // TODO: 实现用户状态切换的 API 调用
            await userAPI.manageUser(user.username, user.status === 1 ? 'disable' : 'enable');
            toast.success(`${user.username} 状态已更新`);
            await fetchUsers(currentPage);
        } catch (error) {
            toast.error(error.message || '更新用户状态失败');
        }
    }

    // 编辑用户
    function handleEdit(user) {
        // TODO: 实现用户编辑功能
        toast.info(`编辑用户: ${user.username}`);
    }

    // 分页变化处理
    const handlePageChange = (newPage) => {
        setCurrentPage(newPage);
        fetchUsers(newPage);
    };

    return (
        <div className="p-4">
            <Dimmer active={loading} inverted>
                <Loader>加载中...</Loader>
            </Dimmer>

            <DataList
                title="用户列表"
                data={users}
                columns={columns}
                columnNames={columnNames}
                renderRow={renderRow}
                customActions={customActions}
                currentPage={currentPage}
                pageSize={pageSize}
                totalItems={totalItems}
                onPageChange={handlePageChange}
            />
        </div>
    );
};

export default UserList;

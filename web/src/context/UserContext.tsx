import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { userAPI } from '../api/userAPI';
import { toast } from 'react-toastify';
import { useAuth } from './AuthContext';
import { getRoleName } from '../utils/roleUtils';

interface UserDTO {
  id: number;
  username: string;
  display_name: string;
  role: number;
  role_name?: string; // 添加角色名称字段便于显示
  status: number;
  email: string;
  created_at?: string;
}

interface AuthContextType {
  isAuthenticated: boolean;
  token?: string;
}

interface UserContextType {
  // 状态
  users: UserDTO[];
  loading: boolean;
  
  // 方法
  fetchUsers: () => Promise<void>;
  addUser: (userData: Partial<UserDTO>) => Promise<void>;
  deleteUser: (userId: number) => Promise<void>;
  updateUserRole: (userId: number, role: number) => Promise<void>;
  toggleUserStatus: (userId: number, action: 'enable' | 'disable') => Promise<void>;
}

// 创建上下文
const UserContext = createContext<UserContextType | undefined>(undefined);

// 使用UserContext的自定义Hook
export const useUsers = (): UserContextType => {
  const context = useContext(UserContext);
  if (!context) {
    throw new Error('useUsers must be used within a UserProvider');
  }
  return context;
};

// 获取带类型的Auth上下文
export const useAuthWithType = () => useAuth() as AuthContextType;

// 用户上下文提供者组件
export const UserProvider = ({ children }: { children?: ReactNode }): ReactNode => {
  const { isAuthenticated } = useAuthWithType();
  const [users, setUsers] = useState<UserDTO[]>([]);
  const [loading, setLoading] = useState<boolean>(false);

  // 处理用户数据，添加角色名称
  const processUserData = (userData: UserDTO[]): UserDTO[] => {
    return userData.map(user => ({
      ...user,
      role_name: getRoleName(user.role)
    }));
  };

  // 获取用户列表
  const fetchUsers = async (): Promise<void> => {
    setLoading(true);
    try {
      const result = await userAPI.getUsers();
      if (Array.isArray(result.data)) {
        setUsers(processUserData(result.data));
      } else {
        console.error('[用户上下文] 获取用户失败: 返回数据格式不正确', result);
        setUsers([]);
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '获取用户列表失败';
      toast.error(errorMessage);
      console.error('[用户上下文] 获取用户失败:', err);
      setUsers([]);
    } finally {
      setLoading(false);
    }
  };

  // 添加用户
  const addUser = async (userData: Partial<UserDTO>): Promise<void> => {
    setLoading(true);
    try {
      const newUser = await userAPI.createUser(userData);
      const processedUser = { ...newUser, role_name: getRoleName(newUser.role) };
      setUsers((prevUsers) => [...prevUsers, processedUser]);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '创建用户失败';
      toast.error(errorMessage);
      console.error('[用户上下文] 创建用户失败:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  // 删除用户
  const deleteUser = async (userId: number): Promise<void> => {
    setLoading(true);
    try {
      await userAPI.deleteUser(userId);
      setUsers((prevUsers) => prevUsers.filter((user) => user.id !== userId));
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '删除用户失败';
      toast.error(errorMessage);
      console.error('[用户上下文] 删除用户失败:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  // 更新用户角色
  const updateUserRole = async (userId: number, role: number): Promise<void> => {
    setLoading(true);
    try {
      const updatedUser = await userAPI.updateUserRole(userId, role);
      const role_name = getRoleName(role);
      setUsers((prevUsers) =>
        prevUsers.map((user) => (user.id === userId ? { ...user, role, role_name } : user))
      );
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '更新用户角色失败';
      toast.error(errorMessage);
      console.error('[用户上下文] 更新用户角色失败:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  // 切换用户状态（启用/禁用）
  const toggleUserStatus = async (userId: number, action: 'enable' | 'disable'): Promise<void> => {
    setLoading(true);
    try {
      const updatedUser = await userAPI.toggleUserStatus(userId, action);
      setUsers((prevUsers) =>
        prevUsers.map((user) => (user.id === userId ? { ...user, status: updatedUser.status } : user))
      );
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : `${action === 'enable' ? '启用' : '禁用'}用户失败`;
      toast.error(errorMessage);
      console.error(`[用户上下文] ${action}用户失败:`, err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  // 当认证状态变化时获取用户列表
  useEffect(() => {
    if (isAuthenticated) {
      fetchUsers();
    }
  }, [isAuthenticated]);

  // 提供的上下文值
  const value: UserContextType = {
    users,
    loading,
    fetchUsers,
    addUser,
    deleteUser,
    updateUserRole,
    toggleUserStatus
  };

  return (
    <UserContext.Provider value={value}>
      {children}
    </UserContext.Provider>
  );
};

export default UserProvider; 
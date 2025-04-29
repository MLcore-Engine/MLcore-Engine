import { API, handleResponse, handlePagedResponse, apiService } from '../helpers/api';

// 基于dto.go定义的UserDTO接口
export interface UserDTO {
  id: number;
  username: string;
  display_name: string;
  role: number;
  status: number;
  email: string;
  created_at?: string; // 对应dto.go中的time.Time
}

// 基于dto.go定义的UsersListData
export interface UsersListData {
  users: UserDTO[];
  total: number;
  page: number;
  limit: number;
}

// 基于dto.go中的PagedData
export interface PagedData {
  total: number;
  page: number;
  limit: number;
}

// 登录请求参数
export interface LoginRequest {
  username: string;
  password: string;
}

// 登录响应中的TokenDTO
export interface TokenDTO {
  token: string;
  user: UserDTO;
}

// API基础路径
const BASE_URL = '/api/user';

/**
 * 用户管理相关API
 */
export const userAPI = {
  /**
   * 获取用户列表
   * @param page 页码（从0开始）
   * @param limit 每页数量
   * @returns 用户列表数据
   */
  getUsers: async (page = 0, limit = 10) => {
    try {
      const response = await API.get(`${BASE_URL}/manage/`, {
        params: { page, limit }
      });
      
      // 使用handleResponse处理响应，提取data字段
      const result = handleResponse<any>(response);
      
      // 适配后端返回的数据结构
      return {
        data: result.users || [],
        total: result.total || 0,
        page: result.page || page,
        limit: result.limit || limit
      };
    } catch (error) {
      console.error('[用户API] 获取用户列表失败:', error);
      throw error;
    }
  },

  /**
   * 搜索用户
   * @param keyword 搜索关键词
   * @param page 页码
   * @param limit 每页数量
   * @returns 搜索结果
   */
  searchUsers: async (keyword = '', page = 0, limit = 10) => {
    try {
      const response = await API.get(`${BASE_URL}/manage/search/`, {
        params: { keyword, page, limit }
      });
      
      const result = handleResponse<any>(response);
      
      return {
        data: result.users || [],
        total: result.total || 0,
        page: result.page || page,
        limit: result.limit || limit
      };
    } catch (error) {
      console.error('[用户API] 搜索用户失败:', error);
      throw error;
    }
  },

  /**
   * 获取用户详情
   * @param userId 用户ID
   * @returns 用户详情
   */
  getUserDetail: async (userId: number) => {
    try {
      const response = await API.get(`${BASE_URL}/manage/${userId}/`);
      return handleResponse<UserDTO>(response);
    } catch (error) {
      console.error(`[用户API] 获取用户详情失败:`, error);
      throw error;
    }
  },

  /**
   * 创建用户
   * @param userData 用户数据
   * @returns 创建的用户
   */
  createUser: async (userData: Partial<UserDTO>) => {
    try {
      const response = await API.post(`${BASE_URL}/manage/`, userData);
      return handleResponse<UserDTO>(response);
    } catch (error) {
      console.error('[用户API] 创建用户失败:', error);
      throw error;
    }
  },

  /**
   * 更新用户
   * @param userId 用户ID
   * @param userData 用户数据
   * @returns 更新后的用户
   */
  updateUser: async (userId: number, userData: Partial<UserDTO>) => {
    try {
      // 构建完整的用户对象用于更新
      const updateData = {
        ID: userId, // 使用大写ID匹配后端格式
        ...userData
      };
      const response = await API.put(`${BASE_URL}/manage/`, updateData);
      return handleResponse<UserDTO>(response);
    } catch (error) {
      console.error(`[用户API] 更新用户失败:`, error);
      throw error;
    }
  },

  /**
   * 更新用户角色
   * @param userId 用户ID
   * @param role 角色ID
   * @returns 更新后的用户
   */
  updateUserRole: async (userId: number, role: number) => {
    try {
      return await userAPI.updateUser(userId, { role });
    } catch (error) {
      console.error(`[用户API] 更新用户角色失败:`, error);
      throw error;
    }
  },

  /**
   * 删除用户
   * @param userId 用户ID
   * @returns 操作结果
   */
  deleteUser: async (userId: number) => {
    try {
      const response = await API.delete(`${BASE_URL}/manage/${userId}/`);
      return handleResponse<boolean>(response);
    } catch (error) {
      console.error(`[用户API] 删除用户失败:`, error);
      throw error;
    }
  },
  
  /**
   * 用户登录
   * @param loginData 登录数据
   * @returns 登录结果与token
   */
  login: async (loginData: LoginRequest) => {
    try {
      const response = await API.post(`${BASE_URL}/login/`, loginData);
      return handleResponse<TokenDTO>(response);
    } catch (error) {
      console.error('[用户API] 登录失败:', error);
      throw error;
    }
  },
  
  /**
   * 启用/禁用用户
   * @param userId 用户ID
   * @param action 操作类型 ('enable'或'disable')
   * @returns 操作结果
   */
  toggleUserStatus: async (userId: number, action: 'enable' | 'disable') => {
    try {
      const payload = {
        username: userId.toString(), // 后端可能使用username字段来接收ID
        action: action
      };
      const response = await API.post(`${BASE_URL}/manage/manage/`, payload);
      return handleResponse<UserDTO>(response);
    } catch (error) {
      console.error(`[用户API] ${action}用户失败:`, error);
      throw error;
    }
  }
}; 
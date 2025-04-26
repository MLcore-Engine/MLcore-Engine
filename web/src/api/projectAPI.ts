import { API, handleResponse } from '../helpers/api';

// API基础路径
const BASE_URL = '/api/project';
const MEMBERS_URL = '/api/project-members';

// 定义项目相关类型 - 与后端保持一致
export interface MemberDTO {
  projectId: string;           // 主键ID，字符串类型
  userId: string;       // 用户ID，字符串类型
  username: string;     // 用户名
  role: number;         // 角色代码
}

export interface Project {
  id: string;            // 项目ID，字符串类型
  name: string;          // 项目名称
  description: string;   // 项目描述
  createdAt?: string;    // 创建时间
  updatedAt?: string;    // 更新时间
  users?: MemberDTO[];   // 项目成员列表
}
/**
 * 项目相关API
 */
export const projectAPI = {
  /***************************** 项目基础API *****************************/
  
  /**
   * 获取项目列表
   * @param page 页码
   * @param limit 每页数量
   * @returns 项目列表
   */
  getProjects: async (page = 1, limit = 10): Promise<{ projects: Project[], total: number }> => {
    try {
      const response = await API.get(`${BASE_URL}/get-all`, {
        params: { page, limit }
      });
      const result = handleResponse<any>(response);
      // console.log('[项目API] 获取项目列表成功:', result);
      return {
        projects: result.projects || [],
        total: result.total || 0
      };
    } catch (error) {
      console.error('[项目API] 获取项目列表失败:', error);
      throw error;
    }
  },

  /**
   * 获取项目详情
   * @param projectId 项目ID
   * @returns 项目详情
   */
  getProject: async (projectId: string): Promise<Project> => {
    try {
      const response = await API.get(`${BASE_URL}/${projectId}/`);
      const result = handleResponse<any>(response);
      return result.data;
    } catch (error) {
      console.error('[项目API] 获取项目详情失败:', error);
      throw error;
    }
  },

  /**
   * 创建项目
   * @param projectData 项目数据
   * @returns 创建的项目
   */
  createProject: async (projectData: Partial<Project>): Promise<Project> => {
    try {
      const response = await API.post(`${BASE_URL}/`, projectData);
      return handleResponse<Project>(response);
    } catch (error) {
      console.error('[项目API] 创建项目失败:', error);
      throw error;
    }
  },

  /**
   * 更新项目
   * @param projectId 项目ID
   * @param projectData 项目数据
   * @returns 更新后的项目
   */
  updateProject: async (projectId: string, projectData: Partial<Project>): Promise<Project> => {
    try {
      const updateData = {
        ...projectData,
        id: projectId
      };
      const response = await API.put(`${BASE_URL}/${projectId}`, updateData);
      return handleResponse<Project>(response);
    } catch (error) {
      console.error('[项目API] 更新项目失败:', error);
      throw error;
    }
  },

  /**
   * 删除项目
   * @param projectId 项目ID
   * @returns 操作结果
   */
  deleteProject: async (projectId: string): Promise<boolean> => {
    try {
      const response = await API.delete(`${BASE_URL}/${projectId}`);
      const result = handleResponse<any>(response);
      return result.data;  
    } catch (error) {
      console.error('[项目API] 删除项目失败:', error);
      throw error;
    }
  },

  /***************************** 项目成员API *****************************/

  /**
   * 获取用户的项目列表
   * @param userId 用户ID
   * @returns 项目列表
   */
  getUserProjects: async (userId: string): Promise<Project[]> => {
    try {
      const response = await API.get(`${MEMBERS_URL}/user/${userId}/`);
      const result = handleResponse<any>(response);
      return result.data || [];
    } catch (error) {
      console.error('[项目API] 获取用户项目列表失败:', error);
      throw error;
    }
  },

  /**
   * 获取项目成员
   * @param projectId 项目ID
   * @returns 成员列表
   */
  getProjectMembers: async (projectId: string): Promise<MemberDTO[]> => {
    try {
      const response = await API.get(`${MEMBERS_URL}/project/${projectId}/`);
      const result = handleResponse<any>(response);
      return result.data || [];
    } catch (error) {
      console.error('[项目API] 获取项目成员失败:', error);
      throw error;
    }
  },

  /**
   * 添加用户到项目
   * @param projectId 项目ID
   * @param userId 用户ID
   * @param role 角色ID
   * @returns 添加的成员信息
   */
  addProjectMember: async (projectId: string, userId: string, role: number): Promise<MemberDTO> => {
    try {
      const payload = { 
        projectId: parseInt(projectId, 10),
        userId: parseInt(userId, 10),
        role 
      };
      const response = await API.post(`${MEMBERS_URL}/`, payload);
      // handleResponse 已返回后端 data 对象
      const raw = handleResponse<any>(response);
      // 映射为前端 MemberDTO
      const memberDTO: MemberDTO = {
        projectId: String(raw.projectId),
        userId:   String(raw.userId),
        username: raw.username,
        role:     raw.role,
      };
      return memberDTO;
    } catch (error) {
      console.error('[项目API] 添加项目成员失败:', error);
      throw error;
    }
  },

  /**
   * 从项目移除用户
   * @param projectId 项目ID
   * @param userId 用户ID
   * @returns 操作结果
   */
  removeProjectMember: async (projectId: string, userId: string): Promise<boolean> => {
    try {
      if (!projectId || !userId) {
        throw new Error('项目ID和用户ID不能为空');
      }
      
      console.log('[项目API] 移除成员请求:', { projectId, userId });
      
      // 转换为数字ID
      const projectIdNum = parseInt(projectId, 10);
      const userIdNum = parseInt(userId, 10);
      
      const response = await API.delete(`${MEMBERS_URL}/${projectIdNum}/${userIdNum}`);
      handleResponse<any>(response); // 检查响应，但不使用返回值
      return true; // 如果没有异常，则认为操作成功
    } catch (error) {
      console.error('[项目API] 移除项目成员失败:', error);
      throw error;
    }
  },

  /**
   * 更新用户在项目中的角色
   * @param projectId 项目ID
   * @param userId 用户ID或用户名
   * @param role 角色ID
   * @returns 更新后的成员信息
   */
  updateProjectMemberRole: async (projectId: string, userId: string, role: number): Promise<MemberDTO> => {
    try {

      const projectIdNum = parseInt(projectId, 10);
      const userIdNum = parseInt(userId, 10);
      
      const payload = { 
        projectId: projectIdNum,
        userId: userIdNum,
        role 
      };
      
      const response = await API.put(`${MEMBERS_URL}/`, payload);
      const result = handleResponse<any>(response);
      return result
    } catch (error) {
      console.error('[项目API] 更新项目成员角色失败:', error);
      throw error;
    }
  }
};
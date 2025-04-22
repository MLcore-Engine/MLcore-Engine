import { API, handleResponse } from '../helpers/api';

const BASE_URL = '/api/project';
const MEMBERS_URL = '/api/project-memberships';

// 定义项目相关类型
export interface MemberDTO {
  ID?: string;
  userID: string;
  username: string;
  role: number;
  [key: string]: any;
}

export interface Project {
  ID: string;
  name: string;
  description: string;
  created_at?: string;
  updated_at?: string;
  users?: MemberDTO[];
  [key: string]: any;
}

// 后端响应结构
interface ProjectsResponse {
  projects: Project[];
  PagedData?: {
    Total: number;
    Page: number;
    Limit: number;
  }
}

export const projectAPI = {


  /***************************** Project related API *****************************/

  //get all projects, todo: add pagination
  getProjects: async (): Promise<Project[]> => {
    try {
      const response = await API.get(`${BASE_URL}/get-all`);
      // 确保处理正确的数据结构 - 从response.data.projects中获取数组
      const responseData = handleResponse<ProjectsResponse>(response);
      return responseData.projects || [];
    } catch (error) {
      console.error('API Error in getProjects:', error);
      throw error; // 错误已在API拦截器中处理，直接抛出
    }
  },

  //create a project
  createProject: async (project: Partial<Project>): Promise<Project> => {
    try {
      const response = await API.post(BASE_URL, project);
      return handleResponse(response);
    } catch (error) {
      console.error('API Error in createProject:', error);
      throw error; // 错误已在API拦截器中处理，直接抛出
    }
  },

  //update a project
  updateProject: async (projectId: string, projectData: Partial<Project>): Promise<Project> => {
    try {
      const response = await API.put(`${BASE_URL}/${projectId}`, projectData);
      return handleResponse(response);
    } catch (error) {
      console.error('API Error in updateProject:', error);
      throw error;
    }
  },

  //delete a project
  deleteProject: async (projectId: string): Promise<void> => {
    try {
      await API.delete(`${BASE_URL}/${projectId}`);
    } catch (error) {
      console.error('API Error in deleteProject:', error);
      throw error;
    }
  },


  //get a project
  getProject: async (projectId: string): Promise<Project> => {
    try {
      const response = await API.get(`${BASE_URL}/${projectId}`);
      return handleResponse(response);
    } catch (error) {
      console.error('API Error in getProject:', error);
      throw error;
    }
  },

  /***************************** Project Members API *****************************/

  //get user's project list
  getUserProjects: async (userId: string): Promise<Project[]> => {
    try {
      const response = await API.get(`${MEMBERS_URL}/user/${userId}`);
      return handleResponse(response);
    } catch (error) {
      console.error('API Error in getUserProjects:', error);
      throw error;
    }
  },

  //get members of a project
  getProjectMembers: async (projectId: string): Promise<MemberDTO[]> => {
    try {
      const response = await API.get(`${MEMBERS_URL}/project/${projectId}`);
      return handleResponse(response);
    } catch (error) {
      console.error('API Error in getProjectMembers:', error);
      throw error;
    }
  },

  //add member to a project
  addProjectMember: async (projectId: string, member: { userID: string, role: number }): Promise<MemberDTO> => {
    try {
      const payload = { projectID: projectId, userID: member.userID, role: member.role };
      const response = await API.post(`${MEMBERS_URL}`, payload);
      return handleResponse(response);
    } catch (error) {
      console.error('API Error in addProjectMember:', error);
      throw error;
    }
  },

  //remove member from a project
  removeProjectMember: async (projectId: string, userId: string): Promise<void> => {
    try {
      await API.delete(`${MEMBERS_URL}/${projectId}/${userId}`);
    } catch (error) {
      console.error('API Error in removeProjectMember:', error);
      throw error;
    }
  },

  //update user's role in a project
  updateUserProjectRole: async (projectId: string, userId: string, role: number): Promise<MemberDTO> => {
    try {
      const payload = { role };   
      const response = await API.put(`${MEMBERS_URL}/${projectId}/${userId}`, payload);
      return handleResponse(response);
    } catch (error) {
      console.error('API Error in updateUserProjectRole:', error);
      throw error;
    }
  },


};
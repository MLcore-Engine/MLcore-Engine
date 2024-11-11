import { API } from '../helpers/api';

const BASE_URL = '/api/project';
const MEMBERS_URL = '/api/project-memberships';

export const projectAPI = {


  /***************************** Project related API *****************************/

  //get all projects, todo: add pagination
  getProjects: async () => {
    try {
      const response = await API.get(`${BASE_URL}/get-all`);
      console.log('API Response in getProjects:', response);
      return response.data.data;
    } catch (error) {
      console.error('API Error in getProjects:', error);
      throw new Error('Failed to fetch projects: ' + error);
    }
  },

  //create a project
  createProject: async (project) => {
    try {
      const response = await API.post(BASE_URL, project);
      return response.data.data;
    } catch (error) {
      console.error('API Error in createProject:', error);
      if (error.response && error.response.data && error.response.data.error) {
        throw new Error(error.response.data.error);
      }
      throw new Error('Failed to create project');
    }
  },

  //update a project
  updateProject: async (projectId, projectData) => {
    try {
      const response = await API.put(`${BASE_URL}/${projectId}`, projectData);
      return response.data.data;
    } catch (error) {
      console.error('API Error in updateProject:', error);
      throw new Error('Failed to update project');
    }
  },

  //delete a project
  deleteProject: async (projectId) => {
    try {
      await API.delete(`${BASE_URL}/${projectId}`);
    } catch (error) {
      console.error('API Error in deleteProject:', error);
      throw new Error('Failed to delete project');
    }
  },


  //get a project
  getProject: async (projectId) => {
    try {
      const response = await API.get(`${BASE_URL}/${projectId}`);
      console.log('API Response in getProject:', response);
      return response.data.data;
    } catch (error) {
      console.error('API Error in getProject:', error);
      throw new Error('Failed to fetch project');
    }
  },

  /***************************** Project Members API *****************************/

  //get user's project list
  getUserProjects: async (userId) => {
    try {
      const response = await API.get(`${MEMBERS_URL}/user/${userId}`);
      return response.data.data;
    } catch (error) {
      console.error('API Error in getUserProjects:', error);
      throw new Error('Failed to fetch user projects');
    }
  },

  //get members of a project
  getProjectMembers: async (projectId) => {
    try {
      const response = await API.get(`${MEMBERS_URL}/project/${projectId}`);
      return response.data.data;
    } catch (error) {
      console.error('API Error in getProjectMembers:', error);
      throw new Error('Failed to fetch project members');
    }
  },

  //add member to a project
  addProjectMember: async (projectId, member) => {
    try {
      const payload = { projectID: projectId, userID: member.userID, role: member.role };
      console.log('Sending payload:', payload);
      const response = await API.post(`${MEMBERS_URL}`, payload);
      return response.data.data;
    } catch (error) {
      console.error('API Error in addProjectMember:', error);
      throw new Error('Failed to add project member');
    }
  },

  //remove member from a project
  removeProjectMember: async (projectId, userId) => {
    try {
      await API.delete(`${MEMBERS_URL}/${projectId}/${userId}`);
    } catch (error) {
      console.error('API Error in removeProjectMember:', error);
      throw new Error('Failed to remove project member');
    }
  },

  //update user's role in a project
  updateUserProjectRole: async (projectId, userId, role) => {
    try {
      const payload = { role };   
      const response = await API.put(`${MEMBERS_URL}/${projectId}/${userId}`, payload);
      return response.data.data;
    } catch (error) {
      console.error('API Error in updateUserProjectRole:', error);
      throw new Error('Failed to update user role in project');
    }
  },


};
import { API } from '../helpers/api';

const BASE_URL = '/api/triton';

export const tritonAPI = {
  /**
   * Create a new Triton deployment
   */
  createTritonDeploy: async (deployData) => {
    try {
      const response = await API.post(BASE_URL, deployData);
      if (!response.data.success) {
        throw new Error(response.data.message || 'Failed to create Triton deployment');
      }
      return response.data.data;
    } catch (error) {
      console.error('API Error in createTritonDeploy:', error);
      throw new Error(error.response?.data?.message || 'Failed to create Triton deployment');
    }
  },

  /**
   * Get Triton deployments with pagination
   */
  getTritonDeploys: async (page = 1, limit = 10) => {
    try {
      const response = await API.get(`${BASE_URL}/get-all`, {
        params: { page, limit }
      });
      if (!response.data.success) {
        throw new Error(response.data.message || 'Failed to fetch Triton deployments');
      }
      return {
        data: response.data.data.deployments,
        total: response.data.data.total,
        page: response.data.data.page,
        limit: response.data.data.limit
      };
    } catch (error) {
      console.error('API Error in getTritonDeploys:', error);
      throw new Error(error.response?.data?.message || 'Failed to fetch Triton deployments');
    }
  },

  /**
   * Delete a Triton deployment
   */
  deleteTritonDeploy: async (id) => {
    try {
      const response = await API.delete(`${BASE_URL}/${id}`);
      if (!response.data.success) {
        throw new Error(response.data.message || 'Failed to delete Triton deployment');
      }
      return response.data.message;
    } catch (error) {
      console.error('API Error in deleteTritonDeploy:', error);
      throw new Error(error.response?.data?.message || 'Failed to delete Triton deployment');
    }
  },

  updateTritonDeploy: async (id, deployData) => {
    try {
      const response = await API.put(`${BASE_URL}/${id}`, deployData);
      if (!response.data.success) {
        throw new Error(response.data.message || 'Failed to update Triton deployment');
      }
      return response.data.data;
    } catch (error) {
      console.error('API Error in updateTritonDeploy:', error);
      throw new Error(error.response?.data?.message || 'Failed to update Triton deployment');
    }
  }
};


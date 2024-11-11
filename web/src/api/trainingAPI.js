import { API } from '../helpers/api';

const BASE_URL = '/api/pytorchtrain';

export const trainingAPI = {
  /**
   * Create a new training job
   */
  createTrainingJob: async (trainingData) => {
    try {
      const response = await API.post(BASE_URL, trainingData);

      if (!response.data.success) {
        throw new Error(response.data.message || 'Failed to create training job');
      }

      return response.data.data;
    } catch (error) {
      console.error('API Error in createTrainingJob:', error);
      throw new Error(error.response?.data?.message || 'Failed to create training job');
    }
  },

  /**
   * Get training jobs with pagination
   */
  getTrainingJobs: async (page = 1, limit = 10) => {
    try {
      const response = await API.get(`${BASE_URL}/get-all`, {
        params: { page, limit }
      });

      if (!response.data.success) {
        throw new Error(response.data.message || 'Failed to fetch training jobs');
      }

      return {
        data: response.data.data.training_jobs,
        total: response.data.data.total,
        page: response.data.data.page,
        limit: response.data.data.limit
      };
    } catch (error) {
      console.error('API Error in getTrainingJobs:', error);
      throw new Error(error.response?.data?.message || 'Failed to fetch training jobs');
    }
  },

  /**
   * Delete a training job
   */
  deleteTrainingJob: async (id) => {
    try {
      const response = await API.delete(`${BASE_URL}/${id}`);

      if (!response.data.success) {
        throw new Error(response.data.message || 'Failed to delete training job');
      }

      return response.data.message;
    } catch (error) {
      console.error('API Error in deleteTrainingJob:', error);
      throw new Error(error.response?.data?.message || 'Failed to delete training job');
    }
  }
};
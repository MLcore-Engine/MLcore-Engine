// src/api/notebookAPI.js

import { API } from '../helpers/api';

const BASE_URL = '/api/notebook';

export const notebookAPI = {
  /**
   * Fetch notebooks with pagination
   * @param {number} page - Page number
   * @param {number} limit - Items per page
   */
  getNotebooks: async (page = 1, limit = 10) => {
    try {
      const response = await API.get(`${BASE_URL}/get-all`, {
        params: {
          page,
          limit,
        },
      });

      if (!response.data.success) {
        throw new Error(response.data.message || 'Failed to fetch notebooks');
      }

      return {
        data: response.data.data.notebooks,
        total: response.data.data.total,
        page: response.data.data.page,
        limit: response.data.data.limit,
        success: response.data.success,
        message: response.data.message,
      };
    } catch (error) {
      console.error('API Error in getNotebooks:', error);
      throw new Error(error.response?.data?.message || 'Failed to fetch notebooks');
    }
  },

  /**
   * Create a new notebook
   * @param {Object} notebookData - Data for the new notebook
   */
  createNotebook: async (notebookData) => {
    try {
      const response = await API.post(`${BASE_URL}/`, notebookData);

      if (!response.data.success) {
        throw new Error(response.data.message || 'Failed to create notebook');
      }

      return response.data.notebook;
    } catch (error) {
      console.error('API Error in createNotebook:', error);
      throw new Error(error.response?.data?.message || 'Failed to create notebook');
    }
  },

  /**
   * Reset a notebook by ID
   * @param {number} id - Notebook ID
   */
  resetNotebook: async (id) => {
    try {
      const response = await API.get(`${BASE_URL}/reset/${id}`);

      if (!response.data.success) {
        throw new Error(response.data.message || 'Failed to reset notebook');
      }

      return response.data.message;
    } catch (error) {
      console.error('API Error in resetNotebook:', error);
      throw new Error(error.response?.data?.message || 'Failed to reset notebook');
    }
  },

  /**
   * Delete a notebook by ID
   * @param {number} id - Notebook ID
   */
  deleteNotebook: async (id) => {
    try {
      const response = await API.delete(`${BASE_URL}/${id}`);

      if (!response.data.success) {
        throw new Error(response.data.message || 'Failed to delete notebook');
      }

      return response.data.message;
    } catch (error) {
      console.error('API Error in deleteNotebook:', error);
      throw new Error(error.response?.data?.message || 'Failed to delete notebook');
    }
  },
};
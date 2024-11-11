import { API } from '../helpers/api';

// 修正基础路由路径，与后端对应
const BASE_URL = '/api/user/manage';

export const userAPI = {

  /**
   * 获取用户列表
   * @param {number} page - 页码，从1开始
   * @param {number} limit - 每页条数
   * @returns {Promise<Array>} 用户列表数据
   */
  getUsers: async (page = 1, limit = 10) => {
    try {
      const response = await API.get(`${BASE_URL}/`, {
        params: {
          page,
          limit
        }
      });

      if (!response.data.success) {
        throw new Error(response.data.message || 'Failed to fetch users');
      }
      console.log('API Response in getUsers:', response);
      return {
        data: response.data.data.users,
        total: response.data.data.total,
        page: response.data.data.page,
        limit: response.data.data.limit,
        success: response.data.success,
        message: response.data.message
        };
    } catch (error) {
      console.error('API Error in getUsers:', error);
      throw new Error(error.response?.data?.message || 'Failed to fetch users');
    }
  },

  /**
   * 搜索用户
   * @param {string} keyword - 搜索关键词
   * @param {number} page - 页码，从1开始
   * @param {number} limit - 每页条数
   * @returns {Promise<Object>} 搜索结果
   */
  searchUsers: async (keyword = '', page = 1, limit = 10) => {
    try {
      const response = await API.get(`${BASE_URL}/search`, {
        params: {
          keyword,
          page,
          limit
        }
      });

      if (!response.data.success) {
        throw new Error(response.data.message || 'Failed to search users');
      }

      return {
        data: response.data.data,
        success: response.data.success,
        message: response.data.message
      };
    } catch (error) {
      console.error('API Error in searchUsers:', error);
      throw new Error(error.response?.data?.message || 'Failed to search users');
    }
  }
};
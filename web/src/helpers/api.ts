import { showError } from './utils';
import axios, { AxiosInstance, AxiosRequestConfig, InternalAxiosRequestConfig, AxiosResponse, AxiosError } from 'axios';
import { useAuth } from '../context/AuthContext'; // 假设有AuthContext

// 定义API响应类型
interface APIResponse<T = any> {
  success: boolean;
  message?: string;
  data: T;
}

// 分页数据结构
interface PagedData {
  Total: number;
  Page: number;
  Limit: number;
}

// 分页响应结构
interface PagedResponse<T = any> {
  [key: string]: T[] | PagedData | any;
  PagedData?: PagedData;
}

// 处理后的分页数据
interface ProcessedPagedData<T = any> {
  list: T[];
  total: number;
  page: number;
  limit: number;
}

// 创建API实例
export const API: AxiosInstance = axios.create({
  baseURL: process.env.REACT_APP_SERVER ? process.env.REACT_APP_SERVER : '',
});

// 创建一个高阶函数配置API
export const configureAPI = (getToken: () => string | null): void => {
  API.interceptors.request.use(
    (config: any): any => {
      const token = getToken();
      if (token) {
        config.headers = config.headers || {};
        config.headers['Authorization'] = 'Bearer ' + token;
      }
      return config;
    },
    (error: AxiosError): Promise<never> => {
      console.error('API Error:', error);
      return Promise.reject(error);
    }
  );
};

// 默认配置 - 使用localStorage
configureAPI(() => localStorage.getItem('token'));

// 提供使用Context的配置函数
export const useAPIWithAuth = (): AxiosInstance => {
  const { token } = useAuth();
  configureAPI(() => token);
  return API;
};

API.interceptors.response.use(
  (response: AxiosResponse): AxiosResponse => response,
  (error: AxiosError): Promise<never> => {
    showError(error);
    return Promise.reject(error);
  }
);

// 标准响应处理函数
export const handleResponse = <T = any>(response: AxiosResponse<APIResponse<T>>): T => {
  if (response.data.success) {
    return response.data.data;
  } else {
    throw new Error(response.data.message || 'API请求失败');
  }
};

// 分页响应处理函数
export const handlePagedResponse = <T = any>(
  response: AxiosResponse<APIResponse<PagedResponse<T>>>
): ProcessedPagedData<T> => {
  if (response.data.success) {
    const { data } = response.data;
    
    // 兼容不同的分页数据结构
    return {
      list: (data.list || data.datasets || data.training_jobs || data.deploys || 
             data.entries || data.notebooks || []) as T[],
      total: data.total || (data.PagedData ? data.PagedData.Total : 0),
      page: data.page || (data.PagedData ? data.PagedData.Page : 1),
      limit: data.limit || (data.PagedData ? data.PagedData.Limit : 10)
    };
  } else {
    throw new Error(response.data.message || 'API请求失败');
  }
};
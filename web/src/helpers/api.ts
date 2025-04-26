import { showError } from './utils';
import axios, { 
  AxiosInstance, 
  AxiosRequestConfig, 
  InternalAxiosRequestConfig, 
  AxiosResponse, 
  AxiosError, 
  AxiosHeaders 
} from 'axios';
import { useAuth } from '../context/AuthContext';

// ========================= API响应类型定义 =========================


/**
 * APIResponse<T> 是通用的API响应数据结构接口，适用于后端返回的数据格式统一的场景。
 * - success: 表示请求是否成功（布尔值）。
 * - message: 可选，返回的提示信息或错误信息（字符串）。
 * - data: 泛型T，实际返回的数据内容，可以是对象、数组等任意类型。
 * 
 * 该接口通常对应后端的BaseResponse结构，方便前端类型推断和数据处理。
 */
export interface APIResponse<T = any> {
  success: boolean;
  message?: string;
  data: T;
}

/**
 * 分页数据
 * 对应dto.go中的PagedData
 */
export interface PagedData {
  total: number;  // 对应dto.go中的Total
  page: number;   // 对应dto.go中的Page
  limit: number;  // 对应dto.go中的Limit
}

/**
 * 处理后的标准分页数据结构
 * 统一前端使用的分页数据格式
 */
export interface StandardListResult<T> {
  items: T[];         // 数据项
  total: number;      // 总数
  page: number;       // 当前页
  limit: number;      // 每页限制
  hasMore: boolean;   // 是否有更多数据
}

// ========================= API 实例配置 =========================

/**
 * 创建并配置全局API实例
 */
export const API: AxiosInstance = axios.create({
  baseURL: process.env.REACT_APP_API_BASE_URL || '',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  }
});

/**
 * 配置API实例的认证令牌获取函数
 * @param getToken 获取认证令牌的函数
 */
export const configureAuth = (getToken: () => string | null): void => {
  API.interceptors.request.use(
    (config: InternalAxiosRequestConfig): InternalAxiosRequestConfig => {
      const token = getToken();
      if (token) {
        if (!config.headers) {
          config.headers = new AxiosHeaders();
        }
        config.headers.set('Authorization', `Bearer ${token}`);
      }
      return config;
    },
    (error: AxiosError): Promise<never> => {
      console.error('[API] 请求配置错误:', error);
      return Promise.reject(error);
    }
  );
};

// 默认使用localStorage中的令牌
configureAuth(() => localStorage.getItem('token'));

// 提供使用Context的配置函数
export const useAPIWithAuth = (): AxiosInstance => {
  const { token } = useAuth();
  configureAuth(() => token);
  return API;
};

// 全局响应拦截器
API.interceptors.response.use(
  (response: AxiosResponse): AxiosResponse => response,
  (error: AxiosError): Promise<never> => {
    // 处理特定HTTP错误
    if (error.response) {
      const status = error.response.status;
      
      // 处理未授权错误
      if (status === 401) {
        console.warn('[API] 未授权请求，可能需要重新登录');
      }
      
      // 处理服务器错误
      if (status >= 500) {
        console.error('[API] 服务器错误:', error.response.data);
      }
    } else if (error.request) {
      // 请求已发送但未收到响应
      console.error('[API] 网络错误:', error.message);
    } else {
      // 其他错误
      console.error('[API] 请求配置错误:', error.message);
    }
    
    // 显示错误提示
    showError(error);
    return Promise.reject(error);
  }
);

// ========================= 响应处理函数 =========================

/**
 * 处理API响应
 * @param response axios响应对象
 * @returns 响应数据
 * @throws 当success为false时抛出包含错误消息的Error
 */
export const handleResponse = <T = any>(response: AxiosResponse<APIResponse<T>>): T => {
  // 检查响应类型，确保是JSON
  const contentType = response.headers['content-type'];
  if (contentType && contentType.includes('application/json')) {
    const { data } = response;
    
    if (!data) {
      throw new Error('后端返回数据为空');
    }
    
    if (data.success) {
      return data.data;
    } else {
      console.error('[API] 后端API请求失败:', {
        url: response.config.url,
        method: response.config.method,
        status: response.status,
        data: data
      });
      throw new Error(data.message || '后端API请求失败');
    }
  } else {
    console.error('[API] 后端返回了非JSON响应:', {
      url: response.config.url,
      status: response.status,
      contentType
    });
    throw new Error('后端API返回了非预期格式的响应');
  }
};

/**
 * 处理分页列表API响应
 * @param response API响应对象
 * @param itemsKey 可选的数据项字段名，用于非标准响应
 * @returns 标准化的分页列表结果
 * @throws 当success为false时抛出包含错误消息的Error
 */
export const handleListResponse = <T = any>(
  response: AxiosResponse<APIResponse<any>>,
  itemsKey?: string
): StandardListResult<T> => {
  const responseData = handleResponse<any>(response);
  
  // 提取数据项数组
  let items: T[] = [];
  
  // 判断数据字段
  const possibleListKeys = [
    'users', 'projects', 'datasets', 'entries', 
    'notebooks', 'training_jobs', 'deploys', 'items'
  ];
  
  // 尝试从已知字段获取数据
  for (const key of possibleListKeys) {
    if (key in responseData && Array.isArray(responseData[key])) {
      items = responseData[key];
      break;
    }
  }
  
  // 如果提供了特定的key且存在该字段
  if (!items.length && itemsKey && itemsKey in responseData) {
    items = responseData[itemsKey] as T[];
  }
  
  // 提取分页信息
  const total = responseData.total || 0;
  const page = responseData.page || 1;
  const limit = responseData.limit || 10;

  return {
    items,
    total,
    page,
    limit,
    hasMore: (page * limit) < total
  };
};

// ========================= 简化API请求函数 =========================

/**
 * API封装类，提供简化的API请求方法
 */
export const apiService = {
  get: <T = any>(url: string, params?: any, config?: AxiosRequestConfig): Promise<T> => {
    return API.get(url, { ...config, params })
      .then(response => handleResponse<T>(response));
  },
  
  getList: <T = any>(
    url: string, 
    params?: any, 
    itemsKey?: string, 
    config?: AxiosRequestConfig
  ): Promise<StandardListResult<T>> => {
    return API.get(url, { ...config, params })
      .then(response => handleListResponse<T>(response, itemsKey));
  },
  
  post: <T = any, D = any>(url: string, data?: D, config?: AxiosRequestConfig): Promise<T> => {
    return API.post(url, data, config)
      .then(response => handleResponse<T>(response));
  },
  
  put: <T = any, D = any>(url: string, data?: D, config?: AxiosRequestConfig): Promise<T> => {
    return API.put(url, data, config)
      .then(response => handleResponse<T>(response));
  },
  
  delete: <T = any>(url: string, config?: AxiosRequestConfig): Promise<T> => {
    return API.delete(url, config)
      .then(response => handleResponse<T>(response));
  },
  
  patch: <T = any, D = any>(url: string, data?: D, config?: AxiosRequestConfig): Promise<T> => {
    return API.patch(url, data, config)
      .then(response => handleResponse<T>(response));
  }
};

// 向后兼容，导出原有的命名
export const handlePagedResponse = handleListResponse;

// 添加一个统一的请求URL处理函数，确保URL格式正确
export const normalizeUrl = (url: string): string => {
  return url; // 不自动添加尾部斜杠，保持URL原样
};

// 在API对象中使用normalizeUrl函数处理所有请求URL
const originalGet = API.get;
API.get = (url: string, config?: any) => {
  return originalGet(normalizeUrl(url), config);
};

const originalPost = API.post;
API.post = (url: string, data?: any, config?: any) => {
  return originalPost(normalizeUrl(url), data, config);
};

const originalPut = API.put;
API.put = (url: string, data?: any, config?: any) => {
  return originalPut(normalizeUrl(url), data, config);
};

const originalDelete = API.delete;
API.delete = (url: string, config?: any) => {
  return originalDelete(normalizeUrl(url), config);
};
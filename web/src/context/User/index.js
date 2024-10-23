import React, { createContext, useContext, useReducer, useEffect } from 'react';
import { reducer, initialState } from './reducer';
import jwtDecode from 'jwt-decode';   

export const UserContext = createContext();

export const UserProvider = ({ children }) => {
  const [state, dispatch] = useReducer(reducer, initialState);

  useEffect(() => {
    const token = localStorage.getItem('token');
    const userString = localStorage.getItem('user');
    const user = userString ? JSON.parse(userString) : null;
    if (token) {
      // 验证 token 有效性
      // TODO: 添加 token 有效性验证
      // isTokenValid(token)
      if (true) {
        dispatch({ type: 'setToken', payload: { token, user } });
      } else {
        localStorage.removeItem('token');
      }
    }
  }, []);

  const value = {
    ...state,
    login: (userData, token) => {
      localStorage.setItem('token', token);
      dispatch({ type: 'login', payload: { user: userData, token } });
    },
    logout: () => {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      dispatch({ type: 'logout' });
    },
  };

  return <UserContext.Provider value={value}>{children}</UserContext.Provider>;
};

export const useUser = () => useContext(UserContext);

// 辅助函数：验证 token
function isTokenValid(token) {
  
  if(!token) return false;

  try {
    const decodedToken = jwtDecode(token);
    const currentTime = Date.now() / 1000; // 转换为秒

    // 检查 token 是否过期
    if (!decodedToken.exp || decodedToken.exp < currentTime) {
      return false;
    }

    // 检查 token 的颁发时间是否在合理范围内
    // 例如，我们可以检查 token 是否是在过去 30 天内颁发的
    const thirtyDaysAgo = currentTime - 30 * 24 * 60 * 60; // 30天前的时间戳（秒）
    if (!decodedToken.iat || decodedToken.iat < thirtyDaysAgo) {
      return false;
    }
    return true;
  } catch (error) {
    console.error('Token validation error:', error);
    return false;
  }
}
import React, { createContext, useContext, useReducer, useEffect } from 'react';
import jwtDecode from 'jwt-decode';
import { toast } from 'react-toastify';

// 初始状态
const initialState = {
  user: null,
  token: null,
  isAuthenticated: false,
};

// 创建 AuthContext
const AuthContext = createContext();

// Reducer 函数
const reducer = (state, action) => {
  switch (action.type) {
    case 'login':
      return {
        ...state,
        user: action.payload.user,
        token: action.payload.token,
        isAuthenticated: true,
      };
    case 'logout':
      return {
        ...state,
        user: null,
        token: null,
        isAuthenticated: false,
      };
    case 'setToken':
      return {
        ...state,
        token: action.payload.token,
        isAuthenticated: true,
        user: action.payload.user,
      };
    default:
      return state;
  }
};

// 辅助函数：验证 token
const isTokenValid = (token) => {
    if (!token) return false;
    try {
        const decodedtoken = jwtDecode(token);
        const currentTime = Date.now() / 1000;
        return decodedtoken.exp && decodedtoken.exp > currentTime;
    }catch(error) {
        console.log('token valdation err: ',  error)
        return false;
    }

}

// AuthProvider 组件
export const AuthProvider = ({ children }) => {
  const [state, dispatch] = useReducer(reducer, initialState);

  useEffect(() => {
    const token = localStorage.getItem('token');
    const userString = localStorage.getItem('user');
    const user = userString ? JSON.parse(userString) : null;

    if (token && isTokenValid(token)) {
      dispatch({ type: 'setToken', payload: { token, user } });
    } else {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
    }
  }, []);

  const login = (userData, token) => {
    localStorage.setItem('token', token);
    localStorage.setItem('user', JSON.stringify(userData));
    dispatch({ type: 'login', payload: { user: userData, token } });
    toast.success('login success！');
  };

  const logout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    dispatch({ type: 'logout' });
    toast.info('logout success！');  
  };

  const value = {
    ...state,
    login,
    logout,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

// 自定义 Hook 以便在组件中使用 AuthContext
export const useAuth = () => useContext(AuthContext);
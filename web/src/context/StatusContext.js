// web/src/context/Status/index.js

import React from 'react';

// Reducer 函数
export const reducer = (state, action) => {
  switch (action.type) {
    case 'set':
      return {
        ...state,
        status: action.payload,
      };
    case 'unset':
      return {
        ...state,
        status: undefined,
      };
    default:
      return state;
  }
};

// 初始状态
export const initialState = {
  status: undefined,
};

// 创建上下文
export const StatusContext = React.createContext({
  state: initialState,
  dispatch: () => null,
});

// 上下文提供者组件
export const StatusProvider = ({ children }) => {
  const [state, dispatch] = React.useReducer(reducer, initialState);

  return (
    <StatusContext.Provider value={[state, dispatch]}>
      {children}
    </StatusContext.Provider>
  );
};
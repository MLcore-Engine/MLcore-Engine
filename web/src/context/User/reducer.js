

export const initialState = {
  user: null,
  token: null,
  isAuthenticated: false,
};

export const reducer = (state, action) => {
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

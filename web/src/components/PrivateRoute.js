import React from 'react';
import { useUser } from '../context/User';
import { Navigate, useLocation } from 'react-router-dom';

const PrivateRoute = ({ children }) => {
  const { isAuthenticated } = useUser();
  const location = useLocation();

  if (!isAuthenticated) {
    return <Navigate to='/login' state={{ from: location }} replace/>;
  }

  return children;
}

export default PrivateRoute;


/// <reference types="react-scripts" />


declare module 'react-toastify';
declare module '@remix-run/router' {
  interface Router {
    basename: string;
    // ...其他属性
  }
}

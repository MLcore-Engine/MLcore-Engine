export {};

// 全局声明React类型
declare global {
  namespace React {
    interface ReactNode {
      children?: React.ReactNode;
    }
  }
}
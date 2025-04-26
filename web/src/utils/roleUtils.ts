/**
 * 角色相关常量和工具函数
 */

// 角色常量定义
export const ROLE = {
  COMMON: 1,    // 普通用户
  ADMIN: 10,    // 管理员
  ROOT: 100     // 超级管理员/拥有者
};

/**
 * 获取用户角色名称
 * @param role 角色值
 * @returns 角色名称
 */
export const getRoleName = (role: number): string => {
  switch (role) {
    case ROLE.COMMON:
      return '普通用户';
    case ROLE.ADMIN:
      return '管理员';
    case ROLE.ROOT:
      return '超级管理员';
    default:
      return '未知';
  }
};

/**
 * 获取项目中用户角色名称
 * @param role 角色值
 * @returns 项目角色名称
 */
export const getProjectRoleName = (role: number): string => {
  switch (role) {
    case ROLE.COMMON:
      return '用户';
    case ROLE.ADMIN:
      return '管理员';
    case ROLE.ROOT:
      return '拥有者';
    default:
      return '未知';
  }
};

/**
 * 角色选项列表
 */
export const roleOptions = [
  { key: '1', text: '普通用户', value: ROLE.COMMON },
  { key: '10', text: '管理员', value: ROLE.ADMIN },
  { key: '100', text: '超级管理员', value: ROLE.ROOT },
];

/**
 * 项目角色选项列表
 */
export const projectRoleOptions = [
  { key: '1', text: '用户', value: ROLE.COMMON },
  { key: '10', text: '管理员', value: ROLE.ADMIN },
  { key: '100', text: '拥有者', value: ROLE.ROOT },
]; 
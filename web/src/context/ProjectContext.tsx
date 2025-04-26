import * as React from 'react';
import { createContext, useContext, useState, useEffect } from 'react';
import { useAuth } from './AuthContext';
import { projectAPI, Project} from '../api/projectAPI';
import { toast } from 'react-toastify';
import type { ReactNode } from 'react';
interface AuthContextType {
  isAuthenticated: boolean;
  token?: string;
}

/**
 * 项目上下文接口
 * 定义了项目管理所需的所有方法和状态
 */
interface ProjectContextType {
  // 状态
  projects: Project[];
  loading: boolean;
  
  // 项目操作方法
  fetchProjects: () => Promise<void>;
  createProject: (projectData: Partial<Project>) => Promise<void>;
  updateProject: (projectId: string, projectData: Partial<Project>) => Promise<void>;
  deleteProject: (projectId: string) => Promise<void>;
  
  // 成员操作方法
  addProjectMember: (projectId: string, member: { userId: string; role: number }) => Promise<void>;
  removeProjectMember: (projectId: string, userId: string) => Promise<void>;
  updateProjectMemberRole: (projectId: string, userId: string, role: number) => Promise<void>;
}

// 创建带默认值的Context，避免null检查
const defaultContextValue: ProjectContextType = {
  projects: [],
  loading: false,
  fetchProjects: async () => {},
  addProjectMember: async () => {},
  removeProjectMember: async () => {},
  updateProjectMemberRole: async () => {},
  createProject: async () => {},
  updateProject: async () => {},
  deleteProject: async () => {},
};

// 创建项目上下文，使用默认值初始化
// 使用类型断言确保类型安全，避免运行时错误
const ProjectContext = createContext<ProjectContextType>(defaultContextValue as ProjectContextType);

/**
 * 使用项目上下文的自定义Hook
 * @returns 项目上下文
 */
export const useProjects = () => useContext(ProjectContext);

/**
 * 类型断言useAuth的返回类型
 * @returns 认证上下文
 */
export const useAuthWithType = () => useAuth() as AuthContextType;


/**
 * 项目上下文提供者组件
 */
export const ProjectProvider = ({ children }: { children?: ReactNode }): ReactNode => {
  const { isAuthenticated } = useAuthWithType();
  const [projects, setProjects] = useState([] as Project[]);
  const [loading, setLoading] = useState(false);

  /**
   * 获取所有项目
   */
  const fetchProjects = async (): Promise<void> => {
    if (!isAuthenticated) return;

    setLoading(true);
    try {
      // 1. 调用 API，拿到原始列表
      const raw = await projectAPI.getProjects(/* 可传 page,limit */) 
      // raw.projects: any[]，这里暂时当成 RawProject[]

      // 2. 逐条映射
      const formatted: Project[] = raw.projects.map((p: any) => ({
        id:          String(p.id),                 // 数字 -> 字符串
        name:        p.name,                       // 原样
        description: p.description,                // 原样
        createdAt:   p.created_at,                 // snake -> camel
        updatedAt:   p.updated_at,                 // 同上
        users: Array.isArray(p.users)             // 可能没 users
          ? p.users.map((u: any) => ({
              id:        String(u.userId),         // MemberDTO.id
              projectId: String(u.projectId),      // projectId
              userId:    String(u.userId),         // userId
              username:  u.username,               // 原样
              role:      u.role                    // 原样数字
            }))
          : []
      }));

      // 3. 赋值到 Context 状态
      setProjects(formatted);

    } catch (err) {
      const msg = err instanceof Error ? err.message : '获取项目列表失败';
      toast.error(msg);
      console.error('[项目上下文] fetchProjects error:', err);

      // 失败时清空列表
      setProjects([]);
    } finally {
      setLoading(false);
    }
  };

  /**
   * 添加项目成员
   * @param projectId 项目ID
   * @param member 成员信息
   */
  const addProjectMember = async (projectId: string, member: { userId: string; role: number }): Promise<void> => {
    setLoading(true);
    try {
      // 转换参数名称以符合API期望
      // const apiMember = {
      //   userId: member.userId,
      //   role: member.role
      // };
      
      const newMember = await projectAPI.addProjectMember(projectId, member.userId, member.role);
      
      // 更新本地状态
      setProjects((prevProjects) =>
        prevProjects.map((project) =>
          project.id === projectId
            ? { ...project, users: [...(project.users || []), newMember] }
            : project
        )
      );
      
      toast.success('成员添加成功');
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '添加成员失败';
      toast.error(errorMessage);
      console.error('[项目上下文] 添加成员失败:', err);
    } finally {
      setLoading(false);
    }
  };

  /**
   * 移除项目成员
   * @param projectId 项目ID
   * @param userId 成员ID
   */
  const removeProjectMember = async (projectId: string, userId: string): Promise<void> => {
    if (!userId || !projectId) {
      console.error('[项目上下文] 移除成员失败: 缺少必要参数', { projectId, userId });
      toast.error('移除成员失败: 缺少必要参数');
      return;
    }
    
    setLoading(true);
    try {
      console.log('[项目上下文] 准备移除成员:', { projectId, userId });
      await projectAPI.removeProjectMember(projectId, userId);
      
      // 更新本地状态，使用userId字段匹配
      setProjects((prevProjects) =>
        prevProjects.map((project) =>
          project.id === projectId
            ? {
                ...project,
                users: project.users?.filter((user) => String(user.userId) !== userId),
              }
            : project
        )
      );
      
      toast.success('成员已移除');
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '移除成员失败';
      toast.info(errorMessage);
      // console.error('[项目上下文] 移除成员失败:', err);
    } finally {
      setLoading(false);
    }
  };

  /**
   * 更新项目成员角色
   * @param projectId 项目ID
   * @param userId 成员ID
   * @param role 角色ID
   */
  const updateProjectMemberRole = async (projectId: string, userId: string, role: number): Promise<void> => {
    setLoading(true);
    try {
      // 使用正确的API方法名
      const updatedMember = await projectAPI.updateProjectMemberRole(projectId, userId, role);  
      
      // 更新本地状态，使用userId字段匹配
      setProjects((prevProjects) =>
        prevProjects.map((project) =>
          project.id === projectId
            ? {
                ...project,
                users: project.users?.map((user) =>
                  String(user.userId) === userId
                    ? { ...user, role: updatedMember.role }
                    : user
                ),
              }
            : project
        )
      );
      toast.success('成员角色已更新');
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '更新成员角色失败';
      toast.error(errorMessage);
      console.error('[项目上下文] 更新成员角色失败:', err);
    } finally {
      setLoading(false);
    }
  };

  /**
   * 创建项目
   * @param projectData 项目数据
   */
  const createProject = async (projectData: Partial<Project>): Promise<void> => {
    setLoading(true);
    try {
      const newProject = await projectAPI.createProject(projectData);
      setProjects((prevProjects) => [...prevProjects, newProject]);
      
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '创建项目失败';
      toast.error(errorMessage);
      console.error('[项目上下文] 创建项目失败:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  /**
   * 更新项目
   * @param projectId 项目ID
   * @param projectData 项目数据
   */
  const updateProject = async (projectId: string, projectData: Partial<Project>): Promise<void> => {
    setLoading(true);
    try {
      const updatedProject = await projectAPI.updateProject(projectId, projectData);
      setProjects((prevProjects) =>
        prevProjects.map((project) =>
          project.id === projectId ? updatedProject : project
        )
      );
      toast.success('项目已更新');
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '更新项目失败';
      toast.error(errorMessage);
      console.error('[项目上下文] 更新项目失败:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  /**
   * 删除项目
   * @param projectId 项目ID
   */
  const deleteProject = async (projectId: string): Promise<void> => {
    setLoading(true);
    try {
      await projectAPI.deleteProject(projectId);
      setProjects((prevProjects) =>
        prevProjects.filter((project) => project.id !== projectId)
      );
      toast.success('项目已删除');
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : '删除项目失败';
      toast.error(errorMessage);
      console.error('[项目上下文] 删除项目失败:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  // 在组件挂载且认证状态变化时获取项目列表
  useEffect(() => {
    if (isAuthenticated) {
      fetchProjects();
    }
  }, [isAuthenticated]);

  // 提供的上下文值
  const value: ProjectContextType = {
    projects,
    loading,
    fetchProjects,
    addProjectMember,
    removeProjectMember,
    updateProjectMemberRole,
    createProject,
    updateProject,
    deleteProject,
  };

  return (  
    <ProjectContext.Provider value={value}>
      {children}
    </ProjectContext.Provider>
  );
};
import React, { createContext, useContext, useState, useEffect } from 'react';
import { useAuth } from './AuthContext';
import { projectAPI, Project, MemberDTO } from '../api/projectAPI';
import { toast } from 'react-toastify';

interface AuthContextType {
  isAuthenticated: boolean;
}

interface ProjectContextType {
  projects: Project[];
  loading: boolean;
  fetchProjects: () => Promise<void>;
  addProjectMember: (projectId: string, member: { userID: string; role: number }) => Promise<void>;
  removeProjectMember: (projectId: string, memberId: string) => Promise<void>;
  updateProjectMemberRole: (projectId: string, memberId: string, role: number) => Promise<void>;
  createProject: (projectData: Partial<Project>) => Promise<void>;
  updateProject: (projectId: string, projectData: Partial<Project>) => Promise<void>;
  deleteProject: (projectId: string) => Promise<void>;
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

const ProjectContext = createContext(defaultContextValue);

// Create a custom Hook to use ProjectContext in components
export const useProjects = () => {
  return useContext(ProjectContext);
};

// 类型断言useAuth的返回类型
export const useAuthWithType = () => useAuth() as AuthContextType;

// ProjectProvider component
export const ProjectProvider = ({ children }: { children: any }) => {
  const { isAuthenticated } = useAuthWithType();
  // 添加明确的类型
  const [projects, setProjects] = useState([] as Project[]);
  const [loading, setLoading] = useState(false);

  /***************************** Project Members API *****************************/
  // Fetch all projects
  const fetchProjects = async () => { 
    
    if (!isAuthenticated) {
      return;
    }
    setLoading(true);
    try {
      // projectAPI.getProjects()现在直接返回项目数组
      const projectsData = await projectAPI.getProjects();
      // 确保projectsData是数组
      if (Array.isArray(projectsData)) {
        setProjects(projectsData);
      } else {
        console.error('Expected array but got:', typeof projectsData);
        setProjects([]);
      }
    } catch (err: any) {
      if (err.response && err.response.status === 401) {
        console.log('Unauthorized. Please login again.');
      } else {
        toast.error(err.message || 'Failed to fetch projects.');
        console.error('API Error in fetchProjects:', err);
      }
      // 确保在出错时设置为空数组
      setProjects([]);
    } finally {
      setLoading(false);
    }
  };

  // 添加项目成员
  
  const addProjectMember = async (projectId: string, member: { userID: string; role: number }) => {
    setLoading(true);
    try {
      const newMember = await projectAPI.addProjectMember(projectId, member);
      setProjects((prevProjects: Project[]) =>
        prevProjects.map((project: Project) =>
          project.ID === projectId
            ? { ...project, users: [...(project.users || []), newMember] }
            : project
        )
      );
      toast.info('Member added successfully.');
    } catch (err: any) {
      console.error("API Error in addProjectMember: ", err);
      toast.error(err.message || 'Failed to add member.');
    } finally {
      setLoading(false);
    }
  };


  // remove project member
  const removeProjectMember = async (projectId: string, memberId: string) => {
    setLoading(true);
    try {
      await projectAPI.removeProjectMember(projectId, memberId);
      setProjects((prevProjects: Project[]) =>
        prevProjects.map((project: Project) =>
          project.ID === projectId
            ? {
                ...project,
                users: project.users?.filter((user: MemberDTO) => user.userID !== memberId),
              }
            : project
        )
      );
      toast.info('Member removed successfully.');
    } catch (err: any) {
      console.error("API Error in removeProjectMember: ", err);
      toast.error(err.message || 'Failed to remove member.');
    } finally {
      setLoading(false);
    }
  };

  const updateProjectMemberRole = async (projectId: string, memberId: string, role: number) => {
    setLoading(true);
    try {
      const updatedMember = await projectAPI.updateUserProjectRole(projectId, memberId, role);
      setProjects((prevProjects: Project[]) =>
        prevProjects.map((project: Project) =>
          project.ID === projectId
            ? {
                ...project,
                users: project.users?.map((user: MemberDTO) =>
                  user.userID === memberId ? { ...user, role: updatedMember.role } : user
                ),
              }
            : project
        )
      );
      toast.success('Member role updated successfully.');
    } catch (err: any) {
      console.error("API Error in updateProjectMemberRole: ", err);
      toast.error(err.message || 'Failed to update member role.');
    } finally {
      setLoading(false);
    }
  };


  /***************************** Project API *****************************/
  // create project
  const createProject = async (projectData: Partial<Project>) => {
    setLoading(true);
    try {
      const newProject = await projectAPI.createProject(projectData);
      setProjects((prevProjects: Project[]) => [...prevProjects, newProject]);
    } catch (err: any) {
      console.error("API Error in createProject: ", err);
      throw err; 
    } finally {
      setLoading(false);
    }
  };

  // update an existingproject
  const updateProject = async (projectId: string, projectData: Partial<Project>) => {
    setLoading(true);
    try {
      const updatedProject = await projectAPI.updateProject(projectId, projectData);
      setProjects((prevProjects: Project[]) =>
        prevProjects.map((project: Project) =>
          project.ID === projectId ? updatedProject : project
        )
      );
    } catch (err: any) {
      console.error('Error in updateProject:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  // delete a project
  const deleteProject = async (projectId: string) => {
    setLoading(true);
    try {
      await projectAPI.deleteProject(projectId);
      setProjects((prevProjects: Project[]) =>
        prevProjects.filter((project: Project) => project.ID !== projectId)
      );
    } catch (err: any) {
      console.error("API Error in deleteProject: ", err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  // 在组件挂载时获取项目列表
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
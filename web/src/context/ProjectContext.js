import React, { createContext, useContext, useState, useEffect } from 'react';
import { projectAPI } from '../api/projectAPI';
import { toast } from 'react-toastify';

// Create ProjectContext
const ProjectContext = createContext(null);

// Create a custom Hook to use ProjectContext in components
export const useProjects = () => {
  return useContext(ProjectContext);
};

// ProjectProvider component
export const ProjectProvider = ({ children }) => {
  
  // projects model state with users
  const [projects, setProjects] = useState([]);
  const [loading, setLoading] = useState(false);

  /***************************** Project Members API *****************************/
  // Fetch all projects
  const fetchProjects = async () => { 
    setLoading(true);
    try {
      const data = await projectAPI.getProjects();
      setProjects(data.projects);
    } catch (err) {
      toast.error(err.message || 'Failed to fetch projects.');
      console.error('API Error in fetchProjects:', err);
    } finally {
      setLoading(false);
    }
  };

  // // 获取用户的项目关系列表
  // const getUserProjects = async (userId) => {
  //   setLoading(true);
  //   try {
  //     const data = await projectAPI.getUserProjects(userId);
  //     setProjectMembers(data);
  //   } catch (err) {
  //     toast.error(err.message || 'Failed to fetch user projects');
  //     console.error('API Error in getUserProjects:', err);
  //   } finally {
  //     setLoading(false);
  //   }
  // };

  // // Fetch project members
  // const fetchProjectMembers = async (projectId) => { 
  //   setLoading(true);
  //   try {
  //     const members = await projectAPI.getProjectMembers(projectId);
  //     setProjectMembers(members);
  //   } catch (err) {
  //     setError(err.message);
  //     console.error('API Error in fetchProjectMembers:', err);
  //   } finally {
  //     setLoading(false);
  //   }
  // };

  // 添加项目成员
  
  const addProjectMember = async (projectId, member) => {
    setLoading(true);
    try {
      const newMember = await projectAPI.addProjectMember(projectId, member);
      setProjects((prevProjects) =>
        prevProjects.map((project) =>
          project.ID === projectId
            ? { ...project, users: [...project.users, newMember] }
            : project
        )
      );
      toast.info('Member added successfully.');
    } catch (err) {
      console.error("API Error in addProjectMember: ", err);
      toast.error(err.message || 'Failed to add member.');
      // throw err;
    } finally {
      setLoading(false);
    }
  };


  // remove project member
  const removeProjectMember = async (projectId, memberId) => {
    setLoading(true);
    // setError(null);
    try {
      await projectAPI.removeProjectMember(projectId, memberId);
      setProjects((prevProjects) =>
        prevProjects.map((project) =>
          project.ID === projectId
            ? {
                ...project,
                users: project.users.filter((user) => user.ID !== memberId),
              }
            : project
        )
      );
      toast.info('Member removed successfully.');
    } catch (err) {
      console.error("API Error in removeProjectMember: ", err);
      toast.error(err.message || 'Failed to remove member.');
      // throw err;
    } finally {
      setLoading(false);
    }
  };

  const updateProjectMemberRole = async (projectId, memberId, role) => {
    setLoading(true);
    // setError(null);
    try {
      const updatedMember = await projectAPI.updateUserProjectRole(projectId, memberId, role);
      setProjects((prevProjects) =>
        prevProjects.map((project) =>
          project.ID === projectId
            ? {
                ...project,
                users: project.users.map((user) =>
                  user.ID === memberId ? updatedMember : user
                ),
              }
            : project
        )
      );
      toast.success('Member role updated successfully.');
    } catch (err) {
      console.error("API Error in updateProjectMemberRole: ", err);
      toast.error(err.message || 'Failed to update member role.');
      // throw err;
    } finally {
      setLoading(false);
    }
  };


  /***************************** Project API *****************************/
  // create project
  const createProject = async (projectData) => {
    setLoading(true);
    // setError(null);
    try {
      const newProject = await projectAPI.createProject(projectData);
      setProjects((prevProjects) => [...prevProjects, newProject]);
    } catch (err) {
      console.error("API Error in createProject: ", err);
      // toast.error(err.message || 'Failed to create project.');
      throw err; 
    } finally {
      setLoading(false);
    }
  };

  // update an existingproject
  const updateProject = async (projectId, projectData) => {
    setLoading(true);
    // setError(null);
    try {
      const updatedProject = await projectAPI.updateProject(projectId, projectData);
      setProjects((prevProjects) =>
        prevProjects.map((project) =>
          project.ID === projectId ? updatedProject : project
        )
      );
    } catch (err) {
      // setError(err.message);  
      console.error('Error in createProject:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  // delete a project
  const deleteProject = async (projectId) => {
    setLoading(true);
    try {
      await projectAPI.deleteProject(projectId);
      setProjects((prevProjects) =>
        prevProjects.filter((project) => project.ID !== projectId)
      );
    } catch (err) {
      console.error("API Error in deleteProject: ", err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  // 在组件挂载时获取项目列表
  useEffect(() => {
    fetchProjects();
  }, []);

  // 提供的上下文值
  const value = {
    projects,
    loading,
    // error,
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
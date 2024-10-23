import ProjectManagementSidebar from '../components/sidebars/ProjectManagementSidebar';
import  ModelDevelopmentSidebar from '../components/sidebars/ModelDevelopmentSidebar';
import ImageManagementSidebar from '../components/sidebars/ImageManagementSidebar';
import ModelTrainingSidebar from '../components/sidebars/ModelTrainingSidebar';
import ModelDeploymentSidebar from '../components/sidebars/ModelDeploymentSidebar';
import DashboardSidebar from '../components/sidebars/DashboardSidebar';


const sidebarMapping = {
  '/dashboard': DashboardSidebar,
  '/project': ProjectManagementSidebar,
  '/project/create': ProjectManagementSidebar,
  '/notebook': ModelDevelopmentSidebar,
  '/notebook/create': ModelDevelopmentSidebar,
  '/image': ImageManagementSidebar,
  '/image/create': ImageManagementSidebar,
  '/train': ModelTrainingSidebar,
  '/train/create': ModelTrainingSidebar,
  '/deploy': ModelDeploymentSidebar,
  '/deploy/create': ModelDeploymentSidebar,
};

export default sidebarMapping;


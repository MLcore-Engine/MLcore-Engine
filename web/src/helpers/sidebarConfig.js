import ProjectManagementSidebar from '../components/sidebars/ProjectManagementSidebar';
import  ModelDevelopSidebar from '../components/sidebars/ModelDevelopSidebar';
import ImageManagementSidebar from '../components/sidebars/ImageManageSidebar';
import ModelTrainingSidebar from '../components/sidebars/ModelTrainingSidebar';
import ModelDeploymentSidebar from '../components/sidebars/ModelDeploySidebar';
import DashboardSidebar from '../components/sidebars/DashboardSidebar';


const sidebarMapping = {
  '/dashboard': DashboardSidebar,
  '/project': ProjectManagementSidebar,
  '/project/create': ProjectManagementSidebar,
  '/notebook': ModelDevelopSidebar,
  '/notebook/create': ModelDevelopSidebar,
  '/image': ImageManagementSidebar,
  '/image/create': ImageManagementSidebar,
  '/train': ModelTrainingSidebar,
  '/train/create': ModelTrainingSidebar,
  '/deploy': ModelDeploymentSidebar,
  '/deploy/create': ModelDeploymentSidebar,
};

export default sidebarMapping;


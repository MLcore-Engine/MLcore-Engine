import React, { useState, useEffect } from 'react';
import { Header, Segment, Pagination } from 'semantic-ui-react';
import { tritonAPI } from '../../api/tritonAPI';
import { toast } from 'react-toastify';
import DataList from '../../components/common/DataList';
import CreateServingForm from './CreateServingForm';
import ServingDetails from './ServingDetails';

const ServingList = () => {
  const [deploys, setDeploys] = useState([]);
  const [loading, setLoading] = useState(false);
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [isDetailsOpen, setIsDetailsOpen] = useState(false);
  const [selectedDeploy, setSelectedDeploy] = useState(null);
  const [pageInfo, setPageInfo] = useState({ page: 1, limit: 10, total: 0 });

  const fetchDeploys = async (page = 1, limit = 10) => {
    setLoading(true);
    try {
      const response = await tritonAPI.getTritonDeploys(page, limit);
      setDeploys(response.data);
      setPageInfo({
        page: response.page,
        limit: response.limit,
        total: response.total
      });
    } catch (error) {
      toast.error(error.message || 'Failed to fetch Triton deployments');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchDeploys();
  }, []);

  const handleDelete = async (item) => {
    try {
      await tritonAPI.deleteTritonDeploy(item.id);
      toast.success('Triton deployment deleted successfully');
      fetchDeploys(pageInfo.page, pageInfo.limit);
    } catch (error) {
      toast.error(error.message || 'Failed to delete model service');
    }
  };

  const handleDetails = (deploy) => {
    setSelectedDeploy(deploy);
    setIsDetailsOpen(true);
  };

  const handlePageChange = (e, { activePage }) => {
    fetchDeploys(activePage, pageInfo.limit);
  };

  const columns = ['name', 'namespace', 'status', 'access_url', 'resources'];
  const columnNames = {
    name: 'Name',
    namespace: 'Namespace',
    status: 'Status',
    access_url: 'Access URL',
    resources: 'Resources'
  };

  const customActions = [
    {
      icon: 'trash',
      color: 'red',
      onClick: handleDelete,
      confirm: {
        content: 'Are you sure you want to delete this model deployment?',
        confirmButton: 'Delete',
        cancelButton: 'Cancel',
      }
    },
    {
      icon: 'info circle',
      color: 'blue',
      onClick: handleDetails
    }
  ];

  return (
    <div className="p-4">
      <Header as="h2">Serving-List</Header>
      <Segment>
        <DataList
          title="Serving-List"
          data={deploys}
          columns={columns}
          columnNames={columnNames}
          onAdd={() => setIsCreateOpen(true)}
          customActions={customActions}
          renderRow={(deploy) => ({
            name: deploy.name,
            namespace: deploy.namespace,
            status: deploy.status,
            image: deploy.image,
            resources: `CPU: ${deploy.cpu}, Memory: ${deploy.memory}, GPU: ${deploy.gpu}`
          })}
        />

        <Pagination
          activePage={pageInfo.page}
          totalPages={Math.ceil(pageInfo.total / pageInfo.limit)}
          onPageChange={handlePageChange}
          floated='right'
          style={{ marginTop: '1em' }}
        />

        {isCreateOpen && (
          <CreateServingForm
            isOpen={isCreateOpen}
            onClose={() => setIsCreateOpen(false)}
            onRefresh={() => fetchDeploys(pageInfo.page, pageInfo.limit)}
          />
        )}

        {isDetailsOpen && selectedDeploy && (
          <ServingDetails  
            deploy={selectedDeploy}
            isOpen={isDetailsOpen}
            onClose={() => setIsDetailsOpen(false)}
          />
        )}
      </Segment>
    </div>
  );
};

export default ServingList;
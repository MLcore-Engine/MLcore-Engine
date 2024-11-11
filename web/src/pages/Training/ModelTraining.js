// src/pages/training/ModelTraining.js

import React, { useState, useEffect } from 'react';
import { Header as SemanticHeader, Segment, Button, Icon, Pagination } from 'semantic-ui-react';
import { trainingAPI } from '../../api/trainingAPI';
import { toast } from 'react-toastify';
import DataList from '../../components/common/DataList';
import CreateTrainingForm from './CreateTrainingForm';
import TrainingJobDetails from './TrainingJobDetails'; // We'll create this component next

const ModelTraining = () => {
  const [jobs, setJobs] = useState([]);
  const [loading, setLoading] = useState(false);
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [isDetailsOpen, setIsDetailsOpen] = useState(false);
  const [selectedJob, setSelectedJob] = useState(null);
  const [pageInfo, setPageInfo] = useState({ page: 1, limit: 10, total: 0 });

  const fetchJobs = async (page = 1, limit = 10) => {
    setLoading(true);
    try {
      const response = await trainingAPI.getTrainingJobs(page, limit);
      setJobs(response.data);
      setPageInfo({ page: response.data.page, limit: response.data.limit, total: response.data.total });
    } catch (error) {
      toast.error(error.message || 'Failed to fetch training jobs');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchJobs();
  }, []);

  const handleDelete = async (item) => {
    try {
      await trainingAPI.deleteTrainingJob(item.id);
      toast.success('Training job deleted successfully');
      fetchJobs(pageInfo.page, pageInfo.limit);
    } catch (error) {
      toast.error(error.message || 'Failed to delete training job');
    }
  };

  const handleDetails = (job) => {
    setSelectedJob(job);
    setIsDetailsOpen(true);
  };

  const handlePageChange = (e, { activePage }) => {
    fetchJobs(activePage, pageInfo.limit);
  };

  const customActions = [
    {
      icon: 'trash',
      color: 'red',
      onClick: handleDelete,
      confirm: {
        content: 'Are you sure you want to delete this training job?',
        confirmButton: 'Delete',
        cancelButton: 'Cancel',
      }
    },
    {
      icon: 'info circle',
      color: 'blue',
      onClick: handleDetails,
      tooltip: 'View Details',
      show: () => true, // Always show the details button
    },
  ];

  const columns = ['name', 'status', 'image', 'resources'];
  const columnNames = {
    name: 'Name',
    status: 'Status',
    image: 'Image',
    resources: 'Resources',
  };

  return (
    <div className="p-4">
      <SemanticHeader as="h2">模型训练</SemanticHeader>
      <Segment>
        
        <DataList
          title="TrainingJob List"
          data={jobs}
          columns={columns}
          columnNames={columnNames}
          onAdd={() => setIsCreateOpen(true)}
          customActions={customActions}
          renderRow={(job) => ({
            name: job.name,
            status: job.status,
            image: job.image,
            resources: `CPU: ${job.cpu_limit}, Memory: ${job.memory_limit}`
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
          <CreateTrainingForm
            isOpen={isCreateOpen}
            onClose={() => setIsCreateOpen(false)}
            onRefresh={() => fetchJobs(pageInfo.page, pageInfo.limit)}
          />
        )}

        {isDetailsOpen && selectedJob && (
          <TrainingJobDetails
            job={selectedJob}
            isOpen={isDetailsOpen}
            onClose={() => setIsDetailsOpen(false)}
          />
        )}
      </Segment>
    </div>
  );
};

export default ModelTraining;
// src/pages/notebook/NotebookList.js

import React, { useState, useEffect } from 'react';
import { Pagination } from 'semantic-ui-react';
import DataList from '../../components/common/DataList';
import { notebookAPI } from '../../api/notebookAPI';
import { toast } from 'react-toastify';
import CreateNotebookForm from './CreateNotebookForm'; // We'll create this component next

const NotebookList = () => {
  const [notebooks, setNotebooks] = useState([]);
  const [loading, setLoading] = useState(false);
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [pageInfo, setPageInfo] = useState({ page: 1, limit: 10, total: 0 });

  const fetchNotebooks = async (page = 1, limit = 10) => {
    setLoading(true);
    try {
      const response = await notebookAPI.getNotebooks(page, limit);
      setNotebooks(response.data);
      // console.log("response notebook : ", response.data)
      setPageInfo({ page: response.page, limit: response.limit, total: response.total });
    } catch (error) {
      toast.error(error.message || 'Failed to fetch notebooks');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchNotebooks();
  }, []);

  const handleCreate = () => {
    setIsCreateOpen(true);
  };

  const handleReset = async (item) => {
    try {
      await notebookAPI.resetNotebook(item.id);
      toast.info('Notebook reset successfully');
      fetchNotebooks(pageInfo.page, pageInfo.limit);
    } catch (error) {
      toast.error(error.message || 'Failed to reset notebook');
    }
  };

  const handleDelete = async (item) => {
    try {
      await notebookAPI.deleteNotebook(item.id);
      toast.info('Notebook deleted successfully');
      fetchNotebooks(pageInfo.page, pageInfo.limit);
    } catch (error) {
      toast.error(error.message || 'Failed to delete notebook');
    }
  };

  const handlePageChange = (e, { activePage }) => {
    fetchNotebooks(activePage, pageInfo.limit);
  };

  const columns = ['name', 'describe', 'resource_memory', 'resource_cpu', 'resource_gpu', 'status'];
  const columnNames = {
    name: 'Name',
    describe: 'Description',
    resource_memory: 'Memory',
    resource_cpu: 'CPU',
    resource_gpu: 'GPU',
    status: 'Status',
    // access_url: 'Access URL',
  };

  const handleAccessEnvironment = (url) => {
    if (!url) {
      toast.error('url is not available');  
      return;
    }
    // new page 
    window.open(url, '_blank');
  };


  if (loading) return <div className='p-4'>Loading...</div>;

  return (
    <div className='p-4'>
      <DataList
        title="Notebook List"
        data={notebooks}
        columns={columns}
        columnNames={columnNames}
        onAdd={handleCreate}
        customActions={[
          {
            icon: 'refresh',
            color: 'orange',
            onClick: handleReset,
            confirm: {
              content: 'Are you sure you want to reset this notebook?',
              confirmButton: 'Reset',
              cancelButton: 'Cancel',
            }
          },
          {
            icon: 'trash',
            color: 'red',
            onClick: handleDelete,
            confirm: {
              content: 'Are you sure you want to delete this notebook?',
              confirmButton: 'Delete',
              cancelButton: 'Cancel',
            }
          },
          {
            icon: 'external',  // 或使用其他合适的图标
            color: 'green',
            onClick: (item) => handleAccessEnvironment(item.access_url),
            tooltip: 'enter environment',  // 鼠标悬停时显示的提示文本  
            // 只在 access_url 存在且状态为 Running 时显示按钮
            show: (item) => item.access_url,
          },
        ]}
      />
      <Pagination
        activePage={pageInfo.page}
        totalPages={Math.ceil(pageInfo.total / pageInfo.limit)}
        onPageChange={handlePageChange}
      />
      {isCreateOpen && (
        <CreateNotebookForm
          isOpen={isCreateOpen}
          onClose={() => setIsCreateOpen(false)}
          onRefresh={() => fetchNotebooks(pageInfo.page, pageInfo.limit)}
        />
      )}
    </div>
  );
};

export default NotebookList;
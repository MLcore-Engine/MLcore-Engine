import React, { useState } from 'react';
import { Table, Input, Button, Icon, Confirm } from 'semantic-ui-react';
import styles from './DataList.module.css';

const DataList = ({ title, data, columns, onAdd, onEdit, onDelete }) => {
  const [searchTerm, setSearchTerm] = useState('');

  const [isConfirmOpen, setIsConfirmOpen] = useState(false);
  const [deleteItemId, setDeleteItemId] = useState(null);

  const handleDeleteclick = (id) => {
    setDeleteItemId(id);
    setIsConfirmOpen(true);
  };

  const handleConfirmDelete = () => {
    if (deleteItemId) {
      onDelete(deleteItemId);
      setIsConfirmOpen(false);
      setDeleteItemId(null);
    }
  };

  const handleSearch = (e, { value }) => setSearchTerm(value);

  const filteredData = data.filter((item) =>
    Object.values(item).some((value) =>
      value.toString().toLowerCase().includes(searchTerm.toLowerCase())
    )
  );

  return (
    <div>
      <h2>{title}</h2>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '1em' }}>
        <Input icon="search" placeholder="搜索..." value={searchTerm} onChange={handleSearch} />
        <Button primary onClick={onAdd}>  
          <Icon name="plus" /> 添加
        </Button>
      </div>
      <Table celled striped>
        <Table.Header>
          <Table.Row>
            {columns.map((column) => (
              <Table.HeaderCell sortable key={column}>{column}</Table.HeaderCell>
            ))}
            <Table.HeaderCell sortable >操作</Table.HeaderCell>
          </Table.Row>
        </Table.Header>
        <Table.Body>
          {filteredData.map((item) => (
            <Table.Row key={item.id}>
              {columns.map((column) => (
                <Table.Cell key={`${item.id}-${column}`}>{item[column]}</Table.Cell>
              ))}
              <Table.Cell>
                <Button icon color="blue" onClick={() => onEdit(item)}>
                  <Icon name="edit" />
                </Button>
                <Button icon color="red" onClick={() => handleDeleteclick(item.id)}>
                  <Icon name="trash" />
                </Button>
              </Table.Cell>
            </Table.Row>
          ))}
          {filteredData.length === 0 && (
            <Table.Row>
              <Table.Cell colSpan={columns.length + 1} textAlign="center">
                无数据
              </Table.Cell>
            </Table.Row>
          )}
        </Table.Body>
      </Table>

      <Confirm
        open={isConfirmOpen}
        content="确定要删除？"
        onCancel={() => setIsConfirmOpen(false)}
        onConfirm={handleConfirmDelete}
        cancelButton="取消"
        confirmButton="确定"
        size='mini'
        className={styles.customConfirmDialog} 
      />
    </div>
  );
};

export default DataList;
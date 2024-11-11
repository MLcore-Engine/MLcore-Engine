import React, { useState } from 'react';
import { Table, Input, Button, Icon, Confirm } from 'semantic-ui-react';

const DataList = ({
  title,
  data,
  columns,
  columnNames,
  onAdd,
  onEdit,
  onDelete,
  customActions,
  renderRow,
}) => {
  
  const [searchTerm, setSearchTerm] = useState('');
  const [isConfirmOpen, setIsConfirmOpen] = useState(false);
  const [deleteItem, setDeleteItem] = useState(null);
  const [confirmAction, setConfirmAction] = useState(null);
  const [confirmContent, setConfirmContent] = useState('');
  const [confirmButtons, setConfirmButtons] = useState({});

  const handleActionClick = (item, action) => {
    if (action.confirm) {
      setConfirmAction(() => () => action.onClick(item));
      setConfirmContent(action.confirm.content || 'are you sure？');
      setConfirmButtons({
        confirmButton: action.confirm.confirmButton || 'yes',
        cancelButton: action.confirm.cancelButton || 'no',
      });
      setIsConfirmOpen(true);
    } else {
      action.onClick(item);
    }
  };

  const handleConfirm = () => {
    if (confirmAction) {
      confirmAction();
      setConfirmAction(null);
    }
    setIsConfirmOpen(false);
    setConfirmContent('');
    setConfirmButtons({});
  };

  // 取消操作
  const handleCancel = () => {
    setIsConfirmOpen(false);
    setConfirmAction(null);
    setConfirmContent('');
    setConfirmButtons({});
  };

  const handleConfirmDelete = () => {
    if (deleteItem) {
      onDelete(deleteItem);
      setIsConfirmOpen(false);
      setDeleteItem(null);
    }
  };


  const renderConfirm = () => {
    
    // if custom action confirm
    if (confirmAction) {
      return (
        <Confirm
          open={isConfirmOpen}
          content={confirmContent}
          onCancel={handleCancel}
          onConfirm={handleConfirm}
          cancelButton={confirmButtons.cancelButton || '取消'}
          confirmButton={confirmButtons.confirmButton || '确定'}
          size="mini"
          dimmer="inverted"
        />
      );
    }

    // if delete confirm
    if (deleteItem) {
      return (
        <Confirm
          open={isConfirmOpen}
          content="are you sure？"
          onCancel={() => setIsConfirmOpen(false)}
          onConfirm={handleConfirmDelete}
          cancelButton="no"
          confirmButton="yes"
          size="mini"
        />
      );
    }

    return null;
  };

  const handleDeleteClick = (item) => {
    setDeleteItem(item);
    setIsConfirmOpen(true);
  };

  const handleSearch = (e, { value }) => setSearchTerm(value);

  const safeFilter = (item) => {
    try {
      return columns.some((column) => {
        const value = item[column];
        if (value == null) return false;
        return value.toString().toLowerCase().includes(searchTerm.toLowerCase());
      });
    } catch (error) {
      console.error('search error :', error);
      return false;
    }
  };



  const filteredData = data.filter(safeFilter);

  return (
    <div>
      <h3>{title}</h3>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '1em' }}>
        <Input icon="search" placeholder="searching..." value={searchTerm} onChange={handleSearch} />
        {onAdd && (
          <Button primary onClick={onAdd}>
            <Icon name="plus" /> Add
          </Button>
        )}
      </div>
      <Table celled striped>
        <Table.Header>
          <Table.Row>
            {columns.map((column) => (
              <Table.HeaderCell key={column}>
                {columnNames[column] || column}
              </Table.HeaderCell>
            ))}
            <Table.HeaderCell>Action</Table.HeaderCell> 
          </Table.Row>
        </Table.Header>
        <Table.Body>
          {filteredData.length > 0 ? (
            filteredData.map((item) => {
              const rowData = renderRow ? renderRow(item) : item;
              return (
                <Table.Row key={item.ID}>
                  {columns.map((column) => (
                    <Table.Cell key={`${item.ID}-${column}`}>
                      {rowData[column]}
                    </Table.Cell>
                  ))}
                  <Table.Cell>
                    {customActions && customActions.length > 0
                      ? customActions.map((action, index) => (
                          <Button
                            key={index}
                            icon
                            color={action.color}
                            onClick={() => handleActionClick(item, action)}
                            aria-label={`${action.color}-action-button`}
                            >
                            <Icon name={action.icon} />
                          </Button>
                        ))
                      : <>
                          {onEdit && (
                            <Button icon color="blue" onClick={() => onEdit(item)}>
                              <Icon name="edit" />
                            </Button>
                          )}
                          {onDelete && (
                            <Button icon color="red" onClick={() => handleDeleteClick(item)}>
                              <Icon name="trash" />
                            </Button>
                          )}
                        </>
                    }
                  </Table.Cell>
                </Table.Row>
              );
            })
          ) : (
            <Table.Row>
              <Table.Cell colSpan={columns.length + 1} textAlign="center">
                No data
              </Table.Cell>
            </Table.Row>
          )}
        </Table.Body>
      </Table>

      {renderConfirm()}
    </div>
  );
};

export default DataList;
import React, { useState } from 'react';
import { Table, Input, Button, Icon, Confirm, Header, Segment } from 'semantic-ui-react';

/**
 * 通用数据列表组件
 * 用于展示结构化数据，支持搜索、添加、编辑、删除等操作
 */
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
  className,
}) => {
  // 状态管理
  const [searchTerm, setSearchTerm] = useState('');
  const [isConfirmOpen, setIsConfirmOpen] = useState(false);
  const [deleteItem, setDeleteItem] = useState(null);
  const [confirmAction, setConfirmAction] = useState(null);
  const [confirmContent, setConfirmContent] = useState('');
  const [confirmButtons, setConfirmButtons] = useState({});

  // 数据安全处理
  const safeData = Array.isArray(data) ? data : [];

  // 事件处理
  const handleActionClick = (item, action) => {
    if (action.confirm) {
      setConfirmAction(() => () => action.onClick(item));
      setConfirmContent(action.confirm.content || '确认此操作?');
      setConfirmButtons({
        confirmButton: action.confirm.confirmButton || '确认',
        cancelButton: action.confirm.cancelButton || '取消',
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
    closeConfirm();
  };

  const handleCancel = () => {
    closeConfirm();
  };

  const closeConfirm = () => {
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

  const handleDeleteClick = (item) => {
    setDeleteItem(item);
    setIsConfirmOpen(true);
  };

  const handleSearch = (e, { value }) => setSearchTerm(value);

  // 安全的数据过滤处理
  const safeFilter = (item) => {
    try {
      if (!item || typeof item !== 'object') return false;
      
      return columns.some((column) => {
        const value = item[column];
        if (value == null) return false;
        return value.toString().toLowerCase().includes(searchTerm.toLowerCase());
      });
    } catch (error) {
      console.error('搜索错误:', error);
      return false;
    }
  };

  const filteredData = safeData.filter(safeFilter);

  // 渲染确认对话框
  const renderConfirm = () => {
    if (confirmAction) {
      return (
        <Confirm
          open={isConfirmOpen}
          content={confirmContent}
          onCancel={handleCancel}
          onConfirm={handleConfirm}
          cancelButton={confirmButtons.cancelButton}
          confirmButton={confirmButtons.confirmButton}
          size="mini"
          dimmer="inverted"
        />
      );
    }

    if (deleteItem) {
      return (
        <Confirm
          open={isConfirmOpen}
          content="确定要删除此项吗?"
          onCancel={() => setIsConfirmOpen(false)}
          onConfirm={handleConfirmDelete}
          cancelButton="取消"
          confirmButton="确认"
          size="mini"
          dimmer="inverted"
        />
      );
    }

    return null;
  };

  // 自定义操作按钮和基础操作按钮渲染
  const renderActionButtons = (item) => {
    if (customActions && customActions.length > 0) {
      return customActions.map((action, index) => (
        <Button
          key={index}
          icon
          className="action-button"
          color={action.color}
          onClick={() => handleActionClick(item, action)}
          aria-label={`${action.color}-action-button`}
          size="tiny"
        >
          <Icon name={action.icon} />
        </Button>
      ));
    }

    return (
      <>
        {onEdit && (
          <Button 
            icon 
            className="action-button" 
            color="blue" 
            onClick={() => onEdit(item)} 
            size="tiny"
          >
            <Icon name="edit" />
          </Button>
        )}
        {onDelete && (
          <Button 
            icon 
            className="action-button" 
            color="red" 
            onClick={() => handleDeleteClick(item)} 
            size="tiny"
          >
            <Icon name="trash" />
          </Button>
        )}
      </>
    );
  };

  // 主体渲染
  return (
    <div className={`data-list-container ${className || ''}`}>
      <div className="datalist-header">
        <Header as="h2" className="title-text">{title}</Header>
        <div className="datalist-actions">
          <Input 
            icon="search" 
            placeholder="搜索..." 
            value={searchTerm} 
            onChange={handleSearch}
            size="small"
            className="search-input"
          />
          {onAdd && (
            <Button 
              primary 
              size="small" 
              onClick={onAdd}
              className="add-button"
            >
              <Icon name="plus" /> 添加
            </Button>
          )}
        </div>
      </div>
      
      <div className="table-container">
        <Table celled selectable>
          <Table.Header>
            <Table.Row>
              {columns.map((column) => (
                <Table.HeaderCell key={column}>
                  {columnNames[column] || column}
                </Table.HeaderCell>
              ))}
              <Table.HeaderCell width={2} textAlign="center">操作</Table.HeaderCell> 
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {filteredData.length > 0 ? (
              filteredData.map((item) => {
                const rowData = renderRow ? renderRow(item) : item;
                const rowKey = item.id || item.ID || `row-${Math.random()}`;
                
                return (
                  <Table.Row key={rowKey} className="data-row">
                    {columns.map((column) => (
                      <Table.Cell key={`${rowKey}-${column}`}>
                        {rowData[column]}
                      </Table.Cell>
                    ))}
                    <Table.Cell textAlign="center">
                      <div className="action-buttons">
                        {renderActionButtons(item)}
                      </div>
                    </Table.Cell>
                  </Table.Row>
                );
              })
            ) : (
              <Table.Row>
                <Table.Cell colSpan={columns.length + 1} textAlign="center">
                  <div className="empty-state">
                    <Icon name="database" size="large" />
                    <div>暂无数据</div>
                  </div>
                </Table.Cell>
              </Table.Row>
            )}
          </Table.Body>
        </Table>
      </div>

      {renderConfirm()}

      <style jsx>{`
        .data-list-container {
          background-color: white;
          border-radius: 0.75rem;
          padding: 1.5rem;
          box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
        }
        .datalist-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: 1.5rem;
          flex-wrap: wrap;
          gap: 1rem;
        }
        .title-text {
          margin: 0 !important;
          color: #333 !important;
          font-weight: 600 !important;
        }
        .datalist-actions {
          display: flex;
          gap: 0.75rem;
          align-items: center;
          flex-wrap: wrap;
        }
        .search-input {
          min-width: 240px;
        }
        .add-button {
          display: flex;
          align-items: center;
          gap: 0.3rem;
          padding: 0.6rem 1rem !important;
        }
        .table-container {
          border-radius: 0.5rem;
          overflow: hidden;
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.03);
          margin-bottom: 1rem;
        }
        .action-buttons {
          display: flex;
          gap: 0.5rem;
          justify-content: center;
        }
        .empty-state {
          padding: 2rem 0;
          color: #888;
          display: flex;
          flex-direction: column;
          align-items: center;
          gap: 0.5rem;
        }
      `}</style>
    </div>
  );
};

export default DataList;
import React from 'react';
import { Button } from 'semantic-ui-react';
import DataList from './DataList';
import { getRoleName } from '../../utils/roleUtils';

// 自定义操作按钮类型
interface CustomAction {
  icon: string;
  color: string;
  onClick: (user: any) => void;
  confirm?: {
    content: (user: any) => string;
    confirmButton: string;
    cancelButton: string;
  };
}

// 组件属性类型
interface UserTableProps {
  users: any[];
  loading?: boolean;
  title: string;
  columns: string[];
  columnNames: Record<string, string>;
  onAdd?: () => void;
  onEdit?: (user: any) => void;
  onDelete?: (user: any) => void;
  customActions?: CustomAction[];
  renderRow?: (user: any) => any;
  className?: string;
}

/**
 * 通用用户表格组件
 */
const UserTable: React.FC<UserTableProps> = ({
  users,
  loading = false,
  title,
  columns,
  columnNames,
  onAdd,
  onEdit,
  onDelete,
  customActions = [],
  renderRow,
  className = ''
}) => {
  // 默认行数据转换函数
  const defaultRenderRow = (user: any) => ({
    ...user,
    role: getRoleName(user.role),
    status: user.status === 1 ? '启用' : '禁用'
  });

  if (loading && users.length === 0) {
    return <div>加载中...</div>;
  }

  return (
    <div className="p-4">
      <h2>{title}</h2>
      <DataList
        title={title}
        data={users}
        columns={columns}
        columnNames={columnNames}
        onAdd={onAdd}
        onEdit={onEdit}
        onDelete={onDelete}
        customActions={customActions}
        renderRow={renderRow || defaultRenderRow}
        className={className}
      />
    </div>
  );
};

export default UserTable; 
import React from 'react';
import { Dropdown, DropdownProps } from 'semantic-ui-react';
import { roleOptions, projectRoleOptions } from '../../utils/roleUtils';

interface RoleSelectorProps {
  value: number;
  onChange: (event: React.SyntheticEvent<HTMLElement>, data: DropdownProps) => void;
  placeholder?: string;
  fluid?: boolean;
  isProjectRole?: boolean; // 是否是项目角色选择
}

/**
 * 通用角色选择器组件
 */
const RoleSelector: React.FC<RoleSelectorProps> = ({
  value,
  onChange,
  placeholder = "选择角色",
  fluid = true,
  isProjectRole = false
}) => {
  const options = isProjectRole ? projectRoleOptions : roleOptions;

  return (
    <Dropdown
      selection
      fluid={fluid}
      options={options}
      value={value}
      onChange={onChange}
      placeholder={placeholder}
    />
  );
};

export default RoleSelector; 
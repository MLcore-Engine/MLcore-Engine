import React, { useState } from 'react';
import { Button, Modal, Form, Message } from 'semantic-ui-react';
import DataList from '../../components/common/DataList';
import { useProjects } from '../../context/ProjectContext';
import { toast } from 'react-toastify';
import { Project } from '../../api/projectAPI';


// Project结构定义:
// interface Project {
//   id: string;           // 项目ID，字符串类型
//   name: string;         // 项目名称
//   description: string;  // 项目描述
//   createdAt?: string;   // 创建时
//   updatedAt?: string;   // 更新时间
//   users?: MemberDTO[];  // 项目成员列表
// }

/**
 * 项目组织管理组件
 * 显示项目列表并提供添加、编辑、删除项目功能
 */
const ProjectGroupOrg = () => {
  const {
    projects,
    loading,
    createProject,
    updateProject,
    deleteProject,
  } = useProjects();

  
  // 控制弹窗（Modal）是否打开的状态
  // false 表示弹窗关闭，true 表示弹窗打开
  const [isModalOpen, setIsModalOpen] = useState(false);

  // 弹窗类型，'add' 表示添加项目，'edit' 表示编辑项目
  // 通过该状态区分弹窗当前是新增还是编辑模式
  const [modalType, setModalType] = useState('');

  // 弹窗内表单的数据，包含项目名称和描述
  // 用于绑定表单输入，实现受控组件
  const [projectData, setProjectData] = useState({ name: '', description: '' });

  // 弹窗内的错误信息
  // 当表单校验失败或接口报错时，显示具体错误内容
  const [modalError, setModalError] = useState('');

  // 当前选中的项目ID，仅在编辑或删除时使用
  // 用于标识需要操作的具体项目
  const [selectedProjectId, setSelectedProjectId] = useState(null);

  // 提交状态，true 表示正在提交（防止重复提交）
  // 用于在提交过程中禁用按钮或显示加载状态
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleAdd = () => {
    setModalType('add');
    setProjectData({ name: '', description: '' });
    setModalError('');
    setIsModalOpen(true);
  };

  /**
   * 处理编辑项目操作
   * @param {Object} project - 项目对象，包含id、name、description等字段
   */
  const handleEdit = (project: Project) => {
    setModalType('edit');
    setProjectData({ name: project.name, description: project.description });
    setSelectedProjectId(project.id);
    setModalError('');
    setIsModalOpen(true);
  };

  /**
   * 处理删除项目操作
   * @param {Object} item - 项目对象，必须包含id字段
   */
  const handleDelete = async (item: Project) => {
    try {
      await deleteProject(item.id);
    } catch (err) {
      // console.error(err.message);
      toast.info(err.message || '项目删除失败');
    }
  };

  const handleSubmit = async () => {
    if (!projectData.name || !projectData.name.trim()) {
      setModalError('项目名称不能为空');  
      return;
    }

    setIsSubmitting(true);
    setModalError('');

    try {
      if (modalType === 'add') {
        await createProject(projectData);
      } else if (modalType === 'edit' && selectedProjectId) {
        await updateProject(selectedProjectId, projectData);
      }
      
      setIsModalOpen(false);
      setProjectData({ name: '', description: '' });
      setSelectedProjectId(null);
    } catch (err) {
      setModalError(err instanceof Error ? err.message : '操作失败');
      toast.error(err instanceof Error ? err.message : '操作失败');
    } finally {
      setIsSubmitting(false);
    }
  };

  const columns = ['name', 'description'];
  const columnNames = { name: '项目名称', description: '描述' };

  if (loading) return <div className='p-4'>加载中...</div>;

  return (
    <div className='p-4'>
      <DataList
        title="项目列表"
        data={projects || []}
        columns={columns}
        columnNames={columnNames}
        onAdd={handleAdd}
        onEdit={handleEdit}
        onDelete={handleDelete}
        customActions={[]}
        renderRow={(item) => item}
        className="project-list"
      />

      <Modal
              open={isModalOpen}
              onClose={() => setIsModalOpen(false)}
              closeOnDimmerClick={!isSubmitting}
              closeOnEscape={!isSubmitting}
              aria-labelledby="modal-header"
      >
        <Modal.Header id="modal-header">{modalType === 'add' ? '添加项目' : '编辑项目'}</Modal.Header>
        <Modal.Content>
          <Form error={!!modalError}>
            <Form.Input
              label="项目名称"
              value={projectData.name}
              onChange={(e, { value }) => setProjectData({ ...projectData, name: value })}
              placeholder="输入项目名称"
              aria-label="project-name"
              required
            />
            <Form.TextArea
              label="描述"
              value={projectData.description}
              onChange={(e, { value }) => setProjectData({ ...projectData, description: String(value) })}
              placeholder="输入项目描述"
              aria-label="project-description"
            />
            {modalError && <Message error content={modalError} />}
          </Form>
        </Modal.Content>
        <Modal.Actions>
          <Button 
            onClick={() => setIsModalOpen(false)}
            disabled={isSubmitting}
            aria-label="cancel-button"
          >
            取消
          </Button>
          <Button 
            primary 
            onClick={handleSubmit}
            loading={isSubmitting}
            disabled={isSubmitting}
            aria-label={modalType === 'add' ? '添加' : '保存'}
          >
            {modalType === 'add' ? '添加' : '保存'}
          </Button>
        </Modal.Actions>
      </Modal>
    </div>
  );
};

export default ProjectGroupOrg;
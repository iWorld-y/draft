import React from 'react';
import type { UploadTask } from '../../services/dictionary';
import './UploadStatus.css';

interface UploadStatusProps {
  task: UploadTask;
}

const UploadStatus: React.FC<UploadStatusProps> = ({ task }) => {
  const getStatusText = () => {
    switch (task.status) {
      case 'pending':
        return '等待处理...';
      case 'processing':
        return '正在处理...';
      case 'completed':
        return '处理完成！';
      case 'failed':
        return '处理失败';
      default:
        return '未知状态';
    }
  };

  const getStatusClass = () => {
    switch (task.status) {
      case 'completed':
        return 'status-completed';
      case 'failed':
        return 'status-failed';
      default:
        return 'status-processing';
    }
  };
  
  return (
    <div className="upload-status">
      <div className={`status-header ${getStatusClass()}`}>
        <span className="status-icon">
          {task.status === 'completed' ? '✓' : 
           task.status === 'failed' ? '✗' : '⟳'}
        </span>
        <span className="status-text">{getStatusText()}</span>
      </div>
      
      <div className="progress-wrapper">
        <div className="progress-bar-bg">
          <div 
            className={`progress-bar-fill ${getStatusClass()}`}
            style={{ width: `${task.progress}%` }}
          />
        </div>
        <span className="progress-text">{task.progress}%</span>
      </div>
      
      {task.message && (
        <div className="status-message">{task.message}</div>
      )}
    </div>
  );
};

export default UploadStatus;

import React from 'react';
import type { UploadTask } from '../../services/dictionary';
import './UploadStatus.css';

interface UploadStatusProps {
  task: UploadTask;
}

const UploadStatus: React.FC<UploadStatusProps> = ({ task }) => {
  const failedCount = task.failed_words?.length ?? 0;
  const failedDetails = task.failed_details ?? [];

  const getStageLabel = (stage: string) => {
    switch (stage) {
      case 'translate':
        return '翻译';
      case 'save':
        return '保存';
      case 'reuse':
        return '复用';
      default:
        return stage;
    }
  };

  const getStatusText = () => {
    switch (task.status) {
      case 'pending':
        return '等待处理...';
      case 'processing':
        return '正在处理...';
      case 'completed':
        return failedCount > 0 ? '处理完成（部分失败）' : '处理完成！';
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

      {failedDetails.length > 0 && (
        <div className="failed-details">
          <div className="failed-details-title">失败详情（{failedDetails.length}）</div>
          <div className="failed-details-list">
            {failedDetails.map((item, index) => (
              <div className="failed-detail-item" key={`${item.word}-${index}`}>
                <div className="failed-detail-main">
                  <span className="failed-detail-word">{item.word}</span>
                  <span className="failed-detail-stage">{getStageLabel(item.stage)}</span>
                </div>
                <div className="failed-detail-reason">{item.reason}</div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default UploadStatus;

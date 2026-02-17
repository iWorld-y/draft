import React, { useState, useRef, useCallback } from 'react';
import { uploadDictionary, getUploadStatus, type UploadTask } from '../../services/dictionary';
import UploadStatus from '../../components/UploadStatus';
import './DictionaryUpload.css';

const POLLING_INTERVAL = 2000; // 2 seconds

const DictionaryUpload: React.FC = () => {
  const [isDragging, setIsDragging] = useState(false);
  const [uploadTask, setUploadTask] = useState<UploadTask | null>(null);
  const [isUploading, setIsUploading] = useState(false);
  const pollingRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const stopPolling = useCallback(() => {
    if (pollingRef.current) {
      clearInterval(pollingRef.current);
      pollingRef.current = null;
    }
  }, []);

  const startPolling = useCallback((taskId: string) => {
    stopPolling();
    
    pollingRef.current = setInterval(async () => {
      try {
        const response = await getUploadStatus(taskId);
        const task = response.data;
        const failedCount = task.failed_words?.length ?? 0;
        const failedPreview = (task.failed_words || []).slice(0, 5).join(', ');
        const message =
          task.status === 'failed'
            ? `导入失败：${failedCount} 个词未成功解析。${failedPreview ? `示例：${failedPreview}` : ''}`
            : undefined;
        setUploadTask({
          ...task,
          message,
        });
        
        if (task.status === 'completed' || task.status === 'failed') {
          stopPolling();
          if (task.status === 'completed') {
            // Navigate to dictionary list after 2 seconds
            setTimeout(() => {
              window.location.href = '/';
            }, 2000);
          }
        }
      } catch (error) {
        console.error('Failed to get upload status:', error);
      }
    }, POLLING_INTERVAL);
  }, [stopPolling]);

  const handleFileUpload = async (file: File) => {
    if (!file.name.endsWith('.txt')) {
      alert('请上传 TXT 格式的文件');
      return;
    }

    setIsUploading(true);
    try {
      const response = await uploadDictionary(file);
      const { task_id } = response.data;
      
      setUploadTask({
        task_id,
        status: 'pending',
        progress: 0
      });
      
      startPolling(task_id);
    } catch (error) {
      console.error('Upload failed:', error);
      const message = error instanceof Error ? error.message : '上传失败，请重试';
      alert(message);
    } finally {
      setIsUploading(false);
    }
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
    
    const files = e.dataTransfer.files;
    if (files.length > 0) {
      handleFileUpload(files[0]);
    }
  };

  const handleFileInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (files && files.length > 0) {
      handleFileUpload(files[0]);
    }
  };

  const handleClick = () => {
    fileInputRef.current?.click();
  };

  return (
    <div className="dictionary-upload-page">
      <div className="upload-container">
        <h1>上传词典</h1>
        <p className="subtitle">支持 TXT 格式文件，每行一个单词</p>
        
        {!uploadTask ? (
          <div
            className={`upload-area ${isDragging ? 'dragging' : ''} ${isUploading ? 'uploading' : ''}`}
            onDragOver={handleDragOver}
            onDragLeave={handleDragLeave}
            onDrop={handleDrop}
            onClick={handleClick}
          >
            <input
              ref={fileInputRef}
              type="file"
              accept=".txt"
              onChange={handleFileInputChange}
              style={{ display: 'none' }}
            />
            <div className="upload-icon">📁</div>
            <p className="upload-text">
              {isUploading ? '上传中...' : '点击或拖拽文件到此处上传'}
            </p>
            <p className="upload-hint">支持 .txt 格式</p>
          </div>
        ) : (
          <UploadStatus task={uploadTask} />
        )}
        
        <button 
          className="back-button"
          onClick={() => window.location.href = '/'}
        >
          返回列表
        </button>
      </div>
    </div>
  );
};

export default DictionaryUpload;

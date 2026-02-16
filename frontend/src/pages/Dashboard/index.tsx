import React, { useState, useEffect } from 'react';
import { getDictionaries, deleteDictionary, type Dictionary } from '../../services/dictionary';
import './Dashboard.css';

const Dashboard: React.FC = () => {
  const [dictionaries, setDictionaries] = useState<Dictionary[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [stats, setStats] = useState({ totalWords: 0, dictCount: 0 });

  useEffect(() => {
    loadDictionaries();
  }, []);

  const loadDictionaries = async () => {
    setIsLoading(true);
    try {
      const response = await getDictionaries();
      const dicts = response.data;
      setDictionaries(dicts);
      
      const totalWords = dicts.reduce((sum, d) => sum + d.word_count, 0);
      setStats({
        totalWords,
        dictCount: dicts.length
      });
    } catch (error) {
      console.error('Failed to load dictionaries:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('确定要删除这个词典吗？')) {
      return;
    }
    
    try {
      await deleteDictionary(id);
      await loadDictionaries();
    } catch (error) {
      console.error('Failed to delete dictionary:', error);
      alert('删除失败，请重试');
    }
  };

  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr);
    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };

  return (
    <div className="dashboard-page">
      <div className="dashboard-container">
        <header className="dashboard-header">
          <h1>我的词典</h1>
          <div className="header-actions">
            <div className="stats">
              <span className="stat-item">
                <strong>{stats.dictCount}</strong> 个词典
              </span>
              <span className="stat-item">
                <strong>{stats.totalWords}</strong> 个单词
              </span>
            </div>
            <button 
              className="upload-btn"
              onClick={() => window.location.href = '/upload'}
            >
              + 上传词典
            </button>
          </div>
        </header>

        {isLoading ? (
          <div className="loading-container">
            <div className="loading-spinner"></div>
            <p>加载中...</p>
          </div>
        ) : dictionaries.length === 0 ? (
          <div className="empty-state">
            <div className="empty-icon">📚</div>
            <h3>还没有词典</h3>
            <p>上传你的第一个词典开始学习吧</p>
            <button 
              className="primary-button"
              onClick={() => window.location.href = '/upload'}
            >
              上传词典
            </button>
          </div>
        ) : (
          <div className="dictionary-list">
            {dictionaries.map(dict => (
              <div key={dict.id} className="dictionary-card">
                <div className="dict-info">
                  <h3 className="dict-name">{dict.name}</h3>
                  <div className="dict-meta">
                    <span className="word-count">{dict.word_count} 词</span>
                    <span className="created-at">创建于 {formatDate(dict.created_at)}</span>
                  </div>
                </div>
                <div className="dict-actions">
                  <button 
                    className="learn-btn"
                    onClick={() => window.location.href = `/learn?dictId=${dict.id}`}
                  >
                    开始学习
                  </button>
                  <button 
                    className="delete-btn"
                    onClick={() => handleDelete(dict.id)}
                  >
                    删除
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default Dashboard;

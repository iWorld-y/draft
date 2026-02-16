import React, { useState } from 'react';
import WordCard from '../../components/WordCard';
import QualityButtons from '../../components/QualityButtons';
import ProgressBar from '../../components/ProgressBar';
import { useLearning } from '../../hooks/useLearning';
import './Learning.css';

interface LearningPageProps {
  dictId?: number;
}

const Learning: React.FC<LearningPageProps> = ({ dictId = 1 }) => {
  const [revealed, setRevealed] = useState(false);
  const { 
    currentWord, 
    progress, 
    isLoading, 
    isFinished, 
    submitAnswer, 
    loadNextWord,
    loadTasks
  } = useLearning(dictId);

  const handleReveal = () => {
    setRevealed(true);
  };

  const handleQualitySelect = async (quality: number) => {
    await submitAnswer(quality);
    setRevealed(false);
    loadNextWord();
  };

  const handleRestart = () => {
    setRevealed(false);
    loadTasks();
  };

  if (isFinished) {
    return (
      <div className="learning-page">
        <div className="completion-card">
          <div className="completion-icon">🎉</div>
          <h2>恭喜完成今日学习！</h2>
          <p>你已经完成了 {progress.total} 个单词的学习</p>
          <div className="completion-actions">
            <button className="primary-button" onClick={handleRestart}>
              再来一组
            </button>
            <button className="secondary-button" onClick={() => window.location.href = '/'}>
              返回首页
            </button>
          </div>
        </div>
      </div>
    );
  }

  if (!currentWord) {
    return (
      <div className="learning-page">
        <div className="loading-container">
          <div className="loading-spinner"></div>
          <p>加载中...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="learning-page">
      <div className="learning-container">
        <ProgressBar current={progress.completed} total={progress.total} />
        
        <div className="card-container">
          <WordCard 
            word={currentWord} 
            onReveal={handleReveal}
          />
        </div>
        
        {revealed && (
          <QualityButtons 
            onSelect={handleQualitySelect}
            disabled={isLoading}
          />
        )}
      </div>
    </div>
  );
};

export default Learning;

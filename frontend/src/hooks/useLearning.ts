import { useState, useEffect, useCallback } from 'react';
import { getTodayTasks, submitLearning, type Word, type SubmitLearningData } from '../services/learning';

interface Progress {
  completed: number;
  total: number;
}

interface UseLearningReturn {
  currentWord: Word | null;
  progress: Progress;
  isLoading: boolean;
  isFinished: boolean;
  loadTasks: () => Promise<void>;
  submitAnswer: (quality: number) => Promise<void>;
  loadNextWord: () => void;
}

export const useLearning = (dictId: number = 1): UseLearningReturn => {
  const [words, setWords] = useState<Word[]>([]);
  const [currentIndex, setCurrentIndex] = useState(0);
  const [progress, setProgress] = useState<Progress>({ completed: 0, total: 0 });
  const [isLoading, setIsLoading] = useState(false);
  const [isFinished, setIsFinished] = useState(false);

  const loadTasks = useCallback(async () => {
    setIsLoading(true);
    try {
      const data = await getTodayTasks({ dict_id: dictId, limit: 20 });
      const { words: taskWords } = data;

      setWords(taskWords);
      setProgress({
        completed: 0,
        total: taskWords.length
      });
      setIsFinished(taskWords.length === 0);
      setCurrentIndex(0);
    } catch (error) {
      console.error('Failed to load tasks:', error);
      alert('加载任务失败，请重试');
    } finally {
      setIsLoading(false);
    }
  }, [dictId]);

  useEffect(() => {
    loadTasks();
  }, [loadTasks]);

  const submitAnswer = async (quality: number) => {
    const currentWord = words[currentIndex];
    if (!currentWord) return;

    setIsLoading(true);
    try {
      const data: SubmitLearningData = {
        word_id: currentWord.id,
        quality,
      };
      await submitLearning(data);
      setProgress(prev => ({
        ...prev,
        completed: prev.completed + 1
      }));
    } catch (error) {
      console.error('Failed to submit answer:', error);
      alert('提交失败，请重试');
    } finally {
      setIsLoading(false);
    }
  };

  const loadNextWord = () => {
    if (currentIndex + 1 >= words.length) {
      setIsFinished(true);
    } else {
      setCurrentIndex(prev => prev + 1);
    }
  };

  return {
    currentWord: words[currentIndex] || null,
    progress,
    isLoading,
    isFinished,
    loadTasks,
    submitAnswer,
    loadNextWord
  };
};

export default useLearning;

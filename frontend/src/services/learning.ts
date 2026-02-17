import request from './request';

export interface Definition {
  pos: string;
  text: string;
}

export interface WordMeaning {
  definitions: Definition[];
}

export interface Word {
  id: number;
  word: string;
  phonetic?: string;
  meaning: WordMeaning;
  example?: string;
}

export interface LearningTask {
  words: Word[];
  review_count: number;
  new_count: number;
}

export interface SubmitLearningData {
  word_id: number;
  quality: number;
  dictionary_id: number;
}

// Get today's learning tasks
export const getTodayTasks = (params: { dict_id: number; limit?: number }): Promise<{ data: LearningTask }> => {
  return request.get('/learning/today-tasks', { params });
};

// Submit learning result
export const submitLearning = (data: SubmitLearningData): Promise<void> => {
  return request.post('/learning/submit', data);
};

// Get learning statistics
export const getLearningStats = (): Promise<{ data: { total_learned: number; streak_days: number } }> => {
  return request.get('/learning/stats');
};

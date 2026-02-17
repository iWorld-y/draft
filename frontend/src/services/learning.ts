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
}

interface WordPB {
  id: number;
  word: string;
  phonetic?: string;
  meaning?: string;
  example?: string;
}

interface LearningTaskPB {
  words?: WordPB[];
  reviewCount: number;
  newCount: number;
}

const parseMeaning = (raw?: string): WordMeaning => {
  if (!raw) {
    return { definitions: [] };
  }

  try {
    const decoded = atob(raw);
    return JSON.parse(decoded) as WordMeaning;
  } catch {
    return { definitions: [] };
  }
};

export const getTodayTasks = async (params: { dict_id: number; limit?: number }): Promise<LearningTask> => {
  const data = (await request.get('/learning/today-tasks', {
    params: {
      dictId: params.dict_id,
      limit: params.limit,
    },
  })) as LearningTaskPB;

  const words = Array.isArray(data.words)
    ? data.words.map((w) => ({
        id: w.id,
        word: w.word,
        phonetic: w.phonetic,
        meaning: parseMeaning(w.meaning),
        example: w.example,
      }))
    : [];

  return {
    words,
    review_count: data.reviewCount || 0,
    new_count: data.newCount || 0,
  };
};

export const submitLearning = async (data: SubmitLearningData): Promise<void> => {
  await request.post('/learning/submit', {
    wordId: data.word_id,
    quality: data.quality,
    timeSpent: 0,
  });
};

export const getLearningStats = async (): Promise<{ total_learned: number; streak_days: number }> => {
  const data = (await request.get('/learning/stats')) as { totalLearned: number; streakDays: number };
  return {
    total_learned: data.totalLearned,
    streak_days: data.streakDays,
  };
};

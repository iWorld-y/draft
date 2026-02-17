import request from './request';

export interface Dictionary {
  id: number;
  name: string;
  total_words: number;
  learned_words: number;
  progress: number;
  description: string;
  created_at: string;
}

export interface UploadTask {
  task_id: string;
  status: 'pending' | 'processing' | 'completed' | 'failed';
  progress: number;
  total?: number;
  processed?: number;
  failed_words?: string[];
  failed_details?: UploadFailedDetail[];
  message?: string;
}

export interface UploadFailedDetail {
  word: string;
  stage: string;
  reason: string;
  at: string;
}

interface DictionaryItemPB {
  id: number;
  name: string;
  description: string;
  totalWords: number;
  learnedWords: number;
  progress: number;
  createdAt: string;
}

interface UploadStatusPB {
  taskId: string;
  status: UploadTask['status'];
  progress: number;
  total: number;
  processed: number;
  failedWords?: string[];
}

const toBase64 = async (file: File): Promise<string> => {
  const buffer = await file.arrayBuffer();
  const bytes = new Uint8Array(buffer);
  let binary = '';
  for (let i = 0; i < bytes.length; i += 1) {
    binary += String.fromCharCode(bytes[i]);
  }
  return btoa(binary);
};

export const uploadDictionary = async (file: File): Promise<{ task_id: string }> => {
  const fileContent = await toBase64(file);
  const data = (await request.post('/dictionaries/upload', {
    fileContent,
    name: file.name.replace(/\.txt$/i, ''),
    description: '',
  })) as { taskId: string };

  return {
    task_id: data.taskId,
  };
};

export const getUploadStatus = async (taskId: string): Promise<UploadTask> => {
  const data = (await request.get(`/dictionaries/upload/status/${taskId}`)) as UploadStatusPB;
  return {
    task_id: data.taskId,
    status: data.status,
    progress: data.progress,
    total: data.total,
    processed: data.processed,
    failed_words: data.failedWords || [],
  };
};

export const getDictionaries = async (): Promise<{ items: Dictionary[] }> => {
  const data = (await request.get('/dictionaries')) as { items?: DictionaryItemPB[] };
  const items = Array.isArray(data.items) ? data.items : [];

  return {
    items: items.map((item: DictionaryItemPB) => ({
      id: item.id,
      name: item.name,
      description: item.description,
      total_words: item.totalWords,
      learned_words: item.learnedWords,
      progress: item.progress,
      created_at: item.createdAt,
    })),
  };
};

export const deleteDictionary = async (_id: number): Promise<void> => {
  throw new Error('当前后端 proto 契约未提供删除词典接口');
};

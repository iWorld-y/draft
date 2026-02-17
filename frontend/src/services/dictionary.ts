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
  message?: string;
}

// Upload dictionary file
export const uploadDictionary = (file: File): Promise<{ data: { task_id: string } }> => {
  const formData = new FormData();
  formData.append('file', file);
  
  return request.post('/dictionaries/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};

// Query upload progress
export const getUploadStatus = (taskId: string): Promise<{ data: UploadTask }> => {
  return request.get(`/dictionaries/upload/status/${taskId}`);
};

// Get dictionary list
export const getDictionaries = (): Promise<{ data: { items: Dictionary[] } }> => {
  return request.get('/dictionaries');
};

// Delete dictionary
export const deleteDictionary = (id: number): Promise<void> => {
  return request.delete(`/dictionaries/${id}`);
};

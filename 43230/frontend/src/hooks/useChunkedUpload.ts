import { useState, useCallback } from 'react';
import { fileApi } from '@/services/api';
import { message } from 'antd';

const CHUNK_SIZE = 5 * 1024 * 1024;

export function useChunkedUpload() {
  const [uploading, setUploading] = useState(false);
  const [progress, setProgress] = useState(0);
  const [uploadId, setUploadId] = useState<string | null>(null);

  const uploadFile = useCallback(async (file: File, onProgress?: (percent: number) => void) => {
    setUploading(true);
    setProgress(0);

    try {
      const totalChunks = Math.ceil(file.size / CHUNK_SIZE);

      const initResponse = await fileApi.initiateUpload({
        file_name: file.name,
        file_size: file.size,
        file_type: file.type,
      });

      const { upload_id } = initResponse.data;
      setUploadId(upload_id);

      for (let i = 0; i < totalChunks; i++) {
        const start = i * CHUNK_SIZE;
        const end = Math.min(start + CHUNK_SIZE, file.size);
        const chunk = file.slice(start, end);

        await fileApi.uploadChunk(upload_id, i, chunk);

        const percent = Math.round(((i + 1) / totalChunks) * 100);
        setProgress(percent);
        onProgress?.(percent);
      }

      const completeResponse = await fileApi.completeUpload(upload_id);
      message.success('文件上传成功');

      setUploading(false);
      return completeResponse.data;
    } catch (error: any) {
      setUploading(false);
      message.error(error.response?.data?.error || '上传失败');
      throw error;
    }
  }, []);

  return {
    uploading,
    progress,
    uploadId,
    uploadFile,
  };
}

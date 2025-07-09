ALTER TABLE transcriptions
ADD COLUMN original_transcription_id BIGINT CHECK (original_transcription_id >= 0);
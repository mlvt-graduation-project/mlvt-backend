ALTER TABLE audios
ADD COLUMN transcription_id BIGINT CHECK (transcription_id >= 0);
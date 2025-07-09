-- 1. Add the title column to the audios table
ALTER TABLE audios
  ADD COLUMN title TEXT NOT NULL;

-- 2. Add the title column to the transcriptions table
ALTER TABLE transcriptions
  ADD COLUMN title TEXT NOT NULL;


-- 1. Add the title column to the audios table
ALTER TABLE audios
  ADD COLUMN title TEXT NOT NULL;

-- 2. Add the title column to the transcriptions table
ALTER TABLE transcriptions
  ADD COLUMN title TEXT NOT NULL;

-- 3. Add is_deleted column to audios table
ALTER TABLE audios
  ADD COLUMN is_deleted BOOLEAN NOT NULL DEFAULT FALSE;

-- 4. Add is_deleted column to transcriptions table
ALTER TABLE transcriptions
  ADD COLUMN is_deleted BOOLEAN NOT NULL DEFAULT FALSE;

-- 5. Add is_deleted column to videos table
ALTER TABLE videos
  ADD COLUMN is_deleted BOOLEAN NOT NULL DEFAULT FALSE;
  
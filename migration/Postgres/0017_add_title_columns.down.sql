-- 1. Drop the title column from the audios table
ALTER TABLE audios DROP COLUMN IF EXISTS title;

-- 2. Drop the title column from the transcriptions table
ALTER TABLE transcriptions DROP COLUMN IF EXISTS title;

-- 3. Drop the is_deleted column from the audios table
ALTER TABLE audios DROP COLUMN IF EXISTS is_deleted;

-- 4. Drop the is_deleted column from the transcriptions table
ALTER TABLE transcriptions DROP COLUMN IF EXISTS is_deleted;

-- 5. Drop the is_deleted column from the videos table
ALTER TABLE videos DROP COLUMN IF EXISTS is_deleted;

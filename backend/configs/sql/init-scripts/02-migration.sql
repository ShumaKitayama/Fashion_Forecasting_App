-- Migration to add updated_at column to trend_records table
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name = 'trend_records' 
                   AND column_name = 'updated_at') THEN
        ALTER TABLE trend_records ADD COLUMN updated_at TIMESTAMP DEFAULT NOW();
    END IF;
END $$; 
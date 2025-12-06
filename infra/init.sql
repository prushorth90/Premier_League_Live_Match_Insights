REATE EXTENSION IF NOT EXISTS vector;

-- 2. Create a table for Matches
-- Stores high-level metadata about the game.
CREATE TABLE IF NOT EXISTS matches (
    id SERIAL PRIMARY KEY,
    home_team VARCHAR(100) NOT NULL,
    away_team VARCHAR(100) NOT NULL,
    match_date TIMESTAMP WITH TIME ZONE NOT NULL,
    competition VARCHAR(100) DEFAULT 'Premier League',
    final_score_home INT,
    final_score_away INT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 3. Create a table for Match Events (The "Context" for RAG)
-- Stores the actual text commentary we will feed into the LLM.
CREATE TABLE IF NOT EXISTS match_events (
    id SERIAL PRIMARY KEY,
    match_id INT REFERENCES matches(id) ON DELETE CASCADE,
    minute_label VARCHAR(10), -- e.g., "45+2'"
    event_type VARCHAR(50),   -- e.g., "Goal", "Yellow Card", "Substitution"
    player_name VARCHAR(100),
    description TEXT NOT NULL, -- The main commentary text
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 4. Create a table for Vector Embeddings
-- Stores the mathematical representation of the event description.
-- We use a separate table to keep the main events table lightweight.
CREATE TABLE IF NOT EXISTS match_embeddings (
    id SERIAL PRIMARY KEY,
    event_id INT REFERENCES match_events(id) ON DELETE CASCADE,
    
    -- VECTOR(1536) is the standard size for OpenAI's text-embedding-3-small
    -- If you use a local model (like all-MiniLM-L6-v2), change this to 384.
    embedding_vector VECTOR(1536), 
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 5. Create an HNSW Index for fast similarity search
-- usage: vector_cosine_ops is for Cosine Similarity (standard for RAG)
-- m=16 and ef_construction=64 are good defaults for performance/accuracy balance.
CREATE INDEX ON match_embeddings 
USING hnsw (embedding_vector vector_cosine_ops) 
WITH (m = 16, ef_construction = 64);

-- Optional: Create a standard index on match_id for faster joins
CREATE INDEX idx_match_events_match_id ON match_events(match_id);
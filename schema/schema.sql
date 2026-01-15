CREATE TABLE venues (
    venue_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    venue_name VARCHAR(100) NOT NULL,
    full_address VARCHAR(200) NOT NULL,
    short_address VARCHAR(100) NOT NULL
);

CREATE TABLE composers (
    composer_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    full_name VARCHAR(200) NOT NULL,
    short_name VARCHAR(200) NOT NULL
);

CREATE TABLE pieces (
    piece_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    piece_title VARCHAR(200) NOT NULL,
    composer_id INT NOT NULL REFERENCES composers(composer_id) ON DELETE CASCADE
);

CREATE TABLE programmes (
    programme_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    programme_title VARCHAR(200) NOT NULL
);

CREATE TABLE programme_pieces (
    programme_id INT NOT NULL REFERENCES programmes(programme_id) ON DELETE CASCADE,
    piece_id INT NOT NULL REFERENCES pieces(piece_id) ON DELETE CASCADE,
    sequence INT NOT NULL CHECK (sequence > 0),
    PRIMARY KEY (programme_id, piece_id),
    UNIQUE (programme_id, sequence)
);

-- Consider extending variants to include 'cancelled' and 'deleted'.
CREATE TYPE event_status AS ENUM ('draft', 'published', 'archived');

CREATE TABLE events (
    event_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    event_title VARCHAR(200) NOT NULL,
    event_date TIMESTAMP,
    ticket_link VARCHAR(500),
    venue_id INT REFERENCES venues(venue_id) ON DELETE CASCADE,
    programme_id INT REFERENCES programmes(programme_id) ON DELETE CASCADE,
    status event_status NOT NULL DEFAULT 'draft',
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create a trigger for updating the updated_at column.
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Attach the trigger to the events table.
CREATE TRIGGER update_events_updated_at
BEFORE UPDATE ON events
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

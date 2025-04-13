-- Create links table
CREATE TABLE links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    original_url TEXT NOT NULL,
    short_url TEXT NOT NULL UNIQUE,
    custom_slug TEXT UNIQUE,
    clicks INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE,
    customer_id UUID NOT NULL REFERENCES customer(id) ON DELETE CASCADE
);

-- Create index on short_url for faster lookups
CREATE INDEX idx_links_short_url ON links(short_url);

-- Create index on custom_slug for faster lookups
CREATE INDEX idx_links_custom_slug ON links(custom_slug);

-- Create index on customer_id for faster queries
CREATE INDEX idx_links_customer_id ON links(customer_id);

-- Create index on expires_at for faster queries
CREATE INDEX idx_links_expires_at ON links(expires_at);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_links_updated_at
    BEFORE UPDATE ON links
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column(); 
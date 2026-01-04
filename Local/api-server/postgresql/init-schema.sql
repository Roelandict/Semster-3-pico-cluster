-- Initialize foodchain_db schema for sensor data
-- Execute as sensor_admin user

-- Create currentTemperature table
CREATE TABLE IF NOT EXISTS public.currentTemperature (
    id BIGSERIAL PRIMARY KEY,
    sensor_id TEXT NOT NULL,
    temperature_avg NUMERIC(5, 2) NOT NULL,
    units TEXT DEFAULT 'Celsius',
    truck_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_currentTemperature_truck_id ON public.currentTemperature(truck_id);
CREATE INDEX IF NOT EXISTS idx_currentTemperature_sensor_id ON public.currentTemperature(sensor_id);
CREATE INDEX IF NOT EXISTS idx_currentTemperature_created_at ON public.currentTemperature(created_at);

-- Grant permissions to sensor_admin role
GRANT ALL PRIVILEGES ON TABLE public.currentTemperature TO sensor_admin;
GRANT ALL PRIVILEGES ON SEQUENCE public.currentTemperature_id_seq TO sensor_admin;

-- Allow web_anon role to select (read-only)
GRANT SELECT ON TABLE public.currentTemperature TO web_anon;

-- Enable Row Level Security (optional but recommended)
ALTER TABLE public.currentTemperature ENABLE ROW LEVEL SECURITY;

-- Create RLS policy for web_anon (can only read)
CREATE POLICY select_currenttemperature ON public.currentTemperature 
    FOR SELECT USING (true);

-- Create RLS policy for sensor_admin (can do all operations)
CREATE POLICY all_currenttemperature ON public.currentTemperature 
    FOR ALL USING (true);

-- Optional: Create a view for public access
CREATE OR REPLACE VIEW public.public_temperature_data AS
    SELECT sensor_id, temperature_avg, units, truck_id, created_at
    FROM public.currentTemperature
    ORDER BY created_at DESC
    LIMIT 100;

GRANT SELECT ON public.public_temperature_data TO web_anon;

-- Display table info
\d public.currentTemperature
\d+ public.currentTemperature

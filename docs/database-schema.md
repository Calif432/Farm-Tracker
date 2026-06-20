-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users & Authentication
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    full_name VARCHAR(255),
    role VARCHAR(50) NOT NULL DEFAULT 'worker' CHECK (role IN ('owner', 'manager', 'worker')),
    farm_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Farms (in case you expand to multiple farms later)
CREATE TABLE farms (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    location TEXT,
    size_hectares DECIMAL(10,2),
    owner_id UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Livestock
CREATE TABLE animals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tag_id VARCHAR(100) UNIQUE NOT NULL,
    type VARCHAR(50) NOT NULL, -- cow, sheep, goat, rabbit, chicken, etc.
    breed VARCHAR(100),
    gender VARCHAR(20) CHECK (gender IN ('male', 'female')),
    date_of_birth DATE NOT NULL,
    status VARCHAR(30) DEFAULT 'alive' CHECK (status IN ('alive', 'dead', 'sold', 'slaughtered')),
    parent_id UUID REFERENCES animals(id), -- for lineage
    mother_id UUID REFERENCES animals(id),
    father_id UUID REFERENCES animals(id),
    farm_id UUID REFERENCES farms(id),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE animal_produce (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    animal_id UUID REFERENCES animals(id),
    produce_type VARCHAR(50) NOT NULL, -- milk, eggs, wool, meat, etc.
    quantity DECIMAL(10,2) NOT NULL,
    unit VARCHAR(20) NOT NULL,
    recorded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE vaccinations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    animal_id UUID REFERENCES animals(id),
    vaccine_name VARCHAR(100) NOT NULL,
    administered_date DATE NOT NULL,
    next_due_date DATE,
    administered_by VARCHAR(100),
    notes TEXT
);

CREATE TABLE deaths (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    animal_id UUID REFERENCES animals(id) UNIQUE,
    death_date DATE NOT NULL,
    reason TEXT NOT NULL,
    prevention_notes TEXT,
    recorded_by UUID REFERENCES users(id)
);

-- Crops & Fields
CREATE TABLE fields (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    size_hectares DECIMAL(10,2),
    soil_type VARCHAR(50),
    location TEXT,
    farm_id UUID REFERENCES farms(id)
);

CREATE TABLE plantings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    field_id UUID REFERENCES fields(id),
    crop_type VARCHAR(50) NOT NULL, -- maize, beans, tomatoes, tea, coffee, etc.
    planting_date DATE NOT NULL,
    expected_harvest_date DATE,
    status VARCHAR(30) DEFAULT 'growing',
    quantity_planted DECIMAL(10,2),
    unit VARCHAR(20)
);

CREATE TABLE harvests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    planting_id UUID REFERENCES plantings(id),
    harvest_date DATE NOT NULL,
    quantity DECIMAL(10,2) NOT NULL,
    unit VARCHAR(20),
    quality_grade VARCHAR(20)
);

-- Inventory
CREATE TABLE inventory_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    category VARCHAR(50), -- feed, fertilizer, seed, pesticide, tool, etc.
    quantity DECIMAL(10,2) DEFAULT 0,
    unit VARCHAR(20),
    min_stock_level DECIMAL(10,2) DEFAULT 0,
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Finance
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type VARCHAR(20) NOT NULL CHECK (type IN ('income', 'expense')),
    amount DECIMAL(12,2) NOT NULL,
    category VARCHAR(100),
    description TEXT,
    related_to_type VARCHAR(50), -- animal, planting, inventory, etc.
    related_to_id UUID,
    transaction_date DATE NOT NULL,
    recorded_by UUID REFERENCES users(id)
);

-- Tasks
CREATE TABLE tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    due_date TIMESTAMP WITH TIME ZONE,
    status VARCHAR(30) DEFAULT 'pending',
    assigned_to UUID REFERENCES users(id),
    related_to_type VARCHAR(50),
    related_to_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_animals_tag ON animals(tag_id);
CREATE INDEX idx_animals_status ON animals(status);
CREATE INDEX idx_plantings_crop ON plantings(crop_type);
CREATE INDEX idx_transactions_date ON transactions(transaction_date);

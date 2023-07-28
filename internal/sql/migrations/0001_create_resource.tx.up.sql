CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY,
    create_time timestamptz NOT NULL,
	update_time timestamptz NOT NULL,
	delete_time timestamptz,
    version int8 NOT NULL DEFAULT 1,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    create_by VARCHAR(100)
);

CREATE FUNCTION increment_version() RETURNS TRIGGER AS $$
BEGIN
	NEW.version = OLD.version+1;
	RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_increment_version
	BEFORE UPDATE
	ON users
	FOR EACH ROW EXECUTE PROCEDURE increment_version();

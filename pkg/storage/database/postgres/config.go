package postgres

type PostgresConfig struct {
	ConnectionString string
	MigrationPath    string
}

func NewConfig() *PostgresConfig {
	return &PostgresConfig{
		MigrationPath: "pkg/storage/database/postgres/migrations", //TODO find a better way to pass the migrations path
	}
}

// SetConnectionString sets the connection string parameter on the Config object
func (c *PostgresConfig) SetConnectionString(connectionString string) *PostgresConfig {
	c.ConnectionString = connectionString

	return c
}

// SetMigrationPath sets the migration string parameter on the Config object
func (c *PostgresConfig) SetMigrationPath(migrationpath string) *PostgresConfig {
	c.MigrationPath = migrationpath

	return c
}
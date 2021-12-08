package pool

import "database/sql"

// Connection Database connection
type Connection struct {
	Pool   *DbPool
	DB     *sql.DB
	Status bool
}

// Close Returning a database connection to the connection pool
func (conn *Connection) Close() {
	conn.Pool.Close(conn)
}

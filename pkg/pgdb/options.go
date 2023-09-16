package pgdb

import "time"

type Option func(*PostgresDB)

func MaxPoolSize(size int) Option {
	return func(c *PostgresDB) {
		c.maxPoolSize = size
	}
}

func ConnAttempts(attempts int) Option {
	return func(c *PostgresDB) {
		c.connAttempts = attempts
	}
}

func ConnTimeout(timeout time.Duration) Option {
	return func(c *PostgresDB) {
		c.connTimeout = timeout
	}
}

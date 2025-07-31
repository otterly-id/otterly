package database

import "github.com/otterly-id/otterly/backend/app/queries"

type Queries struct {
	*queries.UserQueries
}

func OpenDBConnection() (*Queries, error) {
	db, err := PostgreSQLConnection()
	if err != nil {
		return nil, err
	}

	return &Queries{
		UserQueries: &queries.UserQueries{DB: db},
	}, nil
}

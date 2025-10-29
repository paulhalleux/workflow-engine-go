package utils

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func IsDuplicateKeyError(err error) bool {
	var pgErr *pgconn.PgError
	errors.As(err, &pgErr)
	return pgErr.Code == "23505"
}

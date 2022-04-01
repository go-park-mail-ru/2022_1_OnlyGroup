package postgres

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/pkg/errors"
)

func checkError(err *error, whatDo string, where string) error {
	var pgErr *pgconn.PgError
	if errors.As(*err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.ForeignKeyViolation:
			*err = fmt.Errorf(whatDo+" "+where+" failed: %s, %w", err, handlers.ErrBadRequest)
			return *err
		case pgerrcode.NoData:
			*err = fmt.Errorf(whatDo+" "+where+" failed: %s, %w", err, handlers.ErrBaseApp)
			return *err
		}
	}
	return nil
}

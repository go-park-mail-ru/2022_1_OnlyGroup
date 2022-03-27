package postgres

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
)

const maxEmailSize = "128"
const maxPasswordSize = "32"

type postgresUsersRepository struct {
	connection *pgx.Conn
	tableName  string
}

func CreatePostgresUsersRepo(conn *pgx.Conn, tableName string) (*postgresUsersRepository, error) {
	_, err := conn.Exec("CREATE TABLE IF NOT EXISTS " + tableName + "(id bigserial primary key, email varchar(" +
		maxEmailSize + ") unique, password varchar(" + maxPasswordSize + "));")
	if err != nil {
		return nil, err
	}

	return &postgresUsersRepository{connection: conn, tableName: tableName}, nil
}

func (repo *postgresUsersRepository) AddUser(email string, password string) (int, error) {
	var id int
	err := repo.connection.QueryRow("INSERT INTO "+
		repo.tableName+" (email, password) VALUES ($1, $2) RETURNING id;", email, password).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, handlers.ErrAuthEmailUsed
	}
	if err != nil {
		err = errors.Wrap(handlers.ErrBaseApp, err.Error())
		err = errors.Wrap(err, "add user failed")
		return 0, err
	}
	return id, nil
}

func (repo *postgresUsersRepository) Authorize(email string, password string) (int, error) {
	var id int
	err := repo.connection.QueryRow("SELECT id FROM "+
		repo.tableName+" WHERE email=$1 AND password=$2;", email, password).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, handlers.ErrAuthWrongPassword
	}
	if err != nil {
		err = errors.Wrap(handlers.ErrBaseApp, err.Error())
		err = errors.Wrap(err, "auth user failed")
		return 0, err
	}
	return id, nil
}

func (repo *postgresUsersRepository) ChangePassword(id int, newPassword string) error {
	tag, err := repo.connection.Exec("UPDATE "+repo.tableName+" SET password=$1 WHERE id=$2", newPassword, id)
	if err != nil {
		err = errors.Wrap(handlers.ErrBaseApp, err.Error())
		err = errors.Wrap(err, "change password user failed")
		return err
	}
	if tag.RowsAffected() == 0 {
		return handlers.ErrAuthUserNotFound
	}
	return nil
}

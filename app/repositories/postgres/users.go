package postgres

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type postgresUsersRepository struct {
	connection *sqlx.DB
	tableName  string
}

func NewPostgresUsersRepo(conn *sqlx.DB, tableName string) (*postgresUsersRepository, error) {
	_, err := conn.Exec("CREATE TABLE IF NOT EXISTS " + tableName + "(id bigserial primary key, email varchar(128) unique, password varchar(32));")
	if err != nil {
		return nil, fmt.Errorf("create table failed: %w", err)
	}

	return &postgresUsersRepository{connection: conn, tableName: tableName}, nil
}

func (repo *postgresUsersRepository) AddUser(email string, password string) (int, error) {
	var id int
	err := repo.connection.QueryRow("SELECT id FROM "+repo.tableName+" WHERE email=$1;", email).Scan(&id)
	if err == nil {
		return 0, handlers.ErrAuthEmailUsed
	}
	err = repo.connection.QueryRow("INSERT INTO "+repo.tableName+" (email, password) VALUES ($1, $2) RETURNING id;", email, password).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("add user failed: %s, %w", err.Error(), handlers.ErrBaseApp)
	}
	return id, nil
}

func (repo *postgresUsersRepository) Authorize(email string, password string) (int, error) {
	var id int
	var dbPassword string
	err := repo.connection.QueryRow("SELECT id, password FROM "+repo.tableName+" WHERE email=$1;", email).Scan(&id, &dbPassword)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, handlers.ErrAuthWrongPassword
	}
	if err != nil {
		return 0, fmt.Errorf("auth failed: %s, %w", err.Error(), handlers.ErrBaseApp)
	}
	if dbPassword != password {
		return 0, handlers.ErrAuthWrongPassword
	}
	return id, nil
}

func (repo *postgresUsersRepository) ChangePassword(id int, newPassword string) error {
	result, err := repo.connection.Exec("UPDATE "+repo.tableName+" SET password=$1 WHERE id=$2;", newPassword, id)
	if err != nil {
		return fmt.Errorf("changePassword failed: %s, %w", err.Error(), handlers.ErrBaseApp)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("changePassword failed: %s, %w", err.Error(), handlers.ErrBaseApp)
	}
	if affected == 0 {
		return handlers.ErrAuthUserNotFound
	}
	return nil
}

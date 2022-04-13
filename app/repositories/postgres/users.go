package postgres

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"2022_1_OnlyGroup_back/pkg/randomGenerator"
	"database/sql"
	"encoding/base32"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/argon2"
	"strings"
)

const defaultSaltSize = 16
const defaultArgon2Time = 1
const defaultArgon2Memory = 64 * 1024
const defaultArgon2Threads = 4
const defaultArgon2KeyLen = 32
const defaultSaltHashSeparator = "_"

type postgresUsersRepository struct {
	connection    *sqlx.DB
	saltGenerator randomGenerator.RandomGenerator
	tableName     string
}

func NewPostgresUsersRepo(conn *sqlx.DB, tableName string, saltGenerator randomGenerator.RandomGenerator) (*postgresUsersRepository, error) {
	_, err := conn.Exec("CREATE TABLE IF NOT EXISTS " + tableName + "(id bigserial primary key, email varchar(128) unique, password varchar(256));")
	if err != nil {
		return nil, fmt.Errorf("create table failed: %w", err)
	}

	return &postgresUsersRepository{connection: conn, tableName: tableName, saltGenerator: saltGenerator}, nil
}

func (repo *postgresUsersRepository) AddUser(email string, password string) (int, error) {
	var id int
	err := repo.connection.QueryRow("SELECT id FROM "+repo.tableName+" WHERE email=$1;", email).Scan(&id)
	if err == nil {
		return 0, handlers.ErrAuthEmailUsed
	}

	salt, err := repo.saltGenerator.String(defaultSaltSize)
	if err != nil {
		return 0, err
	}

	hashedPassword := base32.StdEncoding.EncodeToString(argon2.IDKey([]byte(password), []byte(salt), defaultArgon2Time, defaultArgon2Memory, defaultArgon2Threads, defaultArgon2KeyLen))
	dbPassword := strings.Join([]string{salt, string(hashedPassword)}, defaultSaltHashSeparator)

	err = repo.connection.QueryRow("INSERT INTO "+repo.tableName+" (email, password) VALUES ($1, $2) RETURNING id;", email, dbPassword).Scan(&id)

	if err != nil {
		return 0, handlers.ErrBaseApp.Wrap(err, "add user failed")
	}
	return id, nil
}

func (repo *postgresUsersRepository) Authorize(email string, password string) (int, error) {
	var id int
	var dbPassword string
	err := repo.connection.QueryRow("SELECT id, password FROM "+repo.tableName+" WHERE email=$1;", email).Scan(&id, &dbPassword)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, handlers.ErrAuthWrongPassword
	}
	if err != nil {
		return 0, handlers.ErrBaseApp.Wrap(err, "auth failed")
	}
	saltAndHashedPassword := strings.Split(dbPassword, defaultSaltHashSeparator)
	salt, err := base32.StdEncoding.DecodeString(saltAndHashedPassword[0])
	if err != nil {
		return 0, handlers.ErrBaseApp.Wrap(err, "decode salt from database failed")
	}
	passwordFromUser := base32.StdEncoding.EncodeToString(argon2.IDKey([]byte(password), salt, defaultArgon2Time, defaultArgon2Memory, defaultArgon2Threads, defaultArgon2KeyLen))

	if saltAndHashedPassword[1] != passwordFromUser {
		return 0, handlers.ErrAuthWrongPassword
	}
	return id, nil
}

func (repo *postgresUsersRepository) ChangePassword(id int, newPassword string) error {
	salt, err := repo.saltGenerator.String(defaultSaltSize)
	if err != nil {
		return err
	}
	hashedPassword := base32.StdEncoding.EncodeToString(argon2.IDKey([]byte(newPassword), []byte(salt), defaultArgon2Time, defaultArgon2Memory, defaultArgon2Threads, defaultArgon2KeyLen))
	dbPassword := strings.Join([]string{salt, hashedPassword}, defaultSaltHashSeparator)

	result, err := repo.connection.Exec("UPDATE "+repo.tableName+" SET password=$1 WHERE id=$2;", dbPassword, id)
	if err != nil {
		return handlers.ErrBaseApp.Wrap(err, "changePassword failed")
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return handlers.ErrBaseApp.Wrap(err, "changePassword failed")
	}
	if affected == 0 {
		return handlers.ErrAuthUserNotFound
	}
	return nil
}

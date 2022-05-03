package postgres

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"2022_1_OnlyGroup_back/app/models"
	"github.com/jmoiron/sqlx"
)

type InterestsPostgres struct {
	dataBase        *sqlx.DB
	staticInterests string
}

func NewInterestsPostgres(dataBase *sqlx.DB, staticInterests string) (*InterestsPostgres, error) {
	_, err := dataBase.Exec("CREATE TABLE IF NOT EXISTS " + staticInterests + "(" +
		"id     bigserial,\n" +
		"title     	varchar(32));")
	if err != nil {
		return nil, handlers.ErrBaseApp.Wrap(err, "create table failed")
	}
	return &InterestsPostgres{dataBase: dataBase, staticInterests: staticInterests}, err
}

func (repo *InterestsPostgres) GetInterests() ([]models.Interest, error) {
	var interests []models.Interest
	err := repo.dataBase.Select(&interests, "SELECT id, title FROM "+repo.staticInterests)
	if err != nil {
		return nil, handlers.ErrBaseApp.Wrap(err, "get interests failed")
	}

	return interests, nil
}

func (repo *InterestsPostgres) CheckInterests(interests []models.Interest) error {
	var findStatus bool
	for _, val := range interests {
		err := repo.dataBase.Select(&findStatus, "select exists(select * from "+repo.staticInterests+" where id = $1);", val.Id)
		if err != nil {
			return handlers.ErrBaseApp.Wrap(err, "failed check interests")
		}
		if !findStatus {
			return handlers.ErrBadRequest
		}
	}
	return nil
}

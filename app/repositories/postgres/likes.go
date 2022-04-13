package postgres

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"2022_1_OnlyGroup_back/app/models"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
)

type LikesPostgres struct {
	dataBase       *sqlx.DB
	tableNameLikes string
	tableNameUsers string
}

func NewLikesPostgres(dataBase *sqlx.DB, tableNameLikes string, tableNameUsers string) (*LikesPostgres, error) {
	_, err := dataBase.Exec("CREATE TABLE IF NOT EXISTS " + tableNameLikes + "(" +
		"who     bigserial references " + tableNameUsers + "(id),\n" +
		"whom     bigserial references " + tableNameUsers + "(id),\n" +
		"action     numeric default -1);")
	if err != nil {
		return nil, handlers.ErrBaseApp.Wrap(err, "get shortProfile failed")
	}

	return &LikesPostgres{dataBase: dataBase, tableNameLikes: tableNameLikes, tableNameUsers: tableNameUsers}, nil
}

func (repo *LikesPostgres) SetAction(profileId int, likes models.Likes) (err error) {
	_, err = repo.dataBase.Exec("DELETE FROM "+repo.tableNameLikes+" WHERE (who=$1 and whom=$2)", profileId, likes.Id)
	if err != nil {
		return handlers.ErrBaseApp.Wrap(err, "get shortProfile failed")
	}
	_, err = repo.dataBase.Exec("INSERT INTO "+repo.tableNameLikes+" (who, whom, action) VALUES ($1, $2, $3)", profileId, likes.Id, likes.Action)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.ForeignKeyViolation:
				return handlers.ErrBadRequest
			case pgerrcode.NoData:
				return handlers.ErrBaseApp
			}
		}
	}
	return
}

func (repo *LikesPostgres) GetMatched(profileId int) (likesVector models.LikesMatched, err error) {
	err = repo.dataBase.Select(&likesVector.VectorId, "select l1.whom from "+repo.tableNameLikes+" as l1 join "+repo.tableNameLikes+" as l2 on l1.whom = l2.who and l1.action=1 where l1.who=l2.whom and l2.action=1 and l1.who=$1", profileId)
	if err != nil {
		return likesVector, handlers.ErrBaseApp.Wrap(err, "get shortProfile failed")
	}
	return
}

package postgres

import (
	"2022_1_OnlyGroup_back/app/handlers/http"
	"2022_1_OnlyGroup_back/app/models"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type ProfilePostgres struct {
	dataBase             *sqlx.DB
	tableNameProfiles    string
	tableUserInterests   string
	tableStaticInterests string
	tableNameFilters     string
	tableNameLikes       string
}

const sizeVectorCandidates = 3
const defaultInterest = "golang programming"

func NewProfilePostgres(dataBase *sqlx.DB, tableNameProfile string, tableNameUsers string, tableNameInterests string, tableStaticInterests string, tableFilters string, tableLikes string) (*ProfilePostgres, error) {
	_, err := dataBase.Exec("CREATE TABLE IF NOT EXISTS " + tableNameProfile + "(" +
		"UserId     bigserial unique,\n" +
		"FirstName  varchar(32) default '',\n" +
		"LastName   text default '',\n" +
		"Birthday   timestamp default now(),\n" +
		"City       varchar(32) default '',\n" +
		"AboutUser  text default '',\n" +
		"Height     numeric default 0,\n" +
		"Gender     numeric default -1 );")
	if err != nil {
		return nil, http.ErrBaseApp.Wrap(err, "create table failed")
	}

	_, err = dataBase.Exec("CREATE TABLE IF NOT EXISTS " + tableStaticInterests + "(" +
		"id bigserial unique primary key,\n" +
		"title varchar(32));")
	if err != nil {
		return nil, http.ErrBaseApp.Wrap(err, "create table failed")
	}

	_, err = dataBase.Exec("CREATE TABLE IF NOT EXISTS " + tableNameInterests + "(" +
		"UserId bigserial references " + tableNameProfile + "(UserId)  ON DELETE CASCADE,\n" +
		"Id bigserial references " + tableStaticInterests + "(id) ON DELETE CASCADE);")
	if err != nil {
		return nil, http.ErrBaseApp.Wrap(err, "create table failed")
	}
	_, err = dataBase.Exec("INSERT INTO "+tableStaticInterests+"(title) VALUES ($1)", defaultInterest)
	if err != nil {
		return nil, http.ErrBaseApp.Wrap(err, "insert empty profile failed")
	}

	_, err = dataBase.Exec("CREATE TABLE IF NOT EXISTS " + tableFilters + "(" +
		"UserId bigserial unique references " + tableNameProfile + "(UserId) ON DELETE CASCADE,\n" +
		"BottomHeightFilter	numeric default 0,\n" +
		"TopHeightFilter	numeric default 0,\n" +
		"GenderFilter   numeric default 0,\n" +
		"BottomAgeFilter	numeric default 0,\n" +
		"TopAgeFilter numeric default 0);")
	if err != nil {
		return nil, http.ErrBaseApp.Wrap(err, "create table failed")
	}
	_, err = dataBase.Exec("CREATE TABLE IF NOT EXISTS " + tableLikes + "(" +
		"who     bigserial references " + tableNameProfile + "(UserId)  ON DELETE CASCADE,\n" +
		"whom     bigserial references " + tableNameProfile + "(UserId)  ON DELETE CASCADE,\n" +
		"action     numeric default -1);")
	if err != nil {
		return nil, http.ErrBaseApp.Wrap(err, "get shortProfile failed")
	}

	return &ProfilePostgres{dataBase, tableNameProfile, tableNameInterests, tableStaticInterests, tableFilters, tableLikes}, nil
}
func (repo *ProfilePostgres) Get(profileId int) (profile models.Profile, err error) {
	err = repo.dataBase.QueryRowx("SELECT firstname, lastname, birthday, city, aboutuser, userid, height,gender FROM "+repo.tableNameProfiles+" WHERE userid=$1", profileId).StructScan(&profile)
	if err != nil {
		return profile, http.ErrBaseApp.Wrap(err, "get profile failed")
	}

	var interests []models.Interest
	err = repo.dataBase.Select(&interests, "select l2.id, l2.title from "+repo.tableUserInterests+" as l1 join "+repo.tableStaticInterests+" as l2 on l1.id = l2.id where userid = $1;", profileId)
	if err != nil {
		return profile, http.ErrBaseApp.Wrap(err, "get interests failed")
	}
	profile.Interests = interests
	return
}

func (repo *ProfilePostgres) GetShort(profileId int) (shortProfile models.ShortProfile, err error) {
	err = repo.dataBase.QueryRowx("SELECT firstName, lastname, city FROM "+repo.tableNameProfiles+" WHERE userid=$1", profileId).StructScan(&shortProfile)
	if err != nil {
		return shortProfile, http.ErrBaseApp.Wrap(err, "get shortProfile failed")
	}
	return
}

func (repo *ProfilePostgres) Change(profileId int, profile models.Profile) (err error) {
	_, err = repo.dataBase.NamedExec("UPDATE "+repo.tableNameProfiles+" SET firstname=:firstname, lastname=:lastname, birthday=:birthday, city=:city, aboutuser=:aboutuser, gender=:gender, height=:height WHERE userid = :userid", profile)
	if err != nil {
		return http.ErrBaseApp.Wrap(err, "change profile failed")
	}
	_, err = repo.dataBase.NamedExec("DELETE FROM "+repo.tableUserInterests+" WHERE userid = :userid", profile)
	if err != nil {
		return http.ErrBaseApp.Wrap(err, "delete interests failed")
	}

	for _, val := range profile.Interests {
		_, err = repo.dataBase.Exec("INSERT INTO "+repo.tableUserInterests+" (UserId, id) VALUES ($1, $2)", profile.UserId, val.Id)
		if err != nil {
			return http.ErrBaseApp.Wrap(err, "change interests failed")
		}
	}
	return
}

func (repo *ProfilePostgres) Delete(profileId int) (err error) {
	_, err = repo.dataBase.Exec("DELETE FROM "+repo.tableNameProfiles+" WHERE userid = $1", profileId)
	if err != nil {
		return http.ErrBaseApp.Wrap(err, "delete profile failed")
	}
	return
}

func (repo *ProfilePostgres) Add(profile models.Profile) (err error) {
	_, err = repo.dataBase.NamedExec("INSERT INTO "+repo.tableNameProfiles+" (firstname, lastname, birthday, city, aboutuser, userid, gender, height) VALUES (:firstname, :lastname, :birthday, :city, :aboutuser, :userid, :gender, :height)", profile)
	if err != nil {
		return http.ErrBaseApp.Wrap(err, "insert profile failed")
	}
	for _, val := range profile.Interests {
		_, err = repo.dataBase.Exec("INSERT INTO "+repo.tableUserInterests+" (UserId, Id) VALUES ($1, $2)", profile.UserId, val.Id)
		if err != nil {
			return http.ErrBaseApp.Wrap(err, "insert interests failed")
		}
	}
	return
}

func (repo *ProfilePostgres) AddEmpty(profileId int) (err error) {
	_, err = repo.dataBase.Exec("INSERT INTO "+repo.tableNameProfiles+"(userid) VALUES ($1)", profileId)
	if err != nil {
		return http.ErrBaseApp.Wrap(err, "insert empty profile failed")
	}
	_, err = repo.dataBase.Exec("INSERT INTO "+repo.tableNameFilters+"(userid) VALUES ($1)", profileId)
	if err != nil {
		return http.ErrBaseApp.Wrap(err, "insert empty filters failed")
	}
	return
}

func (repo *ProfilePostgres) FindCandidate(profileId int) (candidateProfiles models.VectorCandidate, err error) {
	var profilesId []int
	err = repo.dataBase.Select(&profilesId, "SELECT userid FROM "+repo.tableNameProfiles+" WHERE userid !=$1 ORDER BY random() LIMIT 3", profileId)
	if err != nil {
		return candidateProfiles, http.ErrBaseApp.Wrap(err, "get profiles failed")
	}
	candidateProfiles.Candidates = make([]int, sizeVectorCandidates)
	for idx, val := range profilesId {
		candidateProfiles.Candidates[idx] = val
	}
	return
}

func (repo *ProfilePostgres) CheckFiled(profileId int) (err error) {
	var profile models.Profile
	err = repo.dataBase.QueryRowx("SELECT firstname, lastname, birthday, city, aboutuser, userid, gender, height FROM "+repo.tableNameProfiles+" WHERE userid=$1 LIMIT 3 ", profileId).StructScan(&profile)
	if err != nil {
		return http.ErrBaseApp.Wrap(err, "get profile failed")
	}
	var interests []models.Interest
	err = repo.dataBase.Select(&interests, "SELECT interests FROM "+repo.tableUserInterests+" WHERE userid=$1", profileId)
	if err != nil {
		return http.ErrBaseApp.Wrap(err, "get interests failed")
	}
	profile.Interests = interests
	if profile.Gender == -1 ||
		profile.City == "" ||
		profile.LastName == "" ||
		profile.AboutUser == "" ||
		//profile.Birthday == "" ||
		len(profile.Interests) == 0 ||
		profile.FirstName == "" {
		return http.ErrProfileNotFiled
	}
	return
}

func (repo *ProfilePostgres) GetFilters(userId int) (models.Filters, error) {
	var filters models.Filters
	err := repo.dataBase.QueryRow("SELECT BottomHeightFilter, TopHeightFilter, GenderFilter, BottomAgeFilter, TopAgeFilter from "+repo.tableNameFilters+" where userid=$1", userId).Scan(&filters.HeightFilter[0], &filters.HeightFilter[1], &filters.GenderFilter, &filters.AgeFilter[0], &filters.AgeFilter[1])

	if err != nil {
		return models.Filters{}, http.ErrBaseApp.Wrap(err, "failed get filters")
	}
	return filters, nil
}

func (repo *ProfilePostgres) ChangeFilters(userId int, filters models.Filters) error {
	_, err := repo.dataBase.Exec("UPDATE "+repo.tableNameFilters+" SET BottomHeightFilter=$1, TopHeightFilter=$2, GenderFilter=$3, BottomAgeFilter=$4, TopAgeFilter=$5 WHERE userid=$6", filters.HeightFilter[0], filters.HeightFilter[1], filters.GenderFilter, filters.AgeFilter[0], filters.AgeFilter[1], userId)
	if err != nil {
		return http.ErrBaseApp.Wrap(err, "change profile failed")
	}
	return nil
}

func (repo *ProfilePostgres) GetInterests() ([]models.Interest, error) {
	var interests []models.Interest
	err := repo.dataBase.Select(&interests, "SELECT id, title FROM "+repo.tableStaticInterests)
	if err != nil {
		return nil, http.ErrBaseApp.Wrap(err, "get interests failed")
	}

	return interests, nil
}

func (repo *ProfilePostgres) GetDynamicInterest(interest string) ([]models.Interest, error) {
	var interests []models.Interest
	err := repo.dataBase.Select(&interests, "select * from "+repo.tableStaticInterests+" where title ILIKE $1;", "%"+interest+"%")
	if err != nil {
		return nil, http.ErrBaseApp.Wrap(err, "get interests failed")
	}
	return interests, nil
}

func (repo *ProfilePostgres) CheckInterests(interests []models.Interest) error {
	var findStatus []bool
	for _, val := range interests {
		err := repo.dataBase.Select(&findStatus, "select exists(select * from "+repo.tableStaticInterests+" where id = $1);", val.Id)
		if err != nil {
			return http.ErrBaseApp.Wrap(err, "failed check interests")
		}
		if !findStatus[0] {
			return http.ErrBadRequest
		}
	}
	return nil
}

func (repo *ProfilePostgres) SetAction(profileId int, likes models.Likes) (err error) {
	_, err = repo.dataBase.Exec("DELETE FROM "+repo.tableNameLikes+" WHERE (who=$1 and whom=$2)", profileId, likes.Id)
	if err != nil {
		return http.ErrBaseApp.Wrap(err, "get shortProfile failed")
	}
	_, err = repo.dataBase.Exec("INSERT INTO "+repo.tableNameLikes+" (who, whom, action) VALUES ($1, $2, $3)", profileId, likes.Id, likes.Action)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.ForeignKeyViolation:
				return http.ErrBadRequest
			case pgerrcode.NoData:
				return http.ErrBaseApp
			}
		}
	}
	return
}

func (repo *ProfilePostgres) GetMatched(profileId int) (likesVector models.LikesMatched, err error) {
	err = repo.dataBase.Select(&likesVector.VectorId, "select l1.whom from "+repo.tableNameLikes+" as l1 join "+repo.tableNameLikes+" as l2 on l1.whom = l2.who and l1.action=1 where l1.who=l2.whom and l2.action=1 and l1.who=$1", profileId)
	if err != nil {
		return likesVector, http.ErrBaseApp.Wrap(err, "get shortProfile failed")
	}
	return
}

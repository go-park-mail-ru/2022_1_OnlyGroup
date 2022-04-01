package postgres

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"2022_1_OnlyGroup_back/app/models"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type ProfilePostgres struct {
	dataBase           *sqlx.DB
	tableNameProfiles  string
	tableNameUsers     string
	tableNameInterests string
}

const sizeVectorCandidates = 3

func NewProfilePostgres(dataBase *sqlx.DB, tableNameProfile string, tableNameUsers string, tableNameInterests string) (*ProfilePostgres, error) {
	_, err := dataBase.Exec("CREATE TABLE IF NOT EXISTS " + tableNameProfile + "(" +
		"UserId     bigserial unique references " + tableNameUsers + "(id),\n" +
		"FirstName   varchar(32) default '',\n" +
		"LastName   text default '',\n" +
		"Birthday   varchar(32) default '',\n" +
		"City       varchar(32) default '',\n" +
		"AboutUser  text default '',\n" +
		"Height     numeric default 0,\n" +
		"Gender     varchar(32) default '');")

	if err != nil {
		return nil, fmt.Errorf("create table failed: %s, %w", err.Error(), handlers.ErrBaseApp)
	}

	_, err = dataBase.Exec("CREATE TABLE IF NOT EXISTS " + tableNameInterests + "(" +
		"UserId bigserial references " + tableNameUsers + "(id)," +
		"Interests varchar(32) default '');")

	if err != nil {
		return nil, fmt.Errorf("create table failed: %s, %w", err.Error(), handlers.ErrBaseApp)
	}

	return &ProfilePostgres{dataBase, tableNameProfile, tableNameUsers, tableNameInterests}, nil
}

func (repo *ProfilePostgres) Get(profileId int) (profile models.Profile, err error) {
	err = repo.dataBase.QueryRowx("SELECT firstname, lastname, birthday, city, aboutuser, userid, height,gender FROM "+repo.tableNameProfiles+" WHERE userid=$1", profileId).StructScan(&profile)
	if err != nil {
		return profile, fmt.Errorf("get profile failed: %s, %w", err.Error(), handlers.ErrBaseApp)
	}

	var interests []string
	err = repo.dataBase.Select(&interests, "SELECT interests FROM "+repo.tableNameInterests+" WHERE userid=$1", profileId)
	if err != nil {
		return profile, fmt.Errorf("get interests failed: %s, %w", err.Error(), handlers.ErrBaseApp)
	}
	profile.Interests = interests
	return
}

func (repo *ProfilePostgres) GetShort(profileId int) (shortProfile models.ShortProfile, err error) {
	err = repo.dataBase.QueryRowx("SELECT firstName, lastname, city FROM "+repo.tableNameProfiles+" WHERE userid=$1", profileId).StructScan(&shortProfile)
	if err != nil {
		return shortProfile, fmt.Errorf("get shortProfile failed: %s, %w", err.Error(), handlers.ErrBaseApp)
	}
	return
}

func (repo *ProfilePostgres) Change(profileId int, profile models.Profile) (err error) {
	_, err = repo.dataBase.NamedExec("UPDATE "+repo.tableNameProfiles+" SET firstname=:firstname, lastname=:lastname, birthday=:birthday, city=:city, aboutuser=:aboutuser, gender=:gender, height=:height WHERE userid = :userid", profile)
	if err != nil {
		return fmt.Errorf("create table failed: %s, %w", err, handlers.ErrBaseApp)
	}
	_, err = repo.dataBase.Exec("DELETE FROM "+repo.tableNameInterests+" WHERE userid = $1", profileId)
	if err != nil {
		return fmt.Errorf("delete interests failed: %s, %w", err.Error(), handlers.ErrBaseApp)
	}

	for _, val := range profile.Interests {
		_, err = repo.dataBase.Exec("INSERT INTO "+repo.tableNameInterests+" (userid, interests) VALUES ($1, $2)", profile.UserId, val)
		if err != nil {
			return fmt.Errorf("change interests failed: %s, %w", err.Error(), handlers.ErrBaseApp)
		}
	}
	return
}

func (repo *ProfilePostgres) Delete(profileId int) (err error) {
	_, err = repo.dataBase.Exec("DELETE FROM "+repo.tableNameProfiles+" WHERE userid = $1", profileId)
	if err != nil {
		fmt.Errorf("delete profile failed: %s, %w", err.Error(), handlers.ErrBaseApp)
		return err
	}

	_, err = repo.dataBase.Exec("DELETE FROM "+repo.tableNameInterests+" WHERE userid = $1", profileId)
	if err != nil {
		return fmt.Errorf("delete interests failed: %s, %w", err.Error(), handlers.ErrBaseApp)
	}
	return
}

func (repo *ProfilePostgres) Add(profile models.Profile) (err error) {
	_, err = repo.dataBase.NamedExec("INSERT INTO "+repo.tableNameProfiles+" (firstname, lastname, birthday, city, aboutuser, userid, gender, height) VALUES (:firstname, :lastname, :birthday, :city, :aboutuser, :userid, :gender, :height)", profile)
	if err != nil {
		return fmt.Errorf("insert profile failed: %s, %w", err.Error(), handlers.ErrBaseApp)
	}
	for _, val := range profile.Interests {
		_, err = repo.dataBase.Exec("INSERT INTO "+repo.tableNameInterests+" (userid, interests) VALUES ($1, $2)", profile.UserId, val)
		if err != nil {
			return fmt.Errorf("insert interests failed: %s, %w", err.Error(), handlers.ErrBaseApp)
		}
	}
	return
}

func (repo *ProfilePostgres) AddEmpty(profileId int) (err error) {
	_, err = repo.dataBase.Exec("INSERT INTO "+repo.tableNameProfiles+"(userid) VALUES ($1)", profileId)
	if err != nil {
		return fmt.Errorf("insert empty profile failed: %s, %w", err.Error(), handlers.ErrBaseApp)
	}
	return
}

func (repo *ProfilePostgres) FindCandidate(profileId int) (candidateProfiles models.VectorCandidate, err error) {
	var profilesId []int
	err = repo.dataBase.Select(&profilesId, "SELECT userid FROM "+repo.tableNameProfiles+" WHERE userid !=$1 ORDER BY random() LIMIT 3", profileId)
	if err != nil {
		return candidateProfiles, fmt.Errorf("get profiles failed: %s, %w", err.Error(), handlers.ErrBaseApp)
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
		return fmt.Errorf("get profile failed: %s, %w", err.Error(), handlers.ErrBaseApp)
	}
	var interests []string
	err = repo.dataBase.Select(&interests, "SELECT interests FROM "+repo.tableNameInterests+" WHERE userid=$1", profileId)
	if err != nil {
		return fmt.Errorf("get interests failed: %s, %w", err, handlers.ErrBaseApp)
	}
	profile.Interests = interests
	if profile.Gender == "" ||
		profile.City == "" ||
		profile.LastName == "" ||
		profile.AboutUser == "" ||
		profile.Birthday == "" ||
		len(profile.Interests) == 0 ||
		profile.FirstName == "" {
		return handlers.ErrProfileNotFiled
	}
	return
}

package postgres

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"2022_1_OnlyGroup_back/app/models"
	"fmt"
	"github.com/jmoiron/sqlx"
)

const subTableName = "interests"

type ProfilePostgres struct {
	dataBase          *sqlx.DB
	tableNameProfiles string
	tableNameUsers    string
}

var count int

func NewProfilePostgres(dataBase *sqlx.DB, tableNameProfile string, tableNameUsers string) (*ProfilePostgres, error) {
	_, err := dataBase.Exec("CREATE TABLE IF NOT EXISTS " + tableNameProfile + "(" +
		"UserId     bigserial unique references " + tableNameUsers + "(id),\n" +
		"FirstName   varchar(32) default '',\n" +
		"LastName   text default '',\n" +
		"Birthday   varchar(32) default '',\n" +
		"City       varchar(32) default '',\n" +
		"AboutUser  text default '',\n" +
		"Gender     varchar(32) default '');")

	if err != nil {
		return nil, fmt.Errorf("create table failed: %s, %w", err.Error(), handlers.ErrBaseApp)
	}

	_, err = dataBase.Exec("CREATE TABLE IF NOT EXISTS " + subTableName + "(" +
		"UserId bigserial references " + tableNameUsers + "(id)," +
		"Interests varchar(32) default '');")

	if err != nil {
		return nil, fmt.Errorf("create table failed: %s, %w", err.Error(), handlers.ErrBaseApp)
	}

	return &ProfilePostgres{dataBase, tableNameProfile, tableNameUsers}, nil
}

func (repo *ProfilePostgres) GetProfile(profileId int) (profile models.Profile, err error) {
	err = repo.dataBase.QueryRowx("SELECT firstname, lastname, birthday, city, aboutuser, userid, gender FROM "+repo.tableNameProfiles+" WHERE userid=$1", profileId).StructScan(&profile)
	if err != nil {
		return profile, fmt.Errorf("get profile failed: %s, %w", err.Error(), handlers.ErrBaseApp)
	}

	var interests []string
	err = repo.dataBase.Select(&interests, "SELECT interests FROM "+subTableName+" WHERE userid=$1", profileId)
	if err != nil {
		return profile, fmt.Errorf("get interests failed: %s, %w", err.Error(), handlers.ErrBaseApp)
	}
	profile.Interests = interests
	return
}

func (repo *ProfilePostgres) GetShortProfile(profileId int) (shortProfile models.ShortProfile, err error) {
	err = repo.dataBase.QueryRowx("SELECT firstName, lastname, city FROM "+repo.tableNameProfiles+" WHERE userid=$1", profileId).StructScan(&shortProfile)
	if err != nil {
		return shortProfile, fmt.Errorf("get shortProfile failed: %s, %w", err.Error(), handlers.ErrBaseApp)
	}
	return
}

func (repo *ProfilePostgres) ChangeProfile(profileId int, profile models.Profile) (err error) {
	_, err = repo.dataBase.NamedExec("UPDATE "+repo.tableNameProfiles+" SET firstname=:firstname, lastname=:lastname, birthday=:birthday, city=:city, aboutuser=:aboutuser, gender=:gender WHERE userid = :userid", profile)
	if err != nil {
		return fmt.Errorf("create table failed: %s, %w", err, handlers.ErrBaseApp)
	}
	_, err = repo.dataBase.Exec("DELETE FROM "+subTableName+" WHERE userid = $1", profileId)
	if err != nil {
		return fmt.Errorf("delete interests failed: %s, %w", err.Error(), handlers.ErrBaseApp)
	}

	for _, val := range profile.Interests {
		_, err = repo.dataBase.Exec("INSERT INTO "+subTableName+" (userid, interests) VALUES ($1, $2)", profile.UserId, val)
		if err != nil {
			return fmt.Errorf("change interests failed: %s, %w", err.Error(), handlers.ErrBaseApp)
		}
	}
	return
}

func (repo *ProfilePostgres) DeleteProfile(profileId int) (err error) {
	_, err = repo.dataBase.Exec("DELETE FROM "+repo.tableNameProfiles+" WHERE userid = $1", profileId)
	if err != nil {
		fmt.Errorf("delete profile failed: %s, %w", err.Error(), handlers.ErrBaseApp)
		return err
	}

	_, err = repo.dataBase.Exec("DELETE FROM "+subTableName+" WHERE userid = $1", profileId)
	if err != nil {
		return fmt.Errorf("delete interests failed: %s, %w", err.Error(), handlers.ErrBaseApp)
	}
	return
}

func (repo *ProfilePostgres) AddProfile(profile models.Profile) (err error) {
	_, err = repo.dataBase.Exec("INSERT INTO "+repo.tableNameProfiles+" (firstname, lastname, birthday, city, aboutuser, userid, gender) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		profile.FirstName, profile.LastName, profile.Birthday, profile.City, profile.AboutUser, profile.UserId, profile.Gender)
	if err != nil {
		return fmt.Errorf("insert profile failed: %s, %w", err.Error(), handlers.ErrBaseApp)
	}
	for _, val := range profile.Interests {
		_, err = repo.dataBase.Exec("INSERT INTO "+subTableName+" (userid, interests) VALUES ($1, $2)", profile.UserId, val)
		if err != nil {
			return fmt.Errorf("insert interests failed: %s, %w", err.Error(), handlers.ErrBaseApp)
		}
	}
	return
}

func (repo *ProfilePostgres) AddEmptyProfile(profileId int) (err error) {
	_, err = repo.dataBase.Exec("INSERT INTO "+repo.tableNameProfiles+"(userid) VALUES ($1)", profileId)
	if err != nil {
		return fmt.Errorf("insert empty profile failed: %s, %w", err.Error(), handlers.ErrBaseApp)
	}
	return
}

func (repo *ProfilePostgres) FindCandidateProfile(profileId int) (candidateProfiles models.VectorCandidate, err error) {
	var profile []models.Profile
	err = repo.dataBase.Select(&profile, "SELECT * FROM "+repo.tableNameProfiles)
	if err != nil {
		return candidateProfiles, fmt.Errorf("get profiles failed: %s, %w", err.Error(), handlers.ErrBaseApp)
	}
	for i := 0; i < 3; i++ {
		if count == len(profile) {
			count = 0
		}
		candidateProfiles.Candidates = append(candidateProfiles.Candidates, profile[count].UserId)
		count += 1
	}
	return
}

func (repo *ProfilePostgres) CheckProfileFiled(profileId int) (err error) {
	var profile models.Profile
	err = repo.dataBase.QueryRowx("SELECT firstname, lastname, birthday, city, aboutuser, userid, gender FROM "+repo.tableNameProfiles+" WHERE userid=$1", profileId).StructScan(&profile)
	if err != nil {
		return fmt.Errorf("get profile failed: %s, %w", err.Error(), handlers.ErrBaseApp)
	}
	var interests []string
	err = repo.dataBase.Select(&interests, "SELECT interests FROM "+subTableName+" WHERE userid=$1", profileId)
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

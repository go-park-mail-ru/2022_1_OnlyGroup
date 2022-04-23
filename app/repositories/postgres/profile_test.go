package postgres

import (
	"2022_1_OnlyGroup_back/app/models"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zhashkevych/go-sqlxmock"
	"testing"
)

type creatingProfileRepoSuite struct {
	suite.Suite
	db     *sqlx.DB
	dbMock sqlmock.Sqlmock
}

func (suite *creatingProfileRepoSuite) SetupTest() {
	var err error
	suite.db, suite.dbMock, err = sqlmock.Newx(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		panic("creating mock failed")
	}
}
func (suite *creatingProfileRepoSuite) AfterTest() {
	suite.db.Close()
}

func (suite *creatingProfileRepoSuite) TestOk() {
	tableNameProfile := "test_name"
	tableNameUsers := "test_users"
	tableNameInterests := "test_interests"
	suite.dbMock.ExpectExec("CREATE TABLE IF NOT EXISTS " + tableNameProfile + "(" +
		"UserId     bigserial unique references " + tableNameUsers + "(id),\n" +
		"FirstName   varchar(32) default '',\n" +
		"LastName   text default '',\n" +
		"Birthday   timestamp,\n" +
		"City       varchar(32) default '',\n" +
		"AboutUser  text default '',\n" +
		"Height     numeric default 0,\n" +
		"Gender     numeric default -1 );").WillReturnResult(sqlmock.NewResult(0, 0))
	suite.dbMock.ExpectExec("CREATE TABLE IF NOT EXISTS " + tableNameInterests + "(" +
		"UserId bigserial references " + tableNameUsers + "(id)," +
		"Interests varchar(32) default '');").WillReturnResult(sqlmock.NewResult(0, 0))

	_, err := NewProfilePostgres(suite.db, tableNameProfile, tableNameUsers, tableNameInterests)

	if !errors.Is(err, nil) {
		suite.T().Errorf("Wrapped error mismatched, expected: '%v', got '%v'", nil, err)
	}
}

func (suite *creatingProfileRepoSuite) TestDBProfileError() {
	tableNameProfile := "test_name"
	tableNameUsers := "test_users"
	tableNameInterests := "test_interests"
	suite.dbMock.ExpectExec("CREATE TABLE IF NOT EXISTS " + tableNameProfile + "(" +
		"UserId     bigserial unique references " + tableNameUsers + "(id),\n" +
		"FirstName  varchar(32) default '',\n" +
		"LastName   text default '',\n" +
		"Birthday   timestamp,\n" +
		"City       varchar(32) default '',\n" +
		"AboutUser  text default '',\n" +
		"Height     numeric default 0,\n" +
		"Gender     numeric default -1 );").WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(postgresError)
	_, err := NewProfilePostgres(suite.db, tableNameProfile, tableNameUsers, tableNameInterests)
	if !errors.Is(err, postgresError) {
		suite.T().Errorf("Wrapped error mismatched, expected: '%v', got '%v'", postgresError, err)
	}
}

func (suite *creatingProfileRepoSuite) TestDBInterestsError() {
	tableNameProfile := "test_name"
	tableNameUsers := "test_users"
	tableNameInterests := "test_interests"
	suite.dbMock.ExpectExec("CREATE TABLE IF NOT EXISTS " + tableNameProfile + "(" +
		"UserId     bigserial unique references " + tableNameUsers + "(id),\n" +
		"FirstName  varchar(32) default '',\n" +
		"LastName   text default '',\n" +
		"Birthday   timestamp,\n" +
		"City       varchar(32) default '',\n" +
		"AboutUser  text default '',\n" +
		"Height     numeric default 0,\n" +
		"Gender     numeric default -1 );").WillReturnResult(sqlmock.NewResult(0, 0))
	suite.dbMock.ExpectExec("CREATE TABLE IF NOT EXISTS " + tableNameInterests + "(" +
		"UserId bigserial references " + tableNameUsers + "(id)," +
		"Interests varchar(32) default '');").WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(postgresError)

	_, err := NewProfilePostgres(suite.db, tableNameProfile, tableNameUsers, tableNameInterests)
	if !errors.Is(err, postgresError) {
		suite.T().Errorf("Wrapped error mismatched, expected: '%v', got '%v'", postgresError, err)
	}
}

func TestProfileCreatingRepo(t *testing.T) {
	suite.Run(t, new(creatingProfileRepoSuite))
}

func TestGetTableDriven(t *testing.T) {
	var tests = []struct {
		testName         string
		mockPrepare      func(mock *sqlmock.Sqlmock)
		expectedError    error
		testMasterId     int
		profileModel     models.Profile
		testProfileTable string
	}{
		{
			"All ok",
			func(mock *sqlmock.Sqlmock) {
				(*mock).ExpectQuery("SELECT firstname, lastname, birthday, city, aboutuser, userid, height,gender FROM test_profile_table WHERE userid=$1").WithArgs(0).WillReturnRows(sqlmock.NewRows([]string{"firstname", "lastname", "birthday", "city", "aboutuser", "userid", "height", "gender"}).AddRow("testFirstName", "testLastName", "1998-01-1", "testCity", "testAboutUser", 8, 189, 1))
				(*mock).ExpectQuery("SELECT interests FROM test_interests_table WHERE userid=$1").WithArgs(0).WillReturnRows(sqlmock.NewRows([]string{"interests"}).AddRow("football").AddRow("Music"))

			},
			nil,
			0,
			models.Profile{FirstName: "testFirstName", LastName: "testLastName", Birthday: "1998-01-1", City: "testCity", AboutUser: "testAboutUser", UserId: 8, Height: 189, Gender: 1, Interests: []string{"football", "Music"}},
			"test_profile_table",
		},
		{
			"Postgres base err profile",
			func(mock *sqlmock.Sqlmock) {
				(*mock).ExpectQuery("SELECT firstname, lastname, birthday, city, aboutuser, userid, height,gender FROM test_profile_table WHERE userid=$1").WithArgs(0).WillReturnError(postgresError)
				(*mock).ExpectQuery("SELECT interests FROM test_interests_table WHERE userid=$1").WithArgs(0).WillReturnRows(sqlmock.NewRows([]string{"interests"}).AddRow("football").AddRow("Music"))
			},
			postgresError,
			0,
			models.Profile{},
			"test_profile_table",
		},
		{
			"Postgres base err interests",
			func(mock *sqlmock.Sqlmock) {
				(*mock).ExpectQuery("SELECT firstname, lastname, birthday, city, aboutuser, userid, height,gender FROM test_profile_table WHERE userid=$1").WithArgs(0).WillReturnRows(sqlmock.NewRows([]string{"firstname", "lastname", "birthday", "city", "aboutuser", "userid", "height", "gender"}).AddRow("testFirstName", "testLastName", "1998-01-1", "testCity", "testAboutUser", 8, 189, 1))
				(*mock).ExpectQuery("SELECT interests FROM test_interests_table WHERE userid=$1").WithArgs(0).WillReturnError(postgresError)

			},
			postgresError,
			0,
			models.Profile{},
			"test_profile_table",
		},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			db, dbMock, err := sqlmock.Newx(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			defer db.Close()
			if err != nil {
				panic("creating mock failed")
			}
			testingRepo := ProfilePostgres{db, "test_profile_table", "test_user_table", "test_interests_table"}
			test.mockPrepare(&dbMock)
			testModel, err := testingRepo.Get(test.testMasterId)

			if !errors.Is(err, test.expectedError) {
				t.Errorf("Wrapped error mismatched, expected: '%v', got '%v'", test.expectedError, err)
			}
			if err == nil {
				assert.Equal(t, test.profileModel.LastName, testModel.LastName)
				assert.Equal(t, test.profileModel.FirstName, testModel.FirstName)
				assert.Equal(t, test.profileModel.City, testModel.City)
				assert.Equal(t, test.profileModel.AboutUser, testModel.AboutUser)
				assert.Equal(t, test.profileModel.Gender, testModel.Gender)
				assert.Equal(t, test.profileModel.UserId, testModel.UserId)
				for idx, val := range test.profileModel.Interests {
					assert.Equal(t, val, testModel.Interests[idx])
				}
			}
		})
	}
}

func TestGetShortTableDriven(t *testing.T) {
	var tests = []struct {
		testName          string
		mockPrepare       func(mock *sqlmock.Sqlmock)
		expectedError     error
		testMasterId      int
		profileModel      models.ShortProfile
		testProfileTable  string
		testInterestTable string
		testUserTable     string
	}{
		{
			"All ok",
			func(mock *sqlmock.Sqlmock) {
				(*mock).ExpectQuery("SELECT firstName, lastname, city FROM test_profile_table WHERE userid=$1").WithArgs(0).WillReturnRows(sqlmock.NewRows([]string{"firstname", "lastname", "city"}).AddRow("testFirstName", "testLastName", "testCity"))
			},
			nil,
			0,
			models.ShortProfile{FirstName: "testFirstName", LastName: "testLastName", City: "testCity"},
			"test_profile_table",
			"test_user_table",
			"test_interests_table",
		},
		{
			"Postgres base err",
			func(mock *sqlmock.Sqlmock) {
				(*mock).ExpectQuery("SELECT firstName, lastname, city FROM test_profile_table WHERE userid=$1").WithArgs(0).WillReturnError(postgresError)
			},
			postgresError,
			0,
			models.ShortProfile{},
			"test_profile_table",
			"test_user_table",
			"test_interests_table",
		},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			db, dbMock, err := sqlmock.Newx(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			defer db.Close()
			if err != nil {
				panic("creating mock failed")
			}
			testingRepo := ProfilePostgres{db, test.testProfileTable, test.testUserTable, test.testInterestTable}
			test.mockPrepare(&dbMock)
			testModel, err := testingRepo.GetShort(test.testMasterId)

			if !errors.Is(err, test.expectedError) {
				t.Errorf("Wrapped error mismatched, expected: '%v', got '%v'", test.expectedError, err)
			}
			assert.Equal(t, test.profileModel.LastName, testModel.LastName)
			assert.Equal(t, test.profileModel.FirstName, testModel.FirstName)
			assert.Equal(t, test.profileModel.City, testModel.City)

		})
	}
}

func TestChangeTableDriven(t *testing.T) {
	var tests = []struct {
		testName          string
		mockPrepare       func(mock *sqlmock.Sqlmock)
		expectedError     error
		testMasterId      int
		profileModel      models.Profile
		testProfileTable  string
		testInterestTable string
		testUserTable     string
	}{
		{
			"All ok",
			func(mock *sqlmock.Sqlmock) {
				(*mock).ExpectExec("UPDATE test_profile_table SET firstname=?, lastname=?, birthday=?, city=?, aboutuser=?, gender=?, height=? WHERE userid=?").WithArgs("testFirstName", "testLastName", "1998-01-1", "testCity", "testAboutUser", 1, 189, 8).WillReturnResult(sqlmock.NewResult(0, 1))
				(*mock).ExpectExec("DELETE FROM test_interests_table WHERE userid=?").WithArgs(8).WillReturnResult(sqlmock.NewResult(0, 2))
				(*mock).ExpectExec("INSERT INTO test_interests_table (userid, interests) VALUES ($1, $2)").WithArgs(8, "football").WillReturnResult(sqlmock.NewResult(1, 0))
				(*mock).ExpectExec("INSERT INTO test_interests_table (userid, interests) VALUES ($1, $2)").WithArgs(8, "Music").WillReturnResult(sqlmock.NewResult(2, 0))

			},
			nil,
			0,
			models.Profile{FirstName: "testFirstName", LastName: "testLastName", Birthday: "1998-01-1", City: "testCity", AboutUser: "testAboutUser", UserId: 8, Height: 189, Gender: 1, Interests: []string{"football", "Music"}},
			"test_profile_table",
			"test_interests_table",
			"test_user_table",
		},
		{
			"Postgres base err profile",
			func(mock *sqlmock.Sqlmock) {
				(*mock).ExpectExec("UPDATE test_profile_table SET firstname=?, lastname=?, birthday=?, city=?, aboutuser=?, gender=?, height=? WHERE userid=?").WithArgs("testFirstName", "testLastName", "1998-01-1", "testCity", "testAboutUser", 1, 189, 8).WillReturnError(postgresError)
				(*mock).ExpectExec("DELETE FROM test_interests_table WHERE userid=?").WithArgs(8).WillReturnResult(sqlmock.NewResult(0, 2))
				(*mock).ExpectExec("INSERT INTO test_interests_table (userid, interests) VALUES ($1, $2)").WithArgs(8, "football").WillReturnResult(sqlmock.NewResult(1, 0))
				(*mock).ExpectExec("INSERT INTO test_interests_table (userid, interests) VALUES ($1, $2)").WithArgs(8, "Music").WillReturnResult(sqlmock.NewResult(2, 0))

			},
			postgresError,
			0,
			models.Profile{FirstName: "testFirstName", LastName: "testLastName", Birthday: "1998-01-1", City: "testCity", AboutUser: "testAboutUser", UserId: 8, Height: 189, Gender: 1, Interests: []string{"football", "Music"}},
			"test_profile_table",
			"test_interests_table",
			"test_user_table",
		},
		{
			"Postgres base err interests delete",
			func(mock *sqlmock.Sqlmock) {
				(*mock).ExpectExec("UPDATE test_profile_table SET firstname=?, lastname=?, birthday=?, city=?, aboutuser=?, gender=?, height=? WHERE userid=?").WithArgs("testFirstName", "testLastName", "1998-01-1", "testCity", "testAboutUser", 1, 189, 8).WillReturnResult(sqlmock.NewErrorResult(postgresError))
				(*mock).ExpectExec("DELETE FROM test_interests_table WHERE userid=?").WithArgs(8).WillReturnError(postgresError)
				(*mock).ExpectExec("INSERT INTO test_interests_table (userid, interests) VALUES ($1, $2)").WithArgs(8, "football").WillReturnResult(sqlmock.NewResult(1, 0))
				(*mock).ExpectExec("INSERT INTO test_interests_table (userid, interests) VALUES ($1, $2)").WithArgs(8, "Music").WillReturnResult(sqlmock.NewResult(2, 0))
			},
			postgresError,
			0,
			models.Profile{FirstName: "testFirstName", LastName: "testLastName", Birthday: "1998-01-1", City: "testCity", AboutUser: "testAboutUser", UserId: 8, Height: 189, Gender: 1, Interests: []string{"football", "Music"}},
			"test_profile_table",
			"test_interests_table",
			"test_user_table",
		},
		{
			"Postgres base err interests insert",
			func(mock *sqlmock.Sqlmock) {
				(*mock).ExpectExec("UPDATE test_profile_table SET firstname=?, lastname=?, birthday=?, city=?, aboutuser=?, gender=?, height=? WHERE userid=?").WithArgs("testFirstName", "testLastName", "1998-01-1", "testCity", "testAboutUser", 1, 189, 8).WillReturnResult(sqlmock.NewErrorResult(postgresError))
				(*mock).ExpectExec("DELETE FROM test_interests_table WHERE userid=?").WithArgs(8).WillReturnResult(sqlmock.NewResult(0, 2))
				(*mock).ExpectExec("INSERT INTO test_interests_table (userid, interests) VALUES ($1, $2)").WithArgs(8, "football").WillReturnError(postgresError)
				(*mock).ExpectExec("INSERT INTO test_interests_table (userid, interests) VALUES ($1, $2)").WithArgs(8, "Music").WillReturnError(postgresError)
			},
			postgresError,
			0,
			models.Profile{FirstName: "testFirstName", LastName: "testLastName", Birthday: "1998-01-1", City: "testCity", AboutUser: "testAboutUser", UserId: 8, Height: 189, Gender: 1, Interests: []string{"football", "Music"}},
			"test_profile_table",
			"test_interests_table",
			"test_user_table",
		},
		{
			"Postgres base err interests insert",
			func(mock *sqlmock.Sqlmock) {
				(*mock).ExpectExec("UPDATE test_profile_table SET firstname=?, lastname=?, birthday=?, city=?, aboutuser=?, gender=?, height=? WHERE userid=?").WithArgs("testFirstName", "testLastName", "1998-01-1", "testCity", "testAboutUser", 1, 189, 8).WillReturnResult(sqlmock.NewErrorResult(postgresError))
				(*mock).ExpectExec("DELETE FROM test_interests_table WHERE userid=?").WithArgs(8).WillReturnResult(sqlmock.NewResult(0, 2))
				(*mock).ExpectExec("INSERT INTO test_interests_table (userid, interests) VALUES ($1, $2)").WithArgs(8, "football").WillReturnResult(sqlmock.NewResult(1, 0))
				(*mock).ExpectExec("INSERT INTO test_interests_table (userid, interests) VALUES ($1, $2)").WithArgs(8, "Music").WillReturnError(postgresError)
			},
			postgresError,
			0,
			models.Profile{FirstName: "testFirstName", LastName: "testLastName", Birthday: "1998-01-1", City: "testCity", AboutUser: "testAboutUser", UserId: 8, Height: 189, Gender: 1, Interests: []string{"football", "Music"}},
			"test_profile_table",
			"test_interests_table",
			"test_user_table",
		},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			db, dbMock, err := sqlmock.Newx(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			defer db.Close()
			if err != nil {
				panic("creating mock failed")
			}
			testingRepo := ProfilePostgres{db, test.testProfileTable, test.testUserTable, test.testInterestTable}
			test.mockPrepare(&dbMock)
			err = testingRepo.Change(test.profileModel.UserId, test.profileModel)

			if !errors.Is(err, test.expectedError) {
				t.Errorf("Wrapped error mismatched, expected: '%v', got '%v'", test.expectedError, err)
			}
		})
	}
}

func TestDeleteTableDriven(t *testing.T) {
	var tests = []struct {
		testName          string
		mockPrepare       func(mock *sqlmock.Sqlmock)
		expectedError     error
		testMasterId      int
		testProfileTable  string
		testInterestTable string
		testUserTable     string
	}{
		{
			"All ok",
			func(mock *sqlmock.Sqlmock) {
				(*mock).ExpectExec("DELETE FROM test_profile_table WHERE userid = $1").WithArgs(8).WillReturnResult(sqlmock.NewResult(0, 1))
				(*mock).ExpectExec("DELETE FROM test_interests_table WHERE userid = $1").WithArgs(8).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			nil,
			8,
			"test_profile_table",
			"test_interests_table",
			"test_user_table",
		},
		{
			"Postgres base err profile",
			func(mock *sqlmock.Sqlmock) {
				(*mock).ExpectExec("DELETE FROM test_profile_table WHERE userid = $1").WithArgs(8).WillReturnError(postgresError)
				(*mock).ExpectExec("DELETE FROM test_interests_table WHERE userid = $1").WithArgs(8).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			postgresError,
			8,
			"test_profile_table",
			"test_interests_table",
			"test_user_table",
		},
		{
			"Postgres base err interests delete",
			func(mock *sqlmock.Sqlmock) {
				(*mock).ExpectExec("DELETE FROM test_profile_table WHERE userid = $1").WithArgs(8).WillReturnResult(sqlmock.NewResult(0, 1))
				(*mock).ExpectExec("DELETE FROM test_interests_table WHERE userid = $1").WithArgs(8).WillReturnError(postgresError)
			},
			postgresError,
			8,
			"test_profile_table",
			"test_interests_table",
			"test_user_table",
		},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			db, dbMock, err := sqlmock.Newx(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			defer db.Close()
			if err != nil {
				panic("creating mock failed")
			}
			testingRepo := ProfilePostgres{db, test.testProfileTable, test.testUserTable, test.testInterestTable}
			test.mockPrepare(&dbMock)
			err = testingRepo.Delete(test.testMasterId)

			if !errors.Is(err, test.expectedError) {
				t.Errorf("Wrapped error mismatched, expected: '%v', got '%v'", test.expectedError, err)
			}
		})
	}
}

func TestAddEmptyTableDriven(t *testing.T) {
	var tests = []struct {
		testName          string
		mockPrepare       func(mock *sqlmock.Sqlmock)
		expectedError     error
		testMasterId      int
		testProfileTable  string
		testInterestTable string
		testUserTable     string
	}{
		{
			"All ok",
			func(mock *sqlmock.Sqlmock) {
				(*mock).ExpectExec("INSERT INTO test_profile_table (userid) VALUES ($1)").WithArgs(8).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			nil,
			8,
			"test_profile_table",
			"test_interests_table",
			"test_user_table",
		},
		{
			"Postgres base err",
			func(mock *sqlmock.Sqlmock) {
				(*mock).ExpectExec("INSERT INTO test_profile_table (userid) VALUES ($1)").WithArgs(8).WillReturnError(postgresError)
			},
			postgresError,
			8,
			"test_profile_table",
			"test_interests_table",
			"test_user_table",
		},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			db, dbMock, err := sqlmock.Newx(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			defer db.Close()
			if err != nil {
				panic("creating mock failed")
			}
			testingRepo := ProfilePostgres{db, test.testProfileTable, test.testUserTable, test.testInterestTable}
			test.mockPrepare(&dbMock)
			err = testingRepo.AddEmpty(test.testMasterId)

			if !errors.Is(err, test.expectedError) {
				t.Errorf("Wrapped error mismatched, expected: '%v', got '%v'", test.expectedError, err)
			}
		})
	}
}

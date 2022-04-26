package postgres

import (
	"2022_1_OnlyGroup_back/app/handlers"
	randomGenerator "2022_1_OnlyGroup_back/pkg/randomGenerator/impl"
	"database/sql"
	"github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zhashkevych/go-sqlxmock"
	"testing"
)

var postgresError = errors.New("test postgres err")

type creatingRepoTestSuite struct {
	suite.Suite
	db     *sqlx.DB
	dbMock sqlmock.Sqlmock
}

func (suite *creatingRepoTestSuite) SetupTest() {
	var err error
	suite.db, suite.dbMock, err = sqlmock.Newx(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		panic("creating mock failed")
	}
}

func (suite *creatingRepoTestSuite) AfterTest() {
	suite.db.Close()
}

func (suite *creatingRepoTestSuite) TestOk() {
	suite.dbMock.ExpectExec("CREATE TABLE IF NOT EXISTS test_name(id bigserial primary key, email varchar(128) unique, password varchar(128));").WillReturnResult(sqlmock.NewResult(0, 0))
	_, err := NewPostgresUsersRepo(suite.db, "test_name", randomGenerator.NewCryptoRandomGenerator())
	if !errors.Is(err, nil) {
		suite.T().Errorf("Wrapped error mismatched, expected: '%v', got '%v'", nil, err)
	}
}

func (suite *creatingRepoTestSuite) TestDBError() {
	suite.dbMock.ExpectExec("CREATE TABLE IF NOT EXISTS test_name(id bigserial primary key, email varchar(128) unique, password varchar(128));").WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(postgresError)
	_, err := NewPostgresUsersRepo(suite.db, "test_name", randomGenerator.NewCryptoRandomGenerator())
	if !errors.Is(err, postgresError) {
		suite.T().Errorf("Wrapped error mismatched, expected: '%v', got '%v'", postgresError, err)
	}
}

func TestCreatingRepo(t *testing.T) {
	suite.Run(t, new(creatingRepoTestSuite))
}

func TestAddUserTableDriven(t *testing.T) {
	var tests = []struct {
		testName      string
		mockPrepare   func(mock *sqlmock.Sqlmock)
		expectedError error
		testUserId    int
		testEmail     string
		testPass      string
		testUserTable string
	}{
		{
			"All ok",
			func(mock *sqlmock.Sqlmock) {
				(*mock).ExpectQuery("SELECT id FROM test_user_table WHERE email=$1;").WithArgs("test@mail.ru").WillReturnError(pgx.ErrNoRows)
				(*mock).ExpectQuery("INSERT INTO test_user_table (email, password) VALUES ($1, $2) RETURNING id;").WithArgs("test@mail.ru", "TestPass3").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("4"))
			},
			nil,
			4,
			"test@mail.ru",
			"TestPass3",
			"test_user_table",
		},
		{
			"Email already registred",
			func(mock *sqlmock.Sqlmock) {
				(*mock).ExpectQuery("SELECT id FROM test_user_table WHERE email=$1;").WithArgs("test@mail.ru").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("4"))
			},
			handlers.ErrAuthEmailUsed,
			0,
			"test@mail.ru",
			"TestPass3",
			"test_user_table",
		},
		{
			"Postgres internal error",
			func(mock *sqlmock.Sqlmock) {
				(*mock).ExpectQuery("SELECT id FROM test_user_table WHERE email=$1;").WithArgs("test@mail.ru").WillReturnError(postgresError)
			},
			handlers.ErrBaseApp,
			0,
			"test@mail.ru",
			"TestPass3",
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
			testingRepo := postgresUsersRepository{db, randomGenerator.NewCryptoRandomGenerator(), test.testUserTable}
			test.mockPrepare(&dbMock)
			userIdReturned, err := testingRepo.AddUser(test.testEmail, test.testPass)

			if !errors.Is(err, test.expectedError) {
				t.Errorf("Wrapped error mismatched, expected: '%v', got '%v'", test.expectedError, err)
			}
			if test.expectedError == nil {
				assert.Equal(t, test.testUserId, userIdReturned)
			}
		})
	}
}

func TestAuthTableDriven(t *testing.T) {
	var tests = []struct {
		testName      string
		mockPrepare   func(mock *sqlmock.Sqlmock)
		expectedError error
		testUserId    int
		testEmail     string
		testPass      string
		testUserTable string
	}{
		{
			"All ok",
			func(mock *sqlmock.Sqlmock) {
				(*mock).ExpectQuery("SELECT id, password FROM test_user_table WHERE email=$1;").WithArgs("test@mail.ru").WillReturnRows(sqlmock.NewRows([]string{"id", "password"}).AddRow("4", "TestPass3"))
			},
			nil,
			4,
			"test@mail.ru",
			"TestPass3",
			"test_user_table",
		},
		{
			"User not registered",
			func(mock *sqlmock.Sqlmock) {
				(*mock).ExpectQuery("SELECT id, password FROM test_user_table WHERE email=$1;").WithArgs("test@mail.ru").WillReturnError(sql.ErrNoRows)
			},
			handlers.ErrAuthWrongPassword,
			0,
			"test@mail.ru",
			"TestPass3",
			"test_user_table",
		},
		{
			"Wrong password",
			func(mock *sqlmock.Sqlmock) {
				(*mock).ExpectQuery("SELECT id, password FROM test_user_table WHERE email=$1;").WithArgs("test@mail.ru").WillReturnRows(sqlmock.NewRows([]string{"id", "password"}).AddRow("4", "TestPasfwefws3"))
			},
			handlers.ErrAuthWrongPassword,
			0,
			"test@mail.ru",
			"TestPass3",
			"test_user_table",
		},
		{
			"Postgres internal error",
			func(mock *sqlmock.Sqlmock) {
				(*mock).ExpectQuery("SELECT id, password FROM test_user_table WHERE email=$1;").WithArgs("test@mail.ru").WillReturnError(postgresError)
			},
			handlers.ErrBaseApp,
			0,
			"test@mail.ru",
			"TestPass3",
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
			testingRepo := postgresUsersRepository{db, randomGenerator.NewCryptoRandomGenerator(), test.testUserTable}
			test.mockPrepare(&dbMock)
			userIdReturned, err := testingRepo.Authorize(test.testEmail, test.testPass)

			if !errors.Is(err, test.expectedError) {
				t.Errorf("Wrapped error mismatched, expected: '%v', got '%v'", test.expectedError, err)
			}
			if test.expectedError == nil {
				assert.Equal(t, test.testUserId, userIdReturned)
			}
		})
	}
}

func TestChangePswdTableDriven(t *testing.T) {
	var tests = []struct {
		testName      string
		mockPrepare   func(mock *sqlmock.Sqlmock)
		expectedError error
		testUserId    int
		testNewPass   string
		testUserTable string
	}{
		{
			"All ok",
			func(mock *sqlmock.Sqlmock) {
				(*mock).ExpectExec("UPDATE test_user_table SET password=$1 WHERE id=$2;").WithArgs("TestNewPass1", 4).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			nil,
			4,
			"TestNewPass1",
			"test_user_table",
		},
		{
			"User not found",
			func(mock *sqlmock.Sqlmock) {
				(*mock).ExpectExec("UPDATE test_user_table SET password=$1 WHERE id=$2;").WithArgs("TestNewPass1", 4).WillReturnResult(sqlmock.NewResult(0, 0))
			},
			handlers.ErrAuthUserNotFound,
			4,
			"TestNewPass1",
			"test_user_table",
		},
		{
			"Postgres internal error",
			func(mock *sqlmock.Sqlmock) {
				(*mock).ExpectExec("UPDATE test_user_table SET password=$1 WHERE id=$2;").WithArgs("TestNewPass1", 4).WillReturnError(postgresError)
			},
			handlers.ErrBaseApp,
			0,
			"TestNewPass1",
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
			testingRepo := postgresUsersRepository{db, randomGenerator.NewCryptoRandomGenerator(), test.testUserTable}
			test.mockPrepare(&dbMock)
			err = testingRepo.ChangePassword(test.testUserId, test.testNewPass)

			if !errors.Is(err, test.expectedError) {
				t.Errorf("Wrapped error mismatched, expected: '%v', got '%v'", test.expectedError, err)
			}
		})
	}
}

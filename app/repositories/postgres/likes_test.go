package postgres

import (
	"2022_1_OnlyGroup_back/app/handlers"
	"2022_1_OnlyGroup_back/app/models"
	"github.com/jackc/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"github.com/zhashkevych/go-sqlxmock"
	"testing"
)

var ForeignKeyErr = pgconn.PgError{Code: "23503", Severity: "ERROR", Message: "insert or update on table \"test_likes_table\" violates foreign key constraint \"test_likes_table_firstid_fkey\"",
	Detail: "Key (firstid)=(8) is not present in table \"os_users\".", Position: 0, InternalPosition: 0, Line: 2463}

type creatingLikesRepoSuite struct {
	suite.Suite
	db     *sqlx.DB
	dbMock sqlmock.Sqlmock
}

func (suite *creatingLikesRepoSuite) SetupTest() {
	var err error
	suite.db, suite.dbMock, err = sqlmock.Newx(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		panic("creating mock failed")
	}
}
func (suite *creatingLikesRepoSuite) AfterTest() {
	suite.db.Close()
}

func (suite *creatingLikesRepoSuite) TestOk() {
	tableNameLikes := "test_name"
	tableNameUsers := "test_users"

	suite.dbMock.ExpectExec("CREATE TABLE IF NOT EXISTS " + tableNameLikes + "(" +
		"firstId     bigserial references " + tableNameUsers + "(id),\n" +
		"lastId     bigserial references " + tableNameUsers + "(id),\n" +
		"action     numeric default -1);").WillReturnResult(sqlmock.NewResult(0, 0))
	_, err := NewLikesPostgres(suite.db, tableNameLikes, tableNameUsers)
	if !errors.Is(err, nil) {
		suite.T().Errorf("Wrapped error mismatched, expected: '%v', got '%v'", nil, err)
	}
}

func (suite *creatingLikesRepoSuite) TestDBError() {
	tableNameLikes := "test_name"
	tableNameUsers := "test_users"
	suite.dbMock.ExpectExec("CREATE TABLE IF NOT EXISTS " + tableNameLikes + "(" +
		"firstId     bigserial references " + tableNameUsers + "(id),\n" +
		"lastId     bigserial references " + tableNameUsers + "(id),\n" +
		"action     numeric default -1);").WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(postgresError)
	_, err := NewLikesPostgres(suite.db, tableNameLikes, tableNameUsers)
	if !errors.Is(err, postgresError) {
		suite.T().Errorf("Wrapped error mismatched, expected: '%v', got '%v'", postgresError, err)
	}
}

func TestLikesCreatingRepo(t *testing.T) {
	suite.Run(t, new(creatingLikesRepoSuite))
}

func TestSetActionTableDriven(t *testing.T) {
	var tests = []struct {
		testName       string
		mockPrepare    func(mock *sqlmock.Sqlmock)
		expectedError  error
		testMasterId   int
		testSlaveId    int
		testAction     int
		testLikesTable string
	}{
		{
			"All ok",
			func(mock *sqlmock.Sqlmock) {
				(*mock).ExpectExec("DELETE FROM test_likes_table WHERE (firstId=$1 and lastId=$2)").WithArgs(0, 1).WillReturnResult(sqlmock.NewErrorResult(nil))
				(*mock).ExpectExec("INSERT INTO test_likes_table (firstId, lastId, action) VALUES ($1, $2, $3)").WithArgs(0, 1, 1).WillReturnResult(sqlmock.NewErrorResult(nil))
			},
			nil,
			0,
			1,
			1,

			"test_likes_table",
		},
		{
			"Postgres internal err",
			func(mock *sqlmock.Sqlmock) {
				(*mock).ExpectExec("DELETE FROM test_likes_table WHERE (firstId=$1 and lastId=$2)").WithArgs(0, 1).WillReturnError(postgresError)
			},
			handlers.ErrBaseApp,
			0,
			1,
			1,
			"test_likes_table",
		},
		{
			"Postgres bad Request",
			func(mock *sqlmock.Sqlmock) {
				(*mock).ExpectExec("DELETE FROM test_likes_table WHERE (firstId=$1 and lastId=$2)").WithArgs(0, 2).WillReturnResult(sqlmock.NewErrorResult(nil))
				(*mock).ExpectExec("INSERT INTO test_likes_table (firstId, lastId, action) VALUES ($1, $2, $3)").WithArgs(0, 2, 1).WillReturnError(errors.New(ForeignKeyErr.Error()))

			},
			handlers.ErrBadRequest,
			0,
			2,
			1,
			"test_likes_table",
		},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			db, dbMock, err := sqlmock.Newx(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			defer db.Close()
			if err != nil {
				panic("creating mock failed")
			}
			testingRepo := LikesPostgres{db, "test_likes_table", "test_user_table"}
			test.mockPrepare(&dbMock)
			err = testingRepo.SetAction(test.testMasterId, models.Likes{Id: test.testSlaveId, Action: test.testAction})

			if !errors.Is(err, test.expectedError) {
				t.Errorf("Wrapped error mismatched, expected: '%v', got '%v'", test.expectedError, err)
			}
		})
	}
}

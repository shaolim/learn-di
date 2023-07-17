package data

import (
	"database/sql"
	"errors"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSave_getDBError(t *testing.T) {
	defer func(original func() (*sql.DB, error)) {
		getDB = original
	}(getDB)

	getDB = func() (*sql.DB, error) {
		return nil, errors.New("getDB() failed")
	}

	in := &Person{
		FullName: "Jake Blues",
		Phone:    "0123456789",
		Currency: "AUD",
		Price:    123.45,
	}

	resultID, err := Save(in)
	require.Error(t, err)
	assert.Equal(t, 0, resultID)
}

func TestSave(t *testing.T) {
	in := &Person{
		FullName: "Jake Blues",
		Phone:    "0123456789",
		Currency: "AUD",
		Price:    123.45,
	}

	scenarios := map[string]struct {
		configureMockDB func(sqlmock.Sqlmock)
		want            int
		wantErr         bool
	}{
		"happy path": {
			configureMockDB: func(dbMock sqlmock.Sqlmock) {
				queryRegex := convertSQLToRegex(sqlInsert)
				dbMock.ExpectExec(queryRegex).WillReturnResult(sqlmock.NewResult(2, 1))
			},
			want:    2,
			wantErr: false,
		},
		"insert error": {
			configureMockDB: func(dbMock sqlmock.Sqlmock) {
				queryRegex := convertSQLToRegex(sqlInsert)
				dbMock.ExpectExec(queryRegex).WillReturnError(errors.New("failed to insert"))
			},
			wantErr: true,
		},
	}

	for scenario, tt := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			defer func(original *sql.DB) {
				db = original
			}(db)

			testDb, dbMock, err := sqlmock.New()
			require.NoError(t, err)
			defer testDb.Close()

			tt.configureMockDB(dbMock)
			db = testDb

			// call function
			resultID, err := Save(in)

			assert.Equal(t, tt.want, resultID)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.NoError(t, dbMock.ExpectationsWereMet())
		})
	}
}

func TestLoadAll(t *testing.T) {
	scenarios := map[string]struct {
		configureMockDB func(sqlmock.Sqlmock)
		want            []*Person
		wantErr         bool
	}{
		"happy path": {
			configureMockDB: func(dbMock sqlmock.Sqlmock) {
				queryRegex := convertSQLToRegex(sqlLoadAll)
				dbMock.ExpectQuery(queryRegex).WillReturnRows(
					sqlmock.NewRows(strings.Split(sqlAllColumns, ", ")).
						AddRow(1, "John", "0123456789", "AUD", 12.34))
			},
			want: []*Person{
				{
					ID:       1,
					FullName: "John",
					Phone:    "0123456789",
					Currency: "AUD",
					Price:    12.34,
				},
			},
			wantErr: false,
		},
		"load error": {
			configureMockDB: func(dbMock sqlmock.Sqlmock) {
				queryRegex := convertSQLToRegex(sqlLoadAll)
				dbMock.ExpectQuery(queryRegex).
					WillReturnError(errors.New("something failed"))
			},
			wantErr: true,
		},
	}

	for scenario, tt := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			defer func(original *sql.DB) {
				db = original
			}(db)

			testDb, dbMock, err := sqlmock.New()
			require.NoError(t, err)
			defer testDb.Close()

			tt.configureMockDB(dbMock)
			db = testDb

			result, err := LoadAll()

			assert.NoError(t, dbMock.ExpectationsWereMet())
			assert.Equal(t, tt.want, result)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestLoad(t *testing.T) {
	scenarios := map[string]struct {
		configureMockDB func(sqlmock.Sqlmock)
		want            *Person
		wantErr         bool
	}{
		"happy path": {
			configureMockDB: func(dbMock sqlmock.Sqlmock) {
				queryRegex := convertSQLToRegex(sqlLoadAll)
				dbMock.ExpectQuery(queryRegex).
					WillReturnRows(sqlmock.NewRows(strings.Split(sqlAllColumns, ", ")).
						AddRow(2, "Paul", "0123456789", "CAD", 23.45))
			},
			want: &Person{
				ID:       2,
				FullName: "Paul",
				Phone:    "0123456789",
				Currency: "CAD",
				Price:    23.45,
			},
			wantErr: false,
		},
		"load error": {
			configureMockDB: func(dbMock sqlmock.Sqlmock) {
				queryRegex := convertSQLToRegex(sqlLoadAll)
				dbMock.ExpectQuery(queryRegex).WillReturnError(errors.New("something failed"))
			},
			wantErr: true,
		},
	}

	for scenario, tt := range scenarios {
		t.Run(scenario, func(t *testing.T) {
			defer func(original *sql.DB) {
				db = original
			}(db)

			testDb, dbMock, err := sqlmock.New()
			require.NoError(t, err)
			defer testDb.Close()

			tt.configureMockDB(dbMock)
			db = testDb

			result, err := Load(2)
			assert.Equal(t, tt.want, result)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.NoError(t, dbMock.ExpectationsWereMet())

		})
	}
}

func convertSQLToRegex(in string) string {
	return `\Q` + in + `\E`
}

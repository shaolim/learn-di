package data

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const (
	sqlAllColumns = "id, fullname, phone, currency, price"
	sqlInsert     = "INSERT INTO person(fullname, phone, currency, price) VALUES (?, ?, ?, ?)"
	sqlLoadAll    = "SELECT " + sqlAllColumns + " FROM person"
	sqlLoadByID   = "SELECT " + sqlAllColumns + " FROM person WHERE id = ? LIMIT 1"
)

var (
	db *sql.DB

	// ErrNotFound is returned when the no records where matched by the query
	ErrNotFound = errors.New("not found")
)

var getDB = func() (*sql.DB, error) {
	if db == nil {
		var err error
		db, err = sql.Open("mysql", "user:password@/dbname")
		if err != nil {
			panic(err.Error())
		}

		return db, nil
	}
	return db, nil
}

type Person struct {
	ID       int
	FullName string
	Phone    string
	Currency string
	Price    float64
}

func Save(in *Person) (int, error) {
	db, err := getDB()
	if err != nil {
		fmt.Printf("failed to get DB connection: %v \n", err)
		return 0, err
	}

	result, err := db.Exec(sqlInsert, in.FullName, in.Phone, in.Currency, in.Price)
	if err != nil {
		fmt.Printf("failed to save person into DB. err: %v \n", err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Printf("failed to retrieve id of last saved persion. err: %v \n", err)
		return 0, err
	}

	return int(id), nil
}

func LoadAll() ([]*Person, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(sqlLoadAll)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var out []*Person
	for rows.Next() {
		record, err := populatePerson(rows.Scan)
		if err != nil {
			return nil, err
		}

		out = append(out, record)
	}

	if len(out) == 0 {
		return nil, ErrNotFound
	}

	return out, nil
}

func Load(id int) (*Person, error) {
	db, err := getDB()
	if err != nil {
		fmt.Printf("failed to get DB connection. err: %s", err)
		return nil, err
	}

	row := db.QueryRow(sqlLoadByID, id)
	out, err := populatePerson(row.Scan)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("failed to load requested person '%d'. err: %v", id, err)
			return nil, err
		}

		return nil, err
	}

	return out, nil
}

type scanner func(dest ...interface{}) error

func populatePerson(scanner scanner) (*Person, error) {
	out := &Person{}
	err := scanner(&out.ID, &out.FullName, &out.Phone, &out.Currency, &out.Price)
	return out, err
}

func init() {
	_, _ = getDB()
}

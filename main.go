package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type lock struct {
	isLocked    bool
	lockGranted sql.NullString
	lockedBy    sql.NullString
}

const (
	// (LOCKED = b'1') is done to convert tinyint to boolean. See: https://github.com/go-sql-driver/mysql/issues/440
	DB_GET_LOCK_QUERY    = "SELECT (LOCKED = b'1'), LOCKGRANTED, LOCKEDBY from DATABASECHANGELOGLOCK"
	DB_DELETE_LOCK_QUERY = "DELETE FROM DATABASECHANGELOGLOCK"
	DB_ENGINE            = "mysql"
	DB_PORT              = "3306"
	TIME_LAYOUT          = "2006-01-02 15:04:05"
)

var (
	maxLockTime = os.Getenv("MAX_LOCK_TIME")
	DBUser      = os.Getenv("DB_USER")
	DBPass      = os.Getenv("DB_PASS")
	DBURL       = os.Getenv("DB_URL")
	DBName      = os.Getenv("DB_NAME")
)

func main() {
	db := connectToDb()
	defer db.Close()
	result := db.QueryRow(DB_GET_LOCK_QUERY)
	var l = lock{}
	e := result.Scan(&l.isLocked, &l.lockGranted, &l.lockedBy)
	handleError(e)
	if l.isLocked {
		handleDBLock(l, db)
	}
}

func handleDBLock(l lock, db *sql.DB) {
	lockTime, e := time.Parse(TIME_LAYOUT, l.lockGranted.String)
	handleError(e)
	now := time.Now()
	diff := now.Sub(lockTime)
	maxLockTime, e := time.ParseDuration(maxLockTime)
	handleError(e)
	if diff > maxLockTime {
		lockedBy, _ := l.lockedBy.Value()
		log.Printf("Lock from %v older than %v detected, will try to delete...", lockedBy, diff.Minutes())
		res, e := db.Exec(DB_DELETE_LOCK_QUERY)
		handleError(e)
		rows, e := res.RowsAffected()
		handleError(e)
		log.Printf("Deleted locks: %v", rows)
	}
}

func connectToDb() *sql.DB {
	connstr := DBUser + ":" + DBPass + "@tcp(" + DBURL + ":" + DB_PORT + ")/" + DBName
	db, err := sql.Open(DB_ENGINE, connstr)
	handleError(err)
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}

func handleError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

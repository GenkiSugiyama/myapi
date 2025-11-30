package repositories_test

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"testing"
)

var testDB *sql.DB

var (
	dbUser     = "user"
	dbPassword = "pass"
	dbDatabase = "myapi_db"
	dbConn     = fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?parseTime=true", dbUser, dbPassword, dbDatabase)
)

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		os.Exit(1)
	}

	m.Run()

	teardown()
}

func setup() error {
	fmt.Println("connecting to DB for setup...")
	if err := connectDB(); err != nil {
		return err
	}
	fmt.Println("cleaning up DB for setup...")
	if err := cleanupDB(); err != nil {
		fmt.Println("cleanup", err)
		return err
	}
	fmt.Println("setting up test data for setup...")
	if err := setupTestData(); err != nil {
		fmt.Println("setup", err)
		return err
	}
	return nil
}

func connectDB() error {
	var err error
	testDB, err = sql.Open("mysql", dbConn)
	if err != nil {
		return err
	}
	return nil
}

func teardown() {
	fmt.Println("cleaning up DB for teardown...")
	cleanupDB()
	fmt.Println("closing DB connection for teardown...")
	testDB.Close()
}

func setupTestData() error {
	cmd := exec.Command("mysql", "-h", "127.0.0.1", "-u", "user", "myapi_db", "--password=pass")
	file, err := os.Open("./testdata/setupDB.sql")
	if err != nil {
		return err
	}
	defer file.Close()
	cmd.Stdin = file
	return cmd.Run()
}

func cleanupDB() error {
	cmd := exec.Command("mysql", "-h", "127.0.0.1", "-u", "user", "myapi_db", "--password=pass")
	file, err := os.Open("./testdata/cleanupDB.sql")
	if err != nil {
		return err
	}
	defer file.Close()
	cmd.Stdin = file
	return cmd.Run()
}

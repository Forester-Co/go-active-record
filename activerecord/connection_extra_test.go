package activerecord

import (
	"testing"
)

const sqliteDriver = "sqlite3"

func TestConnectAndSetConnection(t *testing.T) {
	db, err := Connect("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	SetConnection(db, "sqlite3")
	if GetConnection() != db {
		t.Error("GetConnection should return the set db")
	}
	if GetDriverName() != sqliteDriver {
		t.Error("GetDriverName should return the set driver name")
	}
}

func TestClose(t *testing.T) {
	SetConnection(nil, "")
	if err := Close(); err != nil {
		t.Error("Close should not fail if db is nil")
	}
	db, _ := Connect("sqlite3", ":memory:")
	SetConnection(db, "sqlite3")
	if err := Close(); err != nil {
		t.Errorf("Close failed: %v", err)
	}
}

func TestBegin_Exec_Query_QueryRow_NoConnection(t *testing.T) {
	SetConnection(nil, "")
	_, err := Begin()
	if err == nil {
		t.Error("Begin should fail if db is nil")
	}
	_, err = Exec("SELECT 1")
	if err == nil {
		t.Error("Exec should fail if db is nil")
	}
	_, err = Query("SELECT 1") //nolint:rowserrcheck // This is expected to fail due to no connection
	if err == nil {
		t.Error("Query should fail if db is nil")
	}
	if rows, ok := err.(interface{ Err() error }); ok {
		if rowsErr := rows.Err(); rowsErr != nil {
			t.Errorf("Rows error: %v", rowsErr)
		}
	}
	row := QueryRow("SELECT 1")
	if row != nil {
		t.Error("QueryRow should return nil if db is nil")
	}
}

func TestExec_Query_QueryRow(t *testing.T) {
	db, _ := Connect("sqlite3", "file:testdb_temp.sqlite?cache=shared&mode=rwc")
	SetConnection(db, "sqlite3")
	if _, err := Exec("DROP TABLE IF EXISTS test"); err != nil {
		t.Fatalf("Failed to drop table: %v", err)
	}
	_, err := Exec("CREATE TABLE test (id INTEGER PRIMARY KEY)")
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
		return
	}
	_, err = Exec("INSERT INTO test (id) VALUES (?)", 1)
	if err != nil {
		t.Errorf("Exec insert failed: %v", err)
	}
	rows, err := Query("SELECT id FROM test")
	if err != nil {
		t.Errorf("Query failed: %v", err)
	}
	defer rows.Close()
	if rowsErr := rows.Err(); rowsErr != nil {
		t.Errorf("Rows error: %v", rowsErr)
	}
	row := QueryRow("SELECT id FROM test WHERE id = ?", 1)
	var id int
	err = row.Scan(&id)
	if err != nil {
		t.Errorf("QueryRow scan failed: %v", err)
	}
	if id != 1 {
		t.Errorf("Expected id 1, got %d", id)
	}
}

func TestConnection(t *testing.T) {
	if GetDriverName() != sqliteDriver {
		t.Skip("This test requires SQLite")
	}
}

func TestSetConnection(t *testing.T) {
	db, _ := Connect("sqlite3", ":memory:")
	SetConnection(db, "sqlite3")
	if GetConnection() != db {
		t.Error("SetConnection should set the connection")
	}
	if GetDriverName() != "sqlite3" {
		t.Error("SetConnection should set the driver name")
	}
}

func TestGetConnection(t *testing.T) {
	db, _ := Connect("sqlite3", ":memory:")
	SetConnection(db, "sqlite3")
	if GetConnection() != db {
		t.Error("GetConnection should return the set connection")
	}
}

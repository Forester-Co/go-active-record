package activerecord

import (
	"testing"
)

type DummyMigration struct {
	upCalled, downCalled *bool
}

func (m *DummyMigration) Version() int64 { return 1 }
func (m *DummyMigration) Up() error      { *m.upCalled = true; return nil }
func (m *DummyMigration) Down() error    { *m.downCalled = true; return nil }

func TestMigrator_Migrate_Status_Rollback(t *testing.T) {
	db, _ := Connect("sqlite3", ":memory:")
	SetConnection(db, "sqlite3")
	migrator := NewMigrator()
	migrator.db = db

	if err := migrator.CreateMigrationsTable(); err != nil {
		t.Fatalf("Failed to create migrations table: %v", err)
	}
	up, down := false, false
	migration := &DummyMigration{&up, &down}
	// Migrate
	err := migrator.Migrate([]MigrationInterface{migration})
	if err != nil {
		t.Errorf("Migrate failed: %v", err)
	}
	if !up {
		t.Error("Up should be called")
	}
	// Status
	err = migrator.Status([]MigrationInterface{migration})
	if err != nil {
		t.Errorf("Status failed: %v", err)
	}
	// Rollback
	err = migrator.Rollback([]MigrationInterface{migration})
	if err != nil {
		t.Errorf("Rollback failed: %v", err)
	}
	if !down {
		t.Error("Down should be called")
	}
}

func TestMigrator_Rollback_NoMigrations(t *testing.T) {
	db, _ := Connect("sqlite3", ":memory:")
	SetConnection(db, "sqlite3")
	migrator := NewMigrator()
	if err := migrator.CreateMigrationsTable(); err != nil {
		t.Fatalf("Failed to create migrations table: %v", err)
	}
	err := migrator.Rollback([]MigrationInterface{})
	if err == nil {
		t.Error("Rollback should fail if no migrations applied")
	}
}

type OtherMigration struct{ DummyMigration }

func (m *OtherMigration) Version() int64 { return 2 }

func TestMigrator_Rollback_NotFound(t *testing.T) {
	db, _ := Connect("sqlite3", ":memory:")
	SetConnection(db, "sqlite3")
	migrator := NewMigrator()
	if err := migrator.CreateMigrationsTable(); err != nil {
		t.Fatalf("Failed to create migrations table: %v", err)
	}
	up, down := false, false
	migration := &DummyMigration{&up, &down}
	if err := migrator.Migrate([]MigrationInterface{migration}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}
	// Try rollback with wrong version
	other := &OtherMigration{DummyMigration{&up, &down}}
	err := migrator.Rollback([]MigrationInterface{other})
	if err == nil {
		t.Error("Rollback should fail if migration not found")
	}
}

func TestCreateTable_DropTable_TableBuilder(t *testing.T) {
	db, _ := Connect("sqlite3", ":memory:")
	SetConnection(db, "sqlite3")
	err := CreateTable("tb", func(tb *TableBuilder) {
		tb.Column("id", "INTEGER", "PRIMARY KEY")
		tb.Column("name", "TEXT")
		tb.Timestamps()
		tb.Index("name")
		// tb.PrimaryKey("id") // Remove for SQLite compatibility
	})
	if err != nil {
		t.Errorf("CreateTable failed: %v", err)
	}
	err = DropTable("tb")
	if err != nil {
		t.Errorf("DropTable failed: %v", err)
	}
}

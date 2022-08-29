package sqlstore

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	ds = Plugin{}
)

func TestCreateDB(t *testing.T) {
	assert.Nil(t, ds.OpenDB("gorm_org.db", "sqlite3"))
}

func TestUnitCreateMember(t *testing.T) {
	ds := Plugin{}
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"sqlite_version()"}).AddRow("3.39.2")
	mock.ExpectQuery("select sqlite_version()").WillReturnRows(rows)
	gdb, err := gorm.Open(sqlite.Dialector{Conn: db}, &gorm.Config{})
	ds.db = gdb
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectExec("CREATE TABLE `members` (`id` integer,`created_at` datetime,`updated_at` datetime,`spiffe_id` text,`description` text,PRIMARY KEY (`id`))").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE UNIQUE INDEX `idx_members_description` ON `members`(`description`)").WillReturnResult(sqlmock.NewResult(1, 1))

	// now we execute our method
	if err = ds.createMemberTableInDB(); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

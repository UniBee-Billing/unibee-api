package pgsql

import (
	"context"
	"database/sql"
	"github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/database/gdb"
)

// PgSqlDriverProxy is a custom database driver, which is used for testing only.
// For simplifying the unit testing case purpose, PgSqlDriverProxy struct inherits the mysql driver
// gdb.Driver and overwrites its functions DoQuery and DoExec.
// So if there's any sql execution, it goes through PgSqlDriverProxy.DoQuery/PgSqlDriverProxy.DoExec firstly
// and then gdb.Driver.DoQuery/gdb.Driver.DoExec.
// You can call it sql "HOOK" or "HiJack" as your will.
type PgSqlDriverProxy struct {
	*pgsql.Driver
}

func init() {
	// It here registers my custom driver in package initialization function "init".
	// You can later use this type in the database configuration.
	if err := gdb.Register("postgres", &PgSqlDriverProxy{}); err != nil {
		panic(err)
	}
}

// New creates and returns a database object for mysql.
// It implements the interface of gdb.Driver for extra database driver installation.
func (d *PgSqlDriverProxy) New(core *gdb.Core, node *gdb.ConfigNode) (gdb.DB, error) {
	return &PgSqlDriverProxy{
		&pgsql.Driver{
			Core: core,
		},
	}, nil
}

// DoCommit commits current sql and arguments to underlying sql driver.
func (d *PgSqlDriverProxy) DoCommit(ctx context.Context, in gdb.DoCommitInput) (out gdb.DoCommitOutput, err error) {
	out, err = d.Driver.DoCommit(ctx, in)
	return
}

func (d *PgSqlDriverProxy) DoInsert(ctx context.Context, link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption) (result sql.Result, err error) {
	if len(list) == 1 {
		// todo mark check id is private key and contain auto increase sequence
		fieldMap := list[0]
		if _, ok := fieldMap["id"]; ok {
			delete(fieldMap, "id")
		}
	}
	return d.Driver.DoInsert(ctx, link, table, list, option)
}

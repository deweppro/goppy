package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/dewep-online/goppy/plugins"
	"github.com/deweppro/go-logger"
	"github.com/deweppro/go-orm"
	"github.com/deweppro/go-orm/schema"
	"github.com/deweppro/go-orm/schema/sqlite"
)

type (
	ConfigSqlite struct {
		Pool []Item `yaml:"sqlite"`
	}
	Item struct {
		Name          string   `yaml:"name"`
		File          string   `yaml:"file"`
		InitMigration []string `yaml:"init_migration"`
	}
)

func (v *ConfigSqlite) Default() {
	if len(v.Pool) == 0 {
		v.Pool = []Item{
			{
				Name: "main",
				File: "./sqlite.db",
				InitMigration: []string{
					"./migration.sql",
				},
			},
		}
	}
}

//List getting all configs
func (v *ConfigSqlite) List() (list []schema.ItemInterface) {
	for _, item := range v.Pool {
		list = append(list, item)
	}
	return
}

//GetName getting config name
func (i Item) GetName() string { return i.Name }

//GetDSN connection params
func (i Item) GetDSN() string { return i.File }

//Setup setting config connections params
func (i Item) Setup(_ schema.SetupInterface) {}

//WithSQLite launch SQLite connection pool
func WithSQLite() plugins.Plugin {
	return plugins.Plugin{
		Config: &ConfigSqlite{},
		Inject: func(conf *ConfigSqlite, log logger.Logger) (*sqliteProvider, SQLite) {
			conn := sqlite.New(conf)
			o := orm.NewDB(conn, orm.Plugins{Logger: log})
			return &sqliteProvider{conn: conn, conf: *conf, log: log}, o
		},
	}
}

type (
	sqliteProvider struct {
		conn schema.Connector
		conf ConfigSqlite
		log  logger.Logger
	}

	//SQLite connection SQLite interface
	SQLite interface {
		Pool(name string) orm.StmtInterface
	}
)

func (v *sqliteProvider) Up() error {
	if err := v.conn.Reconnect(); err != nil {
		return err
	}
	for _, item := range v.conf.Pool {
		p, err := v.conn.Pool(item.Name)
		if err != nil {
			return fmt.Errorf("pool `%s`: %w", item.Name, err)
		}
		if err = p.Ping(); err != nil {
			return fmt.Errorf("pool `%s`: %w", item.Name, err)
		}
		if err = v.migration(p, item.InitMigration); err != nil {
			return fmt.Errorf("pool `%s`: %w", item.Name, err)
		}
		v.log.WithFields(logger.Fields{item.Name: item.File}).Infof("SQLite connect")
	}
	return nil
}

func (v *sqliteProvider) Down() error {
	return v.conn.Close()
}

const sqliteMaster = "select count(*) from `sqlite_master`;"

func (v *sqliteProvider) migration(conn *sql.DB, mig []string) error {
	ctx := context.TODO()
	var count int
	checkDB := func() error {
		row := conn.QueryRowContext(ctx, sqliteMaster)
		if err := row.Scan(&count); err != nil {
			return err
		}
		if err := row.Err(); err != nil {
			return err
		}
		return nil
	}

	if err := checkDB(); err != nil {
		return err
	}

	if count == 0 {
		for _, filename := range mig {
			b, err := os.ReadFile(filename)
			if err != nil {
				return fmt.Errorf("read init migration `%s`: %w", filename, err)
			}
			if _, err = conn.ExecContext(ctx, string(b)); err != nil {
				return fmt.Errorf("exec init migration `%s`: %w", filename, err)
			}
		}
	}

	if err := checkDB(); err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("empty database")
	}

	return nil
}

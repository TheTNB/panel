package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/TheTNB/panel/v2/pkg/io"
	"github.com/TheTNB/panel/v2/pkg/shell"
	"github.com/TheTNB/panel/v2/pkg/systemctl"
	"github.com/TheTNB/panel/v2/pkg/types"
)

type Postgres struct {
	db       *sql.DB
	username string
	password string
	address  string
	port     uint
}

func NewPostgres(username, password, address string, port uint) (*Postgres, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable", address, port, username, password)
	if password == "" {
		dsn = fmt.Sprintf("host=%s port=%d user=%s dbname=postgres sslmode=disable", address, port, username)
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("初始化Postgres连接失败: %w", err)
	}
	if db.Ping() != nil {
		return nil, fmt.Errorf("连接Postgres失败: %w", err)
	}
	return &Postgres{
		db:       db,
		username: username,
		password: password,
		address:  address,
		port:     port,
	}, nil
}

func (m *Postgres) Close() error {
	return m.db.Close()
}

func (m *Postgres) Ping() error {
	return m.db.Ping()
}

func (m *Postgres) Query(query string, args ...any) (*sql.Rows, error) {
	return m.db.Query(query, args...)
}

func (m *Postgres) QueryRow(query string, args ...any) *sql.Row {
	return m.db.QueryRow(query, args...)
}

func (m *Postgres) Exec(query string, args ...any) (sql.Result, error) {
	return m.db.Exec(query, args...)
}

func (m *Postgres) Prepare(query string) (*sql.Stmt, error) {
	return m.db.Prepare(query)
}

func (m *Postgres) DatabaseCreate(name string) error {
	_, err := m.Exec(fmt.Sprintf("CREATE DATABASE %s", name))
	return err
}

func (m *Postgres) DatabaseDrop(name string) error {
	_, err := m.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", name))
	return err
}

func (m *Postgres) UserCreate(user, password string) error {
	_, err := m.Exec(fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s'", user, password))
	if err != nil {
		return err
	}

	return nil
}

func (m *Postgres) UserDrop(user string) error {
	_, err := m.Exec(fmt.Sprintf("DROP USER IF EXISTS %s", user))
	if err != nil {
		return err
	}

	_, _ = shell.Execf(`sed -i '/` + user + `/d' /www/server/postgresql/data/pg_hba.conf`)
	return systemctl.Reload("postgresql")
}

func (m *Postgres) UserPassword(user, password string) error {
	_, err := m.Exec(fmt.Sprintf("ALTER USER %s WITH PASSWORD '%s'", user, password))
	return err
}

func (m *Postgres) PrivilegesGrant(user, database string) error {
	if _, err := m.Exec(fmt.Sprintf("ALTER DATABASE %s OWNER TO %s", database, user)); err != nil {
		return err
	}
	if _, err := m.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s", database, user)); err != nil {
		return err
	}

	return nil
}

func (m *Postgres) PrivilegesRevoke(user, database string) error {
	_, err := m.Exec(fmt.Sprintf("REVOKE ALL PRIVILEGES ON DATABASE %s FROM %s", database, user))
	return err
}

func (m *Postgres) HostAdd(database, user, host string) error {
	config := fmt.Sprintf("host    %s    %s    %s    scram-sha-256", database, user, host)
	if err := io.WriteAppend("/www/server/postgresql/data/pg_hba.conf", config); err != nil {
		return err
	}

	return systemctl.Reload("postgresql")
}

func (m *Postgres) HostRemove(database, user, host string) error {
	regex := fmt.Sprintf(`host\s+%s\s+%s\s+%s`, database, user, host)
	if _, err := shell.Execf(`sed -i '/` + regex + `/d' /www/server/postgresql/data/pg_hba.conf`); err != nil {
		return err
	}

	return systemctl.Reload("postgresql")
}

func (m *Postgres) Users() ([]types.PostgresUser, error) {
	query := `
        SELECT rolname,
               rolsuper,
               rolcreaterole,
               rolcreatedb,
               rolreplication,
               rolbypassrls
        FROM pg_roles
        WHERE rolcanlogin = true;
    `
	rows, err := m.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []types.PostgresUser
	for rows.Next() {
		var user types.PostgresUser
		var super, canCreateRole, canCreateDb, replication, bypassRls bool
		if err = rows.Scan(&user.Role, &super, &canCreateRole, &canCreateDb, &replication, &bypassRls); err != nil {
			return nil, err
		}

		permissions := map[string]bool{
			"超级用户":   super,
			"创建角色":   canCreateRole,
			"创建数据库":  canCreateDb,
			"可以复制":   replication,
			"绕过行级安全": bypassRls,
		}
		for perm, enabled := range permissions {
			if enabled {
				user.Attributes = append(user.Attributes, perm)
			}
		}

		if len(user.Attributes) == 0 {
			user.Attributes = append(user.Attributes, "无")
		}

		users = append(users, user)
	}

	return users, nil
}

func (m *Postgres) Databases() ([]types.PostgresDatabase, error) {
	query := `
        SELECT d.datname, pg_catalog.pg_get_userbyid(d.datdba), pg_catalog.pg_encoding_to_char(d.encoding)
        FROM pg_catalog.pg_database d
        WHERE datistemplate = false;
    `
	rows, err := m.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var databases []types.PostgresDatabase
	for rows.Next() {
		var db types.PostgresDatabase
		if err := rows.Scan(&db.Name, &db.Owner, &db.Encoding); err != nil {
			return nil, err
		}
		databases = append(databases, db)
	}

	return databases, nil
}

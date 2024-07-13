package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/TheTNB/panel/v2/pkg/types"
)

type MySQL struct {
	db       *sql.DB
	username string
	password string
	address  string
}

func NewMySQL(username, password, address string, typ ...string) (*MySQL, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/", username, password, address)
	if len(typ) > 0 && typ[0] == "unix" {
		dsn = fmt.Sprintf("%s:%s@unix(%s)/", username, password, address)
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("初始化MySQL连接失败: %w", err)
	}
	if db.Ping() != nil {
		return nil, fmt.Errorf("连接MySQL失败: %w", err)
	}
	return &MySQL{
		db:       db,
		username: username,
		password: password,
		address:  address,
	}, nil
}

func (m *MySQL) Close() error {
	return m.db.Close()
}

func (m *MySQL) Ping() error {
	return m.db.Ping()
}

func (m *MySQL) Query(query string, args ...any) (*sql.Rows, error) {
	return m.db.Query(query, args...)
}

func (m *MySQL) QueryRow(query string, args ...any) *sql.Row {
	return m.db.QueryRow(query, args...)
}

func (m *MySQL) Exec(query string, args ...any) (sql.Result, error) {
	return m.db.Exec(query, args...)
}

func (m *MySQL) Prepare(query string) (*sql.Stmt, error) {
	return m.db.Prepare(query)
}

func (m *MySQL) DatabaseCreate(name string) error {
	_, err := m.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", name))
	m.flushPrivileges()
	return err
}

func (m *MySQL) DatabaseDrop(name string) error {
	_, err := m.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", name))
	m.flushPrivileges()
	return err
}

func (m *MySQL) UserCreate(user, password string) error {
	_, err := m.Exec(fmt.Sprintf("CREATE USER IF NOT EXISTS '%s'@'localhost' IDENTIFIED BY '%s'", user, password))
	m.flushPrivileges()
	return err
}

func (m *MySQL) UserDrop(user string) error {
	_, err := m.Exec(fmt.Sprintf("DROP USER IF EXISTS '%s'", user))
	m.flushPrivileges()
	return err
}

func (m *MySQL) UserPassword(user, password string) error {
	_, err := m.Exec(fmt.Sprintf("ALTER USER '%s'@'localhost' IDENTIFIED BY '%s'", user, password))
	m.flushPrivileges()
	return err
}

func (m *MySQL) PrivilegesGrant(user, database string) error {
	_, err := m.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON %s.* TO '%s'@'localhost'", database, user))
	m.flushPrivileges()
	return err
}

func (m *MySQL) PrivilegesRevoke(user, database string) error {
	_, err := m.Exec(fmt.Sprintf("REVOKE ALL PRIVILEGES ON %s.* FROM '%s'@'localhost'", database, user))
	m.flushPrivileges()
	return err
}

func (m *MySQL) Users() ([]types.MySQLUser, error) {
	rows, err := m.Query("SELECT user, host FROM mysql.user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []types.MySQLUser
	for rows.Next() {
		var user, host string
		if err := rows.Scan(&user, &host); err != nil {
			continue
		}
		grants, err := m.userGrants(user, host)
		if err != nil {
			continue
		}

		users = append(users, types.MySQLUser{
			User:   user,
			Host:   host,
			Grants: grants,
		})
	}

	return users, nil
}

func (m *MySQL) Databases() ([]types.MySQLDatabase, error) {
	rows, err := m.Query("SHOW DATABASES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var databases []types.MySQLDatabase
	for rows.Next() {
		var database string
		if err := rows.Scan(&database); err != nil {
			continue
		}
		databases = append(databases, types.MySQLDatabase{
			Name: database,
		})
	}
	return databases, nil
}

func (m *MySQL) userGrants(user, host string) ([]string, error) {
	rows, err := m.Query(fmt.Sprintf("SHOW GRANTS FOR '%s'@'%s'", user, host))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var grants []string
	for rows.Next() {
		var grant string
		if err := rows.Scan(&grant); err != nil {
			continue
		}
		grants = append(grants, grant)
	}
	return grants, nil
}

func (m *MySQL) flushPrivileges() {
	_, _ = m.Exec("FLUSH PRIVILEGES")
}

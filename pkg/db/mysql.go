package db

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"github.com/TheTNB/panel/pkg/types"
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
		return nil, fmt.Errorf("init mysql connection failed: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("connect to mysql failed: %w", err)
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

func (m *MySQL) DatabaseExists(name string) (bool, error) {
	rows, err := m.Query("SHOW DATABASES")
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		var database string
		if err := rows.Scan(&database); err != nil {
			continue
		}
		if database == name {
			return true, nil
		}
	}
	return false, nil
}

func (m *MySQL) DatabaseSize(name string) (int64, error) {
	var size int64
	err := m.QueryRow(fmt.Sprintf("SELECT COALESCE(SUM(data_length) + SUM(index_length), 0) FROM information_schema.tables WHERE table_schema = '%s'", name)).Scan(&size)
	return size, err
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

func (m *MySQL) UserPrivileges(user, host string) (map[string][]string, error) {
	rows, err := m.Query(fmt.Sprintf("SHOW GRANTS FOR '%s'@'%s'", user, host))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	privileges := make(map[string][]string)
	for rows.Next() {
		var grant string
		if err = rows.Scan(&grant); err != nil {
			return nil, err
		}
		if !strings.HasPrefix(grant, "GRANT ") {
			continue
		}

		parts := strings.Split(grant, " ON ")
		if len(parts) < 2 {
			continue
		}

		privList := strings.TrimPrefix(parts[0], "GRANT ")
		privs := strings.Split(privList, ", ")

		dbPart := strings.Split(parts[1], " TO")[0]
		// *.* 表示全局权限
		if dbPart == "*.*" {
			dbPart = "*"
		}

		dbPart = strings.Trim(dbPart, "`")
		privileges[dbPart] = append(privileges[dbPart], privs...)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return privileges, nil
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
	query := `
        SELECT 
            SCHEMA_NAME,
            DEFAULT_CHARACTER_SET_NAME,
            DEFAULT_COLLATION_NAME
        FROM INFORMATION_SCHEMA.SCHEMATA
        WHERE SCHEMA_NAME NOT IN ('information_schema', 'performance_schema', 'mysql', 'sys')
    `

	rows, err := m.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var databases []types.MySQLDatabase
	for rows.Next() {
		var db types.MySQLDatabase
		if err = rows.Scan(&db.Name, &db.CharSet, &db.Collation); err != nil {
			return nil, err
		}
		databases = append(databases, db)
	}

	if err = rows.Err(); err != nil {
		return nil, err
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

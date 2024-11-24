package db

import (
	"database/sql"
	"fmt"
	"slices"
	"strings"

	_ "github.com/lib/pq"

	"github.com/TheTNB/panel/pkg/systemctl"
	"github.com/TheTNB/panel/pkg/types"
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
		if username == "" {
			username = "postgres"
		}
		dsn = fmt.Sprintf("host=%s port=%d user=%s dbname=postgres sslmode=disable", address, port, username)
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("init postgres connection failed: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("connect to postgres failed: %w", err)
	}
	return &Postgres{
		db:       db,
		username: username,
		password: password,
		address:  address,
		port:     port,
	}, nil
}

func (r *Postgres) Close() error {
	return r.db.Close()
}

func (r *Postgres) Ping() error {
	return r.db.Ping()
}

func (r *Postgres) Query(query string, args ...any) (*sql.Rows, error) {
	return r.db.Query(query, args...)
}

func (r *Postgres) QueryRow(query string, args ...any) *sql.Row {
	return r.db.QueryRow(query, args...)
}

func (r *Postgres) Exec(query string, args ...any) (sql.Result, error) {
	return r.db.Exec(query, args...)
}

func (r *Postgres) Prepare(query string) (*sql.Stmt, error) {
	return r.db.Prepare(query)
}

func (r *Postgres) DatabaseCreate(name string) error {
	_, err := r.Exec(fmt.Sprintf("CREATE DATABASE %s", name))
	return err
}

func (r *Postgres) DatabaseDrop(name string) error {
	_, err := r.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", name))
	return err
}

func (r *Postgres) DatabaseExist(name string) (bool, error) {
	var count int
	if err := r.QueryRow("SELECT COUNT(*) FROM pg_database WHERE datname = $1", name).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Postgres) DatabaseSize(name string) (int64, error) {
	query := fmt.Sprintf("SELECT pg_database_size('%s')", name)
	var size int64
	if err := r.QueryRow(query).Scan(&size); err != nil {
		return 0, err
	}
	return size, nil
}

func (r *Postgres) UserCreate(user, password string) error {
	_, err := r.Exec(fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s'", user, password))
	if err != nil {
		return err
	}

	return nil
}

func (r *Postgres) UserDrop(user string) error {
	_, err := r.Exec(fmt.Sprintf("DROP USER IF EXISTS %s", user))
	if err != nil {
		return err
	}

	return systemctl.Reload("postgresql")
}

func (r *Postgres) UserPassword(user, password string) error {
	_, err := r.Exec(fmt.Sprintf("ALTER USER %s WITH PASSWORD '%s'", user, password))
	return err
}

func (r *Postgres) UserPrivileges(user string) (map[string][]string, error) {
	query := `
		SELECT 
			table_catalog as database_name,
			string_agg(DISTINCT privilege_type, ',') as privileges
		FROM information_schema.role_database_privileges 
		WHERE grantee = $1
		GROUP BY table_catalog`

	rows, err := r.Query(query, user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	privileges := make(map[string][]string)

	for rows.Next() {
		var db, privilegeStr string
		if err = rows.Scan(&db, &privilegeStr); err != nil {
			return nil, err
		}

		privileges[db] = strings.Split(privilegeStr, ",")
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return privileges, nil
}

func (r *Postgres) PrivilegesGrant(user, database string) error {
	if _, err := r.Exec(fmt.Sprintf("ALTER DATABASE %s OWNER TO %s", database, user)); err != nil {
		return err
	}
	if _, err := r.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s", database, user)); err != nil {
		return err
	}

	return nil
}

func (r *Postgres) PrivilegesRevoke(user, database string) error {
	_, err := r.Exec(fmt.Sprintf("REVOKE ALL PRIVILEGES ON DATABASE %s FROM %s", database, user))
	return err
}

func (r *Postgres) Users() ([]types.PostgresUser, error) {
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
	rows, err := r.Query(query)
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

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Postgres) Databases() ([]types.PostgresDatabase, error) {
	query := `
        SELECT d.datname, pg_catalog.pg_get_userbyid(d.datdba), pg_catalog.pg_encoding_to_char(d.encoding)
        FROM pg_catalog.pg_database d
        WHERE datistemplate = false;
    `
	rows, err := r.Query(query)
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
		if slices.Contains([]string{"template0", "template1", "postgres"}, db.Name) {
			continue
		}
		databases = append(databases, db)
	}

	return databases, nil
}

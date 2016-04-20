package postgresql

import (
	"database/sql"
	"errors"
	"time"

	"gopkg.in/gorp.v1"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tecsisa/authorizr/api"
)

type PostgresRepo struct {
	Dbmap *gorp.DbMap
}

func InitDb(datasourcename string) (*gorp.DbMap, error) {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	db, err := sql.Open("sqlite3", datasourcename)
	if err != nil {
		return nil, err
	}

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	// add a table, setting the table name to 'users' and
	// specifying that the Id property is an auto incrementing PK
	dbmap.AddTableWithName(User{}, "users").SetKeys(true, "Id")

	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		return nil, err
	}

	return dbmap, nil
}

// User database
type User struct {
	Id         int64  `db:"id"`
	Name       string `db:"name"`
	Path       string `db:"path"`
	CreateDate int64  `db:"createdate"`
	Urn        string `db:"urn"`
	Org        string `db:"org"`
}

func (u PostgresRepo) GetUserByID(id uint64) (*api.User, error) {
	obj, err := u.Dbmap.Get(User{}, id)
	if obj != nil {
		user := obj.(*User)
		return &api.User{Id: uint64(user.Id),
			Name: user.Name,
			Path: user.Path,
			Date: time.Unix(0, user.CreateDate),
			Urn:  user.Urn,
			Org:  user.Org,
		}, nil
	}
	return nil, err
}

func (u PostgresRepo) AddUser(user api.User) error {
	userDB := &User{
		Id:         int64(user.Id),
		Name:       user.Name,
		Path:       user.Path,
		CreateDate: time.Now().UTC().UnixNano(),
		Urn:        user.Urn,
		Org:        user.Org,
	}

	return u.Dbmap.Insert(userDB)
}

func (u PostgresRepo) GetUsersByPath(org string, path string) ([]api.User, error) {
	var users []User
	query := "select * from users where org like :org and path like :path"
	_, err := u.Dbmap.Select(&users, query, map[string]interface{}{
		"org":  org,
		"path": path,
	})
	if err != nil {
		return nil, err
	}

	apiusers := make([]api.User, len(users), cap(users))
	for i, u := range users {
		apiusers[i] = api.User{
			Id:   uint64(u.Id),
			Name: u.Name,
			Path: u.Path,
			Date: time.Unix(0, u.CreateDate).UTC(),
			Urn:  u.Urn,
			Org:  u.Org,
		}
	}

	return apiusers, nil
}

func (u PostgresRepo) GetGroupsByUserID(id uint64) ([]api.Group, error) {
	return nil, nil
}

func (u PostgresRepo) RemoveUser(id uint64) error {
	obj, err := u.Dbmap.Get(User{}, id)
	if obj != nil {
		user := obj.(*User)
		_, err := u.Dbmap.Delete(user)
		return err
	} else {
		return errors.New("User not found")
	}
	return err
}

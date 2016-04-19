package postgresql

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tecsisa/authorizr/api"
	"gopkg.in/gorp.v1"
	"log"
	"time"
)

type PostgresRepo struct {
	// TODO: incluir aqui todo lo necesario para conectar a la BD
	Dbmap *gorp.DbMap
}

func InitDb(datasourcename string) *gorp.DbMap {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	db, err := sql.Open("sqlite3", datasourcename)
	checkErr(err, "sql.Open failed")

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	// add a table, setting the table name to 'users' and
	// specifying that the Id property is an auto incrementing PK
	dbmap.AddTableWithName(User{}, "users").SetKeys(true, "Id")

	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
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

func (u PostgresRepo) GetUserByID(id uint64) (api.User, error) {
	obj, err := u.Dbmap.Get(User{}, id)
	if obj != nil {
		user := obj.(*User)
		return api.User{Id: uint64(user.Id),
			Name: user.Name,
			Path: user.Path,
			Date: time.Unix(0, user.CreateDate),
			Urn:  user.Urn,
			Org:  user.Org,
		}, nil
	}
	return api.User{}, err
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
	_, err := u.Dbmap.Select(&users, "select * from users where org like '?' and urn like '%?%'", org, path)
	checkErr(err, "Select users by org and urn failed")

	apiusers := make([]api.User, len(users), cap(users))
	for i, u := range users {
		apiusers[i] = api.User{
			Id:   uint64(u.Id),
			Name: u.Name,
			Path: u.Path,
			Date: time.Unix(u.CreateDate, 0).UTC(),
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
	return nil
}

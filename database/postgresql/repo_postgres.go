package postgresql

import (
	"github.com/tecsisa/authorizr/api"
	"time"
)

type PostgresRepo struct {
	// TODO: incluir aqui todo lo necesario para conectar a la BD
}

func (u PostgresRepo) GetUserByID(id uint64) (api.User, error) {
	return api.User{Id: id,
		Name: "user1",
		Path: "path",
		Date: time.Now(),
		Urn:  "urn:iws:iam:tecsisa:user/path/1",
		Org:  "tecsisa",
	}, nil
}

func (u PostgresRepo) AddUser(user api.User) error {
	return nil
}

func (u PostgresRepo) GetUsersByPath(path string) ([]api.User, error) {
	users := []api.User{
		api.User{
			Id:   1,
			Name: "user1",
			Path: path,
			Date: time.Now(),
			Urn:  "urn:iws:iam:tecsisa:user" + path + "/1",
			Org:  "tecsisa",
		},
		api.User{
			Id:   2,
			Name: "user2",
			Path: path,
			Date: time.Now(),
			Urn:  "urn:iws:iam:tecsisa:user" + path + "/2",
			Org:  "tecsisa",
		},
	}
	return users, nil
}

func (u PostgresRepo) GetGroupsByUserID(id uint64) ([]api.Group, error) {
	return nil, nil
}

func (u PostgresRepo) RemoveUser(id uint64) error {
	return nil
}

package api

import "time"

type UsersAPI struct {

}

func (u *UsersAPI) GetListUsers(path string) ([]User)  {
	users := []User{
		User{
			Id: 1,
			Name: "user1",
			Path: path,
			Date: time.Now(),
			Urn: "urn:iws:iam:tecsisa:user"+path+"/user1",
			Org: "tecsisa",
		},
		User{
			Id: 2,
			Name: "user2",
			Path: path,
			Date: time.Now(),
			Urn: "urn:iws:iam:tecsisa:user"+path+"/user2",
			Org: "tecsisa",
		},
	}
	return users
}

type User struct {
	Id 	uint64 		`json:"ID, omitempty"`
	Name 	string 		`json:"Name, omitempty"`
	Path 	string 		`json:"Path, omitempty"`
	Date 	time.Time 	`json:"Date, omitempty"`
	Urn 	string 		`json:"Urn, omitempty"`
	Org 	string 		`json:"Org, omitempty"`
}
package api

import "time"

type GroupAPI struct {

}

func (u *GroupAPI) GetGroups(path string) (string)  {
	return path
}

type Group struct {
	Id 	uint64 		`json:"ID, omitempty"`
	Name 	string 		`json:"Name, omitempty"`
	Path 	string 		`json:"Path, omitempty"`
	Date 	time.Time 	`json:"Date, omitempty"`
	Urn 	string 		`json:"Urn, omitempty"`
	Org 	string 		`json:"Org, omitempty"`
}


package api

import "time"

type PolicyAPI struct {

}

func (u *PolicyAPI) GetPolicies(path string) (string)  {
	return path
}

type Policy struct {
	Id 		uint64 		`json:"ID, omitempty"`
	Name 		string 		`json:"Name, omitempty"`
	Path 		string 		`json:"Path, omitempty"`
	Date 		time.Time 	`json:"Date, omitempty"`
	Urn 		string 		`json:"Urn, omitempty"`
	Statements 	*[]Statement 	`json:"Statments, omitempty"`
}

type Statement struct {
	Effect 		string 		`json:"Effect, omitempty"`
	Action 		[]string 	`json:"Action, omitempty"`
	Resources 	[]string 	`json:"Resources, omitempty"`
}


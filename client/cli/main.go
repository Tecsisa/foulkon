package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/Tecsisa/foulkon/client"
)

const (
	DEFAULT_ADDRESS = "http://127.0.0.1:8080"
)

func main() {

	help := `Foulkon CLI usage: foulkon [-address=http://1.2.3.4:8080] <command> <action> [<args>]

Available commands:
	user
	group
	policy
	authorize

To get more help, please execute this cli with a <command>

`

	userHelp := `user actions:
	get -id=xxx                   retrieve user xxx
	get-all                       retrieve users
	groups -id=xxx                retrieve user's groups
	create -id=xxx -path=/path    create user 'xxx' with path '/path'
	update -id=xxx -path=/new     update user 'xxx' with path '/new'
	delete -id=xxx                create user 'xxx' with path '/path'

`

	var httpClient http.Client
	var address string
	var command client.Command

	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Printf("%s", help)
		os.Exit(1)
	}

	// get foulkon address
	flag.StringVar(&address, "address", DEFAULT_ADDRESS, "Foulkon Worker address")
	flag.Parse()

	// remove address flag
	r, _ := regexp.Compile(`address`)
	for i, arg := range args {
		if arg == "-address" || arg == "--address" {
			args = append(args[:i], args[i+2:]...)
			break
		}
		if r.MatchString(arg) {
			args = append(args[:i], args[i+1:]...)
			break
		}
	}

	//meta := client.Meta{
	//	address:    address,
	//	httpClient: httpClient,
	//}

	meta := client.Meta{
		Address:    address,
		HttpClient: httpClient,
	}

	// force -h flag
	if len(args) < 2 {
		args = append(args, "-h")
	}

	// parse command/action
	switch args[0] {
	case "user":
		switch args[1] {
		case "get":
			command = &client.GetUserCommand{meta}
		case "get-all":
			command = &client.GetAllUsersCommand{meta}
		case "groups":
			command = &client.GetUserGroupsCommand{meta}
		case "create":
			command = &client.CreateUserCommand{meta}
		case "delete":
			command = &client.DeleteUserCommand{meta}
		case "-h":
			fmt.Printf(userHelp)
			os.Exit(1)
		default:
			fmt.Printf(userHelp)
			os.Exit(1)
		}
	default:
		fmt.Printf("%s", help)
		os.Exit(1)
	}
	msg, err := command.Run(args[2:])
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	println(msg)
	os.Exit(0)
}

package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"strings"

	"github.com/Tecsisa/foulkon/client/api"
)

const (
	DEFAULT_ADDRESS = "http://127.0.0.1:8080"
	FLAG_EXTERNALID = "externalId"
	FLAG_OFFSET     = "offset"
	FLAG_LIMIT      = "limit"
	FLAG_ORDERBY    = "orderBy"
	FLAG_PATHPREFIX = "pathPrefix"
	FLAG_PATH       = "path"
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

	var address string
	var clientApi api.ClientAPI

	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Printf("%s", help)
		os.Exit(1)
	}

	availableFlags := map[string]string{
		FLAG_OFFSET:     "The offset of the items returned",
		FLAG_LIMIT:      "The maximum number of items in the response",
		FLAG_ORDERBY:    "Order data by field",
		FLAG_PATHPREFIX: "Search starts from this path",
		FLAG_EXTERNALID: "User's external identifier",
		FLAG_PATH:       "User location",
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

	clientApi.Address = address

	// force -h flag
	if len(args) < 2 {
		args = append(args, "-h")
	}

	var msg string
	var err error
	// parse command/action
	switch args[0] {
	case "user":
		switch args[1] {
		case "get":
			params := parseFlags(availableFlags, []string{FLAG_EXTERNALID}, args)
			msg, err = clientApi.GetUser(params[FLAG_EXTERNALID])
		case "get-all":
			params := parseFlags(availableFlags, []string{FLAG_PATHPREFIX, FLAG_OFFSET, FLAG_LIMIT, FLAG_ORDERBY}, args)
			msg, err = clientApi.GetAllUsers(params[FLAG_PATHPREFIX], params[FLAG_OFFSET], params[FLAG_LIMIT], params[FLAG_ORDERBY])
		case "groups":
			params := parseFlags(availableFlags, []string{FLAG_EXTERNALID, FLAG_OFFSET, FLAG_LIMIT, FLAG_ORDERBY}, args)
			msg, err = clientApi.GetAllUsers(params[FLAG_EXTERNALID], params[FLAG_OFFSET], params[FLAG_LIMIT], params[FLAG_ORDERBY])
		case "create":
			params := parseFlags(availableFlags, []string{FLAG_EXTERNALID, FLAG_PATH}, args)
			msg, err = clientApi.CreateUser(params[FLAG_EXTERNALID], params[FLAG_PATH])
		case "delete":
			params := parseFlags(availableFlags, []string{FLAG_EXTERNALID}, args)
			msg, err = clientApi.DeleteUser(params[FLAG_EXTERNALID])
		case "update":
			params := parseFlags(availableFlags, []string{FLAG_EXTERNALID, FLAG_PATH}, args)
			msg, err = clientApi.UpdateUser(params[FLAG_EXTERNALID], params[FLAG_PATH])
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
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	println(msg)
	os.Exit(0)
}

// Helper func for updating request params
func parseFlags(availableFlags map[string]string, validFlags []string, cliArgs []string) map[string]string {
	params := make(map[string]string)

	flagSet := flag.NewFlagSet(cliArgs[0]+" "+cliArgs[1], flag.ExitOnError)

	for _, val := range validFlags {
		flagSet.String(val, "", availableFlags[val])
	}

	if err := flagSet.Parse(cliArgs[2:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	for _, v := range validFlags {
		if val := flagSet.Lookup(v); val != nil {
			params[strings.ToUpper(v)] = val.Value.String()
		}
	}

	return params
}

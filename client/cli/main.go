package main

import (
	"flag"
	"fmt"
	"os"

	"strings"

	"github.com/Tecsisa/foulkon/client/api"
)

const (
	DEFAULT_ADDRESS     = "http://127.0.0.1:8000"
	FLAG_EXTERNALID     = "externalId"
	FLAG_ORGANIZATIONID = "orgId"
	FLAG_POLICYNAME     = "policyName"
	FLAG_OFFSET         = "offset"
	FLAG_LIMIT          = "limit"
	FLAG_ORDERBY        = "orderBy"
	FLAG_PATHPREFIX     = "pathPrefix"
	FLAG_PATH           = "path"
	FLAG_EFFECT         = "effect"
	FLAG_ACTIONS        = "actions"
	FLAG_RESOURCES      = "resources"
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
	userHelp := `User actions:
	get -id=xxx                   retrieve user xxx
	get-all                       retrieve users
	groups -id=xxx                retrieve user's groups
	create -id=xxx -path=/path/   create user 'xxx' with path '/path/'
	update -id=xxx -path=/new/    update user 'xxx' with path '/new/'
	delete -id=xxx                delete user 'xxx'
`
	policyHelp := `Policy actions:
	get -policyName=xxx                   	retrieve policy xxx
	get-all                       		retrieve all policies
	groups-policy -id=xxx                	retrieve all group with policy 'xxx' attached to
	policies-organization -orgId=yyy	retrieve all policies that belong to oranization 'yyy'
	create -id=xxx -path=/path/   		create policy 'xxx' with path '/path/'
	update -policyName=xxx -orgId=yyy    	update policy 'xxx' that belong to oranization 'yyy'
	delete -policyName=xxx -orgId=yyy     	delete policy 'xxx' that belong to oranization 'yyy'
`

	var clientApi api.ClientAPI
	availableFlags := map[string]string{
		FLAG_OFFSET:         "The offset of the items returned",
		FLAG_EXTERNALID:     "User's external identifier",
		FLAG_ORGANIZATIONID: "Policy organization",
		FLAG_POLICYNAME:     "Policy name",
		FLAG_LIMIT:          "The maximum number of items in the response",
		FLAG_ORDERBY:        "Order data by field",
		FLAG_PATHPREFIX:     "Search starts from this path",
		FLAG_PATH:           "--- location",
		FLAG_EFFECT:         "flag effect",
		FLAG_ACTIONS:        "flag actions",
		FLAG_RESOURCES:      "flag resources",
	}

	//remove program path name from args
	args := os.Args[1:]

	// get foulkon address
	flag.StringVar(&clientApi.Address, "address", DEFAULT_ADDRESS, "Foulkon Worker address")
	flag.Parse()

	// remove address flag
	for i, arg := range args {
		if arg == "-address" || arg == "--address" {
			args = append(args[:i], args[i+2:]...)
			break
		}
		if strings.Contains(arg, "address=") {
			args = append(args[:i], args[i+1:]...)
			break
		}
	}

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
			fallthrough
		default:
			msg = userHelp
		}
	case "policy":
		switch args[1] {
		case "get":
			params := parseFlags(availableFlags, []string{FLAG_ORGANIZATIONID, FLAG_POLICYNAME}, args)
			msg, err = clientApi.GetPolicy(params[FLAG_ORGANIZATIONID], params[FLAG_POLICYNAME])
		case "get-all":
			params := parseFlags(availableFlags, []string{FLAG_PATHPREFIX, FLAG_OFFSET, FLAG_LIMIT, FLAG_ORDERBY}, args)
			msg, err = clientApi.GetAllPolicy(params[FLAG_PATHPREFIX], params[FLAG_OFFSET], params[FLAG_LIMIT], params[FLAG_ORDERBY])
		case "create":
			params := parseFlags(availableFlags, []string{FLAG_ORGANIZATIONID, FLAG_POLICYNAME, FLAG_PATH, FLAG_EFFECT, FLAG_ACTIONS, FLAG_RESOURCES}, args)
			msg, err = clientApi.CreatePolicy(params[FLAG_ORGANIZATIONID], params[FLAG_POLICYNAME], params[FLAG_PATH], params[FLAG_EFFECT], params[FLAG_ACTIONS], params[FLAG_RESOURCES])
		case "update":
			//params := parseFlags(availableFlags, []string{FLAG_ORGANIZATIONID, FLAG_POLICYNAME}, args)
			//msg, err = clientApi.UpdatePolicy(params[FLAG_ORGANIZATIONID], params[FLAG_POLICYNAME])
		case "delete":
			params := parseFlags(availableFlags, []string{FLAG_ORGANIZATIONID, FLAG_POLICYNAME}, args)
			msg, err = clientApi.DeletePolicy(params[FLAG_ORGANIZATIONID], params[FLAG_POLICYNAME])
		case "groups-policy":
			params := parseFlags(availableFlags, []string{FLAG_ORGANIZATIONID, FLAG_POLICYNAME, FLAG_OFFSET, FLAG_LIMIT, FLAG_ORDERBY}, args)
			msg, err = clientApi.GetGroupsPolicy(params[FLAG_ORGANIZATIONID], params[FLAG_POLICYNAME], params[FLAG_OFFSET], params[FLAG_LIMIT], params[FLAG_ORDERBY])
		case "policies-organization":
			params := parseFlags(availableFlags, []string{FLAG_ORGANIZATIONID, FLAG_PATHPREFIX, FLAG_OFFSET, FLAG_LIMIT, FLAG_ORDERBY}, args)
			msg, err = clientApi.GetPoliciesOrganization(params[FLAG_ORGANIZATIONID], params[FLAG_PATHPREFIX], params[FLAG_OFFSET], params[FLAG_LIMIT], params[FLAG_ORDERBY])
		case "-h":
			fallthrough
		default:
			msg = policyHelp
		}
	default:
		msg = help
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
			params[v] = val.Value.String()
		}
	}

	return params
}

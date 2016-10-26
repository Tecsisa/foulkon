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
	FLAG_EXTERNALID     = "extId"
	FLAG_ORGANIZATIONID = "orgId"
	FLAG_GROUPNAME      = "groupName"
	FLAG_NEWGROUPNAME   = "newGroupName"
	FLAG_POLICYNAME     = "policyName"
	FLAG_USERNAME       = "userName"
	FLAG_OFFSET         = "offset"
	FLAG_LIMIT          = "limit"
	FLAG_ORDERBY        = "orderBy"
	FLAG_PATHPREFIX     = "pathPrefix"
	FLAG_NEWPATH        = "newPath"
	FLAG_PATH           = "path"
	FLAG_STATEMENT      = "statement"
)

type Cli struct {
	UserApi   api.UserAPI
	GroupApi  api.GroupAPI
	PolicyApi api.PolicyAPI
}

// Helper func for updating request params
func parseFlags(availableFlags map[string]string, validFlags, cliArgs []string, requireFlags int) map[string]string {
	params := make(map[string]string)

	flagSet := flag.NewFlagSet(cliArgs[0]+" "+cliArgs[1], flag.ExitOnError)

	for _, val := range validFlags {
		flagSet.String(val, "", availableFlags[val])
	}

	if err := flagSet.Parse(cliArgs[2:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	for i, v := range validFlags {
		val := flagSet.Lookup(v).Value.String()
		if i < requireFlags && val == "" {
			return nil
		}
		params[v] = val
	}

	return params
}

func main() {

	availableFlags := map[string]string{
		FLAG_OFFSET:         "Offset of returned items",
		FLAG_EXTERNALID:     "User's external identifier",
		FLAG_ORGANIZATIONID: "Policy's organization",
		FLAG_POLICYNAME:     "Policy's name",
		FLAG_LIMIT:          "Maximum number of items in response",
		FLAG_ORDERBY:        "Sort the result by specified column",
		FLAG_PATHPREFIX:     "Search starts from this path",
		FLAG_PATH:           "--- location",
		FLAG_STATEMENT:      "policy's statement",
		FLAG_GROUPNAME:      "Group's name",
		FLAG_NEWGROUPNAME:   "New group name",
		FLAG_USERNAME:       "User's Name",
		FLAG_NEWPATH:        "New Path",
	}

	help := `Foulkon CLI usage: foulkon [-address=http://1.2.3.4:8080] <command> <action> [<args>]

Available commands:
	user
	group
	policy
	authorize

To get more help, please execute this cli with a <command>`

	userHelp := `User actions:
	get -id=xxx                   		retrieve user xxx
	get-all <optionalParams>      		retrieve users
	groups -id=xxx <optionalParams> 	retrieve user's groups
	create -id=xxx -path=/path/   		create user 'xxx' with path '/path/'
	update -id=xxx -path=/new/    		update user 'xxx' with path '/new/'
	delete -id=xxx                		delete user 'xxx'

	optionalParams:
	-pathPrefix, -offset, -limit, -orderBy			Control de output in list actions
	`

	groupHelp := `Group actions:
	get -groupName=xxx                                                 retrieve group xxx
	get-all <optionalParams>                                           retrieve all groups
	get-org-groups -orgId=xxx <optionalParams>                         retrieve all groups within an organization xxx
	create -groupName=xxx -orgName=xxx -path=/path/                    create group 'xxx' with path '/path/'
	update -orgId=xxx -groupName=xxx -newGroupName=xxx -newPath=xxx    update group 'xxx' that belong to oranization 'yyy'
	delete -policyName=xxx -orgId=yyy                                  delete group 'xxx' that belong to oranization 'yyy'
	get-members -orgId=xxx -groupName=xxx <optionalParams>             retrieve group members
	add-member -orgId=xxx -groupName=xxx -userName=yyy                 add member 'yyy' to group 'xxx'
	remove-member -orgId=xxx -groupName=xxx -userName=yyy              remove member 'yyy' from group 'xxx'
	get-policies -orgId=xxx -groupName=xxx <optionalParams>            retrieve group policies
	attach-policy -orgId=xxx -groupName=xxx -policyName=yyy            attach policy 'yyy' to group 'xxx'
	detach-policy -orgId=xxx -groupName=xxx -policyName=yyy            detach policy 'yyy' to group 'xxx'

	Optional Parameters:
	-pathPrefix, -offset, -limit, -orderBy			Control de output in list actions
	`

	policyHelp := `Policy actions:
	get -orgId=yyy -policyName=xxx   				retrieve policy xxx that belong to organizaation 'yyy'
	get-all <optionalParams>					retrieve all policies
	groups-attached -orgId=yyy -policyName=xxx <optionalParams>    	retrieve all groups with policy 'xxx', that belong to organizaation 'yyy', attached to
	policies-organization -orgId=yyy <optionalParams>		retrieve all policies that belong to oranization 'yyy'
	create -orgId=yyy -policyName=xxx -path=/path/ -statement=zzz  	create policy 'xxx' with path '/path/' that belong to organizaation 'yyy' and with statements 'zzz'(JSON format)
	update -orgId=yyy -policyName=xxx -statement=zzz	 	update policy 'xxx' with path '/path/' that belong to organizaation 'yyy' and new statements 'zzz'(JSON format)
	delete -orgId=yyy -policyName=xxx 				delete policy 'xxx' that belong to oranization 'yyy'

	Optional Parameters:
	-pathPrefix, -offset, -limit, -orderBy			Control de output in list actions
	`

	var cli Cli
	clientApi := &api.ClientAPI{}
	cli.UserApi = clientApi

	cli.PolicyApi = clientApi

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

	// remove help flag
	for i, arg := range args {
		if arg == "help" {
			args[i] = "--help"
			break
		}
	}

	// force -h flag
	if len(args) < 2 {
		args = append(args, "--help")
	}

	//statement_example := `[{"effect":"allow","actions":["iam:getUser","iam:*"],"resources":["urn:everything:*"]}]`
	var msg string
	var err error

	// parse command/action
	switch args[0] {
	case "user":
		switch args[1] {
		case "get":

			if params := parseFlags(availableFlags, []string{FLAG_EXTERNALID}, args, 1); params == nil {
				msg = userHelp
			} else {
				msg, err = cli.UserApi.GetUser(params[FLAG_EXTERNALID])
			}
		case "get-all":
			if params := parseFlags(availableFlags, []string{FLAG_PATHPREFIX, FLAG_OFFSET, FLAG_LIMIT, FLAG_ORDERBY}, args, 0); params == nil {
				msg = userHelp
			} else {
				msg, err = cli.UserApi.GetAllUsers(params[FLAG_PATHPREFIX], params[FLAG_OFFSET], params[FLAG_LIMIT], params[FLAG_ORDERBY])
			}
		case "groups":
			if params := parseFlags(availableFlags, []string{FLAG_EXTERNALID, FLAG_OFFSET, FLAG_LIMIT, FLAG_ORDERBY}, args, 1); params == nil {
				msg = userHelp
			} else {
				msg, err = cli.UserApi.GetAllUsers(params[FLAG_EXTERNALID], params[FLAG_OFFSET], params[FLAG_LIMIT], params[FLAG_ORDERBY])
			}
		case "create":
			if params := parseFlags(availableFlags, []string{FLAG_EXTERNALID, FLAG_PATH}, args, 2); params == nil {
				msg = userHelp
			} else {
				msg, err = cli.UserApi.CreateUser(params[FLAG_EXTERNALID], params[FLAG_PATH])
			}
		case "delete":
			if params := parseFlags(availableFlags, []string{FLAG_EXTERNALID}, args, 1); params == nil {
				msg = userHelp
			} else {
				msg, err = cli.UserApi.DeleteUser(params[FLAG_EXTERNALID])
			}
		case "update":
			if params := parseFlags(availableFlags, []string{FLAG_EXTERNALID, FLAG_PATH}, args, 2); params == nil {
				msg = userHelp
			} else {
				msg, err = cli.UserApi.UpdateUser(params[FLAG_EXTERNALID], params[FLAG_PATH])
			}
		case "--help":
			fallthrough
		default:
			msg = userHelp
		}
	case "policy":
		switch args[1] {
		case "get":
			if params := parseFlags(availableFlags, []string{FLAG_ORGANIZATIONID, FLAG_POLICYNAME}, args, 2); params == nil {
				msg = policyHelp
			} else {
				msg, err = cli.PolicyApi.GetPolicy(params[FLAG_ORGANIZATIONID], params[FLAG_POLICYNAME])
			}
		case "get-all":
			if params := parseFlags(availableFlags, []string{FLAG_PATHPREFIX, FLAG_OFFSET, FLAG_LIMIT, FLAG_ORDERBY}, args, 0); params == nil {
				msg = policyHelp
			} else {
				msg, err = cli.PolicyApi.GetAllPolicies(params[FLAG_PATHPREFIX], params[FLAG_OFFSET], params[FLAG_LIMIT], params[FLAG_ORDERBY])
			}
		case "groups-attached":
			if params := parseFlags(availableFlags, []string{FLAG_ORGANIZATIONID, FLAG_POLICYNAME, FLAG_OFFSET, FLAG_LIMIT, FLAG_ORDERBY}, args, 2); params == nil {
				msg = policyHelp
			} else {
				msg, err = cli.PolicyApi.GetGroupsAttached(params[FLAG_ORGANIZATIONID], params[FLAG_POLICYNAME], params[FLAG_OFFSET], params[FLAG_LIMIT], params[FLAG_ORDERBY])
			}
		case "policies-organization":
			if params := parseFlags(availableFlags, []string{FLAG_ORGANIZATIONID, FLAG_PATHPREFIX, FLAG_OFFSET, FLAG_LIMIT, FLAG_ORDERBY}, args, 1); params == nil {
				msg = policyHelp
			} else {
				msg, err = cli.PolicyApi.GetPoliciesOrganization(params[FLAG_ORGANIZATIONID], params[FLAG_PATHPREFIX], params[FLAG_OFFSET], params[FLAG_LIMIT], params[FLAG_ORDERBY])
			}
		case "create":
			if params := parseFlags(availableFlags, []string{FLAG_ORGANIZATIONID, FLAG_POLICYNAME, FLAG_PATH, FLAG_STATEMENT}, args, 4); params == nil {
				msg = policyHelp
			} else {
				msg, err = cli.PolicyApi.CreatePolicy(params[FLAG_ORGANIZATIONID], params[FLAG_POLICYNAME], params[FLAG_PATH], params[FLAG_STATEMENT])
			}
		case "update":
			if params := parseFlags(availableFlags, []string{FLAG_ORGANIZATIONID, FLAG_POLICYNAME, FLAG_PATH, FLAG_STATEMENT}, args, 4); params == nil {
				msg = policyHelp
			} else {
				msg, err = cli.PolicyApi.UpdatePolicy(params[FLAG_ORGANIZATIONID], params[FLAG_POLICYNAME], params[FLAG_PATH], params[FLAG_STATEMENT])
			}
		case "delete":
			if params := parseFlags(availableFlags, []string{FLAG_ORGANIZATIONID, FLAG_POLICYNAME}, args, 2); params == nil {
				msg = policyHelp
			} else {
				msg, err = cli.PolicyApi.DeletePolicy(params[FLAG_ORGANIZATIONID], params[FLAG_POLICYNAME])
			}
		case "-help":
			fallthrough
		default:
			msg = policyHelp
		}

	case "group":
		switch args[1] {
		case "get":
			if params := parseFlags(availableFlags, []string{FLAG_ORGANIZATIONID, FLAG_GROUPNAME}, args, 2); params == nil {
				msg = groupHelp
			} else {
				msg, err = cli.GroupApi.GetGroup(params[FLAG_ORGANIZATIONID], params[FLAG_GROUPNAME])
			}
		case "get-all":
			if params := parseFlags(availableFlags, []string{FLAG_PATHPREFIX, FLAG_OFFSET, FLAG_LIMIT, FLAG_ORDERBY}, args, 0); params == nil {
				msg = groupHelp
			} else {
				msg, err = cli.GroupApi.GetAllGroups(params[FLAG_PATHPREFIX], params[FLAG_OFFSET], params[FLAG_LIMIT], params[FLAG_ORDERBY])
			}
		case "get-org-groups":
			if params := parseFlags(availableFlags, []string{FLAG_ORGANIZATIONID, FLAG_PATHPREFIX, FLAG_OFFSET, FLAG_LIMIT, FLAG_ORDERBY}, args, 1); params == nil {
				msg = groupHelp
			} else {
				msg, err = cli.GroupApi.GetGroupsByOrg(params[FLAG_ORGANIZATIONID], params[FLAG_PATHPREFIX], params[FLAG_OFFSET], params[FLAG_LIMIT], params[FLAG_ORDERBY])
			}
		case "create":
			if params := parseFlags(availableFlags, []string{FLAG_ORGANIZATIONID, FLAG_GROUPNAME, FLAG_PATH}, args, 3); params == nil {
				msg = groupHelp
			} else {
				msg, err = cli.GroupApi.CreateGroup(params[FLAG_ORGANIZATIONID], params[FLAG_GROUPNAME], params[FLAG_PATH])
			}
		case "update":
			if params := parseFlags(availableFlags, []string{FLAG_ORGANIZATIONID, FLAG_GROUPNAME, FLAG_NEWGROUPNAME, FLAG_NEWPATH}, args, 4); params == nil {
				msg = groupHelp
			} else {
				msg, err = cli.GroupApi.UpdateGroup(params[FLAG_ORGANIZATIONID], params[FLAG_GROUPNAME], params[FLAG_NEWGROUPNAME], params[FLAG_NEWPATH])
			}
		case "delete":
			if params := parseFlags(availableFlags, []string{FLAG_ORGANIZATIONID, FLAG_GROUPNAME}, args, 2); params == nil {
				msg = groupHelp
			} else {
				msg, err = cli.GroupApi.DeleteGroup(params[FLAG_ORGANIZATIONID], params[FLAG_GROUPNAME])
			}
		case "get-members":
			if params := parseFlags(availableFlags, []string{FLAG_ORGANIZATIONID, FLAG_GROUPNAME, FLAG_PATHPREFIX, FLAG_OFFSET, FLAG_LIMIT, FLAG_ORDERBY}, args, 2); params == nil {
				msg = groupHelp
			} else {
				msg, err = cli.GroupApi.GetGroupMembers(params[FLAG_ORGANIZATIONID], params[FLAG_GROUPNAME], params[FLAG_PATHPREFIX], params[FLAG_OFFSET], params[FLAG_LIMIT], params[FLAG_ORDERBY])
			}
		case "add-member":
			if params := parseFlags(availableFlags, []string{FLAG_ORGANIZATIONID, FLAG_GROUPNAME, FLAG_USERNAME}, args, 3); params == nil {
				msg = groupHelp
			} else {
				msg, err = cli.GroupApi.AddMemberToGroup(params[FLAG_ORGANIZATIONID], params[FLAG_GROUPNAME], params[FLAG_USERNAME])
			}
		case "remove-member":
			if params := parseFlags(availableFlags, []string{FLAG_ORGANIZATIONID, FLAG_GROUPNAME, FLAG_USERNAME}, args, 3); params == nil {
				msg = groupHelp
			} else {
				msg, err = cli.GroupApi.RemoveMemberFromGroup(params[FLAG_ORGANIZATIONID], params[FLAG_GROUPNAME], params[FLAG_USERNAME])
			}
		case "get-policies":
			if params := parseFlags(availableFlags, []string{FLAG_ORGANIZATIONID, FLAG_GROUPNAME, FLAG_PATHPREFIX, FLAG_OFFSET, FLAG_LIMIT, FLAG_ORDERBY}, args, 2); params == nil {
				msg = groupHelp
			} else {
				msg, err = cli.GroupApi.GetGroupPolicies(params[FLAG_ORGANIZATIONID], params[FLAG_GROUPNAME], params[FLAG_OFFSET], params[FLAG_LIMIT], params[FLAG_ORDERBY])
			}
		case "attach-policy":
			if params := parseFlags(availableFlags, []string{FLAG_ORGANIZATIONID, FLAG_GROUPNAME, FLAG_POLICYNAME}, args, 3); params == nil {
				msg = groupHelp
			} else {
				msg, err = cli.GroupApi.AttachPolicyToGroup(params[FLAG_ORGANIZATIONID], params[FLAG_GROUPNAME], params[FLAG_POLICYNAME])
			}
		case "detach-policy":
			if params := parseFlags(availableFlags, []string{FLAG_ORGANIZATIONID, FLAG_GROUPNAME, FLAG_POLICYNAME}, args, 3); params == nil {
				msg = groupHelp
			} else {
				msg, err = cli.GroupApi.DetachPolicyFromGroup(params[FLAG_ORGANIZATIONID], params[FLAG_GROUPNAME], params[FLAG_POLICYNAME])
			}
		case "-h":
			fallthrough
		default:
			msg = groupHelp
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

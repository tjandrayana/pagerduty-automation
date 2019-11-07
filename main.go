package main

import (
	"fmt"
	"os"
	"flag"
	"strings"

	pager "github.com/PagerDuty/go-pagerduty"
	"github.com/tjandrayana/pagerduty-automation/pagerduty"
)

const (
	TYPE_USER = "user"
	DEFAULT_TIMEZONE = "Asia/Jakarta"

	PD_TOKEN = ""
)

type Arguments []string

var m pagerduty.Module

func main() {
	var err error
	m = pagerduty.Init(PD_TOKEN)
	
	// Get User Input Parameter
	command, args := getCommand()

	// Do Command
	switch command {
	case "create_user":
		err = createUser(args)
	case "get_users":
		err = getAllUsers(args)
	case "set_rules":
		err = setNotificationRules(args)
	default:
		err = help("")
	}

	if err != nil {
		fmt.Println("Failed to do command.", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func getCommand() (string, []string) {
	args := Arguments(os.Args)
	if len(args) <= 1 {
		return "", []string{}
	}
	return args[1], args[2:]
}

func help(about string) error {
	if about == "create_user" {
		fmt.Println("This is help page.")
		return nil
	}

	fmt.Println("This is help page.")
	return nil
}

func createUser(args []string) error {
	if len(args) < 2 {
		return help("create_user")
	}

	var role, job string
	flag.StringVar(&role, "role", "", "User Role")
	flag.StringVar(&job, "job", "", "User Job")
	flag.Parse()

	name := args[0]
	email := args[1]
	if role == "" {
		role = "limited_user"
	}

	user, err := m.CreateUser(pager.User{
		Type:        TYPE_USER,
		Name:        name,
		Email:       email,
		Timezone:    DEFAULT_TIMEZONE,
		Role:        role,
		JobTitle:    job,
	})
	if err != nil {
		return fmt.Errorf("Failed to create new user because %s", err)
	}

	fmt.Printf("Success create pagerduty user %s with email %s\n", user.ID, email)
	return nil
}

func getAllUsers(args []string) error {
	if len(args) < 1 {
		return help("get_all_user")
	}

	userType := pagerduty.ConstAllUser
	if args[0] == "miss" {
		userType = pagerduty.ConstUserHasNotBeenValidated
	}

	users, err := m.ListUser("", userType)
	if err != nil {
		return fmt.Errorf("Failed to get list users because %s", err)
	}

	for i, u := range users {
		hasPhone := false
		hasSMS := false
		hasPushNotif := false
		for _, w := range u.ContactMethods {
			if w.Type == "phone_contact_method_reference" {
				hasPhone = true
			}
			if w.Type == "sms_contact_method_reference" {
				hasSMS = true
			}
			if w.Type == "push_notification_contact_method_reference" {
				hasPushNotif = true
			}
		}

		missing := []string{}
		if hasPhone == false {
			missing = append(missing, "phone")
		}
		if hasSMS == false {
			missing = append(missing, "sms")
		}
		if hasPushNotif == false {
			missing = append(missing, "push_notif")
		}

		fmt.Printf("%s => %s: %s\n", rightPad(fmt.Sprintf("%d", i+1), 4), rightPad(u.Email, 50), strings.Join(missing, "-"))
	}

	return nil
}

func setNotificationRules(args []string) error {
	if len(args) < 1 {
		return help("set_notification_rules")
	}

	email := args[0]

	users, err := m.ListUser(email, pagerduty.ConstAllUser)
	if err != nil {
		return err
	}

	if len(users) == 0 {
		return fmt.Errorf("User %s is not found", email)
	} else if len(users) > 1 {
		return fmt.Errorf("Email %s found more than 1. Please detail the email", email)
	}

	err = m.SetDefaultNotification(users[0].ID)
	if err != nil {
		return err
	}

	fmt.Printf("Success set default notification rules for user %s\n", email)
	return nil
}

//------------------------------------------------------------[TOOLS]

func rightPad(str string, count int) string {
	if len(str) > count {
		runes := []rune(str)
		return string(runes[:count])
	} else if len(str) < count {
		for i := len(str); i < count; i++ {
			str = str + " "
		}
	}
	return str
}

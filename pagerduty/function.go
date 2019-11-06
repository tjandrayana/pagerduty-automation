package pagerduty

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	pager "github.com/PagerDuty/go-pagerduty"
	"github.com/tjandrayana/pagerduty-automation/httpreq"
)

func (m Module) GetTeamDetailByName(teamName string) []pager.Team {

	var team TeamQueryResult
	address := "https://api.pagerduty.com"
	path := "teams"

	header := make(map[string]string)
	header["Accept"] = "application/vnd.pagerduty+json;version=2"
	header["Authorization"] = fmt.Sprintf("Token token=%s", m.Token)
	params := url.Values{}
	params.Add("query", teamName)

	agent := httpreq.NewHTTPRequest()
	agent.Url = address
	agent.Path = path
	agent.Method = "GET"
	agent.Headers = header
	agent.Param = params

	body, err := agent.DoReq()
	if err != nil {
		log.Println(err)
		agent.Debug()
		if body != nil {
			log.Println(string(*body))
		}
		return team.Teams
	}
	if err := json.Unmarshal(*body, &team); err != nil {
		log.Println(err)
		log.Println(string(*body))
		return team.Teams
	}

	return team.Teams
}

func (m Module) CreateUser(detail pager.User) (*pager.User, error) {
	user, err := m.Client.CreateUser(detail)
	if err != nil {
		log.Println(err)
		return user, err
	}

	fmt.Println("success")
	return user, err
}

func (m Module) CreateUserNotificationRules(id string, detail pager.NotificationRule) (pager.NotificationRule, error) {
	notif := Notification{
		detail,
	}
	address := "https://api.pagerduty.com"
	path := "users/" + id + "/notification_rules"

	header := make(map[string]string)
	header["Accept"] = "application/vnd.pagerduty+json;version=2"
	header["Authorization"] = fmt.Sprintf("Token token=%s", m.Token)

	agent := httpreq.NewHTTPRequest()
	agent.Url = address
	agent.Path = path
	agent.Method = "POST"
	agent.Headers = header
	agent.IsJson = true
	agent.Json = notif

	body, err := agent.DoReq()
	if err != nil {
		log.Println(err)
		agent.Debug()
		if body != nil {
			log.Println(string(*body))
		}
		return notif.NotificationRule, err
	}
	if err := json.Unmarshal(*body, &notif); err != nil {
		log.Println(err)
		log.Println(string(*body))
		return notif.NotificationRule, err
	}

	return notif.NotificationRule, nil
}

func (m Module) GetUserContactMethods(id string) []pager.ContactMethod {

	var contact ContactMethodQueryResult
	address := "https://api.pagerduty.com"
	path := "users/" + id + "/contact_methods"

	header := make(map[string]string)
	header["Accept"] = "application/vnd.pagerduty+json;version=2"
	header["Authorization"] = fmt.Sprintf("Token token=%s", m.Token)

	agent := httpreq.NewHTTPRequest()
	agent.Url = address
	agent.Path = path
	agent.Method = "GET"
	agent.Headers = header

	body, err := agent.DoReq()
	if err != nil {
		log.Println(err)
		agent.Debug()
		if body != nil {
			log.Println(string(*body))
		}
		return contact.ContactMethods
	}
	if err := json.Unmarshal(*body, &contact); err != nil {
		log.Println(err)
		log.Println(string(*body))
		return contact.ContactMethods
	}

	return contact.ContactMethods
}

func (m Module) SetDefaultNotification(id string) error {
	contact := m.GetUserContactMethods(id)
	mcontact := make(map[string]string)

	if len(contact) > 0 {
		for _, c := range contact {
			mcontact[c.Type] = c.ID
		}
	}

	cJ, _ := json.Marshal(contact)
	fmt.Println(string(cJ))

	setnot, err := m.CreateUserNotificationRules(id, pager.NotificationRule{
		Type:                "assignment_notification_rule",
		StartDelayInMinutes: 8,
		ContactMethod: pager.ContactMethod{
			ID:      mcontact["phone_contact_method"],
			Type:    "phone_contact_method",
			Summary: "Mobile",
		},

		Urgency: "high",
	})
	if err != nil {
		fmt.Println(err)
		return err
	}

	sn, _ := json.Marshal(setnot)
	fmt.Println(string(sn))
	return err
}

func (m Module) ListUser(userCondition int) []pager.User {

	var (
		offset int

		users []pager.User
	)

	for {
		lur, err := m.Client.ListUsers(pager.ListUsersOptions{
			APIListObject: pager.APIListObject{
				Limit:  100,
				Offset: uint(offset),
			},
		})
		if err != nil {
			panic(err)
		}

		for _, u := range lur.Users {

			switch userCondition {
			case ConstAllUser:
				users = append(users, u)
			case ConstUserHasNotBeenValidated:
				if len(u.ContactMethods) < 2 {
					users = append(users, u)
				}
			default:
				fmt.Println("Please set the condition !!!")
				return users
			}

		}

		if len(lur.Users) == 0 {
			break
		}
		offset += 100
	}

	return users
}

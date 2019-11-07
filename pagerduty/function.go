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
	path := "teams"

	header := make(map[string]string)
	header["Accept"] = HeaderAccept
	header["Authorization"] = fmt.Sprintf("Token token=%s", m.Token)
	params := url.Values{}
	params.Add("query", teamName)

	agent := httpreq.NewHTTPRequest()
	agent.Url = PagerDutyURL
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
		return user, err
	}

	return user, err
}

func (m Module) CreateUserNotificationRules(id string, detail pager.NotificationRule) (pager.NotificationRule, error) {
	notif := Notification{
		detail,
	}
	path := "users/" + id + "/notification_rules"

	header := make(map[string]string)
	header["Accept"] = HeaderAccept
	header["Authorization"] = fmt.Sprintf("Token token=%s", m.Token)

	agent := httpreq.NewHTTPRequest()
	agent.Url = PagerDutyURL
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
	path := "users/" + id + "/contact_methods"

	header := make(map[string]string)
	header["Accept"] = HeaderAccept
	header["Authorization"] = fmt.Sprintf("Token token=%s", m.Token)

	agent := httpreq.NewHTTPRequest()
	agent.Url = PagerDutyURL
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

	_, err := m.CreateUserNotificationRules(id, pager.NotificationRule{
		Type:                "assignment_notification_rule",
		StartDelayInMinutes: 16,
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

	return err
}

func (m Module) ListUser(query string, userCondition int) ([]pager.User, error) {

	var (
		offset int
		users  []pager.User
	)

	for {
		lur, err := m.Client.ListUsers(pager.ListUsersOptions{
			APIListObject: pager.APIListObject{
				Limit:  100,
				Offset: uint(offset),
			},
			Includes: []string{"contact_methods"},
			Query:    query,
		})
		if err != nil {
			return users, err
		}

		for _, u := range lur.Users {
			switch userCondition {
			case ConstAllUser:
				users = append(users, u)
			case ConstUserHasNotBeenValidated:
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

				if !(hasPhone && hasSMS && hasPushNotif) {
					users = append(users, u)
				}
			default:
				return users, fmt.Errorf("Please set the condition !!!")
			}

		}
		if len(lur.Users) == 0 {
			break
		}
		offset += 100
	}

	return users, nil
}

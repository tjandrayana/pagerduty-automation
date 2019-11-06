package pagerduty

import pager "github.com/PagerDuty/go-pagerduty"

type User struct {
	DetailUser pager.User `json:"user"`
}

type Contact struct {
	ContactMethod pager.ContactMethod `json:"contact_method,omitempty"`
}

type TeamQueryResult struct {
	Teams []pager.Team `json:"teams, omitempty"`
}

type NotificationRulesQueryResult struct {
	NotificationRules []pager.NotificationRule `json:"notification_rules,omitempty"`
}

type Notification struct {
	pager.NotificationRule `json:"notification_rule"`
}

type ContactMethodQueryResult struct {
	ContactMethods []pager.ContactMethod `json:"contact_methods,omitempty"`
}

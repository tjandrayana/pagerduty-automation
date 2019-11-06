package pagerduty

import pager "github.com/PagerDuty/go-pagerduty"

type Module struct {
	Token  string
	Client *pager.Client
}

func Init() Module {
	token := "pager_duty_token"
	return Module{
		Token:  token,
		Client: pager.NewClient(token),
	}
}

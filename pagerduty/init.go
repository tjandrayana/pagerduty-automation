package pagerduty

import pager "github.com/PagerDuty/go-pagerduty"

type Module struct {
	Token  string
	Client *pager.Client
}

func Init(token string) Module {
	return Module{
		Token:  token,
		Client: pager.NewClient(token),
	}
}

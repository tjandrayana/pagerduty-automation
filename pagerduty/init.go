package pagerduty

import pager "github.com/PagerDuty/go-pagerduty"

type Module struct {
	Token  string
	Client *pager.Client
}

func Init() Module {
	token := "CEs_sysoX1tm5Wur8S8h"
	return Module{
		Token:  token,
		Client: pager.NewClient(token),
	}
}

# PagerDuty Automation

This is an apps that used for lazy person. It help you to automate some features in PagerDuty. To use this apps make sure you have PagerDuty api token and set it to your os environment.
 
**Command :**
```sh
export PagerDutyToken=FillThisWithYourAPIToken
```

#### Get All Users
```sh
clear ; go run main.go -command=get_users
```

#### Get Miss Users (User Is Not Verified)
```sh
clear ; go run main.go -command=get_users
```

#### Create New User
```sh
clear ; go run main.go -command=create_user -name="Iron Man" -email="iron.man@email.com" -role=admin -job="Software Engineer"
```
For [PagerDuty Roles](https://api-reference.pagerduty.com/#!/Users/post_users)  you can follow this documentation.

#### Set User Notification Rule
```sh
clear ; go run main.go -command=set_rules -email="iron.man@email.com"
```
I set some default rules, but you can customize it.

### Finally

Special thanks for : 
    - [PagerDuty](https://www.pagerduty.com/)
    - [PagerDuty Lib](github.com/PagerDuty/go-pagerduty)

API References : https://api-reference.pagerduty.com/

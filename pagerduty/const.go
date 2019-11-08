package pagerduty

const (
	ConstAllUser = iota
	ConstUserHasNotBeenValidated
)

const (
	PagerDutyURL = "https://api.pagerduty.com"
	HeaderAccept = "application/vnd.pagerduty+json;version=2"
)

// Notification Rule Type
const (
	ConstAssignmentNotificationRule string = "assignment_notification_rule"
)

// Contact Method Constanta
const (
	ConstPushNotificationContactMethod string = "push_notification_contact_method"
	ConstPhoneContactMethod            string = "phone_contact_method"
	ConstSmsConstactMehtod             string = "sms_contact_method"
)

package main

import (
	"fmt"
	"log"

	pager "github.com/PagerDuty/go-pagerduty"
	"github.com/tjandrayana/pagerduty-automation/pagerduty"
)

func main() {
	m := pagerduty.Init()
	// TestCase1(m)
	TestCase2(m)
	// TestCase3(m)
}

// Create User
func TestCase1(m pagerduty.Module) {
	count := 20
	user, err := m.CreateUser(pager.User{
		Type:        "user",
		Name:        fmt.Sprintf("Fahrel %d", count),
		Email:       fmt.Sprintf("fahrel-%d@testing.com", count),
		Timezone:    "Asia/Jakarta",
		Role:        "limited_user",
		JobTitle:    "Buruh Ketik",
		Description: "Saya buruh ketik",
	})
	if err != nil {
		log.Println(err)
		panic(err)
	}

	fmt.Println(user)

}

// List of User who not open the link sent by admin
//  ConstAllUser 					=> Get All User
//  ConstUserHasNotBeenValidated 	=> Get User  Has Not Been Validated / Not Enter Phone Number / Not Install the Apps
// query can be a names or email or empty
func TestCase2(m pagerduty.Module) {
	query := "names or email addresses"
	users := m.ListUser(query, pagerduty.ConstAllUser)
	for i, u := range users {
		fmt.Printf("%d => %v\n", i+1, u.Email)
	}
}

// Set Notif Rule For the user
func TestCase3(m pagerduty.Module) {
	id := "XXXXX"
	if err := m.SetDefaultNotification(id); err != nil {
		log.Println(err)
		panic(err)
	}

}

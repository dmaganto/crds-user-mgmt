package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/slack-go/slack"
)

func addUserToSlackGroup(mail, slackGroup string) error {
	fmt.Println("Adding user to Slack group inside slack.go", mail, slackGroup)
	slackAPIToken := os.Getenv("SLACK_API_TOKEN")
	api := slack.New(slackAPIToken)
	// Step 1: Lookup User ID from email
	userID, err := getUserByEmail(mail)
	if err != nil {
		fmt.Printf("Error getting user: %s\n", err)
		return err
	}

	// Step 2: Lookup Group ID from group name
	groups, err := api.GetUserGroups(slack.GetUserGroupsOptionIncludeUsers(true))
	if err != nil {
		fmt.Printf("Error getting groups: %s\n", err)
		return err
	}
	var group slack.UserGroup
	for _, g := range groups {
		if g.Handle == slackGroup {
			group = g
			break
		}
	}
	if group.ID == "" {
		fmt.Println("Group not found")
		os.Exit(1)
	}

	// Get current users of the group
	currentUsers := group.Users
	// Append new user to the list
	currentUsers = append(currentUsers, userID)
	// Users to add as a string
	newUsers := strings.Join(currentUsers, ",")
	// Step 3: Add user to group
	_, err = api.UpdateUserGroupMembers(group.ID, newUsers)
	if err != nil {
		fmt.Printf("Error adding user to group: %s\n", err)
		return err
	}

	fmt.Printf("User %s successfully added to group %s\n", mail, slackGroup)
	return nil
}

func getUserByEmail(mail string) (string, error) {
	fmt.Println("Getting user by email inside slack.go", mail)
	slackAPIToken := os.Getenv("SLACK_API_TOKEN")
	api := slack.New(slackAPIToken)
	user, err := api.GetUserByEmail(mail)
	if err != nil {
		fmt.Printf("Error getting user: %s\n", err)
		return "", err
	}

	return user.ID, nil
}

func createSlackChannel(channelName string) (string, error) {
	fmt.Println("Creating Slack channel inside slack.go", channelName)
	slackAPIToken := os.Getenv("SLACK_API_TOKEN")
	api := slack.New(slackAPIToken)

	channel, err := api.CreateConversation(slack.CreateConversationParams{
		ChannelName: channelName,
		IsPrivate:   false,
		TeamID:      "",
	})
	if err != nil {
		fmt.Printf("Failed to create channel: %s Reason: %s\n", channelName, err)
		return "", err
	}

	fmt.Printf("Channel %s successfully created\n", channel.Name)
	return channel.ID, nil
}

// Function to add a channel to a team
func addChannelToTeam(channelID, teamName string) error {
	// get the current channels in the team
	slackAPIToken := os.Getenv("SLACK_API_TOKEN")
	api := slack.New(slackAPIToken)

	gid, currentChannels, err := getTeamIDandChannels(teamName)
	if err != nil {
		fmt.Printf("Error getting team: %s\n", err)
		return err
	}
	newChannels := append(currentChannels, channelID)
	//groupID of testautomation: S066X9BS3K4
	api.UpdateUserGroup(gid, slack.UpdateUserGroupsOptionChannels(newChannels))
	return nil
}

// Get the id of a team and its channels based on @teamName
func getTeamIDandChannels(teamName string) (string, []string, error) {
	slackAPIToken := os.Getenv("SLACK_API_TOKEN")
	api := slack.New(slackAPIToken)
	userGroups, err := api.GetUserGroups()
	if err != nil {
		fmt.Printf("Error fetching user groups: %s\n", err)
		return "", nil, err
	}

	var userGroupID string
	fmt.Printf("User groups: %v\n", userGroups)
	for _, userGroup := range userGroups {
		if userGroup.Handle == teamName {
			userGroupID = userGroup.ID
			return userGroupID, userGroup.Prefs.Channels, nil
		}
	}
	return "", nil, fmt.Errorf("Group %s not found\n", teamName)

}

// create a new user group
func createNewUserGroup(teamName string, userIDs, channelIDs []string) (string, error) {
	slackAPIToken := os.Getenv("SLACK_API_TOKEN")
	api := slack.New(slackAPIToken)
	userGroup := slack.UserGroup{
		Name:        "Test " + teamName,
		Handle:      teamName,
		IsUserGroup: true,
		Description: "This is a test group",
		Users:       userIDs, //David U01U0J0AMCM
		Prefs: slack.UserGroupPrefs{
			Channels: channelIDs, //monitoring-testautomation-non-critical
		},
	}
	group, err := api.CreateUserGroup(userGroup)
	if err != nil {
		fmt.Printf("Failed to create user group: %s\n", err)
		return "", err
	}
	// For some reason, members are not added to the group with the CreateUserGroup method, adding them here
	_, err = api.UpdateUserGroupMembers(group.ID, strings.Join(userIDs, ","))
	if err != nil {
		fmt.Printf("Error adding user to group: %s\n", err)
		return group.ID, err
	}
	fmt.Printf("User Group created: %+v\n", group)
	return group.ID, nil
}

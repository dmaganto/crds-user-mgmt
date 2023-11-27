package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type Team struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type TeamsResponse struct {
	Data []Team `json:"data"`
}

type Service struct {
	Name string `json:"name"`
	Team Team   `json:"team"`
}

type ServicesResponse struct {
	Data []Service `json:"data"`
}

// Function to check if a service already exists
// Used in createService() to prevent duplicate services
func serviceAlreadyExists(serviceName, teamId string) (bool, error) {
	query := "name:" + url.QueryEscape(serviceName)
	url := "https://api.opsgenie.com/v1/services?query=" + query + "&limit=10&sort=name"
	opsgenieApiToken := os.Getenv("OPSGENIE_API_TOKEN")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", opsgenieApiToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("error reading response body: %v", err)
	}
	var servicesResponse ServicesResponse
	err = json.Unmarshal(body, &servicesResponse)
	if err != nil {
		return false, fmt.Errorf("error unmarshalling response: %v", err)
	}
	for _, service := range servicesResponse.Data {
		if service.Name == serviceName {
			return true, nil
		}
	}

	return false, nil
}

// Function to create a new service if it doesn't exists
func createService(service, teamId string) (string, error) {
	serviceExist, err := serviceAlreadyExists(service, teamId)
	fmt.Printf("Service %s exist: %v\n", service, serviceExist)
	if err != nil {
		return "", fmt.Errorf("error checking if service exists: %v", err)
	} else if serviceExist {
		return "", fmt.Errorf("error service %s already exists: ", service)
	}
	opsgenieApiToken := os.Getenv("OPSGENIE_API_TOKEN")
	url := "https://api.opsgenie.com/v1/services"

	var jsonStr = []byte(fmt.Sprintf(`{  
		"name": "%s",  
		"description": "Only for test new automation",  
		"teamId": "%s"  
	}`, service, teamId))

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", opsgenieApiToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}
	defer resp.Body.Close()
	return "ok", nil
}

// Function to get the teamId from the team name
func getTeamIdFromName(teamName string) (string, error) {
	url := "https://api.opsgenie.com/v2/teams"
	opsgenieApiToken := os.Getenv("OPSGENIE_API_TOKEN")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", opsgenieApiToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	var teamsResponse TeamsResponse
	err = json.Unmarshal(body, &teamsResponse)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling response: %v", err)
	}

	for _, team := range teamsResponse.Data {
		if team.Name == teamName {
			return team.Id, nil
		}
	}

	return "", fmt.Errorf("team not found: %s", teamName)
}

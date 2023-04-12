package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const (
	gitlabToken = ""
	gitlabAPI   = "https://gitlab.com/api/v4/projects/%s/pipeline?ref=develop"
)

type PipelineRequest struct {
	Ref string `json:"ref"`
}

func main() {
	// check if projects are passed as arguments
	if len(os.Args) < 2 {
		fmt.Println("Please provide at least one project ID as an argument")
		os.Exit(1)
	}

	// loop over projects passed as arguments
	for _, projectID := range os.Args[1:] {
		// create pipeline request payload
		payload := PipelineRequest{Ref: "develop"}
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			fmt.Printf("Failed to marshal JSON payload for project %s: %v\n", projectID, err)
			continue
		}

		// create GitLab API request
		url := fmt.Sprintf(gitlabAPI, projectID)
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonPayload))
		if err != nil {
			fmt.Printf("Failed to create API request for project %s: %v\n", projectID, err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("PRIVATE-TOKEN", gitlabToken)

		// send GitLab API request
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Failed to send API request for project %s: %v\n", projectID, err)
			continue
		}

		// check response status code
		if resp.StatusCode != http.StatusCreated {
			fmt.Printf("Failed to trigger pipeline for project %s: %s\n", projectID, resp.Status)
			continue
		}

		fmt.Printf("Pipeline triggered successfully for project %s\n", projectID)
	}
}

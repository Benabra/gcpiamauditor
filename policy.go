package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

type Binding struct {
	Role    string   `json:"role"`
	Members []string `json:"members"`
}

type IAMPolicy struct {
	Bindings []Binding `json:"bindings"`
}

func getIAMPolicy(projectID string) (*IAMPolicy, error) {
	cmd := exec.Command("gcloud", "projects", "get-iam-policy", projectID, "--format=json")
	var out bytes.Buffer
	cmd.Stdout = &out
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%s: %s", err, errBuf.String())
	}
	var iamPolicy IAMPolicy
	if err := json.Unmarshal(out.Bytes(), &iamPolicy); err != nil {
		return nil, err
	}
	return &iamPolicy, nil
}

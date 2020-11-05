package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type change struct {
	Actions []string    `json:"actions"`
	After   interface{} `json:"after"`
}

type resourceChange struct {
	Type         string `json:"type"`
	Name         string `json:"name"`
	ProviderName string `json:"provider_name"`
	Change       change
}

type tfplan struct {
	ResourceChanges []resourceChange `json:"resource_changes"`
}

func (r *resourceChange) ResourceType() string {
	return strings.ToUpper(fmt.Sprintf("TERRAFORM::RESOURCE::%s", strings.Replace(r.Type, "_", "", -1)))
}

func (r *resourceChange) Config() (string, error) {
	cfg, err := json.MarshalIndent(r.Change.After, "  ", "  ")
	if err != nil {
		return "", err
	}

	return string(cfg), nil
}

func realMain() int {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "error: no filename given as parameter\nUsage: %s [PLANFILE.json]\n", os.Args[0])
		return 1
	}

	fileName := os.Args[1]

	planBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed reading Terraform plan %s: %s\n", fileName, err)
		return 1
	}

	var plan tfplan
	if err := json.Unmarshal(planBytes, &plan); err != nil {
		fmt.Fprintf(os.Stderr, "failed parsing Terraform plan %s: %s\n", fileName, err)
		return 1
	}

	for _, r := range plan.ResourceChanges {
		hasUpdate := false
		for _, action := range r.Change.Actions {
			if action == "update" {
				hasUpdate = true
			}
		}
		if !hasUpdate {
			continue
		}

		cfg, err := r.Config()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed parsing config for a resource, got error: %s\n", fileName, err)
			return 1
		}
		fmt.Printf("aws configservice put-resource-config --resource-type %s --resource-id %s --schema-version-id 00000001 --configuration '%s'\n", r.ResourceType(), r.Name, cfg)
	}

	return 0
}

func main() {
	os.Exit(realMain())
}

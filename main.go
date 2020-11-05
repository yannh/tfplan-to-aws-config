package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)


type change struct {
	Actions []string `json:"actions"`
	After interface{} `json:"after"`
}

type resourceChange struct {
	Name string          `json:"name"`
	ProviderName string  `json:"provider_name"`
	Change change
}

type tfplan struct {
	ResourceChanges []resourceChange  `json:"resource_changes"`
}

func realMain () int {
	fileName := "tfplan.json"

	planBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed reading Terraform plan %s: %s", fileName, err)
		return 1
	}

	var plan tfplan
	if err := json.Unmarshal(planBytes, &plan); err != nil {
		fmt.Fprintf(os.Stderr, "failed parsing Terraform plan %s: %s", fileName, err)
		return 1
	}

	fmt.Printf("%+v", plan.ResourceChanges)

	return 0
}

func main() {
	os.Exit(realMain())
}
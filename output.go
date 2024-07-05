package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type RoleInfo struct {
	Project    string   `json:"project" yaml:"project"`
	MemberType string   `json:"member_type" yaml:"member_type"`
	User       string   `json:"user" yaml:"user"`
	Roles      []string `json:"roles" yaml:"roles"`
}

func outputResults(userRoles map[string]map[string]map[string]map[string]bool, outputFormat string) {
	rolesList := []RoleInfo{}
	for project, memberTypes := range userRoles {
		for memberType, users := range memberTypes {
			for user, rolesMap := range users {
				var roles []string
				for role := range rolesMap {
					roles = append(roles, role)
				}
				rolesList = append(rolesList, RoleInfo{
					Project:    project,
					MemberType: memberType,
					User:       user,
					Roles:      roles,
				})
			}
		}
	}

	switch outputFormat {
	case "json":
		outputJSON(rolesList)
	case "yaml":
		outputYAML(rolesList)
	case "csv":
		outputCSV(rolesList)
	default:
		outputTable(rolesList)
	}
}

func outputJSON(rolesList []RoleInfo) {
	jsonData, err := json.MarshalIndent(rolesList, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}
	fmt.Println(string(jsonData))
}

func outputYAML(rolesList []RoleInfo) {
	yamlData, err := yaml.Marshal(rolesList)
	if err != nil {
		log.Fatalf("Failed to marshal YAML: %v", err)
	}
	fmt.Println(string(yamlData))
}

func outputCSV(rolesList []RoleInfo) {
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	writer.Write([]string{"PROJECT", "MEMBER TYPE", "USER", "ROLES"})
	for _, roleInfo := range rolesList {
		writer.Write([]string{roleInfo.Project, roleInfo.MemberType, roleInfo.User, strings.Join(roleInfo.Roles, ", ")})
	}
}

func outputTable(rolesList []RoleInfo) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleColoredBright)

	t.AppendHeader(table.Row{"PROJECT", "MEMBER TYPE", "USER", "ROLES"})

	for _, roleInfo := range rolesList {
		t.AppendRow(table.Row{roleInfo.Project, roleInfo.MemberType, roleInfo.User, strings.Join(roleInfo.Roles, ", ")})
	}

	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Colors: text.Colors{text.FgCyan}},
		{Number: 2, Colors: text.Colors{text.FgGreen}},
		{Number: 3, Colors: text.Colors{text.FgYellow}},
		{Number: 4, Colors: text.Colors{text.FgMagenta}},
	})

	t.Render()
}

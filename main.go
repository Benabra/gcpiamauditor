package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/progress"
)

func main() {
	// Define and parse the projectFile, exclude, services, member, and output arguments
	var projectFile string
	var excludeTypes string
	var services string
	var memberRoles string
	var outputFormat string
	flag.StringVar(&projectFile, "projectFile", "", "Path to the file containing GCP project IDs")
	flag.StringVar(&excludeTypes, "exclude", "", "Comma-separated list of member types to exclude (e.g., serviceAccount, user, group)")
	flag.StringVar(&services, "services", "", "Comma-separated list of GCP services to check (e.g., bigquery)")
	flag.StringVar(&memberRoles, "member", "", "Comma-separated list of IAM roles to include (e.g., admin, owner, editor)")
	flag.StringVar(&outputFormat, "output", "table", "Output format: table, json, yaml, or csv")
	flag.Parse()

	if projectFile == "" {
		log.Fatalf("Usage: %s -projectFile=projects.txt [-exclude=serviceAccount,user,group] [-services=bigquery] [-member=admin,owner,editor] [-output=table,json,yaml,csv]", os.Args[0])
	}

	projectIDs, err := readProjectFile(projectFile)
	if err != nil {
		log.Fatalf("Failed to read project file: %v", err)
	}

	excludeSet := make(map[string]bool)
	if excludeTypes != "" {
		for _, t := range strings.Split(excludeTypes, ",") {
			excludeSet[t] = true
		}
	}

	serviceSet := make(map[string]bool)
	if services != "" {
		for _, s := range strings.Split(services, ",") {
			serviceSet[s] = true
		}
	}

	memberSet := make(map[string]bool)
	if memberRoles != "" {
		for _, m := range strings.Split(memberRoles, ",") {
			memberSet[strings.ToLower(m)] = true
		}
	}

	// Initialize the progress writer
	pw := progress.NewWriter()
	pw.SetAutoStop(true)
	pw.SetTrackerPosition(progress.PositionRight)
	pw.SetMessageWidth(50)
	pw.SetTrackerLength(25)
	pw.ShowOverallTracker(true)
	pw.ShowTime(true)
	pw.ShowTracker(true)

	go pw.Render()

	// Map to store user roles across projects
	userRoles := make(map[string]map[string]map[string]map[string]bool)

	// Create a tracker for each project
	var trackers []*progress.Tracker
	for _, projectID := range projectIDs {
		tracker := &progress.Tracker{Message: fmt.Sprintf("Processing project: %s", projectID), Total: 1}
		trackers = append(trackers, tracker)
		pw.AppendTracker(tracker)
	}

	for i, projectID := range projectIDs {
		tracker := trackers[i]
		// Update tracker progress
		tracker.Increment(1)

		// Get IAM policy for the project
		iamPolicy, err := getIAMPolicy(projectID)
		if err != nil {
			tracker.MarkAsErrored()
			tracker.Message = fmt.Sprintf("Failed to get IAM policy for project %s: %v", projectID, err)
			continue
		}

		// Collect admin roles assigned to user accounts
		collectAdminRoles(userRoles, projectID, iamPolicy, excludeSet, serviceSet, memberSet)

		// Mark the tracker as done
		tracker.MarkAsDone()
		time.Sleep(100 * time.Millisecond) // simulate work
	}

	pw.Stop()

	// Print the results
	outputResults(userRoles, outputFormat)

	fmt.Println("Admin roles have been listed for all user accounts (excluding specified member types) in the specified projects.")
}

func readProjectFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var projectIDs []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		projectID := strings.TrimSpace(scanner.Text())
		if projectID != "" {
			projectIDs = append(projectIDs, projectID)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return projectIDs, nil
}

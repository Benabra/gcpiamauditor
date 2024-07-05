package main

import (
	"strings"
)

func collectAdminRoles(userRoles map[string]map[string]map[string]map[string]bool, projectID string, iamPolicy *IAMPolicy, excludeSet, serviceSet, memberSet map[string]bool) {
	if _, exists := userRoles[projectID]; !exists {
		userRoles[projectID] = make(map[string]map[string]map[string]bool)
	}

	for _, binding := range iamPolicy.Bindings {
		roleLower := strings.ToLower(binding.Role)
		if len(memberSet) > 0 {
			matched := false
			for memberRole := range memberSet {
				if strings.Contains(roleLower, memberRole) {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}

		if strings.Contains(roleLower, "admin") || containsMemberRole(memberSet, roleLower) {
			for _, member := range binding.Members {
				memberType, user := parseMember(member)
				if excludeSet[memberType] {
					continue
				}
				if len(serviceSet) > 0 {
					matched := false
					for service := range serviceSet {
						if strings.Contains(roleLower, strings.ToLower(service)) {
							matched = true
							break
						}
					}
					if !matched {
						continue
					}
				}
				if _, exists := userRoles[projectID][memberType]; !exists {
					userRoles[projectID][memberType] = make(map[string]map[string]bool)
				}
				if _, exists := userRoles[projectID][memberType][user]; !exists {
					userRoles[projectID][memberType][user] = make(map[string]bool)
				}
				userRoles[projectID][memberType][user][binding.Role] = true
			}
		}
	}
}

func containsMemberRole(memberSet map[string]bool, role string) bool {
	for memberRole := range memberSet {
		if strings.Contains(role, memberRole) {
			return true
		}
	}
	return false
}

func parseMember(member string) (string, string) {
	parts := strings.Split(member, ":")
	if len(parts) > 1 {
		return parts[0], parts[1]
	}
	return "", member
}

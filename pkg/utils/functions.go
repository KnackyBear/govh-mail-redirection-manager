package utils

import (
	"fmt"
	"strings"
)

func GetDomain(mail string) (string, error) {
	if mail == "" {
		return "", fmt.Errorf("email required")
	}

	email := strings.Split(mail, "@")
	if len(email) == 2 {
		subdomain := strings.Split(email[1], ".")
		if len(subdomain) > 1 {
			return fmt.Sprintf("%s.%s", subdomain[len(subdomain)-2], subdomain[len(subdomain)-1]), nil
		} else {
			return "", fmt.Errorf("domain not found in %v", mail)
		}
	} else {
		return "", fmt.Errorf("domain not found in %v", mail)
	}
}

func StrArrayContains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	"github.com/mailru/easyjson" //nolint: depguard
)

type User struct {
	Email string `json:"Email"` //nolint
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return countDomains(r, domain)
}

func countDomains(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)

	var user User
	var emailDomain string

	for scanner.Scan() {
		line := scanner.Bytes()
		err := easyjson.Unmarshal(line, &user)
		if err != nil {
			return nil, err
		}

		emailDomain = strings.ToLower(user.Email)
		if !strings.Contains(emailDomain, domain) {
			continue
		}
		atIndex := strings.Index(emailDomain, "@")
		if atIndex != -1 {
			result[emailDomain[atIndex+1:]]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

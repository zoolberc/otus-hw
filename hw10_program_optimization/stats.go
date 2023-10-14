package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/valyala/fastjson"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUserEmails(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

func getUserEmails(r io.Reader) ([]string, error) {
	emails := make([]string, 0)
	scanner := bufio.NewScanner(r)
	p := fastjson.Parser{}
	for scanner.Scan() {
		l := scanner.Bytes()
		v, err := p.ParseBytes(l)
		if err != nil {
			return nil, err
		}
		email := string(v.GetStringBytes("Email"))
		emails = append(emails, email)
	}
	return emails, nil
}

func countDomains(emails []string, domain string) (DomainStat, error) {
	result := make(DomainStat)
	for _, email := range emails {
		domainName := strings.ToLower(strings.SplitN(email, "@", 2)[1])
		matched := strings.Contains(email, "."+domain)
		if matched {
			count := result[domainName]
			count++
			result[domainName] = count
		}
	}
	return result, nil
}

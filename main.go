package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

func main() {
	const length = 3

	in := make(chan string)
	out := make(chan domainStatus)

	go func() {
		perm("abcdefghijklmnopqrstuvwxyz0123456789", length, func(s string) {
			in <- s + ".dev"
		})
		close(in)
	}()

	workers := 50

	var wg sync.WaitGroup
	wg.Add(workers)

	for n := 0; n < workers; n++ {
		go func() {
			defer wg.Done()

			for domain := range in {
				status, err := check(domain)
				if err != nil {
					log.Fatal(err)
				}

				out <- status
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	bufStdout := bufio.NewWriter(os.Stdout)
	enc := json.NewEncoder(bufStdout)
	for status := range out {
		enc.Encode(status)
	}
	bufStdout.Flush()
}

type domainStatus struct {
	Domain    string `json:"domain"`
	Available bool   `json:"available"`
	Reason    string `json:"reason,omitempty"`
	Status    string `json:"status,omitempty"`
	Tier      string `json:"tier,omitempty"`
}

func check(domain string) (domainStatus, error) {
	v := make(url.Values)
	v.Set("domain", domain)

	u, err := url.Parse("https://pubapi-dot-domain-registry.appspot.com/check")
	if err != nil {
		return domainStatus{}, fmt.Errorf("parse: %v", err)
	}

	u.RawQuery = v.Encode()

	c := &http.Client{Timeout: 30 * time.Second}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return domainStatus{}, fmt.Errorf("new request: %v", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return domainStatus{}, fmt.Errorf("query: %v", err)
	}

	var status domainStatus

	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return domainStatus{}, fmt.Errorf("decode: %v", err)
	}

	if status.Status != "success" {
		return domainStatus{}, fmt.Errorf("[%s] status = %q", domain, status.Status)
	}

	status.Domain = domain
	status.Status = ""

	return status, nil
}

// perm generates all permutations, with repeated characters, of s of the given
// length, calling f for each.
func perm(s string, length int, f func(s string)) {
	if s == "" || length == 0 {
		return
	}

	p := make([]rune, length)

	var rec func(idx int)

	rec = func(idx int) {
		for _, r := range s {
			p[idx] = r
			if idx == length-1 {
				f(string(p))
			} else {
				rec(idx + 1)
			}
		}
	}

	rec(0)
}

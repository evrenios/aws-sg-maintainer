package maintainer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"gopkg.in/fatih/set.v0"
)

type githubResponse struct {
	Hooks []string `json:"hooks"`
}

var (
	client = &http.Client{Timeout: 10 * time.Second}
)

const (
	githubEndpoint = "https://api.github.com/meta"
	cloudflareIPV4 = "https://www.cloudflare.com/ips-v4"
	cloudflareIPV6 = "https://www.cloudflare.com/ips-v6"
)

func getGithubIPBlocks() (*set.SetNonTS, error) {
	resp, err := client.Get(githubEndpoint)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	by, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d for endpoint %s", resp.StatusCode, githubEndpoint)
	}

	response := new(githubResponse)
	if err = json.Unmarshal(by, &response); err != nil {
		return nil, err
	}

	ipBlocks := set.NewNonTS()

	for _, ip := range response.Hooks {
		ipBlocks.Add(ip)
	}
	return ipBlocks, nil
}

func getAllCloudFlareIPBlocks() (*set.SetNonTS, error) {
	ipv4, err := getCloudFlareIPBlocks(cloudflareIPV4)
	if err != nil {
		return nil, err
	}
	ipv6, err := getCloudFlareIPBlocks(cloudflareIPV6)
	if err != nil {
		return nil, err
	}
	ipv6.Merge(ipv4)
	return ipv6, nil
}

func getCloudFlareIPBlocks(url string) (*set.SetNonTS, error) {
	ipBlocks := set.NewNonTS()

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	by, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d for endpoint %s", resp.StatusCode, url)
	}

	allIPBlocks := strings.Split(string(by), "\n")

	for _, ip := range allIPBlocks {
		if len(ip) == 0 {
			continue
		}
		ipBlocks.Add(ip)
	}
	return ipBlocks, nil
}

func getServiceIPBlocks(service ServiceProvider) (*set.SetNonTS, error) {
	switch service {
	case Github:
		return getGithubIPBlocks()
	case Cloudflare:
		return getAllCloudFlareIPBlocks()
	}
	return nil, errors.New("unexpected service type")
}

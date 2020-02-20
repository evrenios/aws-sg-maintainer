package main

import (
	"fmt"

	"github.com/evrenios/aws-sg-maintainer/maintainer"
)

func main() {

	config := &maintainer.MaintainerConfig{
		ReadOnly:  true,
		AWSRegion: "eu-west-1",
		Services: []*maintainer.ServiceConfig{
			&maintainer.ServiceConfig{
				Service:          maintainer.Cloudflare,
				SecurityGroupIDs: []string{"sg-0dfc4936e4ed96ff5"},
				Ports:            []int64{443},
			},
			&maintainer.ServiceConfig{
				Service:          maintainer.Github,
				SecurityGroupIDs: []string{"sg-06aefdb239c0aa541"},
				Ports:            []int64{80},
			},
		},
	}

	if err := maintainer.MaintenanceTime(config); err != nil {
		panic(err)
	}

	fmt.Println("done")
}

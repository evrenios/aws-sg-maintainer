package maintainer

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/fatih/color"
	"gopkg.in/fatih/set.v0"
)

// ServiceProvider alias for supported services
// currently it's Github and Cloudflare
type ServiceProvider = string

var (
	Github     ServiceProvider = "github"
	Cloudflare ServiceProvider = "cloudflare"
)

type ServiceConfig struct {
	Service          ServiceProvider
	SecurityGroupIDs []string
	Ports            []int64
}
type MaintainerConfig struct {
	Services  []*ServiceConfig
	ReadOnly  bool
	AWSRegion string
	Ec2Svc    *ec2.EC2
}

func MaintenanceTime(config *MaintainerConfig) error {
	if config == nil {
		return errors.New("config can not be nil")
	}
	if len(config.Services) == 0 {
		return errors.New("there is no service to process")
	}

	if config.Ec2Svc == nil {
		config.Ec2Svc = ec2.New(session.New(), &aws.Config{Region: aws.String(config.AWSRegion)})
	}

	allSecurityGroups, err := getAllSecurityGroups(config.Ec2Svc)
	if err != nil {
		return fmt.Errorf("failed to describe security groups with error: %s", err.Error())
	}

	for _, serviceConfig := range config.Services {

		for _, securityGroupID := range serviceConfig.SecurityGroupIDs {
			sg, found := getServiceSG(allSecurityGroups, securityGroupID)
			if !found {
				return fmt.Errorf("security group with ID %s does not exist", securityGroupID)
			}

			serviceIPBlocks, err := getServiceIPBlocks(serviceConfig.Service)
			if err != nil {
				return err
			}

			for _, port := range serviceConfig.Ports {
				sgIPBlocks := getAllIPBlocksOfSgForPort(sg, port)

				ipBlocksToAdd := set.Difference(serviceIPBlocks, sgIPBlocks)
				ipBlocksToRemove := set.Difference(sgIPBlocks, serviceIPBlocks)

				color.Green("IP Blocks to add %s for service %s  ", ipBlocksToAdd.List(), serviceConfig.Service)
				color.Red("IP Blocks to remove %s for service %s  ", ipBlocksToRemove.List(), serviceConfig.Service)

				if config.ReadOnly {
					color.Cyan("!! THIS IS READ ONLY MODE, NO CHANGE HAS BEEN MADE !!")
					continue
				}

				if ipBlocksToAdd.Size() > 0 {
					if err := addIPBlocks(config.Ec2Svc, sg, set.StringSlice(ipBlocksToAdd), port); err != nil {
						return err
					}
				}

				if ipBlocksToRemove.Size() > 0 {
					if err := removeIPBlocks(config.Ec2Svc, sg, set.StringSlice(ipBlocksToRemove), port); err != nil {
						return err
					}

				}

			}
		}
	}

	return nil
}

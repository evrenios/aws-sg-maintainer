package maintainer

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"gopkg.in/fatih/set.v0"
)

const (
	ipProtocol = "tcp"
)

func getAllSecurityGroups(ec2Svc *ec2.EC2) ([]*ec2.SecurityGroup, error) {
	input := &ec2.DescribeSecurityGroupsInput{
		GroupNames: []*string{},
	}

	result, err := ec2Svc.DescribeSecurityGroups(input)
	if err != nil {
		return nil, err
	}
	return result.SecurityGroups, nil
}

func getServiceSG(securityGroups []*ec2.SecurityGroup, sgID string) (*ec2.SecurityGroup, bool) {
	for _, sg := range securityGroups {
		if *sg.GroupId == sgID {
			return sg, true
		}

	}
	return nil, false
}

func getAllIPBlocksOfSgForPort(sg *ec2.SecurityGroup, port int64) *set.SetNonTS {
	allIPBlocks := set.NewNonTS()
	for _, ipPermission := range sg.IpPermissions {
		for _, rr := range ipPermission.IpRanges {
			if port == *ipPermission.FromPort {
				allIPBlocks.Add(*rr.CidrIp)
			}

		}
		for _, rr := range ipPermission.Ipv6Ranges {
			if port == *ipPermission.FromPort {
				allIPBlocks.Add(*rr.CidrIpv6)
			}
		}
	}
	return allIPBlocks
}

func removeIPBlocks(ec2Svc *ec2.EC2, sg *ec2.SecurityGroup, ipBLocks []string, port int64) error {
	for _, ipBLock := range ipBLocks {
		_, err := ec2Svc.RevokeSecurityGroupIngress(
			&ec2.RevokeSecurityGroupIngressInput{
				CidrIp:     aws.String(ipBLock),
				GroupId:    sg.GroupId,
				IpProtocol: aws.String(ipProtocol),
				FromPort:   aws.Int64(port),
				ToPort:     aws.Int64(port),
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func addIPBlocks(ec2Svc *ec2.EC2, sg *ec2.SecurityGroup, ipBLocks []string, port int64) error {
	var err error
	for _, ipBLock := range ipBLocks {
		if strings.Contains(ipBLock, ":") {
			// ipv6
			_, err = ec2Svc.AuthorizeSecurityGroupIngress(
				&ec2.AuthorizeSecurityGroupIngressInput{
					GroupId: sg.GroupId,
					IpPermissions: []*ec2.IpPermission{
						{
							FromPort:   aws.Int64(port),
							IpProtocol: aws.String(ipProtocol),
							Ipv6Ranges: []*ec2.Ipv6Range{
								{
									CidrIpv6: aws.String(fmt.Sprint(ipBLock)),
								},
							},
							ToPort: aws.Int64(port),
						},
					},
				},
			)

		} else {
			// ipv4
			_, err = ec2Svc.AuthorizeSecurityGroupIngress(
				&ec2.AuthorizeSecurityGroupIngressInput{
					CidrIp:     aws.String(fmt.Sprint(ipBLock)),
					GroupId:    sg.GroupId,
					IpProtocol: aws.String(ipProtocol),
					FromPort:   aws.Int64(port),
					ToPort:     aws.Int64(port),
				},
			)
		}

		if err != nil {
			return err
		}

	}
	return nil
}

# AWS SG Maintainer 

AWS SG Maintainer fetches the ip blocks for `Github` and `Cloudflare` so the new ones will be added to your security groups and the old ones will be removed.

## Quick Example

```go
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
				SecurityGroupIDs: []string{"sg-1dfc1936b4en96mf5"},
				Ports:            []int64{443},
			},
			&maintainer.ServiceConfig{
				Service:          maintainer.Github,
				SecurityGroupIDs: []string{"sg-04aerdp226c4ai541"},
				Ports:            []int64{80},
			},
		},
	}

	if err := maintainer.MaintenanceTime(config); err != nil {
		panic(err)
	}
}
```


### Features
- You can pass your own ec2 service to the `MaintainerConfig` 
- if the `Ec2Svc` is empty, program tries to initialize the ec2 service from the default aws config
-  It has a `ReadOnly` mode so you can just see what will it change before applying

## Todo 

* Add new service ( you can also create issues for them with links or create PR)

## We'r Hiring 

 * This tool and many more of them are being created on a daily basis at Insider. If you want to join, [apply](https://useinsider.com/career/)

## License

The MIT License (MIT) - see [`LICENSE.md`](https://github.com/evrenios/aws-sg-maintainer/blob/master/LICENSE.md) for more details


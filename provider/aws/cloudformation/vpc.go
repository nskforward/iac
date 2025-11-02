package cloudformation

import (
	"encoding/json"
)

type VPC struct {
	Tagger

	id string

	// The IPv4 network range for the VPC, in CIDR notation. For example, 10.0.0.0/16
	// We modify the specified CIDR block to its canonical form
	// for example, if you specify 100.68.0.18/18, we modify it to 100.68.0.0/18.
	ipv4Cidr string

	// Indicates whether the instances launched in the VPC get DNS hostnames
	// If enabled, instances in the VPC get DNS hostnames
	// otherwise, they do not. Disabled by default for nondefault VPCs.
	enableDnsHostnames bool

	// Indicates whether the DNS resolution is supported for the VPC.
	// If enabled, queries to the Amazon provided DNS server at the 169.254.169.253 IP address,
	// or the reserved IP address at the base of the VPC network range "plus two" succeed.
	// If disabled, the Amazon provided DNS service in the VPC that resolves public DNS hostnames to IP addresses is not enabled.
	// Enabled by default.
	enableDnsSupport bool
}

type CFVPS struct {
	Type       string          `json:"Type"`
	Properties CFVPSProperties `json:"Properties"`
}

type CFVPSProperties struct {
	CidrBlock          string `json:"CidrBlock"`
	EnableDnsHostnames bool   `json:"EnableDnsHostnames"`
	EnableDnsSupport   bool   `json:"EnableDnsSupport"`
	InstanceTenancy    string `json:"InstanceTenancy,omitempty"`
	Ipv4IpamPoolId     string `json:"Ipv4IpamPoolId,omitempty"`
	Ipv4NetmaskLength  int    `json:"Ipv4NetmaskLength,omitempty"`
	Tags               []Tag  `json:"Tags"`
}

func NewVPC(id string, ipv4Cidr string, enableDnsHostnames, enableDnsSupport bool) *VPC {
	vpc := &VPC{
		id:                 id,
		ipv4Cidr:           ipv4Cidr,
		enableDnsHostnames: enableDnsHostnames,
		enableDnsSupport:   enableDnsSupport,
	}

	vpc.AddTag("Name", id)

	return vpc
}

func (vpc *VPC) Serialize() SerializeResult {
	data, err := json.Marshal(CFVPS{
		Type: "AWS::EC2::VPC",
		Properties: CFVPSProperties{
			CidrBlock:          vpc.ipv4Cidr,
			EnableDnsHostnames: vpc.enableDnsHostnames,
			EnableDnsSupport:   vpc.enableDnsSupport,
			Tags:               vpc.tags,
		},
	})
	if err != nil {
		panic(err)
	}

	return SerializeResult{
		ResourceID:   vpc.id,
		ResourceBody: data,
	}
}

package template

var DefsExample = map[string]Definition{
	"createsubnet": {
		Action:         "create",
		Entity:         "subnet",
		Api:            "ec2",
		RequiredParams: []*Param{{Name: "cidr"}, {Name: "vpc"}},
		ExtraParams:    []*Param{{Name: "availabilityzone"}, {Name: "name"}},
	},
	"updatesubnet": {
		Action:         "update",
		Entity:         "subnet",
		Api:            "ec2",
		RequiredParams: []*Param{{Name: "id"}},
		ExtraParams:    []*Param{{Name: "public"}},
	},
	"createinstance": {
		Action:         "create",
		Entity:         "instance",
		Api:            "ec2",
		RequiredParams: []*Param{{Name: "count"}, {Name: "image"}, {Name: "subnet"}, {Name: "type"}},
		ExtraParams:    []*Param{{Name: "name"}, {Name: "ip"}, {Name: "keypair"}, {Name: "lock"}, {Name: "role"}, {Name: "securitygroup"}, {Name: "userdata"}},
	},
	"createkeypair": {
		Action:         "create",
		Entity:         "keypair",
		Api:            "ec2",
		RequiredParams: []*Param{{Name: "name"}},
		ExtraParams:    []*Param{{Name: "encrypted"}},
	},
	"createtag": {
		Action:         "create",
		Entity:         "tag",
		Api:            "ec2",
		RequiredParams: []*Param{{Name: "key"}, {Name: "resource"}, {Name: "value"}},
		ExtraParams:    []*Param{},
	},
}

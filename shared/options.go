package shared

// Origin - attributes that uniquely identify a service instance
type Origin struct {
	Region     string
	Zone       string
	SubZone    string
	Service    string
	InstanceId string
}

// SetOrigin - required to track service identification
func SetOrigin(o Origin) {
	opt.origin = o
}

type options struct {
	origin Origin
}

var opt options

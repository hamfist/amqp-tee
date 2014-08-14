package amqptee

var (
	//Version is the git version injected via build flags (see Makefile)
	Version string
	//Rev is the git rev injected via build flags (see Makefile)
	Rev string
)

func init() {
	if Version == "" {
		Version = "<unknown>"
	}
	if Rev == "" {
		Rev = "<unknown>"
	}
}

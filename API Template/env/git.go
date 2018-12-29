package env

//These are set in the go build command (see the Makefile)

var (
	//git describe --always --dirty
	//returns tag if it exists, otherwise the short sha
	//includes -dirty if the working tree has local modification
	Version string
	//git rev-parse --abbrev-ref HEAD
	Branch string
	//date -u
	BuildTime string
	//git rev-list -1 HEAD
	SHA1 string
)

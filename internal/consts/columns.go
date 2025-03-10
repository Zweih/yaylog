package consts

const (
	Date       = "date"
	Name       = "name"
	Reason     = "reason"
	Size       = "size"
	Version    = "version"
	Depends    = "depends"
	RequiredBy = "required-by"
	Provides   = "provides"
)

var DefaultColumns = []string{Date, Name, Reason, Size}

var ValidColumns = []string{Date, Name, Reason, Size, Version, Depends, RequiredBy, Provides}

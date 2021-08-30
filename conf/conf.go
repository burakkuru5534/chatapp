package conf

import (
	"github.com/rs/zerolog"
)

var (
	DateFormat      = "02.01.2006"
	DateTimeFormat  = "02.01.2006 15:04"
	DateTimeTFormat = "02.01.2006 15:04:05"
	TimeFormat      = "15:04"

	Sqls = "./sql"

	DevMode  = false
	LogLevel = zerolog.DebugLevel

	Product        = ""
	CheckLicense   = true
	Audit          = false
	Implementation = ""

	SslEnabled = false
	SslKey     = ""
	SslCert    = ""

	UIPath     = "./ui/dist/spa"
	MainDomain = ""
)

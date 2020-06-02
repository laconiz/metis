package formatter

import "github.com/laconiz/metis/log"

type Formatter interface {
	Format(*log.Log) ([]byte, error)
}

const DefaultTimeLayout = "2006-01-02 15:04:05.000"

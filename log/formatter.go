package log

type Formatter interface {
	Format(*Log) ([]byte, error)
}

const DefaultTimeLayout = "2006-01-02 15:04:05.000"

package rogger

type Formatter interface {
	Format(*Entry) ([]byte, error)
}

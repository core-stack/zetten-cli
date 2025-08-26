package auth

type ListCommand struct {
	Verbose bool `short:"v" help:"Show verbose output"`
}

func (c *ListCommand) Run() error {
	// Implementação do comando list
	return nil
}

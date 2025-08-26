package auth

type RemoveCommand struct {
	Username string `arg:"" help:"Username to remove"`
}

func (c *RemoveCommand) Run() error {
	// Implementação do comando remove
	return nil
}

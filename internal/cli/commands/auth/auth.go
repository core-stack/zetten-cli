package auth

type AuthCommand struct {
	Add    AddCommand    `cmd:"" help:"Add authentication credentials"`
	Remove RemoveCommand `cmd:"" help:"Remove authentication credentials"`
	List   ListCommand   `cmd:"" help:"List authentication credentials"`
}

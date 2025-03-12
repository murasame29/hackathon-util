package create

type CreateRolesOptions struct {
	URL      string
	FilePath string

	DiscordBotToken string
}

func NewCreateRolesOptions() *CreateRolesOptions {
	return &CreateRolesOptions{}
}

func (o *CreateRolesOptions) Complete() error {
	panic("implement me")
}

func (o *CreateRolesOptions) Run() error {
	panic("implement me")
}

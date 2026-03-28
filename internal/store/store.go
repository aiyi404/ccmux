package store

type AppState struct {
	Mode string
}

func New(flagMode string) (*AppState, error) {
	mode := "standalone"
	if flagMode != "" {
		mode = flagMode
	}
	return &AppState{Mode: mode}, nil
}

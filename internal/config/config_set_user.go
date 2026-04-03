package config

func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName
	err := write(*cfg)
	if err != nil {
		return err
	}
	return nil
}
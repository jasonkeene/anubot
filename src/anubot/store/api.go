package store

// SetCredentials stores credentials
func (s *Store) SetCredentials(kind, user, pass string) error {

	err := s.setValueForKey(kind+"-username", user)
	if err != nil {
		return err
	}
	return s.setValueForKey(kind+"-password", pass)
}

// HasCredentials returns true if valid credentials are set
func (s *Store) HasCredentials(kind string) bool {
	user, pass, err := s.Credentials(kind)
	if err != nil || user == "" || pass == "" {
		return false
	}
	return true
}

// Credentials retrieves credentials
func (s *Store) Credentials(kind string) (string, string, error) {
	user, err := s.valueFromKey(kind + "-username")
	if err != nil {
		return "", "", err
	}
	pass, err := s.valueFromKey(kind + "-password")
	if err != nil {
		return "", "", err
	}
	return user, pass, nil
}

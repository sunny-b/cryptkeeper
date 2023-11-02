package shell

// Export represents environment variables to add and remove on the host
// shell.
type Export map[string]*string

// Add represents the addition of a new environment variable
func (e Export) Add(key, value string) {
	e[key] = &value
}

// Remove represents the removal of a given `key` environment variable.
func (e Export) Remove(key string) {
	e[key] = nil
}

package utils

func ToPtr(s string) *string {
	return &s
}

func In[M ~map[string]E, E any](key string, m M) bool {
	_, ok := m[key]
	return ok
}

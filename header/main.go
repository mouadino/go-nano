package header

type Header map[string]string

func (h Header) Get(key string) string {
	value, ok := h[key]
	if !ok {
		return ""
	}
	return value
}

func (h Header) Set(key string, value string) {
	h[key] = value
}

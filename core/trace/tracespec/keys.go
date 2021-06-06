package tracespec

var TracingKey = contextKey("X-Trace")

type contextKey string

func (c contextKey) String() string {
	return "trace/tracespec context key " + string(c)
}

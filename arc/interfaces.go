package arc

// CacheService is interface for ARC
type CacheService interface {
	Get(key interface{}) (value interface{}, ok bool)
	Put(key, value interface{}) bool
	Traverse()
	Len() int
}

// Logger is used for logging
type Logger interface {
	// Debug logging: an informative message that can aid in debugging. Will normally not be captured but may be enabled
	// on demand to diagnose production issues.
	Debug(msg string, keyvals ...interface{})
	// Info logging: a high-level event occurred such as a business transaction completed. These will never generate
	// alerts but are vital to understanding what the system is doing. Stack traces are not appropriate.
	Info(msg string, keyvals ...interface{})
	// Warn logging: the component is not operating optimally but is not in immediate danger of disrupting service.
	// Warnings are likely to generate production alerts if they happen frequently (spikes), otherwise they are
	// prioritised for action the next working day. Must be actionable. Stack trace should be restricted.
	Warn(msg string, keyvals ...interface{})
	// Error logging: The component is not operating correctly and needs urgent assistance. This is worthy of an alert
	// in production as it is causing a disruption of service. Must be actionable. Stack traces are allowed.
	Error(msg string, keyvals ...interface{})
}

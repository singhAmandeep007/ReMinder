package constants

// Environment constants
const (
	EnvProduction  = "production"
	EnvDevelopment = "development"
	EnvTest        = "test"

	DefaultJWTSecret = "default-secret-key"
)

// Rate limiting constants
const (
	DefaultRateLimit     = 100 // Requests per minute
	AuthRateLimit        = 10  // Auth requests per minute
	LowPriorityRateLimit = 50  // Requests per minute for non-critical endpoints
)

// Cache durations
const (
	ShortCacheDuration  = 60          // 1 minute in seconds
	MediumCacheDuration = 60 * 10     // 10 minutes in seconds
	LongCacheDuration   = 60 * 60 * 2 // 2 hours in seconds
)

// Database constants
const (
	SQLite    = "sqlite"
	MongoDB   = "mongodb"
	Firestore = "firestore"
)

// Collection names
const (
	UsersCollection          = "users"
	RemindersCollection      = "reminders"
	ReminderGroupsCollection = "reminder_groups"
)

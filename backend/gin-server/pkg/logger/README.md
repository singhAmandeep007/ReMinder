# Logger Package

The `logger` package provides a configurable logging solution for Go applications. It supports multiple log levels, destinations (console and file), and is designed to be extensible and thread-safe.

---

## **Constants**

### Log Levels
- `LogLevelDebug`: `"DEBUG"`
- `LogLevelInfo`: `"INFO"`
- `LogLevelWarn`: `"WARN"`
- `LogLevelError`: `"ERROR"`
- `LogLevelFatal`: `"FATAL"`
- `LogLevelUnknown`: `"UNKNOWN"`

### Logger Types
- `FileLogger`: `"file"`
- `ConsoleLogger`: `"console"`

---

## **Types**

### **LogLevel**
Represents the severity level of a log message.

#### Values:
- `DebugLevel`
- `InfoLevel`
- `WarnLevel`
- `ErrorLevel`
- `FatalLevel`

#### Methods:
- `String() string`: Returns the string representation of the log level.
- `ToZapLevel() zapcore.Level`: Converts `LogLevel` to `zapcore.Level`.

---

### **LogEntry**
Represents a single log message.

#### Fields:
- `ServiceName string`: Name of the service generating the log.
- `Level LogLevel`: Severity level of the log.
- `Message string`: Log message.
- `Fields map[string]interface{}`: Additional structured data for the log.

---

### **Destination**
Interface for log destinations.

#### Methods:
- `Write(entry LogEntry) error`: Writes a log entry to the destination.
- `Close() error`: Closes the destination.

---

### **ConsoleDestination**
Writes logs to the console.

#### Methods:
- `NewConsoleDestination() *ConsoleDestination`: Creates a new console destination.
- `Write(entry LogEntry) error`: Writes a log entry to the console.
- `Close() error`: Closes the console destination.

---

### **FileDestination**
Writes logs to a file.

#### Methods:
- `NewFileDestination(path string, maxSize int, maxBackups int, maxAge int, compress bool) *FileDestination`: Creates a new file destination.
- `Write(entry LogEntry) error`: Writes a log entry to the file.
- `Close() error`: Closes the file destination.

---

### **Logger**
Main logger interface.

#### Fields:
- `serviceName string`: Name of the service.
- `minLevel LogLevel`: Minimum log level.
- `isProd bool`: Indicates if the logger is in production mode.
- `destinations map[string]Destination`: Map of destinations.
- `defaultDests []string`: Default destinations.

#### Methods:
- `New(options ...Option) *Logger`: Creates a new logger with the given options.
- `AddDestination(name string, dest Destination)`: Adds a destination to the logger.
- `RemoveDestination(name string)`: Removes a destination from the logger.
- `SetDefaultDestinations(dests ...string)`: Sets the default destinations.
- `Debug(msg string, fields map[string]interface{}, dests ...string)`: Logs a debug message.
- `Info(msg string, fields map[string]interface{}, dests ...string)`: Logs an info message.
- `Warn(msg string, fields map[string]interface{}, dests ...string)`: Logs a warning message.
- `Error(msg string, fields map[string]interface{}, dests ...string)`: Logs an error message.
- `Fatal(msg string, fields map[string]interface{}, dests ...string)`: Logs a fatal message and terminates the program.
- `Debugf(format string, args ...interface{})`: Logs a formatted debug message.
- `Infof(format string, args ...interface{})`: Logs a formatted info message.
- `Warnf(format string, args ...interface{})`: Logs a formatted warning message.
- `Errorf(format string, args ...interface{})`: Logs a formatted error message.
- `Fatalf(format string, args ...interface{})`: Logs a formatted fatal message and terminates the program.
- `Close()`: Closes all destinations.

---

## **Options**

### **Option**
Defines a function signature for configuration options.

#### Available Options:
- `WithServiceName(name string) Option`: Sets the service name.
- `WithMinLevel(level LogLevel) Option`: Sets the minimum log level.
- `WithProduction(isProd bool) Option`: Sets the logger to production mode.
- `WithDefaultDestinations(dests ...string) Option`: Sets the default destinations.
- `WithConsoleDestination() Option`: Adds a console destination.
- `WithFileDestination(path string, maxSize, maxBackups, maxAge int, compress bool) Option`: Adds a file destination.

---

## **Usage**

### **Creating a Logger**
```go
logger := logger.New(
    logger.WithServiceName("MyService"),
    logger.WithMinLevel(logger.InfoLevel),
    logger.WithConsoleDestination(),
    logger.WithFileDestination("logs/app.log", 10, 5, 30, true),
    logger.WithDefaultDestinations("console", "file"),
)
```

### **Logging Messages**
```go
logger.Info("Application started", map[string]interface{}{
    "version": "1.0.0",
})

logger.Debug("Debugging details", nil)

logger.Error("An error occurred", map[string]interface{}{
    "error": err,
})
```

### **Closing the Logger**
```go
defer logger.Close()
```

---

## **Notes**
- In production mode, the minimum log level is automatically set to `InfoLevel` if it is lower.
- File logging uses the `lumberjack` library for log rotation and compression.
- Console logging uses `zap` for structured and human-readable logs.

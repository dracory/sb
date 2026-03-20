# Database Connection

## Database Setup

Initialize database connections with different methods and dialects.

## Connection Methods

### From Existing Go DB Instance

```go
myDb := sb.NewDatabase(sqlDb, sb.DIALECT_MYSQL)
```

### From Driver

```go
myDb, err := sb.NewDatabaseFromDriver("sqlite3", "test.db")
if err != nil {
    log.Fatal(err)
}
```

## Supported Databases

### MySQL

```go
import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

// Using existing connection
db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/dbname")
if err != nil {
    log.Fatal(err)
}

myDb := sb.NewDatabase(db, sb.DIALECT_MYSQL)
```

### PostgreSQL

```go
import (
    "database/sql"
    _ "github.com/lib/pq"
)

// Using existing connection
db, err := sql.Open("postgres", "postgres://user:password@localhost/dbname?sslmode=disable")
if err != nil {
    log.Fatal(err)
}

myDb := sb.NewDatabase(db, sb.DIALECT_POSTGRES)
```

### SQLite

```go
import (
    _ "modernc.org/sqlite"
)

// Using driver
myDb, err := sb.NewDatabaseFromDriver("sqlite3", "test.db")
if err != nil {
    log.Fatal(err)
}
```

### MSSQL

```go
import (
    "database/sql"
    _ "github.com/denisenkom/go-mssqldb"
)

// Using existing connection
db, err := sql.Open("sqlserver", "server=localhost;user id=user;password=password;database=dbname")
if err != nil {
    log.Fatal(err)
}

myDb := sb.NewDatabase(db, sb.DIALECT_MSSQL)
```

## Dialect Constants

```go
const (
    DIALECT_MYSQL     = "mysql"
    DIALECT_POSTGRES  = "postgres"
    DIALECT_SQLITE    = "sqlite"
    DIALECT_MSSQL     = "mssql"
)
```

## Connection Examples

### MySQL with Connection Pooling

```go
func SetupMySQL(dsn string) (*sb.Database, error) {
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, err
    }
    
    // Configure connection pool
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)
    
    // Test connection
    if err := db.Ping(); err != nil {
        return nil, err
    }
    
    return sb.NewDatabase(db, sb.DIALECT_MYSQL), nil
}

// Usage
myDb, err := SetupMySQL("user:password@tcp(localhost:3306)/dbname?parseTime=true")
```

### PostgreSQL with SSL

```go
func SetupPostgreSQL(connStr string) (*sb.Database, error) {
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }
    
    // Configure connection pool
    db.SetMaxOpenConns(20)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(10 * time.Minute)
    
    // Test connection
    if err := db.Ping(); err != nil {
        return nil, err
    }
    
    return sb.NewDatabase(db, sb.DIALECT_POSTGRES), nil
}

// Usage
connStr := "postgres://user:password@localhost/dbname?sslmode=require&sslrootcert=/path/to/cert"
myDb, err := SetupPostgreSQL(connStr)
```

### SQLite with WAL Mode

```go
func SetupSQLite(dbPath string) (*sb.Database, error) {
    myDb, err := sb.NewDatabaseFromDriver("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }
    
    // Configure SQLite settings
    ctx := context.Background()
    
    // Enable WAL mode for better concurrency
    if _, err := myDb.Exec(ctx, "PRAGMA journal_mode=WAL"); err != nil {
        return nil, err
    }
    
    // Enable foreign key constraints
    if _, err := myDb.Exec(ctx, "PRAGMA foreign_keys=ON"); err != nil {
        return nil, err
    }
    
    return myDb, nil
}

// Usage
myDb, err := SetupSQLite("./data/app.db")
```

### MSSQL with Connection Timeout

```go
func SetupMSSQL(connStr string) (*sb.Database, error) {
    db, err := sql.Open("sqlserver", connStr)
    if err != nil {
        return nil, err
    }
    
    // Configure connection pool
    db.SetMaxOpenConns(30)
    db.SetMaxIdleConns(10)
    db.SetConnMaxLifetime(15 * time.Minute)
    
    // Test connection
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    if err := db.PingContext(ctx); err != nil {
        return nil, err
    }
    
    return sb.NewDatabase(db, sb.DIALECT_MSSQL), nil
}

// Usage
connStr := "server=localhost;user id=user;password=password;database=dbname;connection timeout=30"
myDb, err := SetupMSSQL(connStr)
```

## Database Operations

### Execute SQL

```go
myDb := sb.NewDatabase(sqlDb, sb.DIALECT_MYSQL)
ctx := context.Background()

sql := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Insert(map[string]string{
        "name": "John Doe",
        "email": "john@example.com",
    })

_, err := myDb.Exec(ctx, sql)
if err != nil {
    log.Fatal(err)
}
```

### Query as Map

```go
// Returns map[string]any
ctx := context.Background()
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Where(&sb.Where{Column: "id", Operator: "=", Value: "1"}).
    Select([]string{"name", "email"})

mapAny, err := myDb.SelectToMapAny(ctx, sql)
if err != nil {
    log.Fatal(err)
}

// Returns map[string]string
mapString, err := myDb.SelectToMapString(ctx, sql)
if err != nil {
    log.Fatal(err)
}
```

## Connection Management

### Connection Pool Configuration

```go
func ConfigureConnectionPool(db *sql.DB) {
    // Set maximum number of open connections
    db.SetMaxOpenConns(25)
    
    // Set maximum number of idle connections
    db.SetMaxIdleConns(5)
    
    // Set maximum lifetime of connections
    db.SetConnMaxLifetime(5 * time.Minute)
    
    // Set maximum idle time for connections
    db.SetConnMaxIdleTime(2 * time.Minute)
}
```

### Health Check

```go
func CheckDatabaseHealth(myDb *sb.Database) error {
    ctx := context.Background()
    
    // Simple ping test
    return myDb.(*sb.Database).Ping()
}

func CheckDatabaseHealthWithTimeout(myDb *sb.Database) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    return myDb.(*sb.Database).PingContext(ctx)
}
```

### Graceful Shutdown

```go
func CloseDatabase(myDb *sb.Database) error {
    if myDb != nil {
        return myDb.(*sb.Database).Close()
    }
    return nil
}

// Usage in main()
func main() {
    myDb, err := SetupMySQL("user:password@tcp(localhost:3306)/dbname")
    if err != nil {
        log.Fatal(err)
    }
    defer CloseDatabase(myDb)
    
    // Application logic...
}
```

## Environment Configuration

### Configuration Structure

```go
type DatabaseConfig struct {
    Driver   string
    Host     string
    Port     int
    Database string
    Username string
    Password string
    SSLMode  string
    MaxOpen  int
    MaxIdle  int
}

func ConnectFromConfig(config DatabaseConfig) (*sb.Database, error) {
    switch config.Driver {
    case "mysql":
        dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
            config.Username, config.Password, config.Host, config.Port, config.Database)
        return SetupMySQL(dsn)
        
    case "postgres":
        connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
            config.Username, config.Password, config.Host, config.Port, config.Database, config.SSLMode)
        return SetupPostgreSQL(connStr)
        
    case "sqlite":
        return SetupSQLite(config.Database)
        
    case "mssql":
        connStr := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s",
            config.Host, config.Username, config.Password, config.Database)
        return SetupMSSQL(connStr)
        
    default:
        return nil, fmt.Errorf("unsupported driver: %s", config.Driver)
    }
}
```

### Environment Variables

```go
import "os"

func LoadDatabaseFromEnv() (*sb.Database, error) {
    config := DatabaseConfig{
        Driver:   getEnv("DB_DRIVER", "mysql"),
        Host:     getEnv("DB_HOST", "localhost"),
        Port:     getEnvInt("DB_PORT", 3306),
        Database: getEnv("DB_NAME", "testdb"),
        Username: getEnv("DB_USER", "root"),
        Password: getEnv("DB_PASSWORD", ""),
        SSLMode:  getEnv("DB_SSLMODE", "disable"),
        MaxOpen:  getEnvInt("DB_MAX_OPEN", 25),
        MaxIdle:  getEnvInt("DB_MAX_IDLE", 5),
    }
    
    return ConnectFromConfig(config)
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}
```

## Best Practices

### Connection Management

1. **Use connection pooling** - Configure appropriate pool sizes
2. **Set timeouts** - Prevent hanging connections
3. **Handle errors gracefully** - Always check for connection errors
4. **Close connections** - Use defer for proper cleanup
5. **Monitor connections** - Track pool usage and health

### Security

1. **Use environment variables** - Never hardcode credentials
2. **Enable SSL/TLS** - Encrypt connections in production
3. **Limit permissions** - Use least privilege principle
4. **Use connection strings** - Avoid building DSNs manually
5. **Rotate credentials** - Regularly update passwords

### Performance

1. **Tune pool sizes** - Based on application load
2. **Set connection lifetimes** - Prevent stale connections
3. **Use prepared statements** - For repeated queries
4. **Monitor metrics** - Track connection usage
5. **Test under load** - Verify performance under stress

## Troubleshooting

### Common Issues

1. **Connection refused** - Check database server status
2. **Authentication failed** - Verify credentials
3. **Timeout errors** - Increase timeout values
4. **Pool exhaustion** - Increase max connections
5. **SSL errors** - Check certificate configuration

### Debug Connections

```go
func DebugConnection(myDb *sb.Database) {
    db := myDb.(*sb.Database)
    
    stats := db.Stats()
    fmt.Printf("Open Connections: %d\n", stats.OpenConnections)
    fmt.Printf("In Use: %d\n", stats.InUse)
    fmt.Printf("Idle: %d\n", stats.Idle)
    fmt.Printf("Wait Count: %d\n", stats.WaitCount)
    fmt.Printf("Wait Duration: %v\n", stats.WaitDuration)
    fmt.Printf("Max Idle Closed: %d\n", stats.MaxIdleClosed)
    fmt.Printf("Max Lifetime Closed: %d\n", stats.MaxLifetimeClosed)
}
```

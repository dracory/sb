# Database Transactions

## Transaction Management

SB provides built-in transaction support for ensuring data consistency across multiple operations.

## ExecInTransaction Method

Execute multiple operations in a single transaction with automatic rollback on error.

### Basic Transaction

```go
import _ "modernc.org/sqlite"

myDb := sb.NewDatabaseFromDriver("sqlite3", "test.db")
ctx := context.Background()

err := myDb.ExecInTransaction(ctx, func(tx *sb.Database) error {
	if _, err := tx.Exec(ctx, sql1); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, sql2); err != nil {
		return err
	}

	return nil
})

if err != nil {
	// handle error - transaction was automatically rolled back
}
```

### Complex Transaction Example

```go
func CreateUserWithProfile(db *sb.Database, name, email, bio string) error {
    ctx := context.Background()
    
    return db.ExecInTransaction(ctx, func(tx *sb.Database) error {
        // Insert user
        userSQL := sb.NewBuilder(sb.DIALECT_MYSQL).
            Table("users").
            Insert(map[string]string{
                "name":       name,
                "email":      email,
                "created_at": time.Now().Format("2006-01-02T15:04:05"),
            })
        
        if _, err := tx.Exec(ctx, userSQL); err != nil {
            return fmt.Errorf("failed to create user: %w", err)
        }
        
        // Get user ID (simplified - in practice you'd use LAST_INSERT_ID())
        userID := "1" // Assume we got this from the previous insert
        
        // Insert profile
        profileSQL := sb.NewBuilder(sb.DIALECT_MYSQL).
            Table("profiles").
            Insert(map[string]string{
                "user_id":    userID,
                "bio":        bio,
                "created_at": time.Now().Format("2006-01-02T15:04:05"),
            })
        
        if _, err := tx.Exec(ctx, profileSQL); err != nil {
            return fmt.Errorf("failed to create profile: %w", err)
        }
        
        return nil
    })
}
```

## Manual Transaction Control

For more complex scenarios, use manual transaction control methods.

### Manual Transaction Example

```go
import _ "modernc.org/sqlite"

myDb := sb.NewDatabaseFromDriver("sqlite3", "test.db")
myDb.BeginTransaction()

ctx := context.Background()

_, err := myDb.Exec(ctx, sql1)

if err != nil {
	myDb.RollbackTransaction()
	return err
}

_, err = myDb.Exec(ctx, sql2)

if err != nil {
	myDb.RollbackTransaction()
	return err
}

myDb.CommitTransaction()
```

### Manual Transaction with Error Handling

```go
func ProcessOrder(db *sb.Database, orderID int, status string) error {
    ctx := context.Background()
    
    // Begin transaction
    db.BeginTransaction()
    defer func() {
        if r := recover(); r != nil {
            db.RollbackTransaction()
        }
    }()
    
    // Update order status
    orderSQL := sb.NewBuilder(sb.DIALECT_MYSQL).
        Table("orders").
        Where(&sb.Where{Column: "id", Operator: "=", Value: strconv.Itoa(orderID)}).
        Update(map[string]string{
            "status":     status,
            "updated_at": time.Now().Format("2006-01-02T15:04:05"),
        })
    
    if _, err := db.Exec(ctx, orderSQL); err != nil {
        db.RollbackTransaction()
        return fmt.Errorf("failed to update order: %w", err)
    }
    
    // Update inventory if order is completed
    if status == "completed" {
        // Get order items first
        itemsSQL := sb.NewBuilder(sb.DIALECT_MYSQL).
            Table("order_items").
            Where(&sb.Where{Column: "order_id", Operator: "=", Value: strconv.Itoa(orderID)}).
            Select([]string{"product_id", "quantity"})
        
        // In practice, you'd execute this query and process results
        // For simplicity, assuming we have the data
        
        updateInventorySQL := sb.NewBuilder(sb.DIALECT_MYSQL).
            Table("inventory").
            Where(&sb.Where{Column: "product_id", Operator: "=", Value: "123"}).
            Update(map[string]string{
                "quantity":   "quantity - 1",  // Simplified
                "updated_at": time.Now().Format("2006-01-02T15:04:05"),
            })
        
        if _, err := db.Exec(ctx, updateInventorySQL); err != nil {
            db.RollbackTransaction()
            return fmt.Errorf("failed to update inventory: %w", err)
        }
    }
    
    // Commit transaction
    if err := db.CommitTransaction(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }
    
    return nil
}
```

## Transaction Methods Reference

### Automatic Transaction Management

- `ExecInTransaction(ctx, fn)` - Execute function within transaction

### Manual Transaction Control

- `BeginTransaction()` - Start a new transaction
- `CommitTransaction()` - Commit the current transaction
- `RollbackTransaction()` - Rollback the current transaction

## Transaction Patterns

### Banking Transfer Example

```go
func TransferFunds(db *sb.Database, fromAccount, toAccount int, amount float64) error {
    ctx := context.Background()
    
    return db.ExecInTransaction(ctx, func(tx *sb.Database) error {
        // Debit from source account
        debitSQL := sb.NewBuilder(sb.DIALECT_MYSQL).
            Table("accounts").
            Where(&sb.Where{Column: "id", Operator: "=", Value: strconv.Itoa(fromAccount)}).
            Where(&sb.Where{Column: "balance", Operator: ">=", Value: strconv.FormatFloat(amount, 'f', 2, 64)}).
            Update(map[string]string{
                "balance": fmt.Sprintf("balance - %.2f", amount),
            })
        
        if result, err := tx.Exec(ctx, debitSQL); err != nil {
            return fmt.Errorf("failed to debit source account: %w", err)
        } else if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
            return fmt.Errorf("insufficient funds or account not found")
        }
        
        // Credit to destination account
        creditSQL := sb.NewBuilder(sb.DIALECT_MYSQL).
            Table("accounts").
            Where(&sb.Where{Column: "id", Operator: "=", Value: strconv.Itoa(toAccount)}).
            Update(map[string]string{
                "balance": fmt.Sprintf("balance + %.2f", amount),
            })
        
        if result, err := tx.Exec(ctx, creditSQL); err != nil {
            return fmt.Errorf("failed to credit destination account: %w", err)
        } else if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
            return fmt.Errorf("destination account not found")
        }
        
        // Record transaction
        transactionSQL := sb.NewBuilder(sb.DIALECT_MYSQL).
            Table("transactions").
            Insert(map[string]string{
                "from_account": strconv.Itoa(fromAccount),
                "to_account":   strconv.Itoa(toAccount),
                "amount":       strconv.FormatFloat(amount, 'f', 2, 64),
                "created_at":   time.Now().Format("2006-01-02T15:04:05"),
            })
        
        if _, err := tx.Exec(ctx, transactionSQL); err != nil {
            return fmt.Errorf("failed to record transaction: %w", err)
        }
        
        return nil
    })
}
```

### Order Processing with Inventory

```go
func ProcessOrderWithInventory(db *sb.Database, orderID int) error {
    ctx := context.Background()
    
    return db.ExecInTransaction(ctx, func(tx *sb.Database) error {
        // Lock order for update
        lockSQL := sb.NewBuilder(sb.DIALECT_MYSQL).
            Table("orders").
            Where(&sb.Where{Column: "id", Operator: "=", Value: strconv.Itoa(orderID)}).
            Where(&sb.Where{Column: "status", Operator: "=", Value: "pending"}).
            Select([]string{"id"})
        
        // In practice, execute this to verify order exists and is pending
        
        // Update order status
        updateOrderSQL := sb.NewBuilder(sb.DIALECT_MYSQL).
            Table("orders").
            Where(&sb.Where{Column: "id", Operator: "=", Value: strconv.Itoa(orderID)}).
            Update(map[string]string{
                "status":     "processing",
                "updated_at": time.Now().Format("2006-01-02T15:04:05"),
            })
        
        if _, err := tx.Exec(ctx, updateOrderSQL); err != nil {
            return fmt.Errorf("failed to update order status: %w", err)
        }
        
        // Get order items and update inventory
        itemsSQL := sb.NewBuilder(sb.DIALECT_MYSQL).
            Table("order_items").
            Where(&sb.Where{Column: "order_id", Operator: "=", Value: strconv.Itoa(orderID)}).
            Select([]string{"product_id", "quantity"})
        
        // Execute query and process each item
        // For each item:
        updateInventorySQL := sb.NewBuilder(sb.DIALECT_MYSQL).
            Table("inventory").
            Where(&sb.Where{Column: "product_id", Operator: "=", Value: "productID"}).
            Where(&sb.Where{Column: "quantity", Operator: ">=", Value: "itemQuantity"}).
            Update(map[string]string{
                "quantity":   fmt.Sprintf("quantity - %s", "itemQuantity"),
                "updated_at": time.Now().Format("2006-01-02T15:04:05"),
            })
        
        // Execute inventory updates for each item
        // If any inventory update fails, return error to rollback
        
        // Update order to completed
        completeOrderSQL := sb.NewBuilder(sb.DIALECT_MYSQL).
            Table("orders").
            Where(&sb.Where{Column: "id", Operator: "=", Value: strconv.Itoa(orderID)}).
            Update(map[string]string{
                "status":     "completed",
                "updated_at": time.Now().Format("2006-01-02T15:04:05"),
            })
        
        if _, err := tx.Exec(ctx, completeOrderSQL); err != nil {
            return fmt.Errorf("failed to complete order: %w", err)
        }
        
        return nil
    })
}
```

## Best Practices

### Transaction Design

1. **Keep transactions short** - Minimize lock duration
2. **Avoid user input in transactions** - Validate before starting
3. **Use appropriate isolation levels** - Balance consistency and performance
4. **Handle errors properly** - Always rollback on error
5. **Log transaction boundaries** - For debugging and auditing

### Error Handling

```go
func SafeTransactionExample(db *sb.Database) error {
    ctx := context.Background()
    
    return db.ExecInTransaction(ctx, func(tx *sb.Database) error {
        // Validate input before transaction
        // ... validation logic ...
        
        // Perform operations
        if err := performOperation1(tx, ctx); err != nil {
            return fmt.Errorf("operation 1 failed: %w", err)
        }
        
        if err := performOperation2(tx, ctx); err != nil {
            return fmt.Errorf("operation 2 failed: %w", err)
        }
        
        // Log success
        log.Printf("Transaction completed successfully")
        
        return nil
    })
}
```

### Nested Transactions (Savepoints)

```go
func ComplexTransactionWithSavepoints(db *sb.Database) error {
    ctx := context.Background()
    
    return db.ExecInTransaction(ctx, func(tx *sb.Database) error {
        // Main operation
        if err := mainOperation(tx, ctx); err != nil {
            return err
        }
        
        // Attempt optional operation
        if err := optionalOperation(tx, ctx); err != nil {
            // Log but don't fail the whole transaction
            log.Printf("Optional operation failed: %v", err)
        }
        
        return nil
    })
}
```

## Performance Considerations

1. **Batch operations** - Group multiple small operations
2. **Avoid long-running transactions** - Can cause lock contention
3. **Use appropriate isolation levels** - Read committed is usually sufficient
4. **Monitor transaction duration** - Set timeouts for long operations
5. **Consider connection pooling** - Reuse connections efficiently

## Database-Specific Behavior

### MySQL
- Supports ACID transactions with InnoDB engine
- Autocommit mode by default
- Supports savepoints for nested transactions

### PostgreSQL
- Full ACID compliance
- MVCC (Multi-Version Concurrency Control)
- Supports savepoints and prepared transactions

### SQLite
- Supports ACID transactions
- Serializes transactions by default
- May have performance limitations with concurrent writes

### MSSQL
- Full ACID compliance
- Supports multiple isolation levels
- Supports savepoints and distributed transactions

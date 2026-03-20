package sb_test

import (
	"testing"

	"github.com/dracory/sb"
)

// BenchmarkSQLGeneration benchmarks basic SQL generation performance
func BenchmarkSQLGeneration(b *testing.B) {
	dialects := []string{
		sb.DIALECT_MYSQL,
		sb.DIALECT_POSTGRES,
		sb.DIALECT_SQLITE,
		sb.DIALECT_MSSQL,
	}

	for _, dialect := range dialects {
		b.Run(dialect, func(b *testing.B) {
			builder := sb.NewBuilder(dialect).Table("users")

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = builder.Select([]string{"id", "name", "email"})
			}
		})
	}
}

// BenchmarkComplexQuery benchmarks complex query generation performance
func BenchmarkComplexQuery(b *testing.B) {
	dialects := []string{
		sb.DIALECT_MYSQL,
		sb.DIALECT_POSTGRES,
		sb.DIALECT_SQLITE,
		sb.DIALECT_MSSQL,
	}

	for _, dialect := range dialects {
		b.Run(dialect, func(b *testing.B) {
			builder := sb.NewBuilder(dialect).
				Table("users").
				Where(&sb.Where{Column: "status", Operator: "=", Value: "active"}).
				Where(&sb.Where{Column: "created_at", Operator: ">", Value: "2023-01-01"}).
				OrderBy("name", "ASC").
				Limit(10)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = builder.Select([]string{"id", "name", "email", "created_at"})
			}
		})
	}
}

// BenchmarkJoinQuery benchmarks JOIN query generation performance
func BenchmarkJoinQuery(b *testing.B) {
	dialects := []string{
		sb.DIALECT_MYSQL,
		sb.DIALECT_POSTGRES,
		sb.DIALECT_SQLITE,
		sb.DIALECT_MSSQL,
	}

	for _, dialect := range dialects {
		b.Run(dialect, func(b *testing.B) {
			builder := sb.NewBuilder(dialect).
				Table("users").
				InnerJoin("profiles", "users.id = profiles.user_id").
				LeftJoin("orders", "users.id = orders.user_id").
				Where(&sb.Where{Column: "users.status", Operator: "=", Value: "active"})

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = builder.Select([]string{"users.name", "profiles.avatar", "orders.total"})
			}
		})
	}
}

// BenchmarkSubquery benchmarks subquery generation performance
func BenchmarkSubquery(b *testing.B) {
	dialects := []string{
		sb.DIALECT_MYSQL,
		sb.DIALECT_POSTGRES,
		sb.DIALECT_SQLITE,
		sb.DIALECT_MSSQL,
	}

	for _, dialect := range dialects {
		b.Run(dialect, func(b *testing.B) {
			subqueryBuilder := sb.NewBuilder(dialect).
				Table("orders").
				Where(&sb.Where{Column: "total", Operator: ">", Value: "1000"})

			builder, err := sb.NewBuilder(dialect).
				Table("users").
				InSubquery(subqueryBuilder)
			if err != nil {
				b.Fatal(err)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = builder.Select([]string{"name", "email"})
			}
		})
	}
}

// BenchmarkErrorHandling benchmarks error handling performance
func BenchmarkErrorHandling(b *testing.B) {
	dialects := []string{
		sb.DIALECT_MYSQL,
		sb.DIALECT_POSTGRES,
		sb.DIALECT_SQLITE,
		sb.DIALECT_MSSQL,
	}

	for _, dialect := range dialects {
		b.Run(dialect, func(b *testing.B) {
			builder := sb.NewBuilder(dialect).
				Table("users").
				Where(&sb.Where{Column: "status", Operator: "=", Value: "active"}).
				OrderBy("name", "ASC").
				Limit(10)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = builder.Select([]string{"id", "name", "email"})
			}
		})
	}
}

// BenchmarkErrorCollection benchmarks error collection overhead
func BenchmarkErrorCollection(b *testing.B) {
	dialects := []string{
		sb.DIALECT_MYSQL,
		sb.DIALECT_POSTGRES,
		sb.DIALECT_SQLITE,
		sb.DIALECT_MSSQL,
	}

	for _, dialect := range dialects {
		b.Run(dialect, func(b *testing.B) {
			// Test with error collection (normal case)
			b.Run("Normal", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					builder := sb.NewBuilder(dialect).
						Table("users").
						Where(&sb.Where{Column: "status", Operator: "=", Value: "active"})
					_, _, _ = builder.Select([]string{"id", "name"})
				}
			})

			// Test with error collection (error case)
			b.Run("Error", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					builder := sb.NewBuilder("invalid_dialect").
						Table("users").
						Where(&sb.Where{Column: "status", Operator: "=", Value: "active"})
					_, _, _ = builder.Select([]string{"id", "name"})
				}
			})
		})
	}
}

// BenchmarkFluentChaining benchmarks fluent chaining performance
func BenchmarkFluentChaining(b *testing.B) {
	dialects := []string{
		sb.DIALECT_MYSQL,
		sb.DIALECT_POSTGRES,
		sb.DIALECT_SQLITE,
		sb.DIALECT_MSSQL,
	}

	for _, dialect := range dialects {
		b.Run(dialect, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = sb.NewBuilder(dialect).
					Table("users").
					Where(&sb.Where{Column: "status", Operator: "=", Value: "active"}).
					Where(&sb.Where{Column: "created_at", Operator: ">", Value: "2023-01-01"}).
					OrderBy("name", "ASC").
					Limit(10).
					Select([]string{"id", "name", "email"})
			}
		})
	}
}

// BenchmarkCreateTable benchmarks CREATE TABLE generation performance
func BenchmarkCreateTable(b *testing.B) {
	dialects := []string{
		sb.DIALECT_MYSQL,
		sb.DIALECT_POSTGRES,
		sb.DIALECT_SQLITE,
		sb.DIALECT_MSSQL,
	}

	for _, dialect := range dialects {
		b.Run(dialect, func(b *testing.B) {
			columns := []sb.Column{
				{Name: "id", Type: "INTEGER", PrimaryKey: true, AutoIncrement: true},
				{Name: "name", Type: "VARCHAR(255)", Nullable: false},
				{Name: "email", Type: "VARCHAR(255)", Unique: true},
				{Name: "created_at", Type: "TIMESTAMP", Default: "CURRENT_TIMESTAMP"},
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				builder := sb.NewBuilder(dialect).Table("users")
				for _, col := range columns {
					builder.Column(col)
				}
				_, _ = builder.Create()
			}
		})
	}
}

// BenchmarkInsert benchmarks INSERT statement generation performance
func BenchmarkInsert(b *testing.B) {
	dialects := []string{
		sb.DIALECT_MYSQL,
		sb.DIALECT_POSTGRES,
		sb.DIALECT_SQLITE,
		sb.DIALECT_MSSQL,
	}

	for _, dialect := range dialects {
		b.Run(dialect, func(b *testing.B) {
			data := map[string]string{
				"name":   "John Doe",
				"email":  "john@example.com",
				"status": "active",
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = sb.NewBuilder(dialect).Table("users").Insert(data)
			}
		})
	}
}

// BenchmarkUpdate benchmarks UPDATE statement generation performance
func BenchmarkUpdate(b *testing.B) {
	dialects := []string{
		sb.DIALECT_MYSQL,
		sb.DIALECT_POSTGRES,
		sb.DIALECT_SQLITE,
		sb.DIALECT_MSSQL,
	}

	for _, dialect := range dialects {
		b.Run(dialect, func(b *testing.B) {
			data := map[string]string{
				"status":     "inactive",
				"updated_at": "2023-12-01",
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = sb.NewBuilder(dialect).
					Table("users").
					Where(&sb.Where{Column: "id", Operator: "=", Value: "1"}).
					Update(data)
			}
		})
	}
}

// BenchmarkDelete benchmarks DELETE statement generation performance
func BenchmarkDelete(b *testing.B) {
	dialects := []string{
		sb.DIALECT_MYSQL,
		sb.DIALECT_POSTGRES,
		sb.DIALECT_SQLITE,
		sb.DIALECT_MSSQL,
	}

	for _, dialect := range dialects {
		b.Run(dialect, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = sb.NewBuilder(dialect).
					Table("users").
					Where(&sb.Where{Column: "status", Operator: "=", Value: "inactive"}).
					Delete()
			}
		})
	}
}

// BenchmarkIndexOperations benchmarks index operation performance
func BenchmarkIndexOperations(b *testing.B) {
	dialects := []string{
		sb.DIALECT_MYSQL,
		sb.DIALECT_POSTGRES,
		sb.DIALECT_SQLITE,
		sb.DIALECT_MSSQL,
	}

	for _, dialect := range dialects {
		b.Run(dialect, func(b *testing.B) {
			b.Run("CreateIndex", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, _ = sb.NewBuilder(dialect).Table("users").CreateIndex("idx_email", "email")
				}
			})

			b.Run("DropIndex", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, _ = sb.NewBuilder(dialect).Table("users").DropIndex("idx_email")
				}
			})
		})
	}
}

// BenchmarkMemoryUsage measures memory allocation patterns
func BenchmarkMemoryUsage(b *testing.B) {
	dialects := []string{
		sb.DIALECT_MYSQL,
		sb.DIALECT_POSTGRES,
		sb.DIALECT_SQLITE,
		sb.DIALECT_MSSQL,
	}

	for _, dialect := range dialects {
		b.Run(dialect, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				builder := sb.NewBuilder(dialect).
					Table("users").
					Where(&sb.Where{Column: "status", Operator: "=", Value: "active"}).
					OrderBy("name", "ASC").
					Limit(10)
				_, _, _ = builder.Select([]string{"id", "name", "email"})
			}
		})
	}
}

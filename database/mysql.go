package database

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"

	"go-fiber-gorm-sample/config"

	gomysql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMySQL() (*gorm.DB, error) {
	// Get the certificate file path from the environment variable
	certFilePath := os.Getenv("MYSQL_CERT_FILE_PATH")

	var useTLS bool
	if certFilePath == "" {
		useTLS = false
	} else {
		useTLS = true
	}

	fmt.Println("useTLS", useTLS)

	var dsn string

	if useTLS {
		// Load CA certificate for TLS configuration
		rootCertPool := x509.NewCertPool()
		pem, err := os.ReadFile(certFilePath) // Path to your CA cert file
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate file: %v", err)
		}

		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			log.Println("Failed to append PEM.")
		}

		// Create a TLS configuration
		tlsConfig := &tls.Config{
			RootCAs: rootCertPool,
		}

		// Register the custom TLS configuration
		err = gomysql.RegisterTLSConfig("custom", tlsConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to register TLS config: %v", err)
		}
	}

	// Construct the DSN using Config.DBConfig values
	dsn = fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true&timeout=15s",
		config.DBConfig.User,
		config.DBConfig.Password,
		config.DBConfig.Host,
		config.DBConfig.Port,
		config.DBConfig.DBName,
	)
	log.Printf("[DEBUG] Using non-TLS DSN: %s\n", sanitizeDSN(dsn))
	// Connect to TiDB (MySQL compatible)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Logger: logger.New(
		// 	log.New(os.Stdout, "\r\n", log.LstdFlags), // use standard log
		// 	logger.Config{
		// 		SlowThreshold:             time.Second,
		// 		LogLevel:                  logger.Info, // or logger.Warn for less verbosity
		// 		IgnoreRecordNotFoundError: true,
		// 		Colorful:                  true,
		// 	},
		// ),
	})
	if err != nil {
		log.Printf("[DEBUG] Failed to connect database: %v\n", err)
		panic("failed to connect database")
	}

	// Automigrate all models
	//if err := db.AutoMigrate(
	//	&entities.Sample{},
	//); err != nil {
	//	return nil, fmt.Errorf("failed to automigrate models: %v", err)
	//}

	// Ping the database to ensure connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}

func sanitizeDSN(dsn string) string {
	redacted := dsn
	if atPos := findRuneFromLeft(redacted, '@'); atPos != -1 {
		if colonPos := findRuneReverse(redacted[:atPos], ':'); colonPos != -1 {
			redacted = redacted[:colonPos+1] + "*****" + redacted[atPos:]
		}
	}
	return redacted
}

func findRuneFromLeft(s string, r rune) int {
	for i, c := range s {
		if c == r {
			return i
		}
	}
	return -1
}

func findRuneReverse(s string, r rune) int {
	for i := len(s) - 1; i >= 0; i-- {
		if rune(s[i]) == r {
			return i
		}
	}
	return -1
}

func AutomigrateModels(db *gorm.DB, models ...interface{}) error {
	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to automigrate models: %v", err)
	}
	return nil
}

func CloseDB(database *gorm.DB) {
	sqlDB, err := database.DB()
	if err != nil {
		log.Fatalf("Failed to retrieve sql.DB: %v", err)
	}
	_ = sqlDB.Close()
}

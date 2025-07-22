package config

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() (*sql.DB, error) {
	// Replace with your database credentials
	dsn := "root:@tcp(127.0.0.1:3306)/finance?parseTime=true&multiStatements=true"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	// Enable connection pooling
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * 60) // 5 minutes

	// Initialize database schema
	////if err := initSchema(db); err != nil {
	//	return nil, fmt.Errorf("error initializing database schema: %v", err)
	//}

	fmt.Println("Successfully connected to and initialized the database")
	return db, nil
}

//func initSchema(db *sql.DB) error {
//	schemas := []string{
//		`CREATE TABLE IF NOT EXISTS users (
//			id INT AUTO_INCREMENT PRIMARY KEY,
//			username VARCHAR(50) NOT NULL UNIQUE,
//			pin VARCHAR(255) NOT NULL,
//			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
//		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
//
//		`CREATE TABLE IF NOT EXISTS transaction_types (
//			id INT AUTO_INCREMENT PRIMARY KEY,
//			name VARCHAR(50) NOT NULL UNIQUE,
//			description TEXT
//		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
//
//		`CREATE TABLE IF NOT EXISTS categories (
//			id INT AUTO_INCREMENT PRIMARY KEY,
//			name VARCHAR(100) NOT NULL,
//			type_id INT NOT NULL,
//			description TEXT,
//			FOREIGN KEY (type_id) REFERENCES transaction_types(id) ON DELETE CASCADE
//		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
//
//		`CREATE TABLE IF NOT EXISTS accounts (
//			id INT(11) NOT NULL AUTO_INCREMENT,
//			name VARCHAR(100) NOT NULL COLLATE 'utf8mb4_uca1400_ai_ci',
//			account_type_id INT(11) NOT NULL,
//			initial_amount DECIMAL(14,2) NOT NULL DEFAULT 0.00,
//			initial_date DATE NOT NULL,
//			remarks TEXT DEFAULT NULL COLLATE 'utf8mb4_uca1400_ai_ci',
//			credit_limit DECIMAL(14,2) DEFAULT NULL,
//			credit_bill_date DATE DEFAULT NULL,
//			credit_payment_date DATE DEFAULT NULL,
//			current_balance DECIMAL(15,2) DEFAULT 0.00,
//			created_at TIMESTAMP NULL DEFAULT current_timestamp(),
//			updated_at TIMESTAMP NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
//			created_by INT(11) NOT NULL,
//			updated_by INT(11) DEFAULT NULL,
//			PRIMARY KEY (id) USING BTREE,
//			KEY account_type_id (account_type_id),
//			CONSTRAINT accounts_ibfk_1 FOREIGN KEY (account_type_id) REFERENCES account_types (id)
//		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;`,
//
//		`CREATE TABLE IF NOT EXISTS transactions (
//			id INT AUTO_INCREMENT PRIMARY KEY,
//			amount DECIMAL(10,2) NOT NULL,
//			description TEXT,
//			type_id INT NOT NULL,
//			category_id INT,
//			account_id INT NOT NULL,
//			user_id INT NOT NULL,
//			transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
//			FOREIGN KEY (type_id) REFERENCES transaction_types(id),
//			FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL,
//			FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,
//			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
//		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`,
//
//		// Insert default transaction types if they don't exist
//		`INSERT IGNORE INTO transaction_types (id, name, description) VALUES
//		(1, 'income', 'Money received'),
//		(2, 'expense', 'Money spent');`,
//
//		// Insert default categories if they don't exist
//		`INSERT IGNORE INTO categories (id, name, type_id, description) VALUES
//		(1, 'Salary', 1, 'Monthly salary'),
//		(2, 'Freelance', 1, 'Freelance income'),
//		(3, 'Investments', 1, 'Investment returns'),
//		(4, 'Other Income', 1, 'Other income sources'),
//		(5, 'Food & Dining', 2, 'Groceries, restaurants, etc.'),
//		(6, 'Shopping', 2, 'Clothing, electronics, etc.'),
//		(7, 'Housing', 2, 'Rent, mortgage, utilities'),
//		(8, 'Transportation', 2, 'Fuel, public transport, etc.'),
//		(9, 'Entertainment', 2, 'Movies, games, hobbies'),
//		(10, 'Other Expenses', 2, 'Other expenses');`,
//
//		// Create a default user if none exists
//		`INSERT IGNORE INTO users (id, user_name, pin) VALUES
//		(1, 'admin', '1234');`,
//
//		// Create default account types if they don't exist
//		`INSERT IGNORE INTO account_types (id, name, description) VALUES
//		(1, 'cash', 'Physical currency in hand'),
//		(2, 'savings', 'Savings bank account'),
//		(3, 'current', 'Current bank account'),
//		(4, 'wallet', 'Mobile or digital wallets'),
//		(5, 'credit_card', 'Credit card account'),
//		(6, 'credit_line', 'Buy Now Pay Later, overdraft, or other short-term borrowing'),
//		(7, 'loan', 'Loan accounts – personal, auto, home, etc.'),
//		(8, 'investment_equity', 'Stocks, mutual funds, ETFs'),
//		(9, 'investment_fixed', 'FDs, RDs, bonds, debentures'),
//		(10, 'investment_retirement', 'NPS, PPF, EPF, etc.'),
//		(11, 'investment_insurance', 'LIC, ULIP, endowment plans'),
//		(12, 'other', 'Catch-all for other account types');`,
//
//		// Create a default account if none exists
//		`INSERT IGNORE INTO accounts (id, name, account_type_id, initial_amount, initial_date, current_balance, created_by) VALUES
//		(1, 'Cash', 1, 1000.00, CURDATE(), 1000.00, 1),
//		(2, 'Savings Account', 2, 5000.00, CURDATE(), 5000.00, 1);`,
//	}
//
//	for _, schema := range schemas {
//		_, err := db.Exec(schema)
//		if err != nil {
//			log.Printf("Error executing schema: %v", err)
//			return fmt.Errorf("error executing schema: %v", err)
//		}
//	}
//
//	return nil
//}

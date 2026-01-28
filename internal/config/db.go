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
	if err := initSchema(db); err != nil {
		return nil, fmt.Errorf("error initializing database schema: %v", err)
	}

	fmt.Println("Successfully connected to and initialized the database")
	return db, nil
}

func initSchema(db *sql.DB) error {
	schemaSQL := `
-- Create database if not exists
CREATE DATABASE IF NOT EXISTS finance;
USE finance;

-- Tables
CREATE TABLE IF NOT EXISTS account_types (
  id int(11) NOT NULL AUTO_INCREMENT,
  name varchar(50) NOT NULL,
  description text DEFAULT NULL,
  created_at timestamp NULL DEFAULT current_timestamp(),
  updated_at timestamp NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  created_by int(11) NOT NULL,
  updated_by int(11) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY name (name)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS accounts (
  id int(11) NOT NULL AUTO_INCREMENT,
  name varchar(100) NOT NULL,
  account_type_id int(11) NOT NULL,
  initial_amount decimal(14,2) NOT NULL DEFAULT 0.00,
  initial_date date NOT NULL,
  remarks text DEFAULT NULL,
  credit_limit decimal(14,2) DEFAULT NULL,
  credit_bill_date date DEFAULT NULL,
  credit_payment_date date DEFAULT NULL,
  current_balance decimal(15,2) DEFAULT 0.00,
  created_at timestamp NULL DEFAULT current_timestamp(),
  updated_at timestamp NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  created_by int(11) NOT NULL,
  updated_by int(11) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY account_type_id (account_type_id),
  CONSTRAINT accounts_ibfk_1 FOREIGN KEY (account_type_id) REFERENCES account_types (id)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS users (
  id int(11) NOT NULL AUTO_INCREMENT,
  user_name varchar(100) NOT NULL UNIQUE,
  pin varchar(4) NOT NULL,
  role varchar(20) DEFAULT 'user',
  created_at timestamp NULL DEFAULT current_timestamp(),
  updated_at timestamp NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (id)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS transaction_types (
  id int(11) NOT NULL AUTO_INCREMENT,
  name varchar(50) NOT NULL,
  description text DEFAULT NULL,
  created_at timestamp NULL DEFAULT current_timestamp(),
  updated_at timestamp NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  created_by int(11) NOT NULL,
  updated_by int(11) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY name (name)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS categories (
  id int(11) NOT NULL AUTO_INCREMENT,
  name varchar(100) NOT NULL,
  type_id int(11) NOT NULL,
  description text DEFAULT NULL,
  created_at timestamp NULL DEFAULT current_timestamp(),
  updated_at timestamp NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  created_by int(11) NOT NULL,
  updated_by int(11) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY type_id (type_id),
  CONSTRAINT categories_ibfk_1 FOREIGN KEY (type_id) REFERENCES transaction_types (id)
) ENGINE=InnoDB AUTO_INCREMENT=24 DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS transactions (
  id int(11) NOT NULL AUTO_INCREMENT,
  type_id int(11) NOT NULL,
  amount decimal(12,2) NOT NULL CHECK (amount >= 0),
  account_id int(11) NOT NULL,
  transaction_date date NOT NULL,
  from_account_id int(11) DEFAULT NULL,
  to_account_id int(11) DEFAULT NULL,
  category_id int(11) DEFAULT NULL,
  remarks text DEFAULT NULL,
  created_at datetime NOT NULL DEFAULT current_timestamp(),
  updated_at datetime NOT NULL DEFAULT current_timestamp(),
  created_by int(11) NOT NULL,
  updated_by int(11) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY from_account_id (from_account_id),
  KEY to_account_id (to_account_id),
  KEY idx_transaction_date (transaction_date),
  KEY idx_account_id (account_id),
  KEY idx_category_id (category_id),
  KEY idx_type_id (type_id),
  KEY idx_amount (amount),
  CONSTRAINT transactions_ibfk_1 FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE NO ACTION,
  CONSTRAINT transactions_ibfk_2 FOREIGN KEY (from_account_id) REFERENCES accounts (id) ON DELETE NO ACTION,
  CONSTRAINT transactions_ibfk_3 FOREIGN KEY (to_account_id) REFERENCES accounts (id) ON DELETE NO ACTION,
  CONSTRAINT transactions_ibfk_4 FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE SET NULL,
  CONSTRAINT transactions_ibfk_5 FOREIGN KEY (type_id) REFERENCES transaction_types (id) ON DELETE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS budgets (
  id int(11) NOT NULL AUTO_INCREMENT,
  user_id int(11) NOT NULL,
  category_id int(11) DEFAULT NULL,
  amount decimal(12,2) NOT NULL,
  month int(11) NOT NULL,
  year int(11) NOT NULL,
  created_at datetime NOT NULL DEFAULT current_timestamp(),
  updated_at datetime NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (id),
  KEY user_id (user_id),
  KEY category_id (category_id),
  CONSTRAINT budgets_ibfk_1 FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
  CONSTRAINT budgets_ibfk_2 FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Dummy Data
INSERT IGNORE INTO users (id, user_name, pin, role, created_at, updated_at) VALUES
  (1, 'admin', '1234', 'admin', NOW(), NOW()),
  (2, 'user1', '1111', 'user', NOW(), NOW());

INSERT IGNORE INTO transaction_types (id, name, description, created_at, updated_at, created_by, updated_by) VALUES
  (1, 'income', NULL, NOW(), NOW(), 0, NULL),
  (2, 'expense', NULL, NOW(), NOW(), 0, NULL),
  (3, 'investment', NULL, NOW(), NOW(), 0, NULL),
  (4, 'loan_payment', NULL, NOW(), NOW(), 0, NULL),
  (5, 'transfer', NULL, NOW(), NOW(), 0, NULL),
  (6, 'correction', NULL, NOW(), NOW(), 0, NULL);

INSERT IGNORE INTO account_types (id, name, description, created_at, updated_at, created_by, updated_by) VALUES
  (1, 'cash', 'Physical currency in hand', NOW(), NOW(), 0, NULL),
  (2, 'savings', 'Savings bank account', NOW(), NOW(), 0, NULL),
  (3, 'current', 'Current bank account', NOW(), NOW(), 0, NULL),
  (4, 'wallet', 'Mobile or digital wallets', NOW(), NOW(), 0, NULL),
  (5, 'credit_card', 'Credit card account', NOW(), NOW(), 0, NULL),
  (6, 'credit_line', 'Buy Now Pay Later, overdraft, or other short-term borrowing', NOW(), NOW(), 0, NULL),
  (7, 'loan', 'Loan accounts – personal, auto, home, etc.', NOW(), NOW(), 0, NULL),
  (8, 'investment_equity', 'Stocks, mutual funds, ETFs', NOW(), NOW(), 0, NULL),
  (9, 'investment_fixed', 'FDs, RDs, bonds, debentures', NOW(), NOW(), 0, NULL),
  (10, 'investment_retirement', 'NPS, PPF, EPF, etc.', NOW(), NOW(), 0, NULL),
  (11, 'investment_insurance', 'LIC, ULIP, endowment plans', NOW(), NOW(), 0, NULL),
  (12, 'other', 'Catch-all for weird shit', NOW(), NOW(), 0, NULL);

INSERT IGNORE INTO accounts (id, name, account_type_id, initial_amount, initial_date, remarks, credit_limit, credit_bill_date, credit_payment_date, current_balance, created_at, updated_at, created_by, updated_by) VALUES
  (1, 'Cash', 1, 1000.00, CURDATE(), 'Hard Cash INR', NULL, NULL, NULL, 1000.00, NOW(), NOW(), 1, NULL),
  (2, 'Savings Account', 2, 5000.00, CURDATE(), 'Main Savings', NULL, NULL, NULL, 5000.00, NOW(), NOW(), 1, NULL);

INSERT IGNORE INTO categories (id, name, type_id, description, created_at, updated_at, created_by, updated_by) VALUES
  (1, 'Salary', 1, 'Monthly salary', NOW(), NOW(), 0, NULL),
  (2, 'Freelance', 1, 'Freelance income', NOW(), NOW(), 0, NULL),
  (3, 'Food', 2, 'Daily meals', NOW(), NOW(), 0, NULL),
  (4, 'Transport', 2, 'Travel expenses', NOW(), NOW(), 0, NULL);

INSERT IGNORE INTO transactions (id, type_id, amount, account_id, transaction_date, from_account_id, to_account_id, category_id, remarks, created_at, updated_at, created_by, updated_by) VALUES
  (1, 1, 5000.00, 1, CURDATE(), NULL, NULL, 1, 'Salary', NOW(), NOW(), 1, NULL),
  (2, 2, 200.00, 1, CURDATE(), NULL, NULL, 3, 'Lunch', NOW(), NOW(), 1, NULL);

INSERT IGNORE INTO budgets (id, user_id, category_id, amount, month, year, created_at, updated_at) VALUES
  (1, 1, 3, 500.00, MONTH(CURDATE()), YEAR(CURDATE()), NOW(), NOW()),
  (2, 1, 6, 1000.00, MONTH(CURDATE()), YEAR(CURDATE()), NOW(), NOW());
`

	_, err := db.Exec(schemaSQL)
	if err != nil {
		return fmt.Errorf("error executing schema: %v", err)
	}

	return nil
}

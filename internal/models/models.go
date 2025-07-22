package models

import "time"

// Account represents the accounts table
	type Account struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	AccountTypeID     int       `json:"account_type_id"`
	InitialAmount     float64   `json:"initial_amount"`
	InitialDate       time.Time `json:"initial_date"`
	Remarks           *string   `json:"remarks"`
	CreditLimit       *float64  `json:"credit_limit"`
	CreditBillDate    *time.Time `json:"credit_bill_date"`
	CreditPaymentDate *time.Time `json:"credit_payment_date"`
	CurrentBalance    float64   `json:"current_balance"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	CreatedBy         int       `json:"created_by"`
	UpdatedBy         *int      `json:"updated_by"`
}

// AccountType represents the account_types table
	type AccountType struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   int       `json:"created_by"`
	UpdatedBy   *int      `json:"updated_by"`
}

// Category represents the categories table
	type Category struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	TypeID      int       `json:"type_id"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   int       `json:"created_by"`
	UpdatedBy   *int      `json:"updated_by"`
}

// Transaction represents the transactions table
	type Transaction struct {
	ID              int       `json:"id"`
	TypeID          int       `json:"type_id"`
	Amount          float64   `json:"amount"`
	AccountID       int       `json:"account_id"`
	AccountName     *string   `json:"account_name"`
	TransactionDate time.Time `json:"transaction_date"`
	FromAccountID   *int      `json:"from_account_id"`
	FromAccountName *string   `json:"from_account_name"`
	ToAccountID     *int      `json:"to_account_id"`
	ToAccountName   *string   `json:"to_account_name"`
	CategoryID      *int      `json:"category_id"`
	CategoryName    *string   `json:"category_name"`
	Remarks         *string   `json:"remarks"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedBy       int       `json:"created_by"`
	UpdatedBy       *int      `json:"updated_by"`
}

type User struct {
	ID        int       `json:"id"`
	Pin       string    `json:"pin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TransactionFilter struct {
	AccountID int    `json:"account_id"`
	Month     string `json:"month"`
	Year      string `json:"year"`
	TypeID    string `json:"type_id"`
}

// TransactionType represents the transaction_types table
	type TransactionType struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   int       `json:"created_by"`
	UpdatedBy   *int      `json:"updated_by"`
}

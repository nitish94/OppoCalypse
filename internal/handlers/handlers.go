package handlers

import (
	"OppoCalypse/internal/config"
	"OppoCalypse/internal/models"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"unsafe"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// ShowLoginPage displays the login form
func ShowLoginPage(c *gin.Context) {
	// Check if already logged in
	session := sessions.Default(c)
	if user := session.Get("user"); user != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	// Set context for the login page
	c.HTML(http.StatusOK, "login.tmpl", gin.H{
		"title":   "Login - OppoCalypse",
		"hideNav": true, // Hide navigation bar on login page
	})
}

// Login handles user login
func Login(c *gin.Context) {
	session := sessions.Default(c)

	// If already logged in, redirect to home
	if user := session.Get("user"); user != nil {
		fmt.Printf("User already logged in with ID: %v\n", user)
		c.Redirect(http.StatusFound, "/")
		return
	}

	// Only handle POST requests
	if c.Request.Method != http.MethodPost {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"hideNav": true,
		})
		return
	}

	pin := c.PostForm("pin")
	fmt.Printf("Login attempt with PIN: %s\n", pin)

	if len(pin) != 4 {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"error":   "PIN must be 4 digits",
			"hideNav": true,
		})
		return
	}

	fmt.Println("Connecting to database...")
	db, err := config.ConnectDB()
	if err != nil {
		fmt.Printf("Database connection error: %v\n", err)
		c.HTML(http.StatusInternalServerError, "login.tmpl", gin.H{
			"error":   "Database connection error",
			"hideNav": true,
		})
		return
	}
	defer db.Close()
	fmt.Println("Database connection successful.")

	var userID int
	var username string
	fmt.Println("Running user query...")
	err = db.QueryRow("SELECT id, user_name FROM users WHERE pin = ?", pin).Scan(&userID, &username)
	fmt.Printf("Query result: err=%v, userID=%v, username=%v\n", err, userID, username)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("No user found with PIN: %s\n", pin)
			c.HTML(http.StatusOK, "login.tmpl", gin.H{
				"error":   "Invalid PIN",
				"hideNav": true,
			})
		} else {
			fmt.Printf("Database error during login: %v\n", err)
			c.HTML(http.StatusInternalServerError, "login.tmpl", gin.H{
				"error":   "Internal server error",
				"hideNav": true,
			})
		}
		return
	}

	// Set session values
	session.Set("user", userID)
	session.Set("username", username)

	// Save the session
	if err := session.Save(); err != nil {
		fmt.Printf("Error saving session: %v\n", err)
		c.HTML(http.StatusInternalServerError, "login.tmpl", gin.H{
			"error":   "Failed to create session",
			"hideNav": true,
		})
		return
	}

	fmt.Printf("Login successful for user ID: %d, redirecting to /\n", userID)

	// Add cache control headers
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	// Redirect to home page with a 303 status code to force a GET request
	c.Redirect(http.StatusSeeOther, "/")
}

// Logout handles user logout
func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/login")
}

// GetTransactions retrieves and displays a list of transactions
func GetTransactions(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user")
	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Debug: Print user ID
	fmt.Printf("Getting transactions for user ID: %v\n", userID)

	db, err := config.ConnectDB()
	if err != nil {
		fmt.Printf("Database connection error: %v\n", err)
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error":   "Error connecting to database",
			"hideNav": false,
		})
		return
	}
	defer db.Close()

	// Parse filter parameters
	filter := models.TransactionFilter{
		AccountID: 0,
	}

	// Get account ID filter
	accountIDStr := c.Query("account_id")
	if accountIDStr != "" {
		accountID, err := strconv.Atoi(accountIDStr)
		if err == nil {
			filter.AccountID = accountID
		}
	}

	// Get month filter
	filter.Month = c.Query("month")

	// Get year filter
	filter.Year = c.Query("year")

	// Get type filter
	filter.TypeID = c.Query("type_id")

	// Initialize transactions as an empty slice
	transactions := []models.Transaction{}

	// Check if transactions table exists
	var tableExists int
	err = db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'transactions'").Scan(&tableExists)
	if err != nil {
		fmt.Printf("Error checking if transactions table exists: %v\n", err)
	}

	// If table exists, try to query transactions
	if tableExists > 0 {
		// Build the query with filters
		query := `
			SELECT t.id, t.type_id, t.amount, t.account_id, t.transaction_date, t.category_id, 
			COALESCE(c.name, '') as category_name, t.remarks, t.created_at, 
			a.name as account_name, 
			from_a.name as from_account_name, 
			to_a.name as to_account_name,
			t.from_account_id, t.to_account_id
			FROM transactions t
			LEFT JOIN categories c ON t.category_id = c.id
			LEFT JOIN accounts a ON t.account_id = a.id
			LEFT JOIN accounts from_a ON t.from_account_id = from_a.id
			LEFT JOIN accounts to_a ON t.to_account_id = to_a.id
			WHERE 1=1`

		// Add filters
		args := []interface{}{}

		// Account filter
		if filter.AccountID > 0 {
			query += ` AND (t.account_id = ? OR t.from_account_id = ? OR t.to_account_id = ?)`
			args = append(args, filter.AccountID, filter.AccountID, filter.AccountID)
		}

		// Month filter
		if filter.Month != "" {
			query += ` AND MONTH(t.transaction_date) = ?`
			args = append(args, filter.Month)
		}

		// Year filter
		if filter.Year != "" {
			query += ` AND YEAR(t.transaction_date) = ?`
			args = append(args, filter.Year)
		}

		// Type filter
		if filter.TypeID != "" {
			query += ` AND t.type_id = ?`
			args = append(args, filter.TypeID)
		}

		// Add ordering
		query += ` ORDER BY t.transaction_date DESC, t.created_at DESC`

		fmt.Printf("Executing query: %s with args: %v\n", query, args)

		rows, err := db.Query(query, args...)
		if err != nil {
			fmt.Printf("Error querying transactions: %v\n", err)
		} else {
			defer rows.Close()

			for rows.Next() {
				var t models.Transaction
				var categoryName, accountName, fromAccountName, toAccountName sql.NullString
				var fromAccountID, toAccountID sql.NullInt64
				if err := rows.Scan(
					&t.ID, &t.TypeID, &t.Amount, &t.AccountID, &t.TransactionDate, 
					&t.CategoryID, &categoryName, &t.Remarks, &t.CreatedAt,
					&accountName, &fromAccountName, &toAccountName,
					&fromAccountID, &toAccountID); err != nil {
					fmt.Printf("Error scanning transaction row: %v\n", err)
					continue
				}
				if categoryName.Valid {
					t.CategoryName = &categoryName.String
				}
				if accountName.Valid {
					t.AccountName = &accountName.String
				}
				if fromAccountName.Valid {
					t.FromAccountName = &fromAccountName.String
				}
				if toAccountName.Valid {
					t.ToAccountName = &toAccountName.String
				}
				if fromAccountID.Valid {
					t.FromAccountID = (*int)(unsafe.Pointer(&fromAccountID.Int64))
				}
				if toAccountID.Valid {
					t.ToAccountID = (*int)(unsafe.Pointer(&toAccountID.Int64))
				}
				transactions = append(transactions, t)
			}

			if err = rows.Err(); err != nil {
				fmt.Printf("Error iterating through transactions: %v\n", err)
			}
		}
	}

	// Get user info for the dashboard
	var username string
	err = db.QueryRow("SELECT user_name FROM users WHERE id = ?", userID).Scan(&username)
	if err != nil {
		fmt.Printf("Error getting username: %v\n", err)
		username = "User"
	}

	// Get accounts for filter dropdown
	accounts := []models.Account{}
	accountRows, err := db.Query("SELECT id, name FROM accounts ORDER BY name")
	if err == nil {
		defer accountRows.Close()
		for accountRows.Next() {
			var acc models.Account
			if err := accountRows.Scan(&acc.ID, &acc.Name); err == nil {
				accounts = append(accounts, acc)
			}
		}
	}

	// Get available years for filter dropdown
	years := []string{}
	yearRows, err := db.Query("SELECT DISTINCT YEAR(transaction_date) FROM transactions ORDER BY YEAR(transaction_date) DESC")
	if err == nil {
		defer yearRows.Close()
		for yearRows.Next() {
			var year string
			if err := yearRows.Scan(&year); err == nil {
				years = append(years, year)
			}
		}
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title":        "Dashboard - OppoCalypse",
		"transactions": transactions,
		"username":     username,
		"user":         userID,
		"hideNav":      false,
		"accounts":     accounts,
		"years":        years,
		"filter":       filter,
	})
}

// NewTransactionForm displays the form for creating a new transaction
func NewTransactionForm(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user")
	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	db, err := config.ConnectDB()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error":   "Error connecting to database",
			"hideNav": false,
		})
		return
	}
	defer db.Close()

	var categories []models.Category
	rows, err := db.Query("SELECT id, name, type_id FROM categories")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error":   "Error querying categories",
			"hideNav": false,
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.TypeID); err != nil {
			c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
				"error":   "Error scanning category",
				"hideNav": false,
			})
			return
		}
		categories = append(categories, cat)
	}

	var accounts []models.Account
	rows, err = db.Query("SELECT id, name FROM accounts")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error":   "Error querying accounts",
			"hideNav": false,
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var acc models.Account
		if err := rows.Scan(&acc.ID, &acc.Name); err != nil {
			c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
				"error":   "Error scanning account",
				"hideNav": false,
			})
			return
		}
		accounts = append(accounts, acc)
	}

	c.HTML(http.StatusOK, "form.tmpl", gin.H{
		"action":     "/transactions",
		"now":        time.Now().Format("2006-01-02"),
		"categories": categories,
		"accounts":   accounts,
		"user":       userID,
		"hideNav":    false,
	})
}

// CreateTransaction creates a new transaction
func CreateTransaction(c *gin.Context) {
	session := sessions.Default(c)
	userID, ok := session.Get("user").(int)
	if !ok {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	db, err := config.ConnectDB()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error connecting to database")
		return
	}
	defer db.Close()

	// Parse common fields
	amount, _ := strconv.ParseFloat(c.PostForm("amount"), 64)
	typeID, _ := strconv.Atoi(c.PostForm("type_id"))
	transactionDate, _ := time.Parse("2006-01-02", c.PostForm("transaction_date"))
	remarks := c.PostForm("remarks")
	categoryIDStr := c.PostForm("category_id")

	var categoryID sql.NullInt64
	if categoryIDStr != "" {
		id, err := strconv.Atoi(categoryIDStr)
		if err == nil {
			categoryID = sql.NullInt64{Int64: int64(id), Valid: true}
		}
	}

	var fromAccountID sql.NullInt64
	fromAccountIDStr := c.PostForm("from_account_id")
	if fromAccountIDStr != "" {
		id, err := strconv.Atoi(fromAccountIDStr)
		if err == nil {
			fromAccountID = sql.NullInt64{Int64: int64(id), Valid: true}
		}
	}

	var toAccountID sql.NullInt64
	toAccountIDStr := c.PostForm("to_account_id")
	if toAccountIDStr != "" {
		id, err := strconv.Atoi(toAccountIDStr)
		if err == nil {
			toAccountID = sql.NullInt64{Int64: int64(id), Valid: true}
		}
	}

	var accountID int64

	// Determine account_id based on transaction type
	switch typeID {
	case 1: // Income
		if !toAccountID.Valid {
			c.String(http.StatusBadRequest, "To Account is required for income")
			return
		}
		accountID = toAccountID.Int64
		fromAccountID.Valid = false // Ensure from is null
	case 2: // Expense
		if !fromAccountID.Valid {
			c.String(http.StatusBadRequest, "From Account is required for expense")
			return
		}
		accountID = fromAccountID.Int64
		toAccountID.Valid = false // Ensure to is null
	case 5: // Transfer
		if !fromAccountID.Valid || !toAccountID.Valid {
			c.String(http.StatusBadRequest, "From and To Accounts are required for transfer")
			return
		}
		// For transfer, let's assume the 'from' account is the primary one for the `account_id` field.
		accountID = fromAccountID.Int64
		categoryID.Valid = false // No category for transfers
	default:
		c.String(http.StatusBadRequest, "Invalid transaction type")
		return
	}

	// Insert the transaction
	_, err = db.Exec(`
		INSERT INTO transactions 
		(type_id, amount, account_id, transaction_date, from_account_id, to_account_id, category_id, remarks, created_by) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, 
		typeID, amount, accountID, transactionDate, fromAccountID, toAccountID, categoryID, remarks, userID)

	if err != nil {
		fmt.Println("Error creating transaction:", err)
		c.String(http.StatusInternalServerError, "Error creating transaction")
		return
	}

	c.Redirect(http.StatusFound, "/")
}

// EditTransactionForm displays the form for editing a transaction
func EditTransactionForm(c *gin.Context) {
	// Implementation for editing a transaction form
}

// UpdateTransaction updates an existing transaction
func UpdateTransaction(c *gin.Context) {
	// Implementation for updating a transaction
}

// DeleteTransaction deletes a transaction
func DeleteTransaction(c *gin.Context) {
	// Implementation for deleting a transaction
}

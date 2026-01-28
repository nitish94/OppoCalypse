package handlers

import (
	"OppoCalypse/internal/config"
	"OppoCalypse/internal/models"
	"database/sql"
	"encoding/csv"
	"encoding/json"
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

	username := c.PostForm("username")
	pin := c.PostForm("pin")
	fmt.Printf("Login attempt with Username: %s, PIN: %s\n", username, pin)

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
	var userName string
	var role string
	fmt.Println("Running user query...")
	err = db.QueryRow("SELECT id, user_name, role FROM users WHERE user_name = ? AND pin = ?", username, pin).Scan(&userID, &userName, &role)
	fmt.Printf("Query result: err=%v, userID=%v, username=%v, role=%v\n", err, userID, userName, role)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("No user found with Username: %s, PIN: %s\n", username, pin)
			c.HTML(http.StatusOK, "login.tmpl", gin.H{
				"error":   "Invalid username or PIN",
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
	session.Set("username", userName)
	session.Set("role", role)

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
	role := session.Get("role")
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
			JOIN user_accounts ua ON ua.account_id = t.account_id AND ua.user_id = ?
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

		// Add user_id to args
		args = append([]interface{}{userID}, args...)

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

	// Get accounts for filter dropdown (only accessible to user)
	accounts := []models.Account{}
	accountRows, err := db.Query(`
		SELECT a.id, a.name FROM accounts a
		JOIN user_accounts ua ON a.id = ua.account_id AND ua.user_id = ?
		ORDER BY a.name`, userID)
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
	yearRows, err := db.Query(`
		SELECT DISTINCT YEAR(t.transaction_date) FROM transactions t
		JOIN user_accounts ua ON ua.account_id = t.account_id AND ua.user_id = ?
		ORDER BY YEAR(t.transaction_date) DESC`, userID)
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
		"role":         role,
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
	role := session.Get("role")
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
	rows, err = db.Query(`
		SELECT a.id, a.name FROM accounts a
		JOIN user_accounts ua ON a.id = ua.account_id AND ua.user_id = ?`, userID)
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
		"role":       role,
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

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error starting transaction")
		return
	}
	defer tx.Rollback()

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
	_, err = tx.Exec(`
		INSERT INTO transactions 
		(type_id, amount, account_id, transaction_date, from_account_id, to_account_id, category_id, remarks, created_by) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, 
		typeID, amount, accountID, transactionDate, fromAccountID, toAccountID, categoryID, remarks, userID)

	if err != nil {
		fmt.Println("Error creating transaction:", err)
		c.String(http.StatusInternalServerError, "Error creating transaction")
		return
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		c.String(http.StatusInternalServerError, "Error committing transaction")
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

// Admin Handlers

// ListUsers displays a list of all users (admin only)
func ListUsers(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user")
	role := session.Get("role")

	db, err := config.ConnectDB()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error":   "Database connection error",
			"hideNav": false,
		})
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, user_name, role, created_at FROM users ORDER BY id")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error":   "Error fetching users",
			"hideNav": false,
		})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.UserName, &u.Role, &u.CreatedAt); err != nil {
			continue
		}
		users = append(users, u)
	}

	c.HTML(http.StatusOK, "admin_users.tmpl", gin.H{
		"title":   "Manage Users - OppoCalypse",
		"users":   users,
		"user":    userID,
		"role":    role,
		"hideNav": false,
	})
}

// NewUserForm displays the form to create a new user
func NewUserForm(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user")
	role := session.Get("role")

	c.HTML(http.StatusOK, "admin_user_form.tmpl", gin.H{
		"title":   "Add New User - OppoCalypse",
		"action":  "/admin/users",
		"user":    userID,
		"role":    role,
		"hideNav": false,
	})
}

// CreateUser creates a new user
func CreateUser(c *gin.Context) {
	userName := c.PostForm("user_name")
	pin := c.PostForm("pin")
	role := c.PostForm("role")

	if role == "" {
		role = "user"
	}

	db, err := config.ConnectDB()
	if err != nil {
		c.String(http.StatusInternalServerError, "Database error")
		return
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO users (user_name, pin, role) VALUES (?, ?, ?)", userName, pin, role)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating user")
		return
	}

	c.Redirect(http.StatusFound, "/admin/users")
}

// EditUserForm displays the form to edit a user
func EditUserForm(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user")
	role := session.Get("role")

	id := c.Param("id")
	db, err := config.ConnectDB()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error":   "Database error",
			"hideNav": false,
		})
		return
	}
	defer db.Close()

	var u models.User
	err = db.QueryRow("SELECT id, user_name, role FROM users WHERE id = ?", id).Scan(&u.ID, &u.UserName, &u.Role)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error":   "User not found",
			"hideNav": false,
		})
		return
	}

	c.HTML(http.StatusOK, "admin_user_form.tmpl", gin.H{
		"title":   "Edit User - OppoCalypse",
		"user":    u,
		"action":  "/admin/users/" + id,
		"current_user": userID,
		"role":    role,
		"hideNav": false,
	})
}

// UpdateUser updates a user
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	userName := c.PostForm("user_name")
	role := c.PostForm("role")

	db, err := config.ConnectDB()
	if err != nil {
		c.String(http.StatusInternalServerError, "Database error")
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE users SET user_name = ?, role = ? WHERE id = ?", userName, role, id)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error updating user")
		return
	}

	c.Redirect(http.StatusFound, "/admin/users")
}

// DeleteUser deletes a user
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	db, err := config.ConnectDB()
	if err != nil {
		c.String(http.StatusInternalServerError, "Database error")
		return
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error deleting user")
		return
	}

	c.Redirect(http.StatusFound, "/admin/users")
}

// ResetUserPIN resets a user's PIN
func ResetUserPIN(c *gin.Context) {
	id := c.Param("id")
	newPIN := c.PostForm("new_pin")

	db, err := config.ConnectDB()
	if err != nil {
		c.String(http.StatusInternalServerError, "Database error")
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE users SET pin = ? WHERE id = ?", newPIN, id)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error resetting PIN")
		return
	}

	c.Redirect(http.StatusFound, "/admin/users")
}

// Admin Account Management

// ListAccounts displays accounts and their assigned users
func ListAccounts(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user")
	role := session.Get("role")

	db, err := config.ConnectDB()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error":   "Database connection error",
			"hideNav": false,
		})
		return
	}
	defer db.Close()

	// Get all accounts with their types
	accounts := []models.Account{}
	accountRows, err := db.Query(`
		SELECT a.id, a.name, at.name as type_name, a.current_balance
		FROM accounts a
		JOIN account_types at ON a.account_type_id = at.id
		ORDER BY a.name`)
	if err == nil {
		defer accountRows.Close()
		for accountRows.Next() {
			var acc models.Account
			var typeName string
			if err := accountRows.Scan(&acc.ID, &acc.Name, &typeName, &acc.CurrentBalance); err == nil {
				acc.Remarks = &typeName // Reuse Remarks for type
				accounts = append(accounts, acc)
			}
		}
	}

	// Get all users
	users := []models.User{}
	userRows, err := db.Query("SELECT id, user_name FROM users ORDER BY user_name")
	if err == nil {
		defer userRows.Close()
		for userRows.Next() {
			var u models.User
			if err := userRows.Scan(&u.ID, &u.UserName); err == nil {
				users = append(users, u)
			}
		}
	}

	// Get account-user assignments
	selected := make(map[int]map[int]bool) // account_id -> user_id -> selected
	assignRows, err := db.Query("SELECT account_id, user_id FROM user_accounts")
	if err == nil {
		defer assignRows.Close()
		for assignRows.Next() {
			var accountID, userID int
			assignRows.Scan(&accountID, &userID)
			if selected[accountID] == nil {
				selected[accountID] = make(map[int]bool)
			}
			selected[accountID][userID] = true
		}
	}

	c.HTML(http.StatusOK, "admin_accounts.tmpl", gin.H{
		"title":     "Account Management - OppoCalypse",
		"accounts":  accounts,
		"users":     users,
		"selected":  selected,
		"user":      userID,
		"role":      role,
		"hideNav":   false,
	})
}

// AssignAccountUsers assigns users to accounts
func AssignAccountUsers(c *gin.Context) {
	accountIDStr := c.PostForm("account_id")
	userIDs := c.PostFormArray("user_ids[]")

	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid account ID")
		return
	}

	db, err := config.ConnectDB()
	if err != nil {
		c.String(http.StatusInternalServerError, "Database error")
		return
	}
	defer db.Close()

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		c.String(http.StatusInternalServerError, "Transaction error")
		return
	}
	defer tx.Rollback()

	// Remove existing assignments for this account
	_, err = tx.Exec("DELETE FROM user_accounts WHERE account_id = ?", accountID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error removing assignments")
		return
	}

	// Add new assignments
	for _, userIDStr := range userIDs {
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			continue
		}
		_, err = tx.Exec("INSERT INTO user_accounts (user_id, account_id) VALUES (?, ?)", userID, accountID)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error assigning user")
			return
		}
	}

	// Commit
	if err = tx.Commit(); err != nil {
		c.String(http.StatusInternalServerError, "Commit error")
		return
	}

	c.Redirect(http.StatusFound, "/admin/accounts")
}

// ShowExportPage displays the export form
func ShowExportPage(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user")
	role := session.Get("role")
	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	db, err := config.ConnectDB()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error":   "Database connection error",
			"hideNav": false,
		})
		return
	}
	defer db.Close()

	// Get accounts (accessible to user)
	accounts := []models.Account{}
	accountRows, err := db.Query(`
		SELECT a.id, a.name FROM accounts a
		JOIN user_accounts ua ON a.id = ua.account_id AND ua.user_id = ?
		ORDER BY a.name`, userID)
	if err == nil {
		defer accountRows.Close()
		for accountRows.Next() {
			var acc models.Account
			if err := accountRows.Scan(&acc.ID, &acc.Name); err == nil {
				accounts = append(accounts, acc)
			}
		}
	}

	// Get categories
	categories := []models.Category{}
	catRows, err := db.Query("SELECT id, name FROM categories ORDER BY name")
	if err == nil {
		defer catRows.Close()
		for catRows.Next() {
			var cat models.Category
			if err := catRows.Scan(&cat.ID, &cat.Name); err == nil {
				categories = append(categories, cat)
			}
		}
	}

	c.HTML(http.StatusOK, "export.tmpl", gin.H{
		"title":      "Export Transactions - OppoCalypse",
		"accounts":   accounts,
		"categories": categories,
		"user":       userID,
		"role":       role,
		"hideNav":    false,
	})
}

// ExportTransactions exports transactions to CSV based on filters
func ExportTransactions(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user")
	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	db, err := config.ConnectDB()
	if err != nil {
		c.String(http.StatusInternalServerError, "Database error")
		return
	}
	defer db.Close()

	// Parse filters
	accountIDStr := c.PostForm("account_id")
	categoryIDStr := c.PostForm("category_id")
	fromDate := c.PostForm("from_date")
	toDate := c.PostForm("to_date")

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
		JOIN user_accounts ua ON ua.account_id = t.account_id AND ua.user_id = ?
		WHERE 1=1`

	args := []interface{}{}

	if accountIDStr != "" {
		accountID, err := strconv.Atoi(accountIDStr)
		if err == nil {
			query += ` AND (t.account_id = ? OR t.from_account_id = ? OR t.to_account_id = ?)`
			args = append(args, accountID, accountID, accountID)
		}
	}

	if categoryIDStr != "" {
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err == nil {
			query += ` AND t.category_id = ?`
			args = append(args, categoryID)
		}
	}

	if fromDate != "" {
		query += ` AND t.transaction_date >= ?`
		args = append(args, fromDate)
	}

	if toDate != "" {
		query += ` AND t.transaction_date <= ?`
		args = append(args, toDate)
	}

	query += ` ORDER BY t.transaction_date DESC`

	// Add user_id to args
	args = append([]interface{}{userID}, args...)

	rows, err := db.Query(query, args...)
	if err != nil {
		c.String(http.StatusInternalServerError, "Query error")
		return
	}
	defer rows.Close()

	// Set headers for CSV download
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=transactions.csv")

	// Write CSV
	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Header
	writer.Write([]string{"ID", "Type", "Amount", "Account", "Date", "Category", "Remarks", "From Account", "To Account", "Created At"})

	for rows.Next() {
		var t models.Transaction
		var categoryName, accountName, fromAccountName, toAccountName sql.NullString
		var fromAccountID, toAccountID sql.NullInt64
		if err := rows.Scan(
			&t.ID, &t.TypeID, &t.Amount, &t.AccountID, &t.TransactionDate,
			&t.CategoryID, &categoryName, &t.Remarks, &t.CreatedAt,
			&accountName, &fromAccountName, &toAccountName,
			&fromAccountID, &toAccountID); err != nil {
			continue
		}

		typeStr := ""
		switch t.TypeID {
		case 1:
			typeStr = "Income"
		case 2:
			typeStr = "Expense"
		case 3:
			typeStr = "Investment"
		case 4:
			typeStr = "Loan Payment"
		case 5:
			typeStr = "Transfer"
		case 6:
			typeStr = "Correction"
		}

		catName := ""
		if categoryName.Valid {
			catName = categoryName.String
		}

		accName := ""
		if accountName.Valid {
			accName = accountName.String
		}

		fromAcc := ""
		if fromAccountName.Valid {
			fromAcc = fromAccountName.String
		}

		toAcc := ""
		if toAccountName.Valid {
			toAcc = toAccountName.String
		}

		remarks := ""
		if t.Remarks != nil {
			remarks = *t.Remarks
		}

		writer.Write([]string{
			strconv.Itoa(t.ID),
			typeStr,
			fmt.Sprintf("%.2f", t.Amount),
			accName,
			t.TransactionDate.Format("2006-01-02"),
			catName,
			remarks,
			fromAcc,
			toAcc,
			t.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
}

// Budget Handlers

// ListBudgets displays a list of budgets
func ListBudgets(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user")
	role := session.Get("role")
	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	db, err := config.ConnectDB()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error":   "Database connection error",
			"hideNav": false,
		})
		return
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT b.id, b.user_id, b.category_id, b.amount, b.month, b.year, b.created_at, b.updated_at, c.name as category_name
		FROM budgets b
		LEFT JOIN categories c ON b.category_id = c.id
		WHERE b.user_id = ?
		ORDER BY b.year DESC, b.month DESC`, userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error":   "Error fetching budgets",
			"hideNav": false,
		})
		return
	}
	defer rows.Close()

	var budgets []models.Budget
	for rows.Next() {
		var b models.Budget
		var catName sql.NullString
		if err := rows.Scan(&b.ID, &b.UserID, &b.CategoryID, &b.Amount, &b.Month, &b.Year, &b.CreatedAt, &b.UpdatedAt, &catName); err != nil {
			continue
		}
		if catName.Valid {
			b.Category = &models.Category{Name: catName.String}
		}
		budgets = append(budgets, b)
	}

	c.HTML(http.StatusOK, "budgets.tmpl", gin.H{
		"title":   "Budgets - OppoCalypse",
		"budgets": budgets,
		"user":    userID,
		"role":    role,
		"hideNav": false,
	})
}

// NewBudgetForm displays the form to create a new budget
func NewBudgetForm(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user")
	role := session.Get("role")
	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	db, err := config.ConnectDB()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error":   "Database error",
			"hideNav": false,
		})
		return
	}
	defer db.Close()

	categories := []models.Category{}
	catRows, err := db.Query("SELECT id, name FROM categories ORDER BY name")
	if err == nil {
		defer catRows.Close()
		for catRows.Next() {
			var cat models.Category
			if err := catRows.Scan(&cat.ID, &cat.Name); err == nil {
				categories = append(categories, cat)
			}
		}
	}

	months := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

	c.HTML(http.StatusOK, "budget_form.tmpl", gin.H{
		"title":      "New Budget - OppoCalypse",
		"action":     "/budgets",
		"categories": categories,
		"months":     months,
		"user":       userID,
		"role":       role,
		"now":        time.Now(),
		"hideNav":    false,
	})
}

// CreateBudget creates a new budget
func CreateBudget(c *gin.Context) {
	session := sessions.Default(c)
	userID, ok := session.Get("user").(int)
	if !ok {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	categoryIDStr := c.PostForm("category_id")
	amount, _ := strconv.ParseFloat(c.PostForm("amount"), 64)
	month, _ := strconv.Atoi(c.PostForm("month"))
	year, _ := strconv.Atoi(c.PostForm("year"))

	var categoryID sql.NullInt64
	if categoryIDStr != "" {
		id, err := strconv.Atoi(categoryIDStr)
		if err == nil {
			categoryID = sql.NullInt64{Int64: int64(id), Valid: true}
		}
	}

	db, err := config.ConnectDB()
	if err != nil {
		c.String(http.StatusInternalServerError, "Database error")
		return
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO budgets (user_id, category_id, amount, month, year) VALUES (?, ?, ?, ?, ?)",
		userID, categoryID, amount, month, year)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating budget")
		return
	}

	c.Redirect(http.StatusFound, "/budgets")
}

// EditBudgetForm displays the form to edit a budget
func EditBudgetForm(c *gin.Context) {
	id := c.Param("id")
	session := sessions.Default(c)
	userID := session.Get("user")
	role := session.Get("role")
	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	db, err := config.ConnectDB()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error":   "Database error",
			"hideNav": false,
		})
		return
	}
	defer db.Close()

	var b models.Budget
	err = db.QueryRow("SELECT id, user_id, category_id, amount, month, year FROM budgets WHERE id = ? AND user_id = ?", id, userID).Scan(
		&b.ID, &b.UserID, &b.CategoryID, &b.Amount, &b.Month, &b.Year)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error":   "Budget not found",
			"hideNav": false,
		})
		return
	}

	categories := []models.Category{}
	catRows, err := db.Query("SELECT id, name FROM categories ORDER BY name")
	if err == nil {
		defer catRows.Close()
		for catRows.Next() {
			var cat models.Category
			if err := catRows.Scan(&cat.ID, &cat.Name); err == nil {
				categories = append(categories, cat)
			}
		}
	}

	months := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

	c.HTML(http.StatusOK, "budget_form.tmpl", gin.H{
		"title":      "Edit Budget - OppoCalypse",
		"budget":     b,
		"action":     "/budgets/" + id,
		"categories": categories,
		"months":     months,
		"user":       userID,
		"role":       role,
		"now":        time.Now(),
		"hideNav":    false,
	})
}

// UpdateBudget updates a budget
func UpdateBudget(c *gin.Context) {
	id := c.Param("id")
	session := sessions.Default(c)
	userID, ok := session.Get("user").(int)
	if !ok {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	categoryIDStr := c.PostForm("category_id")
	amount, _ := strconv.ParseFloat(c.PostForm("amount"), 64)
	month, _ := strconv.Atoi(c.PostForm("month"))
	year, _ := strconv.Atoi(c.PostForm("year"))

	var categoryID sql.NullInt64
	if categoryIDStr != "" {
		id, err := strconv.Atoi(categoryIDStr)
		if err == nil {
			categoryID = sql.NullInt64{Int64: int64(id), Valid: true}
		}
	}

	db, err := config.ConnectDB()
	if err != nil {
		c.String(http.StatusInternalServerError, "Database error")
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE budgets SET category_id = ?, amount = ?, month = ?, year = ? WHERE id = ? AND user_id = ?",
		categoryID, amount, month, year, id, userID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error updating budget")
		return
	}

	c.Redirect(http.StatusFound, "/budgets")
}

// DeleteBudget deletes a budget
func DeleteBudget(c *gin.Context) {
	id := c.Param("id")
	session := sessions.Default(c)
	userID, ok := session.Get("user").(int)
	if !ok {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	db, err := config.ConnectDB()
	if err != nil {
		c.String(http.StatusInternalServerError, "Database error")
		return
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM budgets WHERE id = ? AND user_id = ?", id, userID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error deleting budget")
		return
	}

	c.Redirect(http.StatusFound, "/budgets")
}

// Account Handlers

// NewAccountForm displays the form to create a new account
func NewAccountForm(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user")
	role := session.Get("role")
	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	db, err := config.ConnectDB()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error":   "Database connection error",
			"hideNav": false,
		})
		return
	}
	defer db.Close()

	accountTypes := []models.AccountType{}
	atRows, err := db.Query("SELECT id, name FROM account_types ORDER BY name")
	if err == nil {
		defer atRows.Close()
		for atRows.Next() {
			var at models.AccountType
			if err := atRows.Scan(&at.ID, &at.Name); err == nil {
				accountTypes = append(accountTypes, at)
			}
		}
	}

	c.HTML(http.StatusOK, "account_form.tmpl", gin.H{
		"title":         "New Account - OppoCalypse",
		"action":        "/accounts",
		"account_types": accountTypes,
		"user":          userID,
		"role":          role,
		"hideNav":       false,
	})
}

// CreateAccount creates a new account and assigns to user
func CreateAccount(c *gin.Context) {
	session := sessions.Default(c)
	userID, ok := session.Get("user").(int)
	if !ok {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	name := c.PostForm("name")
	accountTypeIDStr := c.PostForm("account_type_id")
	initialAmount, _ := strconv.ParseFloat(c.PostForm("initial_amount"), 64)
	initialDate := c.PostForm("initial_date")
	remarks := c.PostForm("remarks")

	accountTypeID, err := strconv.Atoi(accountTypeIDStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid account type")
		return
	}

	db, err := config.ConnectDB()
	if err != nil {
		c.String(http.StatusInternalServerError, "Database error")
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		c.String(http.StatusInternalServerError, "Transaction error")
		return
	}
	defer tx.Rollback()

	// Insert account
	result, err := tx.Exec(`
		INSERT INTO accounts (name, account_type_id, initial_amount, initial_date, current_balance, remarks, created_by)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		name, accountTypeID, initialAmount, initialDate, initialAmount, remarks, userID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating account")
		return
	}

	accountID, err := result.LastInsertId()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting account ID")
		return
	}

	// Assign to user
	_, err = tx.Exec("INSERT INTO user_accounts (user_id, account_id) VALUES (?, ?)", userID, accountID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error assigning account")
		return
	}

	if err = tx.Commit(); err != nil {
		c.String(http.StatusInternalServerError, "Commit error")
		return
	}

	c.Redirect(http.StatusFound, "/")
}

// NewAdminAccountForm displays the form to create a new account (admin)
func NewAdminAccountForm(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user")
	role := session.Get("role")

	db, err := config.ConnectDB()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error":   "Database connection error",
			"hideNav": false,
		})
		return
	}
	defer db.Close()

	accountTypes := []models.AccountType{}
	atRows, err := db.Query("SELECT id, name FROM account_types ORDER BY name")
	if err == nil {
		defer atRows.Close()
		for atRows.Next() {
			var at models.AccountType
			if err := atRows.Scan(&at.ID, &at.Name); err == nil {
				accountTypes = append(accountTypes, at)
			}
		}
	}

	users := []models.User{}
	userRows, err := db.Query("SELECT id, user_name FROM users ORDER BY user_name")
	if err == nil {
		defer userRows.Close()
		for userRows.Next() {
			var u models.User
			if err := userRows.Scan(&u.ID, &u.UserName); err == nil {
				users = append(users, u)
			}
		}
	}

	c.HTML(http.StatusOK, "account_form.tmpl", gin.H{
		"title":         "New Account - OppoCalypse",
		"action":        "/admin/accounts",
		"account_types": accountTypes,
		"users":         users,
		"user":          userID,
		"role":          role,
		"hideNav":       false,
	})
}

// CreateAdminAccount creates a new account (admin can assign users)
func CreateAdminAccount(c *gin.Context) {
	name := c.PostForm("name")
	accountTypeIDStr := c.PostForm("account_type_id")
	initialAmount, _ := strconv.ParseFloat(c.PostForm("initial_amount"), 64)
	initialDate := c.PostForm("initial_date")
	remarks := c.PostForm("remarks")
	userIDs := c.PostFormArray("user_ids[]")

	accountTypeID, err := strconv.Atoi(accountTypeIDStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid account type")
		return
	}

	session := sessions.Default(c)
	adminID, ok := session.Get("user").(int)
	if !ok {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	db, err := config.ConnectDB()
	if err != nil {
		c.String(http.StatusInternalServerError, "Database error")
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		c.String(http.StatusInternalServerError, "Transaction error")
		return
	}
	defer tx.Rollback()

	// Insert account
	result, err := tx.Exec(`
		INSERT INTO accounts (name, account_type_id, initial_amount, initial_date, current_balance, remarks, created_by)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		name, accountTypeID, initialAmount, initialDate, initialAmount, remarks, adminID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating account")
		return
	}

	accountID, err := result.LastInsertId()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting account ID")
		return
	}

	// Assign to selected users
	for _, userIDStr := range userIDs {
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			continue
		}
		_, err = tx.Exec("INSERT INTO user_accounts (user_id, account_id) VALUES (?, ?)", userID, accountID)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error assigning user")
			return
		}
	}

	if err = tx.Commit(); err != nil {
		c.String(http.StatusInternalServerError, "Commit error")
		return
	}

	c.Redirect(http.StatusFound, "/admin/accounts")
}

// ShowGraphs displays graphs
func ShowGraphs(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user")
	role := session.Get("role")
	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	db, err := config.ConnectDB()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
			"error":   "Database connection error",
			"hideNav": false,
		})
		return
	}
	defer db.Close()

	// Data for line/bar chart: income and expense by month
	var incomeData []map[string]interface{}
	var expenseData []map[string]interface{}

	incomeRows, err := db.Query(`
		SELECT YEAR(t.transaction_date) as year, MONTH(t.transaction_date) as month, SUM(t.amount) as total
		FROM transactions t
		JOIN user_accounts ua ON ua.account_id = t.account_id AND ua.user_id = ?
		WHERE t.type_id = 1
		GROUP BY YEAR(t.transaction_date), MONTH(t.transaction_date)
		ORDER BY year, month`, userID)
	if err == nil {
		defer incomeRows.Close()
		for incomeRows.Next() {
			var year, month int
			var total float64
			incomeRows.Scan(&year, &month, &total)
			incomeData = append(incomeData, map[string]interface{}{
				"label": fmt.Sprintf("%d-%02d", year, month),
				"value": total,
			})
		}
	}

	expenseRows, err := db.Query(`
		SELECT YEAR(t.transaction_date) as year, MONTH(t.transaction_date) as month, SUM(t.amount) as total
		FROM transactions t
		JOIN user_accounts ua ON ua.account_id = t.account_id AND ua.user_id = ?
		WHERE t.type_id = 2
		GROUP BY YEAR(t.transaction_date), MONTH(t.transaction_date)
		ORDER BY year, month`, userID)
	if err == nil {
		defer expenseRows.Close()
		for expenseRows.Next() {
			var year, month int
			var total float64
			expenseRows.Scan(&year, &month, &total)
			expenseData = append(expenseData, map[string]interface{}{
				"label": fmt.Sprintf("%d-%02d", year, month),
				"value": total,
			})
		}
	}

	// Data for pie charts: income by category
	var incomeCatData []map[string]interface{}
	incomeCatRows, err := db.Query(`
		SELECT c.name, SUM(t.amount) as total
		FROM transactions t
		JOIN categories c ON t.category_id = c.id
		JOIN user_accounts ua ON ua.account_id = t.account_id AND ua.user_id = ?
		WHERE t.type_id = 1
		GROUP BY c.id, c.name
		ORDER BY total DESC`, userID)
	if err == nil {
		defer incomeCatRows.Close()
		for incomeCatRows.Next() {
			var name string
			var total float64
			incomeCatRows.Scan(&name, &total)
			incomeCatData = append(incomeCatData, map[string]interface{}{
				"label": name,
				"value": total,
			})
		}
	}

	// Expense by category
	var expenseCatData []map[string]interface{}
	expenseCatRows, err := db.Query(`
		SELECT c.name, SUM(t.amount) as total
		FROM transactions t
		JOIN categories c ON t.category_id = c.id
		JOIN user_accounts ua ON ua.account_id = t.account_id AND ua.user_id = ?
		WHERE t.type_id = 2
		GROUP BY c.id, c.name
		ORDER BY total DESC`, userID)
	if err == nil {
		defer expenseCatRows.Close()
		for expenseCatRows.Next() {
			var name string
			var total float64
			expenseCatRows.Scan(&name, &total)
			expenseCatData = append(expenseCatData, map[string]interface{}{
				"label": name,
				"value": total,
			})
		}
	}

	incomeDataJSON, _ := json.Marshal(incomeData)
	expenseDataJSON, _ := json.Marshal(expenseData)
	incomeCatDataJSON, _ := json.Marshal(incomeCatData)
	expenseCatDataJSON, _ := json.Marshal(expenseCatData)

	c.HTML(http.StatusOK, "graphs.tmpl", gin.H{
		"title":             "Graphs - OppoCalypse",
		"incomeData":        incomeData,
		"expenseData":       expenseData,
		"incomeCatData":     incomeCatData,
		"expenseCatData":    expenseCatData,
		"incomeDataJSON":    string(incomeDataJSON),
		"expenseDataJSON":   string(expenseDataJSON),
		"incomeCatDataJSON": string(incomeCatDataJSON),
		"expenseCatDataJSON": string(expenseCatDataJSON),
		"user":              userID,
		"role":              role,
		"hideNav":           false,
	})
}

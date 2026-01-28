# OppoCalypse - Personal Finance Manager

A comprehensive web-based personal finance management application built with Go, Gin, and MySQL/MariaDB. Features transaction tracking, budgeting, user management, and data visualization.

## Features

### Core Functionality
- **User Authentication**: Secure login with username and PIN
- **Transaction Management**: Add, edit, delete income, expense, transfer, and correction transactions
- **Account Management**: Manage multiple accounts (cash, savings, credit cards, etc.)
- **Category Management**: Organize transactions with customizable categories
- **Budgeting**: Set monthly budgets per category or overall
- **Data Export**: Export transactions to CSV with filters
- **Graphs & Analytics**: Visualize income/expense trends and category breakdowns

### Admin Features
- **User Management**: Add, edit, delete users, reset PINs
- **Account Access Control**: Assign accounts to users (support for joint accounts)
- **System Administration**: Full control over users and account permissions

### Advanced Features
- **Joint Accounts**: Share accounts between users (admin-controlled)
- **Role-Based Access**: Admin and regular user roles
- **Session Management**: Secure session handling
- **Responsive UI**: Clean, modern interface

## Installation

### Prerequisites
- Go 1.19+
- MySQL/MariaDB 5.7+
- Git

### Setup

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd OppoCalypse
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Configure environment:**
   ```bash
   cp .env-example .env
   # Edit .env with your settings
   ```

4. **Database setup:**
   - Create a MySQL/MariaDB database named `finance`
   - The application will automatically create tables and insert dummy data on first run

5. **Build and run:**
   ```bash
   go build -o oppocalypse ./cmd/server
   ./oppocalypse
   ```

6. **Access the application:**
   - Open http://localhost:3000 in your browser
   - Default admin login: username `admin`, PIN `1234`

## Configuration

### Environment Variables (.env)

| Variable      | Description              | Default          |
|---------------|--------------------------|------------------|
| PORT          | Server port              | 3000             |
| DB_USER       | Database username        | root             |
| DB_PASSWORD   | Database password        | (empty)          |
| DB_HOST       | Database host            | 127.0.0.1       |
| DB_PORT       | Database port            | 3306             |
| DB_NAME       | Database name            | finance          |
| SESSION_SECRET| Session encryption key   | default-secret   |

## Database Schema

### Core Tables
- `users`: User accounts with authentication
- `accounts`: Financial accounts (cash, bank, cards)
- `categories`: Transaction categories
- `transactions`: Financial transactions
- `transaction_types`: Income, expense, transfer, etc.
- `account_types`: Cash, savings, credit, etc.
- `budgets`: Monthly budget settings
- `user_accounts`: Many-to-many user-account relationships (for joint accounts)

## API Endpoints

### Public
- `GET /login` - Login page
- `POST /login` - Authenticate user

### Authenticated
- `GET /` - Dashboard with transactions
- `GET /transactions/new` - New transaction form
- `POST /transactions` - Create transaction
- `GET /transactions/:id/edit` - Edit transaction form
- `POST /transactions/:id` - Update transaction
- `POST /transactions/:id/delete` - Delete transaction
- `GET /budgets` - Budget list
- `GET /budgets/new` - New budget form
- `POST /budgets` - Create budget
- `GET /budgets/:id/edit` - Edit budget form
- `POST /budgets/:id` - Update budget
- `POST /budgets/:id/delete` - Delete budget
- `GET /export` - Export form
- `POST /export` - Export CSV
- `GET /graphs` - Analytics graphs
- `POST /logout` - Logout

### Admin Only
- `GET /admin/users` - User management
- `GET /admin/users/new` - New user form
- `POST /admin/users` - Create user
- `GET /admin/users/:id/edit` - Edit user form
- `POST /admin/users/:id` - Update user
- `POST /admin/users/:id/delete` - Delete user
- `POST /admin/users/:id/reset-pin` - Reset user PIN
- `GET /admin/accounts` - Account management
- `POST /admin/accounts/assign` - Assign users to accounts

## User Roles

### Regular User
- Manage personal transactions
- View accessible accounts
- Set personal budgets
- Export data
- View graphs

### Admin
- All user permissions
- Manage all users (add, edit, delete, reset PIN)
- Manage account access (assign users to accounts)
- Create joint accounts

## Security

- PIN-based authentication (4-digit)
- Session-based authorization
- Role-based access control
- SQL injection prevention with prepared statements
- CSRF protection via sessions

## Development

### Project Structure
```
OppoCalypse/
├── cmd/server/          # Application entry point
├── internal/
│   ├── config/          # Database configuration
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # Gin middleware
│   ├── models/          # Data models
│   └── routes/          # Route definitions
├── static/              # CSS, JS, images
├── templates/           # HTML templates
├── .env-example         # Environment template
├── go.mod               # Go modules
└── README2.md           # This file
```

### Adding New Features
1. Define models in `internal/models/`
2. Add handlers in `internal/handlers/`
3. Define routes in `internal/routes/`
4. Create templates in `templates/`
5. Update database schema in `internal/config/db.go`

### Testing
- Run with `go run cmd/server/main.go`
- Test endpoints with curl or Postman
- Check logs for errors

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make changes with proper documentation
4. Test thoroughly
5. Submit a pull request

## License

MIT License - see LICENSE file for details

## Support

For issues or questions, please create an issue in the repository.
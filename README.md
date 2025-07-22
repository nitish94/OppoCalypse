# OppoCalypse - Personal Finance Manager

OppoCalypse is a lightweight, self-hosted personal finance management application built with Go and MySQL. It helps you track income, expenses, and transfers across multiple accounts.

## Features

- **Simple Authentication**: PIN-based login system
- **Transaction Management**: Track income, expenses, and transfers
- **Multiple Account Types**: Support for cash, bank, credit cards, and more
- **Transaction Categories**: Organize transactions by categories
- **Filtering**: Filter transactions by account, month, year, and type
- **Responsive Design**: Works on desktop and mobile devices

## Screenshots

(Add screenshots here)

## Requirements

- Go 1.18 or higher
- MySQL 5.7 or higher
- Modern web browser

## Installation

1. **Clone the repository**

```bash
git clone https://github.com/yourusername/oppocalypse.git
cd oppocalypse
```

2. **Set up the database**

### Option A: Local MySQL/MariaDB

Create a MySQL database and import the schema:

```bash
mysql -u root -p
CREATE DATABASE finance;
USE finance;
source honefinance_final.sql;
```

### Option B: Using Termux on Android

You can host the database on your Android device using Termux:

#### Android Termux Setup (One-time)

1. **Install MariaDB**

```bash
pkg update && pkg upgrade -y
pkg install mariadb -y
```

2. **Initialize Database**

```bash
mariadb-install-db
```

3. **Allow Remote Access**

Create/edit config file:

```bash
nano $PREFIX/etc/my.cnf.d/network.cnf
```

Add:

```ini
[mysqld]
bind-address=0.0.0.0
port=3306
```

4. **Start MariaDB**

```bash
pkill mariadbd
mariadbd-safe &
```

5. **Create Remote User**

```bash
mysql -u root -e "
DROP USER IF EXISTS 'nitish'@'%';
CREATE USER 'nitish'@'%' IDENTIFIED BY 'admin123';
GRANT ALL PRIVILEGES ON *.* TO 'nitish'@'%' WITH GRANT OPTION;
FLUSH PRIVILEGES;"
```

6. **Create Database and Import Schema**

```bash
mysql -u root -e "CREATE DATABASE finance;"
mysql -u root finance < honefinance_final.sql
```

7. **Get Android IP**

```bash
ip addr show wlan0
```

Look for `inet 192.168.x.x` → **copy this IP**

8. **Restarting Termux MariaDB**

Each time you restart Termux, just run:

```bash
mariadbd-safe &
```

3. **Configure the database connection**

Edit `internal/config/db.go` to match your database credentials:

```go
// For local MySQL/MariaDB
dsn := "username:password@tcp(127.0.0.1:3306)/finance?parseTime=true&multiStatements=true"

// For Termux MariaDB (replace with your Android device IP)
dsn := "nitish:admin123@tcp(192.168.x.x:3306)/finance?parseTime=true&multiStatements=true"
```

4. **Build and run the application**

```bash
go build -o oppocalypse ./cmd/server
./oppocalypse
```

5. **Access the application**

Open your browser and navigate to:
```
http://localhost:8081
```

Default login PIN: `0000` (user: Nitish Haldar)

## Usage

### Adding Transactions

1. Log in with your PIN
2. Click "New Transaction" in the navigation bar
3. Select transaction type (Income, Expense, or Transfer)
4. Fill in the details and click "Create"

### Filtering Transactions

1. Click the "Filter" button on the dashboard
2. Select your desired filters (account, month, year, type)
3. Click "Apply Filter" to see filtered results
4. Click "Reset" to clear all filters

## Project Structure

```
OppoCalypse/
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── config/          # Database configuration
│   ├── db/              # SQL schema files
│   ├── handlers/        # HTTP request handlers
│   ├── middleware/      # Authentication middleware
│   ├── models/          # Data models
│   └── routes/          # Route definitions
├── static/
│   ├── css/             # Stylesheets
│   └── js/              # JavaScript files
├── templates/           # HTML templates
├── go.mod               # Go module definition
└── go.sum               # Go module checksums
```

## Roadmap

### Planned Features

- **Export Functionality**: Export transactions to CSV/Excel
- **Budget Planning**: Set and track monthly budgets
- **Data Visualization**: Charts and graphs for financial analysis
- **Multi-user Support**: Family account management
- **Recurring Transactions**: Schedule regular transactions
- **Mobile App**: Native mobile applications

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
-- --------------------------------------------------------
-- Host:                         192.168.29.78
-- Server version:               11.8.2-MariaDB - MariaDB Server
-- Server OS:                    Android
-- HeidiSQL Version:             12.11.0.7065
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

-- -----------------------------------------------------------------------------
-- Database: `finance`
-- -----------------------------------------------------------------------------
CREATE DATABASE IF NOT EXISTS `finance` /*!40100 DEFAULT CHARACTER SET utf8mb4 */;
USE `finance`;

-- -----------------------------------------------------------------------------
-- Table: `account_types`
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `account_types` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `description` text DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `created_by` int(11) NOT NULL,
  `updated_by` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4;

-- -----------------------------------------------------------------------------
-- Table: `accounts`
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `accounts` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `account_type_id` int(11) NOT NULL,
  `initial_amount` decimal(14,2) NOT NULL DEFAULT 0.00,
  `initial_date` date NOT NULL,
  `remarks` text DEFAULT NULL,
  `credit_limit` decimal(14,2) DEFAULT NULL,
  `credit_bill_date` date DEFAULT NULL,
  `credit_payment_date` date DEFAULT NULL,
  `current_balance` decimal(15,2) DEFAULT 0.00,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `created_by` int(11) NOT NULL,
  `updated_by` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `account_type_id` (`account_type_id`),
  CONSTRAINT `accounts_ibfk_1` FOREIGN KEY (`account_type_id`) REFERENCES `account_types` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4;

-- -----------------------------------------------------------------------------
-- Table: `users`
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_name` varchar(100) NOT NULL,
  `pin` varchar(4) NOT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;

-- -----------------------------------------------------------------------------
-- Table: `transaction_types`
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `transaction_types` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `description` text DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `created_by` int(11) NOT NULL,
  `updated_by` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4;

-- -----------------------------------------------------------------------------
-- Table: `categories`
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `categories` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `type_id` int(11) NOT NULL,
  `description` text DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `created_by` int(11) NOT NULL,
  `updated_by` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `type_id` (`type_id`),
  CONSTRAINT `categories_ibfk_1` FOREIGN KEY (`type_id`) REFERENCES `transaction_types` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=24 DEFAULT CHARSET=utf8mb4;

-- -----------------------------------------------------------------------------
-- Table: `transactions`
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `transactions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `type_id` int(11) NOT NULL,
  `amount` decimal(12,2) NOT NULL CHECK (`amount` >= 0),
  `account_id` int(11) NOT NULL,
  `transaction_date` date NOT NULL,
  `from_account_id` int(11) DEFAULT NULL,
  `to_account_id` int(11) DEFAULT NULL,
  `category_id` int(11) DEFAULT NULL,
  `remarks` text DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT current_timestamp(),
  `updated_at` datetime NOT NULL DEFAULT current_timestamp(),
  `created_by` int(11) NOT NULL,
  `updated_by` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `from_account_id` (`from_account_id`),
  KEY `to_account_id` (`to_account_id`),
  KEY `idx_transaction_date` (`transaction_date`),
  KEY `idx_account_id` (`account_id`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_type_id` (`type_id`),
  KEY `idx_amount` (`amount`),
  CONSTRAINT `transactions_ibfk_1` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`) ON DELETE NO ACTION,
  CONSTRAINT `transactions_ibfk_2` FOREIGN KEY (`from_account_id`) REFERENCES `accounts` (`id`) ON DELETE NO ACTION,
  CONSTRAINT `transactions_ibfk_3` FOREIGN KEY (`to_account_id`) REFERENCES `accounts` (`id`) ON DELETE NO ACTION,
  CONSTRAINT `transactions_ibfk_4` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`) ON DELETE SET NULL,
  CONSTRAINT `transactions_ibfk_5` FOREIGN KEY (`type_id`) REFERENCES `transaction_types` (`id`) ON DELETE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4;

-- -----------------------------------------------------------------------------
-- Data: `users`
-- -----------------------------------------------------------------------------
INSERT INTO `users` (`id`, `user_name`, `pin`, `created_at`, `updated_at`) VALUES
  (1, 'Nitish Haldar', '0000', '2025-07-17 18:25:42', '2025-07-17 18:25:42'),
  (2, 'Swagata Haldar', '1111', '2025-07-17 19:49:48', '2025-07-17 19:49:48');

-- -----------------------------------------------------------------------------
-- Data: `transaction_types`
-- -----------------------------------------------------------------------------
INSERT INTO `transaction_types` (`id`, `name`, `description`, `created_at`, `updated_at`, `created_by`, `updated_by`) VALUES
  (1, 'income', NULL, '2025-07-17 16:09:32', '2025-07-17 16:09:32', 0, NULL),
  (2, 'expense', NULL, '2025-07-17 16:09:32', '2025-07-17 16:09:32', 0, NULL),
  (3, 'investment', NULL, '2025-07-17 16:09:32', '2025-07-17 16:09:32', 0, NULL),
  (4, 'loan_payment', NULL, '2025-07-17 16:09:32', '2025-07-17 16:09:32', 0, NULL),
  (5, 'transfer', NULL, '2025-07-17 16:09:32', '2025-07-17 16:09:32', 0, NULL),
  (6, 'correction', NULL, '2025-07-17 16:09:32', '2025-07-17 16:09:32', 0, NULL);

-- -----------------------------------------------------------------------------
-- Data: `account_types`
-- -----------------------------------------------------------------------------
INSERT INTO `account_types` (`id`, `name`, `description`, `created_at`, `updated_at`, `created_by`, `updated_by`) VALUES
  (1, 'cash', 'Physical currency in hand', '2025-07-17 16:09:32', '2025-07-17 16:09:32', 0, NULL),
  (2, 'savings', 'Savings bank account', '2025-07-17 16:09:32', '2025-07-17 16:09:32', 0, NULL),
  (3, 'current', 'Current bank account', '2025-07-17 16:09:32', '2025-07-17 16:09:32', 0, NULL),
  (4, 'wallet', 'Mobile or digital wallets', '2025-07-17 16:09:32', '2025-07-17 16:09:32', 0, NULL),
  (5, 'credit_card', 'Credit card account', '2025-07-17 16:09:32', '2025-07-17 16:09:32', 0, NULL),
  (6, 'credit_line', 'Buy Now Pay Later, overdraft, or other short-term borrowing', '2025-07-17 16:09:32', '2025-07-17 16:09:32', 0, NULL),
  (7, 'loan', 'Loan accounts – personal, auto, home, etc.', '2025-07-17 16:09:32', '2025-07-17 16:09:32', 0, NULL),
  (8, 'investment_equity', 'Stocks, mutual funds, ETFs', '2025-07-17 16:09:32', '2025-07-17 16:09:32', 0, NULL),
  (9, 'investment_fixed', 'FDs, RDs, bonds, debentures', '2025-07-17 16:09:32', '2025-07-17 16:09:32', 0, NULL),
  (10, 'investment_retirement', 'NPS, PPF, EPF, etc.', '2025-07-17 16:09:32', '2025-07-17 16:09:32', 0, NULL),
  (11, 'investment_insurance', 'LIC, ULIP, endowment plans (a.k.a. legalized theft)', '2025-07-17 16:09:32', '2025-07-17 16:09:32', 0, NULL),
  (12, 'other', 'Catch-all for weird shit', '2025-07-17 16:09:32', '2025-07-17 16:09:32', 0, NULL);

-- -----------------------------------------------------------------------------
-- Data: `accounts`
-- -----------------------------------------------------------------------------
INSERT INTO `accounts` (`id`, `name`, `account_type_id`, `initial_amount`, `initial_date`, `remarks`, `credit_limit`, `credit_bill_date`, `credit_payment_date`, `current_balance`, `created_at`, `updated_at`, `created_by`, `updated_by`) VALUES
  (1, 'Cash', 1, 1000.00, '2025-07-18', 'Hard Cash INR', NULL, NULL, NULL, 1000.00, '2025-07-17 19:11:06', '2025-07-18 09:51:52', 1, NULL),
  (3, 'Bank of Baroda', 2, 35557.00, '2025-04-29', 'Salary | NG5D00000', NULL, NULL, NULL, 0.00, '2025-07-18 07:25:09', '2025-07-18 07:25:09', 1, NULL),
  (4, 'Jupiter', 2, 7133.43, '2025-04-29', '777700000000', NULL, NULL, NULL, 0.00, '2025-07-18 07:25:09', '2025-07-18 07:25:09', 1, NULL),
  (5, 'HDFC Bank', 2, 25167.00, '2025-04-29', '501000000000', NULL, NULL, NULL, 0.00, '2025-07-18 07:25:09', '2025-07-18 07:25:09', 1, NULL),
  (6, 'Axis Flipkart', 5, 99978.00, '2025-04-29', 'XX 1580', 101000.00, NULL, NULL, 0.00, '2025-07-18 07:25:09', '2025-07-18 07:25:09', 1, NULL),
  (7, 'Axis MyZone', 5, 100806.00, '2025-04-29', 'XX 2096', 101000.00, NULL, NULL, 0.00, '2025-07-18 07:25:09', '2025-07-18 07:25:09', 1, NULL),
  (8, 'SBI SimplyClick', 5, 46448.06, '2025-04-29', 'XX 6681', 68000.00, NULL, NULL, 0.00, '2025-07-18 07:25:09', '2025-07-18 07:25:09', 1, NULL),
  (9, 'ICICI Amazon Pay', 5, 50000.00, '2025-04-29', 'XX 6005', 50000.00, NULL, NULL, 0.00, '2025-07-18 07:25:09', '2025-07-18 07:25:09', 1, NULL),
  (10, 'HDFC Money Back', 5, 38662.00, '2025-04-29', 'XX 9930', 40000.00, NULL, NULL, 0.00, '2025-07-18 07:25:09', '2025-07-18 07:25:09', 1, NULL),
  (11, 'Cash Back and Wallet', 4, 3.00, '2025-04-29', 'Amazon Balance, CRED etc', NULL, NULL, NULL, 0.00, '2025-07-18 07:25:09', '2025-07-18 07:25:09', 1, NULL),
  (12, 'Amazon Pay Later', 6, 7000.00, '2025-04-29', '', 5000.00, NULL, NULL, 0.00, '2025-07-18 07:25:09', '2025-07-18 07:25:09', 1, NULL),
  (13, 'Cash', 1, 2000.00, '2025-04-29', '', NULL, NULL, NULL, 0.00, '2025-07-18 07:25:09', '2025-07-18 07:25:09', 1, NULL),
  (14, 'NPS', 10, 1500.00, '2025-04-29', 'PRAN - 1111111111', NULL, NULL, NULL, 0.00, '2025-07-18 07:25:09', '2025-07-18 07:25:09', 1, NULL),
  (15, 'PPF', 10, 1500.00, '2025-04-29', '550000000000', NULL, NULL, NULL, 0.00, '2025-07-18 07:25:09', '2025-07-18 07:25:09', 1, NULL);

-- -----------------------------------------------------------------------------
-- Data: `categories`
-- -----------------------------------------------------------------------------
INSERT INTO `categories` (`id`, `name`, `type_id`, `description`, `created_at`, `updated_at`, `created_by`, `updated_by`) VALUES
  (1, 'Salary', 1, 'Monthly or irregular salary', '2025-07-17 16:37:59', '2025-07-17 16:37:59', 0, NULL),
  (2, 'Cashback', 1, 'Cashback from purchases', '2025-07-17 16:37:59', '2025-07-17 16:37:59', 0, NULL),
  (3, 'Interest', 1, 'Bank interest, FD interest etc.', '2025-07-17 16:37:59', '2025-07-17 16:37:59', 0, NULL),
  (4, 'Repay', 1, 'Refunds or loan repayments received', '2025-07-17 16:37:59', '2025-07-17 16:37:59', 0, NULL),
  (5, 'Credit', 1, 'Money added through credit', '2025-07-17 16:37:59', '2025-07-17 16:37:59', 0, NULL),
  (6, 'Food', 2, 'Daily meals, snacks, beverages', '2025-07-17 16:37:59', '2025-07-17 16:37:59', 0, NULL),
  (7, 'Dining Out', 2, 'Restaurants, cafes, takeout', '2025-07-17 16:37:59', '2025-07-17 16:37:59', 0, NULL),
  (8, 'Phone Recharge', 2, 'Mobile top-ups and plans', '2025-07-17 16:37:59', '2025-07-17 16:37:59', 0, NULL),
  (9, 'Broadband Recharge', 2, 'Home internet bills', '2025-07-17 16:37:59', '2025-07-17 16:37:59', 0, NULL),
  (10, 'Electric Bill', 2, 'Electricity usage', '2025-07-17 16:37:59', '2025-07-17 16:37:59', 0, NULL),
  (11, 'Cloth', 2, 'Clothing and fashion', '2025-07-17 16:37:59', '2025-07-17 16:37:59', 0, NULL),
  (12, 'Travel', 2, 'Transport, trips, fuel', '2025-07-17 16:37:59', '2025-07-17 16:37:59', 0, NULL),
  (13, 'Household Items', 2, 'Non-food essential goods', '2025-07-17 16:37:59', '2025-07-17 16:37:59', 0, NULL),
  (14, 'Family Expense', 2, 'Spending for family needs', '2025-07-17 16:37:59', '2025-07-17 16:37:59', 0, NULL),
  (15, 'Household Commodities', 2, 'Cleaning, kitchen items', '2025-07-17 16:37:59', '2025-07-17 16:37:59', 0, NULL),
  (16, 'Grocery', 2, 'Vegetables, staples, food stock', '2025-07-17 16:37:59', '2025-07-17 16:37:59', 0, NULL),
  (17, 'Toiletry', 2, 'Toothpaste, soap, etc.', '2025-07-17 16:37:59', '2025-07-17 16:37:59', 0, NULL),
  (18, 'Medicine', 2, 'Pharmacy, first-aid, supplements', '2025-07-17 16:37:59', '2025-07-17 16:37:59', 0, NULL),
  (19, 'Other Bills', 2, 'Miscellaneous recurring bills', '2025-07-17 16:37:59', '2025-07-17 16:37:59', 0, NULL),
  (20, 'Skin Care & Cosmetics', 2, 'Beauty and skincare', '2025-07-17 16:37:59', '2025-07-17 16:37:59', 0, NULL),
  (21, 'Entertainment', 2, 'Movies, subscriptions, fun', '2025-07-17 16:37:59', '2025-07-17 16:37:59', 0, NULL),
  (22, 'Gifts', 2, 'Gifting to others', '2025-07-17 16:37:59', '2025-07-17 16:37:59', 0, NULL),
  (23, 'Lending', 2, 'Money lend to thers', '2025-07-17 16:37:59', '2025-07-17 16:41:59', 0, NULL);

-- -----------------------------------------------------------------------------
-- Data: `transactions`
-- -----------------------------------------------------------------------------
INSERT INTO `transactions` (`id`, `type_id`, `amount`, `account_id`, `transaction_date`, `from_account_id`, `to_account_id`, `category_id`, `remarks`, `created_at`, `updated_at`, `created_by`, `updated_by`) VALUES
  (1, 5, 2000.00, 3, '2025-07-18', 3, 4, NULL, 'tranfer', '2025-07-18 14:52:12', '2025-07-18 14:52:12', 1, NULL),
  (3, 2, 60.00, 1, '2025-07-18', 1, NULL, 6, 'Fuchka', '2025-07-18 15:31:39', '2025-07-18 15:31:39', 2, NULL);

/*!40103 SET TIME_ZONE=IFNULL(@OLD_TIME_ZONE, 'system') */;
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
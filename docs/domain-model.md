# Domain Model Documentation

## Overview

The Proletariat Budget system is built around a comprehensive domain model that represents the core financial concepts and their relationships. This document outlines the main domain entities, their attributes, and how they interact within the system.

## Core Domain Entities

### Account

Accounts represent financial containers that hold money, such as bank accounts, cash, credit cards, or investment accounts.

**Key Attributes:**
- ID: Unique identifier
- Name: Account name
- Type: Type of account (checking, savings, credit card, investment, cash)
- Balance: Current balance
- Currency: Currency of the account
- Active: Whether the account is active or archived

**Relationships:**
- Belongs to one or more Household Members
- Has many Transactions
- Can be the source or destination for Transfers
- Can be linked to Savings Goals

### Expenditure

Expenditures represent money spent from accounts.

**Key Attributes:**
- ID: Unique identifier
- Category: Category of the expenditure
- Date: When the expenditure occurred
- Amount: How much was spent
- Description: Description of the expenditure
- Declared: Whether the expenditure has been declared/confirmed
- Planned: Whether the expenditure was planned or unexpected

**Relationships:**
- Belongs to a Category
- Can have multiple Tags
- Linked to one or more Transactions

### Ingress

Ingresses represent money coming into accounts (income).

**Key Attributes:**
- ID: Unique identifier
- Category: Category of the income
- Source: Source of the income
- Date: When the income was received
- Amount: How much was received
- Description: Description of the income
- Is Recurring: Whether this income recurs regularly

**Relationships:**
- Belongs to a Category
- Can have multiple Tags
- Can have a Recurrence Pattern
- Linked to one or more Transactions

### Savings Goal

Savings Goals represent financial targets that users are saving towards.

**Key Attributes:**
- ID: Unique identifier
- Name: Name of the savings goal
- Category: Category of the goal
- Target Amount: Amount to be saved
- Current Amount: Current amount saved
- Target Date: When the goal should be achieved
- Priority: Relative importance of the goal
- Auto Contribute: Whether automatic contributions are enabled

**Relationships:**
- Belongs to a Category
- Linked to an Account
- Has many Contributions
- Has many Withdrawals
- Can have multiple Tags

### Transaction

Transactions represent the movement of money within the system.

**Key Attributes:**
- ID: Unique identifier
- Account ID: The account involved
- Amount: Amount of the transaction
- Transaction Date: When the transaction occurred
- Description: Description of the transaction
- Transaction Type: Type of transaction (expenditure, ingress, transfer, etc.)
- Balance After: Account balance after the transaction
- Status: Status of the transaction (pending, completed, failed, cancelled)

**Relationships:**
- Belongs to an Account
- Must be linked to an Expenditure, Ingress, Transfer, Savings Contribution, or Savings Withdrawal

### Household Member

Household Members represent individuals who are part of the financial household.

**Key Attributes:**
- ID: Unique identifier
- Name: Name of the household member
- Email: Email address
- Role: Role within the household

**Relationships:**
- Can own multiple Accounts
- Can be responsible for Expenditures
- Can be the source of Ingresses
- 
## Domain Relationships Diagram
                              +----------------+
                              | Household      |
                              | Member         |
                              +-------+--------+
                                      |
                                      | owns
                                      |
              +---------------------+-+--+--------------------+
              |                     |    |                    |
              v                     v    v                    v
     +--------+-------+    +--------+----+---+    +-----------+--------+
     |                |    |                 |    |                    |
     |    Account     |    |   Expenditure   |    |      Ingress       |
     |                |    |                 |    |                    |
     +--------+-------+    +---------+-------+    +-----------+--------+
              |                      |                        |
              | contains             | categorized by         | categorized by
              v                      v                        v
     +--------+-------+    +---------+-------+    +-----------+--------+
     |                |    |                 |    |                    |
     | Transaction    |    |    Category     |    |     Category       |
     |                |    |                 |    |                    |
     +----------------+    +-----------------+    +--------------------+
              |
              | can be
              v
     +--------+-------+
     |                |
     |   Transfer     |
     |                |
     +----------------+
              |
              | can fund
              v
     +--------+-------+
     |                |
     | Savings Goal   |
     |                |
     +--------+-------+
              |
              | has
      +-------+--------+
      |                |
      | Contribution/  |
      | Withdrawal     |
      +----------------+

## Value Objects

### Money

Money is a value object that represents a monetary amount in a specific currency.

**Attributes:**
- Amount: Decimal value
- Currency: Reference to a Currency entity

### Date Range

Date Range is a value object that represents a period between two dates.

**Attributes:**
- Start Date: Beginning of the range
- End Date: End of the range

### Recurrence Pattern

Recurrence Pattern is a value object that defines how often an event repeats.

**Attributes:**
- Frequency: How often the event occurs (daily, weekly, monthly, yearly)
- Interval: Number of frequency units between occurrences
- End Date: Optional date when the recurrence ends

## Domain Services

### Budget Calculation Service

Calculates budget information based on income, expenses, and savings goals.

### Financial Projection Service

Projects future financial status based on recurring income, expenses, and savings goals.

### Transaction Categorization Service

Automatically categorizes transactions based on patterns and rules.

### Currency Conversion Service

Handles conversion between different currencies for reporting and analysis.

## Aggregates

### Account Aggregate

The Account entity serves as the root of the Account aggregate, which includes Transactions related to the account.

### Savings Goal Aggregate

The Savings Goal entity serves as the root of the Savings Goal aggregate, which includes Contributions and Withdrawals.

## Domain Events

- AccountCreated
- TransactionRecorded
- SavingsGoalAchieved
- BudgetExceeded
- RecurringIncomeReceived

## Domain Rules

1. An account's balance must always reflect the sum of its transactions
2. A transfer must have both a source and destination account
3. A savings goal's current amount must not exceed its target amount
4. A withdrawal from a savings goal cannot exceed its current amount
5. A transaction must be associated with exactly one account
6. Recurring transactions must have a valid recurrence pattern
7. A household member must have at least one account
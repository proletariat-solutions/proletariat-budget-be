
-- Insert household members
INSERT INTO household_members (name, surname, nickname, role, active) VALUES
('John', 'Smith', 'Johnny', 'head_of_household', TRUE),
('Sarah', 'Smith', 'Sari', 'spouse', TRUE),
('Michael', 'Smith', 'Mike', 'child', TRUE);

-- Insert categories
INSERT INTO categories (name, description, color, background_color, active, category_type) VALUES
-- Expenditure categories
('Groceries', 'Food and household items', '#FF6B6B', '#FFE5E5', TRUE, 'expenditure'),
('Transportation', 'Car, gas, public transport', '#4ECDC4', '#E5F9F7', TRUE, 'expenditure'),
('Entertainment', 'Movies, games, dining out', '#45B7D1', '#E5F4FD', TRUE, 'expenditure'),
('Utilities', 'Electricity, water, internet', '#96CEB4', '#F0F9F4', TRUE, 'expenditure'),
('Healthcare', 'Medical expenses, insurance', '#FFEAA7', '#FFFBF0', TRUE, 'expenditure'),
-- Ingress categories
('Salary', 'Monthly salary income', '#74B9FF', '#E8F4FF', TRUE, 'ingress'),
('Freelance', 'Freelance work income', '#A29BFE', '#F1F0FF', TRUE, 'ingress'),
('Investment Returns', 'Dividends and returns', '#FD79A8', '#FFF0F6', TRUE, 'ingress'),
-- Saving goal categories
('Emergency Fund', 'Emergency savings', '#00B894', '#E8F8F5', TRUE, 'saving_goal'),
('Vacation', 'Holiday savings', '#FDCB6E', '#FFF8E8', TRUE, 'saving_goal'),
('Home Down Payment', 'Saving for house', '#E17055', '#FDF2F0', TRUE, 'saving_goal');

-- Insert tags
INSERT INTO tags (name, description, color, background_color, type) VALUES
-- Expenditure tags
('Essential', 'Necessary expenses', '#E74C3C', '#FADBD8', 'expenditure'),
('Luxury', 'Non-essential purchases', '#9B59B6', '#F4ECF7', 'expenditure'),
('Recurring', 'Regular monthly expenses', '#3498DB', '#EBF5FB', 'expenditure'),
('One-time', 'Single occurrence expense', '#F39C12', '#FEF9E7', 'expenditure'),
-- Ingress tags
('Primary Income', 'Main source of income', '#27AE60', '#E8F8F5', 'ingress'),
('Secondary Income', 'Additional income source', '#16A085', '#E8F6F3', 'ingress'),
('Bonus', 'Extra income', '#F1C40F', '#FEFBEA', 'ingress'),
-- Saving goal tags
('Short-term', 'Goals within 1 year', '#E67E22', '#FDF2E9', 'saving_goal'),
('Long-term', 'Goals over 1 year', '#8E44AD', '#F8F4FD', 'saving_goal'),
('High Priority', 'Important goals', '#C0392B', '#FADBD8', 'saving_goal'),
-- Transfer tags
('Internal', 'Between own accounts', '#34495E', '#EAEDED', 'transfer'),
('External', 'To external accounts', '#2C3E50', '#D5DBDB', 'transfer'),
-- Saving contribution tags
('Automatic', 'Auto contributions', '#1ABC9C', '#E8F8F5', 'saving_contribution'),
('Manual', 'Manual contributions', '#E74C3C', '#FADBD8', 'saving_contribution'),
-- Saving withdrawal tags
('Emergency', 'Emergency withdrawals', '#E74C3C', '#FADBD8', 'saving_withdrawal'),
('Goal Achieved', 'Goal completion withdrawal', '#27AE60', '#E8F8F5', 'saving_withdrawal');

-- Insert accounts (using USD currency id = 168)
INSERT INTO accounts (name, type, owner, institution, currency, initial_balance, current_balance, active, description, account_number) VALUES
('Main Checking', 'bank', 1, 'Chase Bank', 168, 5000.00, 7250.50, TRUE, 'Primary checking account', 'CHK-001-12345'),
('Savings Account', 'bank', 1, 'Chase Bank', 168, 10000.00, 12500.75, TRUE, 'High-yield savings account', 'SAV-001-67890'),
('Cash Wallet', 'cash', 2, NULL, 168, 200.00, 150.00, TRUE, 'Daily cash expenses', NULL);

-- Insert savings goals
INSERT INTO savings_goals (name, category_id, description, target_amount, currency, target_date, initial_amount, current_amount, percent_complete, account_id, priority, auto_contribute, auto_contribute_amount, auto_contribute_frequency, status) VALUES
('Emergency Fund', 9, 'Six months of expenses', 15000.00, 168, '2025-12-31', 2000.00, 5500.00, 36.67, 2, 1, TRUE, 500.00, 'monthly', 'active'),
('Summer Vacation 2025', 10, 'Family trip to Europe', 8000.00, 168, '2025-06-01', 1000.00, 3200.00, 40.00, 2, 2, TRUE, 300.00, 'monthly', 'active'),
('House Down Payment', 11, '20% down payment for new home', 50000.00, 168, '2027-01-01', 5000.00, 12000.00, 24.00, 2, 1, TRUE, 800.00, 'monthly', 'active');

-- Insert savings goal tags
INSERT INTO savings_goal_tags (savings_goal_id, tag_id) VALUES
(1, 8), (1, 10), -- Emergency Fund: Short-term, High Priority
(2, 8), -- Summer Vacation: Short-term
(3, 9), (3, 10); -- House Down Payment: Long-term, High Priority

-- Insert transactions
INSERT INTO transactions (account_id, amount, currency, transaction_date, description, transaction_type, balance_after, status) VALUES
-- Expenditures
(1, 150.75, 168, '2024-01-15 10:30:00', 'Weekly grocery shopping at Walmart', 'expenditure', 7099.75, 'completed'),
(1, 45.20, 168, '2024-01-16 14:20:00', 'Gas station fill-up', 'expenditure', 7054.55, 'completed'),
(1, 89.99, 168, '2024-01-17 19:45:00', 'Dinner at Italian restaurant', 'expenditure', 6964.56, 'completed'),
(3, 25.00, 168, '2024-01-18 12:00:00', 'Coffee and lunch', 'expenditure', 125.00, 'completed'),
(1, 120.00, 168, '2024-01-19 09:15:00', 'Monthly internet bill', 'expenditure', 6844.56, 'completed'),
-- Ingresses
(1, 3500.00, 168, '2024-01-01 08:00:00', 'Monthly salary deposit', 'ingress', 6500.00, 'completed'),
(1, 750.00, 168, '2024-01-10 16:30:00', 'Freelance project payment', 'ingress', 7250.00, 'completed'),
(2, 125.50, 168, '2024-01-05 10:00:00', 'Investment dividend', 'ingress', 12125.50, 'completed'),
-- Savings contributions
(2, 0.00, 168, '2024-01-02 09:00:00', 'Monthly emergency fund contribution', 'transfer', 12625.50, 'completed'),
(2, 0.00, 168, '2024-01-02 09:05:00', 'Monthly vacation fund contribution', 'transfer', 12925.50, 'completed'),
-- Transfers
(1, 200.00, 168, '2024-01-20 11:00:00', 'Transfer to savings account', 'transfer', 6644.56, 'completed'),-- Emergency withdrawal to use on medical services
(2, 0.00, 168, '2024-01-25 15:00:00', 'Emergency withdrawal to use on medical services', 'transfer', 12500.7, 'completed'),
(2, 100.00, 168, '2024-01-25 15:00:00', 'Emergency withdrawal to use on medical services', 'expenditure', 12400.7, 'completed');


-- Insert transfers
INSERT INTO transfers (source_account_id, destination_account_id, destination_amount, exchange_rate_multiplier, fees, transaction_id) VALUES
    (1, 2, 200.00, 1.000000, 0.00, 11), -- Transfer from checking to savings
    (2, 2, 200.00, 0, 0, 12), -- Transfer from savings to
    (2, 2, 200.00, 0, 0, 9),
    (2, 2, 200.00, 0, 0, 10);
-- Insert expenditures
INSERT INTO expenditures (category_id, declared, planned, transaction_id) VALUES
(1, TRUE, TRUE, 1),   -- Groceries
(2, TRUE, FALSE, 2),  -- Transportation (gas)
(3, FALSE, FALSE, 3), -- Entertainment (dinner)
(3, FALSE, FALSE, 4), -- Entertainment (coffee/lunch)
(4, TRUE, TRUE, 5),   -- Utilities (internet)
(1, true, false, 13);


-- Insert expenditure tags
INSERT INTO expenditure_tags (expenditure_id, tag_id) VALUES
(1, 1), (1, 3), -- Groceries: Essential, Recurring
(2, 1), -- Gas: Essential
(3, 2), (3, 4), -- Dinner: Luxury, One-time
(4, 2), -- Coffee/lunch: Luxury
(5, 1), (5, 3); -- Internet: Essential, Recurring

-- Insert ingresses
INSERT INTO ingresses (category_id, source, is_recurring, transaction_id) VALUES
(6, 'ABC Corporation', TRUE, 6),   -- Salary
(7, 'Freelance Client XYZ', FALSE, 7), -- Freelance
(8, 'Investment Portfolio', TRUE, 8);   -- Investment returns

-- Insert ingress tags
INSERT INTO ingress_tags (ingress_id, tag_id) VALUES
(1, 5), -- Salary: Primary Income
(2, 6), -- Freelance: Secondary Income
(3, 6), (3, 7); -- Investment: Secondary Income, Bonus

-- Insert ingress recurrence patterns
INSERT INTO ingress_recurrence_patterns (ingress_id, frequency, interval_value, amount) VALUES
(1, 'monthly', 1, 3500.00), -- Monthly salary
(3, 'monthly', 1, 125.50);  -- Monthly dividends

-- Insert savings contributions
INSERT INTO savings_contributions (savings_goal_id, transfer_id, date) VALUES
(1, 3, '2024-01-02'),
(2, 4, '2024-01-02');

insert into savings_withdrawals (savings_goal_id, date, reason, transfer_id) VALUES
(1, '2024-01-25', 'Emergency withdrawal', 11);


-- Insert savings contribution tags
INSERT INTO savings_contribution_tags (contribution_id, tag_id) VALUES
(1, 13), -- Emergency fund contribution: Automatic
(2, 13); -- Vacation fund contribution: Automatic

-- Update account balances and savings goal amounts based on transactions
UPDATE accounts SET current_balance = 6644.56 WHERE id = 1; -- Main Checking after all transactions
UPDATE accounts SET current_balance = 12500.75 WHERE id = 2; -- Savings after transfer and contributions
UPDATE accounts SET current_balance = 125.00 WHERE id = 3; -- Cash wallet after coffee expense

UPDATE savings_goals SET current_amount = 5500.00, percent_complete = 36.67 WHERE id = 1; -- Emergency fund
UPDATE savings_goals SET current_amount = 3200.00, percent_complete = 40.00 WHERE id = 2; -- Vacation fund
UPDATE savings_goals SET current_amount = 12000.00, percent_complete = 24.00 WHERE id = 3; -- House down payment
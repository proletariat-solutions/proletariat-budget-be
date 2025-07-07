SET FOREIGN_KEY_CHECKS = 0;
TRUNCATE TABLE proletariat_budget.accounts;
TRUNCATE TABLE proletariat_budget.categories;
TRUNCATE TABLE proletariat_budget.exchange_rates;
TRUNCATE TABLE proletariat_budget.expenditure_tags;
TRUNCATE TABLE proletariat_budget.expenditures;
TRUNCATE TABLE proletariat_budget.household_members;
TRUNCATE TABLE proletariat_budget.ingress_recurrence_patterns;
TRUNCATE TABLE proletariat_budget.ingress_tags;
TRUNCATE TABLE proletariat_budget.ingresses;
TRUNCATE TABLE proletariat_budget.roles;
TRUNCATE TABLE proletariat_budget.savings_contribution_tags;
TRUNCATE TABLE proletariat_budget.savings_contributions;
TRUNCATE TABLE proletariat_budget.savings_goal_tags;
TRUNCATE TABLE proletariat_budget.savings_goals;
TRUNCATE TABLE proletariat_budget.savings_withdrawal_tags;
TRUNCATE TABLE proletariat_budget.savings_withdrawals;
TRUNCATE TABLE proletariat_budget.tags;
TRUNCATE TABLE proletariat_budget.transaction_rollbacks;
TRUNCATE TABLE proletariat_budget.transactions;
TRUNCATE TABLE proletariat_budget.transfers;
TRUNCATE TABLE proletariat_budget.user_roles;
TRUNCATE TABLE proletariat_budget.users;
SET FOREIGN_KEY_CHECKS = 1;
-- Drop the stored procedures
DROP PROCEDURE IF EXISTS CreateExpenditure;
DROP PROCEDURE IF EXISTS CreateIngress;
DROP PROCEDURE IF EXISTS CreateTransfer;
DROP PROCEDURE IF EXISTS RollbackTransaction;


-- Create stored procedure for expenditure creation
DELIMITER $$

CREATE PROCEDURE CreateExpenditure(
    IN p_account_id BIGINT,
    IN p_amount DECIMAL(15, 2),
    IN p_currency INT,
    IN p_transaction_date DATETIME,
    IN p_description TEXT,
    IN p_category_id BIGINT,
    IN p_declared BOOLEAN,
    IN p_planned BOOLEAN,
    IN p_tag_ids TEXT, -- Comma-separated list of tag IDs
    OUT p_expenditure_id BIGINT,
    OUT p_transaction_id BIGINT
)
BEGIN
    DECLARE v_current_balance DECIMAL(15, 2);
    DECLARE v_new_balance DECIMAL(15, 2);
    DECLARE v_tag_id BIGINT;
    DECLARE v_pos INT;
    DECLARE v_remaining_tags TEXT;
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
        BEGIN
            ROLLBACK;
            RESIGNAL;
        END;

    START TRANSACTION;

    -- Get current account balance
    SELECT current_balance
    INTO v_current_balance
    FROM accounts
    WHERE id = p_account_id
      AND active = TRUE
        FOR
    UPDATE;

    -- Check if account exists
    IF v_current_balance IS NULL THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Account not found or inactive';
    END IF;

    -- Calculate new balance (subtract expenditure amount)
    SET v_new_balance = v_current_balance - p_amount;

    -- Create transaction record
    INSERT INTO transactions (account_id,
                              amount,
                              currency,
                              transaction_date,
                              description,
                              transaction_type,
                              balance_after,
                              status)
    VALUES (p_account_id,
            p_amount, -- Negative amount for expenditure
            p_currency,
            p_transaction_date,
            p_description,
            'expenditure',
            v_new_balance,
            'completed');

    SET p_transaction_id = LAST_INSERT_ID();

    -- Create expenditure record
    INSERT INTO expenditures (category_id,
                              declared,
                              planned,
                              transaction_id)
    VALUES (p_category_id,
            p_declared,
            p_planned,
            p_transaction_id);

    SET p_expenditure_id = LAST_INSERT_ID();

    -- Update account balance
    UPDATE accounts
    SET current_balance = v_new_balance,
        updated_at      = NOW()
    WHERE id = p_account_id;

    -- Process tags if provided
    IF p_tag_ids IS NOT NULL AND LENGTH(TRIM(p_tag_ids)) > 0 THEN
        SET v_remaining_tags = CONCAT(p_tag_ids, ',');

        WHILE LENGTH(v_remaining_tags) > 0
            DO
                SET v_pos = LOCATE(',', v_remaining_tags);
                IF v_pos > 0 THEN
                    SET v_tag_id = CAST(TRIM(SUBSTRING(v_remaining_tags, 1, v_pos - 1)) AS UNSIGNED);
                    SET v_remaining_tags = SUBSTRING(v_remaining_tags, v_pos + 1);

                    -- Insert tag relationship if tag_id is valid
                    IF v_tag_id > 0 THEN
                        INSERT INTO expenditure_tags (expenditure_id,
                                                      tag_id)
                        VALUES (p_expenditure_id,
                                v_tag_id);
                    END IF;
                ELSE
                    SET v_remaining_tags = '';
                END IF;
            END WHILE;
    END IF;
    COMMIT;
END$$

DELIMITER ;


DELIMITER $$

CREATE PROCEDURE CreateIngress(
    IN p_account_id BIGINT,
    IN p_amount DECIMAL(15, 2),
    IN p_currency INT,
    IN p_transaction_date DATETIME,
    IN p_description TEXT,
    IN p_category_id BIGINT,
    IN p_source VARCHAR(255),
    IN p_is_recurring BOOLEAN,
    IN p_tag_ids JSON -- JSON array of tag IDs, e.g., '[1, 2, 3]'
)
BEGIN
    DECLARE v_transaction_id BIGINT;
    DECLARE v_ingress_id BIGINT;
    DECLARE v_current_balance DECIMAL(15, 2);
    DECLARE v_new_balance DECIMAL(15, 2);
    DECLARE v_tag_count INT;
    DECLARE v_counter INT DEFAULT 0;
    DECLARE v_tag_id BIGINT;
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
        BEGIN
            ROLLBACK;
            RESIGNAL;
        END;

    START TRANSACTION;

    -- Get current account balance
    SELECT current_balance
    INTO v_current_balance
    FROM accounts
    WHERE id = p_account_id
      AND active = TRUE;

    IF v_current_balance IS NULL THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Account not found or inactive';
    END IF;

    -- Calculate new balance
    SET v_new_balance = v_current_balance + p_amount;

    -- Create transaction record
    INSERT INTO transactions (account_id,
                              amount,
                              currency,
                              transaction_date,
                              description,
                              transaction_type,
                              balance_after,
                              status)
    VALUES (p_account_id,
            p_amount,
            p_currency,
            p_transaction_date,
            p_description,
            'ingress',
            v_new_balance,
            'completed');

    SET v_transaction_id = LAST_INSERT_ID();

    -- Update account balance
    UPDATE accounts
    SET current_balance = v_new_balance,
        updated_at      = CURRENT_TIMESTAMP
    WHERE id = p_account_id;

    -- Create ingress record
    INSERT INTO ingresses (category_id,
                           source,
                           is_recurring,
                           transaction_id)
    VALUES (p_category_id,
            p_source,
            p_is_recurring,
            v_transaction_id);

    SET v_ingress_id = LAST_INSERT_ID();

    -- Handle tags if provided
    IF p_tag_ids IS NOT NULL AND JSON_LENGTH(p_tag_ids) > 0 THEN
        SET v_tag_count = JSON_LENGTH(p_tag_ids);

        WHILE v_counter < v_tag_count
            DO
                SET v_tag_id = JSON_EXTRACT(p_tag_ids, CONCAT('$[', v_counter, ']'));

                -- Insert tag relationship
                INSERT INTO ingress_tags (ingress_id,
                                          tag_id)
                VALUES (v_ingress_id,
                        v_tag_id);

                SET v_counter = v_counter + 1;
            END WHILE;
    END IF;
    COMMIT;

    -- Return the created IDs for reference
    SELECT v_transaction_id as transaction_id, v_ingress_id as ingress_id, v_new_balance as new_balance;

END$$

DELIMITER ;

DELIMITER $$

CREATE PROCEDURE CreateTransfer(
    IN p_source_account_id BIGINT,
    IN p_destination_account_id BIGINT,
    IN p_source_amount DECIMAL(15, 2),
    IN p_destination_amount DECIMAL(15, 2),
    IN p_exchange_rate_multiplier DECIMAL(15, 6),
    IN p_fees DECIMAL(15, 2),
    IN p_transaction_date DATETIME,
    IN p_description TEXT
)
BEGIN
    DECLARE v_outgoing_transaction_id BIGINT;
    DECLARE v_incoming_transaction_id BIGINT;
    DECLARE v_transfer_id BIGINT;
    DECLARE v_source_current_balance DECIMAL(15, 2);
    DECLARE v_dest_current_balance DECIMAL(15, 2);
    DECLARE v_source_new_balance DECIMAL(15, 2);
    DECLARE v_dest_new_balance DECIMAL(15, 2);
    DECLARE v_source_currency INT;
    DECLARE v_dest_currency INT;
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
        BEGIN
            ROLLBACK;
            RESIGNAL;
        END;

    START TRANSACTION;

    -- Get source account details
    SELECT current_balance, currency
    INTO v_source_current_balance, v_source_currency
    FROM accounts
    WHERE id = p_source_account_id
      AND active = TRUE;

    IF v_source_current_balance IS NULL THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Source account not found or inactive';
    END IF;

    -- Get destination account details
    SELECT current_balance, currency
    INTO v_dest_current_balance, v_dest_currency
    FROM accounts
    WHERE id = p_destination_account_id
      AND active = TRUE;

    IF v_dest_current_balance IS NULL THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Destination account not found or inactive';
    END IF;

    -- Validate source account has sufficient funds
    IF v_source_current_balance < p_source_amount THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Insufficient funds in source account';
    END IF;

    -- Validate amounts are positive
    IF p_source_amount <= 0 OR p_destination_amount <= 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Transfer amounts must be positive';
    END IF;

    -- Calculate new balances
    SET v_source_new_balance = v_source_current_balance - p_source_amount;
    SET v_dest_new_balance = v_dest_current_balance + p_destination_amount;

    -- Create outgoing transaction (debit from source account)
    INSERT INTO transactions (account_id,
                              amount,
                              currency,
                              transaction_date,
                              description,
                              transaction_type,
                              balance_after,
                              status)
    VALUES (p_source_account_id,
            p_source_amount,
            v_source_currency,
            p_transaction_date,
            CONCAT('Transfer to account ', p_destination_account_id,
                   CASE WHEN p_description IS NOT NULL THEN CONCAT(' - ', p_description) ELSE '' END),
            'transfer',
            v_source_new_balance,
            'completed');

    SET v_outgoing_transaction_id = LAST_INSERT_ID();

    -- Create incoming transaction (credit to destination account)
    INSERT INTO transactions (account_id,
                              amount,
                              currency,
                              transaction_date,
                              description,
                              transaction_type,
                              balance_after,
                              status)
    VALUES (p_destination_account_id,
            p_destination_amount,
            v_dest_currency,
            p_transaction_date,
            CONCAT('Transfer from account ', p_source_account_id,
                   CASE WHEN p_description IS NOT NULL THEN CONCAT(' - ', p_description) ELSE '' END),
            'transfer',
            v_dest_new_balance,
            'completed');

    SET v_incoming_transaction_id = LAST_INSERT_ID();

    -- Update source account balance
    UPDATE accounts
    SET current_balance = v_source_new_balance,
        updated_at      = CURRENT_TIMESTAMP
    WHERE id = p_source_account_id;

    -- Update destination account balance
    UPDATE accounts
    SET current_balance = v_dest_new_balance,
        updated_at      = CURRENT_TIMESTAMP
    WHERE id = p_destination_account_id;

    -- Create transfer record
    INSERT INTO transfers (source_account_id,
                           destination_account_id,
                           destination_amount,
                           exchange_rate_multiplier,
                           fees,
                           outgoing_transaction_id,
                           incoming_transaction_id)
    VALUES (p_source_account_id,
            p_destination_account_id,
            p_destination_amount,
            p_exchange_rate_multiplier,
            p_fees,
            v_outgoing_transaction_id,
            v_incoming_transaction_id);

    SET v_transfer_id = LAST_INSERT_ID();
    COMMIT;

    -- Return the created IDs and new balances for reference
    SELECT v_transfer_id             as transfer_id,
           v_outgoing_transaction_id as outgoing_transaction_id,
           v_incoming_transaction_id as incoming_transaction_id,
           v_source_new_balance      as source_new_balance,
           v_dest_new_balance        as destination_new_balance;

END$$

DELIMITER ;


DELIMITER $$

CREATE PROCEDURE RollbackTransaction(
    IN p_transaction_id BIGINT,
    IN p_rollback_reason TEXT
)
BEGIN
    DECLARE v_transaction_type ENUM ('expenditure', 'ingress', 'transfer', 'rollback');
    DECLARE v_account_id BIGINT;
    DECLARE v_amount DECIMAL(15, 2);
    DECLARE v_currency INT;
    DECLARE v_transaction_date DATETIME;
    DECLARE v_description TEXT;
    DECLARE v_current_balance DECIMAL(15, 2);
    DECLARE v_new_balance DECIMAL(15, 2);
    DECLARE v_rollback_transaction_id BIGINT;

    -- Transfer specific variables
    DECLARE v_transfer_id BIGINT;
    DECLARE v_source_account_id BIGINT;
    DECLARE v_dest_account_id BIGINT;
    DECLARE v_outgoing_transaction_id BIGINT;
    DECLARE v_incoming_transaction_id BIGINT;
    DECLARE v_source_amount DECIMAL(15, 2);
    DECLARE v_dest_amount DECIMAL(15, 2);
    DECLARE v_source_currency INT;
    DECLARE v_dest_currency INT;
    DECLARE v_source_current_balance DECIMAL(15, 2);
    DECLARE v_dest_current_balance DECIMAL(15, 2);
    DECLARE v_source_new_balance DECIMAL(15, 2);
    DECLARE v_dest_new_balance DECIMAL(15, 2);
    DECLARE v_source_rollback_id BIGINT;
    DECLARE v_dest_rollback_id BIGINT;

    DECLARE EXIT HANDLER FOR SQLEXCEPTION
        BEGIN
            ROLLBACK;
            RESIGNAL;
        END;

    START TRANSACTION;

    -- Get transaction details
    SELECT transaction_type, account_id, amount, currency, transaction_date, description
    INTO v_transaction_type, v_account_id, v_amount, v_currency, v_transaction_date, v_description
    FROM transactions
    WHERE id = p_transaction_id
      AND status = 'completed';

    IF v_transaction_type IS NULL THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Transaction not found or not completed';
    END IF;

    -- Exclude rollback type transactions
    IF v_transaction_type = 'rollback' THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Cannot rollback a rollback transaction';
    END IF;

    -- Check if transaction is already rolled back
    IF EXISTS (SELECT 1 FROM transaction_rollbacks WHERE transaction_id = p_transaction_id) THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Transaction has already been rolled back';
    END IF;

    -- Handle different transaction types
    CASE v_transaction_type
        WHEN 'expenditure' THEN -- Get current account balance
        SELECT current_balance
        INTO v_current_balance
        FROM accounts
        WHERE id = v_account_id
          AND active = TRUE;

        IF v_current_balance IS NULL THEN
            SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Account not found or inactive';
        END IF;

        -- For expenditure rollback: add money back to account (negative expenditure amount)
        SET v_new_balance = v_current_balance + v_amount;

        -- Create rollback transaction with negative amount
        INSERT INTO transactions (account_id,
                                  amount,
                                  currency,
                                  transaction_date,
                                  description,
                                  transaction_type,
                                  balance_after,
                                  status)
        VALUES (v_account_id,
                -v_amount,
                v_currency,
                CURRENT_TIMESTAMP,
                CONCAT('Rollback of expenditure transaction #', p_transaction_id,
                       CASE
                           WHEN p_rollback_reason IS NOT NULL THEN CONCAT(' - Reason: ', p_rollback_reason)
                           ELSE '' END),
                'rollback',
                v_new_balance,
                'completed');

        SET v_rollback_transaction_id = LAST_INSERT_ID();

        -- Update account balance
        UPDATE accounts
        SET current_balance = v_new_balance,
            updated_at      = CURRENT_TIMESTAMP
        WHERE id = v_account_id;

        -- Record rollback relationship
        INSERT INTO transaction_rollbacks (transaction_id,
                                           rollback_transaction_id,
                                           rollback_reason)
        VALUES (p_transaction_id,
                v_rollback_transaction_id,
                p_rollback_reason);

        WHEN 'ingress' THEN -- Get current account balance
        SELECT current_balance
        INTO v_current_balance
        FROM accounts
        WHERE id = v_account_id
          AND active = TRUE;

        IF v_current_balance IS NULL THEN
            SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Account not found or inactive';
        END IF;

        -- Check if account has sufficient funds for rollback
        IF v_current_balance < v_amount THEN
            SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Insufficient funds to rollback ingress transaction';
        END IF;

        -- For ingress rollback: remove money from account (negative ingress amount)
        SET v_new_balance = v_current_balance - v_amount;

        -- Create rollback transaction with negative amount
        INSERT INTO transactions (account_id,
                                  amount,
                                  currency,
                                  transaction_date,
                                  description,
                                  transaction_type,
                                  balance_after,
                                  status)
        VALUES (v_account_id,
                -v_amount,
                v_currency,
                CURRENT_TIMESTAMP,
                CONCAT('Rollback of ingress transaction #', p_transaction_id,
                       CASE
                           WHEN p_rollback_reason IS NOT NULL THEN CONCAT(' - Reason: ', p_rollback_reason)
                           ELSE '' END),
                'rollback',
                v_new_balance,
                'completed');

        SET v_rollback_transaction_id = LAST_INSERT_ID();

        -- Update account balance
        UPDATE accounts
        SET current_balance = v_new_balance,
            updated_at      = CURRENT_TIMESTAMP
        WHERE id = v_account_id;

        -- Record rollback relationship
        INSERT INTO transaction_rollbacks (transaction_id,
                                           rollback_transaction_id,
                                           rollback_reason)
        VALUES (p_transaction_id,
                v_rollback_transaction_id,
                p_rollback_reason);

        WHEN 'transfer' THEN -- Get transfer details
        SELECT id,
               source_account_id,
               destination_account_id,
               outgoing_transaction_id,
               incoming_transaction_id
        INTO
            v_transfer_id, v_source_account_id, v_dest_account_id,
            v_outgoing_transaction_id, v_incoming_transaction_id
        FROM transfers
        WHERE outgoing_transaction_id = p_transaction_id
           OR incoming_transaction_id = p_transaction_id;

        IF v_transfer_id IS NULL THEN
            SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Transfer record not found for this transaction';
        END IF;

        -- Get outgoing transaction details
        SELECT account_id, amount, currency
        INTO v_source_account_id, v_source_amount, v_source_currency
        FROM transactions
        WHERE id = v_outgoing_transaction_id;

        -- Get incoming transaction details
        SELECT account_id, amount, currency
        INTO v_dest_account_id, v_dest_amount, v_dest_currency
        FROM transactions
        WHERE id = v_incoming_transaction_id;

        -- Get current balances
        SELECT current_balance
        INTO v_source_current_balance
        FROM accounts
        WHERE id = v_source_account_id
          AND active = TRUE;

        SELECT current_balance
        INTO v_dest_current_balance
        FROM accounts
        WHERE id = v_dest_account_id
          AND active = TRUE;

        IF v_source_current_balance IS NULL OR v_dest_current_balance IS NULL THEN
            SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'One or both accounts not found or inactive';
        END IF;

        -- Check if destination account has sufficient funds for rollback
        IF v_dest_current_balance < v_dest_amount THEN
            SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Insufficient funds in destination account to rollback transfer';
        END IF;

        -- Calculate new balances (reverse the transfer)
        SET v_source_new_balance = v_source_current_balance + v_source_amount; -- Add back to source
        SET v_dest_new_balance = v_dest_current_balance - v_dest_amount;
        -- Remove from destination

        -- Create rollback transaction for source account (negative outgoing amount)
        INSERT INTO transactions (account_id,
                                  amount,
                                  currency,
                                  transaction_date,
                                  description,
                                  transaction_type,
                                  balance_after,
                                  status)
        VALUES (v_source_account_id,
                -v_source_amount,
                v_source_currency,
                CURRENT_TIMESTAMP,
                CONCAT('Rollback of transfer (outgoing) #', v_outgoing_transaction_id,
                       CASE
                           WHEN p_rollback_reason IS NOT NULL THEN CONCAT(' - Reason: ', p_rollback_reason)
                           ELSE '' END),
                'rollback',
                v_source_new_balance,
                'completed');

        SET v_source_rollback_id = LAST_INSERT_ID();

        -- Create rollback transaction for destination account (negative incoming amount)
        INSERT INTO transactions (account_id,
                                  amount,
                                  currency,
                                  transaction_date,
                                  description,
                                  transaction_type,
                                  balance_after,
                                  status)
        VALUES (v_dest_account_id,
                -v_dest_amount,
                v_dest_currency,
                CURRENT_TIMESTAMP,
                CONCAT('Rollback of transfer (incoming) #', v_incoming_transaction_id,
                       CASE
                           WHEN p_rollback_reason IS NOT NULL THEN CONCAT(' - Reason: ', p_rollback_reason)
                           ELSE '' END),
                'rollback',
                v_dest_new_balance,
                'completed');

        SET v_dest_rollback_id = LAST_INSERT_ID();

        -- Update account balances
        UPDATE accounts
        SET current_balance = v_source_new_balance,
            updated_at      = CURRENT_TIMESTAMP
        WHERE id = v_source_account_id;

        UPDATE accounts
        SET current_balance = v_dest_new_balance,
            updated_at      = CURRENT_TIMESTAMP
        WHERE id = v_dest_account_id;

        -- Record rollback relationships for both transactions
        INSERT INTO transaction_rollbacks (transaction_id,
                                           rollback_transaction_id,
                                           rollback_reason)
        VALUES (v_outgoing_transaction_id,
                v_source_rollback_id,
                p_rollback_reason),
               (v_incoming_transaction_id,
                v_dest_rollback_id,
                p_rollback_reason);

        ELSE SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Unsupported transaction type for rollback';
        END CASE;
    COMMIT;

    -- Return success information
    SELECT CASE v_transaction_type
               WHEN 'transfer' THEN CONCAT('Transfer rolled back successfully. Source rollback ID: ',
                                           v_source_rollback_id, ', Destination rollback ID: ', v_dest_rollback_id)
               ELSE CONCAT('Transaction rolled back successfully. Rollback transaction ID: ', v_rollback_transaction_id)
               END as result_message;

END$$

DELIMITER ;


-- Insert household members
INSERT INTO household_members (name,
                               surname,
                               nickname,
                               role,
                               active)
VALUES ('John',
        'Smith',
        'Johnny',
        'head_of_household',
        TRUE),
       ('Jane',
        'Smith',
        'Janie',
        'spouse',
        TRUE),
       ('Alex',
        'Smith',
        'Al',
        'child',
        TRUE);

-- Insert accounts
INSERT INTO accounts (name,
                      type,
                      owner,
                      institution,
                      currency,
                      initial_balance,
                      current_balance,
                      active,
                      description,
                      account_number)
VALUES
-- Cash accounts
('Cash ARS',
 'cash',
 1,
 NULL,
 (SELECT id FROM currencies WHERE symbol = 'ARS'),
 500000.00,
 450000.00,
 TRUE,
 'Argentine Peso cash wallet',
 NULL),
('Cash USD',
 'cash',
 1,
 NULL,
 (SELECT id FROM currencies WHERE symbol = 'USD'),
 5000.00,
 4500.00,
 TRUE,
 'US Dollar cash wallet',
 NULL),
-- Bank accounts
('John Main Account',
 'bank',
 1,
 'Banco Nación',
 (SELECT id FROM currencies WHERE symbol = 'ARS'),
 1500000.00,
 1250000.00,
 TRUE,
 'Primary checking account',
 '1234567890'),
('Jane Savings Account',
 'bank',
 2,
 'Banco Santander',
 (SELECT id FROM currencies WHERE symbol = 'ARS'),
 800000.00,
 850000.00,
 TRUE,
 'Personal savings account',
 '0987654321'),
('Alex Student Account',
 'bank',
 3,
 'Banco Galicia',
 (SELECT id FROM currencies WHERE symbol = 'ARS'),
 250000.00,
 220000.00,
 TRUE,
 'Student checking account',
 '1122334455');

-- Insert tags for ingresses
INSERT INTO tags (name,
                  description,
                  color,
                  background_color,
                  type)
VALUES ('Salary',
        'Regular employment income',
        '#FFFFFF',
        '#4CAF50',
        'income'),
       ('Bonus',
        'Performance or holiday bonuses',
        '#FFFFFF',
        '#FF9800',
        'income'),
       ('Freelance',
        'Independent contractor work',
        '#FFFFFF',
        '#2196F3',
        'income'),
       ('Investment',
        'Returns from investments',
        '#FFFFFF',
        '#9C27B0',
        'income'),
       ('Gift',
        'Money received as gifts',
        '#FFFFFF',
        '#E91E63',
        'income'),
       ('Refund',
        'Tax refunds or purchase returns',
        '#FFFFFF',
        '#00BCD4',
        'income'),
       ('Side Business',
        'Income from side business activities',
        '#FFFFFF',
        '#795548',
        'income'),
       ('Rental',
        'Income from property rentals',
        '#FFFFFF',
        '#607D8B',
        'income');

-- Insert categories for ingresses
INSERT INTO categories (name,
                        description,
                        color,
                        background_color,
                        active,
                        category_type)
VALUES ('Employment Income',
        'Regular salary and wages',
        '#FFFFFF',
        '#4CAF50',
        TRUE,
        'income'),
       ('Business Income',
        'Income from business activities',
        '#FFFFFF',
        '#FF9800',
        TRUE,
        'income'),
       ('Investment Income',
        'Returns from investments and dividends',
        '#FFFFFF',
        '#9C27B0',
        TRUE,
        'income'),
       ('Other Income',
        'Miscellaneous income sources',
        '#FFFFFF',
        '#607D8B',
        TRUE,
        'income'),
       ('Government Benefits',
        'Social security and government assistance',
        '#FFFFFF',
        '#2196F3',
        TRUE,
        'income');

-- Insert tags for expenditures
INSERT INTO tags (name,
                  description,
                  color,
                  background_color,
                  type)
VALUES ('Groceries',
        'Food and household items',
        '#FFFFFF',
        '#FF5722',
        'expense'),
       ('Utilities',
        'Electricity, gas, water, internet',
        '#FFFFFF',
        '#FFC107',
        'expense'),
       ('Transportation',
        'Fuel, public transport, maintenance',
        '#FFFFFF',
        '#2196F3',
        'expense'),
       ('Healthcare',
        'Medical expenses and insurance',
        '#FFFFFF',
        '#E91E63',
        'expense'),
       ('Entertainment',
        'Movies, games, dining out',
        '#FFFFFF',
        '#9C27B0',
        'expense'),
       ('Education',
        'Books, courses, school fees',
        '#FFFFFF',
        '#00BCD4',
        'expense'),
       ('Clothing',
        'Apparel and accessories',
        '#FFFFFF',
        '#795548',
        'expense'),
       ('Home Maintenance',
        'Repairs and home improvements',
        '#FFFFFF',
        '#607D8B',
        'expense'),
       ('Insurance',
        'Life, health, auto insurance',
        '#FFFFFF',
        '#FF9800',
        'expense'),
       ('Taxes',
        'Income tax, property tax',
        '#FFFFFF',
        '#F44336',
        'expense'),
       ('Debt Payment',
        'Credit card, loan payments',
        '#FFFFFF',
        '#9E9E9E',
        'expense'),
       ('Emergency',
        'Unexpected urgent expenses',
        '#FFFFFF',
        '#FF1744',
        'expense'),
       ('Subscription',
        'Monthly/yearly service subscriptions',
        '#FFFFFF',
        '#3F51B5',
        'expense'),
       ('Personal Care',
        'Hygiene, beauty, wellness',
        '#FFFFFF',
        '#4CAF50',
        'expense'),
       ('Technology',
        'Electronics, software, gadgets',
        '#FFFFFF',
        '#00E676',
        'expense');

-- Insert categories for expenditures
INSERT INTO categories (name,
                        description,
                        color,
                        background_color,
                        active,
                        category_type)
VALUES ('Food & Dining',
        'Groceries, restaurants, food delivery',
        '#FFFFFF',
        '#FF5722',
        TRUE,
        'expense'),
       ('Housing',
        'Rent, mortgage, utilities, maintenance',
        '#FFFFFF',
        '#795548',
        TRUE,
        'expense'),
       ('Transportation & Travel',
        'Vehicle expenses, public transport, travel',
        '#FFFFFF',
        '#2196F3',
        TRUE,
        'expense'),
       ('Health & Medical',
        'Healthcare, medications, insurance',
        '#FFFFFF',
        '#E91E63',
        TRUE,
        'expense'),
       ('Entertainment & Recreation',
        'Movies, hobbies, sports, dining out',
        '#FFFFFF',
        '#9C27B0',
        TRUE,
        'expense'),
       ('Education & Learning',
        'School fees, books, courses, training',
        '#FFFFFF',
        '#00BCD4',
        TRUE,
        'expense'),
       ('Personal & Family',
        'Clothing, personal care, family needs',
        '#FFFFFF',
        '#4CAF50',
        TRUE,
        'expense'),
       ('Financial Services',
        'Banking fees, insurance, taxes, debt payments',
        '#FFFFFF',
        '#FF9800',
        TRUE,
        'expense'),
       ('Technology & Communication',
        'Phone, internet, software, electronics',
        '#FFFFFF',
        '#3F51B5',
        TRUE,
        'expense'),
       ('Miscellaneous',
        'Other expenses not categorized elsewhere',
        '#FFFFFF',
        '#607D8B',
        TRUE,
        'expense');


-- Insert tags for savings goals
INSERT INTO tags (name,
                  description,
                  color,
                  background_color,
                  type)
VALUES ('Emergency Fund',
        'Emergency financial reserves',
        '#FFFFFF',
        '#FF5722',
        'savings'),
       ('Vacation',
        'Travel and holiday savings',
        '#FFFFFF',
        '#4CAF50',
        'savings'),
       ('Home Purchase',
        'Saving for buying a house',
        '#FFFFFF',
        '#2196F3',
        'savings'),
       ('Car Purchase',
        'Saving for vehicle acquisition',
        '#FFFFFF',
        '#FF9800',
        'savings'),
       ('Retirement',
        'Long-term retirement planning',
        '#FFFFFF',
        '#9C27B0',
        'savings'),
       ('Education Fund',
        'Educational expenses savings',
        '#FFFFFF',
        '#00BCD4',
        'savings'),
       ('Wedding',
        'Wedding ceremony and reception',
        '#FFFFFF',
        '#E91E63',
        'savings'),
       ('Home Improvement',
        'House renovation and upgrades',
        '#FFFFFF',
        '#795548',
        'savings'),
       ('Medical Fund',
        'Healthcare and medical expenses',
        '#FFFFFF',
        '#607D8B',
        'savings'),
       ('Technology Upgrade',
        'Electronics and gadgets',
        '#FFFFFF',
        '#3F51B5',
        'savings'),
       ('Short Term',
        'Goals achievable within 1 year',
        '#FFFFFF',
        '#FFC107',
        'savings'),
       ('Long Term',
        'Goals requiring multiple years',
        '#FFFFFF',
        '#8BC34A',
        'savings'),
       ('High Priority',
        'Most important financial goals',
        '#FFFFFF',
        '#F44336',
        'savings'),
       ('Low Priority',
        'Nice-to-have financial goals',
        '#FFFFFF',
        '#9E9E9E',
        'savings'),
       ('Family Goal',
        'Shared family objectives',
        '#FFFFFF',
        '#CDDC39',
        'savings');

-- Insert categories for savings goals
INSERT INTO categories (name,
                        description,
                        color,
                        background_color,
                        active,
                        category_type)
VALUES ('Emergency & Security',
        'Emergency funds and financial security',
        '#FFFFFF',
        '#FF5722',
        TRUE,
        'savings'),
       ('Travel & Leisure',
        'Vacation and recreational activities',
        '#FFFFFF',
        '#4CAF50',
        TRUE,
        'savings'),
       ('Real Estate',
        'Property purchase and improvements',
        '#FFFFFF',
        '#2196F3',
        TRUE,
        'savings'),
       ('Vehicle & Transportation',
        'Car purchase and transportation needs',
        '#FFFFFF',
        '#FF9800',
        TRUE,
        'savings'),
       ('Retirement & Future',
        'Long-term financial planning',
        '#FFFFFF',
        '#9C27B0',
        TRUE,
        'savings'),
       ('Education & Development',
        'Learning and skill development',
        '#FFFFFF',
        '#00BCD4',
        TRUE,
        'savings'),
       ('Life Events',
        'Weddings, celebrations, major events',
        '#FFFFFF',
        '#E91E63',
        TRUE,
        'savings'),
       ('Health & Wellness',
        'Medical and healthcare expenses',
        '#FFFFFF',
        '#607D8B',
        TRUE,
        'savings'),
       ('Technology & Electronics',
        'Gadgets and tech upgrades',
        '#FFFFFF',
        '#3F51B5',
        TRUE,
        'savings'),
       ('General Savings',
        'Miscellaneous savings goals',
        '#FFFFFF',
        '#795548',
        TRUE,
        'savings');


-- Insert savings goals
INSERT INTO savings_goals (name,
                           category_id,
                           description,
                           target_amount,
                           currency,
                           target_date,
                           initial_amount,
                           current_amount,
                           percent_complete,
                           account_id,
                           priority,
                           auto_contribute,
                           auto_contribute_amount,
                           auto_contribute_frequency,
                           status,
                           projected_completion_date)
VALUES
-- Emergency fund
('Emergency Fund 6 Months',
 (SELECT id FROM categories WHERE name = 'Emergency & Security'),
 'Emergency fund covering 6 months of expenses',
 600000.00,
 (SELECT id FROM currencies WHERE symbol = 'ARS'),
 '2024-12-31',
 50000.00,
 125000.00,
 20.83,
 4,
 1,
 TRUE,
 15000.00,
 'monthly',
 'active',
 '2024-11-30'),

-- Family vacation
('Europe Family Trip 2025',
 (SELECT id FROM categories WHERE name = 'Travel & Leisure'),
 'Summer vacation to Europe for the whole family',
 2500.00,
 (SELECT id FROM currencies WHERE symbol = 'USD'),
 '2025-06-01',
 200.00,
 450.00,
 18.00,
 2,
 2,
 TRUE,
 150.00,
 'monthly',
 'active',
 '2025-05-15'),

-- Home down payment
('House Down Payment',
 (SELECT id FROM categories WHERE name = 'Real Estate'),
 'Down payment for family home purchase',
 1500000.00,
 (SELECT id FROM currencies WHERE symbol = 'ARS'),
 '2026-03-01',
 100000.00,
 280000.00,
 18.67,
 3,
 1,
 TRUE,
 50000.00,
 'monthly',
 'active',
 '2025-12-01'),

-- Car replacement
('New Family Car',
 (SELECT id FROM categories WHERE name = 'Vehicle & Transportation'),
 'Replace current car with newer model',
 800000.00,
 (SELECT id FROM currencies WHERE symbol = 'ARS'),
 '2024-08-01',
 75000.00,
 150000.00,
 18.75,
 4,
 3,
 FALSE,
 NULL,
 NULL,
 'active',
 '2024-10-01'),

-- Alex education fund
('Alex University Fund',
 (SELECT id FROM categories WHERE name = 'Education & Development'),
 'University education expenses for Alex',
 1200000.00,
 (SELECT id FROM currencies WHERE symbol = 'ARS'),
 '2028-03-01',
 25000.00,
 85000.00,
 7.08,
 5,
 2,
 TRUE,
 8000.00,
 'monthly',
 'active',
 '2027-08-01'),

-- Home renovation
('Kitchen Renovation',
 (SELECT id FROM categories WHERE name = 'Real Estate'),
 'Complete kitchen remodel and upgrade',
 400000.00,
 (SELECT id FROM currencies WHERE symbol = 'ARS'),
 '2024-10-01',
 30000.00,
 95000.00,
 23.75,
 3,
 4,
 FALSE,
 NULL,
 NULL,
 'active',
 '2024-12-01'),

-- Technology upgrade
('Home Tech Upgrade',
 (SELECT id FROM categories WHERE name = 'Technology & Electronics'),
 'New laptops and smart home devices',
 300000.00,
 (SELECT id FROM currencies WHERE symbol = 'ARS'),
 '2024-07-01',
 20000.00,
 65000.00,
 21.67,
 4,
 5,
 TRUE,
 12000.00,
 'monthly',
 'active',
 '2024-09-01'),

-- Wedding fund (completed goal)
('Jane Sister Wedding Gift',
 (SELECT id FROM categories WHERE name = 'Life Events'),
 'Special wedding gift for Jane sister',
 50000.00,
 (SELECT id FROM currencies WHERE symbol = 'ARS'),
 '2024-01-15',
 10000.00,
 50000.00,
 100.00,
 4,
 2,
 FALSE,
 NULL,
 NULL,
 'completed',
 '2024-01-10'),

-- Medical fund
('Family Health Fund',
 (SELECT id FROM categories WHERE name = 'Health & Wellness'),
 'Medical emergencies and health expenses',
 200000.00,
 (SELECT id FROM currencies WHERE symbol = 'ARS'),
 '2024-12-31',
 15000.00,
 45000.00,
 22.50,
 3,
 3,
 TRUE,
 8000.00,
 'monthly',
 'active',
 '2024-11-01'),

-- Retirement savings
('John Retirement Fund',
 (SELECT id FROM categories WHERE name = 'Retirement & Future'),
 'Long-term retirement savings for John',
 5000000.00,
 (SELECT id FROM currencies WHERE symbol = 'ARS'),
 '2040-01-01',
 200000.00,
 350000.00,
 7.00,
 3,
 1,
 TRUE,
 25000.00,
 'monthly',
 'active',
 '2038-06-01');

-- Insert savings goal tags relationships
INSERT INTO savings_goal_tags (savings_goal_id,
                               tag_id)
VALUES
-- Emergency fund
(1,
 (SELECT id FROM tags WHERE name = 'Emergency Fund')),
(1,
 (SELECT id FROM tags WHERE name = 'High Priority')),
(1,
 (SELECT id FROM tags WHERE name = 'Family Goal')),

-- Europe trip
(2,
 (SELECT id FROM tags WHERE name = 'Vacation')),
(2,
 (SELECT id FROM tags WHERE name = 'Family Goal')),
(2,
 (SELECT id FROM tags WHERE name = 'Long Term')),

-- House down payment
(3,
 (SELECT id FROM tags WHERE name = 'Home Purchase')),
(3,
 (SELECT id FROM tags WHERE name = 'High Priority')),
(3,
 (SELECT id FROM tags WHERE name = 'Long Term')),
(3,
 (SELECT id FROM tags WHERE name = 'Family Goal')),

-- New car
(4,
 (SELECT id FROM tags WHERE name = 'Car Purchase')),
(4,
 (SELECT id FROM tags WHERE name = 'Short Term')),
(4,
 (SELECT id FROM tags WHERE name = 'Family Goal')),

-- Alex education
(5,
 (SELECT id FROM tags WHERE name = 'Education Fund')),
(5,
 (SELECT id FROM tags WHERE name = 'Long Term')),
(5,
 (SELECT id FROM tags WHERE name = 'High Priority')),
(5,
 (SELECT id FROM tags WHERE name = 'Family Goal')),

-- Kitchen renovation
(6,
 (SELECT id FROM tags WHERE name = 'Home Improvement')),
(6,
 (SELECT id FROM tags WHERE name = 'Short Term')),
(6,
 (SELECT id FROM tags WHERE name = 'Low Priority')),

-- Tech upgrade
(7,
 (SELECT id FROM tags WHERE name = 'Technology Upgrade')),
(7,
 (SELECT id FROM tags WHERE name = 'Short Term')),
(7,
 (SELECT id FROM tags WHERE name = 'Low Priority')),

-- Wedding gift (completed)
(8,
 (SELECT id FROM tags WHERE name = 'Wedding')),
(8,
 (SELECT id FROM tags WHERE name = 'Family Goal')),
(8,
 (SELECT id FROM tags WHERE name = 'Short Term')),

-- Medical fund
(9,
 (SELECT id FROM tags WHERE name = 'Medical Fund')),
(9,
 (SELECT id FROM tags WHERE name = 'Emergency Fund')),
(9,
 (SELECT id FROM tags WHERE name = 'High Priority')),
(9,
 (SELECT id FROM tags WHERE name = 'Family Goal')),

-- Retirement
(10,
 (SELECT id FROM tags WHERE name = 'Retirement')),
(10,
 (SELECT id FROM tags WHERE name = 'Long Term')),
(10,
 (SELECT id FROM tags WHERE name = 'High Priority'));


-- Insert tags for savings contributions
INSERT INTO tags (name,
                  description,
                  color,
                  background_color,
                  type)
VALUES ('Monthly Contribution',
        'Regular monthly savings deposits',
        '#FFFFFF',
        '#4CAF50',
        'savings_contribution'),
       ('Bonus Allocation',
        'Savings from work bonuses',
        '#FFFFFF',
        '#FF9800',
        'savings_contribution'),
       ('Tax Refund Savings',
        'Savings from tax refunds',
        '#FFFFFF',
        '#2196F3',
        'savings_contribution'),
       ('Windfall',
        'Unexpected money saved',
        '#FFFFFF',
        '#9C27B0',
        'savings_contribution'),
       ('Salary Increase',
        'Additional savings from pay raise',
        '#FFFFFF',
        '#00BCD4',
        'savings_contribution'),
       ('Side Income',
        'Savings from additional income sources',
        '#FFFFFF',
        '#795548',
        'savings_contribution'),
       ('Expense Reduction',
        'Money saved from cutting expenses',
        '#FFFFFF',
        '#607D8B',
        'savings_contribution'),
       ('Goal Acceleration',
        'Extra contributions to reach goals faster',
        '#FFFFFF',
        '#E91E63',
        'savings_contribution'),
       ('Automatic Transfer',
        'Automated savings contributions',
        '#FFFFFF',
        '#3F51B5',
        'savings_contribution'),
       ('Manual Deposit',
        'Manually made savings deposits',
        '#FFFFFF',
        '#8BC34A',
        'savings_contribution');

-- Insert tags for savings withdrawals
INSERT INTO tags (name,
                  description,
                  color,
                  background_color,
                  type)
VALUES ('Goal Achievement',
        'Withdrawal upon reaching savings goal',
        '#FFFFFF',
        '#4CAF50',
        'savings_withdrawal'),
       ('Emergency Use',
        'Withdrawal for emergency situations',
        '#FFFFFF',
        '#F44336',
        'savings_withdrawal'),
       ('Partial Use',
        'Partial withdrawal for specific need',
        '#FFFFFF',
        '#FF9800',
        'savings_withdrawal'),
       ('Goal Reallocation',
        'Moving funds between savings goals',
        '#FFFFFF',
        '#2196F3',
        'savings_withdrawal'),
       ('Investment Opportunity',
        'Withdrawal to invest elsewhere',
        '#FFFFFF',
        '#9C27B0',
        'savings_withdrawal'),
       ('Unexpected Expense',
        'Withdrawal for unplanned costs',
        '#FFFFFF',
        '#FF5722',
        'savings_withdrawal'),
       ('Goal Cancellation',
        'Withdrawal due to cancelled goal',
        '#FFFFFF',
        '#9E9E9E',
        'savings_withdrawal'),
       ('Better Opportunity',
        'Withdrawal for better savings option',
        '#FFFFFF',
        '#00BCD4',
        'savings_withdrawal'),
       ('Family Need',
        'Withdrawal for family requirements',
        '#FFFFFF',
        '#E91E63',
        'savings_withdrawal'),
       ('Planned Purchase',
        'Withdrawal for intended purchase',
        '#FFFFFF',
        '#795548',
        'savings_withdrawal');


-- Sample ingresses using CreateIngress procedure
CALL CreateIngress(1, 85000.00, 7, '2024-01-15 09:00:00', 'January salary payment', 1, 'Employer Direct Deposit', TRUE,
                   '[
                     1,
                     9
                   ]');
CALL CreateIngress(2, 45000.00, 7, '2024-01-20 14:30:00', 'Freelance web development project', 3, 'Client Payment',
                   FALSE, '[
          3,
          6
        ]');
CALL CreateIngress(1, 12000.00, 7, '2024-01-25 16:45:00', 'Performance bonus Q4 2023', 1, 'Company Bonus', FALSE, '[
  2,
  9
]');
CALL CreateIngress(3, 8500.00, 7, '2024-02-01 10:15:00', 'Part-time job payment', 1, 'Part-time Employer', TRUE, '[
  1,
  9
]');
CALL CreateIngress(2, 25000.00, 7, '2024-02-05 11:20:00', 'Tax refund 2023', 4, 'Government Refund', FALSE, '[
  6
]');
CALL CreateIngress(1, 15000.00, 7, '2024-02-10 13:00:00', 'Investment dividends', 3, 'Stock Dividends', FALSE, '[
  4
]');
CALL CreateIngress(4, 30000.00, 7, '2024-02-14 09:30:00', 'Rental property income', 3, 'Tenant Payment', TRUE, '[
  8
]');
CALL CreateIngress(1, 5000.00, 7, '2024-02-18 15:45:00', 'Birthday gift from parents', 4, 'Family Gift', FALSE, '[
  5
]');
CALL CreateIngress(2, 18000.00, 7, '2024-02-22 12:10:00', 'Side business sales', 2, 'Online Store', FALSE, '[
  7,
  6
]');
CALL CreateIngress(3, 75.00, 150, '2024-02-28 16:00:00', 'Tutoring sessions payment', 2, 'Private Students', FALSE, '[
  3,
  6
]');

-- Sample expenditures using CreateExpenditure procedure
CALL CreateExpenditure(1, 35000.00, 7, '2024-01-16 18:30:00', 'Weekly grocery shopping at supermarket', 1, TRUE, TRUE,
                       '1,14', @exp_id, @trans_id);
CALL CreateExpenditure(2, 8500.00, 7, '2024-01-18 20:15:00', 'Dinner at Italian restaurant', 5, TRUE, FALSE, '5',
                       @exp_id, @trans_id);
CALL CreateExpenditure(1, 15000.00, 7, '2024-01-20 14:00:00', 'Monthly electricity bill', 2, TRUE, TRUE, '2', @exp_id,
                       @trans_id);
CALL CreateExpenditure(3, 100.00, 150, '2024-01-22 16:45:00', 'University textbooks for semester', 6, TRUE, TRUE, '6',
                       @exp_id, @trans_id);
CALL CreateExpenditure(1, 45000.00, 7, '2024-01-25 10:30:00', 'Car maintenance and oil change', 3, TRUE, FALSE, '3,8',
                       @exp_id, @trans_id);
CALL CreateExpenditure(2, 6500.00, 7, '2024-01-28 19:20:00', 'Movie tickets and popcorn', 5, TRUE, FALSE, '5', @exp_id,
                       @trans_id);
CALL CreateExpenditure(4, 25000.00, 7, '2024-02-02 11:15:00', 'Health insurance premium', 4, TRUE, TRUE, '4,9', @exp_id,
                       @trans_id);
CALL CreateExpenditure(1, 18000.00, 7, '2024-02-05 15:40:00', 'New winter jacket', 7, TRUE, FALSE, '7', @exp_id,
                       @trans_id);
CALL CreateExpenditure(3, 95.00, 7, '2024-02-08 13:25:00', 'Netflix and Spotify subscriptions', 9, TRUE, TRUE, '13',
                       @exp_id, @trans_id);
CALL CreateExpenditure(2, 32000.00, 7, '2024-02-12 17:10:00', 'Emergency dental treatment', 4, TRUE, FALSE, '4,12',
                       @exp_id, @trans_id);

-- Sample transfers using CreateTransfer procedure
CALL CreateTransfer(1, 4, 50000.00, 50000.00, 1.000000, 0.00, '2024-01-30 10:00:00',
                    'Monthly savings transfer to Jane account');
CALL CreateTransfer(2, 1, 200.00, 82000.00, 410.000000, 500.00, '2024-02-03 14:30:00',
                    'USD to ARS conversion for expenses');
CALL CreateTransfer(3, 1, 15000.00, 15000.00, 1.000000, 0.00, '2024-02-07 09:15:00', 'Alex allowance transfer');
CALL CreateTransfer(4, 5, 25000.00, 25000.00, 1.000000, 0.00, '2024-02-15 16:20:00',
                    'Transfer to Alex student account');
CALL CreateTransfer(1, 2, 100.00, 41000.00, 410.000000, 200.00, '2024-02-20 11:45:00',
                    'USD purchase for Europe trip savings');

-- Sample rollbacks using RollbackTransaction procedure
CALL RollbackTransaction(12, 'Duplicate transaction - correcting accounting error');
CALL RollbackTransaction(18, 'Incorrect amount charged - merchant refund processed');
CALL RollbackTransaction(6, 'Cancelled freelance project - payment returned');
CALL RollbackTransaction(21, 'Wrong account transfer - reversing transaction');
CALL RollbackTransaction(23, 'Exchange rate error - recalculating transfer');
-- Drop the stored procedures
DROP PROCEDURE IF EXISTS CreateExpenditure;
DROP PROCEDURE IF EXISTS CreateIngress;
DROP PROCEDURE IF EXISTS CreateTransfer;
DROP PROCEDURE IF EXISTS RollbackTransaction;
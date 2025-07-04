Create SCHEMA proletariat_budget;

use proletariat_budget;

create table household_members
(
    id         BIGINT AUTO_INCREMENT PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    surname    VARCHAR(255) NOT NULL,
    nickname   VARCHAR(255),
    role       varchar(255) NOT NULL,
    active     BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);


create table tags
(
    id               BIGINT auto_increment,
    name             VARCHAR(255) NOT NULL unique,
    description      TEXT,
    color            VARCHAR(255),
    background_color VARCHAR(255),
    type             varchar(255) NOT NULL,
    created_at       TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    index idx_name (name),
    index idx_type (type),
    PRIMARY KEY (id, name)

);

create table currencies
(
    id         INT PRIMARY KEY AUTO_INCREMENT,
    name       VARCHAR(50) NOT NULL,
    symbol     VARCHAR(10) NOT NULL,
    created_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

INSERT INTO currencies (name, symbol)
VALUES ('United Arab Emirates Dirham', 'AED'),
       ('Afghan Afghani', 'AFN'),
       ('Albanian Lek', 'ALL'),
       ('Armenian Dram', 'AMD'),
       ('Netherlands Antillean Guilder', 'ANG'),
       ('Angolan Kwanza', 'AOA'),
       ('Argentine Peso', 'ARS'),
       ('Australian Dollar', 'AUD'),
       ('Aruban Florin', 'AWG'),
       ('Azerbaijani Manat', 'AZN'),
       ('Bosnia-Herzegovina Convertible Mark', 'BAM'),
       ('Barbadian Dollar', 'BBD'),
       ('Bangladeshi Taka', 'BDT'),
       ('Bulgarian Lev', 'BGN'),
       ('Bahraini Dinar', 'BHD'),
       ('Burundian Franc', 'BIF'),
       ('Bermudan Dollar', 'BMD'),
       ('Brunei Dollar', 'BND'),
       ('Bolivian Boliviano', 'BOB'),
       ('Brazilian Real', 'BRL'),
       ('Bahamian Dollar', 'BSD'),
       ('Bitcoin', 'BTC'),
       ('Bhutanese Ngultrum', 'BTN'),
       ('Botswanan Pula', 'BWP'),
       ('Belarusian Ruble', 'BYN'),
       ('Belize Dollar', 'BZD'),
       ('Canadian Dollar', 'CAD'),
       ('Congolese Franc', 'CDF'),
       ('Swiss Franc', 'CHF'),
       ('Chilean Unit of Account (UF)', 'CLF'),
       ('Chilean Peso', 'CLP'),
       ('Chinese Yuan (Offshore)', 'CNH'),
       ('Chinese Yuan', 'CNY'),
       ('Colombian Peso', 'COP'),
       ('Costa Rican Colón', 'CRC'),
       ('Cuban Convertible Peso', 'CUC'),
       ('Cuban Peso', 'CUP'),
       ('Cape Verdean Escudo', 'CVE'),
       ('Czech Republic Koruna', 'CZK'),
       ('Djiboutian Franc', 'DJF'),
       ('Danish Krone', 'DKK'),
       ('Dominican Peso', 'DOP'),
       ('Algerian Dinar', 'DZD'),
       ('Egyptian Pound', 'EGP'),
       ('Eritrean Nakfa', 'ERN'),
       ('Ethiopian Birr', 'ETB'),
       ('Euro', 'EUR'),
       ('Fijian Dollar', 'FJD'),
       ('Falkland Islands Pound', 'FKP'),
       ('British Pound Sterling', 'GBP'),
       ('Georgian Lari', 'GEL'),
       ('Guernsey Pound', 'GGP'),
       ('Ghanaian Cedi', 'GHS'),
       ('Gibraltar Pound', 'GIP'),
       ('Gambian Dalasi', 'GMD'),
       ('Guinean Franc', 'GNF'),
       ('Guatemalan Quetzal', 'GTQ'),
       ('Guyanaese Dollar', 'GYD'),
       ('Hong Kong Dollar', 'HKD'),
       ('Honduran Lempira', 'HNL'),
       ('Croatian Kuna', 'HRK'),
       ('Haitian Gourde', 'HTG'),
       ('Hungarian Forint', 'HUF'),
       ('Indonesian Rupiah', 'IDR'),
       ('Israeli New Sheqel', 'ILS'),
       ('Manx pound', 'IMP'),
       ('Indian Rupee', 'INR'),
       ('Iraqi Dinar', 'IQD'),
       ('Iranian Rial', 'IRR'),
       ('Icelandic Króna', 'ISK'),
       ('Jersey Pound', 'JEP'),
       ('Jamaican Dollar', 'JMD'),
       ('Jordanian Dinar', 'JOD'),
       ('Japanese Yen', 'JPY'),
       ('Kenyan Shilling', 'KES'),
       ('Kyrgystani Som', 'KGS'),
       ('Cambodian Riel', 'KHR'),
       ('Comorian Franc', 'KMF'),
       ('North Korean Won', 'KPW'),
       ('South Korean Won', 'KRW'),
       ('Kuwaiti Dinar', 'KWD'),
       ('Cayman Islands Dollar', 'KYD'),
       ('Kazakhstani Tenge', 'KZT'),
       ('Laotian Kip', 'LAK'),
       ('Lebanese Pound', 'LBP'),
       ('Sri Lankan Rupee', 'LKR'),
       ('Liberian Dollar', 'LRD'),
       ('Lesotho Loti', 'LSL'),
       ('Libyan Dinar', 'LYD'),
       ('Moroccan Dirham', 'MAD'),
       ('Moldovan Leu', 'MDL'),
       ('Malagasy Ariary', 'MGA'),
       ('Macedonian Denar', 'MKD'),
       ('Myanma Kyat', 'MMK'),
       ('Mongolian Tugrik', 'MNT'),
       ('Macanese Pataca', 'MOP'),
       ('Mauritanian Ouguiya', 'MRU'),
       ('Mauritian Rupee', 'MUR'),
       ('Maldivian Rufiyaa', 'MVR'),
       ('Malawian Kwacha', 'MWK'),
       ('Mexican Peso', 'MXN'),
       ('Malaysian Ringgit', 'MYR'),
       ('Mozambican Metical', 'MZN'),
       ('Namibian Dollar', 'NAD'),
       ('Nigerian Naira', 'NGN'),
       ('Nicaraguan Córdoba', 'NIO'),
       ('Norwegian Krone', 'NOK'),
       ('Nepalese Rupee', 'NPR'),
       ('New Zealand Dollar', 'NZD'),
       ('Omani Rial', 'OMR'),
       ('Panamanian Balboa', 'PAB'),
       ('Peruvian Nuevo Sol', 'PEN'),
       ('Papua New Guinean Kina', 'PGK'),
       ('Philippine Peso', 'PHP'),
       ('Pakistani Rupee', 'PKR'),
       ('Polish Zloty', 'PLN'),
       ('Paraguayan Guarani', 'PYG'),
       ('Qatari Rial', 'QAR'),
       ('Romanian Leu', 'RON'),
       ('Serbian Dinar', 'RSD'),
       ('Russian Ruble', 'RUB'),
       ('Rwandan Franc', 'RWF'),
       ('Saudi Riyal', 'SAR'),
       ('Solomon Islands Dollar', 'SBD'),
       ('Seychellois Rupee', 'SCR'),
       ('Sudanese Pound', 'SDG'),
       ('Swedish Krona', 'SEK'),
       ('Singapore Dollar', 'SGD'),
       ('Saint Helena Pound', 'SHP'),
       ('Sierra Leonean Leone', 'SLL'),
       ('Somali Shilling', 'SOS'),
       ('Surinamese Dollar', 'SRD'),
       ('South Sudanese Pound', 'SSP'),
       ('São Tomé and Príncipe Dobra (pre-2018)', 'STD'),
       ('São Tomé and Príncipe Dobra', 'STN'),
       ('Salvadoran Colón', 'SVC'),
       ('Syrian Pound', 'SYP'),
       ('Swazi Lilangeni', 'SZL'),
       ('Thai Baht', 'THB'),
       ('Tajikistani Somoni', 'TJS'),
       ('Turkmenistani Manat', 'TMT'),
       ('Tunisian Dinar', 'TND'),
       ('Tongan Pa\'anga', 'TOP'),
       ('Turkish Lira', 'TRY'),
       ('Trinidad and Tobago Dollar', 'TTD'),
       ('New Taiwan Dollar', 'TWD'),
       ('Tanzanian Shilling', 'TZS'),
       ('Ukrainian Hryvnia', 'UAH'),
       ('Ugandan Shilling', 'UGX'),
       ('United States Dollar', 'USD'),
       ('Uruguayan Peso', 'UYU'),
       ('Uzbekistan Som', 'UZS'),
       ('Venezuelan Bolívar Fuerte (Old)', 'VEF'),
       ('Venezuelan Bolívar Soberano', 'VES'),
       ('Vietnamese Dong', 'VND'),
       ('Vanuatu Vatu', 'VUV'),
       ('Samoan Tala', 'WST'),
       ('CFA Franc BEAC', 'XAF'),
       ('Silver Ounce', 'XAG'),
       ('Gold Ounce', 'XAU'),
       ('East Caribbean Dollar', 'XCD'),
       ('Special Drawing Rights', 'XDR'),
       ('CFA Franc BCEAO', 'XOF'),
       ('Palladium Ounce', 'XPD'),
       ('CFP Franc', 'XPF'),
       ('Platinum Ounce', 'XPT'),
       ('Yemeni Rial', 'YER'),
       ('South African Rand', 'ZAR'),
       ('Zambian Kwacha', 'ZMW'),
       ('Zimbabwean Dollar', 'ZWL');

CREATE TABLE accounts
(
    id                  BIGINT auto_increment PRIMARY KEY,
    name                VARCHAR(255)                                           NOT NULL,
    type                ENUM ('bank', 'cash', 'investment', 'crypto', 'other') NOT NULL,
    owner               bigint                                                 NOT NULL,
    institution         VARCHAR(255),
    currency            int                                                    NOT NULL,
    initial_balance     DECIMAL(15, 2)                                         NOT NULL,
    current_balance     DECIMAL(15, 2)                                         NOT NULL,
    active              BOOLEAN                                                NOT NULL DEFAULT TRUE,
    description         TEXT,
    account_number      VARCHAR(255),
    account_information TEXT,
    created_at          TIMESTAMP                                              NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP                                              NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

alter table accounts
    add constraint fk_accounts_owner FOREIGN KEY (owner) REFERENCES household_members (id);
alter table accounts
    add constraint fk_accounts_currency FOREIGN KEY (currency) REFERENCES currencies (id);


CREATE INDEX idx_accounts_type ON accounts (type);
CREATE INDEX idx_accounts_currency ON accounts (currency);
CREATE INDEX idx_accounts_active ON accounts (active);

CREATE TABLE transactions
(
    id               BIGINT PRIMARY KEY auto_increment,
    account_id       BIGINT                                                                                    NOT NULL,
    amount           DECIMAL(15, 2)                                                                            NOT NULL,
    currency         int                                                                                       NOT NULL,
    transaction_date DATETIME                                                                                  NOT NULL,
    description      TEXT,
    transaction_type ENUM ('expenditure', 'ingress', 'transfer', 'savings_contribution', 'savings_withdrawal') NOT NULL,
    balance_after    DECIMAL(15, 2)                                                                            NOT NULL,
    status           ENUM ('pending', 'completed', 'failed', 'cancelled')                                      NOT NULL DEFAULT 'completed',
    created_at       TIMESTAMP                                                                                 NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP                                                                                 NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    -- Foreign key to accounts table
    CONSTRAINT fk_transaction_account FOREIGN KEY (account_id) REFERENCES accounts (id),
    constraint fk_transaction_currency FOREIGN KEY (currency) REFERENCES currencies (id)
);

-- Create indexes for faster lookups
CREATE INDEX idx_transactions_account_id ON transactions (account_id);
CREATE INDEX idx_transactions_date ON transactions (transaction_date);
CREATE INDEX idx_transactions_type ON transactions (transaction_type);


create table categories
(
    id               BIGINT auto_increment PRIMARY KEY,
    name             VARCHAR(255) NOT NULL,
    description      TEXT,
    color            VARCHAR(255),
    background_color VARCHAR(255),
    active           BOOLEAN      NOT NULL DEFAULT TRUE,
    category_type    VARCHAR(255) NOT NULL,
    index idx_name (name)
);


CREATE TABLE expenditures
(
    id             BIGINT auto_increment PRIMARY KEY,
    category_id    BIGINT    NOT NULL,
    declared       BOOLEAN   NOT NULL DEFAULT FALSE,
    planned        BOOLEAN   NOT NULL DEFAULT FALSE,
    transaction_id BIGINT    NOT NULL unique REFERENCES transactions (id),
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_category (category_id),
    INDEX idx_declared (declared),
    INDEX idx_planned (planned),
    FOREIGN KEY (category_id) REFERENCES categories (id)
);

CREATE TABLE expenditure_tags
(
    expenditure_id BIGINT NOT NULL
        REFERENCES expenditures (id),
    tag_id         BIGINT NOT NULL
        REFERENCES tags (id),
    index (expenditure_id, tag_id)
);

ALTER TABLE expenditure_tags
    ADD CONSTRAINT fk_expenditure_id FOREIGN KEY (expenditure_id) REFERENCES expenditures (id);

ALTER TABLE expenditure_tags
    ADD CONSTRAINT fk_tag_id FOREIGN KEY (tag_id) REFERENCES tags (id);


-- Create ingresses table
CREATE TABLE ingresses
(
    id             BIGINT auto_increment PRIMARY KEY,
    category_id       bigint    NOT NULL,
    source         VARCHAR(255),
    is_recurring   BOOLEAN            DEFAULT FALSE,
    transaction_id BIGINT    NOT NULL unique REFERENCES transactions (id),
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_category (category_id),
    INDEX idx_source (source),
    foreign key (category_id) REFERENCES categories (id)
);

-- Create table for ingress tags junction with ingresses
Create table ingress_tags
(
    ingress_id BIGINT NOT NULL,
    tag_id     BIGINT NOT NULL,
    index (ingress_id, tag_id),
    FOREIGN KEY (ingress_id) REFERENCES ingresses (id),
    FOREIGN KEY (tag_id) REFERENCES tags (id)
);

-- Create table for recurrence patterns
CREATE TABLE ingress_recurrence_patterns
(
    id             BIGINT auto_increment PRIMARY KEY,
    ingress_id     BIGINT                                        NOT NULL,
    frequency      ENUM ('daily', 'weekly', 'monthly', 'yearly') NOT NULL,
    interval_value INT                                           NOT NULL DEFAULT 1,
    amount         DECIMAL(15, 2),
    end_date       DATE,
    foreign key (ingress_id) REFERENCES ingresses (id) ON DELETE CASCADE
);

-- Create index for common query patterns
CREATE INDEX idx_ingresses_category ON ingresses (category_id);
CREATE INDEX idx_ingresses_is_recurring ON ingresses (is_recurring);

-- Create savings goals table
CREATE TABLE savings_goals
(
    id                        BIGINT PRIMARY KEY auto_increment,
    name                      VARCHAR(255)                              NOT NULL,
    category_id                  BIGINT                                    NOT NULL,
    description               TEXT,
    target_amount             DECIMAL(15, 2)                            NOT NULL,
    currency                  int                                       NOT NULL,
    target_date               DATE,
    initial_amount            DECIMAL(15, 2)                            NOT NULL DEFAULT 0.00,
    current_amount            DECIMAL(15, 2)                            NOT NULL DEFAULT 0.00,
    percent_complete          DECIMAL(5, 2)                             NOT NULL DEFAULT 0.00,
    account_id                BIGINT                                    NOT NULL,
    priority                  INT,
    auto_contribute           BOOLEAN                                            DEFAULT FALSE,
    auto_contribute_amount    DECIMAL(15, 2),
    auto_contribute_frequency ENUM ('daily', 'weekly', 'monthly', 'yearly'),
    status                    ENUM ('active', 'completed', 'abandoned', 'inactive') NOT NULL DEFAULT 'active',
    projected_completion_date DATE,
    created_at                TIMESTAMP                                 NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                TIMESTAMP                                 NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (account_id) REFERENCES accounts (id),
    FOREIGN KEY (category_id) REFERENCES categories (id),
    foreign key (currency) REFERENCES currencies (id)
);


-- Create savings goal tags junction table
CREATE TABLE savings_goal_tags
(
    savings_goal_id BIGINT NOT NULL,
    tag_id          BIGINT NOT NULL,
    PRIMARY KEY (savings_goal_id, tag_id),
    FOREIGN KEY (tag_id) REFERENCES tags (id) ON DELETE CASCADE,
    FOREIGN KEY (savings_goal_id) REFERENCES savings_goals (id) ON DELETE CASCADE,
    INDEX (savings_goal_id, tag_id)
);

-- Create savings contributions table
CREATE TABLE savings_contributions
(
    id                BIGINT PRIMARY KEY auto_increment,
    savings_goal_id   BIGINT    NOT NULL,
    date              DATE      NOT NULL,
    source_account_id BIGINT    NOT NULL,
    notes             TEXT,
    created_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (savings_goal_id) REFERENCES savings_goals (id) ON DELETE CASCADE,
    FOREIGN KEY (source_account_id) REFERENCES accounts (id)
);

-- Create savings contribution tags table
CREATE TABLE savings_contribution_tags
(
    contribution_id BIGINT NOT NULL,
    tag_id          BIGINT NOT NULL,
    PRIMARY KEY (contribution_id, tag_id),
    FOREIGN KEY (contribution_id) REFERENCES savings_contributions (id),
    foreign key (tag_id) REFERENCES tags (id)
);

-- Create savings withdrawals table
CREATE TABLE savings_withdrawals
(
    id                     BIGINT auto_increment PRIMARY KEY,
    savings_goal_id        BIGINT       NOT NULL,
    date                   DATE         NOT NULL,
    destination_account_id BIGINT,
    reason                 VARCHAR(255) NOT NULL,
    notes                  TEXT,
    created_at             TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at             TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (savings_goal_id) REFERENCES savings_goals (id) ON DELETE CASCADE,
    FOREIGN KEY (destination_account_id) REFERENCES accounts (id)
);

-- Create savings withdrawal tags table
CREATE TABLE savings_withdrawal_tags
(
    withdrawal_id BIGINT NOT NULL,
    tag_id        BIGINT NOT NULL,
    PRIMARY KEY (withdrawal_id, tag_id),
    FOREIGN KEY (withdrawal_id) REFERENCES savings_withdrawals (id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags (id)
);

-- Create indexes for common query patterns
CREATE INDEX idx_savings_goals_account_id ON savings_goals (account_id);
CREATE INDEX idx_savings_goals_status ON savings_goals (status);
CREATE INDEX idx_savings_goals_currency ON savings_goals (currency);
CREATE INDEX idx_savings_contributions_goal_id ON savings_contributions (savings_goal_id);
CREATE INDEX idx_savings_contributions_date ON savings_contributions (date);
CREATE INDEX idx_savings_withdrawals_goal_id ON savings_withdrawals (savings_goal_id);
CREATE INDEX idx_savings_withdrawals_date ON savings_withdrawals (date);

CREATE TABLE exchange_rates
(
    id              BIGINT PRIMARY KEY auto_increment,
    base_currency   int            NOT NULL,
    target_currency int            NOT NULL,
    rate            DECIMAL(15, 6) NOT NULL,
    date            DATE           NOT NULL,
    created_at      TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    -- Indexes for faster lookups
    INDEX idx_base_currency (base_currency),
    INDEX idx_target_currency (target_currency),
    INDEX idx_date (date),
    -- Unique constraint to prevent duplicate entries for the same currency pair on the same date
    UNIQUE KEY uk_currency_pair_date (base_currency, target_currency, date),
    foreign key (base_currency) REFERENCES currencies (id),
    foreign key (target_currency) REFERENCES currencies (id)
);

-- Link to transfers (assuming a transfers table exists or will be created)
CREATE TABLE transfers
(
    id                       BIGINT PRIMARY KEY auto_increment,
    source_account_id        BIGINT    NOT NULL,
    destination_account_id   BIGINT    NOT NULL,
    destination_amount       DECIMAL(15, 2),
    exchange_rate_multiplier DECIMAL(15, 6),
    fees                     DECIMAL(15, 2),
    transaction_id           BIGINT    NOT NULL unique REFERENCES transactions (id),
    created_at               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at               TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_transfer_source_account FOREIGN KEY (source_account_id) REFERENCES accounts (id),
    CONSTRAINT fk_transfer_destination_account FOREIGN KEY (destination_account_id) REFERENCES accounts (id)
);

create table users
(
    id            BIGINT PRIMARY KEY AUTO_INCREMENT,
    username      VARCHAR(50) UNIQUE  NOT NULL,
    password_hash VARCHAR(255)        NOT NULL,
    salt          VARCHAR(255)        NOT NULL,
    email         VARCHAR(255) UNIQUE NOT NULL,
    name          VARCHAR(50)         NOT NULL,
    surname       VARCHAR(50)         NOT NULL,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

create table roles
(
    id         BIGINT PRIMARY KEY AUTO_INCREMENT,
    name       VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

create table user_roles
(
    user_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,

    PRIMARY KEY (user_id, role_id),
    CONSTRAINT fk_user_roles_user FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT fk_user_roles_role FOREIGN KEY (role_id) REFERENCES roles (id)
);
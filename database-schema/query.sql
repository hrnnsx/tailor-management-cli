DROP TRIGGER IF EXISTS trg_after_user_insert;
DROP TRIGGER IF EXISTS trg_after_user_update;

DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS size_requirements;
DROP TABLE IF EXISTS fabric_patterns;
DROP TABLE IF EXISTS patterns;
DROP TABLE IF EXISTS fabrics;
DROP TABLE IF EXISTS user_measurements;
DROP TABLE IF EXISTS workers;
DROP TABLE IF EXISTS users;

-- 1. Users
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role ENUM('customer', 'admin', 'worker') NOT NULL DEFAULT 'customer',
    phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Workers (Subtype 1:1)
CREATE TABLE workers (
    user_id INT PRIMARY KEY,
    availability BOOLEAN NOT NULL DEFAULT TRUE,

    CONSTRAINT fk_workers_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);


-- 3. User Measurements
CREATE TABLE user_measurements (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    title VARCHAR(50) NOT NULL,
    height_cm DECIMAL(5, 2),
    chest_circumference DECIMAL(5, 2) NOT NULL,
    waist_circumference DECIMAL(5, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_measurements_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- 4. Fabrics
CREATE TABLE fabrics (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    price_per_cm DECIMAL(10, 2) NOT NULL
);

-- 5. Patterns
CREATE TABLE patterns (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

-- 6. Fabric Patterns (Stok Kombinasi - Sentimeter)
CREATE TABLE fabric_patterns (
    id INT AUTO_INCREMENT PRIMARY KEY,
    fabric_id INT NOT NULL,
    pattern_id INT NOT NULL,
    stock_cm INT NOT NULL DEFAULT 0,

    CONSTRAINT fk_fp_fabric
        FOREIGN KEY (fabric_id)
        REFERENCES fabrics(id),

    CONSTRAINT fk_fp_pattern
        FOREIGN KEY (pattern_id)
        REFERENCES patterns(id),

    UNIQUE KEY uk_fabric_pattern (fabric_id, pattern_id)
);

-- 7. Size Requirements (Kebutuhan Kain per Ukuran - Sentimeter)
CREATE TABLE size_requirements (
    id INT AUTO_INCREMENT PRIMARY KEY,
    size ENUM('XS', 'S', 'M', 'L', 'XL', 'XXL') NOT NULL UNIQUE,
    required_cm INT NOT NULL
);


-- 8. Orders
CREATE TABLE orders (
    id INT AUTO_INCREMENT PRIMARY KEY,

    order_code VARCHAR(50) NOT NULL UNIQUE,

    customer_id INT NOT NULL,

    assigned_worker_id INT DEFAULT NULL,

    user_measurement_id INT NOT NULL,

    fabric_pattern_id INT NOT NULL,

    determined_size ENUM(
        'XS',
        'S',
        'M',
        'L',
        'XL',
        'XXL'
    ) NOT NULL,

    cm_used INT NOT NULL,

    price_per_cm_snapshot DECIMAL(10, 2) NOT NULL,

    total_price DECIMAL(10, 2) NOT NULL,

    payment_status ENUM(
        'unpaid',
        'paid'
    ) NOT NULL DEFAULT 'unpaid',

    status ENUM(
        'pending',
        'in progress',
        'waiting payment',
        'finished'
    ) NOT NULL DEFAULT 'pending',

    progress ENUM(
        'Order diterima',
        'Persiapan bahan',
        'Pemotongan kain',
        'Proses jahit',
        'Finishing',
        'Selesai'
    ) NOT NULL DEFAULT 'Order diterima',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_orders_customer
        FOREIGN KEY (customer_id)
        REFERENCES users(id),

    CONSTRAINT fk_orders_worker
        FOREIGN KEY (assigned_worker_id)
        REFERENCES workers(user_id)
        ON DELETE SET NULL,

    CONSTRAINT fk_orders_measurement
        FOREIGN KEY (user_measurement_id)
        REFERENCES user_measurements(id),

    CONSTRAINT fk_orders_fabric_pattern
        FOREIGN KEY (fabric_pattern_id)
        REFERENCES fabric_patterns(id)
);

-- 9. Payments
CREATE TABLE payments (
    id INT AUTO_INCREMENT PRIMARY KEY,

    order_id INT NOT NULL,

    amount DECIMAL(10, 2) NOT NULL,

    status ENUM(
        'pending',
        'verified',
        'rejected'
    ) NOT NULL DEFAULT 'pending',

    verified_by INT DEFAULT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_payments_order
        FOREIGN KEY (order_id)
        REFERENCES orders(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_payments_verifier
        FOREIGN KEY (verified_by)
        REFERENCES users(id)
        ON DELETE SET NULL
);

-- TRIGGERS: Auto Sinkronisasi users -> workers
DELIMITER $$

CREATE TRIGGER trg_after_user_insert
AFTER INSERT ON users
FOR EACH ROW
BEGIN
    IF NEW.role = 'worker' THEN
        INSERT INTO workers (
            user_id,
            availability
        )
        VALUES (
            NEW.id,
            TRUE
        );
    END IF;
END$$


CREATE TRIGGER trg_after_user_update
AFTER UPDATE ON users
FOR EACH ROW
BEGIN

    -- Jika user berubah menjadi worker
    IF NEW.role = 'worker'
       AND OLD.role != 'worker' THEN

        INSERT INTO workers (
            user_id,
            availability
        )
        VALUES (
            NEW.id,
            TRUE
        )
        ON DUPLICATE KEY UPDATE
            availability = TRUE;

    -- Jika worker berubah menjadi role lain
    ELSEIF NEW.role != 'worker'
       AND OLD.role = 'worker' THEN

        DELETE FROM workers
        WHERE user_id = NEW.id;

    END IF;

END$$

DELIMITER ;

INSERT INTO size_requirements (size, required_cm)
VALUES
    ('XS', 200),
    ('S', 210),
    ('M', 220),
    ('L', 230),
    ('XL', 240),
    ('XXL', 250);
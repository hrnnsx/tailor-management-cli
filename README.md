# Tailor Management CLI

A comprehensive command-line application for managing a tailor shop's operations, including customer orders, fabric inventory, measurements, and worker assignments.

## Overview

This system manages the complete workflow of a tailor shop, from customer measurements and fabric selection to order tracking and payment verification.

## Database Schema

### Tables

#### 1. **users**
Core user accounts table supporting three roles: customer, admin, and worker.
- `id` (INT, PK): Auto-incremented user ID
- `name` (VARCHAR): User's full name
- `email` (VARCHAR, UNIQUE): User email address
- `password` (VARCHAR): Hashed password
- `role` (ENUM): One of `customer`, `admin`, or `worker`
- `phone` (VARCHAR): Contact number
- `created_at` (TIMESTAMP): Account creation timestamp

#### 2. **workers**
Extension table for worker users (1:1 relationship with users).
- `user_id` (INT, PK, FK): Foreign key to users
- `availability` (BOOLEAN): Whether the worker is available for new assignments

#### 3. **user_measurements**
Stores multiple measurement profiles for each customer.
- `id` (INT, PK): Auto-incremented measurement ID
- `user_id` (INT, FK): Reference to customer
- `title` (VARCHAR): Name of the measurement profile
- `height_cm` (DECIMAL): Height in centimeters
- `chest_circumference` (DECIMAL): Chest circumference
- `waist_circumference` (DECIMAL): Waist circumference
- `hip_circumference` (DECIMAL): Hip circumference
- `shoulder_width` (DECIMAL): Shoulder width
- `arm_length` (DECIMAL): Arm length
- `created_at` (TIMESTAMP): Creation timestamp

#### 4. **fabrics**
Inventory of available fabrics.
- `id` (INT, PK): Auto-incremented fabric ID
- `name` (VARCHAR): Fabric name/description
- `price_per_meter` (DECIMAL): Price per meter in currency units

#### 5. **patterns**
Available clothing patterns/designs.
- `id` (INT, PK): Auto-incremented pattern ID
- `name` (VARCHAR): Pattern name/description

#### 6. **fabric_patterns**
Junction table for fabric and pattern combinations with stock tracking (stock measured in centimeters).
- `id` (INT, PK): Auto-incremented ID
- `fabric_id` (INT, FK): Reference to fabric
- `pattern_id` (INT, FK): Reference to pattern
- `stock_cm` (INT): Available stock in centimeters
- **UNIQUE**: `(fabric_id, pattern_id)` - ensures each combination is unique

#### 7. **size_requirements**
Standard size chart defining fabric requirement per size.
- `id` (INT, PK): Auto-incremented ID
- `size` (ENUM): Standard size: `XS`, `S`, `M`, `L`, `XL`, `XXL`
- `required_cm` (INT): Fabric requirement in centimeters for this size

#### 8. **orders**
Customer orders linking measurements, fabric selections, and assignments.
- `id` (INT, PK): Auto-incremented order ID
- `order_code` (VARCHAR, UNIQUE): Unique order identifier
- `customer_id` (INT, FK): Reference to customer (users table)
- `assigned_worker_id` (INT, FK): Reference to assigned worker (workers table, nullable)
- `user_measurement_id` (INT, FK): Reference to customer's measurement profile
- `fabric_pattern_id` (INT, FK): Reference to selected fabric-pattern combination
- `determined_size` (ENUM): Final size determined for the order (`XS`-`XXL`)
- `cm_used` (INT): Actual fabric used in centimeters
- `price_per_meter_snapshot` (DECIMAL): Price snapshot at order time
- `total_price` (DECIMAL): Total order price
- `payment_status` (ENUM): `unpaid` or `paid`
- `status` (ENUM): Order status: `pending`, `in progress`, or `finished`
- `created_at` (TIMESTAMP): Order creation timestamp

#### 9. **payments**
Payment records for orders.
- `id` (INT, PK): Auto-incremented payment ID
- `order_id` (INT, FK): Reference to order
- `amount` (DECIMAL): Payment amount
- `proof_image_url` (VARCHAR): URL to payment proof image
- `status` (ENUM): Payment status: `pending`, `verified`, or `rejected`
- `verified_by` (INT, FK): Reference to verifying admin (users table, nullable)
- `created_at` (TIMESTAMP): Payment timestamp

### Triggers

#### `trg_after_user_insert`
**Purpose**: Auto-sync worker role to workers table
- **Event**: After INSERT on users table
- **Logic**: If new user has role = 'worker', automatically create a workers record with availability = TRUE

#### `trg_after_user_update`
**Purpose**: Keep workers table in sync with user roles
- **Event**: After UPDATE on users table
- **Logic**:
  - If user role changed TO 'worker', create workers record
  - If user role changed FROM 'worker', delete workers record

## Entity Relationship Diagram

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         DATABASE SCHEMA DIAGRAM                         │
└─────────────────────────────────────────────────────────────────────────┘

                              ┌──────────┐
                              │  users   │
                              ├──────────┤
                              │ id (PK)  │
                              │ name     │
                              │ email    │
                              │ password │
                              │ role     │
                              │ phone    │
                              │created_at│
                              └──────────┘
                                   │
                  ┌────────────────┼────────────────┐
                  │                │                │
                  ▼                ▼                ▼
           ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
           │   workers    │ │ user_        │ │  payments    │
           │              │ │ measurements │ │              │
           ├──────────────┤ ├──────────────┤ ├──────────────┤
           │user_id (FK)  │ │id (PK)       │ │id (PK)       │
           │availability  │ │user_id (FK)  │ │order_id (FK) │
           └──────────────┘ │title         │ │amount        │
                            │height_cm     │ │proof_image_url
                            │chest_circ    │ │status        │
                            │waist_circ    │ │verified_by   │
                            │hip_circ      │ │created_at    │
                            │shoulder_width│ └──────────────┘
                            │arm_length    │      │
                            │created_at    │      │
                            └──────────────┘      │
                                   │              │
                                   └──────────┐   │
                                              ▼   ▼
                                         ┌──────────────┐
                                         │   orders     │
                                         ├──────────────┤
                                         │id (PK)       │
                                         │order_code    │
                                         │customer_id   │
                                         │assigned_worker│
                                         │measurement_id│
                                         │fabric_pattern│
                                         │determined_sz │
                                         │cm_used       │
                                         │price_snapshot│
                                         │total_price   │
                                         │payment_status│
                                         │status        │
                                         │created_at    │
                                         └──────────────┘
                                              │
                                              │
                           ┌──────────────────┘
                           ▼
                    ┌──────────────────┐
                    │ fabric_patterns  │
                    ├──────────────────┤
                    │id (PK)           │
                    │fabric_id (FK)    │
                    │pattern_id (FK)   │
                    │stock_cm          │
                    └──────────────────┘
                         │      │
            ┌────────────┘      └────────────┐
            ▼                                 ▼
      ┌──────────────┐              ┌──────────────┐
      │   fabrics    │              │   patterns   │
      ├──────────────┤              ├──────────────┤
      │id (PK)       │              │id (PK)       │
      │name          │              │name          │
      │price_per_m   │              └──────────────┘
      └──────────────┘

      ┌──────────────────────┐
      │ size_requirements    │
      ├──────────────────────┤
      │id (PK)               │
      │size (UNIQUE)         │
      │required_cm           │
      └──────────────────────┘
      (Referenced by orders)
```

## Key Relationships

- **Users → Workers**: One-to-one for worker users
- **Users → Measurements**: One-to-many (customers have multiple measurements)
- **Users → Orders**: One-to-many (customers place multiple orders)
- **Users → Payments**: One-to-many via orders (admins verify payments)
- **Measurements → Orders**: One-to-many (a measurement can be used in multiple orders)
- **Fabrics ↔ Patterns**: Many-to-many through fabric_patterns junction table
- **Orders → Fabric_patterns**: Many-to-one (one fabric-pattern combination per order)

## Features

### Order Management
- Create and track customer orders
- Assign workers to orders
- Track order status (pending → in progress → finished)
- Calculate pricing based on size and fabric selection

### Inventory Management
- Track fabric and pattern combinations
- Stock management in centimeters
- Size-based requirement calculation

### Measurement System
- Store multiple measurement profiles per customer
- Reference measurements when creating orders

### Payment Tracking
- Record payment details and proof images
- Track payment status (pending → verified/rejected)
- Admin verification workflow

### Worker Management
- Track worker availability
- Assign workers to orders
- Role-based access control

## Installation

1. Ensure MySQL/MariaDB is installed
2. Run the database schema script:
   ```bash
   mysql -u root -p < database-schema/query.sql
   ```

## Database Schema File

Located in `database-schema/query.sql`

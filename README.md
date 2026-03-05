# BizKeeper

BizKeeper is a digital ledger platform where business owners record their bulk product purchases in structured purchase sessions .

The system preserves the context of each buying event and transforms raw records into organized financial insight while maintaining strict data privacy.

This directory contains the Go-based backend server which powers the entire BizKeeper ecosystem.


## Architectural Identity

BizKeeper is a multi-tenant, purchase session-based ledger system that provides structured financial tracking for business product purchases while enforcing strict data privacy and role-based access control.

It is simple, focused, and purpose-built for secure digital purchase management.

## Users & Roles
The system operates strictly on a multi-tenant RBAC (Role-Based Access Control) system:


### Business Owner

- Registers and logs in
- Creates purchase sessions
- Adds multiple product items per session
- Manages only their own records
- Cannot see other users' data

### Admin (CEO)

- Views platform-level statistics only.
- Restricted to high-level platform analytics (total users, macro system engagement). _Cannot_ access individual business ledgers.

## Core Structure

### Purchase Session (Parent Record)

Represents one complete buying event.
Contains:

- Purchase Date
- Supplier Name
- Payment Method
- Invoice Reference
- Optional General Note
- Automatically linked to Owner

**Example:**

- Owner goes to market on March 4
- Creates one Purchase Session
- Adds all purchased products inside that session

### Product Items (Child Records)

Inside each purchase session, the owner can add multiple products.
Each product item contains:

- Name
- Quantity
- Unit Price
- Category
- Notes (optional)

_If the owner buys 20 different products:_

- 1 Purchase Session is created
- 20 Product Items are saved under that session

This structure preserves context and ensures clean organization.

## Core Action Flow

1. Owner logs in
2. Creates new Purchase Session
3. Adds multiple Product Items
4. Clicks save

**System Actions:**

- Saves one Purchase Session
- Saves all Product Items linked to that session
- Associates everything strictly with the authenticated owner

## Extra Features

### Spending Summary

- Weekly breakdown
- Monthly breakdown
- Category-based spending
- Supplier-based spending
- **Calculated per owner only**

### Price History

Tracked automatically because each product item stores:

- Name
- Unit Price
- Date
  _(Price history is derived from historical product records)._

### Export

Owner can export:

- Purchase Sessions
- Product Items
- Spending Summaries

_Supported Formats: PDF, CSV_

### Reorder Reminder

System checks the last recorded purchase date of a product and flags products not restocked within a defined time window. Calculated dynamically from existing data.


## Core Routing Flow

The API routes are managed strictly by authentication and role middleware.

1. **Public Routes** (`/api/auth`)
   - `/register`: Allows Business Owners to create accounts.
   - `/login`: Generates JWT tokens for secure authentication.

2. **Protected User Routes** (`/api/user`)
   - `/update`: Updates the currently authenticated user's profile.
   - `/delete`: Deletes the user profile.

3. **Protected Purchase Routes** (`/api/purchases`)
   - `POST ""`: Creates a new Purchase Session.
   - `GET ""`: Lists user's Purchase Sessions.
   - `GET "/:id"`: Retrieves a specific Purchase Session.
   - `PUT "/:id"`: Updates a specific Purchase Session.
   - `DELETE "/:id"`: Deletes a specific Purchase Session.
   - **Product Items Management:** Manage products within sessions dynamically (`PUT "/:id/items/:itemId"`, `DELETE "/:id/items/:itemId"`).

## Backend Tech Stack

- **Framework:** Gin
- **Database:** PostgreSQL with GORM
- **Live Reload:** Air
- **Authentication:** JWT
- **Configuration:** godotenv
- **Routing Structure:** Clean architecture with grouped API routers (`/api`)


## Server Directory Structure

```text
server/
├── cmd/
│   └── server/
│       └── main.go     # Entrypoint, DB Init, Route Setup
├── config/
│   └── config.go       # Load env vars
├── internal/
│   ├── db/
│   │   └── db.go       # Postgres configuration
│   ├── handler/
│   │   ├── auth_handler.go     # Handles auth JSON requests
│   │   ├── product_handler.go
│   │   └── purchase_handler.go     # Handles ledger JSON requests
│   ├── middleware/
│   │   ├── auth_middleware.go      # JWT Verification
│   │   └── role_middleware.go      # Role gating (e.g., Admin stats)
│   ├── model/
│   │   ├── User.go                 # Owner/Admin models
│   │   ├── product_item.go         # Product schema
│   │   └── purchase_session.go     # Purchase schema
│   ├── repository/
│   │   ├── product_repo.go         # GORM interactions for Products
│   │   ├── purchase_repo.go        # GORM interactions for Sessions
│   │   └── user_repo.go            # GORM interactions for Users
│   ├── routes/
│   │   ├── auth_route.go           # Auth routes
│   │   ├── product_route.go        # Product routes
│   │   └── purchase_route.go       # Purchase routes
│   └── services/
│       ├── auth_service.go
│       ├── product_service.go
│       └── purchase_service.go
├── pkg/
│   ├── response/
│   └── utils/
│       ├── hashedPassword.go
│       └── jwt.go
├── .air.toml
├── .env                        # Environment variables (Excluded from repo)
├── .gitignore                  # Git ignore file
├── go.mod                      # Go module file
└── go.sum                      # Go module sum file
```


## Security & Privacy Posture

Every protected route requires validation via the `auth_middleware`. The middleware securely derives the active `Owner ID` directly from the validated token. Repositories automatically enforce this `Owner ID` on every database transaction, guaranteeing end-to-end data isolation with zero risk of cross-tenant exposure.
# BizKeeper

BizKeeper is a digital ledger platform where business owners record their bulk product purchases in structured purchase sessions instead of writing them in a physical book.

The system preserves the context of each buying event and transforms raw records into organized financial insight while maintaining strict data privacy.

## Architectural Identity

BizKeeper is a multi-tenant, purchase session-based ledger system that provides structured financial tracking for business product purchases while enforcing strict data privacy and role-based access control.

It is simple, focused, and purpose-built for secure digital purchase management.

## Users & Roles

### Business Owner

- Registers and logs in
- Creates purchase sessions
- Adds multiple product items per session
- Manages only their own records
- Cannot see other users' data

### Admin (CEO)

- Views platform-level statistics only
- Sees number of users
- Sees total purchase sessions
- Sees total product items recorded
- Never accesses individual business records
- Fully role-gated

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

## Security Model

- Authentication required for every endpoint
- Owner ID always derived from authentication token
- Every record strictly tied to its owner
- All database queries filtered by Owner ID
- Admin endpoints restricted to aggregate queries only
- Strict multi-tenant data isolation

## Server Folder Structure

```text
server/
├── cmd/
├── config/
├── internal/
│   ├── db/
│   ├── handler/
│   ├── middleware/
│   ├── model/
│   ├── repository/
│   ├── routes/
│   └── services/
├── pkg/
│   ├── response/
│   └── utils/
│       └── jwt.go
├── .env
├── .gitignore
├── go.mod
└── go.sum
```

# Evermos Internship — REST API

Final project for the Rakamin Project Based Internship at Evermos.

A RESTful back-end service for a marketplace platform built with **Go**, **Fiber v2**, and **GORM** (MySQL). Features user authentication, multi-store product management, address management, and order transactions.

---

## Tech Stack

| Layer | Library |
|---|---|
| HTTP framework | [Fiber v2](https://github.com/gofiber/fiber) |
| ORM | [GORM](https://gorm.io) + MySQL driver |
| Auth | JWT (`dgrijalva/jwt-go` + `gofiber/jwt`) |
| External data | [emsifa Indonesia Region API](https://github.com/emsifa/api-wilayah-indonesia) |

---

## Getting Started

### Local (with .env)

```bash
cp .env.example .env   # fill in your values
go run main.go
```

`.env` variables:

```
PORT=3000
DB_USER=your_mysql_user
DB_PASSWORD=your_mysql_password
DB_HOST=127.0.0.1:3306
JWT_SECRET=your_jwt_secret
```

Database name `evermos` must already exist in MySQL.

### Docker

```bash
docker compose up --build -d
```

App → `http://localhost:3000` · MySQL → `localhost:3306`

Credentials and `JWT_SECRET` are defined in `docker-compose.yml` — change `JWT_SECRET` before any real deployment.

```bash
docker compose down       # stop
docker compose down -v    # stop + wipe DB volume
```

---

## API Endpoints

All protected routes require a `token` header containing a valid JWT.

### Auth
| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/auth/register` | — | Register a new user (auto-creates a toko) |
| POST | `/auth/login` | — | Login, returns JWT token |

### User & Addresses
| Method | Path | Description |
|---|---|---|
| GET | `/user` | Get my profile |
| PUT | `/user` | Update my profile |
| GET | `/user/alamat` | List my addresses |
| GET | `/user/alamat/:id` | Get address by ID |
| POST | `/user/alamat` | Add address |
| PUT | `/user/alamat/:id` | Update address |
| DELETE | `/user/alamat/:id` | Delete address |

### Toko (Store)
| Method | Path | Description |
|---|---|---|
| GET | `/toko/my` | Get my store |
| GET | `/toko` | List all stores |
| GET | `/toko/:id_toko` | Get store by ID |
| PUT | `/toko/:id_toko` | Update store profile |

### Category (admin only)
| Method | Path | Description |
|---|---|---|
| GET | `/category` | List all categories |
| GET | `/category/:id_category` | Get category by ID |
| POST | `/category` | Create category |
| DELETE | `/category/:id` | Delete category |

### Product
| Method | Path | Description |
|---|---|---|
| GET | `/product` | List products (paginated) |
| GET | `/product/:id` | Get product by ID |
| POST | `/product` | Create product (multipart, supports photo upload) |
| PUT | `/product/:id` | Update product |
| DELETE | `/product/:id` | Delete product |

### Transactions
| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/trx` | Admin | List all transactions (paginated) |
| GET | `/trx/:id` | User | Get transaction by ID |
| POST | `/trx` | User | Create transaction |

### Region (no auth)
| Method | Path | Description |
|---|---|---|
| GET | `/provcity/listprovincies` | List all provinces |
| GET | `/provcity/listcities/:prov_id` | List cities by province |
| GET | `/provcity/detailprovince/:prov_id` | Province detail |
| GET | `/provcity/detailcity/:city_id` | City detail |

---

## Testing

```bash
# Unit tests (no DB required)
go test ./helpers/... ./controllers/...

# Integration tests (requires local PostgreSQL — see tests/.env.test)
cd tests && go test -tags=integration .
```

Unit tests use `go-sqlmock` to mock the database. Integration tests spin up against a real PostgreSQL instance, auto-migrate, seed test data, and tear down on completion.

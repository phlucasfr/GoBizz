# GoBizz - Business Management System

GoBizz is a modern business management system, built with a distributed architecture using microservices. The project consists of a Next.js frontend and multiple backend services communicating via gRPC.

## 🏗️ Project Architecture

The project is organized into the following main services:

### Frontend (`/frontend`)
- Developed with **Next.js 14**
- Modern interface with **Tailwind CSS** and **shadcn/ui**
- Directory structure:
  - `/app` - Application routes and pages
  - `/components` - Reusable components
  - `/context` - React contexts
  - `/api` - API integrations
  - `/utils` - Utility functions
  - `/hooks` - Custom hooks
  - `/lib` - Libraries and utilities
  - `/styles` - Global styles
  - `/public` - Static files

### Authentication Service (`/auth-service`) – Go
- Developed in **Go**
- Handles user authentication and authorization
- Communicates with other services via **gRPC**
- Directory structure:
  - `/cmd` - Application entry points
  - `/internal` - Service internal code
  - `/middleware` - Middlewares
  - `/migrations` - Database migrations
  - `/proto` - gRPC protocol definitions
  - `/utils` - Utility functions

### Links Service (`/links-service`) – Go
- Developed in **Go**
- Manages business links and connections
- Uses **gRPC** for service communication

### Recurring Events Service (`/recurring-service`) – **Rust**
- Developed in **Rust** (Tonic + Prost + SQLx)
- Manages **recurring events** and **occurrence generation**
- Exposes a **gRPC** API consumed by other services (e.g., `auth-service` gateway)
- Key endpoints (gRPC):
  - `CreateEvent`, `GetEvent`, `ListEvents`, `UpdateEvent`, `DeleteEvent`
  - `CutEventFrom` (set `stop_at` to a cutoff date)
  - `ListOccurrences` (compute occurrences between `start`/`end`)
- Environment variables:
  - `DATABASE_URL` – Postgres connection string
  - `PORT` – gRPC port (default: `50053`)
- Example `.env`:
  ```env
  DATABASE_URL=postgres://user:pass@localhost:5432/gobizz_recurring
  PORT=50053
  ```
- Run locally:
  ```bash
  cd recurring-service
  # Run the gRPC server
  cargo run
  ```

## 🚀 Features

### Authentication & Security
- User Registration and Login
- Cryptographic Communication
- JWT-based Authentication
- Protected Routes

### Dashboard
- Overview Statistics (Total Links, Active Links, Expired Links, Total Clicks)
- Data Visualization (Link Performance Charts)
- Most Clicked Links List
- Recent Links List

### Business Links Management
- Create Short Links
- Custom URL Slugs
- Link Expiration Dates
- Edit Existing Links
- Delete Links
- Copy Links to Clipboard
- View Link Analytics
- Sort and Filter Links

### Recurring Events Management
- Create recurring events with `start_date`, `interval_days`, and optional `stop_at`
- Update, delete, or **cut** an event from a specific date
- Compute **occurrences** efficiently for calendar views (server-side)
- Name-based filtering

### System Features
- Toast Notifications
- Responsive Design
- Dark/Light Mode

## 🚀 How to Run

### Prerequisites
- Docker and Docker Compose
- Node.js (LTS)
- Go **1.21+**
- **Rust (nightly)** + Cargo
- PostgreSQL
- Protocol Buffers (handled by build-time vendoring or system `protoc`)
- SQLx CLI (for recurring-service migrations): `cargo install sqlx-cli --no-default-features --features native-tls,postgres`

### Environment Setup

1. Configure environment variables:
   ```bash
   # Frontend
   cp frontend/.env.example frontend/.env

   # Auth Service (Go)
   cp auth-service/.env.example auth-service/.env

   # Links Service (Go)
   cp links-service/.env.example links-service/.env

   # Recurring Service (Rust)
   cp recurring-service/.env.example recurring-service/.env
   ```

2. Start base services using Docker Compose:
   ```bash
   cd project
   make run
   ```

3. Run the frontend:
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

4. Run the authentication service:
   ```bash
   cd auth-service
   make server
   ```

5. Run the links service:
   ```bash
   cd links-service
   make server
   ```

6. Run the recurring events service (Rust):
   ```bash
   cd recurring-service
   # run
   cargo run
   ```
## 🛠️ Technologies Used

### Frontend
- Next.js 14
- TypeScript
- Tailwind CSS
- shadcn/ui
- React Hook Form
- Zod
- Recharts
- React Day Picker
- Sonner
- Framer Motion
- React PDF Renderer
- JWT Decode
- Axios
- Lucide React
- Vaul
- Embla Carousel
- React Resizable Panels

### Backend Services (Go)
- Go 1.23
- gRPC + Protobuf
- PostgreSQL / DynamoDB (Links)
- Redis
- SQLC
- JWT
- Docker
- Fiber
- SendGrid
- Golang Migrate
- Testcontainers

### Recurring Service (Rust)
- Rust (nightly)
- **tonic** (gRPC), **prost**, **prost-types**
- **tonic-health**
- **sqlx** (PostgreSQL + chrono + uuid)
- **rayon** (parallel occurrence computation)
- **tracing**, **tracing-subscriber**
- Docker (multi-stage build)

## 🔐 Security

All application routes between backend and frontend use encryption or service networking isolation to ensure communication security. For this reason, a Postman collection for testing is not currently available.

## 🤝 Contributing

1. Fork the project
2. Create a branch for your feature (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 License

This project is under the MIT license. See the `LICENSE` file for more details.

## 👥 Authors

- Phelipe Lucas - Lead Developer

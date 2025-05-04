# GoBizz - Business Management System

## 📋 About the Project

GoBizz is a modern business management system, built with a distributed architecture using microservices. The project consists of a Next.js frontend and multiple Go backend services communicating via gRPC.

## 🏗️ Project Architecture

The project is organized into the following main services:

### Frontend (`/frontend`)

- Developed with Next.js 14
- Modern interface with Tailwind CSS and Shadcn/ui
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

### Authentication Service (`/auth-service`)

- Developed in Go
- Handles user authentication and authorization
- Uses gRPC for service communication
- Directory structure:
  - `/cmd` - Application entry points
  - `/internal` - Service internal code
  - `/middleware` - Middlewares
  - `/migrations` - Database migrations
  - `/proto` - gRPC protocol definitions
  - `/utils` - Utility functions

### Links Service (`/links-service`)

- Developed in Go
- Manages business links and connections
- Uses gRPC for service communication

### Project Configuration (`/project`)

- `docker-compose.yaml` - Services configuration
- `Makefile` - Automation scripts

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

### System Features

- Toast Notifications
- Responsive Design
- Dark/Light Mode

## 🚀 How to Run

### Prerequisites

- Docker and Docker Compose
- Node.js (LTS version)
- Go 1.21+
- PostgreSQL
- SQLC
- Protocol Buffers (protoc)

### Environment Setup

1. Configure environment variables:

```bash
# In frontend
cp frontend/.env.example frontend/.env

# In auth-service
cp auth-service/.env.example auth-service/.env

# In links-service
cp links-service/.env.example links-service/.env
```

2. Start services using Docker Compose:

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

## 🛠️ Technologies Used

### Frontend

- Next.js 14
- TypeScript
- Tailwind CSS
- Shadcn/ui (Component Library)
- React Hook Form (Form Management)
- Zod (Schema Validation)
- Recharts (Data Visualization)
- React Day Picker (Date Handling)
- Sonner (Toast Notifications)
- Framer Motion (Animations)
- React PDF Renderer (PDF Generation)
- JWT Decode (Authentication)
- Axios (HTTP Requests)
- Lucide React (Icons)
- Vaul (Drawer Component)
- Embla Carousel
- React Resizable Panels

### Backend Services

- Go 1.23
- gRPC (Service Communication)
- Protocol Buffers
- PostgreSQL (Database)
- DynamoDB (Links)
- Redis (Caching)
- SQLC (SQL Query Generation)
- JWT (Authentication)
- Docker (Containerization)
- Fiber (HTTP Framework)
- SendGrid (Email Service)
- Golang Migrate (Database Migrations)
- Testcontainers (Testing)

## 🔐 Security

All application routes between backend and frontend use encryption to ensure communication security. For this reason, a Postman collection for testing is not currently available.

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

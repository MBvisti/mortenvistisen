# AGENT.md

This file provides guidance to Amp Code when working with code in this repository.

## Project Overview

Grafto is a starter template for full-stack Go web applications using server-side rendering. It follows an opinionated approach:
- No ORM, using SQL directly with sqlc
- Server-side HTML rendering with templ
- Session-based authentication
- Email verification and password reset flows
- Vanilla CSS for styling
- Never use inline styles for views, ONLY for emails

The target audience is solo developers building side-projects or bootstrapping new businesses.

## Workflow
Always start by reading the task file (usually task.md if not specified) and create a plan of action to implement/complete the task. This must be presented to the user and accepted before you start to implement it, so the user can make changes to the plan if necessary.

Always run: `go fmt ./...`

## Development Commands

### Templates and Assets
- `just compile-templates` - Compile templ templates
- `just fmt-templates` - Format templ templates

### Database
- `just gdf` to generate Go code after adding/update/delete queries in psql/queries

### Code Quality
- `go fmt ./...` - Format go code

## Architecture

### Core Components

1. **HTTP Server**
   - Based on Echo framework
   - Located in `server/http.go`
   - Handles graceful shutdown

2. **Router**
   - Defined in `router/router.go`
   - Sets up all routes through `SetupRoutes`
   - Uses middleware for sessions, contexts, etc.
   - Routes are organized into groups in `router/routes/`

3. **Handlers**
   - Located in `handlers/`
   - Organized by feature domains (authentication, dashboard, etc.)
   - Handle HTTP requests and map to service operations
   - Use templ for HTML rendering

4. **Models**
   - Located in `models/`
   - Represent domain entities (User, Token)
   - Include validation and business logic
   - Use sqlc-generated code for database operations

5. **Database**
   - PostgreSQL with sqlc for type-safe queries
   - Migrations using goose
   - Model definitions in `psql/migrations/`

6. **Services**
   - Located in `services/`
   - Business logic layer
   - Authentication, registration, etc.

7. **Background Processing**
   - Uses River queue (PostgreSQL-backed)
   - Workers for email sending, etc.
   - Defined in `psql/queue/`

8. **Templates**
   - Uses templ for type-safe templates
   - Views in `views/`
   - Components in `views/internal/components/`
   - Layouts in `views/internal/layouts/`

9. **Configuration**
   - Environment-based configuration in `config/`
   - Uses env variables with `caarlos0/env/v10`

### Request Flow

1. HTTP request → Echo router
2. Router middleware (sessions, context, etc.)
3. Route handler
4. Service/model business logic
5. Database operations
6. HTML rendering with templ
7. HTTP response

### Authentication Flow

1. Registration with email/password
2. Email verification using tokens/code
3. Login creates authenticated session
4. Protected routes check session via middleware

## Key Dependencies

- Echo - Web framework
- templ - Type-safe templates
- sqlc - Type-safe SQL
- River - Background job processing
- gorilla/sessions - Session management
- pgx - PostgreSQL driver

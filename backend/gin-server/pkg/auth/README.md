- [Auth Package](#auth-package)
  - [Overview](#overview)
  - [Request and Logic Flow](#request-and-logic-flow)
    - [Authentication Flow](#authentication-flow)


# Auth Package

Functional Options Pattern

## Overview

The `auth` package provides JWT-based authentication for Go applications, supporting:
- Access and refresh token generation
- Token validation and verification
- Configurable token parameters (expiry times, token head name )

## Request and Logic Flow

### Authentication Flow

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Client    │     │   Server    │     │ Auth Package│
└──────┬──────┘     └──────┬──────┘     └──────┬──────┘
       │                   │                   │
       │  Login Request    │                   │
       │────────────────────>                  │
       │                   │                   │
       │                   │  GenerateToken()  │
       │                   │────────────────────>
       │                   │                   │
       │                   │  Access + Refresh │
       │                   │<────────────────────
       │ Tokens Response   │                   │
       │<─────────────────────                 │
       │                   │                   │
       │  API Request +    │                   │
       │  Access Token     │                   │
       │────────────────────>                  │
       │                   │  ValidateToken()  │
       │                   │────────────────────>
       │                   │                   │
       │                   │  Valid/Invalid    │
       │                   │<────────────────────
       │ API Response      |
       |  (token expired)  │                   │
       │<─────────────────────                 │
       │                   │                   │
       │  Refresh Request  │                   │
       │────────────────────>                  │
       │                   │ValidateRefreshToken()
       │                   │────────────────────>
       │                   │                   │
       │                   │  New Access +     |
       |                   |  Refresh Token    │
       │                   │<────────────────────
       │  New Access Token │                   │
       │<─────────────────────                 │
       │                   │                   │
```

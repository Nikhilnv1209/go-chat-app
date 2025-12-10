# Free Deployment Guide (2025)

This guide walks you through deploying the Go Chat App for **FREE** using **Render** (for the backend) and **Neon** (for the database).

## Prerequisites

1.  **GitHub Account**: You must push this project to a public or private GitHub repository.
2.  **Neon Account**: [Sign up for Neon](https://neon.tech/) (Free Tier).
3.  **Render Account**: [Sign up for Render](https://render.com/) (Free Tier).

---

## Step 1: Push Code to GitHub

If you haven't already, initialize a git repository and push your code:

```bash
git init
git add .
git commit -m "Initial commit"
# Replace with your actual repo URL
git remote add origin https://github.com/YOUR_USERNAME/go-chat-app.git
git push -u origin main
```

---

## Step 2: Set up Free Database (Neon)

1.  Log in to the **Neon Console**.
2.  Click **"New Project"**.
3.  Name it `chat-db` and select your region (e.g., `US East`).
4.  Neon will generate a connection string that looks like:
    `postgres://user:password@ep-cool-frog-123456.us-east-2.aws.neon.tech/neondb?sslmode=require`
5.  **Copy this Connection String**. You will need it for Render.

---

## Step 3: Deploy Backend (Render)

1.  Log in to the **Render Dashboard**.
2.  Click **"New +"** -> **"Web Service"**.
3.  Select **"Build and deploy from a Git repository"**.
4.  Connect your GitHub account and select your `go-chat-app` repository.
5.  **Configure the Service**:
    *   **Name**: `go-chat-backend`
    *   **Region**: Choose the one closest to your Neon DB (e.g., `US East`).
    *   **Branch**: `main`
    *   **Runtime**: `Docker` (Render will automatically detect the `Dockerfile`).
    *   **Instance Type**: `Free` (0.1 CPU, 512MB RAM).

6.  **Environment Variables**:
    Scroll down to "Environment Variables" and add the following:

    | Key | Value |
    |-----|-------|
    | `SERVER_PORT` | `8080` |
    | `JWT_SECRET` | Generate a strong random string (e.g. `openssl rand -hex 32`) |
    | `JWT_EXPIRATION_HOURS` | `24` |
    | `DB_HOST` | The host from your Neon string (e.g., `ep-cool-frog-123456.us-east-2.aws.neon.tech`) |
    | `DB_PORT` | `5432` |
    | `DB_USER` | The user from your Neon string |
    | `DB_PASSWORD` | The password from your Neon string |
    | `DB_NAME` | `neondb` (or whatever database name Neon gave you) |

    *Note: If your code uses a single `DSN` variable, construct it using the values above.*

7.  Click **"Create Web Service"**.

---

## Step 4: Verification

Render will start building your Docker image. This might take a few minutes.

1.  Watch the **Logs** tab in Render.
2.  Once "Live", Render will give you a URL like `https://go-chat-backend.onrender.com`.
3.  Test the health endpoint:
    ```bash
    curl https://go-chat-backend.onrender.com/health
    ```
4.  It should return `{"status": "ok"}`.

---

## Important Notes on Free Tier

*   **Spin Down**: Render's free instances "sleep" after 15 minutes of inactivity. The first request after sleeping might take 30-50 seconds to respond.
*   **Database Limits**: Neon's free tier allows 0.5 GB of storage, which is plenty for a text-based chat app MVP.

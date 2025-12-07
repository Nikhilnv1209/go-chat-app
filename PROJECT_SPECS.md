# Project Specifications: Go Chat Backend

## 1. Project Overview
**Goal**: Build a scalable Go backend for real-time chat.
**Phase 1**: Minimum Viable Product (MVP) - Core messaging and basic groups.
**Phase 2**: Advanced Features (Replies, Edits, Media).

## 2. Technology Stack
*   Go, PostgreSQL, WebSockets.

## 3. High-Level Requirements

| Feature | Status | Notes |
| :--- | :--- | :--- |
| **Auth** | **MVP** | Register, Login (JWT). |
| **Real-time Messaging** | **MVP** | 1-on-1 and Group text messages. |
| **Message History** | **MVP** | Fetch past messages. |
| **Online Status** | **MVP** | Online/Offline indicators. |
| **Read Receipts** | **MVP** | Using `MessageReceipt` table. |

*(See `specs/01_MVP_Feature_Spec.md` for detailed requirements and `specs/02_Technical_Design.md` for schema/API)*

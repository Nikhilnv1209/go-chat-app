# Bug Report
**Project**: Go Chat Application
**Report Period**: December 14, 2025 - December 28, 2025 (Past 2 Weeks)
**Last Updated**: December 28, 2025

---

## Overview
This document tracks all bugs identified and resolved during the post-implementation phase of the chat application. The threshold for this report begins two weeks ago (December 14, 2025), where we transitioned from feature implementation to bug fixes and refinements.

**Note**: Issues before this date are considered part of the initial implementation phase and are not tracked here.

---

## Table of Contents
1. [Backend Bugs](#backend-bugs)
2. [Frontend Bugs](#frontend-bugs)
3. [Summary Statistics](#summary-statistics)

---

## Backend Bugs

### B001: Authentication Error Handling Not Differentiated
**Severity**: Medium
**Status**: ✅ Fixed
**Date Fixed**: December 28, 2025
**Commit ID**: `89f233679aa88c22ad50f0aaf8c76af4d8b635f9`

**Description**:
The authentication system was not properly differentiating between "user not found" and "invalid credentials" scenarios. Both cases returned generic error messages, making it difficult for the frontend to provide appropriate user feedback.

**Root Cause**:
- Generic `ErrNotFound` message was too vague
- Login service did not explicitly return `ErrNotFound` when email does not exist
- No distinction between 404 (user not found) and 401 (invalid credentials)

**Fix Details**:
- Updated `ErrNotFound` message from generic text to 'User not found' for clarity
- Modified Login service to return `ErrNotFound` when email does not exist
- Differentiated HTTP status codes: 404 for user not found, 401 for invalid credentials

**Files Changed**:
- `internal/errors/errors.go`

---

### B002: WebSocket Connection Issues and Presence Broadcasting
**Severity**: High
**Status**: ✅ Fixed
**Date Fixed**: December 28, 2025
**Commit ID**: `7ce2c428826e55ba1b1e5da77e7ba38ee76b19cb`

**Description**:
Real-time presence system had multiple critical issues:
1. No initial presence sync on WebSocket connection
2. Presence broadcasts not properly filtered to user's contacts
3. JSON syntax errors due to message batching in write pump
4. Clients missing initial presence events

**Root Cause**:
- Missing `sendInitialPresence()` call on client connection
- Presence broadcasts sent to all users instead of filtering to contacts
- Write pump batching multiple messages caused malformed JSON
- Race condition between connection and event listener setup

**Fix Details**:
- Added initial presence synchronization on connection
- Implemented proper presence broadcasting to contacts only
- Disabled message batching in write pump to prevent JSON errors
- Added separate presence broadcasting logic in Hub

**Files Changed**:
- `cmd/server/main.go`
- `internal/websocket/client.go`
- `internal/websocket/hub.go`

---

### B003: Conversation Type Case Sensitivity Mismatch
**Severity**: Medium
**Status**: ✅ Fixed
**Date Fixed**: December 16, 2025
**Commit ID**: `3cf8d1c682a3f4f09001748269e9b93c50573365`

**Description**:
The conversation handler was checking for uppercase 'DM'/'GROUP' types while MessageService was creating conversations with lowercase 'private'/'group' types. This mismatch caused empty `target_name` fields in the `GetConversations` response, preventing proper display of group names and user names.

**Root Cause**:
- Inconsistent type string values between handlers and services
- Case-sensitive string comparison in handler logic
- No standardization of conversation type values

**Fix Details**:
- Changed type checks from 'DM'/'GROUP' to 'private'/'group' (lowercase)
- Matched actual values used by MessageService during conversation creation
- Ensured group names and user names are properly fetched and populated

**Files Changed**:
- `internal/handlers/chat_handler.go`

---

### B004: Duplicate Conversation Creation Constraint Error
**Severity**: High
**Status**: ✅ Fixed
**Date Fixed**: December 16, 2025
**Commit ID**: `165d667b8adb9aef87d02ed0d1a43ac28ff0ce8b`

**Description**:
The conversation model lacked a composite unique index, allowing duplicate conversations for the same user-target-type combination. This caused `ON CONFLICT` errors during upsert operations with SQL error `SQLSTATE 42P10` (constraint matching error).

**Root Cause**:
- No database constraint preventing duplicate user-target-type combinations
- Missing composite unique index on `(UserID, Type, TargetID)`
- GORM model definition did not enforce uniqueness

**Fix Details**:
- Added `uniqueIndex:idx_user_conversation` to UserID, Type, and TargetID fields
- Prevents duplicate conversations for same user-target-type combination
- Resolved SQL constraint errors in upsert operations

**Files Changed**:
- `internal/models/conversation.go`

---

### B005: Refresh Token Cookie Configuration Preventing Token Refresh
**Severity**: High
**Status**: ✅ Fixed
**Date Fixed**: December 29, 2025

**Description**:
The refresh token flow was failing with 401 errors because the refresh token cookie was not being sent with `/auth/refresh` requests. The cookie was configured with `SameSite=Strict` and `Path=/auth`, which prevented it from being sent in cross-origin requests (frontend on port 3000, backend on port 8080).

**Root Cause**:
- `SameSite=Strict` blocked cookies from being sent with cross-origin requests
- Frontend (localhost:3000) and backend (localhost:8080) are different origins
- `Path=/auth` was unnecessarily restrictive
- `Secure` flag was hardcoded to `false` instead of being environment-aware

**Fix Details**:
- Changed `SameSite` from `Strict` to `Lax` to allow same-site cross-origin requests
- Changed `Path` from `/auth` to `/` to make cookie available for all endpoints
- Made `Secure` flag environment-aware: `true` in production (GIN_MODE=release), `false` in development
- Tested and verified: refresh flow now works correctly when frontend and backend use same hostname

**Files Changed**:
- `internal/handlers/auth_handler.go`

---

### B006: Sender's Messages Not Syncing Across Multiple Devices in Real-Time
**Severity**: High
**Status**: ✅ Fixed
**Date Fixed**: December 29, 2025

**Description**:
When a user (e.g., Bob) is logged in on multiple devices (laptop and mobile) and sends a message from one device, only that specific device receives the real-time update. The sender's other devices only sync the sent message after a page refresh. However, the recipient (e.g., Alice) correctly receives the message in real-time on all their devices.

**Root Cause**:
- In `SendDirectMessage` and `SendGroupMessage` functions, only the sending client receives the `message_sent` acknowledgment via `client.Send <- ack`
- The message is broadcast to all of the receiver's devices via `hub.SendToUser(receiverID, payload)`
- The sender's other devices are not notified because the acknowledgment is only sent to the current client, not to all of the sender's connected devices
- Missing broadcast to sender's other devices using `hub.SendToUser(senderID, payload)`

**Fix Details**:
- Added `hub.SendToUser(senderID, payload)` after broadcasting to receiver in `SendDirectMessage`
- Added `hub.SendToUser(senderID, senderPayload)` after broadcasting to group members in `SendGroupMessage`
- Both DM and group messages now sync across all sender devices in real-time
- Frontend already handles `new_message` events correctly for both sent and received messages
- No frontend changes required - existing logic properly handles messages from current user

**Files Changed**:
- `internal/service/message_service.go`

---

## Frontend Bugs


### F001: WebSocket Connection and Listener Race Condition
**Severity**: High
**Status**: ✅ Fixed
**Date Fixed**: December 28, 2025
**Commit ID**: `193570e37cf1d75ee508ab2cfef085a6d3d14a70`

**Description**:
Two separate React effects in `useSocketConnection` created a race condition where the WebSocket connection was established before event listeners were attached. This caused the client to miss initial presence events sent by `sendInitialPresence`.

**Root Cause**:
- Separate `useEffect` hooks for connection and listener setup
- No guarantee of execution order between effects
- Socket connected before listeners attached
- Initial presence events lost

**Fix Details**:
- Merged two separate effects into one consolidated effect
- Ensured event listeners are attached before calling `connect()`
- Added debug logging for connection lifecycle troubleshooting
- Eliminated race condition completely

**Files Changed**:
- `frontend/hooks/useSocketConnection.ts`

---

### F002: Conversation State Race Condition (HTTP vs WebSocket)
**Severity**: Medium
**Status**: ✅ Fixed
**Date Fixed**: December 28, 2025
**Commit ID**: `fd6c213c0d2d4b4b4481af0efc14b0d6c254d6ab`

**Description**:
Complex merging logic in the `setConversations` Redux reducer caused race conditions between HTTP-fetched data and WebSocket real-time updates. This led to stale `is_online` status and inconsistent presence display.

**Root Cause**:
- `setConversations` tried to merge with existing state
- Complex logic created potential for stale state
- No clear separation between initial load and real-time updates
- Race condition between HTTP and WebSocket state updates

**Fix Details**:
- Removed complex merging logic from `setConversations` reducer
- Backend `is_online` data now trusted as source of truth on fetch
- `setUserOnlineStatus` handles only real-time WebSocket updates
- Clear separation: HTTP for initial/refresh, WebSocket for real-time

**Files Changed**:
- `frontend/store/features/conversationSlice.ts`

---

### F003: Duplicate WebSocket Connections in React StrictMode
**Severity**: Medium
**Status**: ✅ Fixed
**Date Fixed**: December 27, 2025
**Commit ID**: `55633be197e1a914a500364e81f3b70b777ebfa2`

**Description**:
React's StrictMode in development caused duplicate WebSocket connection attempts. Additionally, the current user's profile showed stale status instead of always displaying "Online".

**Root Cause**:
- StrictMode double-renders components in development
- No connection state guard in socket service
- Profile component fetched status from potentially stale Redux state
- Current user status should always be "Online" (implicit)

**Fix Details**:
- Added connection state guard to prevent duplicate connections
- Force 'Online' status display for current user's profile
- Prevent rendering stale state for authenticated user
- Handle StrictMode double-mounting gracefully

**Files Changed**:
- `frontend/components/chat/UserProfile.tsx`
- `frontend/lib/socketService.ts`

---

### F004: Auth Initialization LocalStorage Crash
**Severity**: High
**Status**: ✅ Fixed
**Date Fixed**: December 27, 2025
**Commit ID**: `266e0f59b0d31e7849d09e5e43d0e9b36f6746a6`

**Description**:
The application crashed during auth initialization when localStorage contained the string `'undefined'` instead of valid JSON or null. This occurred when the app tried to parse credentials from localStorage.

**Root Cause**:
- No error handling for invalid JSON in localStorage
- `setCredentials` could store undefined values as string 'undefined'
- JSON.parse() threw uncaught exception on malformed data

**Fix Details**:
- Added try-catch block for safe JSON parsing in auth initialization
- Added presence check in `setCredentials` to avoid storing undefined values
- Graceful fallback when localStorage contains invalid data
- Prevent app crash on initialization

**Files Changed**:
- `frontend/store/features/authSlice.ts`

---

### F005: Chat UI Layout and Scrolling Failure
**Severity**: High
**Status**: ✅ Fixed
**Date Fixed**: December 27, 2025
**Commit ID**: `d661691891d823d9872c55407a87ec34e881a06e`

**Description**:
The chat interface had critical layout issues:
1. Message list container did not scroll properly
2. Chat header and input shrunk when messages grew
3. Custom ScrollArea component had poor scroll performance
4. Messages displaced header and input elements

**Root Cause**:
- Incorrect flexbox configuration in message container
- Missing `min-height: 0` on flex children
- Chat input had no `flex-shrink-0` protection
- Custom ScrollArea wrapper added unnecessary complexity

**Fix Details**:
- Fixed message list container with flexbox `min-height` for proper scrolling
- Added `flex-shrink-0` to chat input to prevent shrinking
- Replaced custom ScrollArea with native `overflow-auto` for better performance
- Refined message bubbles, avatars, and empty states for polish

**Files Changed**:
- `frontend/app/dashboard/chat/[type]/[targetId]/page.tsx`
- `frontend/components/chat/MessageList.tsx`

---

### F006: Auth Error Display and Redirect Issues
**Severity**: Medium
**Status**: ✅ Fixed
**Date Fixed**: December 16, 2025
**Commit ID**: `751d29342460e0251926c79e17259c441fbd51d1`

**Description**:
Two issues in authentication pages:
1. React error: "Objects are not valid as a React child" when displaying errors
2. Incorrect redirect from `/` to `/dashboard` after successful login

**Root Cause**:
- Attempted to render error object directly in JSX instead of `error.message`
- Backend returns `{error: {code, message}}` format
- Hardcoded redirect to `/` instead of `/dashboard`

**Fix Details**:
- Updated error display to access `error.message` from backend AppError structure
- Changed login/register success redirect from `/` to `/dashboard`
- Ensured proper navigation flow after authentication
- Prevented rendering error objects directly in JSX

**Files Changed**:
- `frontend/app/(auth)/login/page.tsx`
- `frontend/app/(auth)/register/page.tsx`

---

### F007: Mobile Background Scroll Artifacts
**Severity**: Low
**Status**: ✅ Fixed
**Date Fixed**: December 16, 2025
**Commit ID**: `e84da84ee4de4e3fd1ad883337a3b2879403ef2f`

**Description**:
On mobile devices, the background layer caused visual artifacts and layout shifts during scroll. The fixed background position created jarring visual effects on mobile browsers.

**Root Cause**:
- Using viewport height (vh) units on mobile causes issues with browser chrome
- Fixed positioning of background without proper viewport handling
- Base background color not applied to fixed layer
- Standard vh doesn't account for mobile browser UI

**Fix Details**:
- Updated Home, Login, and Register pages to use `lvh` (large viewport height)
- Moved base background color to fixed layer to eliminate visual artifacts
- Adjusted vertical spacing with dynamic viewport units (dvh)
- Better mobile layout stability

**Files Changed**:
- `frontend/app/(auth)/login/page.tsx`
- `frontend/app/(auth)/register/page.tsx`
- `frontend/app/page.tsx`

---

## Summary Statistics

### Backend Bugs
- **Total**: 6
- **Fixed**: 6 (100%)
- **Severity Breakdown**:
  - High: 4 (66.7%)
  - Medium: 2 (33.3%)
  - Low: 0 (0%)

### Frontend Bugs
- **Total**: 7
- **Fixed**: 7 (100%)
- **Severity Breakdown**:
  - High: 3 (42.9%)
  - Medium: 3 (42.9%)
  - Low: 1 (14.2%)

### Overall
- **Total Bugs**: 13
- **Total Fixed**: 13 (100%)
- **Average Resolution Time**: < 1 day
- **Most Common Issue Category**: State Management & Race Conditions (4 bugs)


---

## Lessons Learned

### State Management
- **Clear separation of concerns** between HTTP and WebSocket state updates prevents race conditions
- **Trust backend as source of truth** for initial state loads
- **Real-time updates should be incremental**, not full state replacements

### WebSocket Architecture
- **Event listener setup must happen before connection** to avoid missing events
- **Message batching can corrupt JSON** in WebSocket write pumps
- **React StrictMode requires connection guards** to prevent duplicates

### Database Design
- **Composite unique indexes are critical** for preventing duplicate data
- **Case sensitivity matters** in string-based type fields
- **Normalize type values** across all layers of the application

### Frontend Reliability
- **Always handle localStorage edge cases** (null, undefined, invalid JSON)
- **Mobile viewport units (lvh, dvh) are more stable** than standard vh
- **Error objects cannot be rendered directly** in React JSX
- **Flexbox scrolling requires min-height: 0** on flex children

---

## Next Steps
- Continue monitoring for edge cases in presence system
- Add integration tests for WebSocket connection lifecycle
- Implement error boundary for better error handling
- Add telemetry for tracking future bugs proactively

---

*This bug report is maintained as part of the BMAD documentation framework.*

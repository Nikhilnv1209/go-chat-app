# Frontend Implementation Roadmap

**Status**: Ready
**Owner**: Lead Developer Agent

This document outlines the step-by-step execution plan for the coding agents.

## Phase 1: Foundation & Scaffolding
*   [x] **Task 1.1**: Initialize Next.js 16 project in `frontend/` directory. **[Trace: S-SYS-01]**
    *   `npx create-next-app@latest frontend --typescript --tailwind --eslint`
    *   **UPDATE**: Actually using Next.js 16.0.8 with React 19.2.1 and Tailwind v4
*   [x] **Task 1.2**: Install Dependencies. **[Trace: S-SYS-02]**
    *   **COMPLETED**: `@tanstack/react-query`, `@reduxjs/toolkit`, `react-redux`, `lucide-react`, `clsx`, `axios`, `tailwind-merge`
*   [x] **Task 1.3**: Configure Providers (`app/providers.tsx`). **[Trace: S-SYS-03]**
    *   **COMPLETED**: App wrapped in `ReduxProvider` and `QueryClientProvider`
*   [x] **Task 1.4**: Setup Shadcn/UI. **[Trace: S-SYS-04]**
    *   **COMPLETED**: Shadcn initialized with Slate theme
    *   Components added: Button, Input, Card, Avatar, Badge, ScrollArea, Textarea
    *   `cn()` utility created in `lib/utils.ts`

## Phase 2: Authentication Module
*   [x] **Task 2.1**: Implement `authSlice` (Redux) and `api` utility. **[Trace: F-FAU-02]**
    *   **COMPLETED**: Axios wrapper with JWT handling implemented in `lib/api.ts`
*   [x] **Task 2.2**: Build Login Page (`app/(auth)/login/page.tsx`). **[Trace: F-FAU-02, F-FAU-03]**
    *   **COMPLETED**: Form UI with validation, API integration, and Redux dispatch
*   [x] **Task 2.3**: Build Register Page (`app/(auth)/register/page.tsx`). **[Trace: F-FAU-01, F-FAU-03]**
    *   **COMPLETED**: Registration form with validation and error handling
*   [x] **Task 2.4**: Implement Protected Route Layout (`app/(dashboard)/layout.tsx`). **[Trace: F-FAU-02]**
    *   **COMPLETED**: AuthGuard that redirects to `/login` when not authenticated

## Phase 3: The Dashboard Skeleton
*   [x] **Task 3.1**: Build `ChatSidebar` component. **[Trace: F-FDB-01]**
    *   **COMPLETED**: Fetches conversations using React Query, includes search functionality, unread badges, and online status indicators
    *   **BONUS**: Integrated with `conversationApi` and Redux for state management
*   [x] **Task 3.2**: Build `UserProfile` component (Bottom of sidebar). **[Trace: F-FAU-02]**
    *   **COMPLETED**: UserProfile component created with avatar, user details, online status, and logout functionality
    *   **NOTE**: Actually integrated into ChatSidebar as a mini profile section
*   [x] **Task 3.3**: Create Empty State for Dashboard (`app/(dashboard)/page.tsx`). **[Trace: F-FDB-01]**
    *   **COMPLETED**: "Select a chat" placeholder page with feature tips and navigation
*   [x] **Task 3.4**: Modernize Dashboard Layout. **[Trace: F-FDB-01, F-FDB-04]**
    *   **COMPLETED**: Implement `NavigationRail` component with floating sidebar design.
    *   **COMPLETED**: Refactor `DashboardLayout` for responsive sizing and mobile/desktop switching.

## Phase 4: Core Messaging Features
*   [ ] **Task 4.1**: Build `MessageList` component. **[Trace: F-FCH-01]**
    *   Fetch messages for active chat ID.
    *   Render bubbles (Me vs Them).
*   [ ] **Task 4.2**: Build `ChatInput` component. **[Trace: F-FCH-01]**
    *   Textarea with auto-resize.
    *   Send mutation.
*   [ ] **Task 4.3**: Integrate `SocketService`. **[Trace: F-FCH-01, F-FCH-02]**
    *   Connect on mount.
    *   Listen for `new_message` -> Update Query Cache.

## Phase 5: Refinement & Polish
*   [ ] **Task 5.1**: Implement Infinite Scroll for history. **[Trace: F-FHS-01]**
*   [ ] **Task 5.2**: Add Group Creation Modal. **[Trace: F-FGR-01]**
*   [ ] **Task 5.3**: Add Unread Counts to Sidebar. **[Trace: F-FDB-02]**
*   [ ] **Task 5.4**: Implement Read Receipt Logic & UI. **[Trace: F-FRR-01]**
    *   Show icons (sent/delivered/read) in `MessageBubble`.
    *   Emit `message_delivered` on receive.
    *   Call `POST /read` on view.
*   [ ] **Task 5.5**: Implement Typing Indicators. **[Trace: F-FTI-01]**
    *   Debounce `typing_start` in input.
    *   Show animation in `MessageList`.

---

## Traceability Matrix

Complete mapping from Feature → Story → Task(s).

| Feature ID | Story Summary | Implementing Task(s) | Status |
|------------|---------------|----------------------|--------|
| **F-FAU-01** | Sign Up | Task 2.3 | ✅ **DONE** |
| **F-FAU-02** | Login + Session | Task 2.1, 2.2, 2.4, 3.2 | ✅ **DONE** |
| **F-FAU-03** | Loading States | Task 2.2, 2.3 | ✅ **DONE** |
| **F-FDB-01** | Conversation List | Task 3.1, 3.3 | ✅ **DONE** |
| **F-FDB-02** | Unread Badges | Task 5.3 | 🚧 **IN PROGRESS** (Partially done in sidebar) |
| **F-FDB-03** | Theme Toggle | Task 1.3 (Provider) | ⬜ **TODO** |
| **F-FDB-04** | Modern Layout | Task 3.4 | ✅ **DONE** |
| **F-FCH-01** | Instant Messages | Task 4.1, 4.2, 4.3 | ⬜ **TODO** |
| **F-FCH-02** | Online Status | Task 4.3 | ✅ **DONE** (In sidebar) |
| **F-FHS-01** | Infinite Scroll | Task 5.1 | ⬜ **TODO** |
| **F-FGR-01** | Create Group | Task 5.2 | ⬜ **TODO** |
| **F-FRR-01** | Read Receipts | Task 5.4 | ⬜ **TODO** |
| **F-FTI-01** | Typing Indicators | Task 5.5 | ⬜ **TODO** |
| **S-SYS-01** | Init Project | Task 1.1 | ✅ **DONE** |
| **S-SYS-02** | Install Deps | Task 1.2 | ✅ **DONE** |
| **S-SYS-03** | Configure Providers | Task 1.3 | ✅ **DONE** |
| **S-SYS-04** | Setup Shadcn/UI | Task 1.4 | ✅ **DONE** |


---

## Definition of Done (DoD)

A task is considered **DONE** when:

1.  ✅ **Code Complete**: All acceptance criteria from the PRD story are met.
2.  ✅ **Type Safe**: No TypeScript errors (`tsc --noEmit` passes).
3.  ✅ **Linted**: No ESLint warnings/errors.
4.  ✅ **Tested**: Unit tests pass (if applicable).
5.  ✅ **Responsive**: Works on mobile, tablet, and desktop breakpoints.
6.  ✅ **Accessible**: Keyboard navigable, proper ARIA labels.
7.  ✅ **Documented**: Complex logic has inline comments.
8.  ✅ **Committed**: Changes committed with conventional commit message.

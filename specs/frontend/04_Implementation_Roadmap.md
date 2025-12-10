# Frontend Implementation Roadmap

**Status**: Ready
**Owner**: Lead Developer Agent

This document outlines the step-by-step execution plan for the coding agents.

## Phase 1: Foundation & Scaffolding
*   [ ] **Task 1.1**: Initialize Next.js 15 project in `frontend/` directory. **[Trace: S-SYS-01]**
    *   `npx create-next-app@latest frontend --typescript --tailwind --eslint`
*   [ ] **Task 1.2**: Install Dependencies. **[Trace: S-SYS-02]**
    *   `@tanstack/react-query`, `@reduxjs/toolkit`, `react-redux`, `lucide-react`, `clsx`.
*   [ ] **Task 1.3**: Configure Providers (`app/providers.tsx`). **[Trace: S-SYS-03]**
    *   Wrap app in `ReduxProvider` and `QueryClientProvider`.
*   [ ] **Task 1.4**: Setup Shadcn/UI (Manual or CLI). **[Trace: S-SYS-04]**
    *   Create `cn()` utility.
    *   Add `Button`, `Input`, `Card` components.

## Phase 2: Authentication Module
*   [ ] **Task 2.1**: Implement `authSlice` (Redux) and `api` utility. **[Trace: F-FAU-02]**
    *   Setup Axios/Fetch with JWT handling.
*   [ ] **Task 2.2**: Build Login Page (`app/(auth)/login/page.tsx`). **[Trace: F-FAU-02, F-FAU-03]**
    *   Form UI -> Call API -> Dispatch Login -> Redirect.
*   [ ] **Task 2.3**: Build Register Page (`app/(auth)/register/page.tsx`). **[Trace: F-FAU-01, F-FAU-03]**
*   [ ] **Task 2.4**: Implement Protected Route Layout (`app/(dashboard)/layout.tsx`). **[Trace: F-FAU-02]**
    *   Check for auth token; redirect if missing.

## Phase 3: The Dashboard Skeleton
*   [ ] **Task 3.1**: Build `ChatSidebar` component. **[Trace: F-FDB-01]**
    *   Fetch conversations using `useQuery`.
    *   Render list with active states.
*   [ ] **Task 3.2**: Build `UserProfile` component (Bottom of sidebar). **[Trace: F-FAU-02]**
    *   Logout functionality.
*   [ ] **Task 3.3**: Create Empty State for Dashboard (`app/(dashboard)/page.tsx`). **[Trace: F-FDB-01]**

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
| **F-FAU-01** | Sign Up | Task 2.3 | ⬜ Pending |
| **F-FAU-02** | Login + Session | Task 2.1, 2.2, 2.4, 3.2 | ⬜ Pending |
| **F-FAU-03** | Loading States | Task 2.2, 2.3 | ⬜ Pending |
| **F-FDB-01** | Conversation List | Task 3.1, 3.3 | ⬜ Pending |
| **F-FDB-02** | Unread Badges | Task 5.3 | ⬜ Pending |
| **F-FDB-03** | Theme Toggle | Task 1.3 (Provider) | ⬜ Pending |
| **F-FCH-01** | Instant Messages | Task 4.1, 4.2, 4.3 | ⬜ Pending |
| **F-FCH-02** | Online Status | Task 4.3 | ⬜ Pending |
| **F-FHS-01** | Infinite Scroll | Task 5.1 | ⬜ Pending |
| **F-FGR-01** | Create Group | Task 5.2 | ⬜ Pending |
| **F-FRR-01** | Read Receipts | Task 5.4 | ⬜ Pending |
| **F-FTI-01** | Typing Indicators | Task 5.5 | ⬜ Pending |
| **S-SYS-01** | Init Project | Task 1.1 | ⬜ Pending |
| **S-SYS-02** | Install Deps | Task 1.2 | ⬜ Pending |
| **S-SYS-03** | Configure Providers | Task 1.3 | ⬜ Pending |
| **S-SYS-04** | Setup Shadcn/UI | Task 1.4 | ⬜ Pending |


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

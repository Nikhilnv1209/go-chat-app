# Design Language System (DLS)

**Status**: Approved
**Version**: 1.0
**Owner**: UX/UI Agent

## 1. Design Philosophy
**"Invisible UI"**: The interface should recede, letting the content (the conversation) take center stage.
**Keywords**: crisp, fast, native-feel.

## 2. Color System (Tailwind)

We utilize the **Slate** scale for neutrals to provide a "warmer" tech feel than varying shades of pure gray.

| Context | Light Mode | Dark Mode | logic |
|---|---|---|---|
| **Background** | `bg-white` | `bg-slate-950` | Pure contrast |
| **Surface** | `bg-slate-50` | `bg-slate-900` | Sidebars/Cards |
| **Border** | `border-slate-200` | `border-slate-800` | Subtle separation |
| **Primary** | `indigo-600` | `indigo-500` | Actions/Brand |
| **Text Primary** | `text-slate-900` | `text-slate-50` | High readability |
| **Text Secondary** | `text-slate-500` | `text-slate-400` | Metadata |

## 3. Typography
**Font Family**: `Inter` (Variable).

| Style | Specs | Usage |
|---|---|---|
| **H1** | `text-2xl font-bold tracking-tight` | Page Titles |
| **H2** | `text-lg font-semibold tracking-tight` | Section Headers |
| **Body** | `text-sm font-normal leading-relaxed` | Messages |
| **Caption** | `text-xs font-medium text-slate-500` | Time/Status |

## 4. Components & Interactive States

### Buttons
*   **Base**: `h-10 px-4 py-2 rounded-md font-medium transition-colors focus-visible:ring-2`
*   **Primary**: `bg-primary text-primary-foreground hover:bg-primary/90`
*   **Ghost**: `hover:bg-accent hover:text-accent-foreground` (Used for sidebar items)

### Inputs
*   **Base**: `h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2`

### Message Bubbles
*   **Me**: `bg-indigo-600 text-white rounded-2xl rounded-tr-none`.
*   **Them**: `bg-slate-100 dark:bg-slate-800 text-slate-900 dark:text-slate-100 rounded-2xl rounded-tl-none`.
*   **Status Icons (Me)**:
    *   **Sent**: Single gray check (icon-sm, text-slate-300).
    *   **Delivered**: Double gray check.
    *   **Read**: Double blue check (text-blue-300).
*   **Typing Indicator**:
    *   Three scrolling dots animation.
    *   Container: `bg-slate-50 dark:bg-slate-800/50 rounded-full px-3 py-1`.
    *   Text: "Alice is typing..." (text-xs text-slate-500).

## 5. Motion (CSS)
No JavaScript animations.
*   **Fade In**: `@keyframes enter { from { opacity: 0; transform: translateY(10px); } }`
*   **Duration**: `duration-200` standard.

## 6. Spacing Scale (Tailwind)

Consistent spacing using Tailwind's 4px base unit.

| Token | Value | Usage |
|-------|-------|-------|
| `space-1` | 4px | Icon gaps |
| `space-2` | 8px | Inline element spacing |
| `space-3` | 12px | Small padding |
| `space-4` | 16px | Standard padding |
| `space-6` | 24px | Section margins |
| `space-8` | 32px | Card gaps |

## 7. Responsive Breakpoints

Mobile-first approach using Tailwind's default breakpoints.

| Breakpoint | Min Width | Layout Behavior |
|------------|-----------|-----------------|
| `default` | 0px | Mobile: Sidebar hidden (drawer), Single column |
| `sm` | 640px | Small tablets: Same as mobile |
| `md` | 768px | Tablets: Split view (sidebar + chat) |
| `lg` | 1024px | Desktop: Full sidebar (280px width) |
| `xl` | 1280px | Wide desktop: Centered container, max-width |

### Key Responsive Patterns
*   **Sidebar**: Hidden on mobile (`hidden md:flex`), revealed as drawer via hamburger menu.
*   **Chat Input**: Full width on mobile, constrained on desktop.
*   **Message Bubbles**: Max-width `75%` on desktop, `90%` on mobile.

## 8. Additional Components

| Component | Light Mode | Dark Mode |
|-----------|------------|-----------|
| **Card** | `bg-white border-slate-200 shadow-sm` | `bg-slate-900 border-slate-800` |
| **Modal Overlay** | `bg-black/50` | `bg-black/70` |
| **Avatar** | `rounded-full bg-slate-200` | `rounded-full bg-slate-700` |
| **Badge (Unread)** | `bg-rose-500 text-white text-xs` | Same |

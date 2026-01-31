# Design Language System (DLS)

**Status**: Approved
**Version**: 2.0
**Owner**: UX/UI Agent
**Last Updated**: 2025-12-27

## 1. Design Philosophy
**"Clean & Vibrant"**: A modern, pastel-based design with lavender backgrounds and purple/coral accents that creates a welcoming, friendly atmosphere while maintaining professional clarity.

**Keywords**: clean, modern, approachable, colorful, friendly.

## 2. Color System (Design Palette)

We use a custom color palette based on the approved design reference with Poppins as the primary font.

### Core Colors

| Color Name | Hex Code | HSL | Usage |
|------------|----------|-----|-------|
| **Dark** | `#202022` | `240 2% 13%` | Primary text, dark elements |
| **Purple** | `#7678ed` | `239 79% 70%` | Primary accent, interactive elements, sent messages |
| **Coral** | `#ff7a55` | `13 100% 67%` | Secondary accent, highlights, user avatar |
| **Light** | `#f9fafc` | `220 33% 98%` | Light surfaces, cards |
| **Lavender BG** | `#d8d8ec` | `235 25% 92%` | Main background |
| **Sidebar Dark** | `#202022` | `240 2% 13%` | Navigation rail background |

### Color Application

| Context | Color | CSS Variable | Usage |
|---------|-------|--------------|-------|
| **Background** | Lavender `#d8d8ec` | `--background` | Main app background |
| **Surface** | White `#ffffff` | `--card` | Sidebar, cards, received messages |
| **Border** | Purple 10% opacity | `--border` | Subtle separations |
| **Primary** | Purple `#7678ed` | `--primary` | Buttons, links, active states |
| **Secondary** | Coral `#ff7a55` | `--secondary` | Accents, user elements |
| **Text Primary** | Dark `#202022` | `--foreground` | Main text |
| **Text Secondary** | Dark 50% opacity | `--muted-foreground` | Metadata, timestamps |

## 3. Typography

**Font Family**: `Poppins` (Google Fonts)
**Weights**: 300 (Light), 400 (Regular), 500 (Medium), 600 (SemiBold), 700 (Bold)

| Style | Specs | Usage |
|-------|-------|-------|
| **H1** | `text-2xl sm:text-3xl font-bold` | Page titles, hero headings |
| **H2** | `text-lg font-semibold` | Section headers |
| **Body** | `text-sm font-normal leading-relaxed` | Message content |
| **Small** | `text-xs font-medium` | Sender names, timestamps |
| **Caption** | `text-[10px] font-medium` | Inline timestamps, metadata |

## 4. Components & Interactive States

### Buttons
- **Primary**: `bg-[#7678ed] hover:bg-[#5a5cd9] text-white shadow-lg shadow-[#7678ed]/25`
- **Secondary**: `bg-[#ff7a55] hover:bg-[#e66a47] text-white shadow-lg shadow-[#ff7a55]/25`
- **Ghost**: `hover:bg-[#7678ed]/10 text-[#202022]`
- **Border Radius**: `rounded-xl` (12px)

### Inputs
- **Search Bar**:
  - Background: `bg-[#e8e8f5]` (light lavender)
  - Focus: `bg-[#dcdcf0]` (darker lavender)
  - Border: None
  - Height: `h-11`
  - Border Radius: `rounded-[12px]`
  - Icon: Purple tint `text-[#7678ed]/50`

- **Text Input**:
  - Background: `bg-[#f9fafc]`
  - Border: `border-[#7678ed]/20`
  - Focus: `border-[#7678ed] ring-[#7678ed]/20`
  - Border Radius: `rounded-xl`

### Message Bubbles

#### Grouping Behavior
- Consecutive messages from the same sender are grouped
- Avatar shows at the bottom of the message group
- Sender name shows at the top of the first message in a group (inside bubble)

#### Styling
- **My Messages** (Right-aligned):
  - Background: `bg-[#7678ed]` (Purple)
  - Text: `text-white`
  - Avatar: Coral gradient `from-[#ff7a55] to-[#e66a47]`
  - Last message: `rounded-[18px] rounded-br-[4px]` (tail bottom-right)
  - Earlier messages: `rounded-[18px]` (fully rounded)
  - Padding: `px-3 py-2`
  - Font size: `text-[13px]`

- **Their Messages** (Left-aligned):
  - Background: `bg-white` with `shadow-sm`
  - Text: `text-[#202022]`
  - Avatar: Purple gradient `from-[#7678ed] to-[#5a5cd9]`
  - Sender name: `text-[#7678ed] text-[11px] font-semibold` (inside bubble)
  - Last message: `rounded-[18px] rounded-bl-[4px]` (tail bottom-left)
  - Earlier messages: `rounded-[18px]` (fully rounded)
  - Padding: `px-3 py-2`
  - Font size: `text-[13px]`

- **Timestamp**:
  - Position: Inside bubble, bottom-right
  - Style: `text-[10px] opacity-60`
  - Shows on ALL messages
  - Format: `HH:mm` (24-hour)

### Chat Input
- **Container**:
  - Background: `bg-[#f9fafc]` (light surface)
  - Border: `border-[#7678ed]/10`
  - Border Radius: `rounded-2xl`
  - Focus: `border-[#7678ed]/30`

- **Textarea**:
  - Background: `bg-transparent`
  - No border, no ring
  - Placeholder: `text-[#202022]/40`

- **Send Button**:
  - Background: `bg-[#7678ed]`
  - Hover: `bg-[#5a5cd9]`
  - Shadow: `shadow-md shadow-[#7678ed]/20`
  - Border Radius: `rounded-xl`

### Navigation Rail
- **Background**: `bg-[#202022]` (Dark)
- **Active Item**: `bg-[#7678ed]` with `shadow-lg shadow-[#7678ed]/30`
- **Inactive Item**: `text-white/60 hover:bg-white/10`
- **Logo**: Purple gradient `from-[#7678ed] to-[#5a5cd9]`
- **Profile Avatar**: Coral gradient `from-[#ff7a55] to-[#e66a47]`

### Sidebar (Chat List)

#### Desktop Sidebar
- **Background**: `bg-white`
- **Search**: Lavender background `bg-[#e8e8f5]`
- **Conversation Item**:
  - Hover: `hover:bg-[#7678ed]/5`
  - Active: `bg-[#7678ed]/10`
  - Border: `border-[#7678ed]/5`

#### Mobile Sidebar (Premium Design)
- **Header Background**: `bg-gradient-to-b from-[#7678ed] via-[#6d6fe0] to-[#9ca3af]`
- **Floating Content Card**:
  - Background: `bg-[#f9fafc]`
  - Border Radius: `rounded-t-3xl` (top only)
  - Shadow: `shadow-2xl`
  - Spacing: `mt-2` (small gap from gradient)

- **Date Display**:
  - Format: `DD.MM WEEKDAY` (e.g., "01.02 SUN")
  - Size: `text-5xl font-normal`
  - Color: `text-white`
  - Leading zero padding for day/month

- **Filter Tabs** (WhatsApp-style):
  - Background: `bg-white/70` (inactive), transparent (active)
  - Text: `text-[#202022]/55` (inactive), `text-[#8a8cf5]` (active)
  - Indicator: Bottom underline (`w-6 h-0.5 bg-[#8a8cf5]`) for active tab
  - Padding: `px-2 py-2`
  - Border Radius: `rounded-lg`

- **Conversation List**:
  - Item Padding: `px-4 py-2.5`
  - Separator: `h-px bg-gray-200/60 mx-4` (subtle dividers)
  - Hover: `hover:bg-white/50`
  - Active: `bg-[#8a8cf5]/12`
  - Avatar Size: `w-11 h-11`
  - Avatar Shape: `rounded-full` (circular)

- **Daily Quote Card**:
  - Background: `bg-white/10 backdrop-blur-sm`
  - Border: `border border-white/20`
  - Border Radius: `rounded-2xl`
  - Padding: `p-4`
  - Quote Text: `text-white/95 text-sm italic`
  - Author: `text-white/60 text-xs text-right`

- **Unread Badge**: `bg-[#ff7a55]` (Coral), smaller `min-w-[16px] h-4.5`
- **Folder Badges**:
  - Work: `bg-[#7678ed]/20 text-[#8a8cf5]`
  - Friends: `bg-green-500/20 text-green-600`

- **Avatars**:
  - DM: Purple gradient `from-[#8a8cf5] to-[#7678ed]`
  - Group: Coral gradient `from-[#ff7a55] to-[#e66a47]`

## 5. Motion & Animations

### Keyframes (CSS)
```css
@keyframes float {
  0%, 100% { transform: translateY(0) translateX(0); }
  25% { transform: translateY(-20px) translateX(10px); }
  50% { transform: translateY(-10px) translateX(-10px); }
  75% { transform: translateY(-30px) translateX(5px); }
}

@keyframes enter {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
```

### Transitions
- **Standard**: `transition-colors` (200ms)
- **All Properties**: `transition-all` (300ms)
- **Hover Effects**: Smooth color and background transitions

## 6. Spacing Scale

| Token | Value | Usage |
|-------|-------|-------|
| `gap-0.5` | 2px | Message group spacing |
| `gap-1` | 4px | Tight element spacing |
| `gap-2` | 8px | Icon gaps |
| `gap-2.5` | 10px | Message-avatar gap |
| `gap-3` | 12px | Standard spacing |
| `gap-4` | 16px | Section spacing |
| `p-3` | 12px | Message bubble padding |
| `p-4` | 16px | Container padding |

## 7. Responsive Breakpoints

Mobile-first approach using Tailwind's default breakpoints.

| Breakpoint | Min Width | Layout Behavior |
|------------|-----------|-----------------|
| `default` | 0px | Mobile: Sidebar as drawer, full-width chat |
| `md` | 768px | Tablet: Split view (sidebar + chat) |
| `lg` | 1024px | Desktop: Full sidebar (384px width) |

### Key Responsive Patterns
- **Sidebar**: Drawer on mobile, persistent on desktop
- **Navigation Rail**: Hidden on mobile, visible on desktop
- **Message Bubbles**: Max-width `70%` on desktop, `65%` on mobile
- **Chat Input**: Full width with responsive padding

## 8. Scrollbars

### Custom Scrollbar Styling
- **Track**: `bg-[#c8c8e8]/30` (Light lavender)
- **Thumb**: `bg-[#7678ed]/40` (Purple)
- **Thumb Hover**: `bg-[#7678ed]/60`
- **Thumb Active**: `bg-[#7678ed]/80`
- **Width**: `8px`
- **Border Radius**: `4px`

## 9. Shadows

- **Cards**: `shadow-sm` (subtle)
- **Message Bubbles**: `shadow-sm` for received, `shadow-sm shadow-[#7678ed]/20` for sent
- **Buttons**: `shadow-lg shadow-[color]/25`
- **Navigation Items**: `shadow-lg shadow-[#7678ed]/30` for active

## 10. Additional Components

| Component | Styling |
|-----------|---------|
| **Avatar** | `rounded-full` with gradient backgrounds |
| **Badge (Unread)** | `bg-[#ff7a55] text-white text-xs font-bold` |
| **Card** | `bg-white border-[#7678ed]/10 shadow-sm rounded-xl` |
| **Empty State** | `bg-[#f9fafc]` with purple icon accents |
| **Loading Spinner** | `border-[#7678ed] border-t-transparent` |

## 11. Background Patterns

### Wavy Pattern (SVG)
Used on auth pages and landing page for visual interest:
- Color: Purple `#7678ed` with 12-15% opacity
- Position: Bottom of viewport
- Creates subtle topographic effect

### Floating Orbs
- Purple orb: `rgba(118,120,237,0.25)`
- Coral orb: `rgba(255,122,85,0.2)`
- Animations: `animate-float` and `animate-float-reverse`

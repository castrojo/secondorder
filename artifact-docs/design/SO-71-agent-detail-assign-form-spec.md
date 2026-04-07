# Design Specification: SO-71 Persistent Assignment Form on Agent Detail

## Overview
The Agent Detail page must always provide a way to assign work to an agent, regardless of their current inbox status. Previously, the assignment form was hidden when the inbox was empty, creating a "dead end" for users. This design ensures the form is always visible and the inbox provides clear feedback when empty.

## User Flow
1. User navigates to Agent Detail (`/agents/{slug}`).
2. User sees the Agent's profile, usage metrics, and configuration.
3. User sees the "Assign" form, even if the agent has no currently assigned work.
4. User can select an available issue from the dropdown and assign it.
5. Below the form, the "Inbox" section displays current assignments or a "No issues assigned" empty state.

## UI Components

### 1. Assignment Form
- **Placement:** Positioned above the "Inbox" section for high visibility.
- **Select Dropdown:**
  - **Source:** Should pull from all "Available Issues" (statuses: `todo`, `in_progress`, `in_review`).
  - **Placeholder:** "Assign an issue..."
  - **Option Format:** `[KEY]: [TITLE] (assigned to [ASSIGNEE_NAME])` (if already assigned).
- **Submit Button:**
  - **Label:** "Assign"
  - **Style:** Primary action button (e.g., `bg-ac` background).

### 2. Inbox Section
- **Header:** "Inbox" (always visible).
- **List Container:** Card with border and internal dividers.
- **Empty State:**
  - When no issues are assigned, show: "No issues assigned to this agent's inbox."
  - Style: Centered text, italic, using `text-ink3/50` (muted) style.

## Visual Design Details (Tailwind/CSS)

```html
<!-- Assign Form -->
<form class="flex items-center gap-2">
  <select class="flex-1 bg-sf border border-bd rounded-md px-3 py-1.5 text-xs text-ink outline-none focus:ring-1 focus:ring-ac">
    <option>Assign an issue...</option>
    <!-- Available issues range here -->
  </select>
  <button type="submit" class="px-3 py-1.5 rounded-md text-xs font-medium bg-ac text-ink hover:bg-ac-h transition-colors">
    Assign
  </button>
</form>

<!-- Inbox Section -->
<div>
  <h3 class="text-[13px] font-medium text-ink2 mb-2">Inbox</h3>
  <div class="bg-card border border-bd rounded-lg divide-y divide-sf overflow-hidden">
    <!-- IF .Issues -->
    <!-- List items here -->
    <!-- ELSE -->
    <div class="px-4 py-8 text-center text-ink3/50 text-xs italic">
      No issues assigned to this agent's inbox.
    </div>
    <!-- ENDIF -->
  </div>
</div>
```

## Rationale
- **Discoverability:** Placing the form above the inbox ensures it is seen before any empty state or long list of issues.
- **Clarity:** Always showing the "Inbox" header provides context for what the empty state message refers to.
- **Consistency:** Matches the pattern of other detail pages where actions are persistently available.

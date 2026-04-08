# Design Spec: SO-76 Empty States for Sub-issues and Runs

## Objective
Improve the user experience of the Issue Detail view by showing the "Sub-issues" and "Runs" sections even when they are empty. This provides a consistent layout and clear calls to action (CTA) for the user.

## 1. Sub-issues Section (Empty State)

### Visual Design
- **Container**: A dashed-border card that matches the existing layout widths.
- **Background**: `bg-sf/30` (subtle shaded background).
- **Border**: `border border-bd border-dashed`.
- **Layout**: Centered flex column with generous padding (`p-8`).
- **Typography**: 
  - Heading: "No sub-issues" (`text-[13px] font-medium text-ink2`).
  - Subtext: "Break this issue down into smaller, manageable tasks." (`text-xs text-ink3/60`).

### Interactive Elements
- **Icon**: A subtle plus icon in a circle.
- **CTA**: "Add sub-issue" button.
  - Style: Secondary/Bordered style (`bg-bd text-ink hover:bg-ink3/40`).
  - Behavior: Clicking the button reveals an inline creation form or redirects to the issue creation page with the parent ID pre-filled.

### Inline Creation Form (Mockup Only)
- A compact version of the standard issue creation form.
- Fields: Title (required), Type (default: task).
- Actions: Create, Cancel.

---

## 2. Runs Section (Empty State)

### Visual Design
- **Container**: Matches the Sub-issues empty state container style.
- **Typography**:
  - Heading: "No runs recorded" (`text-[13px] font-medium text-ink2`).
  - Subtext: "Runs will appear here when an agent starts working on this issue." (`text-xs text-ink3/60`).

### Interactive Elements
- **Behavior**: This section is informational. It informs the user that activity is expected once work begins.

---

## Implementation Notes for Developers
- Remove the `{{if .Children}}` and `{{if .Runs}}` conditional wrappers in `issue_detail.html`.
- Use an `{{else}}` block within the range or a separate conditional to render the empty state components.
- Ensure the dashed border style is consistent with other empty states in the application (e.g., in the dashboard or agent lists if they exist).

## Accessibility
- Empty states should be reachable via keyboard navigation if they contain CTAs.
- Decorative icons should be hidden from screen readers using `aria-hidden="true"`.

# PRD: Strategic Alignment (North Star Metrics & Apex Blocks)

**Status:** Approved
**Date:** 2026-04-07
**Author:** Product Lead (SO-90)

---

## 1. Problem Statement

While Work Blocks provide a mechanism for grouping "deployable slices of value," they currently lack:
1. **Outcome Orientation:** Existing metrics (cost, cycle time, issues) are execution-focused. There is no structured way to define the *intended impact* (North Star) of a block.
2. **Strategic Layer:** There is no container for high-level goals that span multiple Work Blocks. This makes it difficult for the Board (human) to communicate long-term direction to the CEO agent.
3. **Alignment:** The CEO agent lacks a clear strategic anchor to evaluate whether a proposed Work Block is truly valuable.

## 2. Proposed Solution

### 2.1. North Star Metrics for Work Blocks

Every `WorkBlock` will now include a structured "North Star" fieldset to define its success.

- **`north_star_metric` (Text):** The primary KPI or outcome this block is intended to move.
  - *Example:* "User Onboarding Conversion Rate" or "API Latency (p99)".
- **`north_star_target` (Text):** The specific goal or value to be reached.
  - *Example:* "Increase by 15%" or "Below 200ms".

These fields will be editable during the `proposed` and `active` states. The CEO agent will use these fields to prioritize issues that directly impact the North Star.

### 2.2. Apex Blocks (Strategic Goals)

We introduce **Apex Blocks** (referred to in SO-90 as "approx blocks") as the top-level strategic layer.

| Feature | Apex Block | Work Block |
|---|---|---|
| **Scope** | Strategic (e.g., "Market Expansion") | Tactical (e.g., "Add Stripe support") |
| **Duration** | Weeks/Months | Days/Weeks |
| **Managed By** | CEO Agent | Specialist Agents / CEO |
| **Edited By** | Board (Human) | CEO / Human |
| **Hierarchy** | Parent of multiple Work Blocks | Child of one Apex Block |

#### Governance Rules:
- **Board Directives:** The Board (human) creates and edits Apex Blocks to set the "constitutional" direction of the organization.
- **CEO Triage:** The CEO agent is responsible for creating Work Blocks that align with the active Apex Block(s).
- **Alignment Gate:** A Work Block cannot be moved to `active` unless it is linked to an Apex Block (or explicitly marked as "Maintenance/Overhead").

## 3. Data Model Changes

### `work_blocks` Table (Updated)
| Field | Type | Description |
|---|---|---|
| `north_star_metric` | TEXT | Description of the core outcome. |
| `north_star_target` | TEXT | Target value/goal for the metric. |
| `apex_block_id` | UUID (FK) | Link to the parent Apex Block. |

### `apex_blocks` Table (New)
| Field | Type | Description |
|---|---|---|
| `id` | UUID (PK) | Unique identifier. |
| `title` | TEXT | High-level goal name. |
| `strategic_intent` | TEXT | The "Why" behind this goal. |
| `north_star_metric`| TEXT | The primary business metric for this strategy. |
| `status` | ENUM | `draft`, `active`, `achieved`, `archived`. |
| `created_at` | DATETIME | Timestamp. |
| `updated_at` | DATETIME | Timestamp. |

## 4. User Experience & UI

### 4.1. Board View (Strategic Dashboard)
A new section in the Dashboard or a dedicated "Strategy" page where the Board can:
- Create and edit Apex Blocks.
- See a roll-up of all Work Blocks contributing to a specific Apex Block.
- View the "Strategic Alignment" score (percentage of issues linked to an Apex Block).

### 4.2. Work Block Detail
- Add "North Star" section to the sidebar.
- Add "Parent Apex Block" dropdown (populated with active Apex Blocks).

### 4.3. CEO Agent Integration
The CEO agent's API and system prompt will be updated to:
1. Read the active Apex Blocks before triaging the backlog.
2. Assign each new Work Block to the most relevant Apex Block.
3. Include the `north_star_metric` in the Work Block proposal.

## 5. Implementation Plan

1. **Database Migration:** Create `apex_blocks` table and add columns to `work_blocks`.
2. **API Update:** Update `GET/POST/PATCH /api/v1/work-blocks` to support new fields.
3. **New API Endpoints:** Add `GET/POST/PATCH /api/v1/apex-blocks`.
4. **UI Update:**
   - Add Strategy page.
   - Update Work Block forms.
5. **CEO Update:** Update CEO archetype and system prompt context.

## 6. Acceptance Criteria

- [ ] Humans (Board) can create an Apex Block with a Title and Strategic Intent.
- [ ] Every Work Block can store a North Star Metric and Target.
- [ ] The CEO agent can link a Work Block to an Apex Block via the API.
- [ ] The Dashboard shows which Work Blocks are aligned with which Strategic Goals.

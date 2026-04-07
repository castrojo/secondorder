# UI/UX Spec: Strategic Alignment (SO-107)

## 1. Overview
Introduce the **Strategy** layer to Second Order, allowing the Board (human) to set high-level goals (**Apex Blocks**) and visualize how **Work Blocks** align with these goals.

## 2. New Components

### 2.1. Strategy Page (`/strategy`)
A new top-level navigation item and page for managing Apex Blocks.

- **Apex Block List**: Cards showing active and completed strategic goals.
- **Apex Block Detail (Roll-up)**: Clicking an Apex Block shows all linked Work Blocks, their status, and aggregated metrics (cost, issues completed).
- **Create Apex Block Form**: Fields: `Title`, `Strategic Intent` (Textarea).
- **Strategic Alignment Score**: A global metric showing the % of issues currently assigned to a Work Block that is linked to an Apex Block.

### 2.2. Updated Work Block Forms
Enhance the Work Block creation and detail views to support strategic alignment.

- **Create/Propose Work Block Form**:
    - **Parent Apex Block**: Searchable dropdown of active Apex Blocks.
    - **North Star Metric**: Text input (e.g., "Active Users", "API Latency").
    - **Target Value**: Text input (e.g., "10,000", " < 200ms").
- **Work Block Detail View**:
    - **Strategic Context Card**: Display the parent Apex Block and North Star Metric prominently.
    - **Alignment Status Badge**: "Aligned" (linked to Apex Block) vs. "Standalone" (unlinked).

### 2.3. Dashboard Updates
Visualize alignment on the main dashboard.

- **Strategic Alignment Card**: A new stat card showing the global alignment score.
- **Apex Progress Overview**: A section showing active Apex Blocks and a progress bar based on linked Work Blocks.

## 3. Visual Language
- Adhere to the existing **Inter** typography and grayscale color palette.
- Use **Indigo/Blue** accents for strategic elements to distinguish from the standard green/emerald "action" colors.
- Maintain the "Zero Human Company" aesthetic: minimal, data-dense, but high-contrast.

## 4. User Flows
1. **Setting Strategy**: Human goes to `/strategy` -> Proposes Apex Block "Scale to 1M Users" -> Sets Intent.
2. **Aligning Work**: CEO/Human proposes Work Block "Optimize DB Indexes" -> Selects "Scale to 1M Users" as Parent -> Sets North Star "Query Time" -> Target "50ms".
3. **Monitoring**: Human views Dashboard -> Sees "Scale to 1M Users" progress -> Clicks through to see contributing Work Blocks.

## 5. Acceptance Criteria
- [ ] Mockup for Strategy page with Apex Blocks list.
- [ ] Mockup for Strategy roll-up (Apex Block detail).
- [ ] Mockup for updated Work Block Propose form.
- [ ] Mockup for updated Dashboard with alignment metrics.

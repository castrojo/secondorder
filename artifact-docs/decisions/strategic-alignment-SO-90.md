# Decision: Strategic Alignment (North Star Metrics & Apex Blocks) (SO-90)

**Date:** 2026-04-07
**Status:** Approved
**Context:** The system needs a way to align tactical Work Blocks with long-term strategic goals (Apex Blocks) and track outcome-based success via North Star metrics.

## Decision
1. **Apex Blocks:** We will introduce a top-level strategic layer. These are managed by the CEO agent but edited by the Board (Human) to set constitutional direction.
2. **North Star Metrics:** Every Work Block will store its intended outcome (e.g., "API Latency") and target (e.g., "< 200ms").
3. **Alignment Gate:** Work Blocks should be linked to an Apex Block to ensure all execution contributes to strategic goals.

## Implementation Plan
- SO-106: Backend database migrations and API endpoints.
- SO-107: UI/UX design for the Strategy dashboard and updated Work Block forms.
- SO-108: Frontend implementation of the Strategy UI and integration.
- SO-109: QA validation of the strategic alignment features.

## Rationale
This structure transitions the system from "execution-only" (tracking runs and issues) to "outcome-oriented" (tracking strategic impact), providing a clear anchor for the CEO agent to prioritize the backlog.

# SO-24: Live-Refresh timeAgo Timestamps in Activity Feed

## Problem

The activity feed rendered `timeAgo` timestamps server-side as static HTML strings (e.g. "2m ago"). Once the page was loaded, these strings never updated — an entry that showed "2m ago" on load still showed "2m ago" an hour later without a full page reload.

## Solution

Two-part fix:

### 1. Go template helper — `iso` function (`internal/templates/templates.go`)

Added `isoTime(t time.Time) string` that returns the timestamp as UTC RFC3339 (ISO 8601):

```go
func isoTime(t time.Time) string {
    return t.UTC().Format(time.RFC3339)
}
```

Registered in `funcMap` as `"iso"` alongside `timeAgo`.

### 2. Template + JS (`internal/templates/activity.html`)

**`activity_entry` template** — timestamp span now carries the machine-readable ISO value:

```html
<!-- Before -->
<span class="text-[11px] text-ink3/50 tabular-nums shrink-0 w-16">{{timeAgo .CreatedAt}}</span>

<!-- After -->
<span class="time-ago text-[11px] text-ink3/50 tabular-nums shrink-0 w-16"
      data-ts="{{iso .CreatedAt}}">{{timeAgo .CreatedAt}}</span>
```

**`<script>` block** (IIFE, placed before `{{template "foot" .}}`):

```js
(function () {
  function timeAgo(isoStr) {
    var d = new Date(isoStr);
    var sec = Math.round((Date.now() - d.getTime()) / 1000);
    if (sec < 5)  return 'just now';
    if (sec < 60) return sec + 's ago';
    var min = Math.round(sec / 60);
    if (min < 60) return min + 'm ago';
    var hr = Math.round(min / 60);
    if (hr < 24)  return hr + 'h ago';
    var day = Math.round(hr / 24);
    if (day === 1) return '1d ago';
    return day + 'd ago';
  }
  function refreshTimestamps() {
    document.querySelectorAll('.time-ago[data-ts]').forEach(function (el) {
      el.textContent = timeAgo(el.dataset.ts);
    });
  }
  setInterval(refreshTimestamps, 30000);
  document.addEventListener('htmx:afterSwap', refreshTimestamps);
})();
```

## Design Decisions

| Decision | Rationale |
|----------|-----------|
| 30s interval | Meets AC2 (≤60s); matches typical "just now" → "1m ago" threshold |
| `data-ts` ISO attribute | Allows JS to recompute without storing state in memory |
| IIFE pattern | Self-contained; no global variables; no memory leak risk |
| `htmx:afterSwap` listener | SSE-injected activity entries are immediately correct without waiting for the next 30s tick |
| JS `timeAgo()` mirrors Go `timeAgo()` | Consistent UX between initial render and live updates |
| Static HTML seed value `{{timeAgo .CreatedAt}}` | Page renders correctly even with JS disabled |

## Acceptance Criteria

| AC | Status | Evidence |
|----|--------|---------|
| AC1: timeAgo values update without page reload | ✅ | `setInterval(refreshTimestamps, 30000)` |
| AC2: Update interval ≤60s | ✅ | 30 000 ms |
| AC3: No memory leaks | ✅ | IIFE; one document-level listener; no per-element closure |
| AC4: Works in secondorder dashboard frontend | ✅ | Targets `activity.html` activity feed |
| AC5: PR filed | ✅ | https://github.com/castrojo/secondorder/pull/4 |

## Files Changed

- `internal/templates/templates.go` — `isoTime` func + `"iso"` funcMap entry
- `internal/templates/activity.html` — `data-ts` attribute + `<script>` refresh block

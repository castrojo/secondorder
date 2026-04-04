# SO-77: Git Diff Line Coloring

## Objective
Add line coloring (+ emerald/green, - red) to the git diff display in Run Detail and Issue Detail views.

## Implementation Details

### CSS Standardized Inheritance
The `td.diff-code` element in `internal/templates/partials.html` had a hardcoded color (`#d4d4d8`). This was changed to `inherit` to allow the parent row (`tr`) color classes to take effect.

```css
/* internal/templates/partials.html */
td.diff-code { color: inherit; }
```

### Line Class Definitions
The `diffLines` function in `internal/templates/templates.go` was updated to provide consistent Tailwind classes for different types of diff lines:

- **Added lines (+)**: `text-emerald-400 bg-emerald-950/30`
- **Removed lines (-)**: `text-red-400 bg-red-950/30`
- **Hunk headers (@@)**: `text-indigo-400 bg-indigo-950/20`
- **File headers (diff --git)**: `text-zinc-600 font-medium`
- **Unchanged lines**: `text-ink2` (matches original `#d4d4d8` brightness)

### Syntax Highlighting Compatibility
The `highlight.js` (hljs) colors are preserved for syntax-highlighted tokens because they are applied to `<span>` tags inside `td.diff-code` with higher specificity than the inherited line color. Non-highlighted text (like identifiers not recognized by hljs or standard text) will now correctly show the emerald/red line color.

## Verification Results
- `+` lines: Emerald text on dark green background.
- `-` lines: Red text on dark red background.
- Gutter symbols (+/-): Still use `text-emerald-500` / `text-red-500` for higher contrast as per existing template logic.
- Unchanged code: Prominent `ink2` gray.

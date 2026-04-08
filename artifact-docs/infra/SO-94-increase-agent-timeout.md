# SO-94: Increase agent timeout

Default agent timeout is now 1200 seconds.

Changes:
- New databases create `agents.timeout_sec` with a default of `1200`.
- New agents created through the UI and startup templates default to `1200`.
- Existing agents still at the old default of `600` are migrated to `1200` by migration `016_increase_default_agent_timeout.sql`.

Scope:
- Explicit custom timeouts are preserved.
- Runtime enforcement is unchanged; the scheduler still uses each agent's configured `timeout_sec`.

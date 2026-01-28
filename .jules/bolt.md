## 2024-05-22 - O(N^2) Validation on Startup
**Learning:** The `CheckRepeatedItems` function was O(N^2) and ran on every command execution via `LoadGPaths`. This scales poorly for users with many paths.
**Action:** Always check complexity of functions in the startup path. Use hash maps for uniqueness checks instead of nested loops.

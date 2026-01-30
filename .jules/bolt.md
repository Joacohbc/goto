## 2026-01-29 - O(N) Slice Indexing
**Learning:** Found a pattern where slice lookups by index were implemented as O(N) loops iterating through the entire slice to match the index, instead of direct O(1) access.
**Action:** Look for `for i, v := range slice { if i == target ... }` patterns in other parts of the codebase.

## 2025-02-23 - Quadratic Validator
**Learning:** The `CheckRepeatedItems` validator ran O(N^2) checks on every startup (via `LoadGPathsFile`), causing noticeable lag for users with many aliases.
**Action:** Use map-based O(N) validation for large configs.

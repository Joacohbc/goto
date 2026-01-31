## 2024-05-22 - Prevent RCE in Self-Update
**Vulnerability:** Execution of downloaded, unverified binary to check version number during self-update.
**Learning:** Relying on binary execution for metadata (version) introduces RCE risk if the download is compromised.
**Prevention:** Fetch metadata from the release API instead of extracting it from the binary.

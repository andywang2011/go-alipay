// Package alipay is a testable Alipay OpenAPI client for mini-programs.
//
// Design constraints:
//   - Configuration via Config (no global singletons).
//   - Config via Config struct (no viper; map keys in the application).
//   - HTTP calls use a shared client with configurable timeout and retry.
//   - Injectable httpDo hook for unit tests (httptest.Server).
//
// Production: construct one API client per process and reuse it.
package alipay

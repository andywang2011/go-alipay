# Changelog

## [Unreleased]

## [0.2.0] - 2026-11-07

### Breaking

- See repository diff for API changes in this release (secure defaults, renamed packages where applicable, context-aware constructors).

### Security

- Hardened defaults: TLS, timeouts, body limits, PII redaction where logging applies.

### Fixed

- Bug fixes from the v0.2.0 hardening pass.

### Changed

- CI uses go-version-file; added race, staticcheck, govulncheck, golangci-lint (gosec).



- Security: with `VerifySign` enabled, responses without a `sign` field are rejected (`ErrMissingSign`).
- Security: the legacy RSA (SHA1) sign type is rejected (`ErrUnsupportedSignType`); only RSA2 is allowed.

## [0.1.0]

- Initial release: Alipay OpenAPI client with config injection, HTTP retry, and response signature verification.
- Injectable HTTP transport for `httptest` and custom gateways.

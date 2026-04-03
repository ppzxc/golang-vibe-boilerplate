---
status: accepted
date: 2026-04-04
decision-makers: ppzxc
---

# 에러 처리 전략: 도메인 sentinel error + RFC 9457 ProblemDetail

## Context and Problem Statement

Go REST API에서 에러 응답 형식이 일관되지 않으면 클라이언트 측 오류 처리가 복잡해진다.
domain/app 레이어에서 HTTP 개념(상태 코드)을 직접 사용하면 레이어 경계가 오염된다.
Java 템플릿의 ProblemDetail + ErrorCode enum 조합에 대응하는 Go 관용적 에러 처리 전략이 필요하다.

## Decision Drivers

* HTTP API 에러 응답의 업계 표준 준수 (RFC 9457)
* domain/app 레이어의 순수 Go 유지 (net/http 의존 없음)
* 클라이언트가 에러 종류를 `errors.Is()` 패턴으로 식별 가능
* Request-Id 헤더와 traceId 필드 연동
* `pkg/problemdetail` 패키지로 HTTP 에러 변환 로직 중앙화

## Considered Options

* A. 도메인 sentinel error + `pkg/problemdetail` (RFC 9457)
* B. 커스텀 error 타입 계층 (errors.As 기반)
* C. HTTP 상태 코드를 포함한 단일 AppError 타입

## Decision Outcome

Chosen option: "A. 도메인 sentinel error + `pkg/problemdetail`", because
Go 관용적 sentinel error 패턴으로 도메인 레이어 순수성을 유지하면서,
httphandler 레이어에서만 RFC 9457 ProblemDetail로 변환하여 레이어 경계를 보호한다.

**에러 계층:**
```
도메인 sentinel error (errors.New)
  → 앱 레이어 래핑 (fmt.Errorf("%w"))
    → httphandler 변환 → RFC 9457 ProblemDetail (application/problem+json)
```

**도메인 에러 정의 위치:** `internal/domain/{name}/errors.go`

**HTTP 에러 매핑:** `internal/adapter/httphandler/error.go`
```go
switch {
case errors.Is(err, domain.ErrNotFound):   → 404
case errors.Is(err, domain.ErrNotAllowed): → 422
default:                                    → 500
}
```

**traceId:** `Request-Id` 요청 헤더 값을 ProblemDetail `traceId` 확장 필드에 포함

### Consequences

* Good, because domain 레이어가 `net/http` 의존 없이 순수 Go를 유지한다
* Good, because `errors.Is()` 체인으로 래핑 후에도 원본 에러 타입 식별 가능
* Good, because RFC 9457 표준 준수로 클라이언트 라이브러리 호환성 확보
* Bad, because 새 도메인 에러 추가 시 httphandler 매핑 테이블도 업데이트 필요
* Bad, because Java ErrorCode enum 대비 에러 목록 중앙 관리가 어렵다

### Confirmation

- domain 패키지 파일에 `net/http` import가 없는지 확인
- ProblemDetail 응답 `Content-Type`이 `application/problem+json`인지 통합 테스트 검증
- 매핑되지 않은 에러가 500으로 처리되는지 테스트

## Pros and Cons of the Options

### B. 커스텀 error 타입 계층

* Good, because 에러에 추가 메타데이터(코드, 파라미터) 포함 가능
* Bad, because 타입 어설션(`errors.As`) 패턴이 sentinel error보다 복잡하다
* Bad, because 레이어 간 커스텀 타입 전파로 의존성 증가 가능

### C. HTTP 상태 코드 포함 AppError

* Good, because httphandler에서 매핑 로직 불필요
* Bad, because domain/app 레이어에 HTTP 개념 침투 — 헥사고날 아키텍처 위반

## More Information

* 관련 규칙: `.claude/rules/error-handling.md`
* 관련 ADR: [ADR-0001](0001-hexagonal-architecture.md) — 도메인 순수성 원칙
* 참고: [RFC 9457 — Problem Details for HTTP APIs](https://www.rfc-editor.org/rfc/rfc9457)

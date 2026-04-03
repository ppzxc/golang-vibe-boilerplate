# Error Handling Rules

## 에러 계층 구조

```
도메인 에러 (sentinel errors)
  → 애플리케이션 에러 (fmt.Errorf("%w"))
    → HTTP 에러 (RFC 9457 ProblemDetail)
```

## 도메인 에러

`internal/domain/{name}/` 패키지에 sentinel error 정의:

```go
var (
    ErrNotFound   = errors.New("todo: not found")
    ErrNotAllowed = errors.New("todo: operation not allowed")
)
```

## 애플리케이션 레이어 에러 처리

- 항상 `fmt.Errorf("operation context: %w", err)` 로 래핑
- 도메인 에러를 위로 전파 (sentinel error 타입 유지)
- 스택 트레이스 없음 — `errors.Is()` / `errors.As()` 로 체크

## HTTP 에러 응답 (RFC 9457 Problem Details)

`pkg/problemdetail` 사용:
```go
// Content-Type: application/problem+json
{
  "type": "https://api.example.com/errors/todo-not-found",
  "title": "Todo Not Found",
  "status": 404,
  "detail": "Todo with id '123' does not exist",
  "instance": "/todos/123",
  "traceId": "<Request-Id 헤더 값>"
}
```

## 에러 매핑 (httphandler)

`internal/adapter/httphandler/error.go`에서 도메인 에러 → HTTP 상태 코드 매핑:

```go
switch {
case errors.Is(err, domain.ErrNotFound):     → 404
case errors.Is(err, domain.ErrNotAllowed):   → 422
default:                                      → 500
}
```

## 금지 패턴

- 스택 트레이스 노출 금지 (`err.Error()` 직접 반환)
- DB 에러 메시지 노출 금지
- HTTP 핸들러에서 도메인 에러 직접 생성 금지
- `panic` 사용 금지 (recovery 미들웨어가 처리)

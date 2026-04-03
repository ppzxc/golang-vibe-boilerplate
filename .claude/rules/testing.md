# Testing Rules

## 테스트 유형

| 유형 | 파일 패턴 | 도구 | 위치 |
|------|----------|------|------|
| 단위 테스트 | `*_test.go` | testing + testify | 각 패키지 |
| 통합 테스트 | `*_integration_test.go` | testcontainers | adapter/postgresrepo |
| E2E 테스트 | `main_integration_test.go` | testcontainers + httptest | cmd/server |

## Build Tags

통합 테스트에 build tag 사용:
```go
//go:build integration
```

실행: `go test -tags=integration ./...`

## testify 사용 규칙

- `assert.*` — 계속 실행 (여러 검증)
- `require.*` — 즉시 중단 (nil 포인터 방지)
- `assert.NoError(t, err)` 패턴 우선

## 테스트 네이밍

```go
func TestService_Create(t *testing.T) { ... }
func TestHandler_CreateTodo(t *testing.T) { ... }
func TestTodoRepository_Save_Integration(t *testing.T) { ... }
```

## 테이블 드리븐 테스트 (단위)

```go
tests := []struct {
    name    string
    input   CreateCommand
    want    *domain.Todo
    wantErr bool
}{...}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) { ... })
}
```

## 금지 패턴

- 실제 DB 없이 DB 테스트 작성 금지 (mock DB 사용 금지, testcontainers 필수)
- `t.Sleep()` 금지 — 타이밍 의존 테스트 금지
- 글로벌 상태 수정 금지 — 각 테스트는 독립적
- 프로덕션 코드에서 테스트 도우미 호출 금지

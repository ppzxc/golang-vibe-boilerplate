---
status: accepted
date: 2026-04-04
decision-makers: ppzxc
---

# Chi v5 라우터 선택

## Context and Problem Statement

Go REST API 보일러플레이트의 HTTP 라우터를 결정해야 한다.
`net/http` 표준 라이브러리와의 호환성, 미들웨어 체이닝, 코론-액션 라우팅(`/todos/{id}:complete`) 지원이 필요하다.
ppzxc RESTful API Guidelines는 커스텀 액션에 콜론 구문을 요구한다.

## Decision Drivers

* `net/http` 표준 `http.Handler` 인터페이스 완전 호환
* 경량 미들웨어 체이닝 (인증, 로깅, 추적 등)
* 콜론-액션 라우팅 지원 (`/todos/{id}:complete`)
* 외부 의존성 최소화
* 활성 유지보수 및 커뮤니티

## Considered Options

* A. Chi v5
* B. `net/http` 표준 라이브러리 (Go 1.22+ ServeMux)
* C. Echo v4
* D. Gin v1

## Decision Outcome

Chosen option: "A. Chi v5", because
`net/http` 완전 호환으로 어댑터 교체 시 핸들러 재작성이 불필요하고,
콜론-액션 라우팅을 네이티브 지원하여 ppzxc RESTful Guidelines의 커스텀 액션 패턴을 구현할 수 있다.

**콜론-액션 라우팅 예시:**
```go
r.Post("/todos/{id}:complete", h.CompleteTodo)
r.Post("/todos/{id}:cancel", h.CancelTodo)
```

### Consequences

* Good, because `net/http` 표준 미들웨어 생태계를 그대로 활용할 수 있다
* Good, because Chi 미들웨어 없이 표준 `http.Handler` 래퍼로 대체 가능하다
* Good, because 콜론-액션 라우팅으로 REST 커스텀 액션을 명시적으로 표현한다
* Bad, because Echo/Gin 대비 기본 제공 기능(바인딩, 유효성 검사)이 적어 직접 구현 필요

### Confirmation

- `internal/adapter/httphandler/` 패키지가 `net/http` 타입만 사용하는지 코드 리뷰
- Chi 라우터에서 콜론-액션 경로가 등록되는지 라우팅 테이블 확인

## Pros and Cons of the Options

### B. `net/http` 표준 라이브러리 (Go 1.22+ ServeMux)

* Good, because 외부 의존성 없음
* Good, because Go 1.22 ServeMux가 경로 파라미터를 네이티브 지원
* Bad, because 콜론-액션 라우팅(`/todos/{id}:complete`) 미지원
* Bad, because 미들웨어 체이닝이 수동 구현 필요

### C. Echo v4

* Good, because 풍부한 기본 기능 (바인딩, 유효성 검사, 렌더링)
* Bad, because 자체 `echo.Context` 타입으로 `net/http` 핸들러와 비호환
* Bad, because 어댑터 교체 시 Echo 컨텍스트 의존성 제거 필요

### D. Gin v1

* Good, because 성능 최적화 및 광범위한 사용 사례
* Bad, because 자체 `gin.Context` 타입으로 `net/http` 비호환
* Bad, because `net/http` 생태계 미들웨어 직접 사용 불가

## More Information

* 관련 규칙: `.claude/rules/architecture.md` (URL path versioning 금지)
* 관련 ADR: [ADR-0007](0007-restful-api-guidelines.md) — 콜론-액션 라우팅 정책
* 참고: [Chi v5 문서](https://go-chi.io/)

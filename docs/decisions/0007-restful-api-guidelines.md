---
status: accepted
date: 2026-04-04
decision-makers: ppzxc
---

# ppzxc RESTful API Guidelines 전면 채택

## Context and Problem Statement

Go REST API 보일러플레이트의 API 설계 규칙을 결정해야 한다.
URL 구조, 버전 관리, JSON 형식, 에러 응답, 페이지네이션, 커스텀 액션 등 일관된 가이드라인이 필요하다.
Java 템플릿과 동일한 ppzxc RESTful API Guidelines를 Go 프로젝트에 적용한다.

## Decision Drivers

* Java 템플릿과 API 설계 일관성 유지
* 업계 표준(RFC 9457, HTTP 의미론) 준수
* 헤더 기반 버전 관리로 URL 클린 유지
* 커서 기반 페이지네이션으로 대용량 데이터 처리
* 콜론-액션 구문으로 RPC 스타일 커스텀 액션 표현

## Considered Options

* A. ppzxc RESTful API Guidelines 전면 채택
* B. Google API Design Guide 채택
* C. 제약 없는 팀 컨벤션

## Decision Outcome

Chosen option: "A. ppzxc RESTful API Guidelines 전면 채택", because
Java 템플릿과 동일한 가이드라인으로 프로젝트 간 일관성을 보장하고,
Chi 라우터가 콜론-액션 라우팅을 네이티브 지원하여 구현 비용이 낮다.

**핵심 결정 사항:**

| 항목 | 결정 |
|------|------|
| URL 버전 관리 | 헤더 기반 (`Api-Version` 헤더), URL path 버전(`/v1/`) 금지 |
| URL 형식 | 케밥케이스 (`/todo-items`), 소문자 |
| JSON 필드 | camelCase (`todoId`, `createdAt`) |
| 에러 형식 | RFC 9457 ProblemDetail (`application/problem+json`) |
| 페이지네이션 | 커서 기반 (`Link` 헤더 + `Total-Count` 헤더) |
| 커스텀 액션 | 콜론 구문 (`POST /todos/{id}:complete`) |
| 컬렉션 응답 | 배열 래퍼 객체 (`{"items": [...], "total": N}`) |

**Chi 라우팅 예시:**
```go
r.Route("/todos", func(r chi.Router) {
    r.Get("/", h.ListTodos)
    r.Post("/", h.CreateTodo)
    r.Get("/{id}", h.GetTodo)
    r.Put("/{id}", h.UpdateTodo)
    r.Delete("/{id}", h.DeleteTodo)
    r.Post("/{id}:complete", h.CompleteTodo)
})
```

### Consequences

* Good, because Java 템플릿과 Go 보일러플레이트 간 API 설계 일관성 확보
* Good, because 커서 기반 페이지네이션으로 OFFSET 방식 대비 대용량 데이터 성능 개선
* Good, because 헤더 기반 버전 관리로 URL이 단순하게 유지된다
* Bad, because 헤더 기반 버전 관리가 브라우저 직접 접근 및 캐싱에 불리할 수 있다
* Bad, because 콜론-액션 구문이 일부 클라이언트 라이브러리와 비호환될 수 있다

### Confirmation

- 라우팅 등록 시 `/v1/`, `/v2/` 패턴 없는지 검토
- JSON 응답 필드가 camelCase인지 통합 테스트 검증
- 에러 응답 `Content-Type: application/problem+json` 확인

## Pros and Cons of the Options

### B. Google API Design Guide

* Good, because 광범위한 채택으로 개발자 친숙도가 높다
* Good, because gRPC-HTTP 트랜스코딩을 고려한 설계
* Bad, because ppzxc 가이드라인과 세부 규칙 충돌 가능 (필드 네이밍 등)
* Bad, because Java 템플릿과 불일치 발생

### C. 팀 컨벤션 (제약 없음)

* Good, because 유연하게 상황별 최적 설계 가능
* Bad, because 프로젝트 간 일관성 없어 코드 리뷰 기준 모호
* Bad, because 보일러플레이트의 가이드 역할을 수행하지 못함

## More Information

* 관련 규칙: `.claude/rules/architecture.md` (URL path versioning 금지)
* 관련 ADR: [ADR-0002](0002-chi-router.md) — Chi 콜론-액션 라우팅
* 관련 ADR: [ADR-0006](0006-error-handling.md) — RFC 9457 ProblemDetail
* 참고: [ppzxc RESTful API Guidelines](https://github.com/ppzxc/restful-api-guidelines)

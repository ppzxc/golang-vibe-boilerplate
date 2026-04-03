---
status: accepted
date: 2026-04-04
decision-makers: ppzxc
---

# sqlc 영속성 계층: SQL→Go 코드 생성 + golang-migrate 마이그레이션

## Context and Problem Statement

Go 보일러플레이트의 PostgreSQL 영속성 계층 기술을 결정해야 한다.
Java 템플릿의 jOOQ(타입세이프 SQL 생성)와 Flyway(마이그레이션)에 대응하는 Go 도구가 필요하다.
도메인 모델의 순수성을 유지하면서 타입세이프한 DB 접근이 요구된다.

## Decision Drivers

* 타입세이프한 SQL 쿼리 (런타임 오류 최소화)
* 도메인 모델에 DB 태그(`db:`) 침투 금지
* SQL을 직접 작성하여 쿼리 가시성 확보
* Java Flyway에 대응하는 마이그레이션 도구
* sqlc 자동 생성 코드는 `internal/adapter/postgresrepo/db/`에 격리

## Considered Options

* A. sqlc + golang-migrate
* B. GORM (ORM)
* C. sqlx (경량 SQL 매퍼)

## Decision Outcome

Chosen option: "A. sqlc + golang-migrate", because
SQL 파일에서 타입세이프한 Go 코드를 자동 생성하여 런타임 오류를 방지하고,
DB 모델이 도메인 모델과 완전히 분리되어 헥사고날 아키텍처 원칙을 준수한다.

**디렉토리 구조:**
```
internal/adapter/postgresrepo/
  db/           ← sqlc 자동 생성 (gitignore 포함)
  queries/      ← SQL 쿼리 파일 (수동 작성)
  migrations/   ← golang-migrate 마이그레이션 파일
  *_repo.go     ← Repository interface 구현 (수동 작성)
```

**sqlc 설정:** `sqlc.yaml` — PostgreSQL, pgx/v5 드라이버 사용

### Consequences

* Good, because SQL 쿼리가 `.sql` 파일에 가시적으로 관리된다
* Good, because 생성된 DB 모델이 도메인 모델과 분리되어 레이어 오염이 없다
* Good, because golang-migrate가 Flyway와 동일한 버전 관리 철학을 제공한다
* Bad, because 쿼리 변경 시 `make sqlc` 재실행이 필요하다
* Bad, because 복잡한 동적 쿼리는 sqlc로 표현하기 어렵다

### Confirmation

- `internal/adapter/postgresrepo/db/`가 `.gitignore`에 포함되는지 확인
- 도메인 모델 파일에 `db:` 태그가 없는지 코드 리뷰
- `make sqlc` 명령으로 코드 생성 검증

## Pros and Cons of the Options

### B. GORM

* Good, because 빠른 CRUD 구현, 자동 마이그레이션 지원
* Bad, because 도메인 모델에 GORM 태그(`gorm:`) 침투
* Bad, because 생성 SQL이 숨겨져 쿼리 최적화가 어렵다
* Bad, because 복잡한 조인 쿼리에서 N+1 문제 발생 가능

### C. sqlx

* Good, because `database/sql` 기반으로 경량하며 SQL 직접 작성 가능
* Bad, because 타입세이프 코드 생성 없음 — 런타임에서 SQL 오류 발견
* Bad, because sqlc 대비 보일러플레이트 스캔 코드 증가

## More Information

* 관련 규칙: `.claude/rules/coding-style.md` (DB 모델 수동 작성 금지)
* 관련 ADR: [ADR-0001](0001-hexagonal-architecture.md) — 도메인 순수성 원칙
* 참고: [sqlc 문서](https://sqlc.dev/), [golang-migrate](https://github.com/golang-migrate/migrate)

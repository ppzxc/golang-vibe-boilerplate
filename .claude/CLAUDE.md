# golang-vibe-boilerplate

Go + Chi + Hexagonal Architecture 기반 REST API 보일러플레이트.
Three Dots Labs 스타일, ppzxc RESTful API Guidelines 적용.

## 모듈

`github.com/ppzxc/golang-vibe-boilerplate`

## 패키지 구조

- `cmd/server/` — 엔트리포인트
- `internal/domain/` — 도메인 모델 + 비즈니스 로직, 아웃바운드 포트
- `internal/app/` — 유스케이스 (Service struct)
- `internal/adapter/httphandler/` — HTTP 핸들러 (Chi)
- `internal/adapter/postgresrepo/` — PostgreSQL 어댑터 (sqlc)
- `pkg/` — 공유 패키지 (problemdetail, pagination, httputil)

## 빠른 참조

- `make run` — 서버 실행
- `make lint` — golangci-lint 실행
- `make fmt` — gofumpt 포맷팅
- `make test` — 단위 테스트
- `make test-integration` — 통합 테스트 (testcontainers)
- `make migrate-up` — DB 마이그레이션
- `make sqlc` — sqlc 코드 생성

## 규칙

| 작업 | 규칙 파일 |
|------|----------|
| 아키텍처 | `.claude/rules/architecture.md` |
| 코딩 스타일 | `.claude/rules/coding-style.md` |
| 테스트 | `.claude/rules/testing.md` |
| CI/도구 | `.claude/rules/ci-tools.md` |
| 에러 처리 | `.claude/rules/error-handling.md` |
| 규칙 관리 | `.claude/rules/rules-maintenance.md` |

정규 규칙 소스는 CLAUDE.md, `.claude/rules/*`, `docs/decisions/*`뿐이다.

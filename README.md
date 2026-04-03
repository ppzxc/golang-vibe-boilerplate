# golang-vibe-boilerplate

Go + Chi + Hexagonal Architecture 기반 REST API 보일러플레이트.

Three Dots Labs 스타일의 Hexagonal Architecture와 ppzxc RESTful API Guidelines를 적용한 프로덕션 수준의 Go 프로젝트 템플릿.
Java Spring Template [`java-spring-template`](https://github.com/ppzxc/java-spring-template)과 동일한 수준의 문서 시스템 및 품질 파이프라인을 Go 관용 방식으로 구현.

## 기술 스택

| 카테고리 | 기술 |
|---------|------|
| 언어 | Go 1.24+ |
| 라우터 | [chi/v5](https://github.com/go-chi/chi) |
| DB 접근 | [sqlc](https://sqlc.dev/) + [database/sql](https://pkg.go.dev/database/sql) |
| DB 드라이버 | [lib/pq](https://github.com/lib/pq) (PostgreSQL) |
| 마이그레이션 | [golang-migrate](https://github.com/golang-migrate/migrate) |
| 로깅 | `slog` (stdlib, Go 1.21+) |
| 설정 | 환경변수 (stdlib `os`) |
| 테스트 | [testify](https://github.com/stretchr/testify) + [testcontainers-go](https://golang.testcontainers.org/) |
| 린트 | [golangci-lint](https://golangci-lint.run/) + [gofumpt](https://github.com/mvdan/gofumpt) |
| Git Hooks | [lefthook](https://github.com/evilmartians/lefthook) |
| CI | GitHub Actions |

## 패키지 구조

```
golang-vibe-boilerplate/
├── cmd/server/          # 엔트리포인트 (수동 DI)
├── internal/
│   ├── domain/todo/     # 도메인 모델 + 비즈니스 로직 + Repository 포트
│   ├── app/todo/        # 유스케이스 Service + DTO
│   ├── adapter/
│   │   ├── httphandler/ # Chi HTTP 핸들러 + 미들웨어
│   │   └── postgresrepo/# PostgreSQL 어댑터 (database/sql)
│   └── config/          # 환경변수 기반 설정
├── pkg/
│   ├── problemdetail/   # RFC 9457 Problem Details
│   ├── pagination/      # Cursor pagination + Link 헤더
│   └── httputil/        # Request-Id 헤더 유틸
├── docs/
│   └── decisions/       # ADR (MADR 4.0 형식)
├── deployments/         # Dockerfile, docker-compose.yml
└── .claude/             # Claude Code 규칙 및 설정
```

## 의존성 방향

```
adapter/httphandler ──→ app ──→ domain
adapter/postgresrepo ──→ domain (Repository interface 구현)
cmd/server ──→ 수동 DI (모든 레이어 와이어링)
```

## 빠른 시작

### 요구사항

- Go 1.24+
- Docker (통합 테스트 및 로컬 개발)
- [golangci-lint](https://golangci-lint.run/usage/install/)
- [gofumpt](https://github.com/mvdan/gofumpt#installation)
- [lefthook](https://github.com/evilmartians/lefthook#install)

### 로컬 개발

```bash
# 의존성 설치
go mod download

# Git hooks 설치
lefthook install

# PostgreSQL 시작
make docker-up

# DB 마이그레이션
make migrate-up

# 서버 실행
make run
```

### 환경 변수

`.env.example`을 `.env`로 복사하여 설정:

```bash
cp .env.example .env
```

| 변수 | 기본값 | 설명 |
|------|--------|------|
| `SERVER_PORT` | `:8080` | 서버 포트 |
| `DATABASE_HOST` | `localhost` | PostgreSQL 호스트 |
| `DATABASE_PORT` | `5432` | PostgreSQL 포트 |
| `DATABASE_USER` | `postgres` | DB 사용자 |
| `DATABASE_PASSWORD` | `postgres` | DB 패스워드 |
| `DATABASE_NAME` | `boilerplate` | DB 이름 |
| `DATABASE_SSLMODE` | `disable` | SSL 모드 |

## API

ppzxc RESTful API Guidelines 적용:
- **URL versioning 없음** — `Api-Version: 2026-04-04` 헤더 사용
- **RFC 9457 Problem Details** 에러 응답
- **camelCase** JSON 필드
- **Cursor-based pagination** — `pageSize`, `pageToken` 쿼리 파라미터
- **Request-Id** 헤더 — 모든 응답에 UUID v4 포함

### Todo API

| Method | Path | 상태 코드 | 설명 |
|--------|------|----------|------|
| `POST` | `/todos` | 201 + Location | Todo 생성 |
| `GET` | `/todos` | 200 + Total-Count | 목록 조회 |
| `GET` | `/todos/{todoId}` | 200 | 단건 조회 |
| `PATCH` | `/todos/{todoId}` | 200 | 부분 수정 |
| `DELETE` | `/todos/{todoId}` | 204 | 삭제 |
| `POST` | `/todos/{todoId}:complete` | 200 | 완료 처리 |

### 예시

```bash
# Todo 생성
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"Buy milk","description":"2% milk"}'

# 목록 조회 (페이지네이션)
curl "http://localhost:8080/todos?pageSize=10"

# 완료 처리 (커스텀 액션)
curl -X POST http://localhost:8080/todos/{id}:complete
```

### 에러 응답 예시

```json
{
  "type": "about:blank",
  "title": "Not Found",
  "status": 404,
  "instance": "/todos/nonexistent",
  "traceId": "550e8400-e29b-41d4-a716-446655440000"
}
```

## 개발 명령어

```bash
make run              # 서버 실행
make build            # 바이너리 빌드
make test             # 단위 테스트 (race detector 포함)
make test-integration # 통합 테스트 (Docker 필요)
make test-all         # 전체 테스트 + 커버리지
make lint             # golangci-lint
make fmt              # gofumpt 포맷팅
make vet              # go vet
make vuln             # govulncheck 보안 스캔
make migrate-up       # DB 마이그레이션 적용
make migrate-down     # DB 마이그레이션 롤백
make sqlc             # sqlc 코드 생성
make docker-up        # PostgreSQL + app 시작
make docker-down      # 컨테이너 중지
```

## 품질 파이프라인

| 단계 | 도구 | Java Spring 대응 |
|------|------|-----------------|
| 포맷팅 | gofumpt | Spotless |
| 스타일 | golangci-lint | Checkstyle + ErrorProne |
| 정적 분석 | go vet + staticcheck | SonarQube |
| 보안 | govulncheck | OWASP Dependency Check |
| 동시성 | go test -race | — |
| Git Hooks | lefthook | lefthook |
| CI | GitHub Actions | GitHub Actions |

### Git Hooks (lefthook)

- **pre-commit**: gofumpt 체크, go vet, golangci-lint (병렬)
- **pre-push**: go test -race, govulncheck

## 새 도메인 추가

1. `internal/domain/{name}/` — 모델 + Repository interface
2. `internal/app/{name}/` — Service + DTO
3. `internal/adapter/postgresrepo/{name}_repo.go` — Repository 구현
4. `internal/adapter/httphandler/{name}_handler.go` — HTTP 핸들러
5. `cmd/server/main.go` — DI 와이어링 추가

자세한 내용은 `.claude/rules/architecture.md` 참조.

## 문서

- [아키텍처 결정 기록 (ADR)](docs/decisions/README.md)
- [설계 스펙](docs/superpowers/specs/2026-04-04-golang-vibe-boilerplate-design.md)
- [Claude Code 규칙](.claude/CLAUDE.md)

## 라이선스

MIT

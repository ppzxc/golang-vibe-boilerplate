# Go Vibe Boilerplate 설계 스펙

## 개요

Go + Chi + Hexagonal Architecture 기반 REST API 보일러플레이트.
Java Spring Template(`java-spring-template`)의 3단계 문서 시스템과 품질 파이프라인을 Go 커뮤니티 표준으로 번역한 프로젝트.

**핵심 원칙:**
- Go 커뮤니티 표준 레이아웃 우선, Hexagonal 개념 차용
- ppzxc RESTful API Guidelines 완전 적용
- `.claude/` + `docs/decisions/` 문서 시스템은 Java 템플릿과 동일

---

## 기술 스택

| 카테고리 | 선택 | 근거 |
|---------|------|------|
| 언어 | Go 1.24+ | 최신 안정 버전 |
| 라우터 | chi/v5 | net/http 호환, 미들웨어 체이닝, 콜론 라우팅 지원 |
| 설정 | viper | 환경변수 + YAML 바인딩 |
| DB 접근 | sqlc | SQL → Go 코드 생성 (Java의 jOOQ 대응) |
| 마이그레이션 | golang-migrate | Flyway 대응 |
| 로깅 | slog (stdlib) | Go 1.21+ 표준 구조화 로깅 |
| 검증 | go-playground/validator | 구조체 태그 기반 |
| 테스트 | testing + testify | assertion 헬퍼 |
| 통합 테스트 | testcontainers-go | PostgreSQL 컨테이너 |
| 린트 | golangci-lint + gofumpt | Java의 Spotless+Checkstyle+ErrorProne 대응 |
| 보안 검사 | govulncheck | 취약점 스캔 |
| Git Hooks | lefthook | Java 템플릿과 동일 도구 |
| CI | GitHub Actions | lint → unit → integration → coverage |
| 컨테이너 | Multi-stage Dockerfile | golang → distroless |

---

## 프로젝트 구조

```
golang-vibe-boilerplate/
├── .claude/
│   ├── CLAUDE.md                         # 프로젝트 진입점 (50줄 이내)
│   ├── settings.local.json               # 권한 설정
│   └── rules/                            # 규칙 파일 (각 100줄 이내)
│       ├── architecture.md               # 의존성 방향, 금지 패턴
│       ├── coding-style.md               # Go 네이밍, error wrapping
│       ├── testing.md                    # 테스트 유형, 네이밍 규칙
│       ├── ci-tools.md                   # 도구 목록, CLI, lefthook, CI
│       ├── error-handling.md             # RFC 9457, 에러 코드 체계
│       └── rules-maintenance.md          # ADR↔Rules 동기화
├── cmd/
│   └── server/
│       └── main.go                       # 엔트리포인트: 설정 로딩 → DI → 서버 시작
├── internal/
│   ├── domain/
│   │   └── todo/
│   │       ├── todo.go                  # 도메인 모델 + 비즈니스 로직 (순수 Go)
│   │       └── repository.go            # 아웃바운드 포트 (interface)
│   ├── app/
│   │   └── todo/
│   │       ├── service.go               # 유스케이스 구현
│   │       └── dto.go                   # CreateCommand, UpdateCommand, ListQuery
│   ├── adapter/
│   │   ├── httphandler/                 # 인바운드 어댑터 (Chi)
│   │   │   ├── router.go               # Chi 라우터 + 미들웨어 구성
│   │   │   ├── todo_handler.go          # Todo HTTP 핸들러
│   │   │   ├── request.go              # 요청 DTO (camelCase JSON)
│   │   │   ├── response.go             # 응답 DTO (camelCase JSON)
│   │   │   ├── middleware.go            # logging, recovery, requestid, versioning
│   │   │   └── error.go                # ProblemDetail 에러 응답 변환
│   │   └── postgresrepo/               # 아웃바운드 어댑터 (sqlc)
│   │       ├── todo_repo.go             # Repository interface 구현
│   │       ├── migrations/              # SQL 마이그레이션 파일
│   │       │   ├── 000001_create_todos.up.sql
│   │       │   └── 000001_create_todos.down.sql
│   │       └── queries/                 # sqlc SQL 파일
│   │           └── todo.sql
│   └── config/                          # 설정 로딩
│       ├── config.go
│       └── config.yaml
├── pkg/
│   ├── problemdetail/                   # RFC 9457 Problem Details 구현
│   │   └── problem.go
│   ├── pagination/                      # Cursor pagination + Link 헤더
│   │   └── cursor.go
│   └── httputil/                        # Request-Id 전파, 공통 HTTP 유틸
│       └── requestid.go
├── docs/
│   ├── decisions/                       # ADR (MADR 4.0 형식)
│   │   ├── README.md                    # ADR 인덱스 테이블
│   │   ├── 0000-template.md             # MADR 4.0 전체 템플릿
│   │   ├── 0000-template-minimal.md     # 최소 템플릿
│   │   ├── 0001-hexagonal-architecture.md
│   │   ├── 0002-chi-router.md
│   │   ├── 0003-sqlc-persistence.md
│   │   ├── 0004-quality-toolchain.md
│   │   ├── 0005-ci-pipeline.md
│   │   ├── 0006-error-handling-rfc9457.md
│   │   └── 0007-restful-api-guidelines.md
│   └── superpowers/
│       └── specs/
├── scripts/                             # 개발 헬퍼 스크립트
├── deployments/
│   ├── Dockerfile                       # Multi-stage: golang → distroless
│   └── docker-compose.yml               # PostgreSQL + app
├── .github/
│   ├── workflows/
│   │   └── test.yml                     # 4-job CI: lint → unit → integration → coverage
│   └── dependabot.yml
├── .golangci.yml                        # golangci-lint 설정
├── lefthook.yml                         # Git hooks
├── sqlc.yaml                            # sqlc 코드 생성 설정
├── Makefile                             # 표준 명령어 (lint, test, run, migrate, sqlc)
├── .env.example                         # 환경변수 템플릿
├── .editorconfig
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

---

## Hexagonal Architecture 규칙

### 의존성 방향 (Three Dots Labs 스타일)

```
adapter/httphandler → app (via service interface)
                       app → domain
adapter/postgresrepo → domain (implements repository interface)
```

### Go 관용적 인터페이스 원칙

- **"Accept interfaces, return structs"** — 인터페이스는 소비자 측에서 정의
- **작은 인터페이스** — 1-2개 메서드가 Go 관용적. 필요에 따라 합성
- **암묵적 구현** — `implements` 키워드 없이 자동 만족
- **수동 DI** — `cmd/server/main.go`에서 명시적 생성자 주입 (DI 프레임워크 없음)

### 금지 패턴

| 금지 | 이유 |
|------|------|
| domain → 외부 패키지 import | 도메인 순수성 (stdlib `errors`, `time`, `context` 만 허용) |
| app → adapter import | 의존성 역전 위반 |
| httphandler → postgresrepo import | 어댑터 간 직접 의존 금지 |
| domain에 `database/sql`, `net/http` 등 | 인프라 의존성 진입 금지 |

### 포트 정의 위치

- **아웃바운드 포트**: `internal/domain/{domain}/repository.go` — 저장소 인터페이스
- **인바운드 포트**: Go 관용에 따라 별도 포트 파일 없이, `app/todo/service.go`의 `Service` struct가 직접 유스케이스를 노출. 어댑터가 필요 시 소비자 측에서 인터페이스 정의

### 트랜잭션 처리 (UpdateFn 패턴)

```go
// domain/todo/repository.go
type Repository interface {
    Save(ctx context.Context, todo *Todo) error
    FindByID(ctx context.Context, id string) (*Todo, error)
    Update(ctx context.Context, id string, fn func(*Todo) (*Todo, error)) error  // 트랜잭션 내 업데이트
    Delete(ctx context.Context, id string) error
}
```

`Update`의 `fn` 클로저 안에서 도메인 로직을 실행하면, 어댑터가 트랜잭션으로 감싸줌.

---

## RESTful API 설계 (ppzxc Guidelines 적용)

### Todo API 엔드포인트

| Method | Path | Status | 설명 |
|--------|------|--------|------|
| POST | `/todos` | 201 + Location | Todo 생성 |
| GET | `/todos` | 200 + `[]` | 목록 (cursor pagination) |
| GET | `/todos/{todoId}` | 200 | 단건 조회 |
| PATCH | `/todos/{todoId}` | 200 | 부분 수정 |
| DELETE | `/todos/{todoId}` | 204 | 삭제 |
| POST | `/todos/{todoId}:complete` | 200 | 완료 처리 (커스텀 액션) |

### 요청/응답 규칙

- **URL**: kebab-case, 복수형 명사, trailing slash 없음
- **JSON**: camelCase 필드명, null 필드 생략
- **날짜**: RFC 3339 UTC (`2026-04-04T12:00:00Z`)
- **표준 필드**: `id`, `createdAt`, `updatedAt`
- **서버는 요청의 unknown 필드 무시**, 클라이언트는 응답의 unknown 필드 무시

### 에러 응답 (RFC 9457 Problem Details)

```json
{
  "type": "https://api.example.com/errors/todo-not-found",
  "title": "Todo Not Found",
  "status": 404,
  "detail": "Todo with id '123' does not exist",
  "instance": "/todos/123",
  "traceId": "550e8400-e29b-41d4-a716-446655440000"
}
```

- `Content-Type: application/problem+json`
- `traceId` = `Request-Id` 헤더 값

### 페이지네이션

- Cursor-based: `pageSize` + `pageToken` 쿼리 파라미터
- 기본 `pageSize`: 20, 최대: 100
- `Total-Count` 응답 헤더
- `Link` 헤더 (RFC 8288): `rel="next"`, `prev`, `first`, `last`
- 빈 컬렉션: `200 OK` + `[]` + `Total-Count: 0`

### 버전 관리

- **URL path versioning 금지** (`/v1/...` 불가)
- `Api-Version: 2026-04-04` 헤더 (ISO 8601)
- 버전 없는 요청은 최신 안정 버전 적용
- 응답에 항상 적용된 버전 포함

### 헤더

- `Request-Id`: 모든 응답에 UUID v4 포함, 서비스 간 전파
- `Content-Type: application/json`
- 커스텀 헤더에 `X-` 접두사 금지

---

## .claude/ 디렉토리 상세

### CLAUDE.md (50줄 이내)

```markdown
# golang-vibe-boilerplate

Go + Chi + Hexagonal Architecture REST API 보일러플레이트.

## 모듈

- `github.com/ppzxc/golang-vibe-boilerplate`

## 패키지 구조

- `cmd/server/` — 엔트리포인트
- `internal/domain/` — 도메인 모델, 아웃바운드 포트
- `internal/app/` — 유스케이스 (Three Dots Labs 스타일)
- `internal/adapter/httphandler/` — HTTP 핸들러 (Chi)
- `internal/adapter/postgresrepo/` — PostgreSQL 어댑터 (sqlc)

## 빠른 참조

- `make lint` — golangci-lint 실행
- `make fmt` — gofumpt 포맷팅
- `make test` — 단위 테스트
- `make test-integration` — 통합 테스트 (testcontainers)
- `make run` — 서버 실행
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
```

### settings.local.json

```json
{
  "permissions": {
    "allow": [
      "Bash(git add:*)",
      "Bash(git commit:*)",
      "Bash(gh issue:*)",
      "Bash(gh pr:*)",
      "Bash(make:*)",
      "Bash(go:*)"
    ]
  }
}
```

### rules/ 파일 요약

| 파일 | 핵심 내용 |
|------|----------|
| `architecture.md` | 의존성 방향 (adapter→app→domain), 금지 패턴, 포트 위치, Three Dots Labs 스타일 |
| `coding-style.md` | Go 네이밍 (Effective Go 준수), error wrapping (`fmt.Errorf("%w")`), 구조체 태그 camelCase |
| `testing.md` | `*_test.go` 단위, `*_integration_test.go` 통합 (build tag), testcontainers 규칙, testify 사용 |
| `ci-tools.md` | golangci-lint, gofumpt, govulncheck, go vet, go test -race, lefthook, GitHub Actions |
| `error-handling.md` | 도메인 에러 타입, ProblemDetail 변환, 에러 코드 체계, 스택 트레이스 노출 금지 |
| `rules-maintenance.md` | ADR↔Rules 동기화, MADR 4.0, 규칙 파일 100줄 제한, CLAUDE.md 50줄 제한 |

---

## docs/decisions/ 상세

### ADR 형식

Java 템플릿과 **동일한 MADR 4.0** 형식:
- frontmatter: status, date, decision-makers
- 섹션: Context, Decision Drivers, Considered Options, Decision Outcome, Pros/Cons

### 초기 ADR 목록

| # | 제목 | 핵심 결정 |
|---|------|----------|
| 0001 | Hexagonal Architecture | 패키지 기반 레이어 분리, internal/ 사용 |
| 0002 | Chi Router | net/http 호환성, 콜론 라우팅 지원 |
| 0003 | sqlc Persistence | SQL 우선 접근, golang-migrate 마이그레이션 |
| 0004 | Quality Toolchain | golangci-lint + gofumpt + govulncheck + go vet + race detector |
| 0005 | CI Pipeline | Lefthook (pre-commit/pre-push) + GitHub Actions 4-job |
| 0006 | Error Handling | RFC 9457 ProblemDetail, 도메인 에러 코드 enum |
| 0007 | RESTful API Guidelines | ppzxc guidelines 전면 적용, header versioning |

---

## 품질 파이프라인

### Java 템플릿 대응표

| Java | Go | 역할 |
|------|-----|------|
| Spotless (Google Java Format) | gofumpt | 코드 포맷팅 |
| Checkstyle | golangci-lint (stylecheck, revive) | 스타일 검사 |
| ErrorProne | go vet + golangci-lint (staticcheck) | 정적 분석 |
| NullAway | golangci-lint (nilnil, nilerr) | nil 안전성 |
| ArchUnit | 테스트 코드에서 import 검증 | 아키텍처 규칙 |
| OpenRewrite | — (Go는 자동 리팩토링 도구 없음) | — |
| JaCoCo | go test -coverprofile | 커버리지 |

### Lefthook 설정

```yaml
pre-commit:
  parallel: true
  commands:
    fmt-check:
      glob: "*.go"
      run: gofumpt -l -d .
    lint:
      glob: "*.go"
      run: golangci-lint run
    vet:
      run: go vet ./...

pre-push:
  commands:
    test:
      run: go test -race ./...
    vuln:
      run: govulncheck ./...
```

### GitHub Actions (4-job)

| Job | 트리거 | 내용 |
|-----|-------|------|
| `lint` | push/PR to main | gofumpt check, golangci-lint, go vet |
| `unit-test` | push/PR to main | go test -race (adapter/output 제외) |
| `integration-test` | push/PR to main | go test -tags=integration (testcontainers) |
| `coverage` | unit+integration 후 | go test -coverprofile, 리포트 생성 |

---

## Todo CRUD 샘플 상세

### Domain Layer (Three Dots Labs 스타일)

```go
// internal/domain/todo/todo.go — 모델 + 비즈니스 로직을 한 파일에
type Todo struct {
    ID          string
    Title       string
    Description string
    Completed   bool
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

func NewTodo(title, description string) *Todo { ... }
func (t *Todo) Complete() error { ... }           // 비즈니스 로직
func (t *Todo) UpdateTitle(title string) error { ... }
```

```go
// internal/domain/todo/repository.go — 아웃바운드 포트
type Repository interface {
    Save(ctx context.Context, todo *Todo) error
    FindByID(ctx context.Context, id string) (*Todo, error)
    FindAll(ctx context.Context, cursor string, pageSize int) ([]*Todo, string, error)
    Update(ctx context.Context, id string, fn func(*Todo) (*Todo, error)) error
    Delete(ctx context.Context, id string) error
    Count(ctx context.Context) (int64, error)
}
```

### Application Layer

```go
// internal/app/todo/service.go — 유스케이스 구현 (별도 인터페이스 없이 struct 직접 노출)
type Service struct {
    repo domain.Repository
}

func NewService(repo domain.Repository) *Service { ... }

func (s *Service) Create(ctx context.Context, cmd CreateCommand) (*domain.Todo, error) { ... }
func (s *Service) FindByID(ctx context.Context, id string) (*domain.Todo, error) { ... }
func (s *Service) FindAll(ctx context.Context, query ListQuery) (*PageResult, error) { ... }
func (s *Service) Update(ctx context.Context, id string, cmd UpdateCommand) (*domain.Todo, error) { ... }
func (s *Service) Delete(ctx context.Context, id string) error { ... }
func (s *Service) Complete(ctx context.Context, id string) (*domain.Todo, error) { ... }
```

### Adapter Layer

- `httphandler`는 `app.Service` struct를 직접 받거나, 필요 시 소비자 측 인터페이스를 정의
- `postgresrepo`는 `domain.Repository` 인터페이스를 구현

---

## DI & 서버 구성

`cmd/server/main.go`에서 수동 DI — Go 관용적 명시적 와이어링:

```go
func main() {
    cfg := config.Load()
    db := postgresrepo.Connect(cfg.Database)

    // Outbound adapters
    todoRepo := postgresrepo.NewTodoRepository(db)

    // Application services
    todoService := todo.NewService(todoRepo)

    // Inbound adapters + router
    router := httphandler.NewRouter(todoService)

    // Graceful shutdown
    srv := &http.Server{Addr: cfg.Server.Addr, Handler: router}
    go func() { srv.ListenAndServe() }()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    srv.Shutdown(ctx)
}
```

---

## Makefile 명령어

| 명령어 | 설명 |
|-------|------|
| `make run` | 서버 실행 |
| `make build` | 바이너리 빌드 |
| `make test` | 단위 테스트 |
| `make test-integration` | 통합 테스트 (testcontainers) |
| `make test-all` | 전체 테스트 |
| `make lint` | golangci-lint 실행 |
| `make fmt` | gofumpt 포맷팅 |
| `make vet` | go vet 실행 |
| `make vuln` | govulncheck 실행 |
| `make migrate-up` | DB 마이그레이션 적용 |
| `make migrate-down` | DB 마이그레이션 롤백 |
| `make migrate-create` | 새 마이그레이션 파일 생성 |
| `make sqlc` | sqlc 코드 생성 |
| `make docker-build` | Docker 이미지 빌드 |
| `make docker-up` | docker-compose up |

---

## Dockerfile (Multi-stage)

```dockerfile
# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server ./cmd/server

# Run stage
FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/server /server
COPY --from=builder /app/internal/adapter/output/postgres/migrations /migrations
ENTRYPOINT ["/server"]
```

---

## 검증 방법

1. `make lint` — 린트 통과 확인
2. `make test` — 단위 테스트 통과
3. `make test-integration` — PostgreSQL testcontainers 통합 테스트 통과
4. `make docker-up` — docker-compose로 전체 스택 기동
5. Todo CRUD API 호출 테스트:
   - `POST /todos` → 201 + Location
   - `GET /todos` → 200 + `[]` + Link + Total-Count
   - `GET /todos/{id}` → 200
   - `PATCH /todos/{id}` → 200
   - `POST /todos/{id}:complete` → 200
   - `DELETE /todos/{id}` → 204
6. 에러 응답 확인: `GET /todos/nonexistent` → 404 ProblemDetail
7. `Request-Id` 헤더 전파 확인

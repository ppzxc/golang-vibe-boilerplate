# golang-vibe-boilerplate

Go + Chi + Hexagonal Architecture 기반 REST API 보일러플레이트.

Three Dots Labs 스타일의 Hexagonal Architecture와 ppzxc RESTful API Guidelines를 적용한 프로덕션 수준의 Go 프로젝트 템플릿.
2026년 업계 표준(OpenAPI 3.0.3, DDD Aggregate, Dependency Injection)을 반영하여 설계되었습니다.

## 기술 스택

| 카테고리 | 기술 | 비고 |
|---------|------|------|
| 언어 | Go 1.24+ | 최신 안정 버전 |
| 라우터 | [chi/v5](https://github.com/go-chi/chi) | net/http 호환 |
| API 규약 | [OpenAPI 3.0.3](https://swagger.io/specification/) | [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) 자동 생성 |
| DI | [Google Wire](https://github.com/google/wire) | 컴파일 타임 의존성 주입 |
| DB 접근 | [sqlc](https://sqlc.dev/) + [database/sql](https://pkg.go.dev/database/sql) | 타입 안전한 SQL |
| DB 드라이버 | [lib/pq](https://github.com/lib/pq) (PostgreSQL) | |
| 마이그레이션 | [golang-migrate](https://github.com/golang-migrate/migrate) | |
| 로깅 | `slog` (stdlib, Go 1.21+) | 구조화 로깅 |
| 관측 가능성 | 도메인 이벤트 (Aggregate 내장) | InMemory Event Bus |
| 문서화 | Swagger UI (`/docs`) | |

## 패키지 구조 (2026 Industry Standard)

```
golang-vibe-boilerplate/
├── api/                 # OpenAPI 3.0.3 스펙 (SSOT)
├── cmd/server/          # 엔트리포인트 (di.InitializeServer 사용)
├── internal/
│   ├── domain/todo/     # Aggregate (Events 포함) + Repository 포트
│   ├── app/todo/        # 유스케이스 Service (Event 발행)
│   ├── adapter/
│   │   ├── httphandler/ # Inbound Adapter (Generated ServerInterface 구현)
│   │   ├── postgresrepo/# Outbound Adapter (Repository Port 구현)
│   │   └── eventbus/    # Event Dispatcher (InMemory)
│   ├── di/              # Google Wire 설정 및 생성 코드
│   └── config/          # 환경변수 기반 설정
├── pkg/
│   ├── problemdetail/   # RFC 9457 Problem Details
│   ├── pagination/      # Cursor pagination + Link 헤더
│   └── httputil/        # Request-Id 헤더 유틸
├── docs/
│   └── decisions/       # ADR (MADR 4.0 형식)
└── deployments/         # Dockerfile, docker-compose.yml
```

## 아키텍처 원칙

### 1. Hexagonal & Strict Ports
- **Inbound Adapter (HTTP)**: `httphandler` 패키지가 자신이 사용할 `TodoService` 인터페이스(Port)를 직접 소유합니다.
- **Outbound Port**: 도메인 레이어에서 인터페이스를 정의하고, 인프라 레이어(`postgresrepo`)가 이를 구현합니다.
- **의존성 역전**: 모든 의존성은 도메인을 향합니다.

### 2. DDD (Domain-Driven Design)
- **Aggregate**: `Todo` 객체는 비즈니스 불변성을 유지하며, 상태 변경 시 도메인 이벤트를 생성하여 내부에 보관합니다.
- **Domain Events**: `Complete()`, `Update()` 등의 액션은 이벤트를 발생시키며, 서비스 레이어에서 이를 발행(Publish)합니다.

## 개발 명령어

```bash
make install-tools    # oapi-codegen, wire 설치
make generate         # OpenAPI 및 Wire 코드 생성
make run              # 서버 실행
make build            # 바이너리 빌드
make test             # 전체 테스트
make docker-up        # 인프라 기동
```

## API 문서

서버 실행 후 브라우저에서 접속:
- **Swagger UI**: [http://localhost:8080/docs](http://localhost:8080/docs)
- **OpenAPI Spec**: [http://localhost:8080/openapi.yaml](http://localhost:8080/openapi.yaml)

## 라이선스

MIT

# Architecture Rules

## Three Dots Labs 스타일 Hexagonal Architecture

### 의존성 방향

```
adapter/httphandler → app → domain
adapter/postgresrepo → domain (Repository interface 구현)
cmd/server → 수동 DI (모든 레이어 와이어링)
```

### Go 관용적 인터페이스 원칙

- **"Accept interfaces, return structs"** — 인터페이스는 소비자 측에서 정의
- **작은 인터페이스** — 1-2개 메서드가 Go 관용적. 필요 시 합성
- **암묵적 구현** — `implements` 키워드 없이 자동 만족
- **수동 DI** — DI 프레임워크 없이 `cmd/server/main.go`에서 명시적 생성자 주입

### 포트 정의 위치

- **아웃바운드 포트**: `internal/domain/{domain}/repository.go` — 저장소 인터페이스
- **인바운드**: 별도 포트 파일 없음. `app/{domain}/service.go`의 Service struct를 어댑터가 직접 사용

### 트랜잭션: UpdateFn 패턴

Repository `Update` 메서드는 클로저를 받아 트랜잭션으로 감쌈:
```go
Update(ctx context.Context, id string, fn func(*Todo) (*Todo, error)) error
```

### 금지 패턴

| 금지 | 이유 |
|------|------|
| `domain` → 외부 패키지 import | 도메인 순수성. stdlib(`errors`,`time`,`context`) 만 허용 |
| `app` → `adapter` import | 의존성 역전 위반 |
| `httphandler` → `postgresrepo` import | 어댑터 간 직접 의존 금지 |
| `domain`에 `database/sql`, `net/http` | 인프라 의존성 진입 금지 |
| DI 프레임워크 (wire, fx) | 수동 DI가 명시적이고 이해하기 쉬움 |
| URL path versioning (`/v1/...`) | ppzxc RESTful Guidelines 위반 |

### 새 도메인 추가 절차

1. `internal/domain/{name}/` — 모델 + Repository interface
2. `internal/app/{name}/` — Service + DTO
3. `internal/adapter/postgresrepo/{name}_repo.go` — Repository 구현
4. `internal/adapter/httphandler/{name}_handler.go` — HTTP 핸들러
5. `cmd/server/main.go` — DI 와이어링 추가

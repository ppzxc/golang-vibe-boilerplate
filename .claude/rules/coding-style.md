# Coding Style Rules

## Go 네이밍 규칙 (Effective Go 기준)

| 위치 | 패턴 | 예시 |
|------|------|------|
| 인터페이스 | 동사+er 또는 명사 | `Repository`, `Reader` |
| Service struct | `Service` | `todo.Service` |
| DTO | 명사+Command/Query | `CreateCommand`, `ListQuery` |
| 핸들러 | `Handler` | `TodoHandler` |
| 생성자 | `New*` | `NewService`, `NewHandler` |

## 핵심 코딩 규칙

- **error wrapping**: `fmt.Errorf("context: %w", err)` 사용
- **nil 반환**: 에러 시 zero value + error 반환. `(nil, nil)` 금지
- **컨텍스트**: 함수 첫 번째 파라미터는 항상 `ctx context.Context`
- **로깅**: `slog` 표준 라이브러리만 사용. `fmt.Print*` 로깅 금지
- **초기화**: 생성자 함수(`New*`)에서만 초기화. 전역 변수 금지

## JSON 태그 규칙 (ppzxc RESTful Guidelines)

- camelCase 필드명: `json:"todoId"`, `json:"createdAt"`
- null 필드 생략: `json:"description,omitempty"`
- 읽기 전용 필드: `json:"id" validate:"-"`

## 구조체 설계

- **도메인 모델**: 순수 Go struct. 외부 태그 없음 (`json:`, `db:` 금지)
- **HTTP DTO**: JSON 태그 포함. 도메인 모델과 분리
- **DB 모델**: sqlc 자동 생성. 수동 작성 금지

## 패키지 설계

- 순환 의존성 금지
- `internal/` — 외부 모듈 import 방지
- `pkg/` — 외부 공개 가능한 유틸리티
- 패키지명은 단수: `todo` (not `todos`), `httphandler` (not `handlers`)

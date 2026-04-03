---
status: accepted
date: 2026-04-04
decision-makers: ppzxc
---

# Hexagonal Architecture 적용: Three Dots Labs 스타일 패키지 기반 레이어 분리

## Context and Problem Statement

Go 단일 모듈 보일러플레이트의 아키텍처 패턴을 결정해야 한다.
비즈니스 로직이 HTTP/DB 인프라에 결합되면 테스트가 어렵고, 인프라 교체 시 비즈니스 로직까지 수정해야 한다.
Go 생태계는 Java와 달리 DI 프레임워크나 모듈 시스템 없이 관용적으로 의존성을 관리해야 한다.

## Decision Drivers

* 비즈니스 로직의 인프라 독립성
* `net/http` 호환성을 유지하면서 어댑터 교체 가능
* Go 관용적 인터페이스 패턴 ("accept interfaces, return structs")
* 수동 DI로 명시적 의존성 관리
* 단일 모듈 내 `internal/` 패키지로 외부 접근 차단

## Considered Options

* A. Three Dots Labs 스타일 패키지 기반 Hexagonal (단일 모듈)
* B. Go workspace 멀티모듈 (domain/app/adapter 각 모듈 분리)
* C. 플랫 구조 (Controller → Service → Repository, 단순 레이어드)

## Decision Outcome

Chosen option: "A. Three Dots Labs 스타일 패키지 기반 Hexagonal", because
단일 모듈 내 패키지 경계로 의존성 방향을 강제하면서 Go 관용적 인터페이스 패턴을 최대한 활용하고,
멀티모듈의 복잡도 없이 `internal/` 접근 제어로 레이어 보호가 가능하다.

**의존성 방향:**

```
adapter/httphandler → app → domain
adapter/postgresrepo → domain (Repository interface 구현)
cmd/server → 수동 DI (모든 레이어 와이어링)
```

**UpdateFn 트랜잭션 패턴:**
Repository `Update` 메서드는 클로저를 받아 트랜잭션으로 감쌈:
```go
Update(ctx context.Context, id string, fn func(*Todo) (*Todo, error)) error
```

### Consequences

* Good, because domain 단위 테스트에 DB/HTTP 컨텍스트가 불필요하여 테스트 속도와 격리성이 높다
* Good, because adapter를 교체/추가하는 것만으로 인프라 기술 스택을 변경할 수 있다
* Good, because Go workspace 멀티모듈 대비 빌드/도구 복잡도가 낮다
* Bad, because 아키텍처 위반을 자동 검증하는 ArchUnit 상당 도구가 Go에 없어 코드 리뷰 의존
* Bad, because Port/Adapter 보일러플레이트 코드가 플랫 구조 대비 증가한다

### Confirmation

- `.claude/rules/architecture.md` 금지 패턴 목록으로 코드 리뷰 시 수동 검증
- golangci-lint `depguard` 규칙으로 금지 import 일부 자동화 가능

## Pros and Cons of the Options

### B. Go workspace 멀티모듈

* Good, because 모듈 경계가 컴파일러 수준에서 강제된다
* Bad, because Go workspace 설정, replace directive 관리 복잡도가 크다
* Bad, because 보일러플레이트로서 fork 사용자의 진입 장벽이 높아진다

### C. 플랫 구조 (Layered)

* Good, because 단순하고 학습 곡선이 낮다
* Bad, because 비즈니스 로직이 HTTP/DB 프레임워크에 결합된다
* Bad, because 인프라 교체 시 비즈니스 로직 수정 범위가 넓어진다

## More Information

* 관련 규칙: `.claude/rules/architecture.md`
* 참고: [Three Dots Labs — Go Hexagonal Architecture](https://threedots.tech/post/introducing-clean-architecture/)

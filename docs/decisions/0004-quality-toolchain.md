---
status: accepted
date: 2026-04-04
decision-makers: ppzxc
---

# 코드 품질 도구 체인: gofumpt + golangci-lint + govulncheck

## Context and Problem Statement

Go 보일러플레이트에서 일관된 코드 스타일과 정적 분석을 자동으로 강제할 도구를 선택해야 한다.
Java 템플릿의 Spotless + Checkstyle + ErrorProne + NullAway 조합에 대응하는 Go 생태계 도구가 필요하다.
Go에는 Java ArchUnit에 해당하는 아키텍처 자동 검증 도구가 없다는 점을 인정하고 대안을 수립해야 한다.

## Decision Drivers

* 자동 수정 가능한 포맷팅 도구 (Java Spotless 대응)
* 네이밍 규칙 및 금지 패턴 강제 (Java Checkstyle 대응)
* 정적 분석 및 버그 패턴 탐지 (Java ErrorProne 대응)
* 보안 취약점 스캔
* 동시성 버그 검출 (Go 특화)
* 아키텍처 위반 검출 전략 (ArchUnit 부재 대안)

## Considered Options

* A. gofumpt + golangci-lint + govulncheck + go test -race
* B. gofmt + staticcheck (golangci-lint 없이)
* C. gofumpt + golangci-lint + revive (golangci-lint 내장)

## Decision Outcome

Chosen option: "A. gofumpt + golangci-lint + govulncheck + go test -race", because
각 도구가 포맷팅, 정적 분석, 보안, 동시성 역할을 분담하여 Java 도구 체인과 동등한 품질을 보장한다.

**도구 역할 분담:**

| Go 도구 | Java 대응 | 역할 | 자동수정 |
|--------|----------|------|---------|
| gofumpt | Spotless | 코드 포맷팅 (gofmt 상위 호환) | ✓ |
| golangci-lint | Checkstyle + ErrorProne | 스타일, 버그 패턴, import 순서 | ✗ |
| govulncheck | — | 알려진 CVE 스캔 | ✗ |
| go vet | 정적 분석 기본 | 의심스러운 코드 패턴 | ✗ |
| go test -race | — | 동시성 데이터 레이스 검출 | ✗ |

**ArchUnit 부재 대안:**
Go에는 아키텍처 의존성을 자동 검증하는 ArchUnit 상당 도구가 없다.
이를 보완하기 위해 `.claude/rules/architecture.md`에 금지 패턴을 명시하고 코드 리뷰로 강제한다.
향후 `golangci-lint depguard` 규칙으로 일부 import 제한을 자동화할 수 있다.

### Consequences

* Good, because gofumpt는 gofmt보다 엄격하여 팀 내 스타일 논쟁을 제거한다
* Good, because golangci-lint가 다수의 linter를 통합 실행하여 설정 단순화
* Good, because govulncheck가 Go 모듈 의존성의 알려진 취약점을 탐지한다
* Bad, because 아키텍처 위반 자동 검증 불가 — 코드 리뷰 의존도가 높다
* Bad, because golangci-lint 설정이 복잡하여 초기 세팅에 시간이 소요된다

### Confirmation

- `make fmt` 실행 후 diff가 없는지 확인
- `make lint` CI 통과 여부
- `make vuln` 취약점 0건 확인

## Pros and Cons of the Options

### B. gofmt + staticcheck

* Good, because 외부 의존성이 최소화된다
* Bad, because golangci-lint 대비 검사 범위가 좁다
* Bad, because import 순서, 네이밍 등 추가 규칙 적용이 어렵다

### C. golangci-lint revive 단독

* Good, because golangci-lint 내장으로 별도 설치 불필요
* Bad, because gofumpt의 엄격한 포맷팅 규칙을 대체하지 못함

## More Information

* 관련 규칙: `.claude/rules/ci-tools.md`
* 관련 ADR: [ADR-0005](0005-ci-pipeline.md) — CI 파이프라인 통합

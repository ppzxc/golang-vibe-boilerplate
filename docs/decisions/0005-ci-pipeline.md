---
status: accepted
date: 2026-04-04
decision-makers: ppzxc
---

# CI 파이프라인 전략: Lefthook + GitHub Actions 4-job 구성

## Context and Problem Statement

Go 보일러플레이트의 CI/CD 파이프라인 전략을 결정해야 한다.
Java 템플릿의 Lefthook + GitHub Actions 조합이 Go 프로젝트에서도 동일하게 적용 가능한지 검토한다.
단위 테스트와 통합 테스트(testcontainers)를 분리 실행하고 커버리지를 수집해야 한다.

## Decision Drivers

* 빠른 피드백: pre-commit 단계에서 포맷팅/린팅 실패 조기 감지
* 단위 테스트와 통합 테스트 분리 (testcontainers는 Docker 필요)
* 커버리지 리포트 수집
* Java 템플릿과 동일한 도구(Lefthook, GitHub Actions) 유지
* 동시성 버그 검출 (`go test -race`)

## Considered Options

* A. Lefthook (pre-commit/pre-push) + GitHub Actions 4-job
* B. Husky + lint-staged + GitHub Actions
* C. GitHub Actions 단독 (로컬 훅 없음)

## Decision Outcome

Chosen option: "A. Lefthook (pre-commit/pre-push) + GitHub Actions 4-job", because
Java 템플릿과 동일한 Lefthook을 사용하여 일관성을 유지하고,
4-job 병렬 구성으로 빠른 CI 피드백과 단계별 실패 격리가 가능하다.

**Lefthook 구성:**

```yaml
pre-commit:
  parallel: true
  commands:
    fmt-check:
      run: gofumpt -l . | grep . && exit 1 || exit 0
    lint:
      run: golangci-lint run ./...
    vet:
      run: go vet ./...

pre-push:
  commands:
    test-race:
      run: go test -race ./...
    vuln:
      run: govulncheck ./...
```

**GitHub Actions 4-job:**

1. `lint` — gofumpt check, golangci-lint, go vet
2. `unit-test` — go test -race (integration 빌드 태그 제외)
3. `integration-test` — go test -tags=integration (testcontainers, Docker 필요)
4. `coverage` — unit + integration 커버리지 병합 리포트 (3번 완료 후 실행)

### Consequences

* Good, because pre-commit에서 포맷팅/린팅을 잡아 CI 불필요한 실패를 줄인다
* Good, because integration-test job이 분리되어 Docker 환경 의존성이 격리된다
* Good, because coverage job이 전체 테스트 완료 후 실행되어 정확한 커버리지 수집
* Bad, because Lefthook 설치가 필요하여 신규 기여자의 초기 설정 단계 증가

### Confirmation

- `lefthook install` 실행 후 `.git/hooks/pre-commit` 생성 확인
- GitHub Actions 워크플로우 4개 job이 순서에 맞게 의존성 설정되었는지 확인

## Pros and Cons of the Options

### B. Husky + lint-staged

* Good, because Node.js 생태계에서 널리 사용되어 기여자 친숙도가 높다
* Bad, because Node.js 의존성이 Go 프로젝트에 추가됨 — 불필요한 이종 의존성
* Bad, because Java 템플릿과 도구 불일치

### C. GitHub Actions 단독

* Good, because 로컬 설정 불필요
* Bad, because CI 실패까지 푸시 후 대기 — 느린 피드백 루프

## More Information

* 관련 규칙: `.claude/rules/ci-tools.md`
* 관련 ADR: [ADR-0004](0004-quality-toolchain.md) — 품질 도구 체인

# Rules Maintenance

## 문서 계층 구조

```
CLAUDE.md              — 프로젝트 진입점 (≤50줄)
.claude/rules/*.md     — 제약 규칙 (각 ≤100줄)
docs/decisions/*.md    — 결정 근거 (MADR 4.0)
```

## ADR ↔ Rules 동기화 규칙

- 새 ADR 작성 시: 관련 규칙 파일 업데이트
- 규칙 변경 시: 기존 ADR과 충돌 여부 확인
- 규칙 파일이 100줄 초과 시: 분리 또는 ADR로 이동

## ADR 형식 (MADR 4.0)

```yaml
---
status: accepted | proposed | deprecated | superseded
date: YYYY-MM-DD
decision-makers: [ppzxc]
---

# ADR 제목

## Context and Problem Statement

## Decision Drivers

## Considered Options

## Decision Outcome

### Consequences

### Confirmation

## Pros and Cons of the Options

## More Information
```

## ADR 네이밍

`NNNN-kebab-case-title.md` (0으로 채운 4자리 숫자)

## 규칙 파일 작성 원칙

- 제약 사항만 (구현 세부사항 없음)
- 예시 코드 포함 시 최소화
- "왜"보다 "무엇을 하면/하지 말면 안 되는가"에 집중
- 100줄 초과 = 분리 신호

## CLAUDE.md 유지

- 50줄 초과 금지
- 규칙 라우팅 테이블 최신 상태 유지
- 새 규칙 파일 추가 시 테이블에 행 추가

# CI Tools Rules

## 도구 목록 (Java 템플릿 대응)

| Go 도구 | Java 대응 | 역할 |
|--------|----------|------|
| gofumpt | Spotless | 코드 포맷팅 |
| golangci-lint | Checkstyle + ErrorProne | 스타일 + 정적 분석 |
| govulncheck | — | 보안 취약점 스캔 |
| go vet | 정적 분석 기본 | 정적 분석 |
| go test -race | — | 동시성 버그 검출 |
| lefthook | lefthook | Git hooks |
| GitHub Actions | GitHub Actions | CI |

## Makefile 명령어

```bash
make run            # 서버 실행
make build          # 바이너리 빌드
make test           # 단위 테스트 (go test -race ./...)
make test-integration # 통합 테스트 (go test -tags=integration ./...)
make test-all       # 전체 테스트
make lint           # golangci-lint run
make fmt            # gofumpt -w .
make vet            # go vet ./...
make vuln           # govulncheck ./...
make migrate-up     # DB 마이그레이션 적용
make migrate-down   # DB 마이그레이션 롤백
make migrate-create # 새 마이그레이션 파일 생성
make sqlc           # sqlc generate
make docker-build   # Docker 이미지 빌드
make docker-up      # docker-compose up
make docker-down    # docker-compose down
```

## Lefthook 훅

- **pre-commit** (병렬): gofumpt 체크, golangci-lint, go vet
- **pre-push**: go test -race, govulncheck

## GitHub Actions Jobs

1. `lint` — gofumpt check, golangci-lint, go vet
2. `unit-test` — go test -race (integration 제외)
3. `integration-test` — go test -tags=integration (testcontainers)
4. `coverage` — 커버리지 리포트 (unit + integration 후)

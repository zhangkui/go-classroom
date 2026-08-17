# go-classroom candidate 4

## 项目说明
`zhangkui/go-classroom` 是使用内存存储的课程、成员、作业、提交、评分与统计 Web API。Go 工具链为 `golang:1.22`，无前端工具链。

## 标准构建、运行和测试命令
```bash
cd '/app' && GOTOOLCHAIN=local go build ./...
cd '/app' && GOTOOLCHAIN=local go run ./cmd
cd '/app' && GOTOOLCHAIN=local go test ./...
```

## Docker 构建和进入容器
```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh go-classroom-bug4-candidate linux/amd64
./build_benzhi_docker.sh go-classroom-bug4-candidate-arm64 linux/arm64
docker run --rm -it go-classroom-bug4-candidate bash
```

## 题目验证命令
```bash
go test -buildvcs=false -count=1 -run "TestSubmissionAtDeadlineIsRejected" ./internal/classroom/
go test -buildvcs=false -count=1 ./...
go build -buildvcs=false ./...
```

## Bug 复现
见 `BUG_REPRO.md`。

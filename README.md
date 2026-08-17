# go-classroom

## 项目说明
一个使用 Go 标准库实现的小型课程和作业管理 Web API。教师可以管理课程、作业与评分，学生可以加入课程并提交答案。数据仅保存在进程内存中。

## 标准命令
```bash
go build ./...
go test ./...
go vet ./...
go run ./cmd
```

## 使用方式
启动后服务监听 `:8080`。通过 `POST /courses` 创建课程，通过 `GET /courses?page=1&page_size=20` 查询课程，并使用 `/courses/{id}/members` 管理成员。

# ============ 第一阶段：构建阶段 ============
# 关键修改1：将版本改为更通用的标签，或匹配你的Go版本
FROM golang:1.25.2-alpine3.22 AS build

# 设置工作目录
WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制所有源代码
COPY . .

# 关键修改2：使用更通用的构建命令
# 假设你的主包在项目根目录，编译当前目录的所有go文件
RUN go build -o serve

# ============ 第二阶段：运行阶段 ============
# 关键修改3：提供两个选项，注释掉其中一个

# 选项A：使用scratch（最小，但无法调试）
#FROM scratch

# 选项B：使用alpine（稍大，但可调试，推荐初次使用）
FROM alpine:3.22
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
# 从构建阶段复制必要的文件
COPY --from=build /app/serve .

# 设置时区（可选）
ENV TZ=Asia/Shanghai

# 暴露端口
EXPOSE 8080

# 运行程序
ENTRYPOINT ["./serve"]
CMD ["--account", "root", "--password", "123456"]
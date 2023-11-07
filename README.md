# QiYin Backend

柒音后端服务仓库

### 设计文档和DEMO

文档地址和演示地址：https://eq2pyit41ih.feishu.cn/docx/M6L8dYYg6oq3cvxsBSpcoteLnuc

### 运行服务

1. 安装docker和docker-compose

2. 准备mysql/redis ... etc
```bash
docker-compose up
```

3. 迁移数据库
```bash
cd cmd/dbmigration
go run main.go -c xxx.yaml # 后面为你的配置文件
```

3. 运行后端程序
```bash
# 根目录下
go run main.go -c xxx.yaml # 后面为你的配置文件
```
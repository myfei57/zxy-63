基于 Go 实现的自来水厂加药与滤池控制平台项目，一款后端服务，完成进水流量加矾加氯、滤池轮换与反冲洗、清水池液位调节与加药审计。

## 项目简介

WaterPlant 是自来水厂制水过程控制平台，按「取水 → 加矾/加氯 → 滤池 → 清水池 → 计量 → 审计」的工艺线管理制水过程。服务使用文件型持久化保存流量、配比、余氯目标、清水池液位、加氯累计量与审计记录，并通过 HTTP JSON API 提供控制与查询能力。

## 运行环境

- Go 1.23
- 依赖已 vendor：github.com/go-chi/chi/v5、github.com/google/uuid

## 本地运行

```bash
go build -mod=vendor -o waterplant ./cmd/waterplant
./waterplant -store waterplant.json -addr :8080
```

启动后可访问以下端点：

- GET /health
- GET /snapshot
- GET /report
- GET /catalog
- POST /cycle
- POST /simulate

## Docker 构建

```bash
sh build_benzhi_docker.sh
docker run -p 8080:8080 waterplant-control
```

## 持久化

控制台使用单个 JSON 文件保存运行时状态，默认路径为 waterplant.json，可通过 -store 参数指定。

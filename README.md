# ✨ 虚拟心愿墙社交平台（wishwall）

> 一个用善意连接彼此的心愿墙：发布心愿、认领心愿成为圆梦人、送祝福与虚拟礼物、封存时光胶囊、收集成就徽章。

## 快速启动（Docker Compose，推荐）

```bash
cd /Users/gaobo/repositories/gitlab/评审项目/0-1代码生成提示词/golang-改编提示词/趣味社交主题项目提示词/lp-153
docker compose up -d --build
```

启动后访问：

| 入口 | 地址 |
| --- | --- |
| 前端页面 | http://localhost:18403 |
| 后端 API | http://localhost:19403/api/v1 |
| 健康检查 | http://localhost:19403/healthz |
| MinIO 控制台 | http://localhost:47032 （账号 minioadmin / minioadmin） |

默认管理员账号：`admin / Admin@123456`（首次启动自动播种）。

关闭并清理（含数据卷）：

```bash
docker compose down -v --remove-orphans
```

## 项目主要功能

1. **心愿发布**：文字 + 图片，分类（学习成长/旅行探险/情感陪伴/职业发展/生活小确幸/其他），可见范围（公开/好友可见/匿名），期望完成时间 + 难度标签。
2. **心愿认领与进度追踪**：心愿广场浏览并认领心愿成为「圆梦人」，更新进度（百分比 + 文字），支持里程碑打卡。
3. **发布者验收**：圆梦人送交完成说明后心愿进入「待验收」，仅发布者可确认或驳回（驳回必填原因、保留说明可重新送交）；确认后心愿才完成并发放徽章，同一心愿仅一条验收记录，并发验收只成功一次。
4. **祝福留言板**：每个心愿专属留言板，送祝福与虚拟礼物（🎁 表情包）；心愿完成自动转为庆祝页。
5. **时光胶囊**：定时解锁的文字 + 图片 + 音频胶囊；解锁前内容打码，到期自动解锁并播放解锁动画。
6. **心愿成就徽章**：首次许愿、首次认领、首次祝福、十次圆梦、圆梦大师；展示在个人主页。
7. **搜索与发现广场**：按标签/关键词搜索，热门圆梦人排行榜 + 最新完成的心愿故事。

## 技术栈

| 层 | 技术 |
| --- | --- |
| 前端 | Next.js 14 + TypeScript + Tailwind CSS |
| 后端 | Go 1.22 + Gin + GORM |
| 数据库 | PostgreSQL 16 |
| 缓存 | Redis 7 |
| 对象存储 | MinIO（心愿图片 / 胶囊音频） |
| 认证 | JWT + RBAC |
| 日志 | `log/slog` 结构化日志 |
| 参数校验 | `github.com/go-playground/validator/v10` |
| 接口文档 | 本 README 完整 API 清单 |

## 目录结构

```
lp-153/
├── backend/
│   ├── cmd/server/main.go          # 装配入口：配置/依赖/启动/优雅退出
│   ├── internal/
│   │   ├── config/                 # 环境变量配置
│   │   ├── database/               # PostgreSQL / Redis / MinIO 连接
│   │   ├── model/                  # 8 个实体（user/wish/claim/acceptance/blessing/capsule/badge/audit）
│   │   ├── dto/                    # 每个实体一个 DTO 文件（含 validator 校验）
│   │   ├── repository/             # 每个实体一个仓储文件（哨兵错误）
│   │   ├── service/                # 每个实体一个服务文件（事务/状态机）
│   │   ├── handler/                # 每个实体一个处理器
│   │   ├── router/                 # 每个实体一个路由注册文件
│   │   ├── middleware/             # auth/rbac/request_id/error_handler/audit/ratelimit
│   │   ├── constants/              # 枚举/错误码/日志模板/文案
│   │   └── util/                   # jwt/password/logger/formatters/app_error
│   ├── migrations/                 # 数据库初始化 SQL（参考）
│   ├── pkg/response/               # 统一响应
│   ├── Dockerfile
│   ├── go.mod / go.sum
│   └── api/                        # OpenAPI 预留目录
├── frontend/
│   ├── src/
│   │   ├── api/                    # 每个实体一个 API 文件
│   │   ├── components/             # 共享组件（StatusBadge/WishCard/ProgressBar/...）
│   │   ├── pages/                  # 按模块拆分的页面
│   │   ├── stores/                 # zustand 按实体拆分（auth/wish/capsule/badge）
│   │   ├── hooks/                  # useAuth / usePagination
│   │   ├── utils/                  # request.ts / format.ts
│   │   └── constants/              # 与后端对应的枚举
│   ├── Dockerfile                  # 多阶段构建 → Nginx 托管静态导出
│   └── nginx.conf                  # 前端路由 + /api 反向代理
├── database/init.sql               # 数据库脚本
├── docker-compose.yml
├── .env / .env.example
└── README.md
```

**严禁合并职责到单一文件**：每个实体均按 model / dto / repository / service / handler / router / constants 独立文件拆分，前端按 api / stores / components / pages / hooks / utils / constants 拆分，禁止把多个实体或全部页面写进同一文件。

## 环境变量说明

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| `COMPOSE_PROJECT_NAME` | `wishwall` | Compose 项目名与容器名前缀 |
| `FRONTEND_PORT` | `18403` | 前端宿主端口（容器 80） |
| `BACKEND_PORT` | `19403` | 后端宿主端口（容器 8080） |
| `DB_PORT` | `57607` | PostgreSQL 宿主端口（容器 5432） |
| `REDIS_PORT` | `46325` | Redis 宿主端口（容器 6379） |
| `MINIO_PORT` | `47031` | MinIO API 宿主端口（容器 9000） |
| `MINIO_CONSOLE_PORT` | `47032` | MinIO 控制台宿主端口（容器 9001） |
| `DB_NAME` / `DB_USER` / `DB_PASSWORD` | `wishwall_db` / `wishwall_user` / `wishwall_pwd` | PostgreSQL 连接 |
| `JWT_SECRET` | 随机长字符串 | JWT 签名密钥（生产务必修改） |
| `ADMIN_PASSWORD` | `Admin@123456` | 默认管理员密码 |
| `REDIS_ADDR` / `REDIS_PASSWORD` | `redis:6379` / 空 | 缓存连接 |
| `MINIO_ACCESS_KEY` / `MINIO_SECRET_KEY` | `minioadmin` / `minioadmin` | MinIO 凭据 |
| `UPLOAD_BASE_URL` | `http://localhost:47028` | 上传文件对外访问基地址 |

## API 调用示例（curl）

### 1. 注册

```bash
curl -sS -X POST http://localhost:19403/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","email":"alice@example.com","password":"secret123","nickname":"爱丽丝"}'
```

### 2. 登录（获取 JWT）

```bash
TOKEN=$(curl -sS -X POST http://localhost:19403/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"account":"alice","password":"secret123"}' | python3 -c 'import sys,json;print(json.load(sys.stdin)["data"]["token"])')
```

### 3. 发布心愿（携带 JWT）

```bash
curl -sS -X POST http://localhost:19403/api/v1/wishes \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title":"去冰岛看极光","content":"想在极夜里看到绿色极光","category":"travel","visibility":"public","difficulty":"hard","expected_deadline":"2027-12-31T23:59:59+08:00"}'
```

### 4. 心愿广场列表

```bash
curl -sS "http://localhost:19403/api/v1/wishes?page=1&page_size=10&sort=hot"
```

### 5. 认领心愿（成为圆梦人）

```bash
curl -sS -X POST http://localhost:19403/api/v1/wishes/1/claim \
  -H "Authorization: Bearer $TOKEN"
```

### 6. 更新圆梦进度

```bash
curl -sS -X PUT http://localhost:19403/api/v1/claims/1/progress \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"progress":60,"note":"已经订好机票","is_milestone":true}'
```

### 7. 送交验收（圆梦人，说明必填）

```bash
curl -sS -X POST http://localhost:19403/api/v1/claims/1/acceptance \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"note":"极光真的出现了！照片已送达"}'
```

### 8. 发布者验收（确认 / 驳回，驳回必须写原因）

```bash
# 确认验收：心愿完成并发放徽章（仅心愿发布者）
curl -sS -X POST http://localhost:19403/api/v1/acceptances/1/confirm \
  -H "Authorization: Bearer $TOKEN"

# 驳回验收：记录回到圆梦中并保留说明，可重新送交（仅心愿发布者）
curl -sS -X POST http://localhost:19403/api/v1/acceptances/1/reject \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"reason":"照片不清晰，请补充后重新送交"}'
```

### 9. 送祝福

```bash
curl -sS -X POST http://localhost:19403/api/v1/wishes/1/blessings \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"content":"祝你梦想成真！","gift_emoji":"🎁"}'
```

### 10. 封存时光胶囊

```bash
curl -sS -X POST http://localhost:19403/api/v1/capsules \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title":"给一年后的自己","content":"要更勇敢","unlock_at":"2027-08-17T00:00:00+08:00"}'
```

## API 清单（统一前缀 `/api/v1`，响应统一 `{"code":0,"message":"ok","data":...}`）

### 认证与用户

| 方法 | 路径 | 说明 | 鉴权 |
| --- | --- | --- | --- |
| POST | `/auth/register` | 注册 | 公开 |
| POST | `/auth/login` | 登录（返回 JWT） | 公开 |
| GET | `/users/me` | 当前用户信息 | JWT |
| PUT | `/users/me` | 更新个人资料 | JWT |
| GET | `/users/:id` | 用户公开信息 | 公开 |
| GET | `/healthz` | 健康检查（根路径） | 公开 |

### 心愿

| 方法 | 路径 | 说明 | 鉴权 |
| --- | --- | --- | --- |
| POST | `/wishes` | 发布心愿 | JWT |
| GET | `/wishes` | 心愿广场列表（keyword/category/status/visibility/sort/page/page_size） | 公开 |
| GET | `/wishes/mine` | 我发布的心愿 | JWT |
| GET | `/wishes/:id` | 心愿详情（含认领摘要、祝福数） | 公开 |
| PUT | `/wishes/:id` | 更新心愿（仅作者） | JWT |
| DELETE | `/wishes/:id` | 删除心愿（仅作者） | JWT |
| POST | `/wishes/:id/like` | 点赞 | JWT |

### 认领 / 圆梦

| 方法 | 路径 | 说明 | 鉴权 |
| --- | --- | --- | --- |
| POST | `/wishes/:id/claim` | 认领心愿（事务 + 行锁防并发重复认领） | JWT |
| GET | `/wishes/:id/claim` | 心愿的认领记录 | JWT |
| GET | `/claims/mine` | 我认领的心愿 | JWT |
| PUT | `/claims/:id/progress` | 更新进度（里程碑打卡；待验收期间禁止更新，进度 100% 需送交验收） | JWT |

### 验收（发布者确认/驳回）

| 方法 | 路径 | 说明 | 鉴权 |
| --- | --- | --- | --- |
| POST | `/claims/:id/acceptance` | 圆梦人送交验收（说明必填；同一心愿仅一条记录，唯一索引兜底） | JWT |
| GET | `/wishes/:id/acceptance` | 心愿的验收记录（仅发布者/圆梦人；心愿详情亦内嵌验收摘要） | JWT |
| POST | `/acceptances/:id/confirm` | 发布者确认验收（行锁保证并发只成功一次；确认后心愿完成并发放徽章） | JWT |
| POST | `/acceptances/:id/reject` | 发布者驳回验收（原因必填；记录回到圆梦中并保留说明，可重新送交） | JWT |

### 祝福留言板

| 方法 | 路径 | 说明 | 鉴权 |
| --- | --- | --- | --- |
| POST | `/wishes/:id/blessings` | 发送祝福（可带虚拟礼物） | JWT |
| GET | `/wishes/:id/blessings` | 祝福列表 | 公开 |

### 时光胶囊

| 方法 | 路径 | 说明 | 鉴权 |
| --- | --- | --- | --- |
| POST | `/capsules` | 封存胶囊 | JWT |
| GET | `/capsules/mine` | 我的胶囊（未解锁内容打码） | JWT |
| GET | `/capsules/:id` | 胶囊详情（本人） | JWT |
| DELETE | `/capsules/:id` | 删除胶囊（本人） | JWT |

### 成就徽章 / 发现

| 方法 | 路径 | 说明 | 鉴权 |
| --- | --- | --- | --- |
| GET | `/badges/mine` | 我的徽章 | JWT |
| GET | `/badges/leaderboard` | 圆梦人排行榜（复用 `BadgeService.Leaderboard`） | 公开 |
| GET | `/discover` | 最新完成的心愿故事（复用 `WishService.List`） | 公开 |
| GET | `/discover/leaderboard` | 发现广场排行榜（复用 `BadgeService.Leaderboard`） | 公开 |

### 审计 / 上传

| 方法 | 路径 | 说明 | 鉴权 |
| --- | --- | --- | --- |
| GET | `/audit-logs` | 操作审计日志 | 管理员 |
| POST | `/uploads?kind=image\|audio` | 文件上传（multipart 字段 `file`） | JWT |

**服务/仓储复用标注**：
- `GET /wishes` 与 `GET /discover` 复用 `WishService.List` / `WishRepository.Count`。
- `GET /discover/leaderboard` 与 `GET /badges/leaderboard` 复用 `BadgeService.Leaderboard`。
- `GET /claims/mine` 与个人主页的认领展示复用 `WishClaimService.ListMine`。
- `GET /wishes/:id`（详情内嵌验收摘要）与 `GET /wishes/:id/acceptance` 复用 `WishAcceptanceRepository.FindByWishID`。

## 横切关注点

1. **JWT 认证 + RBAC 权限**：数据库 `users.role` 字段 → `backend/internal/middleware/auth.go`、`middleware/rbac.go`、`util/jwt.go` → 前端 `RequireAuth` 路由守卫、`useAuth` 按钮显隐、`/audit` 管理员页。
2. **操作审计日志**：`audit_logs` 表 → `middleware/audit.go`（写操作自动埋点）+ 各 service 业务埋点（`audit.Record`）→ 前端 `/audit` 审计页。
3. **全局错误处理与请求追踪**：`middleware/request_id.go`、`middleware/error_handler.go`（Recovery）、`util/app_error.go`、`constants/error_codes.go` → 前端 `utils/request.ts` 拦截器（401 跳转、统一解包）。

## 共享枚举出现位置清单

### 枚举 1：心愿状态（pending / claimed / in_progress / completed）

| 层 | 位置 |
| --- | --- |
| 后端 constants | `backend/internal/constants/wish_status.go` |
| 后端模型 | `backend/internal/model/wish.go`（Status 字段）、`model/wish_claim.go`（Status 字段） |
| 后端 DTO | `backend/internal/dto/wish_dto.go`（查询参数） |
| 后端状态机 | `backend/internal/service/wish_claim_service.go`（认领/进度/完成流转） |
| 后端 handler 校验 | `backend/internal/handler/wish_handler.go`、`handler/wish_claim_handler.go` |
| 后端日志模板 | `backend/internal/constants/log_templates.go`（`LogWishClaimed`/`LogClaimCompleted` 等） |
| 后端错误码 | `backend/internal/constants/error_codes.go`（`CodeWishAlreadyClaimed`/`CodeWishStatusInvalid` 等） |
| 后端 formatters | `backend/internal/util/formatters.go`（`FormatWishStatus`） |
| 前端 constants | `frontend/src/constants/index.ts`（`WISH_STATUS`/`WISH_STATUS_TEXT`/`WISH_STATUS_STYLE`） |
| 前端筛选/徽标 | `frontend/src/pages/index.tsx`（状态筛选）、`src/components/StatusBadge.tsx` |

### 枚举 2：可见范围（public / friend / anonymous）

| 层 | 位置 |
| --- | --- |
| 后端 constants | `backend/internal/constants/visibility.go` |
| 后端模型 | `backend/internal/model/wish.go`（Visibility 字段） |
| 后端 DTO | `backend/internal/dto/wish_dto.go`（Create/Update/Query） |
| 后端 handler 校验 | `backend/internal/handler/wish_handler.go`（validator oneof） |
| 后端 formatters | `backend/internal/util/formatters.go`（`FormatVisibility`） |
| 前端 constants | `frontend/src/constants/index.ts`（`VISIBILITY`/`VISIBILITY_TEXT`） |
| 前端页面 | `frontend/src/pages/wishes/create.tsx`（可见范围选择） |

### 枚举 3：难度标签（easy / medium / hard）

| 层 | 位置 |
| --- | --- |
| 后端 constants | `backend/internal/constants/difficulty.go` |
| 后端模型 | `backend/internal/model/wish.go`（Difficulty 字段） |
| 后端 DTO | `backend/internal/dto/wish_dto.go` |
| 后端 formatters | `backend/internal/util/formatters.go`（`FormatDifficulty`） |
| 前端 constants | `frontend/src/constants/index.ts`（`DIFFICULTY`/`DIFFICULTY_TEXT`） |
| 前端页面 | `frontend/src/pages/wishes/create.tsx`、`src/components/WishCard.tsx` |

### 枚举 4：角色（user / admin）

| 层 | 位置 |
| --- | --- |
| 后端 constants | `backend/internal/constants/role.go` |
| 后端模型 | `backend/internal/model/user.go`（Role 字段） |
| 后端 RBAC | `backend/internal/middleware/rbac.go`、`router/router.go`（`RequireRole(RoleAdmin)`） |
| 后端 JWT | `backend/internal/util/jwt.go`（Claims.Role） |
| 后端日志模板 | `backend/internal/constants/log_templates.go` |
| 前端 constants | `frontend/src/constants/index.ts`（`ROLE`/`ROLE_TEXT`） |
| 前端守卫 | `frontend/src/components/RequireAuth.tsx`（adminOnly）、`src/pages/audit.tsx` |

### 枚举 5：时光胶囊状态（locked / unlocked）

| 层 | 位置 |
| --- | --- |
| 后端 constants | `backend/internal/constants/capsule_status.go` |
| 后端模型 | `backend/internal/model/time_capsule.go`（Status 字段） |
| 后端状态机 | `backend/internal/service/time_capsule_service.go`（GetByID/UnlockDue 自动解锁） |
| 后端 formatters | `backend/internal/util/formatters.go`（`FormatCapsuleStatus`） |
| 前端 constants | `frontend/src/constants/index.ts`（`CAPSULE_STATUS`/`CAPSULE_STATUS_TEXT`） |
| 前端页面 | `frontend/src/pages/capsules.tsx`、`src/components/StatusBadge.tsx` |

### 枚举 6：成就徽章类型（first_wish / first_claim / first_blessing / ten_completions / wish_master）

| 层 | 位置 |
| --- | --- |
| 后端 constants | `backend/internal/constants/badge_type.go` |
| 后端模型 | `backend/internal/model/badge.go`（Type 字段） |
| 后端发放逻辑 | `backend/internal/service/badge_service.go` |
| 后端 formatters | `backend/internal/util/formatters.go`（`FormatBadgeType`） |
| 前端 constants | `frontend/src/constants/index.ts`（`BADGE_TYPE`/`BADGE_TYPE_TEXT`） |
| 前端页面 | `frontend/src/pages/profile.tsx`（徽章展示） |

### 枚举 7：验收状态（pending / confirmed / rejected）

| 层 | 位置 |
| --- | --- |
| 后端 constants | `backend/internal/constants/acceptance_status.go` |
| 后端模型 | `backend/internal/model/wish_acceptance.go`（Status 字段） |
| 后端 DTO | `backend/internal/dto/wish_acceptance_dto.go`（送交/驳回入参与验收返回结构） |
| 后端状态机 | `backend/internal/service/wish_acceptance_service.go`（送交/确认/驳回流转）、`service/wish_claim_service.go`（待验收禁止更新进度） |
| 后端仓储 | `backend/internal/repository/wish_acceptance_repository.go`（`wish_id` 唯一索引 + `FOR UPDATE` 行锁） |
| 后端 handler 校验 | `backend/internal/handler/wish_acceptance_handler.go`（说明/原因必填校验） |
| 后端日志模板 | `backend/internal/constants/log_templates.go`（`LogAcceptanceSubmitted`/`LogAcceptanceConfirmed`/`LogAcceptanceRejected`） |
| 后端错误码 | `backend/internal/constants/error_codes.go`（`CodeAcceptanceNotFound`/`CodeAcceptancePending`/`CodeAcceptanceAlreadySubmitted`/`CodeAcceptanceStatusInvalid`/`CodeAcceptanceNotPublisher`/`CodeAcceptanceRequired`） |
| 后端 formatters | `backend/internal/util/formatters.go`（`FormatAcceptanceStatus`） |
| 前端 constants | `frontend/src/constants/index.ts`（`ACCEPTANCE_STATUS`/`ACCEPTANCE_STATUS_TEXT`/`ACCEPTANCE_STATUS_STYLE`） |
| 前端 API/类型 | `frontend/src/api/acceptance.ts`、`frontend/src/api/wish.ts`（`WishDetail.acceptance`） |
| 前端徽标/页面 | `frontend/src/components/StatusBadge.tsx`（kind=acceptance）、`frontend/src/pages/wishes/detail.tsx`（验收状态、驳回原因、确认/驳回/重新送交入口） |

## 屎山代码设计说明（跨文件协同约束）

为验证跨文件协同改动能力，本实现刻意保留了以下「合理耦合」：

1. **日志模块全栈引用**：`internal/util/logger.go` 封装 `log/slog`，所有 handler/service/middleware 均引用；25+ 条日志模板集中在 `internal/constants/log_templates.go`。
2. **异常信息分散透传**：错误码集中在 `internal/constants/error_codes.go`，各 service/handler 手动拼接 message（包含实体名/字段名/角色名），handler 再次包装 service 错误。
3. **常量/工具多处耦合**：`internal/util/formatters.go` 同时提供日期、状态文本、类型文本格式化；`internal/constants/messages.go` 同时承载接口返回文案、日志文案、错误提示。
4. **状态机跨多处定义**：心愿/胶囊状态流转规则同时存在于 service 状态机、前端按钮显隐、日志模板、错误码、formatters。
5. **枚举多处重复定义**：核心枚举在 constants、DTO、模型、日志模板、错误码、formatters、前端 constants 中同时出现（见上表）。

> 给核心实体新增字段/状态时，需同步修改 model / dto / constants / service / repository / handler / formatters / 日志模板 / 错误码 / 前端类型与页面等 ≥ 10 个文件。

## Docker 部署说明

- 端口映射见上表；容器名统一带 `${COMPOSE_PROJECT_NAME:-wishwall}` 前缀。
- 数据卷：`db_data`（PostgreSQL）、`redis_data`（Redis）、`minio_data`（MinIO），删除容器不丢数据；`docker compose down -v` 会清空。
- 所有服务配置了 `healthcheck`，后端通过 `depends_on.condition: service_healthy` 等待数据库/Redis/MinIO 就绪。
- 常见问题：
  - 端口被占用：修改 `.env` 中对应端口后 `docker compose up -d` 重建。
  - 前端页面打不开：检查 `frontend` 容器日志，确认 `backend` healthy。
  - 图片不显示：确认 MinIO 容器 healthy，`UPLOAD_BASE_URL` 指向可达地址。

## 本地开发

```bash
# 后端
cd backend
go mod tidy
go run ./cmd/server

# 前端
cd frontend
npm install
npm run dev
```

后端构建命令：`cd backend && go build ./...`；单元测试：`go test ./...`。

## License

MIT License

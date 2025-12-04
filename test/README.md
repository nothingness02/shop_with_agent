# API 测试指南

本目录包含针对 Shop 项目 API 的完整测试套件。

## 文件说明

### 1. `api_test.go`
Go 标准测试文件，使用 `testing` 包编写单个接口的测试用例。

**特点**：
- 每个接口一个测试函数
- 基础的 HTTP 请求测试
- 适合集成到 CI/CD 流程

**运行方式**：
```bash
# 在项目根目录运行
go test ./test -v

# 运行特定测试
go test ./test -v -run TestCreateShop

# 运行所有 Shop 相关测试
go test ./test -v -run TestShop*
```

### 2. `integration_test.go`
集成测试程序，包含完整的业务流程测试。

**特点**：
- 包含 `APIClient` 结构体，提供可复用的 HTTP 请求方法
- 测试完整的业务流程（创建→修改→查询→删除）
- 更详细的错误信息和成功提示
- 可作为独立程序运行

**运行方式**：
```bash
# 编译并运行集成测试
cd test
go run integration_test.go

# 或编译后运行
go build -o test.exe
./test.exe
```

**输出示例**：
```
╔══════════════════════════════════════════╗
║     API Integration Tests Started        ║
╚══════════════════════════════════════════╝

========== Testing Shop Workflow ==========

[1] Creating a new shop...
✓ Shop created
Response: {"ID":1,"CreatedAt":"...","Name":"TechMart Store",...}
Shop ID: 1

[2] Getting shop details...
✓ Shop retrieved
...
```

### 3. `postman_collection.json`
Postman API 测试集合，可直接导入到 Postman 中。

**导入步骤**：
1. 打开 Postman
2. 点击 `Import` 按钮
3. 选择 `postman_collection.json` 文件
4. 选择 `Import as a copy`

**功能**：
- 按模块组织（Shop、Product、Order）
- 自动提取和保存响应中的 ID
- 包含预定义的请求体
- 支持环境变量

**使用**：
1. 在 Postman 中创建环境变量：
   - `baseURL`: http://localhost:8080
   - `shopId`: （会自动填充）
   - `productId`: （会自动填充）
   - `orderId`: （会自动填充）

2. 按顺序执行请求：
   - 先创建资源（Create）
   - 再查询资源（Get/List）
   - 然后修改资源（Update）
   - 最后删除资源（Delete）

## API 端点快速参考

### Shop 管理（v2）
| 方法 | 路由 | 功能 |
|------|------|------|
| POST | `/api/v2/shops` | 创建商店 |
| GET | `/api/v2/shops` | 获取商店列表 |
| GET | `/api/v2/shops/:id` | 获取商店详情 |
| PATCH | `/api/v2/shops/:id` | 更新商店 |
| DELETE | `/api/v2/shops/:id` | 删除商店 |
| DELETE | `/api/v2/shops` | 批量删除商店 |

### Product 管理（v2）
| 方法 | 路由 | 功能 |
|------|------|------|
| POST | `/api/v2/shops/:shop_id/products` | 创建商品 |
| GET | `/api/v2/shops/:shop_id/products` | 获取商品列表 |
| GET | `/api/v2/products/:id` | 获取商品详情 |
| GET | `/api/v2/shops/:shop_id/products/search?name=xxx` | 按名称搜索商品 |
| PATCH | `/api/v2/products/:id` | 更新商品 |
| DELETE | `/api/v2/products/:id` | 删除商品 |
| DELETE | `/api/v2/products` | 批量删除商品 |

### Order 管理（v1）
| 方法 | 路由 | 功能 |
|------|------|------|
| POST | `/api/v1/orders` | 创建订单 |
| GET | `/api/v1/orders` | 获取订单列表 |
| GET | `/api/v1/orders/:id` | 获取订单详情 |
| PATCH | `/api/v1/orders/:id/status` | 更新订单状态 |
| DELETE | `/api/v1/orders/:id` | 删除订单 |

## 测试前的准备工作

### 1. 启动数据库
```bash
# PostgreSQL 必须运行在 localhost:5432
# 数据库配置见 cmd/api.go
```

### 2. 启动应用服务器
```bash
cd shop_project
go run ./cmd
```

服务器将在 `http://localhost:8080` 启动。

### 3. 初始化测试数据（可选）
```bash
# 导入测试数据
psql -U postgres -d app_db -f ./pkg/database/sql/01_init_schema.sql
psql -U postgres -d app_db -f ./pkg/database/sql/02_insert_test_data.sql
```

## 常见问题

### Q: 运行测试时连接被拒绝
**A**: 确保：
1. 应用服务器已启动（`go run ./cmd`）
2. 服务器运行在 `http://localhost:8080`
3. 数据库连接正常

### Q: 如何查看详细的请求/响应日志
**A**: 修改 `integration_test.go` 中的 `DoRequest` 方法，添加日志输出：
```go
fmt.Printf("Request: %s %s\n", method, url)
fmt.Printf("Response: %s (Status: %d)\n", string(data), status)
```

### Q: 如何修改测试数据
**A**: 
- 对于 `api_test.go`，直接编辑测试函数中的请求体
- 对于 `integration_test.go`，修改 `TestFullShopWorkflow()` 等函数中的请求数据
- 对于 Postman，直接编辑集合中的请求体

## 性能测试建议

使用 `bombardier` 或 `wrk` 进行压力测试：

```bash
# 安装 bombardier
go install github.com/codesenberg/bombardier@latest

# 测试 GET 请求
bombardier -c 10 -n 1000 http://localhost:8080/api/v2/shops

# 测试 POST 请求
bombardier -c 10 -n 100 -m POST -f ./test/shop_payload.json http://localhost:8080/api/v2/shops
```

## 自动化测试流程

创建 GitHub Actions 工作流或 GitLab CI 流程自动运行测试：

```yaml
# .github/workflows/test.yml
name: API Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: app_db
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: 1.24.7
      
      - name: Run tests
        run: |
          go test ./test -v
          cd test && go run integration_test.go
```

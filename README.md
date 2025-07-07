# Utils Store

一个基于Go语言的工具库，提供缓存和数据存储的通用解决方案。

## 功能特性

### 🚀 缓存模块 (Cache)
- **KV缓存接口**：支持泛型的键值对缓存操作
- **Redis支持**：完整的Redis客户端封装
- **内存缓存**：简单的内存缓存实现
- **缓存单例**：单例模式的缓存管理
- **自动编码/解码**：支持结构体的自动序列化

### 🗄️ 存储模块 (Store)
- **GORM集成**：基于GORM的数据库操作封装
- **泛型支持**：类型安全的数据库操作
- **完整CRUD**：增删改查操作的完整实现
- **事务支持**：完整的事务管理
- **分页查询**：内置分页功能
- **搜索功能**：全文搜索和条件搜索
- **软删除**：支持软删除操作
- **历史记录**：数据变更历史追踪

## 安装

```bash
go get github.com/mengri/utils-store
```

## 依赖

- Go 1.23.10+
- GORM v1.30.0+
- Redis Go客户端 v9.11.0+
- MySQL驱动 v1.6.0+

## 快速开始

### 缓存使用示例

```go
package main

import (
    "context"
    "time"
    
    "github.com/mengri/utils-store/cache"
    "github.com/mengri/utils-store/cache/cache_redis"
)

type User struct {
    ID   int64  `json:"id"`
    Name string `json:"name"`
}

func main() {
    // 创建Redis缓存客户端
    redisClient := cache_redis.NewRedisClient(&cache_redis.Config{
        Addr: "localhost:6379",
    })
    
    // 创建KV缓存
    userCache := cache.CreateKvCache[User, int64](
        redisClient,
        time.Minute*10, // 10分钟过期
    )
    
    ctx := context.Background()
    
    // 设置缓存
    user := &User{ID: 1, Name: "张三"}
    err := userCache.Set(ctx, user.ID, user)
    if err != nil {
        panic(err)
    }
    
    // 获取缓存
    cached, err := userCache.Get(ctx, 1)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("缓存的用户: %+v\n", cached)
}
```

### 存储使用示例

```go
package main

import (
    "context"
    
    "github.com/mengri/utils-store/store"
    "github.com/mengri/utils-store/store/store_mysql"
)

type User struct {
    ID   int64  `gorm:"primaryKey" json:"id"`
    UUID string `gorm:"uniqueIndex" json:"uuid"`
    Name string `json:"name"`
}

func (u User) TableName() string {
    return "users"
}

func (u User) IdValue() int64 {
    return u.ID
}

func main() {
    // 创建MySQL连接
    db := store_mysql.NewMysqlDB(&store_mysql.Config{
        Host:     "localhost",
        Port:     3306,
        User:     "root",
        Password: "password",
        Database: "test",
    })
    
    // 创建存储实例
    userStore := &store.Store[User]{}
    userStore.SetDB(db)
    userStore.OnComplete() // 自动创建表
    
    ctx := context.Background()
    
    // 创建用户
    user := &User{
        UUID: "user-123",
        Name: "李四",
    }
    err := userStore.Save(ctx, user)
    if err != nil {
        panic(err)
    }
    
    // 查询用户
    found, err := userStore.Get(ctx, user.ID)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("查询到的用户: %+v\n", found)
    
    // 分页查询
    users, total, err := userStore.ListPageWhere(ctx, map[string]any{}, 1, 10, "id desc")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("用户列表: %+v, 总数: %d\n", users, total)
}
```

### 搜索功能使用示例

```go
package main

import (
    "context"
    
    "github.com/mengri/utils-store/store/search"
)

func main() {
    // 创建搜索存储
    searchStore := &search.Store[User]{}
    searchStore.SetDB(db)
    searchStore.OnComplete()
    
    ctx := context.Background()
    
    // 设置搜索标签
    err := searchStore.SetLabels(ctx, 1, "张三", "管理员", "北京")
    if err != nil {
        panic(err)
    }
    
    // 搜索
    results, err := searchStore.Search(ctx, "张三", map[string]interface{}{}, "id desc")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("搜索结果: %+v\n", results)
    
    // 分页搜索
    users, total, err := searchStore.SearchByPage(ctx, "管理员", nil, 1, 10, "id desc")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("搜索结果: %+v, 总数: %d\n", users, total)
}
```

## API文档

### 缓存模块

#### IKVCache接口
```go
type IKVCache[T any, K comparable] interface {
    Get(ctx context.Context, k K) (*T, error)
    Set(ctx context.Context, k K, t *T) error
    Delete(ctx context.Context, keys ...K) error
}
```

#### ICommonCache接口
```go
type ICommonCache interface {
    Get(ctx context.Context, key string) ([]byte, error)
    GetInt(ctx context.Context, key string) (int64, error)
    Del(ctx context.Context, keys ...string) error
    Set(ctx context.Context, key string, val []byte, expiration time.Duration) error
    HMSet(ctx context.Context, key string, value map[string][]byte, expiration time.Duration) error
    HGetAll(ctx context.Context, key string) (map[string]string, error)
    HDel(ctx context.Context, key string, fields ...string) error
    Incr(ctx context.Context, key string, expiration time.Duration) error
    IncrBy(ctx context.Context, key string, val int64, expiration time.Duration) error
    SetNX(ctx context.Context, key string, val interface{}, expiration time.Duration) (bool, error)
    Clone() ICommonCache
}
```

### 存储模块

#### IBaseStore接口
```go
type IBaseStore[T any] interface {
    Get(ctx context.Context, id int64) (*T, error)
    GetByUUID(ctx context.Context, uuid string) (*T, error)
    Save(ctx context.Context, t *T) error
    UpdateByUnique(ctx context.Context, t *T, uniques []string) error
    Delete(ctx context.Context, id ...int64) (int, error)
    UpdateWhere(ctx context.Context, w map[string]interface{}, m map[string]interface{}) (int64, error)
    Update(ctx context.Context, t *T) (int, error)
    UpdateField(ctx context.Context, field string, value interface{}, sql string, args ...interface{}) (int64, error)
    DeleteWhere(ctx context.Context, m map[string]interface{}) (int64, error)
    DeleteUUID(ctx context.Context, uuid string) error
    DeleteQuery(ctx context.Context, sql string, args ...interface{}) (int64, error)
    CountWhere(ctx context.Context, m map[string]interface{}) (int64, error)
    CountQuery(ctx context.Context, sql string, args ...interface{}) (int64, error)
    CountByGroup(ctx context.Context, m map[string]interface{}, group string) (map[string]int64, error)
    SoftDelete(ctx context.Context, where map[string]interface{}) error
    SoftDeleteQuery(ctx context.Context, sql string, args ...interface{}) error
    Insert(ctx context.Context, t ...*T) error
    List(ctx context.Context, m map[string]interface{}, order ...string) ([]*T, error)
    ListQuery(ctx context.Context, sql string, args []interface{}, order string) ([]*T, error)
    First(ctx context.Context, m map[string]interface{}, order ...string) (*T, error)
    FirstQuery(ctx context.Context, sql string, args []interface{}, order string) (*T, error)
    ListPage(ctx context.Context, sql string, pageNum, pageSize int, args []interface{}, order string) ([]*T, int64, error)
    ListPageWhere(ctx context.Context, where map[string]any, pageNum, pageSize int, order string) ([]*T, int64, error)
    Name() string
}
```

#### ISearchStore接口
```go
type ISearchStore[M any] interface {
    Search(ctx context.Context, keyword string, condition map[string]interface{}, sortRule ...string) ([]*M, error)
    SetLabels(ctx context.Context, id int64, label ...string) error
    Count(ctx context.Context, keyword string, condition map[string]interface{}) (int64, error)
    SearchByPage(ctx context.Context, keyword string, condition map[string]interface{}, page int, pageSize int, sortRule ...string) ([]*M, int64, error)
}
```

## 配置

### Redis配置
```go
type Config struct {
    Addr     string // Redis地址
    Password string // 密码
    DB       int    // 数据库编号
}
```

### MySQL配置
```go
type Config struct {
    Host     string // 主机地址
    Port     int    // 端口
    User     string // 用户名
    Password string // 密码
    Database string // 数据库名
}
```

## 项目结构

```
utils-store/
├── cache/                  # 缓存模块
│   ├── cache_redis/       # Redis缓存实现
│   ├── cache.go           # KV缓存接口
│   ├── common.go          # 通用缓存接口
│   ├── encode.go          # 编码解码
│   └── ...
├── store/                  # 存储模块
│   ├── store_mysql/       # MySQL存储实现
│   ├── search/            # 搜索功能
│   ├── history/           # 历史记录
│   ├── store.go           # 存储接口
│   ├── base.go            # 基础实现
│   └── ...
├── go.mod
└── README.md
```

## 特性说明

### 泛型支持
本项目充分利用Go 1.18+的泛型特性，提供类型安全的API接口。

### 事务支持
存储模块支持完整的事务管理，确保数据一致性。

### 自动迁移
支持数据库表的自动创建和迁移。

### 搜索功能
内置全文搜索功能，支持标签搜索和条件筛选。

### 软删除
支持软删除操作，数据不会被物理删除。

## 许可证

MIT License

## 贡献

欢迎提交Issue和Pull Request来改进这个项目。

## 版本历史

- v1.0.0: 初始版本，包含基础缓存和存储功能 
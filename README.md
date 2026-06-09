# ossx

对象存储客户端管理工具，基于 minio-go/v7，支持多数据源。

## 安装

```bash
go get github.com/go-xuan/ossx
```

## 快速开始

在 `conf/oss.yaml` 中配置：

```yaml
source: "default"
driver: "minio"
enable: true
endpoint: "127.0.0.1:9000"
accessKeyId: "minioadmin"
accessKeySecret: "minioadmin"
secure: false
bucket: "my-bucket"
```

```go
import "github.com/go-xuan/ossx"

func main() {
    ossx.Initialize()
    client := ossx.GetClient("default")
    url, _ := client.GetUrl(ctx, "path/to/file", time.Hour)
}
```

## 主要功能

- **多数据源** — 支持同时连接多个 MinIO 实例
- **完整操作** — CreateBucket / Upload / Download / Remove / Exist / GetUrl
- **可扩展** — 通过 RegisterClientBuilder 注册自定义后端（如 S3、COS）

## 扩展自定义后端

ossx 通过 Builder 注册机制支持扩展任意对象存储后端。内置 `minio`，可通过 `RegisterClientBuilder` 注册 S3、COS 等。

### 核心接口

所有后端必须实现 `Client` 接口（`client.go`）：

```go
type Client interface {
    GetInstance() interface{}
    GetConfig() *Config
    Close() error

    CreateBucket(ctx context.Context, name string, options ...interface{}) error
    Get(ctx context.Context, key string, options ...interface{}) (io.ReadCloser, error)
    Download(ctx context.Context, key string, options ...interface{}) error
    Upload(ctx context.Context, key string, reader io.Reader, options ...interface{}) error
    Exist(ctx context.Context, key string, options ...interface{}) (bool, error)
    Remove(ctx context.Context, key string, options ...interface{}) error
    GetUrl(ctx context.Context, key string, expires time.Duration, options ...interface{}) (string, error)
}
```

以及一个 `ClientBuilder` 构造函数：

```go
type ClientBuilder func(*Config) (Client, error)
```

### 示例：接入 S3

```go
package main

import (
    "context"
    "io"
    "time"

    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/go-xuan/ossx"
)

// 1. 定义 S3Client，实现 ossx.Client 接口
type S3Client struct {
    config *ossx.Config
    client *s3.Client
}

func (c *S3Client) GetInstance() interface{}  { return c.client }
func (c *S3Client) GetConfig() *ossx.Config   { return c.config }
func (c *S3Client) Close() error              { return nil }

func (c *S3Client) CreateBucket(ctx context.Context, name string, options ...interface{}) error {
    // 实现 S3 建桶逻辑
    return nil
}

func (c *S3Client) Upload(ctx context.Context, key string, reader io.Reader, options ...interface{}) error {
    // 实现 S3 上传逻辑
    return nil
}

// ... 实现其余接口方法

// 2. 定义 ClientBuilder
func S3ClientBuilder(config *ossx.Config) (ossx.Client, error) {
    cfg, _ := awsconfig.LoadDefaultConfig(ctx)
    client := s3.NewFromConfig(cfg)
    return &S3Client{config: config, client: client}, nil
}

// 3. init() 中注册
func init() {
    ossx.RegisterClientBuilder("s3", S3ClientBuilder)
}
```

注册后在 `conf/oss.yaml` 中通过 `driver` 字段指定：

```yaml
# MinIO 数据源
- source: "default"
  driver: "minio"
  enable: true
  endpoint: "127.0.0.1:9000"
  accessKeyId: "minioadmin"
  accessKeySecret: "minioadmin"

# S3 数据源
- source: "s3-backup"
  driver: "s3"
  enable: true
  endpoint: "s3.amazonaws.com"
  accessKeyId: "AKIAxxx"
  accessKeySecret: "xxx"
```

### Config 字段说明

`Config` 结构体对所有后端通用，未使用的字段可忽略：

```go
type Config struct {
    Source          string // 数据源名称（用于 GetClient("name") 定位）
    Builder         string // 客户端选型（minio/s3/cos/...，默认 minio）
    Enable          bool   // 是否启用
    Endpoint        string // 服务端点（MinIO/S3 地址）
    AccessKeyId     string // 访问密钥 ID
    AccessKeySecret string // 访问密钥
    AccessToken     string // 访问 Token（可选）
    Secure          bool   // 是否 HTTPS
    Bucket          string // 默认 Bucket
}
```

# tracfox

# 一款轻量级的HTTP(S)负载均衡器和微服务网关

- [x] 支持虚拟主机
- [x] 后端健康检查
- [x] 负载均衡
- [x] 插件机制, 支持自定义插件（如自定义微服务的鉴权插件等等）
- [x] TLS(https) 支持
- [x] 优雅关闭
- [ ] SNI 支持 (开发中)
- [ ] k8s ingress 支持（开发中）


# 编译

    go build -o tracfox cmd/tracfox/main.go


# 帮助

    # ./tracfox --help
    tracfox service, is the api gateway micro service component 

    Usage:
      tracfox [flags]

    Flags:
          --alsologtostderr                  log to standard error as well as files
          --config string                    The Tracfox Server will load its initial configuration from this file. The path may be absolute or relative; relative paths start at the Tracfox's current working directory. Omit this flag to use the built-in default configuration values. Command-line flags override configuration from this file. (default "/etc/tracfox/config.yaml")
          --dry-run                          If true, only check the configuration file and exit.
      -h, --help                             show more information about tracfox
          --log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
          --log_dir string                   If non-empty, write log files in this directory
          --logtostderr                      log to standard error instead of files
          --metrics string                   Metric address for tracfox server
          --stderrthreshold severity         logs at or above this threshold go to stderr (default 2)
      -v, --v Level                          log level for V logs
          --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging


# 检查配置文件正确性

    # ./tracfox  --config=examples/config.yaml --logtostderr  -v 10  --dry-run

# 运行


    # ./tracfox  --config=examples/config.yaml --logtostderr  -v 10

# 简单示例配置


    default:
      graceTimeOut: 100
      readTimeout: 100
      idleTimeout: 100
      writeTimeout: 100
      maxHeaderBytes: 10000
      readHeaderTimeout: 200

    frontends:
    - name: account
      address: localhost:8080
      # certificate: ssl/account.pem
      # certificateKey: ssl/account-key.pem
      virtualHosts:
        - domains: ["localhost"]
          filters:
          - name: cors
            config:
              allowHeaders: "Content-Type, Authorization"
              allowOrigin: "*"
              allowMethods: "GET, POST, PUT, DELETE, PATCH, OPTIONS"
          rules:
          - name: account
            locationRegexp: ^/v1/account/(.*)
            matchMethods: [GET,POST,PUT,DELETE,PATCH]
            rewitePath: /$1
            backend: account
            filters:
            - name: accessByAccount
              config:
                serverName: account
                address: 127.0.0.1:8082
                certificate: ssl/ca.pem

    backends:
    - name: account
      balance: roundrobin
      servers:
      - name: account1
        weight: 2
        failTimeout: 10
        maxFails: 2
        tcpTimeout: 100
        tcpKeepAlive: 100
        idleConnTimeout: 100
        maxIdleConnsPerHost: 100
        target: https://127.0.0.1:8083/
        

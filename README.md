# tracfox

# 一款轻量级的HTTP(S)负载均衡器和微服务网关

- [x] 支持虚拟主机
- [x] 后端健康检查
- [x] 负载均衡
- [x] 插件机制, 支持自定义插件（如自定义微服务的鉴权插件等等）
- [x] TLS(https) 支持
- [ ] SNI 支持 (开发中)


# 简单示例配置


    default:
      metricAddr: localhost:8000
      graceTimeOut: 100
      readTimeout: 100
      idleTimeout: 100
      writeTimeout: 100
      maxHeaderBytes: 10000
      readHeaderTimeout: 200

    frontends:
    - name: account
      address: localhost:8080
      # certificate: ssl/tracfox.pem
      # certificateKey: ssl/tracfox-key.pem
      virtualHosts:
        - domains: [ "localhost", "api.labchan.com" ]
          rules:
          - name: account
            locationRegexp: ^/v1/account/(.*)
            matchMethods: [GET,POST,PUT,DELETE,PATCH]
            rewitePath: /$1
            backend: account
            plugins:
            - name: cors
              config:
                allowHeaders: "Content-Type,Accept"
                allowOrigin: "*"
                allowMethods: "GET"
            - name: accessByAccount
              config:
                serverName: account
                address: localhost:8082
                certificate: /Users/aapeli/labchan/ssl/ca.pem

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
        target: http://127.0.0.1:8083/

      - name: account2
        weight: 2
        maxFails: 2
        failTimeout: 10
        tcpTimeout: 100
        tcpKeepAlive: 100
        idleConnTimeout: 100
        maxIdleConnsPerHost: 100
        target: http://127.0.0.1:8084/


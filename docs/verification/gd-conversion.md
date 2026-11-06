# 转换来源与验证

GD源码基线：a86d1c0ec1635b24a09010fc06f29c3458cf80f4。仅导出提交源码到隔离目录，不带.env/数据/历史；GD未提交工作未迁入。Plus规划先提交4fb6f39。

源测试：隔离副本go test ./...及go vet ./...退出0；DB/向量/模型集成测试需独立服务，不能据此称已全通过。Plus实施前make check退出0。原始日志留在维护者仓库外隔离目录。

第一阶段：密码哈希与JWT校验测试转换后通过；独立Compose MySQL/Go/Nginx的API三组测试通过，覆盖退出撤销、跨班404、学生403、审计读授权以及审计故障下业务回滚。后续材料/检索/AI/MCP仍待验证。

来源摘要：

```json
[
  {
    "path": "cmd/api/auth.go",
    "sha256": "2bfa82d36a88a016e8ff3a0185a8b02cb94f3f71a7cf8ffa3032772dec570be3"
  },
  {
    "path": "cmd/api/password.go",
    "sha256": "d9990482217a43f8e984a3e0ca4518990ccbc7438422d30f314a3a523f8d6b4f"
  },
  {
    "path": "cmd/api/chunks.go",
    "sha256": "08cc6a928a105ad2285641470fccad4233671367c56de869331a6c474308c197"
  },
  {
    "path": "cmd/api/retrieval.go",
    "sha256": "f70b72b799ccbc8bb1ec32091abd09af72d0760dad3dfba2c616552e520c1a68"
  },
  {
    "path": "cmd/api/agents.go",
    "sha256": "1af21590a354af1df4323b58c7e593df2b36e229c9c7f1445495675072e4960d"
  },
  {
    "path": "cmd/api/skills.go",
    "sha256": "983cceee47eb2620675135fd40e0043896c450fa2eee58248c5c2671b36578fd"
  },
  {
    "path": "cmd/api/conversations.go",
    "sha256": "d33f078297ff04624c9f65b8822788d116d4b696df547c754cc9eac88443f071"
  },
  {
    "path": "cmd/api/stream.go",
    "sha256": "856976db8e8c5e02d4f02a9e8d1bb852b2da862546d54d31174c190958b6f3ad"
  },
  {
    "path": "cmd/api/mcp.go",
    "sha256": "afb05baff869814e902dc04f7dcb849c90fd7e321f280e55635cedb2863df6bb"
  }
]
```

结构变化：旧main包的HTTP/SQL拆为httpapi→service→domain与repository适配器；统一snake_case。独立人类评审、异机复验未执行。

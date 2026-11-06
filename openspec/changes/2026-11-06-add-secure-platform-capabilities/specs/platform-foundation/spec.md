## MODIFIED Requirements

### Requirement: F2 Minimal public API and default deny

**最小公开接口与默认拒绝**

系统 SHALL 显式登记公开入口：工程GET `/health`、GET `/live`、页面`/`、GET `/openapi.json`与`/swagger/`文档资产，以及POST `/api/auth/login`。其他已注册业务入口必须Bearer鉴权。
- 未知API和附件静态路径返回404；health/live非GET返回405，且不接受query。
- 登录/资源授权按secure-platform-capabilities验收，不能由工程探活通过推导业务安全。

#### Scenario: S04 Invalid filter failure
- **GIVEN** 前置条件：系统已就绪。
- **WHEN** 操作：GET `/health?class_id=other`。
- **THEN** 结果：400/INVALID_ARGUMENT，无额外依赖详情或数据，也无数据库写入。

#### Scenario: S05 Unregistered protected operation denied
- **GIVEN** 前置条件：匿名用户没有业务授权。
- **WHEN** 操作：GET `/api/teachers/materials`。
- **THEN** 结果：404/NOT_FOUND，不返回教师数据。
- **AND** 这只验证未登记业务默认不可访问，不替代独立的JWT/班级授权测试。

#### Scenario: S06 Missing resource failure
- **GIVEN** 前置条件：系统运行中。
- **WHEN** 操作：GET `/uploads/guessed.txt`或GET `/api/missing`。
- **THEN** 结果：404，不返回文件、前端index或私有数据。

#### Scenario: S07 Unsupported method failure
- **GIVEN** 前置条件：系统运行中。
- **WHEN** 操作：POST `/health`。
- **THEN** 结果：405/METHOD_NOT_ALLOWED，Allow包含GET，没有写入副作用。

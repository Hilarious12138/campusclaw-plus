## ADDED Requirements

### Requirement: Public README boundary
README SHALL 不包含个人姓名和本机配置逐项填写教程，并提供可用的部署文档与验收入口。

#### Scenario: Read project documentation
- **GIVEN** 用户打开仓库README
- **WHEN** 阅读并点击部署说明
- **THEN** README不展示个人姓名或配置填写表，部署说明文件存在，Docker部署与验收命令仍可找到。

# README内容边界（public-documentation）

## Purpose

让公开入口简洁说明项目用途、真实交付范围、部署与验收方式。个人姓名和逐项配置填写教程不放README，详细部署要求由独立文档提供。

## Requirements

### Requirement: Public README boundary

**入口简洁，说明可查**

README SHALL 不含个人姓名或本机配置逐项填写教程，并提供有效的部署文档和验收入口。

#### Scenario: Read project documentation
- **GIVEN** 前置条件：用户打开README。
- **WHEN** 操作：阅读并点击部署说明。
- **THEN** 结果：README没有姓名/配置填写表，部署文件存在，Docker部署和验收命令可找到。

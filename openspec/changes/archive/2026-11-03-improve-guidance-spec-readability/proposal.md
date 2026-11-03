## Why

AGENTS的长段落混合了操作顺序、边界与索引；长期能力规约也常把多个验收结果写在同一句中。用户要求继续采用已确认的短句、分节和简短解释方式，优化总纲及具体能力文件，使人和AI都容易找到执行步骤与判定依据。

## What Changes

- 重组AGENTS为项目范围、阅读顺序、实施步骤、核心边界及验证交付，保留八类规范索引和全部既有约束。
- 改写openspec/specs下7个长期能力文件，增加中文标题/要求概括，拆分多项要求及验收结果。
- 保留Purpose/Requirements/Requirement/Scenario和GIVEN/WHEN/THEN等英文骨架、原Requirement/Scenario标题与编号，避免破坏定位。
- 将现有standards-readability的表达约束扩展至AGENTS及长期能力规约；不新增平行规范或业务能力。

## Capabilities

### New Capabilities
无。

### Modified Capabilities
- `standards-readability`: 统一总纲、工程规范和能力规约的表达方式，并保留各自格式与职责。

## Non-goals

不修改业务代码、配置、技术栈和八份已完成的工程规范，不改变接口/错误码/数值/失败条件，不把未验证事项写成已完成，不重写旧归档原文，不增加个人用途公开资料。

## Alternatives

放弃只替换少量术语：仍难找到步骤。也不把规范内容复制进AGENTS，避免重复正文和额外上下文。

## Impact

范围为AGENTS和7个长期spec，以及本次规划/核对记录。用户已明确要求改写，沿用此前确认的表达风格；推送GitHub前说明实际改动及原因，仍经PR和check/grade门禁。

## Context

原product/backend/security/testing/operations五份规则已存在，AI/data/review从AGENTS迁出形成八份独立文件。AGENTS保留核心边界、读取索引和四文档优先流程。

## Decisions

一份条款保留一份权威正文；只读取相关规则。文档写明适用条件、动作与检查/失败判据。八份文件齐全不代表业务评测、数据恢复或人类评审已完成。

## Verification

文件非空、AGENTS引用路径、Skill元数据、迁移正文、strict、既有门禁与验证记录已检查。AI/数据业务仍待对应change实施和验证；旧执行证据只证明对应提交。

## Current documentation

项目导航由README与docs/roadmap.md提供，验证摘要位于docs/verification/platform-foundation.md；个人材料与产品文档分开保存。原始执行上下文可由Git提交追溯，当前表述整理不改原始结果。

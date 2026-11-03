## Decisions

项目路线图只维护阶段能力、建设目的、审核重点和完成证据；具体设计留在对应change。review规则与模板统一七类/七维、高风险依据、逐文件分类及阻断条件。

## Review boundary

自动检查负责可重复判断事项，人工复核权限、并发、事务和一致性。自审、AI辅助、CI、独立人类评审分别记录，未执行项保持未完成。例外须维护者明确记录理由与补齐条件。

## Stack / Verification

技术栈与业务运行行为不变。文档链接、静态/历史检查、make check和strict实际通过后归档；原始CI/提交记录可追溯。公开文档归入README、docs/roadmap.md、docs/reviews与docs/verification，外部个人材料单独保存。

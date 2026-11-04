## Context
此变更仅调整工程文档。用户已批准上一轮提出的全流程导航，不重复索取确认。

## Decisions
在AGENTS现有“每次变更怎样执行”前添加一段导航，引用review与verify-change；不复制八份规范正文。长期能力新增E5与两个可判定场景。

## Verification
核对导航覆盖阶段和现有规则一致；git diff --check、静态检查、OpenSpec strict、make check及PR现有CI。无业务或依赖变化，不新增运行测试与依赖；已有CI照常运行。审核记录保留独立评审未完成，不改变审批规则。

## Risks
其他会话有未提交变更，本次仅显式暂存自己的文件。

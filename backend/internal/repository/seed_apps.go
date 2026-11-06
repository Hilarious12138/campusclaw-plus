package repository

import (
	"context"
	"encoding/json"
)

func (m *MySQL) SeedApplications(ctx context.Context) error {
	var done int
	if e := m.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version='applications-v1'").Scan(&done); e != nil {
		return e
	}
	if done > 0 {
		return nil
	}
	tx, e := m.db.BeginTx(ctx, nil)
	if e != nil {
		return e
	}
	defer tx.Rollback()
	for _, letter := range []string{"a", "b"} {
		u, e := m.LoginUser(ctx, "teacher_"+letter)
		if e != nil {
			return e
		}
		skills := []string{}
		for _, name := range []string{"数学解题引导", "语文阅读引导"} {
			id := stableID(u.ClassID + name)
			skills = append(skills, id)
			conditions := []string{"方程", "面积", "等式", "未知数"}
			if name == "语文阅读引导" {
				conditions = []string{"语文", "阅读", "诗歌"}
			}
			body, _ := json.Marshal(map[string]any{"name": name, "conditions": conditions, "instructions": "先询问学生已知条件和已尝试步骤，再给下一步提示，不直接给最终答案。引用本班资料说明方法。", "tools": []string{"knowledge_search"}, "enabled": true})
			if _, e = tx.ExecContext(ctx, "INSERT INTO resources VALUES(?,?,?,'skills',1,?,UTC_TIMESTAMP(6),UTC_TIMESTAMP(6)) ON DUPLICATE KEY UPDATE id=id", id, u.ClassID, u.ID, body); e != nil {
				return e
			}
		}
		for i, name := range []string{"知识问答", "数学助手", "语文助手"} {
			bound := []string{}
			mode := "knowledge"
			prompt := "你是本班知识库助手，根据资料作答。没有依据时明确未找到。"
			if i > 0 {
				bound = []string{skills[i-1]}
				mode = "guidance"
				prompt = "你是本班学习引导助手，用简体中文清楚回答，只引导下一步，不直接给最终答案。"
			}
			body, _ := json.Marshal(map[string]any{"name": name, "prompt": prompt, "mode": mode, "skill_ids": bound, "material_ids": []string{}, "enabled": true})
			if _, e = tx.ExecContext(ctx, "INSERT INTO resources VALUES(?,?,?,'agents',1,?,UTC_TIMESTAMP(6),UTC_TIMESTAMP(6)) ON DUPLICATE KEY UPDATE id=id", stableID(u.ClassID+name), u.ClassID, u.ID, body); e != nil {
				return e
			}
		}
	}
	if _, e = tx.ExecContext(ctx, "INSERT INTO schema_migrations VALUES('applications-v1',UTC_TIMESTAMP(6)) ON DUPLICATE KEY UPDATE version=version"); e != nil {
		return e
	}
	return tx.Commit()
}

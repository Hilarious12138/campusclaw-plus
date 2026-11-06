package service

import (
	"campusclaw-plus/internal/domain"
	"strings"
	"unicode/utf8"
)

func (s Conversations) ParseSkill(raw string) (SkillBody, error) {
	v := SkillBody{Tools: []string{}, Conditions: []string{}}
	if len(raw) > 65536 || !utf8.ValidString(raw) {
		return v, domain.ErrInvalid
	}
	raw = strings.ReplaceAll(raw, "\r\n", "\n")
	parts := strings.SplitN(raw, "\n\n", 2)
	if len(parts) != 2 {
		return v, domain.ErrInvalid
	}
	seen := map[string]bool{}
	for _, line := range strings.Split(parts[0], "\n") {
		p := strings.SplitN(line, "：", 2)
		if len(p) != 2 {
			p = strings.SplitN(line, ":", 2)
		}
		if len(p) != 2 {
			return v, domain.ErrInvalid
		}
		k, value := strings.TrimSpace(p[0]), strings.TrimSpace(p[1])
		if seen[k] {
			return v, domain.ErrInvalid
		}
		seen[k] = true
		switch k {
		case "名称", "name":
			v.Name = value
		case "适用条件", "conditions":
			for _, s := range strings.FieldsFunc(value, func(r rune) bool { return r == ',' || r == '，' }) {
				v.Conditions = append(v.Conditions, strings.TrimSpace(s))
			}
		case "工具", "tools":
			for _, s := range strings.FieldsFunc(value, func(r rune) bool { return r == ',' || r == '，' }) {
				v.Tools = append(v.Tools, strings.TrimSpace(s))
			}
		case "启用", "enabled":
			if value != "是" && value != "否" && value != "true" && value != "false" {
				return v, domain.ErrInvalid
			}
			v.Enabled = value == "是" || value == "true"
		default:
			return v, domain.ErrInvalid
		}
	}
	v.Instructions = strings.TrimSpace(parts[1])
	if v.Name == "" || len([]rune(v.Name)) > 100 || v.Instructions == "" || len(v.Instructions) > 32000 || len(v.Tools) > 2 {
		return v, domain.ErrInvalid
	}
	v.Enabled = false
	return v, nil
}

package utils

import "student-attendance/internal/pkg/code"

// GetSubjectCode 根据传来的字符串返回学科代码
func GetSubjectCode(subject string) int8 {
	switch subject {
	case "语文":
		return code.Chinese
	case "数学":
		return code.Mathmatics
	case "英语":
		return code.English
	case "物理":
		return code.Physics
	case "化学":
		return code.Chemistry
	case "生物":
		return code.Biology
	case "道法":
		return code.Politics
	case "历史":
		return code.History
	case "地理":
		return code.Geography
	case "心理":
		return code.Psychology
	case "音乐":
		return code.Music
	case "信息":
		return code.ComputerScience
	case "美术":
		return code.Arts
	case "体育":
		return code.Sports
	case "劳动教育":
		return code.LaborEducation
	case "班会":
		return code.ClassMeeting
	default:
		return 0
	}
}

// GetSubjectFromCode 根据学科代码返回学科名字
func GetSubjectFromCode(num int8) string {
	switch num {
	case code.Chinese:
		return "语文"
	case code.Mathmatics:
		return "数学"
	case code.English:
		return "英语"
	case code.Physics:
		return "物理"
	case code.Chemistry:
		return "化学"
	case code.Biology:
		return "生物"
	case code.Politics:
		return "道法"
	case code.History:
		return "历史"
	case code.Geography:
		return "地理"
	case code.Psychology:
		return "心理"
	case code.Music:
		return "音乐"
	case code.ComputerScience:
		return "信息"
	case code.LaborEducation:
		return "劳动教育"
	case code.Arts:
		return "美术"
	case code.Sports:
		return "体育"
	case code.ClassMeeting:
		return "班会"
	default:
		return ""
	}
}

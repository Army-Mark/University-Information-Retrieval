package models

// University 学校模型
// 完整的学校信息结构体，包含学校的所有详细信息
type University struct {
	// SchoolID 学校ID
	SchoolID string `json:"学校ID"`
	// SchoolName 学校名称
	SchoolName string `json:"学校名称"`
	// Address 地址
	Address string `json:"地址"`
	// Category 类别
	Category string `json:"类别"`
	// Nature 性质
	Nature string `json:"性质"`
	// Affiliation 归属部门
	Affiliation string `json:"归属部门"`
	// Tags 标签
	Tags string `json:"标签"`
	// FoundedYear 建校时间
	FoundedYear string `json:"建校时间"`
	// Area 占地面积
	Area string `json:"占地面积"`
	// BaoyanStar 保研星级
	BaoyanStar string `json:"保研星级"`
	// PhDPrograms 博士点数量
	PhDPrograms string `json:"博士点数量"`
	// MasterPrograms 硕士点数量
	MasterPrograms string `json:"硕士点数量"`
	// NationalKeyDisciplines 国家重点学科数量
	NationalKeyDisciplines string `json:"国家重点学科数量"`
	// SoftScienceRank 软科综合排名
	SoftScienceRank string `json:"软科综合排名"`
	// AlumniRank 校友会综合排名
	AlumniRank string `json:"校友会综合排名"`
	// QSRank QS世界排名
	QSRank string `json:"QS世界排名"`
	// USNewsRank US世界排名
	USNewsRank string `json:"US世界排名"`
	// TimesRank 泰晤士排名
	TimesRank string `json:"泰晤士排名"`
	// PopularityRank 人气值排名
	PopularityRank string `json:"人气值排名"`
	// BasicInfo 基本信息
	BasicInfo string `json:"基本信息"`
	// SchoolForm 办学形式
	SchoolForm string `json:"办学形式"`
	// LogoPath  logo路径
	LogoPath string `json:"logo_path"`
}

// ScrollingUniversity 滚动显示的学校信息
// 用于页面滚动展示的简化学校信息结构体
type ScrollingUniversity struct {
	// SchoolID 学校ID
	SchoolID string `json:"学校ID"`
	// SchoolName 学校名称
	SchoolName string `json:"学校名称"`
	// Address 地址
	Address string `json:"地址"`
	// Category 类别
	Category string `json:"类别"`
	// Nature 性质
	Nature string `json:"性质"`
}

// SearchResult 搜索结果
// 搜索学校时返回的结果结构体
type SearchResult struct {
	// ID 学校ID
	ID string `json:"id"`
	// Name 学校名称
	Name string `json:"name"`
	// MatchType 匹配类型 (exact_id, partial_id, name)
	MatchType string `json:"match_type"`
}

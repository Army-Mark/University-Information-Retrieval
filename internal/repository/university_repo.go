package repository

import (
	"database/sql"
	"school-go/internal/models"
)

// UniversityRepository 学校数据仓库接口
// 定义了学校数据的CRUD操作和搜索功能
// 实现了与数据库的交互逻辑

type UniversityRepository interface {
	// GetAll 获取所有学校数据
	// 返回学校列表和可能的错误
	GetAll() ([]models.University, error)

	// GetByID 根据学校ID获取学校信息
	// id: 学校ID
	// 返回学校信息和可能的错误
	GetByID(id string) (*models.University, error)

	// Create 创建新学校
	// university: 学校信息
	// 返回可能的错误
	Create(university *models.University) error

	// Update 更新学校信息
	// university: 学校信息
	// 返回可能的错误
	Update(university *models.University) error

	// Delete 删除学校
	// id: 学校ID
	// 返回可能的错误
	Delete(id string) error

	// GetScrollingData 获取滚动数据
	// 返回滚动数据列表和可能的错误
	GetScrollingData() ([]models.ScrollingUniversity, error)

	// SaveScrollingPosition 保存滚动位置
	// position: 滚动位置
	// 返回可能的错误
	SaveScrollingPosition(position int) error

	// GetScrollingPosition 获取滚动位置
	// 返回滚动位置和可能的错误
	GetScrollingPosition() (int, error)

	// Search 搜索学校
	// keyword: 搜索关键词
	// 返回搜索结果和可能的错误
	Search(keyword string) ([]models.SearchResult, error)

	// SearchWithPagination 分页搜索学校
	// keyword: 搜索关键词
	// page: 页码
	// pageSize: 每页大小
	// 返回搜索结果、总数和可能的错误
	SearchWithPagination(keyword string, page, pageSize int) ([]models.SearchResult, int, error)
}

// universityRepository 学校数据仓库实现
// 基于SQLite数据库

type universityRepository struct {
	// db 数据库连接
	db *sql.DB
}

// NewUniversityRepository 创建学校数据仓库实例
// 返回 UniversityRepository 接口实现
func NewUniversityRepository() UniversityRepository {
	return &universityRepository{
		db: GetDB(),
	}
}

// selectUniversities 学校表查询字段
const selectUniversities = `学校ID, 学校名称, 地址, 类别, 性质, 归属部门, 标签, 建校时间, 占地面积, 保研星级, 博士点数量, 硕士点数量, 国家重点学科数量, 软科综合排名, 校友会综合排名, QS世界排名, US世界排名, 泰晤士排名, 人气值排名, 基本信息, 办学形式, logo_path`

// scanUniversity 从数据库扫描器中解析学校信息
// scanner: 数据库扫描器
// 返回学校信息和可能的错误
func scanUniversity(scanner interface{ Scan(...interface{}) error }) (*models.University, error) {
	var u models.University
	err := scanner.Scan(
		&u.SchoolID, &u.SchoolName, &u.Address, &u.Category, &u.Nature,
		&u.Affiliation, &u.Tags, &u.FoundedYear, &u.Area, &u.BaoyanStar,
		&u.PhDPrograms, &u.MasterPrograms, &u.NationalKeyDisciplines,
		&u.SoftScienceRank, &u.AlumniRank, &u.QSRank, &u.USNewsRank,
		&u.TimesRank, &u.PopularityRank, &u.BasicInfo, &u.SchoolForm, &u.LogoPath,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetAll 获取所有学校数据
func (r *universityRepository) GetAll() ([]models.University, error) {
	// 先获取总记录数，用于预分配切片容量
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM universities").Scan(&count)
	if err != nil {
		return nil, err
	}

	// 执行查询
	rows, err := r.db.Query("SELECT " + selectUniversities + " FROM universities")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 预分配切片容量
	universities := make([]models.University, 0, count)
	// 遍历结果集
	for rows.Next() {
		u, err := scanUniversity(rows)
		if err != nil {
			return nil, err
		}
		universities = append(universities, *u)
	}

	return universities, nil
}

// GetByID 根据学校ID获取学校信息
func (r *universityRepository) GetByID(id string) (*models.University, error) {
	// 执行查询
	row := r.db.QueryRow("SELECT "+selectUniversities+" FROM universities WHERE 学校ID = ?", id)
	u, err := scanUniversity(row)
	if err != nil {
		if err == sql.ErrNoRows {
			// 未找到记录
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

// Create 创建新学校
func (r *universityRepository) Create(u *models.University) error {
	// 执行插入操作
	_, err := r.db.Exec(
		"INSERT INTO universities ("+selectUniversities+") VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		u.SchoolID, u.SchoolName, u.Address, u.Category, u.Nature,
		u.Affiliation, u.Tags, u.FoundedYear, u.Area, u.BaoyanStar,
		u.PhDPrograms, u.MasterPrograms, u.NationalKeyDisciplines,
		u.SoftScienceRank, u.AlumniRank, u.QSRank, u.USNewsRank,
		u.TimesRank, u.PopularityRank, u.BasicInfo, u.SchoolForm, u.LogoPath,
	)
	return err
}

// Update 更新学校信息
func (r *universityRepository) Update(u *models.University) error {
	// 执行更新操作
	_, err := r.db.Exec(
		"UPDATE universities SET 学校名称=?,地址=?,类别=?,性质=?,归属部门=?,标签=?,建校时间=?,占地面积=?,保研星级=?,博士点数量=?,硕士点数量=?,国家重点学科数量=?,软科综合排名=?,校友会综合排名=?,QS世界排名=?,US世界排名=?,泰晤士排名=?,人气值排名=?,基本信息=?,办学形式=?,logo_path=? WHERE 学校ID=?",
		u.SchoolName, u.Address, u.Category, u.Nature,
		u.Affiliation, u.Tags, u.FoundedYear, u.Area, u.BaoyanStar,
		u.PhDPrograms, u.MasterPrograms, u.NationalKeyDisciplines,
		u.SoftScienceRank, u.AlumniRank, u.QSRank, u.USNewsRank,
		u.TimesRank, u.PopularityRank, u.BasicInfo, u.SchoolForm, u.LogoPath,
		u.SchoolID,
	)
	return err
}

// Delete 删除学校
func (r *universityRepository) Delete(id string) error {
	// 执行删除操作
	_, err := r.db.Exec("DELETE FROM universities WHERE 学校ID = ?", id)
	return err
}

// GetScrollingData 获取滚动数据
func (r *universityRepository) GetScrollingData() ([]models.ScrollingUniversity, error) {
	// 先获取总记录数，用于预分配切片容量
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM universities").Scan(&count)
	if err != nil {
		return nil, err
	}

	// 执行查询
	rows, err := r.db.Query("SELECT 学校ID, 学校名称, 地址, 类别, 性质 FROM universities")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 预分配切片容量
	list := make([]models.ScrollingUniversity, 0, count)
	// 遍历结果集
	for rows.Next() {
		var s models.ScrollingUniversity
		if err := rows.Scan(&s.SchoolID, &s.SchoolName, &s.Address, &s.Category, &s.Nature); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, nil
}

// SaveScrollingPosition 保存滚动位置
func (r *universityRepository) SaveScrollingPosition(position int) error {
	// 确保表存在
	if _, err := r.db.Exec("CREATE TABLE IF NOT EXISTS scrolling_position (id INTEGER PRIMARY KEY, position INTEGER)"); err != nil {
		return err
	}
	// 插入或更新位置
	_, err := r.db.Exec("INSERT OR REPLACE INTO scrolling_position (id, position) VALUES (1, ?)", position)
	return err
}

// GetScrollingPosition 获取滚动位置
func (r *universityRepository) GetScrollingPosition() (int, error) {
	// 确保表存在
	if _, err := r.db.Exec("CREATE TABLE IF NOT EXISTS scrolling_position (id INTEGER PRIMARY KEY, position INTEGER)"); err != nil {
		return 0, err
	}
	var pos int
	// 查询位置
	err := r.db.QueryRow("SELECT position FROM scrolling_position WHERE id = 1").Scan(&pos)
	if err != nil {
		if err == sql.ErrNoRows {
			// 未找到记录，返回默认值0
			return 0, nil
		}
		return 0, err
	}
	return pos, nil
}

// Search 搜索学校
// 实现了SQL级别的搜索，支持精确ID匹配、部分ID匹配、名称匹配和其他字段匹配
func (r *universityRepository) Search(keyword string) ([]models.SearchResult, error) {
	// 构建模糊匹配字符串
	keywordLike := "%" + keyword + "%"

	// 先获取总匹配数，用于预分配切片容量
	var count int
	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM universities WHERE 学校ID = ? OR 学校ID LIKE ? OR 学校名称 LIKE ? OR 地址 LIKE ? OR 类别 LIKE ? OR 性质 LIKE ? OR 归属部门 LIKE ?",
		keyword, keywordLike, keywordLike, keywordLike, keywordLike, keywordLike, keywordLike,
	).Scan(&count)
	if err != nil {
		return nil, err
	}

	// 1. 精确ID匹配
	exactIDRows, err := r.db.Query("SELECT 学校ID, 学校名称 FROM universities WHERE 学校ID = ?", keyword)
	if err != nil {
		return nil, err
	}
	defer exactIDRows.Close()

	// 预分配切片容量
	results := make([]models.SearchResult, 0, count)
	// 处理精确ID匹配结果
	for exactIDRows.Next() {
		var id, name string
		if err := exactIDRows.Scan(&id, &name); err != nil {
			return nil, err
		}
		// 精确匹配结果排在最前面
		results = append([]models.SearchResult{{ID: id, Name: name, MatchType: "exact_id"}}, results...)
	}

	// 2. 部分ID匹配
	partialIDRows, err := r.db.Query("SELECT 学校ID, 学校名称 FROM universities WHERE 学校ID LIKE ? AND 学校ID != ?", keywordLike, keyword)
	if err != nil {
		return nil, err
	}
	defer partialIDRows.Close()

	// 处理部分ID匹配结果
	for partialIDRows.Next() {
		var id, name string
		if err := partialIDRows.Scan(&id, &name); err != nil {
			return nil, err
		}
		results = append(results, models.SearchResult{ID: id, Name: name, MatchType: "partial_id"})
	}

	// 3. 名称匹配
	nameRows, err := r.db.Query("SELECT 学校ID, 学校名称 FROM universities WHERE 学校名称 LIKE ? AND 学校ID NOT LIKE ?", keywordLike, keywordLike)
	if err != nil {
		return nil, err
	}
	defer nameRows.Close()

	// 处理名称匹配结果
	for nameRows.Next() {
		var id, name string
		if err := nameRows.Scan(&id, &name); err != nil {
			return nil, err
		}
		results = append(results, models.SearchResult{ID: id, Name: name, MatchType: "name"})
	}

	// 4. 其他字段匹配（地址、类别、性质、归属部门）
	otherRows, err := r.db.Query(
		"SELECT 学校ID, 学校名称 FROM universities WHERE (地址 LIKE ? OR 类别 LIKE ? OR 性质 LIKE ? OR 归属部门 LIKE ?) AND 学校ID != ? AND 学校ID NOT LIKE ? AND 学校名称 NOT LIKE ?",
		keywordLike, keywordLike, keywordLike, keywordLike, keyword, keywordLike, keywordLike,
	)
	if err != nil {
		return nil, err
	}
	defer otherRows.Close()

	// 处理其他字段匹配结果
	for otherRows.Next() {
		var id, name string
		if err := otherRows.Scan(&id, &name); err != nil {
			return nil, err
		}
		results = append(results, models.SearchResult{ID: id, Name: name, MatchType: "other"})
	}

	return results, nil
}

// SearchWithPagination 分页搜索学校
// 实现了带分页的SQL级别搜索，支持精确ID匹配、部分ID匹配、名称匹配和其他字段匹配
func (r *universityRepository) SearchWithPagination(keyword string, page, pageSize int) ([]models.SearchResult, int, error) {
	// 构建模糊匹配字符串
	keywordLike := "%" + keyword + "%"

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 计算总结果数
	var total int
	totalRows, err := r.db.Query(
		"SELECT COUNT(*) FROM universities WHERE 学校ID = ? OR 学校ID LIKE ? OR 学校名称 LIKE ? OR 地址 LIKE ? OR 类别 LIKE ? OR 性质 LIKE ? OR 归属部门 LIKE ?",
		keyword, keywordLike, keywordLike, keywordLike, keywordLike, keywordLike, keywordLike,
	)
	if err != nil {
		return nil, 0, err
	}
	defer totalRows.Close()

	if totalRows.Next() {
		totalRows.Scan(&total)
	}

	// 1. 精确ID匹配
	exactIDRows, err := r.db.Query("SELECT 学校ID, 学校名称 FROM universities WHERE 学校ID = ?", keyword)
	if err != nil {
		return nil, 0, err
	}
	defer exactIDRows.Close()

	// 预分配切片容量，最多为 pageSize + 1（精确匹配可能有1条）
	results := make([]models.SearchResult, 0, pageSize+1)
	// 处理精确ID匹配结果
	for exactIDRows.Next() {
		var id, name string
		if err := exactIDRows.Scan(&id, &name); err != nil {
			return nil, 0, err
		}
		// 精确匹配结果排在最前面
		results = append([]models.SearchResult{{ID: id, Name: name, MatchType: "exact_id"}}, results...)
	}

	// 2. 部分ID匹配（带分页）
	partialIDRows, err := r.db.Query(
		"SELECT 学校ID, 学校名称 FROM universities WHERE 学校ID LIKE ? AND 学校ID != ? LIMIT ? OFFSET ?",
		keywordLike, keyword, pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer partialIDRows.Close()

	// 处理部分ID匹配结果
	for partialIDRows.Next() {
		var id, name string
		if err := partialIDRows.Scan(&id, &name); err != nil {
			return nil, 0, err
		}
		results = append(results, models.SearchResult{ID: id, Name: name, MatchType: "partial_id"})
	}

	// 3. 名称匹配（带分页）
	nameRows, err := r.db.Query(
		"SELECT 学校ID, 学校名称 FROM universities WHERE 学校名称 LIKE ? AND 学校ID NOT LIKE ? LIMIT ? OFFSET ?",
		keywordLike, keywordLike, pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer nameRows.Close()

	// 处理名称匹配结果
	for nameRows.Next() {
		var id, name string
		if err := nameRows.Scan(&id, &name); err != nil {
			return nil, 0, err
		}
		results = append(results, models.SearchResult{ID: id, Name: name, MatchType: "name"})
	}

	// 4. 其他字段匹配（地址、类别、性质、归属部门）
	otherRows, err := r.db.Query(
		"SELECT 学校ID, 学校名称 FROM universities WHERE (地址 LIKE ? OR 类别 LIKE ? OR 性质 LIKE ? OR 归属部门 LIKE ?) AND 学校ID != ? AND 学校ID NOT LIKE ? AND 学校名称 NOT LIKE ? LIMIT ? OFFSET ?",
		keywordLike, keywordLike, keywordLike, keywordLike, keyword, keywordLike, keywordLike, pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer otherRows.Close()

	// 处理其他字段匹配结果
	for otherRows.Next() {
		var id, name string
		if err := otherRows.Scan(&id, &name); err != nil {
			return nil, 0, err
		}
		results = append(results, models.SearchResult{ID: id, Name: name, MatchType: "other"})
	}

	return results, total, nil
}

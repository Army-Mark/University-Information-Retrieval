package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"school-go/internal/models"
	"school-go/internal/pkg/errors"
	"school-go/internal/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// sanitizeInput 清理输入，防止XSS攻击
func sanitizeInput(input string) string {
	// 替换HTML特殊字符
	replacements := map[string]string{
		"<":  "&lt;",
		">":  "&gt;",
		"&":  "&amp;",
		"\"": "&quot;",
		"'":  "&#39;",
	}
	result := input
	for old, new := range replacements {
		result = strings.ReplaceAll(result, old, new)
	}
	return result
}

// validateInput 验证输入，防止SQL注入
func validateInput(input string) bool {
	// 检查是否包含SQL注入攻击特征
	sqlInjectionPatterns := []string{
		"'", "\"", ";", "--", "OR", "AND", "UNION", "SELECT", "INSERT", "UPDATE", "DELETE", "DROP", "CREATE", "ALTER",
	}
	lowerInput := strings.ToLower(input)
	for _, pattern := range sqlInjectionPatterns {
		if strings.Contains(lowerInput, strings.ToLower(pattern)) {
			return false
		}
	}
	return true
}

// validateSchoolID 验证学校ID格式
func validateSchoolID(id string) bool {
	// 学校ID应该是数字
	match, _ := regexp.MatchString(`^[0-9]+$`, id)
	return match
}

// UniversityHandler 学校处理器
// 处理学校相关的 HTTP 请求

type UniversityHandler struct {
	// service 学校服务
	service service.UniversityService
}

// NewUniversityHandler 创建学校处理器实例
// 返回 UniversityHandler 实例
func NewUniversityHandler() *UniversityHandler {
	return &UniversityHandler{
		service: service.NewUniversityService(),
	}
}

// Search 搜索学校
// 处理 GET /search 请求
// 查询参数：keyword - 搜索关键词，page - 页码，pageSize - 每页大小
// 返回 JSON 格式的搜索结果，包含分页信息
func (h *UniversityHandler) Search(c *gin.Context) {
	keyword := c.Query("keyword")

	// 验证和清理输入
	if !validateInput(keyword) {
		c.Error(errors.BadRequest("搜索关键词包含非法字符", nil))
		return
	}
	keyword = sanitizeInput(keyword)

	// 获取分页参数，默认值为第1页，每页10条
	page := 1
	pageSize := 10

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := c.Query("pageSize"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	results, total := h.service.SearchWithPagination(keyword, page, pageSize)

	c.JSON(http.StatusOK, gin.H{
		"results":  results,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GetUniversity 获取学校详情
// 处理 GET /university/:id 请求
// 路径参数：id - 学校ID或名称
// 返回学校详情页面或404页面
func (h *UniversityHandler) GetUniversity(c *gin.Context) {
	id := c.Param("id")

	// 验证和清理输入
	if !validateInput(id) {
		c.HTML(http.StatusOK, "not_found.html", nil)
		return
	}
	id = sanitizeInput(id)

	uni, err := h.service.GetByIDOrName(id)
	if err != nil || uni == nil {
		c.HTML(http.StatusOK, "not_found.html", nil)
		return
	}

	// 处理标签
	tags := []string{}
	if uni.Tags != "" {
		// 简单分割标签，实际项目中可能需要更复杂的处理
		for _, tag := range split(uni.Tags, ", ") {
			tags = append(tags, tag)
		}
	}

	// 处理排名
	rankings := []map[string]string{
		{"name": "软科综合", "value": uni.SoftScienceRank},
		{"name": "校友会综合", "value": uni.AlumniRank},
		{"name": "QS世界", "value": uni.QSRank},
		{"name": "US世界", "value": uni.USNewsRank},
		{"name": "泰晤士 (大陆)", "value": uni.TimesRank},
		{"name": "人气值排名", "value": uni.PopularityRank},
	}

	c.HTML(http.StatusOK, "university.html", gin.H{
		"university": uni,
		"tags":       tags,
		"rankings":   rankings,
	})
}

// EditUniversity 编辑学校页面
// 处理 GET /edit/:id 请求
// 路径参数：id - 学校ID或名称
// 返回学校编辑页面或404页面
func (h *UniversityHandler) EditUniversity(c *gin.Context) {
	id := c.Param("id")

	// 验证和清理输入
	if !validateInput(id) {
		c.HTML(http.StatusOK, "not_found.html", nil)
		return
	}
	id = sanitizeInput(id)

	uni, err := h.service.GetByIDOrName(id)
	if err != nil || uni == nil {
		c.HTML(http.StatusOK, "not_found.html", nil)
		return
	}

	c.HTML(http.StatusOK, "university_edit.html", gin.H{
		"university": uni,
	})
}

// SaveUniversity 保存学校信息
// 处理 POST /save 请求
// 请求体：学校信息 JSON
// 返回保存结果
func (h *UniversityHandler) SaveUniversity(c *gin.Context) {
	var uni models.University
	if err := c.ShouldBindJSON(&uni); err != nil {
		c.Error(errors.BadRequest("请求数据错误", err))
		return
	}

	if uni.SchoolID == "" {
		c.Error(errors.BadRequest("学校ID不能为空", nil))
		return
	}

	// 验证和清理输入
	if !validateSchoolID(uni.SchoolID) {
		c.Error(errors.BadRequest("学校ID格式错误", nil))
		return
	}
	if !validateInput(uni.SchoolName) {
		c.Error(errors.BadRequest("学校名称包含非法字符", nil))
		return
	}

	// 清理输入
	uni.SchoolName = sanitizeInput(uni.SchoolName)
	uni.Address = sanitizeInput(uni.Address)
	uni.Category = sanitizeInput(uni.Category)
	uni.Nature = sanitizeInput(uni.Nature)
	uni.Affiliation = sanitizeInput(uni.Affiliation)
	uni.Tags = sanitizeInput(uni.Tags)
	uni.FoundedYear = sanitizeInput(uni.FoundedYear)
	uni.Area = sanitizeInput(uni.Area)
	uni.BasicInfo = sanitizeInput(uni.BasicInfo)

	err := h.service.Update(uni.SchoolID, &uni)
	if err != nil {
		c.Error(errors.InternalServerError("保存失败", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "保存成功"})
}

// AddSchoolPage 添加学校页面
// 处理 GET /add_school 请求
// 返回添加学校页面，包含可用的学校ID
func (h *UniversityHandler) AddSchoolPage(c *gin.Context) {
	availableID := h.service.GetAvailableSchoolID()
	c.HTML(http.StatusOK, "add_school.html", gin.H{
		"available_id": availableID,
	})
}

// UploadLogo 上传logo
// 处理 POST /upload_logo 请求
// 表单数据：logo - 图片文件，available_id - 学校ID
// 返回上传结果
func (h *UniversityHandler) UploadLogo(c *gin.Context) {
	file, err := c.FormFile("logo")
	if err != nil {
		c.Error(errors.BadRequest("请选择文件", err))
		return
	}

	availableID := c.PostForm("available_id")
	if availableID == "" {
		c.Error(errors.BadRequest("无效的学校ID", nil))
		return
	}

	// 确保logo目录存在
	os.MkdirAll("static/logo", 0755)

	// 保存文件
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", availableID, ext)
	filepath := filepath.Join("static/logo", filename)

	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.Error(errors.InternalServerError("上传失败", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"file_path": filepath,
		"logo_id":   availableID,
	})
}

// AddSchool 添加学校
// 处理 POST /add_school 请求
// 请求体：学校信息 JSON
// 返回添加结果
func (h *UniversityHandler) AddSchool(c *gin.Context) {
	var uni models.University
	if err := c.ShouldBindJSON(&uni); err != nil {
		c.Error(errors.BadRequest("请求数据错误", err))
		return
	}

	if uni.SchoolID == "" || uni.SchoolName == "" {
		c.Error(errors.BadRequest("学校ID和名称不能为空", nil))
		return
	}

	// 验证和清理输入
	if !validateSchoolID(uni.SchoolID) {
		c.Error(errors.BadRequest("学校ID格式错误", nil))
		return
	}
	if !validateInput(uni.SchoolName) {
		c.Error(errors.BadRequest("学校名称包含非法字符", nil))
		return
	}

	// 清理输入
	uni.SchoolName = sanitizeInput(uni.SchoolName)
	uni.Address = sanitizeInput(uni.Address)
	uni.Category = sanitizeInput(uni.Category)
	uni.Nature = sanitizeInput(uni.Nature)
	uni.Affiliation = sanitizeInput(uni.Affiliation)
	uni.Tags = sanitizeInput(uni.Tags)
	uni.FoundedYear = sanitizeInput(uni.FoundedYear)
	uni.Area = sanitizeInput(uni.Area)
	uni.BasicInfo = sanitizeInput(uni.BasicInfo)

	if h.service.SchoolExists(uni.SchoolID, uni.SchoolName) {
		c.Error(errors.BadRequest("学校ID或名称已存在", nil))
		return
	}

	err := h.service.Create(&uni)
	if err != nil {
		c.Error(errors.InternalServerError("添加失败", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "添加成功"})
}

// DeleteSchool 删除学校
// 处理 POST /delete_school 请求
// 请求体：{"school_id": "学校ID"}
// 返回删除结果
func (h *UniversityHandler) DeleteSchool(c *gin.Context) {
	var req struct {
		SchoolID string `json:"school_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest("请求数据错误", err))
		return
	}

	// 验证和清理输入
	if !validateSchoolID(req.SchoolID) {
		c.Error(errors.BadRequest("学校ID格式错误", nil))
		return
	}

	err := h.service.Delete(req.SchoolID)
	if err != nil {
		c.Error(errors.InternalServerError("删除失败", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "删除成功"})
}

// GetScrollingData 获取滚动数据
// 处理 GET /api/scrolling_data 请求
// 返回滚动数据和位置
func (h *UniversityHandler) GetScrollingData(c *gin.Context) {
	data, err := h.service.GetScrollingData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"data": []models.ScrollingUniversity{}, "position": 0})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data, "position": 0})
}

// UpdateScrollingPosition 更新滚动位置
// 处理 POST /api/scrolling_position 请求
// 请求体：{"position": 滚动位置}
// 返回更新结果
func (h *UniversityHandler) UpdateScrollingPosition(c *gin.Context) {
	var req struct {
		Position int `json:"position"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.BadRequest("请求数据错误", err))
		return
	}

	// 这里可以保存滚动位置到数据库或缓存
	// 目前只是返回成功
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// split 分割字符串
// s: 要分割的字符串
// sep: 分隔符
// 返回分割后的字符串数组
func split(s, sep string) []string {
	var result []string
	start := 0
	sepLen := len(sep)
	for i := 0; i <= len(s)-sepLen; i++ {
		if s[i:i+sepLen] == sep {
			result = append(result, s[start:i])
			start = i + sepLen
			i += sepLen - 1
		}
	}
	if start < len(s) {
		result = append(result, s[start:])
	}
	return result
}

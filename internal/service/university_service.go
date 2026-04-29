package service

import (
	"fmt"
	"school-go/internal/models"
	"school-go/internal/repository"
	"sync"
)

// UniversityService 学校服务接口
// 定义了学校相关的业务逻辑操作
// 作为 repository 层和 handler 层之间的桥梁

type UniversityService interface {
	// GetAll 获取所有学校数据
	// 返回学校列表和可能的错误
	GetAll() ([]models.University, error)

	// GetByID 根据学校ID获取学校信息
	// id: 学校ID
	// 返回学校信息和可能的错误
	GetByID(id string) (*models.University, error)

	// GetByIDOrName 根据ID或名称获取学校信息
	// idOrName: 学校ID或名称
	// 返回学校信息和可能的错误
	GetByIDOrName(idOrName string) (*models.University, error)

	// Create 创建新学校
	// uni: 学校信息
	// 返回可能的错误
	Create(uni *models.University) error

	// Update 更新学校信息
	// id: 学校ID
	// uni: 学校信息
	// 返回可能的错误
	Update(id string, uni *models.University) error

	// Delete 删除学校
	// id: 学校ID
	// 返回可能的错误
	Delete(id string) error

	// GetScrollingData 获取滚动数据
	// 返回滚动数据列表和可能的错误
	GetScrollingData() ([]models.ScrollingUniversity, error)

	// Search 搜索学校
	// keyword: 搜索关键词
	// 返回搜索结果
	Search(keyword string) []models.SearchResult

	// SearchWithPagination 分页搜索学校
	// keyword: 搜索关键词
	// page: 页码
	// pageSize: 每页大小
	// 返回搜索结果和总数
	SearchWithPagination(keyword string, page, pageSize int) ([]models.SearchResult, int)

	// GetAvailableSchoolID 获取可用的学校ID
	// 返回可用的学校ID字符串
	GetAvailableSchoolID() string

	// SchoolExists 检查学校是否存在
	// id: 学校ID
	// name: 学校名称
	// 返回学校是否存在
	SchoolExists(id, name string) bool
}

// universityService 学校服务实现
// 封装了学校相关的业务逻辑

type universityService struct {
	// repo 学校数据仓库
	repo repository.UniversityRepository
	// 搜索结果缓存
	searchCache map[string][]models.SearchResult
	// 缓存锁
	cacheMutex sync.RWMutex
	// 初始化标志
	initialized bool
	// 初始化锁
	initMutex sync.Mutex
}

// NewUniversityService 创建学校服务实例
// 返回 UniversityService 接口实现
func NewUniversityService() UniversityService {
	return &universityService{
		repo: repository.NewUniversityRepository(),
		searchCache: make(map[string][]models.SearchResult),
		initialized: false,
	}
}

// initialize 初始化服务
// 实现延迟加载机制，首次请求时加载数据
func (s *universityService) initialize() {
	s.initMutex.Lock()
	defer s.initMutex.Unlock()
	
	if !s.initialized {
		// 首次加载数据，预热缓存
		_, _ = s.repo.GetAll()
		s.initialized = true
	}
}

// GetAll 获取所有学校数据
func (s *universityService) GetAll() ([]models.University, error) {
	s.initialize()
	return s.repo.GetAll()
}

// GetByID 根据学校ID获取学校信息
func (s *universityService) GetByID(id string) (*models.University, error) {
	return s.repo.GetByID(id)
}

// GetByIDOrName 根据ID或名称获取学校信息
// 并发尝试通过ID和名称获取学校信息
func (s *universityService) GetByIDOrName(idOrName string) (*models.University, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var result *models.University
	var err error

	// 尝试通过ID获取
	wg.Add(1)
	go func() {
		defer wg.Done()
		uni, e := s.repo.GetByID(idOrName)
		if e == nil && uni != nil {
			mu.Lock()
			result = uni
			mu.Unlock()
		}
	}()

	// 尝试通过名称获取
	wg.Add(1)
	go func() {
		defer wg.Done()
		// 先检查结果是否已经找到
		mu.Lock()
		if result != nil {
			mu.Unlock()
			return
		}
		mu.Unlock()

		// 获取所有学校并查找名称匹配的
		universities, e := s.repo.GetAll()
		if e != nil {
			mu.Lock()
			err = e
			mu.Unlock()
			return
		}

		for _, u := range universities {
			if u.SchoolName == idOrName {
				mu.Lock()
				if result == nil {
					result = &u
				}
				mu.Unlock()
				return
			}
		}
	}()

	wg.Wait()

	return result, err
}

// Create 创建新学校
func (s *universityService) Create(uni *models.University) error {
	err := s.repo.Create(uni)
	if err == nil {
		// 记录操作日志
		logService := NewOperationLogService()
		logService.LogOperation("create", "universities", uni.SchoolID, 0, "system", "", "", nil, uni)
	}
	return err
}

// Update 更新学校信息
// 确保更新时使用正确的学校ID
func (s *universityService) Update(id string, uni *models.University) error {
	// 获取旧数据
	oldUni, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	
	uni.SchoolID = id
	err = s.repo.Update(uni)
	if err == nil {
		// 记录操作日志
		logService := NewOperationLogService()
		logService.LogOperation("update", "universities", id, 0, "system", "", "", oldUni, uni)
	}
	return err
}

// Delete 删除学校
func (s *universityService) Delete(id string) error {
	// 获取旧数据
	oldUni, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	
	err = s.repo.Delete(id)
	if err == nil {
		// 记录操作日志
		logService := NewOperationLogService()
		logService.LogOperation("delete", "universities", id, 0, "system", "", "", oldUni, nil)
	}
	return err
}

// GetScrollingData 获取滚动数据
func (s *universityService) GetScrollingData() ([]models.ScrollingUniversity, error) {
	return s.repo.GetScrollingData()
}

// Search 搜索学校
// 调用仓库层的搜索方法，处理错误情况，并添加缓存功能
func (s *universityService) Search(keyword string) []models.SearchResult {
	s.initialize()
	
	// 检查缓存
	s.cacheMutex.RLock()
	if results, ok := s.searchCache[keyword]; ok {
		s.cacheMutex.RUnlock()
		return results
	}
	s.cacheMutex.RUnlock()
	
	// 缓存未命中，执行搜索
	results, err := s.repo.Search(keyword)
	if err != nil {
		// 搜索失败时返回空结果
		return []models.SearchResult{}
	}
	
	// 更新缓存
	s.cacheMutex.Lock()
	s.searchCache[keyword] = results
	s.cacheMutex.Unlock()
	
	return results
}

// SearchWithPagination 分页搜索学校
// 调用仓库层的分页搜索方法，处理错误情况，并添加缓存功能
func (s *universityService) SearchWithPagination(keyword string, page, pageSize int) ([]models.SearchResult, int) {
	s.initialize()
	
	// 生成缓存键
	cacheKey := fmt.Sprintf("%s_%d_%d", keyword, page, pageSize)
	
	// 检查缓存
	s.cacheMutex.RLock()
	if results, ok := s.searchCache[cacheKey]; ok {
		s.cacheMutex.RUnlock()
		// 对于分页搜索，我们只缓存结果，不缓存总数
		// 总数需要实时查询，以确保准确性
		_, total, _ := s.repo.SearchWithPagination(keyword, page, pageSize)
		return results, total
	}
	s.cacheMutex.RUnlock()
	
	// 缓存未命中，执行搜索
	results, total, err := s.repo.SearchWithPagination(keyword, page, pageSize)
	if err != nil {
		// 搜索失败时返回空结果
		return []models.SearchResult{}, 0
	}
	
	// 更新缓存
	s.cacheMutex.Lock()
	s.searchCache[cacheKey] = results
	s.cacheMutex.Unlock()
	
	return results, total
}

// GetAvailableSchoolID 获取可用的学校ID
// 查找最小的未使用ID，如果所有ID都已使用则返回最大ID+1
func (s *universityService) GetAvailableSchoolID() string {
	universities, err := s.repo.GetAll()
	if err != nil {
		// 发生错误时返回默认值"1"
		return "1"
	}

	usedIDs := make(map[int]bool)
	maxID := 0

	// 收集已使用的ID并找出最大ID
	for _, uni := range universities {
		id := parseInt(uni.SchoolID)
		if id > 0 {
			usedIDs[id] = true
			if id > maxID {
				maxID = id
			}
		}
	}

	// 查找最小的未使用ID
	for i := 1; i <= maxID; i++ {
		if !usedIDs[i] {
			return toString(i)
		}
	}

	// 所有ID都已使用，返回最大ID+1
	return toString(maxID + 1)
}

// SchoolExists 检查学校是否存在
// 检查ID或名称是否已存在
func (s *universityService) SchoolExists(id, name string) bool {
	universities, err := s.repo.GetAll()
	if err != nil {
		// 发生错误时返回false
		return false
	}

	for _, uni := range universities {
		if uni.SchoolID == id || uni.SchoolName == name {
			return true
		}
	}

	return false
}

// parseInt 将字符串转换为整数
// s: 字符串
// 返回转换后的整数，转换失败返回0
func parseInt(s string) int {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	if err != nil {
		return 0
	}
	return n
}

// toString 将整数转换为字符串
// n: 整数
// 返回转换后的字符串
func toString(n int) string {
	return fmt.Sprintf("%d", n)
}

// 搜索功能
const searchInput = document.getElementById('search-input');
const searchBtn = document.getElementById('search-btn');
const suggestions = document.getElementById('suggestions');
let selectedId = null; // 当前选中的学校ID

// 搜索建议功能
async function showSearchSuggestions() {
    const keyword = searchInput.value.trim();
    console.log('触发搜索建议:', keyword);
    selectedId = null; // 输入变化时清除选中的ID
    
    if (keyword.length < 1) {
        console.log('搜索内容为空，隐藏建议');
        suggestions.innerHTML = '';
        suggestions.style.display = 'none';
        return;
    }
    
    try {
        console.log('发送搜索请求:', `/search?keyword=${encodeURIComponent(keyword)}`);
        const response = await fetch(`/search?keyword=${encodeURIComponent(keyword)}`);
        if (!response.ok) {
            throw new Error(`网络请求失败: ${response.status}`);
        }
        const data = await response.json();
        console.log('搜索响应数据:', data);
        
        // 提取搜索结果数组
        const results = data.results || [];
        
        if (results.length > 0) {
            console.log('显示搜索建议，共', results.length, '条');
            // 显示建议列表
            suggestions.innerHTML = results.map(item => {
                let matchLabel = '';
                // 检查匹配类型
                if (item.match_type === 'exact_id') {
                    matchLabel = '<span class="match-label exact">精确匹配</span>';
                } else if (item.match_type === 'partial_id') {
                    matchLabel = '<span class="match-label partial">ID匹配</span>';
                }
                return `<div class="suggestion-item" data-id="${item.id}" data-name="${item.name}"><span class="suggestion-id">${item.id}</span> ${item.name} ${matchLabel}</div>`;
            }).join('');
            suggestions.style.display = 'block';
            console.log('建议容器显示状态:', suggestions.style.display);
        } else {
            console.log('无搜索结果，隐藏建议');
            suggestions.innerHTML = '';
            suggestions.style.display = 'none';
        }
    } catch (error) {
        console.error('搜索建议错误:', error);
        // 即使出错也不隐藏建议容器，以便用户看到错误状态
        suggestions.innerHTML = `<div class="suggestion-item">搜索失败，请重试</div>`;
        suggestions.style.display = 'block';
    }
}

// 输入事件触发搜索建议
searchInput.addEventListener('input', showSearchSuggestions);

// 获得焦点时触发搜索建议
searchInput.addEventListener('focus', showSearchSuggestions);

// 点击建议项 - 直接跳转详情页
suggestions.addEventListener('click', function(e) {
    const item = e.target.closest('.suggestion-item');
    if (item) {
        const id = item.getAttribute('data-id');
        suggestions.style.display = 'none';
        // 跳转到详情页
        window.location.href = `/university/${id}`;
    }
});

// 执行搜索跳转
function doSearch() {
    const keyword = searchInput.value.trim();
    if (keyword) {
        // 如果有选中的ID，使用ID跳转；否则使用输入的关键字
        const searchValue = selectedId || keyword;
        window.location.href = `/university/${searchValue}`;
    }
}

// 点击搜索按钮
searchBtn.addEventListener('click', function() {
    doSearch();
});

// 按回车键搜索
searchInput.addEventListener('keypress', function(e) {
    if (e.key === 'Enter') {
        doSearch();
    }
});

// 点击页面其他地方关闭建议
document.addEventListener('click', function(e) {
    if (!e.target.closest('.search-container')) {
        suggestions.style.display = 'none';
    }
});

// 登录相关功能
let isLoggedIn = false;
let currentUsername = '';

// 页面加载时检查登录状态
async function checkLoginStatus() {
    try {
        const response = await fetch('/check_login');
        const data = await response.json();
        if (data.logged_in) {
            isLoggedIn = true;
            currentUsername = data.username;
            updateLoginUI();
        }
    } catch (error) {
        console.error('检查登录状态错误:', error);
    }
}

// 更新登录按钮UI
function updateLoginUI() {
    const loginContainer = document.getElementById('loginContainer');
    if (isLoggedIn) {
        loginContainer.innerHTML = `
            <div class="user-dropdown">
                <span class="username" onclick="window.location.href='/account'">${currentUsername}</span>
                <div class="dropdown-menu">
                    <div class="dropdown-item" onclick="window.location.href='/add_school'">添加院校</div>
                    <div class="dropdown-item" onclick="doLogout()">退出登录</div>
                </div>
            </div>
        `;
    } else {
        loginContainer.innerHTML = `
            <button class="login-btn" id="loginBtn" onclick="showLoginModal()">登录</button>
        `;
    }
}

// 显示登录弹窗
function showLoginModal() {
    document.getElementById('loginModal').classList.add('show');
    document.getElementById('errorMessage').style.display = 'none';
    document.getElementById('username').value = '';
    document.getElementById('password').value = '';
    document.getElementById('username').focus();
}

// 隐藏登录弹窗
function hideLoginModal() {
    document.getElementById('loginModal').classList.remove('show');
}

// 执行登录
async function doLogin() {
    const username = document.getElementById('username').value.trim();
    const password = document.getElementById('password').value;
    const errorMessage = document.getElementById('errorMessage');
    
    if (!username || !password) {
        errorMessage.textContent = '请输入用户名和密码';
        errorMessage.style.display = 'block';
        return;
    }
    
    try {
        const response = await fetch('/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, password })
        });
        
        const data = await response.json();
        
        if (data.success) {
            isLoggedIn = true;
            currentUsername = username;
            updateLoginUI();
            hideLoginModal();
        } else {
            errorMessage.textContent = data.message;
            errorMessage.style.display = 'block';
        }
    } catch (error) {
        console.error('登录错误:', error);
        errorMessage.textContent = '登录失败，请重试';
        errorMessage.style.display = 'block';
    }
}

// 执行退出
async function doLogout() {
    try {
        await fetch('/logout', { method: 'POST' });
        isLoggedIn = false;
        currentUsername = '';
        updateLoginUI();
    } catch (error) {
        console.error('退出错误:', error);
    }
}

// 点击弹窗外部关闭
document.getElementById('loginModal').addEventListener('click', function(e) {
    if (e.target === this) {
        hideLoginModal();
    }
});

// 按回车键登录
document.getElementById('password').addEventListener('keypress', function(e) {
    if (e.key === 'Enter') {
        doLogin();
    }
});

// 滚动数据功能
let scrollingData = [];
let currentPosition = 0;
let scrollInterval = null;
let isPaused = false;
const scrollSpeed = 0.5; // 滚动速度（像素/帧）
const itemHeight = 45; // 每项高度

// 加载滚动数据
async function loadScrollingData() {
    try {
        const response = await fetch('/api/scrolling_data');
        const result = await response.json();
        scrollingData = result.data;
        currentPosition = result.position;
        renderScrollingData();
        startScrolling();
    } catch (error) {
        console.error('加载滚动数据错误:', error);
    }
}

// 渲染滚动数据
function renderScrollingData() {
    const container = document.getElementById('scrollingContent');
    if (!scrollingData || !scrollingData.length) return;
    
    // 复制数据以实现无缝滚动
    const displayData = [...scrollingData, ...scrollingData];
    
    container.innerHTML = displayData.map(item => `
        <div class="scrolling-item" onclick="window.location.href='/university/${item['学校ID']}'">
            <div class="scrolling-cell id">${item['学校ID']}</div>
            <div class="scrolling-cell name">${item['学校名称']}</div>
            <div class="scrolling-cell">${item['地址']}</div>
            <div class="scrolling-cell">${item['类别']}</div>
            <div class="scrolling-cell">${item['性质']}</div>
        </div>
    `).join('');
    
    // 设置初始位置
    updateScrollPosition();
}

// 更新滚动位置
function updateScrollPosition() {
    const container = document.getElementById('scrollingContent');
    const offset = currentPosition * itemHeight;
    container.style.transform = `translateY(-${offset}px)`;
}

// 开始滚动
function startScrolling() {
    if (scrollInterval) return;
    
    scrollInterval = setInterval(() => {
        if (!isPaused && scrollingData.length > 0) {
            currentPosition += scrollSpeed / itemHeight;
            
            // 当滚动到一半时（复制数据的位置），重置到开头
            if (currentPosition >= scrollingData.length) {
                currentPosition = 0;
            }
            
            updateScrollPosition();
        }
    }, 16); // 约60fps
}

// 停止滚动
function stopScrolling() {
    if (scrollInterval) {
        clearInterval(scrollInterval);
        scrollInterval = null;
    }
}

// 保存滚动位置到服务器
async function saveScrollPosition() {
    try {
        await fetch('/api/scrolling_position', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ position: Math.floor(currentPosition) })
        });
    } catch (error) {
        console.error('保存滚动位置错误:', error);
    }
}

// 鼠标悬停事件
const scrollingContainer = document.getElementById('scrollingContainer');
scrollingContainer.addEventListener('mouseenter', () => {
    isPaused = true;
});

scrollingContainer.addEventListener('mouseleave', () => {
    isPaused = false;
});

// 页面离开前保存位置
window.addEventListener('beforeunload', () => {
    saveScrollPosition();
});

// 页面可见性变化时保存位置
document.addEventListener('visibilitychange', () => {
    if (document.hidden) {
        saveScrollPosition();
    }
});

// 页面加载完成后执行
window.addEventListener('DOMContentLoaded', function() {
    // 检查登录状态
    checkLoginStatus();
    
    // 清空搜索框
    searchInput.value = '';
    
    // 加载滚动数据
    loadScrollingData();
});

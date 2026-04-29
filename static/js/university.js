let isLoggedIn = false;
let currentUsername = '';
let currentUserRole = '';
// schoolId 由HTML页面动态设置
let loginAction = '';

// 页面加载时检查登录状态
async function checkLoginStatus() {
    try {
        const response = await fetch('/check_login');
        const data = await response.json();
        if (data.logged_in) {
            isLoggedIn = true;
            currentUsername = data.username;
            currentUserRole = data.role || '';
            updateLoginUI();
        }
        // 无论是否登录，都更新操作按钮
        updateActionButtons();
    } catch (error) {
        console.error('检查登录状态错误:', error);
        // 即使检查登录状态失败，也更新操作按钮
        updateActionButtons();
    }
}

// 根据用户角色更新操作按钮
function updateActionButtons() {
    // 获取编辑和删除按钮
    const editBtn = document.getElementById('editBtn');
    const deleteBtn = document.getElementById('deleteBtn');
    
    // 普通用户和管理员都可以看到编辑按钮
    if (editBtn) {
        editBtn.style.display = 'inline-block';
    }
    
    // 未登录状态和普通用户都可以看到删除按钮
    if (deleteBtn) {
        deleteBtn.style.display = 'inline-block';
    }
}

// 更新登录按钮UI
function updateLoginUI() {
    const loginContainer = document.getElementById('loginContainer');
    if (isLoggedIn) {
        loginContainer.innerHTML = `
            <span class="username" onclick="window.location.href='/account'" style="cursor: pointer;">${currentUsername}</span>
        `;
    } else {
        loginContainer.innerHTML = '';
    }
}

// 处理编辑按钮点击
function handleEdit() {
    if (isLoggedIn) {
        window.location.href = `/edit/${schoolId}`;
    } else {
        loginAction = 'edit';
        showLoginModal();
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
            
            if (loginAction === 'edit') {
                window.location.href = `/edit/${schoolId}`;
            } else if (loginAction === 'delete') {
                document.getElementById('deleteModal').classList.add('show');
            }
            loginAction = '';
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

// 显示删除确认弹窗
function showDeleteModal() {
    if (isLoggedIn) {
        if (currentUserRole === 'admin') {
            document.getElementById('deleteModal').classList.add('show');
        } else {
            showErrorModal('权限不足，无法执行此操作');
        }
    } else {
        loginAction = 'delete';
        showLoginModal();
    }
}

// 隐藏删除确认弹窗
function hideDeleteModal() {
    document.getElementById('deleteModal').classList.remove('show');
}

// 显示错误弹窗
function showErrorModal(message) {
    // 创建错误弹窗元素
    const errorModal = document.createElement('div');
    errorModal.id = 'errorModal';
    errorModal.className = 'modal-overlay show';
    errorModal.innerHTML = `
        <div class="login-modal" style="max-width: 400px;">
            <div class="modal-title">操作失败</div>
            <div class="delete-message">${message}</div>
            <div class="modal-actions">
                <button class="modal-btn modal-btn-primary" onclick="hideErrorModal()">确定</button>
            </div>
        </div>
    `;
    document.body.appendChild(errorModal);
    
    // 点击弹窗外部关闭
    errorModal.addEventListener('click', function(e) {
        if (e.target === this) {
            hideErrorModal();
        }
    });
}

// 隐藏错误弹窗
function hideErrorModal() {
    const errorModal = document.getElementById('errorModal');
    if (errorModal) {
        errorModal.remove();
    }
}

// 确认删除
async function confirmDelete() {
    try {
        const response = await fetch('/delete_school', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest'
            },
            body: JSON.stringify({ school_id: schoolId })
        });
        
        // 检查响应状态
        if (!response.ok) {
            // 如果是 403 错误，显示权限不足的错误
            if (response.status === 403) {
                showErrorModal('权限不足，无法执行此操作');
            } else {
                showErrorModal('删除失败，请稍后重试');
            }
            return;
        }
        
        const data = await response.json();
        
        if (data.success) {
            hideDeleteModal();
            // 删除成功后跳转到首页
            window.location.href = '/';
        } else {
            showErrorModal('删除失败: ' + data.message);
        }
    } catch (error) {
        console.error('删除错误:', error);
        showErrorModal('删除失败，请稍后重试');
    }
}

// 页面加载完成后执行
window.addEventListener('DOMContentLoaded', function() {
    // 检查登录状态
    checkLoginStatus();
    
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
    
    // 点击删除弹窗外部关闭
    document.getElementById('deleteModal').addEventListener('click', function(e) {
        if (e.target === this) {
            hideDeleteModal();
        }
    });
});

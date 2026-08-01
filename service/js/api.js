// API 服务配置
const API_BASE_URL = 'http://localhost:8080/api';

// API 请求封装
class ApiService {
    // 获取存储的 Token
    static getToken() {
        return localStorage.getItem('token');
    }

    // 存储 Token
    static setToken(token) {
        localStorage.setItem('token', token);
    }

    // 清除 Token
    static clearToken() {
        localStorage.removeItem('token');
    }

    // 通用请求方法
    static async request(endpoint, options = {}) {
        const url = `${API_BASE_URL}${endpoint}`;
        const token = this.getToken();
        
        const headers = {
            'Content-Type': 'application/json',
            ...options.headers
        };
        
        if (token) {
            headers['Authorization'] = token;
        }

        try {
            const response = await fetch(url, {
                ...options,
                headers
            });

            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.message || '请求失败');
            }

            return data;
        } catch (error) {
            if (error.message === 'Failed to fetch') {
                // 服务器不可达，返回 null 表示需要降级到脱机模式
                return null;
            }
            throw error;
        }
    }

    // 检查服务器是否可达
    static async isServerAvailable() {
        try {
            const response = await fetch(`${API_BASE_URL}/user/info`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json'
                }
            });
            // 即使返回 401（未授权），也表示服务器在线
            return response.status === 200 || response.status === 401;
        } catch {
            return false;
        }
    }

    // 用户注册
    static async register(username, emailOrPhone, password) {
        return this.request('/user/register', {
            method: 'POST',
            body: JSON.stringify({ username, emailOrPhone, password })
        });
    }

    // 用户登录
    static async login(loginId, password) {
        return this.request('/user/login', {
            method: 'POST',
            body: JSON.stringify({ loginId, password })
        });
    }

    // 获取用户信息
    static async getUserInfo() {
        return this.request('/user/info', {
            method: 'GET'
        });
    }

    // 更新用户信息
    static async updateUser(userData) {
        return this.request('/user/info', {
            method: 'PUT',
            body: JSON.stringify(userData)
        });
    }

    // 修改密码
    static async changePassword(oldPassword, newPassword) {
        return this.request('/user/password', {
            method: 'PUT',
            body: JSON.stringify({ oldPassword, newPassword })
        });
    }

    // 用户登出
    static async logout() {
        return this.request('/user/logout', {
            method: 'POST'
        });
    }

    // ==================== 数据存取 API ====================

    // 保存单条数据（存在则更新）
    static async saveData(category, content, dataKey = '') {
        return this.request('/data/save', {
            method: 'POST',
            body: JSON.stringify({ category, dataKey, content })
        });
    }

    // 批量保存数据
    static async batchSaveData(items) {
        return this.request('/data/batch-save', {
            method: 'POST',
            body: JSON.stringify({ items })
        });
    }

    // 获取数据列表（支持按分类和键筛选）
    static async getDataList(category = '', dataKey = '') {
        const params = new URLSearchParams();
        if (category) params.set('category', category);
        if (dataKey) params.set('dataKey', dataKey);
        const query = params.toString() ? `?${params.toString()}` : '';
        return this.request(`/data/list${query}`, {
            method: 'GET'
        });
    }

    // 获取所有数据分类
    static async getDataCategories() {
        return this.request('/data/categories', {
            method: 'GET'
        });
    }

    // 根据 ID 获取单条数据
    static async getDataById(id) {
        return this.request(`/data/${id}`, {
            method: 'GET'
        });
    }

    // 根据 ID 删除数据
    static async deleteData(id) {
        return this.request(`/data/${id}`, {
            method: 'DELETE'
        });
    }

    // 按分类删除数据
    static async deleteDataByCategory(category) {
        return this.request(`/data/category/${category}`, {
            method: 'DELETE'
        });
    }

    // ==================== 数据同步工具 ====================

    // 将本地 localStorage 数据同步到服务器
    static async syncLocalToServer(category, localKey, dataKey = '') {
        const localData = localStorage.getItem(localKey);
        if (!localData) return null;
        return this.saveData(category, localData, dataKey);
    }

    // 从服务器拉取数据到本地 localStorage
    static async pullServerToLocal(category, localKey, dataKey = '') {
        const result = await this.getDataList(category, dataKey);
        if (result && result.data && result.data.length > 0) {
            const latest = result.data[result.data.length - 1];
            localStorage.setItem(localKey, latest.content);
            return latest;
        }
        return null;
    }
}

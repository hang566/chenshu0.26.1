/**
 * 云同步模块 - 自动同步本地数据到服务器
 * 只需引入此脚本即可自动工作，无需额外配置
 * <script src="../js/cloud-sync.js"></script>
 */
const CloudSync = {
    category: '',
    isOnline: false,
    isSyncing: false,
    lastSyncTime: null,
    autoSyncTimer: null,
    pendingKeys: new Set(),
    syncBtn: null,
    syncDot: null,
    syncText: null,

    // 系统级 key，不参与同步
    SYSTEM_KEYS: ['token', 'user', 'theme', 'accounts', '_offlineData', '_offlineDataBackup'],

    // ==================== 初始化 ====================
    async init() {
        this.category = this.detectCategory();
        this.injectStyles();
        this.createSyncButton();

        // 监听 localStorage 变化
        this.monitorStorage();

        // 检查服务器状态
        this.isOnline = await ApiService.isServerAvailable();
        this.updateUI();

        if (this.isOnline && ApiService.getToken()) {
            // 自动拉取服务器数据
            await this.pull();
            // 启动自动同步
            this.startAutoSync();
        }
    },

    detectCategory() {
        const path = location.pathname.replace(/\\/g, '/');
        const match = path.match(/service\/([^\/]+)/);
        if (match) {
            return match[1].toLowerCase().replace(/\s+/g, '_');
        }
        return 'general';
    },

    // ==================== UI ====================
    injectStyles() {
        if (document.getElementById('cloud-sync-styles')) return;
        const style = document.createElement('style');
        style.id = 'cloud-sync-styles';
        style.textContent = `
            .cs-fab {
                position: fixed; bottom: 20px; right: 20px; z-index: 9998;
                display: flex; align-items: center; gap: 8px;
                background: var(--card-bg, #fff);
                border: 1px solid var(--border, #e0e0e0);
                border-radius: 24px; padding: 8px 16px 8px 12px;
                box-shadow: 0 2px 12px rgba(0,0,0,0.12);
                cursor: pointer; transition: all 0.2s;
                font-size: 13px; user-select: none;
            }
            .cs-fab:hover { box-shadow: 0 4px 16px rgba(0,0,0,0.18); transform: translateY(-1px); }
            .cs-fab:active { transform: translateY(0); }
            .cs-dot {
                width: 8px; height: 8px; border-radius: 50%;
                flex-shrink: 0; transition: background 0.3s;
            }
            .cs-dot.online { background: #34c759; }
            .cs-dot.offline { background: #ff9500; }
            .cs-dot.syncing { background: #007aff; animation: cs-pulse 1s infinite; }
            .cs-dot.error { background: #ff3b30; }
            .cs-text { color: var(--text, #333); white-space: nowrap; }
            .cs-toast {
                position: fixed; bottom: 64px; right: 20px; z-index: 9999;
                background: rgba(0,0,0,0.8); color: #fff;
                padding: 8px 16px; border-radius: 8px; font-size: 13px;
                animation: cs-fadeIn 0.2s; max-width: 280px; word-break: break-all;
            }
            @keyframes cs-pulse { 0%,100% { opacity: 1; } 50% { opacity: 0.4; } }
            @keyframes cs-fadeIn { from { opacity: 0; transform: translateY(4px); } to { opacity: 1; transform: translateY(0); } }
        `;
        document.head.appendChild(style);
    },

    createSyncButton() {
        this.syncBtn = document.createElement('div');
        this.syncBtn.className = 'cs-fab';
        this.syncBtn.innerHTML = `
            <span class="cs-dot offline"></span>
            <span class="cs-text">检查中...</span>
        `;
        this.syncDot = this.syncBtn.querySelector('.cs-dot');
        this.syncText = this.syncBtn.querySelector('.cs-text');
        this.syncBtn.onclick = () => this.manualSync();
        document.body.appendChild(this.syncBtn);
    },

    updateUI() {
        if (this.isSyncing) {
            this.syncDot.className = 'cs-dot syncing';
            this.syncText.textContent = '同步中...';
        } else if (this.isOnline) {
            this.syncDot.className = 'cs-dot online';
            const time = this.lastSyncTime ? this.formatTime(this.lastSyncTime) : '就绪';
            this.syncText.textContent = `已同步 ${time}`;
        } else {
            this.syncDot.className = 'cs-dot offline';
            this.syncText.textContent = '离线模式';
        }
    },

    formatTime(ts) {
        const diff = Date.now() - ts;
        if (diff < 60000) return '刚刚';
        if (diff < 3600000) return Math.floor(diff / 60000) + '分钟前';
        return new Date(ts).toLocaleTimeString();
    },

    toast(msg) {
        const t = document.createElement('div');
        t.className = 'cs-toast';
        t.textContent = msg;
        document.body.appendChild(t);
        setTimeout(() => t.remove(), 2500);
    },

    // ==================== 数据同步 ====================
    isSystemKey(key) {
        return this.SYSTEM_KEYS.includes(key) || key.startsWith('cs_');
    },

    getSyncableKeys() {
        const keys = [];
        for (let i = 0; i < localStorage.length; i++) {
            const key = localStorage.key(i);
            if (!this.isSystemKey(key)) {
                keys.push(key);
            }
        }
        return keys;
    },

    // 拉取服务器数据到本地
    async pull() {
        if (!this.isOnline || !ApiService.getToken()) return;

        try {
            const result = await ApiService.getDataList(this.category);
            if (!result || !result.data) return;

            let pulled = 0;
            for (const item of result.data) {
                const key = item.dataKey || item.category;
                // 只更新本地没有的数据，避免覆盖本地新修改
                const localValue = localStorage.getItem(key);
                if (!localValue) {
                    localStorage.setItem(key, item.content);
                    pulled++;
                }
            }
            if (pulled > 0) {
                this.toast(`从云端恢复 ${pulled} 项数据`);
            }
            this.lastSyncTime = Date.now();
            this.updateUI();
        } catch (e) {
            console.error('[CloudSync] Pull failed:', e);
        }
    },

    // 推送单个 key 到服务器
    async pushKey(key) {
        if (!this.isOnline || !ApiService.getToken()) return;
        const value = localStorage.getItem(key);
        if (!value) return;

        try {
            await ApiService.saveData(this.category, value, key);
        } catch (e) {
            console.error(`[CloudSync] Push key "${key}" failed:`, e);
        }
    },

    // 推送所有本地数据到服务器
    async pushAll() {
        if (!this.isOnline || !ApiService.getToken()) {
            this.toast('离线模式，无法同步');
            return false;
        }

        this.isSyncing = true;
        this.updateUI();

        const keys = this.getSyncableKeys();
        if (keys.length === 0) {
            this.isSyncing = false;
            this.lastSyncTime = Date.now();
            this.updateUI();
            this.toast('没有需要同步的数据');
            return true;
        }

        try {
            const items = keys.map(key => ({
                category: this.category,
                dataKey: key,
                content: localStorage.getItem(key)
            }));

            const result = await ApiService.batchSaveData(items);
            this.isSyncing = false;
            this.lastSyncTime = Date.now();

            if (result) {
                this.toast(`已同步 ${keys.length} 项数据到云端`);
            }
            this.updateUI();
            return true;
        } catch (e) {
            this.isSyncing = false;
            this.syncDot.className = 'cs-dot error';
            this.syncText.textContent = '同步失败';
            this.toast('同步失败：' + (e.message || '未知错误'));
            setTimeout(() => this.updateUI(), 2000);
            return false;
        }
    },

    // 手动同步
    async manualSync() {
        // 先检查服务器状态
        this.isOnline = await ApiService.isServerAvailable();
        if (!this.isOnline) {
            this.updateUI();
            this.toast('服务器不可达，当前为离线模式');
            return;
        }
        if (!ApiService.getToken()) {
            this.toast('请先登录以使用云同步');
            return;
        }
        await this.pushAll();
    },

    // ==================== 自动同步 ====================
    monitorStorage() {
        // 拦截 setItem
        const originalSetItem = localStorage.setItem.bind(localStorage);
        const self = this;
        localStorage.setItem = function(key, value) {
            originalSetItem(key, value);
            if (!self.isSystemKey(key)) {
                self.pendingKeys.add(key);
            }
        };

        // 拦截 removeItem
        const originalRemoveItem = localStorage.removeItem.bind(localStorage);
        localStorage.removeItem = function(key) {
            originalRemoveItem(key);
            if (!self.isSystemKey(key)) {
                self.pendingKeys.add(key);
            }
        };
    },

    startAutoSync() {
        if (this.autoSyncTimer) clearInterval(this.autoSyncTimer);
        this.autoSyncTimer = setInterval(() => {
            if (this.pendingKeys.size > 0 && this.isOnline && ApiService.getToken()) {
                this.syncPending();
            }
        }, 5000); // 每 5 秒检查一次待同步队列
    },

    async syncPending() {
        if (this.isSyncing || this.pendingKeys.size === 0) return;

        const keys = Array.from(this.pendingKeys);
        this.pendingKeys.clear();

        this.isSyncing = true;
        this.updateUI();

        try {
            const items = keys.map(key => ({
                category: this.category,
                dataKey: key,
                content: localStorage.getItem(key) || ''
            }));
            await ApiService.batchSaveData(items);
            this.lastSyncTime = Date.now();
        } catch (e) {
            // 失败则重新加入队列
            keys.forEach(k => this.pendingKeys.add(k));
            console.error('[CloudSync] Auto sync failed:', e);
        }

        this.isSyncing = false;
        this.updateUI();
    }
};

// 自动初始化（DOM 加载完成后）
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => CloudSync.init());
} else {
    CloudSync.init();
}

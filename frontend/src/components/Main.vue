<template>
  <!-- Same template as before -->
  <div class="app-container">
    <div class="window">
      <!-- Toolbar -->
      <div class="toolbar">
        <button @click="editRules" title="Edit Rules">
          <i class="icon-file-text"></i>
        </button>
        <div class="spacer"></div>
        <button @click="settingsMenu" title="Settings">
          <i class="icon-settings"></i>
        </button>
        <button @click="helpMenu" title="Help">
          <i class="icon-help"></i>
        </button>
      </div>

      <div class="divider"></div>

      <!-- CA Certificate Button -->
      <div class="ca-button-container">
        <button
            v-if="showCaButton"
            @click="installCaCertificate"
            class="ca-button warning"
        >
          {{ CA_STATUS }}
        </button>
      </div>

      <div class="spacer-20"></div>

      <!-- Status Indicators -->
      <div class="status-row">
        <div class="status-label">{{ PROXY_TITLE }}</div>
        <div class="status-value" :class="proxyStatusClass">{{ proxyStatus }}</div>
      </div>

      <div class="status-row">
        <div class="status-label">{{ HTTPS_TITLE }}</div>
        <div class="status-value" :class="httpsStatusClass">{{ httpsStatus }}</div>
      </div>

      <div class="spacer-50"></div>

      <!-- Start/Stop Button -->
      <button
          @click="toggleProxy"
          class="action-button"
          :class="proxyButtonClass"
      >
        {{ proxyButtonText }}
      </button>

      <!-- Settings Menu (shown when settings button is clicked) -->
      <div v-if="showSettingsMenu" class="popup-menu" :style="settingsMenuPosition">
        <div class="menu-item" @click="toggleAutoProxy">
          <span>自动开启系统代理</span>
          <input type="checkbox" v-model="autoProxyEnabled">
        </div>
        <div class="menu-item" @click="toggleMinimizeExit">
          <span>最小化退出</span>
          <input type="checkbox" v-model="minimizeExitEnabled">
        </div>
      </div>

      <!-- Help Menu (shown when help button is clicked) -->
      <div v-if="showHelpMenu" class="popup-menu" :style="helpMenuPosition">
        <div class="menu-item" @click="installCaCertificate">
          <span>安装证书</span>
        </div>
        <div class="menu-item" @click="showAbout">
          <span>关于</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import {ref, computed, onMounted} from 'vue';

// Constants
const APP_WIDTH = 400;
const APP_HEIGHT = 500;
const PROXY_TITLE = "系统代理: ";
const HTTPS_TITLE = "HTTPS: ";
const PROXY_BTN_START = "开始";
const PROXY_BTN_STOP = "停止";
const PROXY_STATUS_RUNNING = "启动";
const PROXY_STATUS_ABNORMAL = "异常";
const PROXY_STATUS_OFF = "关闭";
const CA_STATUS = "证书未安装";

// Reactive state
const showCaButton = ref(false);
const proxyStatus = ref(PROXY_STATUS_OFF);
const httpsStatus = ref(PROXY_STATUS_OFF);
const proxyButtonText = ref(PROXY_BTN_START);
const appRunning = ref(false);
const showSettingsMenu = ref(false);
const showHelpMenu = ref(false);
const autoProxyEnabled = ref(false);
const minimizeExitEnabled = ref(false);
const settingsMenuPosition = ref({top: '10', left: '-10'});
const helpMenuPosition = ref({top: '10', left: '-10'});

// Computed properties
const proxyStatusClass = computed(function () {
  return {
    'warning': proxyStatus.value === PROXY_STATUS_RUNNING ||
        proxyStatus.value.includes(PROXY_STATUS_RUNNING),
    'error': proxyStatus.value === PROXY_STATUS_ABNORMAL
  };
});

const httpsStatusClass = computed(function () {
  return {
    'warning': httpsStatus.value === PROXY_STATUS_RUNNING
  };
});

const proxyButtonClass = computed(function () {
  return {
    'warning': proxyButtonText.value === PROXY_BTN_STOP,
    'medium': proxyButtonText.value === PROXY_BTN_START
  };
});

// Methods
function toggleProxy() {
  if (proxyButtonText.value === PROXY_BTN_START) {
    // Start proxy
    proxyButtonText.value = PROXY_BTN_STOP;
    appRunning.value = true;

    if (autoProxyEnabled.value) {
      // Simulate system proxy on
      proxyStatus.value = `${PROXY_STATUS_RUNNING}(:8080)`; // Assuming port 8080
    }
  } else {
    // Stop proxy
    proxyButtonText.value = PROXY_BTN_START;
    appRunning.value = false;

    if (autoProxyEnabled.value) {
      // Simulate system proxy off
      proxyStatus.value = PROXY_STATUS_OFF;
    }
  }
}

function installCaCertificate() {
  // Simulate CA certificate installation
  console.log("Installing CA certificate...");
  showCaButton.value = false;
}

function editRules() {
  console.log("Opening rules editor...");
}

function settingsMenu(event) {
  showSettingsMenu.value = !showSettingsMenu.value;
  showHelpMenu.value = false;

  if (showSettingsMenu.value) {
    const rect = event.target.getBoundingClientRect();
    settingsMenuPosition.value = {
      top: `${rect.bottom + 5}px`,
      left: `${rect.left}px`
    };
  }
}

function helpMenu(event) {
  showHelpMenu.value = !showHelpMenu.value;
  showSettingsMenu.value = false;

  if (showHelpMenu.value) {
    const rect = event.target.getBoundingClientRect();
    helpMenuPosition.value = {
      top: `${rect.bottom + 5}px`,
      left: `${rect.left}px`
    };
  }
}

function toggleAutoProxy() {
  autoProxyEnabled.value = !autoProxyEnabled.value;
  // In a real app, this would save to config
}

function toggleMinimizeExit() {
  minimizeExitEnabled.value = !minimizeExitEnabled.value;
  // In a real app, this would save to config
}

function showAbout() {
  // Show about dialog
  console.log("Showing about dialog...");
}

// Initialization (simulating initTask)
onMounted(function () {
  // Check if HTTPS is enabled
  httpsStatus.value = PROXY_STATUS_RUNNING;

  // Check if CA certificate is installed
  // In this simulation, we'll randomly show the CA button
  showCaButton.value = Math.random() > 0.5;
});
</script>

<style>
/* Same styles as before */
.app-container {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen,
  Ubuntu, Cantarell, "Open Sans", "Helvetica Neue", sans-serif;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f5f5f5;
}

.window {
  width: 400px;
  height: 500px;
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  padding: 16px;
  position: relative;
}

/* Toolbar styles */
.toolbar {
  display: flex;
  align-items: center;
  padding: 8px 0;
}

.toolbar button {
  background: none;
  border: none;
  cursor: pointer;
  padding: 8px;
  font-size: 16px;
  color: #333;
}

.toolbar button:hover {
  background-color: #f0f0f0;
  border-radius: 4px;
}

.spacer {
  flex-grow: 1;
}

/* Divider */
.divider {
  height: 1px;
  background-color: #808080;
  margin: 8px 0;
}

/* Status indicators */
.status-row {
  display: flex;
  align-items: center;
  margin: 8px 0;
  padding-left: 120px;
}

.status-label {
  margin-right: 8px;
}

.status-value {
  font-weight: bold;
}

.status-value.warning {
  color: #ff9800;
}

.status-value.error {
  color: #f44336;
}

/* CA Button */
.ca-button-container {
  display: flex;
  justify-content: center;
  margin: 10px 0;
}

.ca-button {
  padding: 8px 16px;
  border-radius: 4px;
  border: none;
  cursor: pointer;
}

.ca-button.warning {
  background-color: #ff9800;
  color: white;
}

/* Action Button */
.action-button {
  margin: 50px auto;
  padding: 12px 24px;
  border-radius: 4px;
  border: none;
  cursor: pointer;
  font-size: 16px;
  font-weight: bold;
}

.action-button.warning {
  background-color: #ff9800;
  color: white;
}

.action-button.medium {
  background-color: #2196f3;
  color: white;
}

/* Spacer utilities */
.spacer-20 {
  height: 20px;
}

.spacer-50 {
  height: 50px;
}

/* Popup menus */
.popup-menu {
  position: absolute;
  background-color: white;
  border-radius: 4px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  z-index: 100;
  min-width: 200px;
}

.menu-item {
  padding: 8px 16px;
  cursor: pointer;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.menu-item:hover {
  background-color: #f5f5f5;
}

/* Icon styles (using Unicode as simple replacement) */
.icon-file-text::before {
  content: "📄";
}

.icon-settings::before {
  content: "⚙️";
}

.icon-help::before {
  content: "❓";
}
</style>
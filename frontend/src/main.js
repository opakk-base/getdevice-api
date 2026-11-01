import { Browser } from "/wails/runtime.js";
import * as App from "./bindings/getdevice-api/app.js";

// Track which fields are exposed
let exposedFields = new Set([
    "device_id",
    "device_name",
    "client_key",
    "hostname",
    "os",
    "arch",
    "mac_address",
    "ip_address",
    "timestamp",
]);

// Port change debounce timer
let portDebounceTimer = null;

// ========== Initialization ==========

document.addEventListener("DOMContentLoaded", async () => {
    setupPortInput();
    setupLockScreen();
    setupPasswordSettings();
    document.getElementById("lock-password").focus();
});

async function initState() {
    try {
        // Get server status and set port input
        const status = await App.GetServerStatus();
        updateServerUI(status);
        document.getElementById("port-input").value = status.port;
    } catch (err) {
        console.error("Failed to init state:", err);
    }
}

// ========== Tab Navigation ==========

function switchTab(tabName) {
    // Hide all tab contents
    document.querySelectorAll(".tab-content").forEach((el) => {
        el.classList.remove("active");
    });

    // Deactivate all tabs
    document.querySelectorAll(".tab-bar .tab").forEach((el) => {
        el.classList.remove("active");
    });

    // Show selected tab content
    const tabContent = document.getElementById("tab-" + tabName);
    if (tabContent) {
        tabContent.classList.add("active");
    }

    // Activate selected tab button
    const tabs = document.querySelectorAll(".tab-bar .tab");
    const tabIndex = { device: 0, settings: 1, about: 2 };
    if (tabs[tabIndex[tabName]]) {
        tabs[tabIndex[tabName]].classList.add("active");
    }

    // Load data for the tab if needed
    if (tabName === "settings") {
        loadSettings();
    } else if (tabName === "about") {
        loadAboutInfo();
    }
}

// ========== Server Controls ==========

async function toggleServer() {
    try {
        const status = await App.GetServerStatus();

        if (status.running) {
            await App.StopServer();
        } else {
            await App.StartServer();
        }

        await new Promise((r) => setTimeout(r, 200));

        const newStatus = await App.GetServerStatus();
        updateServerUI(newStatus);
    } catch (err) {
        console.error("Server toggle error:", err);
        try {
            const s = await App.GetServerStatus();
            updateServerUI(s);
        } catch (_) {}
    }
}

function updateServerUI(status) {
    const btn = document.getElementById("server-btn");
    const btnText = document.getElementById("server-btn-text");
    const statusBar = document.getElementById("status-bar");
    const statusText = document.getElementById("status-text");

    if (status.running) {
        btn.className = "server-btn running";
        btnText.textContent = "Hentikan Server";
        statusBar.className = "status-bar";
        statusText.innerHTML = 'Berjalan pada port <span id="port-display">' + status.port + "</span>";
    } else {
        btn.className = "server-btn stopped";
        btnText.textContent = "Jalankan Server";
        statusBar.className = "status-bar stopped";
        statusText.innerHTML = "Server berhenti";
    }
    updateEndpointLinks(status.port);
}

function updateEndpointLinks(port) {
    for (const [id, path, label] of [
        ["getdevice-link", "/getdevice", "Data perangkat"],
        ["health-link", "/health", "Status API"],
    ]) {
        const link = document.getElementById(id);
        const url = `http://localhost:${port}${path}`;
        link.href = url;
        link.textContent = `${label} — ${url}`;
    }
}

// ========== Port Input ==========

function setupPortInput() {
    const input = document.getElementById("port-input");

    input.addEventListener("input", () => {
        input.value = input.value.replace(/[^0-9]/g, "");

        const port = parseInt(input.value, 10);
        if (input.value && (isNaN(port) || port < 1 || port > 65535)) {
            input.classList.add("invalid");
        } else {
            input.classList.remove("invalid");
        }

        clearTimeout(portDebounceTimer);
        if (input.value && !isNaN(port) && port >= 1 && port <= 65535) {
            portDebounceTimer = setTimeout(() => applyPort(input.value), 800);
        }
    });

    input.addEventListener("keydown", (e) => {
        if (e.key === "Enter") {
            clearTimeout(portDebounceTimer);
            const port = parseInt(input.value, 10);
            if (!isNaN(port) && port >= 1 && port <= 65535) {
                applyPort(input.value);
            }
        }
    });
}

async function applyPort(port) {
    try {
        await App.SetPort(port);
        await new Promise((r) => setTimeout(r, 300));

        const status = await App.GetServerStatus();
        updateServerUI(status);
    } catch (err) {
        console.error("Failed to set port:", err);
        const statusBar = document.getElementById("status-bar");
        const statusText = document.getElementById("status-text");
        statusBar.className = "status-bar error";
        statusText.textContent = "Port error: " + err;
    }
}

// ========== Settings ==========

async function loadSettings() {
    try {
        // Clear password fields
        document.getElementById("current-password").value = "";
        document.getElementById("new-password").value = "";
        clearPasswordValidation();

        // Load exposed fields
        const fields = await App.GetExposedFields();
        exposedFields = new Set(fields);

        // Sync checkboxes
        document.querySelectorAll(".field-toggle input[type='checkbox']").forEach((cb) => {
            const field = cb.getAttribute("data-field");
            if (field) {
                cb.checked = exposedFields.has(field);
            }
        });

        const background = await App.GetBackgroundSettings();
        document.getElementById("run-in-background").checked = background.runInBackground;
        document.getElementById("start-at-login").checked = background.startAtLogin;
    } catch (err) {
        console.error("Failed to load settings:", err);
    }
}

async function setRunInBackground(enabled) {
    const checkbox = document.getElementById("run-in-background");
    const status = document.getElementById("background-status");
    try {
        await App.SetRunInBackground(enabled);
        status.textContent = "Pengaturan tray disimpan.";
        status.className = "password-status ok";
    } catch (err) {
        checkbox.checked = !enabled;
        status.textContent = "Gagal menyimpan pengaturan: " + err;
        status.className = "password-status error";
    }
}

async function setStartAtLogin(enabled) {
    const checkbox = document.getElementById("start-at-login");
    const status = document.getElementById("background-status");
    try {
        await App.SetStartAtLogin(enabled);
        status.textContent = enabled ? "Pegasus akan dimulai saat Anda login." : "Mulai otomatis dinonaktifkan.";
        status.className = "password-status ok";
    } catch (err) {
        checkbox.checked = !enabled;
        status.textContent = "Gagal mengubah mulai otomatis: " + err;
        status.className = "password-status error";
    }
}

async function toggleField(field, checked) {
    if (checked) {
        exposedFields.add(field);
    } else {
        exposedFields.delete(field);
    }

    try {
        await App.SetExposedFields(Array.from(exposedFields));
    } catch (err) {
        console.error("Failed to update exposed fields:", err);
    }
}

function setPasswordStatus(message, isOk) {
    const status = document.getElementById("password-status");
    status.className = "password-status " + (isOk ? "ok" : "error");
    status.innerHTML = (isOk ? "<span>✓</span> " : "<span>✕</span> ") + message;
}

function resetPasswordStatus() {
    const status = document.getElementById("password-status");
    status.textContent = "";
    status.className = "password-status";
}

function clearPasswordValidation() {
    document.getElementById("current-password").classList.remove("invalid");
    document.getElementById("new-password").classList.remove("invalid");
    resetPasswordStatus();
}

function setupPasswordSettings() {
    const currentInput = document.getElementById("current-password");
    const newInput = document.getElementById("new-password");

    currentInput.addEventListener("input", () => {
        currentInput.classList.remove("invalid");
        resetPasswordStatus();
    });
    newInput.addEventListener("input", () => {
        newInput.classList.remove("invalid");
        resetPasswordStatus();
    });
    currentInput.addEventListener("keydown", (e) => {
        if (e.key === "Enter") {
            newInput.focus();
        }
    });
    newInput.addEventListener("keydown", (e) => {
        if (e.key === "Enter") {
            changePassword();
        }
    });
}

async function changePassword() {
    const currentInput = document.getElementById("current-password");
    const newInput = document.getElementById("new-password");
    const currentPassword = currentInput.value;
    const newPassword = newInput.value.trim();

    clearPasswordValidation();

    if (!currentPassword) {
        currentInput.classList.add("invalid");
        currentInput.focus();
        setPasswordStatus("Masukkan kata sandi saat ini.", false);
        return;
    }

    if (newPassword.length < 4) {
        newInput.classList.add("invalid");
        newInput.focus();
        setPasswordStatus("Kata sandi baru minimal 4 karakter.", false);
        return;
    }

    if (currentPassword === newPassword) {
        newInput.classList.add("invalid");
        newInput.focus();
        setPasswordStatus("Kata sandi baru tidak boleh sama dengan kata sandi saat ini.", false);
        return;
    }

    try {
        await App.ChangePassword(currentPassword, newPassword);
        currentInput.value = "";
        newInput.value = "";
        clearPasswordValidation();
        setPasswordStatus("Kata sandi berhasil diubah.", true);
    } catch (err) {
        const errStr = String(err);
        if (errStr.includes("current password is incorrect")) {
            currentInput.classList.add("invalid");
            currentInput.focus();
            setPasswordStatus("Kata sandi saat ini salah.", false);
        } else if (errStr.includes("at least 4 characters")) {
            newInput.classList.add("invalid");
            newInput.focus();
            setPasswordStatus("Kata sandi baru minimal 4 karakter.", false);
        } else {
            setPasswordStatus("Gagal mengubah kata sandi: " + err, false);
        }
    }
}

function togglePasswordVisibility() {
    const currentInput = document.getElementById("current-password");
    const newInput = document.getElementById("new-password");
    const btn = document.getElementById("password-toggle");
    if (newInput.type === "password") {
        currentInput.type = "text";
        newInput.type = "text";
        btn.textContent = "Sembunyikan";
    } else {
        currentInput.type = "password";
        newInput.type = "password";
        btn.textContent = "Tampilkan";
    }
}

// ========== App Lock ==========

function showLockError(message) {
    const input = document.getElementById("lock-password");
    const error = document.getElementById("lock-error");

    error.textContent = message;
    input.classList.remove("shake");
    // Trigger reflow to restart shake animation
    void input.offsetWidth;
    input.classList.add("invalid", "shake");
    input.focus();
}

function clearLockError() {
    const input = document.getElementById("lock-password");
    const error = document.getElementById("lock-error");
    input.classList.remove("invalid", "shake");
    error.textContent = "";
}

async function unlockApp() {
    const input = document.getElementById("lock-password");
    const password = input.value;

    if (!password) {
        showLockError("Masukkan kata sandi terlebih dahulu.");
        return;
    }

    try {
        const ok = await App.VerifyPassword(password);
        if (ok) {
            clearLockError();
            input.value = "";
            document.getElementById("lock-screen").classList.add("hidden");
            await initState();
        } else {
            showLockError("Kata sandi salah. Silakan coba lagi.");
            input.value = "";
        }
    } catch (err) {
        showLockError("Gagal memverifikasi kata sandi: " + err);
    }
}

async function toggleLockHelper() {
    const box = document.getElementById("lock-helper-box");
    const isHidden = box.classList.toggle("hidden");

    if (!isHidden) {
        try {
            const configPath = await App.GetConfigPath();
            document.getElementById("lock-config-path").textContent = configPath || ".env";
        } catch (_) {
            document.getElementById("lock-config-path").textContent = ".env";
        }
    }
}

function setupLockScreen() {
    const input = document.getElementById("lock-password");
    input.addEventListener("keydown", (e) => {
        if (e.key === "Enter") {
            unlockApp();
        }
    });
    input.addEventListener("input", () => {
        clearLockError();
    });
}

// ========== About ==========

async function loadAboutInfo() {
    try {
        const info = await App.GetAppInfo();
        document.getElementById("about-name").textContent = info.name;
        document.getElementById("about-version").textContent = "v" + info.version;
        document.getElementById("about-desc").textContent = info.description;
        document.getElementById("about-author").textContent = info.author;
        document.getElementById("about-license").textContent = info.license;

        const link = document.getElementById("about-github");
        link.href = info.github;
        link.textContent = info.github.replace("https://github.com/", "");
    } catch (err) {
        console.error("Failed to load about info:", err);
    }
}

function openGitHub(event) {
    event.preventDefault();
    const url = document.getElementById("about-github").href;
    if (url && url !== "#") {
        Browser.OpenURL(url);
    }
}

function openEndpoint(event) {
    event.preventDefault();
    Browser.OpenURL(event.currentTarget.href);
}

Object.assign(window, {
    switchTab,
    toggleServer,
    toggleField,
    setRunInBackground,
    setStartAtLogin,
    togglePasswordVisibility,
    changePassword,
    unlockApp,
    toggleLockHelper,
    openGitHub,
    openEndpoint,
});

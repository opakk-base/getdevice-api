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
        resetPasswordStatus();

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

async function changePassword() {
    const currentInput = document.getElementById("current-password");
    const newInput = document.getElementById("new-password");
    const status = document.getElementById("password-status");
    const currentPassword = currentInput.value;
    const newPassword = newInput.value.trim();

    if (!currentPassword) {
        status.textContent = "Masukkan kata sandi saat ini.";
        status.className = "password-status error";
        return;
    }

    if (newPassword.length < 4) {
        status.textContent = "Kata sandi baru minimal 4 karakter.";
        status.className = "password-status error";
        return;
    }

    try {
        await App.ChangePassword(currentPassword, newPassword);
        currentInput.value = "";
        newInput.value = "";
        status.textContent = "Kata sandi berhasil diubah.";
        status.className = "password-status ok";
    } catch (err) {
        status.textContent = "Gagal mengubah kata sandi: " + err;
        status.className = "password-status error";
    }
}

function resetPasswordStatus() {
    const status = document.getElementById("password-status");
    status.textContent = "";
    status.className = "password-status";
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

async function unlockApp() {
    const input = document.getElementById("lock-password");
    const error = document.getElementById("lock-error");

    try {
        const ok = await App.VerifyPassword(input.value);
        if (ok) {
            input.value = "";
            document.getElementById("lock-screen").classList.add("hidden");
            await initState();
        } else {
            error.textContent = "Incorrect password.";
            input.value = "";
            input.focus();
        }
    } catch (err) {
        error.textContent = "Failed to verify password: " + err;
    }
}

function setupLockScreen() {
    const input = document.getElementById("lock-password");
    const error = document.getElementById("lock-error");
    input.addEventListener("keydown", (e) => {
        if (e.key === "Enter") {
            unlockApp();
        }
    });
    input.addEventListener("input", () => {
        error.textContent = "";
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
    openGitHub,
    openEndpoint,
});

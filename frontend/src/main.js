// Store the latest device info for copy functionality
let currentDeviceInfo = null;

// Field IDs that map to DeviceInfo struct fields
const fields = [
    "device_id",
    "device_name",
    "client_key",
    "hostname",
    "os",
    "arch",
    "mac_address",
    "ip_address",
    "timestamp",
];

// Load device info from Go backend via Wails binding
async function loadDeviceInfo() {
    const btn = document.getElementById("refresh-btn");
    const errorEl = document.getElementById("error-msg");

    btn.classList.add("loading");
    btn.textContent = "Loading...";
    errorEl.classList.add("hidden");

    try {
        // Call Go method via Wails binding
        const info = await window.go.main.App.GetDeviceInfo();
        currentDeviceInfo = info;

        // Update each field
        fields.forEach((field) => {
            const el = document.getElementById(field);
            if (el) {
                el.textContent = info[field] || "N/A";
                el.classList.remove("loading");
            }
        });

        // Update port display
        const port = await window.go.main.App.GetPort();
        document.getElementById("port").textContent = port;
    } catch (err) {
        errorEl.textContent = "Failed to load device info: " + err;
        errorEl.classList.remove("hidden");

        // Set status bar to error
        const statusBar = document.getElementById("status-bar");
        statusBar.classList.add("error");
        document.getElementById("status-text").textContent = "Error connecting";
    } finally {
        btn.classList.remove("loading");
        btn.textContent = "Refresh";
    }
}

// Copy device info as JSON to clipboard
async function copyAsJSON() {
    if (!currentDeviceInfo) return;

    const btn = document.getElementById("copy-btn");

    try {
        const json = JSON.stringify(currentDeviceInfo, null, 2);
        await navigator.clipboard.writeText(json);

        btn.classList.add("copied");
        btn.textContent = "Copied!";

        setTimeout(() => {
            btn.classList.remove("copied");
            btn.textContent = "Copy JSON";
        }, 1500);
    } catch (err) {
        btn.textContent = "Failed";
        setTimeout(() => {
            btn.textContent = "Copy JSON";
        }, 1500);
    }
}

// Load on startup
document.addEventListener("DOMContentLoaded", loadDeviceInfo);

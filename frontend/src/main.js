// Current device info (full, for display)
let currentDeviceInfo = null;

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

// All field IDs
const allFields = [
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

// Port change debounce timer
let portDebounceTimer = null;

// ========== Initialization ==========

document.addEventListener("DOMContentLoaded", async () => {
    await initState();
    await loadDeviceInfo();
    setupPortInput();
});

// Load initial state from Go backend
async function initState() {
    try {
        // Get current exposed fields
        const fields = await window.go.main.App.GetExposedFields();
        exposedFields = new Set(fields);

        // Sync checkboxes with backend state
        allFields.forEach((field) => {
            const card = document.querySelector(`[data-field="${field}"]`);
            // OS and arch share a card
            if (!card && (field === "os" || field === "arch")) return;
            const checkbox = card
                ? card.querySelector(`input[onchange*="${field}"]`)
                : null;
            if (checkbox) {
                checkbox.checked = exposedFields.has(field);
            }
            updateCardState(field);
        });

        // Also sync os/arch checkboxes in the split card
        const osArchCard = document.querySelector('[data-field="os_arch"]');
        if (osArchCard) {
            const checkboxes = osArchCard.querySelectorAll('input[type="checkbox"]');
            checkboxes.forEach((cb) => {
                const match = cb.getAttribute("onchange").match(/toggleField\('(\w+)'/);
                if (match) {
                    cb.checked = exposedFields.has(match[1]);
                }
            });
        }

        // Get server status
        const status = await window.go.main.App.GetServerStatus();
        updateServerUI(status);

        // Set port input
        document.getElementById("port-input").value = status.port;
    } catch (err) {
        console.error("Failed to init state:", err);
    }
}

// ========== Device Info ==========

async function loadDeviceInfo() {
    const btn = document.getElementById("refresh-btn");
    const errorEl = document.getElementById("error-msg");

    btn.classList.add("loading");
    btn.textContent = "Loading...";
    errorEl.classList.add("hidden");

    try {
        const info = await window.go.main.App.GetDeviceInfo();
        currentDeviceInfo = info;

        allFields.forEach((field) => {
            const el = document.getElementById(field);
            if (el) {
                el.textContent = info[field] || "N/A";
                el.classList.remove("loading");
            }
        });
    } catch (err) {
        errorEl.textContent = "Failed to load device info: " + err;
        errorEl.classList.remove("hidden");
    } finally {
        btn.classList.remove("loading");
        btn.textContent = "Refresh";
    }
}

// ========== Field Toggles ==========

async function toggleField(field, checked) {
    if (checked) {
        exposedFields.add(field);
    } else {
        exposedFields.delete(field);
    }

    updateCardState(field);

    // Send updated fields to Go backend
    try {
        await window.go.main.App.SetExposedFields(Array.from(exposedFields));
    } catch (err) {
        console.error("Failed to update exposed fields:", err);
    }
}

function updateCardState(field) {
    // Find the card for this field
    let card = document.querySelector(`[data-field="${field}"]`);

    // OS and arch are in a shared card — handle individually
    if (field === "os" || field === "arch") {
        const valueEl = document.getElementById(field);
        if (valueEl) {
            if (exposedFields.has(field)) {
                valueEl.style.opacity = "1";
                valueEl.style.textDecoration = "none";
            } else {
                valueEl.style.opacity = "0.4";
                valueEl.style.textDecoration = "line-through";
            }
        }
        return;
    }

    if (card) {
        if (exposedFields.has(field)) {
            card.classList.remove("disabled");
        } else {
            card.classList.add("disabled");
        }
    }
}

// ========== Server Controls ==========

async function toggleServer() {
    const btn = document.getElementById("server-btn");
    const btnText = document.getElementById("server-btn-text");

    try {
        const status = await window.go.main.App.GetServerStatus();

        if (status.running) {
            await window.go.main.App.StopServer();
        } else {
            await window.go.main.App.StartServer();
        }

        // Small delay for server to start/stop
        await new Promise((r) => setTimeout(r, 200));

        // Refresh status
        const newStatus = await window.go.main.App.GetServerStatus();
        updateServerUI(newStatus);
    } catch (err) {
        console.error("Server toggle error:", err);
        // Refresh status anyway
        try {
            const s = await window.go.main.App.GetServerStatus();
            updateServerUI(s);
        } catch (_) {}
    }
}

function updateServerUI(status) {
    const btn = document.getElementById("server-btn");
    const btnText = document.getElementById("server-btn-text");
    const statusBar = document.getElementById("status-bar");
    const statusText = document.getElementById("status-text");
    const portDisplay = document.getElementById("port-display");

    if (status.running) {
        btn.className = "server-btn running";
        btnText.textContent = "Stop Server";
        statusBar.className = "status-bar";
        statusText.innerHTML = 'Running on port <span id="port-display">' + status.port + "</span>";
    } else {
        btn.className = "server-btn stopped";
        btnText.textContent = "Start Server";
        statusBar.className = "status-bar stopped";
        statusText.innerHTML = "Server stopped";
    }
}

// ========== Port Input ==========

function setupPortInput() {
    const input = document.getElementById("port-input");

    input.addEventListener("input", () => {
        // Only allow digits
        input.value = input.value.replace(/[^0-9]/g, "");

        // Validate
        const port = parseInt(input.value, 10);
        if (input.value && (isNaN(port) || port < 1 || port > 65535)) {
            input.classList.add("invalid");
        } else {
            input.classList.remove("invalid");
        }

        // Debounce the port change
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
        await window.go.main.App.SetPort(port);

        // Small delay for restart
        await new Promise((r) => setTimeout(r, 300));

        const status = await window.go.main.App.GetServerStatus();
        updateServerUI(status);
    } catch (err) {
        console.error("Failed to set port:", err);
        const statusBar = document.getElementById("status-bar");
        const statusText = document.getElementById("status-text");
        statusBar.className = "status-bar error";
        statusText.textContent = "Port error: " + err;
    }
}

// ========== Copy JSON ==========

async function copyAsJSON() {
    if (!currentDeviceInfo) return;

    const btn = document.getElementById("copy-btn");

    try {
        // Build filtered JSON based on exposed fields
        const filtered = {};
        allFields.forEach((field) => {
            if (exposedFields.has(field)) {
                filtered[field] = currentDeviceInfo[field];
            }
        });

        const json = JSON.stringify(filtered, null, 2);
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

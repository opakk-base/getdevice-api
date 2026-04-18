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
    await initState();
    setupPortInput();
});

async function initState() {
    try {
        // Get server status and set port input
        const status = await window.go.main.App.GetServerStatus();
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
        const status = await window.go.main.App.GetServerStatus();

        if (status.running) {
            await window.go.main.App.StopServer();
        } else {
            await window.go.main.App.StartServer();
        }

        await new Promise((r) => setTimeout(r, 200));

        const newStatus = await window.go.main.App.GetServerStatus();
        updateServerUI(newStatus);
    } catch (err) {
        console.error("Server toggle error:", err);
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
        await window.go.main.App.SetPort(port);
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

// ========== Settings ==========

async function loadSettings() {
    try {
        // Load close behavior
        const behavior = await window.go.main.App.GetCloseBehavior();
        const radios = document.querySelectorAll('input[name="close-behavior"]');
        radios.forEach((radio) => {
            radio.checked = radio.value === behavior;
        });

        // Load exposed fields
        const fields = await window.go.main.App.GetExposedFields();
        exposedFields = new Set(fields);

        // Sync checkboxes
        document.querySelectorAll(".field-toggle input[type='checkbox']").forEach((cb) => {
            const field = cb.getAttribute("data-field");
            if (field) {
                cb.checked = exposedFields.has(field);
            }
        });
    } catch (err) {
        console.error("Failed to load settings:", err);
    }
}

async function toggleField(field, checked) {
    if (checked) {
        exposedFields.add(field);
    } else {
        exposedFields.delete(field);
    }

    try {
        await window.go.main.App.SetExposedFields(Array.from(exposedFields));
    } catch (err) {
        console.error("Failed to update exposed fields:", err);
    }
}

async function setCloseBehavior(value) {
    try {
        await window.go.main.App.SetCloseBehavior(value);
    } catch (err) {
        console.error("Failed to set close behavior:", err);
    }
}

// ========== About ==========

async function loadAboutInfo() {
    try {
        const info = await window.go.main.App.GetAppInfo();
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
        window.runtime.BrowserOpenURL(url);
    }
}

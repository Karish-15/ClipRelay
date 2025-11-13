// ==========================================
// ============= GLOBAL CONFIG ==============
// ==========================================

// Detect whether FE is being accessed locally or via Cloudflare domain.
let BASE_URL;
let SSE_DOMAINS;

if (location.hostname.endsWith(".karishch.online")) {
    // ===== PRODUCTION / CLOUDFLARE =====
    BASE_URL = "https://api.karishch.online";

    SSE_DOMAINS = {
        1: "https://sse1.karishch.online",
        2: "https://sse2.karishch.online",
        3: "https://sse3.karishch.online"
    };

} else {
    // ===== LOCAL DEVELOPMENT =====
    BASE_URL = "http://localhost:8081";

    SSE_DOMAINS = {
        1: "http://localhost:8081",
        2: "http://localhost:8081",
        3: "http://localhost:8081"
    };
}


// ==========================================
// =============== LOGIN =====================
// ==========================================

async function login() {
    const user = document.getElementById("username").value
    const pwd = document.getElementById("password").value

    const res = await fetch(`${BASE_URL}/login`, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({username: user, password: pwd})
    })

    if (!res.ok) {
        document.getElementById("error").innerText = "Invalid login"
        return
    }

    const data = await res.json()
    localStorage.setItem("jwt", data.token)

    window.location = "index.html"
}


// ==========================================
// =============== UTILS =====================
// ==========================================

function authHeaders() {
    return {
        "Authorization": localStorage.getItem("jwt"),
        "Content-Type": "application/json"
    }
}

function addClipToUI(clip, append = false) {
    const div = document.createElement("div")
    div.className = "clip"

    if (clip.in_blob) {
        const img = document.createElement("img")
        img.src = clip.blob_url
        div.appendChild(img)
    } else {
        div.innerText = clip.content
    }

    if (append)
        document.getElementById("clip-list").appendChild(div)
    else
        document.getElementById("clip-list").prepend(div)
}


// ==========================================
// =========== HOME PAGE INIT ================
// ==========================================

let lastClipId = null
let sse = null

async function init() {
    if (!localStorage.getItem("jwt")) {
        window.location = "login.html"
        return
    }

    // Fetch assigned gateway ID
    const gw = await fetch(`${BASE_URL}/sse`, {
        headers: authHeaders()
    }).then(r => r.json())

    const gatewayId = gw.gateway

    // Build SSE URL using domain-based logic
    const token = localStorage.getItem("jwt")
    const esBase = SSE_DOMAINS[gatewayId]
    const sseUrl = `${esBase}/events/${gatewayId}?token=${token}`

    console.log("Connecting SSE →", sseUrl)

    sse = new EventSource(sseUrl)

    // Default unnamed SSE events
    sse.onmessage = (event) => {
        try {
            const clip = JSON.parse(event.data)
            addClipToUI(clip)
        } catch (_) {}
    }

    // Named SSE event (new clip)
    sse.addEventListener("new_clip", (event) => {
        const clip = JSON.parse(event.data)
        addClipToUI(clip)
    })

    // Load latest clip on page load
    const latestClip = await fetch(`${BASE_URL}/clips/latest`, {
        headers: authHeaders()
    }).then(r => r.json())

    addClipToUI(latestClip)
    lastClipId = latestClip.id
}

// Only run init on the index page
if (location.pathname.includes("index.html")) {
    window.onload = init
}


// ==========================================
// ========== HISTORY LOADING ================
// ==========================================

async function loadHistory() {
    if (!lastClipId) return

    const url = `${BASE_URL}/clips?before=${lastClipId}`
    const data = await fetch(url, { headers: authHeaders() }).then(r => r.json())

    if (data.length === 0) {
        alert("No more history")
        return
    }

    data.forEach(c => addClipToUI(c, true))

    lastClipId = data[data.length - 1].id
}


// ==========================================
// ======= PASTE IMAGE/TEXT HANDLING =========
// ==========================================

document.addEventListener("paste", async (event) => {
    const jwt = localStorage.getItem("jwt")
    if (!jwt) return

    const items = event.clipboardData.items

    // IMAGE
    for (const item of items) {
        if (item.type.startsWith("image/")) {
            const file = item.getAsFile()

            // STEP 1: Init blob
            const initRes = await fetch(`${BASE_URL}/clips/blob/init`, {
                method: "POST",
                headers: authHeaders(),
                body: JSON.stringify({
                    mime_type: file.type
                })
            })

            const initData = await initRes.json()
            const uploadUrl = initData.upload_url
            const clipId = initData.id

            // STEP 2: Upload to B2 / S3
            await fetch(uploadUrl, {
                method: "PUT",
                body: file
            })

            // STEP 3: Finalize
            await fetch(`${BASE_URL}/clips/blob`, {
                method: "POST",
                headers: authHeaders(),
                body: JSON.stringify({id: clipId})
            })

            return
        }
    }

    // TEXT
    const text = event.clipboardData.getData("text")
    if (text.trim().length > 0) {
        await fetch(`${BASE_URL}/clips/text`, {
            method: "POST",
            headers: authHeaders(),
            body: JSON.stringify({content: text})
        })
    }
})

window.showToast = function showToast(message, type) {
    const toast = document.createElement("div")
    toast.className = `toast toast-${type}`
    toast.innerHTML = message
    toast.style.cssText = `
        margin-bottom: 0.5rem;
        padding: 1rem;
        border-radius: 0.375rem;
        min-width: 20rem;
        animation: slideIn 0.2s ease-in-out;
        color: white;
        box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
    `
    // Set background color based on type
    switch (type) {
        case "success":
            toast.style.backgroundColor = "#10B981"
            break
        case "error":
            toast.style.backgroundColor = "#EF4444"
            break
        case "warning":
            toast.style.backgroundColor = "#F59E0B"
            break
        case "info":
            toast.style.backgroundColor = "#3B82F6"
            break
    }

    document.getElementById("toast-container").appendChild(toast)

    // Auto remove after 5 seconds
    setTimeout(() => {
        // 250 > 200 to prevent flickering
        toast.style.animation = "slideOut 0.25s ease-in-out"
        setTimeout(() => toast.remove(), 200)
    }, 5000)
}

// Add required CSS animations
if (!document.getElementById("toast-styles")) {
    const style = document.createElement("style")
    style.id = "toast-styles"
    style.textContent = `
        @keyframes slideIn {
            from { transform: translateX(100%); opacity: 0; }
            to { transform: translateX(0); opacity: 1; }
        }
        @keyframes slideOut {
            from { transform: translateX(0); opacity: 1; }
            to { transform: translateX(100%); opacity: 0; }
        }
    `
    document.head.appendChild(style)
}

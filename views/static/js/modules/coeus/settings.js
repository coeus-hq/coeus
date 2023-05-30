if (window.location.pathname === "/settings") {
    setSelectedTimezone();
}

export function setSelectedTimezone() {
    const timezoneOffset = document.getElementById("setting-timezone-offset").value;
    const timezoneSelect = document.getElementById("timezone");

    for (let i = 0; i < timezoneSelect.options.length; i++) {
        if (timezoneSelect.options[i].value === timezoneOffset) {
            timezoneSelect.selectedIndex = i;
            break;
        }
    }
}
if (window.location.pathname === "/settings") {
    setSelectedTimezone();

    // Get the value of the timezone offset from the hidden input
let timezoneOffset = document.getElementById("setting-timezone-offset").value;
// Convert the offset to a text string based on the offset value
let timezoneOffsetText = convertTimezoneOffsetToText(timezoneOffset);
// Set the text of the timezone info span to the offset text
document.getElementById("settings-timezone-info").innerText = timezoneOffsetText;
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


// make a switch statement to convert the offset to a text string
export function convertTimezoneOffsetToText(offset) {
    switch (offset) {
        case "-480":
            return "PST (UTC-8)";
        case "-420":
            return "MST (UTC-7)";
        case "-360":
            return "CST (UTC-6)";
        case "-300":
            return "CDT (UTC-5)";
        case "-240":
            return "EDT (UTC-4)";
        case "0":
            return "UTC";
        case "60":
            return "CET (UTC+1)";
        case "120":
            return "EET (UTC+2)";
        case "330":
            return "IST (UTC+5:30)";
        default:
            return "UTC";
    }
}

// When the edit settings button is clicked, toggle the display of the settings form
export function toggleSettingsForm() {
    document.getElementById("settings-form").classList.toggle("hidden");
    document.getElementById("settings-info").classList.toggle("hidden");
}


// When the dark mode toggle is clicked send a post request to the server to toggle the dark mode setting
export function toggleDarkTheme() {
    let darkModeToggle = document.getElementById("darkModeToggle");

    fetch("/api/settings/dark-theme", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
    })
        .then(response => {
            if (response.ok) {
                location.reload();
            } else {
                console.error('Server response was not ok.');
            }
        })
        .catch(err => {
            console.error('There was a problem with the fetch operation: ' + err.message);
        });
}

// updateUserSettings() is called when the "Save Changes" button in Settings is clicked
export function updateUserSettings(e) {
    e.preventDefault();

    const userId = document.getElementById("updateUserSettingsButton").getAttribute("data-user-id");
    const userFirstName = document.getElementById("firstName").value;
    const userLastName = document.getElementById("lastName").value;
    const userEmail = document.getElementById("email").value;

    // Create a JSON object to send in the request body
    const userData = {
        userId: userId,
        firstName: userFirstName,
        lastName: userLastName,
        email: userEmail,

    };

    // Send a PUT request to the server
    fetch("/api/user", {
        method: "PUT",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(userData)
    })
        .then((response) => {
            if (response.ok) {

                updateTimezone()
                // If the response is OK, reload the page

                // wait for the timezone to update before reloading the page
                setInterval(function () {
                    location.reload();
                }, 500);


            } else {
                // Otherwise, display an error message
                alert("Unable to update user settings");
            }
        });
}

export function updateTimezone() {
    const userTimezone = document.getElementById("timezone").value;

    // Create a JSON object to send in the request body
    const userData = {
        timezone: userTimezone,
    };

    // Send a PUT request to the server
    fetch("/api/settings/timezone", {
        method: "PUT",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(userData)
    })

}

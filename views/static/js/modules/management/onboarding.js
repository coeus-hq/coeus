// logoInputChange is used for the logo upload preview
export function logoInputChange(e) {
    var file = e.target.files[0];
    if (file) {
        var reader = new FileReader();
        reader.onloadend = function () {
            document.querySelector('#logo-display .logo-preview').src = reader.result;
            document.getElementById('logo-display').style.display = 'flex';
            document.getElementById('upload-logo-input-container').style.display = 'none';
        }
        reader.readAsDataURL(file);
    }
}

// toggleAPIKeyInfo toggles the API key info
export function toggleAPIKeyInfo() {

    var infoElement = document.getElementById('api-key-helper-info-hidden');

    if (infoElement.style.display === "none") {
        infoElement.style.display = "block";
    } else {
        infoElement.style.display = "none";
    }
}

// toggleAPIKeyInput toggles the API key input
export function toggleAPIKeyInput() {
    var apiKeyInput = document.getElementById('apiKeyInput');
    var checkbox = document.getElementById('toggleAPIKey');

    // If checked
    if (checkbox.checked) {
        // Show the API key input
        apiKeyInput.classList.remove('invisible');
        apiKeyInput.classList.add('visible');
    } else {
        // Hide the API key input
        apiKeyInput.classList.remove('visible');
        apiKeyInput.classList.add('invisible');
    }
}

// createAdminAccount handles the form submission for creating an admin account
export function createAdminAccount(e) {
    e.preventDefault();

    let firstName = document.getElementById("create-admin-account").elements.namedItem("firstName").value.trim();
    let lastName = document.getElementById("create-admin-account").elements.namedItem("lastName").value.trim();
    let email = document.getElementById("create-admin-account").elements.namedItem("email").value.trim();
    let password = document.getElementById("create-admin-account").elements.namedItem("password").value.trim();
    let confirmPassword = document.getElementById("create-admin-account").elements.namedItem("confirmPassword").value.trim();

    // Match passwords
    if (password !== confirmPassword) {
        alert("Passwords do not match!");
        return;
    }

    // Make sure all fields are filled
    if (firstName === "" || lastName === "" || email === "" || password === "" || confirmPassword === "") {

        document.getElementById('admin-form-error-message').style.display = 'block';

        setTimeout(function () {
            document.getElementById('admin-form-error-message').style.display = 'none';
        }, 3000);

        return;
    }

    let newAdmin = {
        firstName: firstName,
        lastName: lastName,
        email: email,
        password: password
    };

    fetch('/api/admin', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(newAdmin)
    })
        .then(response => response.json())
        .then(response => {
            // Handle the response
            if (response.success) {
                // Show the next step in the onboarding process
                document.getElementById('create-admin-account').style.display = 'none';
                document.getElementById('onboarding-form').style.display = 'flex';

                let firstName = response.adminName
                document.getElementById('create-account-success-header').textContent = `${firstName}, admin status is set`;
            } else {
                // Admin-account-error-message
                document.getElementById('admin-account-error-message').style.display = 'block';

                setTimeout(function () {
                    document.getElementById('admin-account-error-message').style.display = 'none';
                }, 3000);

            }
        })
}

// submitOnboarding handles the form submission for onboarding
export function submitOnboarding(e) {
    e.preventDefault();

    // Get the values of the form
    let organizationName = document.getElementById("onboarding-organization-name").value.trim();
    let organizationTimezone = document.getElementById("onboarding-organization-timezone").value.trim();
    let organizationLogo = document.getElementById("upload-logo-input").files[0]; // Get the File object
    let apiKey = document.getElementById("apiKey").value.trim();
    let email = document.getElementById("sendGridEmail").value.trim();

    // Make sure all fields are filled except for the optional API key
    if (organizationName === "" || organizationTimezone === "" || !organizationLogo) {
        document.getElementById('onboarding-form-error-message').style.display = 'block';
        setTimeout(function () {
            document.getElementById('onboarding-form-error-message').style.display = 'none';
        }, 3000);

        return;
    }

    // Create a FormData object to send in the request body
    let formData = new FormData();
    formData.append('name', organizationName);
    formData.append('timezone', organizationTimezone);
    formData.append('logoPath', organizationLogo);  
    formData.append('apiKey', apiKey);
    formData.append('email', email);


    fetch('/api/onboarding', {
        method: 'POST',
        body: formData
    })
        .then(response => response.json())
        .then(response => {
            // Handle the response
            if (response.success) {
                window.location.href = "/sign-in";
            } else {
                // Error handling
                alert("Error: " + response.message);
            }
        })
}
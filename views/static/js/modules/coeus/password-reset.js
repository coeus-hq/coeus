// These functions hide and show the alerts for the password reset forms
// PASSWORD RESET FORMS START
function showEmailSuccessAlert() {
    document.getElementById("email-success-alert").style.display = "block";

    setTimeout(function () {
        document.getElementById("email-success-alert").style.display = "none";
    }, 5000);
}

function showEmailErrorAlert() {
    document.getElementById("email-error-alert").style.display = "block";

    setTimeout(function () {
        document.getElementById("email-error-alert").style.display = "none";
    }, 3000);
}

function showPinSuccessAlert() {
    document.getElementById("pin-success-alert").style.display = "block";

    setTimeout(function () {
        document.getElementById("pin-success-alert").style.display = "none";
    }, 5000);
}

function showPinErrorAlert() {
    document.getElementById("pin-error-alert").style.display = "block";

    setTimeout(function () {
        document.getElementById("pin-error-alert").style.display = "none";
    }, 3000);
}

function showPasswordResetSuccessAlert() {
    document.getElementById("password-reset-success-alert").style.display = "block";

    setTimeout(function () {
        document.getElementById("password-reset-success-alert").style.display = "none";
    }, 5000);
}

function showPasswordResetErrorAlert() {
    document.getElementById("password-reset-error-alert").style.display = "block";

    setTimeout(function () {
        document.getElementById("password-reset-error-alert").style.display = "none";
    }, 3000);
}
// PASSWORD RESET FORMS END

// sendResetEmail sends the email to the server to get verified the the server sends a pin to the email 
export function sendResetEmail(e) {
    e.preventDefault();

    // Get the email from the form
    const email = document.getElementById("forgot-password-email").value;

    // check if the email is valid
    if (!email || !email.includes("@") || !email.includes(".")) {
        showEmailErrorAlert()
        return;
    }

    // Submit the form data to the server
    fetch("/api/password-reset/send-email", {
        method: "POST",
        body: JSON.stringify({ email }),
        headers: {
            "Content-Type": "application/json",
        },
    })
        .then((response) => {
            // If the server returns a 200 status code, show the success alert
            if (response.status == 200) {

                // Hide the first form and show the second form
                document.getElementById("reset-password-form-1").style.display = "none";
                document.getElementById("reset-password-form-2").style.display = "block";
                showEmailSuccessAlert();
            } else {

                // If the server returns a 400 status code, show the error alert
                showEmailErrorAlert()
            }
        }
        )
        .catch((error) => console.log(error));
}

// verifyPin sends the pin to the server to get verified and if the pin is correct, the user can reset their password
export function verifyPin(e) {
    e.preventDefault();

    // Get the pin from the form
    const pinInputs = document.querySelectorAll(".pin-input");
    let pin = "";

    pinInputs.forEach((input) => {
        pin += input.value;
    });

    // check if the pin is valid
    if (!pin) {
        showPinErrorAlert()
        return;
    }

    // Submit the form data to the server
    fetch("/api/password-reset/verify-pin", {
        method: "POST",
        body: JSON.stringify({ pin }),
        headers: {
            "Content-Type": "application/json",
        },
    })
        .then((response) => {

            // If the server returns a 200 status code, show the success alert
            if (response.status == 200) {

                // Hide the second form and show the third form
                document.getElementById("reset-password-form-2").style.display = "none";
                document.getElementById("reset-password-form-3").style.display = "block";

                showPinSuccessAlert();
            } else {

                // If the server returns a 400 status code, show the error alert
                showPinErrorAlert()
            }
        }
        )
        .catch((error) => console.log(error));
}

// resetPasswordAndSignIn sends the new password to the server to be reset and then logs the user in
export function resetPasswordAndSignIn(e) {
    e.preventDefault();

    // Get the email from the form
    const email = document.getElementById("forgot-password-email").value;

    // Get the password from the form
    const password = document.getElementById("reset-password").value;
    const confirmPassword = document.getElementById("confirm-reset-password").value;

    // check if the password is valid
    if (!password || password.length < 3 || password !== confirmPassword) {
        showPasswordResetErrorAlert()
        return;
    }

    // Submit the form data to the server
    fetch("/api/password-reset", {
        method: "POST",
        body: JSON.stringify({ password, email }),
        headers: {
            "Content-Type": "application/json",
        },
    })
        .then((response) => {
            // If the server returns a 200 status code, show the success alert
            if (response.status == 200) {

                // Set a flag in sessionStorage to show the success alert on the sign-in page
                sessionStorage.setItem('passwordResetSuccess', 'true');
                window.location.href = "/sign-in";
            } else {

                // If the server returns a 400 status code, show the error alert
                showPasswordResetErrorAlert()
            }
        }
        )
        .catch((error) => console.log(error));
}

// Check if the current URL is "http://localhost:8080/sign-in" and if it is, check if the password reset was successful and show the success alert
if (window.location.href.includes("/sign-in")) {
    // Check if the password reset was successful
    if (sessionStorage.getItem('passwordResetSuccess') === 'true') {
        // Show the success alert
        showPasswordResetSuccessAlert();
        // Remove the flag from sessionStorage
        sessionStorage.removeItem('passwordResetSuccess');
    }
}


// Check if the current URL is "http://localhost:8080/forgot-password"
if (window.location.href.includes("/forgot-password")) {
    // This function adds an event listener to each pin input and when the user types in a number, it focuses on the next input
    const pinInputs = document.querySelectorAll('.pin-input');

    pinInputs.forEach(input => {
        input.addEventListener('input', (e) => {
            const index = parseInt(e.target.dataset.index);
            if (index < pinInputs.length - 1 && e.target.value !== '') {
                pinInputs[index + 1].focus();
            }
            if (e.target.value === '') {
                if (index > 0) {
                    pinInputs[index - 1].focus();
                }
            }
        });
    });

    // Add a 'paste' event listener on the first pin input
    pinInputs[0].addEventListener('paste', (e) => {
        // Prevent the default paste behavior
        e.preventDefault();
        // Get the pasted text from the clipboard
        const pastedText = (e.clipboardData || window.clipboardData).getData('text');
        // Split the pasted text into individual digits
        const digits = pastedText.split('');
        // Populate the input fields with the digits
        for (let i = 0; i < pinInputs.length && i < digits.length; i++) {
            if (!isNaN(digits[i])) {
                pinInputs[i].value = digits[i];
            }
        }
    });
}


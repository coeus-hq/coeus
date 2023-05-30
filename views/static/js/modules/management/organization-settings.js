// onUploadLogoInputChange displays the delete logo button if the upload logo input has a value
export function onUploadLogoInputChange() {
    const deleteLogoBtn = document.getElementById("delete-logo-btn");
    const uploadLogoInput = document.getElementById("upload-logo-input");

    if (uploadLogoInput.value) {
        deleteLogoBtn.style.display = "block";
    } else {
        deleteLogoBtn.style.display = "none";
    }
}

// onDeleteLogoButtonClick deletes the logo input value
export function onDeleteLogoButtonClick() {
    const deleteLogoBtn = document.getElementById("delete-logo-btn");
    const uploadLogoInput = document.getElementById("upload-logo-input");

    uploadLogoInput.value = "";
    deleteLogoBtn.style.display = "none";
}

// submitForm submits the settings form data to the server
export async function submitOrgSettingsForm() {
    const form = document.getElementById('settings-form');
    const formData = new FormData(form);
    const uploadLogoInput = document.getElementById("upload-logo-input");

    try {
        const response = await fetch('/api/organization', {
            method: 'POST',
            body: formData
        });

        const successToast = document.getElementById("settingsSavedSuccess");
        const failToast = document.getElementById("settingsSavedFail");

        if (response.ok) {
            // Clear choose file input
            const fileInput = document.getElementById("upload-logo-input");
            fileInput.value = "";
            successToast.classList.add("show");
            setTimeout(function () {
                successToast.classList.remove("show");
            }, 3000);
        } else {
            failToast.classList.add("show");
            setTimeout(function () {
                failToast.classList.remove("show");
            }, 3000);
            console.error('Error:', response.statusText);
        }
    } catch (error) {
        const failToast = document.getElementById("settingsSavedFail");
        failToast.classList.add("show");
        setTimeout(function () {
            failToast.classList.remove("show");
        }, 3000);
        console.error('Error:', error);
    }
}
// databaseReset resets the database fetching to /api/database/reset
export async function databaseReset() {
    const response = await fetch('/api/database/reset');

    if (response.ok) {
        const successToast = document.getElementById("successToast");
        const successToastBody = successToast.querySelector('.toast-body');
        successToast.classList.add("show");
        successToastBody.textContent = "Database reset success. Relocating to onboarding page.";
        setTimeout(function () {
            successToast.classList.remove("show");
            // Relocate window to the onboarding page
            window.location.href = "/onboarding";
        }, 3000);
    } else {
        const failToast = document.getElementById("failToast");
        const failToastBody = failToast.querySelector('.toast-body');
        failToast.classList.add("show");
        failToastBody.textContent = "Database reset failed. Please try again.";
        setTimeout(function () {
            failToast.classList.remove("show");
        }, 3000);
    }
}

// databaseSeed resets the database fetching to /api/database/reset
export async function databaseSeed() {
    const response = await fetch('/api/database/seed');

    if (response.ok) {
        const successToast = document.getElementById("successToast");
        const successToastBody = successToast.querySelector('.toast-body');
        successToast.classList.add("show");
        successToastBody.textContent = "Database seed success.";
        setTimeout(function () {
            successToast.classList.remove("show");
        }, 3000);

        // Relocate window to the onboarding page to re-login
        // window.location.href = "/logout";
    } else {
        const failToast = document.getElementById("failToast");
        const failToastBody = failToast.querySelector('.toast-body');
        failToast.classList.add("show");
        failToastBody.textContent = "Database seed failed. Please try again.";
        setTimeout(function () {
            failToast.classList.remove("show");
        }, 3000);
    }
}

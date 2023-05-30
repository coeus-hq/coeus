// openEditUserModal() is called when the "Edit" button in the Users table is clicked
export function openEditUserModal(button) {
    const userId = button.getAttribute("data-user-id");
    const userFirstName = button.getAttribute("data-user-firstname");
    const userLastName = button.getAttribute("data-user-lastname");
    const userEmail = button.getAttribute("data-user-email");
    const userPassword = button.getAttribute("data-user-password");
    const userPasswordConfirm = button.getAttribute("data-user-password-confirm");

    // Set user data in the Edit User Modal
    document.getElementById("editModalUserName").textContent = `${userFirstName} ${userLastName}`;
    document.getElementById("editModalEmail").textContent = userEmail;
    document.getElementById("editUserFirstName").value = userFirstName;
    document.getElementById("editUserLastName").value = userLastName;
    document.getElementById("editUserEmail").value = userEmail;

    // Add user ID to the "Save Changes" button (or another appropriate element) in the Edit User Modal
    document.getElementById("updateUserModalButton").setAttribute("data-user-id", userId);
}

// openDeleteUserModal() is called when the "Delete" button in the Users table is clicked
export function openDeleteUserModal(button) {
    const userId = button.getAttribute("data-user-id");
    const userFirstName = button.getAttribute("data-user-firstname");
    const userLastName = button.getAttribute("data-user-lastname");
    const userEmail = button.getAttribute("data-user-email");

    // Set user data in the Delete User Modal
    const userNameElement = document.querySelector("#deleteModalName");
    const userEmailElement = document.querySelector("#deleteModalEmail");

    userNameElement.textContent = `${userFirstName} ${userLastName}`;
    userEmailElement.textContent = userEmail;

    // Add user ID to the "Delete" button (or another appropriate element) in the Delete User Modal
    document.getElementById("modalDeleteBtn").setAttribute("data-user-id", userId);
}

// deleteUser() is called when the "Delete" button in the Delete User Modal is clicked
export function deleteUser() {
    const userId = document.getElementById("modalDeleteBtn").getAttribute("data-user-id");

    // Send a DELETE request to the server
    fetch(`/api/user/${userId}`, {
        method: "DELETE",
    })
        .then((response) => {
            if (response.ok) {
                // If the response is OK, reload the page
                location.reload();
                document.getElementById("deleteUserConfirmToast").style.display = "block";
                setTimeout(function () {
                    document.getElementById("deleteUserConfirmToast").style.display = "none";
                }, 3000);
            } else {
                // Otherwise, display an error message
                document.getElementById("deleteUserErrorToast").style.display = "block";
                setTimeout(function () {
                    document.getElementById("deleteUserErrorToast").style.display = "none";
                }, 3000);
            }
        });
}

// updateUserModal() is called when the "Save Changes" button in the Edit User Modal is clicked
export function updateUserModal(e) {
    e.preventDefault();

    const userId = document.getElementById("updateUserModalButton").getAttribute("data-user-id");
    const userFirstName = document.getElementById("editUserFirstName").value;
    const userLastName = document.getElementById("editUserLastName").value;
    const userEmail = document.getElementById("editUserEmail").value;
    const userPassword = document.getElementById("editUserPassword").value;
    const userPasswordConfirm = document.getElementById("editUserPasswordConfirm").value;
    const moderatorType = document.getElementById("moderatorSelect").value;

    // Check if the password and password confirmation fields match
    if (userPassword !== userPasswordConfirm) {
        document.getElementById("editUserPasswordError").classList.add("show");
        setTimeout(function () {
            document.getElementById("editUserPasswordError").classList.remove("show");
        }, 3000);
        return;
    }

    // Create a JSON object to send in the request body
    const userData = {
        userId: userId,
        firstName: userFirstName,
        lastName: userLastName,
        email: userEmail,
        password: userPassword,
        moderatorType: moderatorType
    };

    // Only include the password in the JSON object if it's not empty
    if (userPassword) {
        userData.password = userPassword;
    }

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
                // If the response is OK, reload the page
                location.reload();
                document.getElementById("editUserConfirmToast").style.display = "block";
                setTimeout(function () {
                    document.getElementById("editUserConfirmToast").style.display = "none";
                }, 3000);
            } else {
                // Otherwise, display an error message
                document.getElementById("editUserErrorToast").style.display = "block";

                setTimeout(function () {
                    document.getElementById("editUserErrorToast").style.display = "none";
                }, 3000);
            }
        });
}

// addUser() is called when the "Add User" button in the Add User Modal is clicked
export function addUser(e) {
    e.preventDefault();

    const userFirstName = document.getElementById("addUserFirstName").value;
    const userLastName = document.getElementById("addUserLastName").value;
    const userEmail = document.getElementById("addUserEmail").value;
    const userPassword = document.getElementById("addUserPassword").value;
    const userPasswordConfirm = document.getElementById("addUserPasswordConfirm").value;
    const moderatorType = document.getElementById("addModeratorSelect").value;

    // Check if the password and password confirmation fields match
    if (userPassword !== userPasswordConfirm) {
        document.getElementById("addUserPasswordError").classList.add("show");
        setTimeout(function () {
            document.getElementById("addUserPasswordError").classList.remove("show");
        }, 3000);
        return;
    }

    // If form is not filled out correctly, return
    if (userFirstName === "" || userLastName === "" || userEmail === "" || userPassword === "" || userPasswordConfirm === "") {
        document.getElementById("addUserFormFillError").classList.add("show");
        setTimeout(function () {
            document.getElementById("addUserFormFillError").classList.remove("show");
        }, 3000);
        return;
    }

    // set data-mdb-dismiss="modal" in the button to close the modal after the user is added
    document.getElementById("add-user-modal").setAttribute("data-mdb-dismiss", "modal");

    // Create a JSON object to send in the request body
    const userData = {
        firstName: userFirstName,
        lastName: userLastName,
        email: userEmail,
        password: userPassword,
        moderatorType: moderatorType
    };

    // Send a POST request to the server
    fetch("/api/user", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(userData)
    })
        .then((response) => {
            if (response.ok) {
                // If the response is OK, reload the page
                location.reload();
                document.getElementById("addUserConfirmToast").classList.add("show");
                setTimeout(function () {
                    document.getElementById("addUserConfirmToast").classList.remove("show");
                }, 3000);
            } else {
                // Otherwise, display an error message
                document.getElementById("addUserErrorToast").style.display = "block";
                setTimeout(function () {
                    document.getElementById("addUserErrorToast").style.display = "none";
                }, 3000);
            }
        });
}
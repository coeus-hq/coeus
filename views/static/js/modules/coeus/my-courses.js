// Functions related to course moderation
export function courseRowClick(courseRow) {
    // Get the section ID from the clicked course row
    const sectionID = courseRow.dataset.sectionid;

    // Hide the my-courses-section div
    document.getElementById("my-courses-section").style.display = "none";

    // Remove the enter-live-session class from the chat div for this section
    document.querySelector(`#enter-live-session-${sectionID}`).classList.remove("enter-live-session");

}

// Toggle the edit mode for the course rows
export function toggleTrashCourse() {

    // find all divs with class my-course-card-trash-hidden
    let trashDivs = document.querySelectorAll('.my-course-card-trash');
    let wrapper = document.querySelectorAll('.card-row-trash-wrapper');

    // toggle the hidden class on all trash divs
    trashDivs.forEach(div => {
        //remove the hidden class
        div.classList.toggle('my-course-card-trash-hidden');
    });
    wrapper.forEach(div => {
        div.classList.toggle('d-flex');
    });

}

export function getAllModerators(param) {
    let sectionID;

    // param can be a DOM element, an ID string, or a number
    if (typeof param === 'object') {
        sectionID = param.dataset.sectionid;
    } else if (typeof param === 'string' || typeof param === 'number') {
        sectionID = param;
    } else {
        console.error('Unexpected parameter type in getAllModerators:', typeof param);
        return;
    }

    // Fetch all the moderators for the given section
    fetch(`/api/moderators/${sectionID}`, {
        method: 'GET',
    })
        .then(function (response) {
            // If the API request was successful, display the moderators in the chat
            if (response.status == 200) {
                // Remove double quotes from the sectionID string and convert it to a number


                // Get the chat moderator list and clear it
                const moderatorList = document.getElementById(`chat-moderator-list-${sectionID}`);
                moderatorList.innerHTML = '';

                // Parse the response data and display the moderators
                response.json().then(function (data) {
                    const moderators = data.moderators;

                    if (!moderators || moderators.length == 0) {
                        // Display a message if there are no moderators for this section
                        const newModerator = document.createElement("li");
                        newModerator.innerHTML = "No moderators";
                        moderatorList.appendChild(newModerator);
                    } else {
                        moderators.forEach(function (moderator) {
                            // Create a new moderator element and add it to the chat list
                            const newModerator = document.createElement("div");
                            newModerator.innerHTML = `
                                <li class="mod-li chat-mod-font-secondary list-group-item d-flex justify-content-between align-items-start">
                                    <div class="ms-2 me-auto me-2">
                                        ${moderator.FirstName} ${moderator.LastName}
                                    </div>
                                    <span class="badge-font badge badge-primary rounded-pill">${moderator.UserType}</span>
                                    <span class="remove-mod-button" data-section-id="${moderator.SectionID}" data-user-id="${moderator.ID}" onclick="removeModeratorClick(this)">X</span>
                                </li>
                            `;
                            moderatorList.appendChild(newModerator);
                        });
                    }
                });
            }
        });
}


// Remove a moderator from the section
export function removeModeratorClick(removeModeratorButton) {
    const sectionID = removeModeratorButton.getAttribute("data-section-id");
    const userID = removeModeratorButton.getAttribute("data-user-id");

    fetch(`/api/remove-moderator/${userID}/${sectionID}`, {
        method: 'DELETE',
    })
        .then(function (response) {
            if (response.status == 200) {
                // Remove the moderator from the list of moderators
                getAllModerators(sectionID);
            } else {
                // Add an error message
                console.log("error removing moderator");
            }
        });
}

// Add a moderator to the section by email
export function addModeratorClick(addModeratorButton) {
    const sectionID = addModeratorButton.value;
    const email = document.querySelector(`#add-moderator-${sectionID}`).value.trim();

    // Check to see if the email is valid
    if (email == "" || email == null) {
        // Add an error message
        let errorMessage = document.getElementById(`add-moderator-form-error-message-${sectionID}`);
        errorMessage.style.display = "block";
        setTimeout(function () {
            errorMessage.style.display = "none";
        }, 3000);
        return;
    }

    // Clear the input field
    document.querySelector(`#add-moderator-${sectionID}`).value = '';

    fetch(`/api/add-moderator/${email}/${sectionID}`, {
        method: 'PUT',
    })
        .then(function (response) {
            if (response.status == 200) {
                // Add a success message
                let successMessage = document.getElementById(`add-moderator-success-message-${sectionID}`);
                successMessage.style.display = "block";
                setTimeout(function () {
                    successMessage.style.display = "none";
                }, 3000);

                // Add the new moderator to the list of moderators
                getAllModerators(sectionID);
            } else {
                // Add an error message
                let errorMessage = document.getElementById(`add-moderator-error-message-${sectionID}`);
                errorMessage.style.display = "block";
                setTimeout(function () {
                    errorMessage.style.display = "none";
                }, 3000);
            }
        });
}

// Start a class session
export function startSessionClick(startSessionButton, classSessionID, sectionID) {

    fetch(`/api/start-session/${sectionID}`, {
        method: 'POST',
    }).then(function (response) {
        if (response.status == 200) {
            // Relocate the user after a successful API call
            window.location.href = `/class-session/${sectionID}/${classSessionID}`;
        } else {
            // Add an error message
            let errorMessage = document.getElementById(`start-session-error-message-${sectionIDValue}`);
            errorMessage.style.display = "block";
            setTimeout(function () {
                errorMessage.style.display = "none";
            }, 3000);
        }
    })
        .catch((error) => {
            console.error('Error:', error);
        });
}

// End a class session
export function endSession(classSessionIDInt) {

    fetch(`/api/end-session/${classSessionIDInt}`, {
        method: 'POST',
    }).then(function (response) {
        if (response.status == 200) {
            // Relocate the user after a successful API call
            window.location.href = `/`;
        } else {
            // Add an error message
            let errorMessage = document.getElementById(`end-session-error-message-${classSessionIDInt}`);
            errorMessage.style.display = "block";
            setTimeout(function () {
                errorMessage.style.display = "none";
            }, 3000);
        }
    })
        .catch((error) => {
            console.error('Error:', error);
        });
}

// getAttendanceID gets the attendance id for the section that was clicked then marks the student present
export function getAttendanceID(joinBtn) {
    let attendanceID

    // get the value of the data attribute of the <a> tag that was clicked
    const sectionID = joinBtn.getAttribute("data-sectionid");

    // make a fetch to /api/attendance/id/:sectionID" to get the attendanceID
    fetch(`/api/attendance/id/${sectionID}`, {
        method: 'GET',
    })
        .then(function (response) {
            // Check if the response is ok
            if (response.status == 200) {
                // Convert the response to JSON

                return response.json();
            } else {
                // Add an error message
                throw new Error("Error getting attendanceID");
            }
        })
        .then(function (jsonData) {
            // Display the attendanceID in an alert
            attendanceID = jsonData.attendanceID

            markPresent(attendanceID);
        })
        .catch(function (error) {
            // Handle any errors that occurred during the fetch or JSON conversion
            console.error("Error:", error);
        });
}

// markPresent marks the student present for the attendanceID
export function markPresent(attendanceID) {

    // fetch request to the server to post the attendanceID
    fetch(`/api/attendance/mark-present/${attendanceID}`, {
        method: 'POST',
    })
        .then(function (response) {
            if (response.status == 200) {
                // placeholder for now
            } else {
                // Add an error message
            }
        }
        );
}

// Get all delete buttons
export function deleteButtonClick(button) {
    const sectionID = button.value;
    // Send DELETE request to server
    fetch(`/${sectionID}`, {
        method: 'DELETE',
    })
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to delete course');
            }

            // Reload page to reflect changes
            location.reload();
        })
        .catch(error => {
            console.error(error);
            alert('Failed to delete course');
        });
}
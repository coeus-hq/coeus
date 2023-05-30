// initializeMDBInputs re-initializes the mdb input classes because I can't get them to work with dynamically generated inputs
export function initializeMDBInputs() {
    const inputs = document.querySelectorAll('.form-outline');
    for (let input of inputs) {
        new mdb.Input(input);
    }
}

// generateScheduleInputs to dynamically generate schedule inputs in the add course modal
export function generateScheduleInputs() {
    let courseSections = parseInt(document.getElementById("course-sections").value);
    let scheduleInputsContainer = document.getElementById("scheduleInputsContainer");
    scheduleInputsContainer.innerHTML = "<h4 class='my-5'>Add section scheduals:</h4> "; // Clear the container

    for (let i = 1; i <= courseSections; i++) {
        let scheduleInputGroup = document.createElement("div");
        scheduleInputGroup.classList.add("d-flex", "flex-column", "mb-4");
        scheduleInputGroup.innerHTML = `
        <h5>Section ${i} Schedule</h5>
        <div class="d-flex">
            <div class="form-outline mr10">
                <input required  type="text" id="section-${i}-day" name="section-${i}-day" class="form-control" />
                <label class="form-label" for="section-${i}-day">Days (Ex: M, T, W, Th, F)</label>
            </div>
            <div class="form-outline">
                <input required  type="text" id="section-${i}-timeslot" name="section-${i}-timeslot" class="form-control" />
                <label class="form-label" for="section-${i}-timeslot">Timeslot (Ex: 12pm - 1pm)</label>
            </div>
        </div>
        `;
        scheduleInputsContainer.appendChild(scheduleInputGroup);
    }
    // Initialize the newly added MDB input elements
    initializeMDBInputs();
}

// On click of add course button run function addCourse
export function addCourse(event) {
    event.preventDefault();

    // Get the values of the form and trim them
    let courseNumber = document.getElementById("add-course-number").value.trim();
    let courseTitle = document.getElementById("add-course-title").value.trim();
    let courseSemester = document.getElementById("add-course-semester").value.trim();
    let courseYear = document.getElementById("add-course-year").value.trim();
    let courseSections = document.getElementById("course-sections").value.trim();
    let courseStartDate = document.getElementById("add-course-start").value;
    let courseEndDate = document.getElementById("add-course-end").value;

    // Create a schedules array
    let schedules = [];
    for (let i = 1; i <= courseSections; i++) {
        let day = document.getElementById(`section-${i}-day`).value.trim();
        let timeslot = document.getElementById(`section-${i}-timeslot`).value.trim();
        schedules.push(`${day} | ${timeslot}`);
    }

    // Create a new course object
    let newCourse = {
        courseNumber: courseNumber,
        courseTitle: courseTitle,
        courseSemester: courseSemester,
        courseStartDate: courseStartDate,
        courseEndDate: courseEndDate,
        courseYear: courseYear,
        courseSections: courseSections,
        schedules: schedules
    };

    // Check if all schedual fields are filled out 
    for (let i = 1; i <= courseSections; i++) {
        let day = document.getElementById(`section-${i}-day`).value;
        let timeslot = document.getElementById(`section-${i}-timeslot`).value;
        if (day == "" || timeslot == "") {
            document.getElementById("addCourseFormFillError").style.visibility = "visible";
            setTimeout(function () {
                document.getElementById("addCourseFormFillError").style.visibility = "hidden";
            }, 3000);

            return;
        }
    }

    // Check if all fields are filled out
    if (courseNumber == "" || courseTitle == "" || courseSemester == "" || courseYear == "" || courseSections == "" || courseStartDate == "" || courseEndDate == "") {
        document.getElementById("addCourseFormFillError").style.visibility = "visible";
        setTimeout(function () {
            document.getElementById("addCourseFormFillError").style.visibility = "hidden";
        }, 3000);
    } else {
        // Send the new course object to the server
        fetch('/api/course', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(newCourse)
        }).then(response => {
            if (response.status == 200) {
                // If the course was added successfully, reload the page
                location.reload();
            } else {
                // If the course was not added successfully, display an error message
                alert("Error adding course");
                // document.getElementById("addCourseError").style.display = "block";
            }
        })
    }
}

// openDeleteSectionModal() is called when the "Delete" button in the Delete User Modal is clicked
export function openDeleteSectionModal(button) {
    const courseId = button.getAttribute("data-course-id");
    const sectionNumber = button.getAttribute("data-section-number");
    const courseTitle = button.getAttribute("data-course-title");

    // Set user data in the Delete User Modal
    const courseTitleElement = document.querySelector("#delete-modal-title");
    const sectionNumberElement = document.querySelector("#delete-modal-section-number");

    courseTitleElement.textContent = `Course: ${courseTitle}`;
    sectionNumberElement.textContent = `Section: ${sectionNumber}`;

    // Add course ID to the "Delete" button (or another appropriate element) in the Delete User Modal
    document.getElementById("delete-course-section-modal").setAttribute("data-course-id", courseId);
    document.getElementById("delete-course-section-modal").setAttribute("data-section-number", sectionNumber);
}

// deleteCourseSection is called when the "Delete" button in the Delete Course Modal is clicked
export function deleteCourseSection() {
    const courseId = document.getElementById("delete-course-section-modal").getAttribute("data-course-id");
    const sectionNumber = document.getElementById("delete-course-section-modal").getAttribute("data-section-number");

    // Send a DELETE request to the server to /api/course/section
    fetch(`/api/course/section/${courseId}/${sectionNumber}`, {
        method: 'DELETE'
    }).then(response => {
        if (response.status == 200) {
            // If the user was deleted successfully, reload the page
            location.reload();
        } else {
            // If the user was not deleted successfully, display an error message
            document.getElementById("deleteCourseSectionError").style.display = "block";
        }
    })
}

// openEditCourseModal() is called when the "Edit" button in the Edit Course Modal is clicked
export function openEditCourseModal(button) {
    const courseID = button.getAttribute("data-course-id");
    const sectionID = button.getAttribute("data-section-id");
    const courseNumber = button.getAttribute("data-course-number");
    const courseTitle = button.getAttribute("data-course-title");
    const semester = button.getAttribute("data-course-semester");
    const year = button.getAttribute("data-course-year");
    const sectionName = button.getAttribute("data-section-number");
    const courseStartDate = button.getAttribute("data-course-start");
    const courseEndDate = button.getAttribute("data-course-end");
    const schedule = button.getAttribute("data-schedule");

    // split the schedule string into two strings (day and timeslot) split by the pipe character
    let scheduleDay = schedule.split(" | ")[0];
    let scheduleTime = schedule.split(" | ")[1];

    // Set user data in the Edit User Modal for better user experience
    document.getElementById("edit-modal-course-number").textContent = courseNumber;
    document.getElementById("edit-modal-course-number-t2").textContent = courseNumber;
    document.getElementById("edit-modal-course-title").textContent = courseTitle;
    document.getElementById("edit-modal-section").textContent = `Semester: ${semester}`;
    document.getElementById("edit-course-number").value = courseNumber;
    document.getElementById("edit-course-title").value = courseTitle;
    document.getElementById("edit-course-semester").value = semester;
    document.getElementById("edit-course-year").value = year;
    document.getElementById("edit-section-number").value = sectionName;
    document.getElementById("edit-section-number-t2").value = sectionName;
    document.getElementById("edit-section-number-t2-p").textContent = `Section: ${sectionName}`;
    document.getElementById("edit-course-start").value = courseStartDate;
    document.getElementById("edit-course-end").value = courseEndDate;
    document.getElementById("edit-schedule-days").value = scheduleDay;
    document.getElementById("edit-schedule-days-t2").value = scheduleDay;
    document.getElementById("edit-schedule-time").value = scheduleTime;
    document.getElementById("edit-schedule-time-t2").value = scheduleTime;
    document.getElementById("edit-course-and-section-btn").setAttribute("data-course-id", courseID);
    document.getElementById("edit-course-and-section-btn").setAttribute("data-section-id", sectionID);
}

// editCourseAndSection() is called when the button in the Edit Modal is clicked
export function editCourseAndSection(e) {
    e.preventDefault();

    let data;
    const courseID = parseInt(document.getElementById("edit-course-and-section-btn").getAttribute("data-course-id"), 10);
    const sectionID = parseInt(document.getElementById("edit-course-and-section-btn").getAttribute("data-section-id"), 10);
    const courseNumber = document.getElementById("edit-course-number").value.trim();
    const courseTitle = document.getElementById("edit-course-title").value.trim();
    const semester = document.getElementById("edit-course-semester").value.trim();
    const year = document.getElementById("edit-course-year").value.trim();
    const courseStartDate = document.getElementById("edit-course-start").value.trim();
    const courseEndDate = document.getElementById("edit-course-end").value.trim();
    
    const editTab1 = document.getElementById("edit-tab-1")

    if (editTab1.classList.contains("active")) {
        const sectionName = document.getElementById("edit-section-number").value.trim();
        const scheduleDays = document.getElementById("edit-schedule-days").value.trim();
        const scheduleTime = document.getElementById("edit-schedule-time").value.trim();

        // Create a JSON object
        data = {
            courseID,
            sectionID,
            courseNumber,
            courseTitle,
            semester,
            year,
            sectionName,
            courseStartDate,
            courseEndDate,
            scheduleDays,
            scheduleTime,
        };
    }
    else {
        const sectionName = document.getElementById("edit-section-number-t2").value.trim();
        const scheduleDays = document.getElementById("edit-schedule-days-t2").value.trim();
        const scheduleTime = document.getElementById("edit-schedule-time-t2").value.trim();

        // Create a JSON object
        data = {
            courseID,
            sectionID,
            courseNumber,
            courseTitle,
            semester,
            year,
            sectionName,
            courseStartDate,
            courseEndDate,
            scheduleDays,
            scheduleTime,
        };
    }
    // Send a PUT request to the server
    fetch(`/api/course/section`, {
        method: "PUT",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
    })
        .then((response) => {
            if (response.ok) {
                // If the response is OK, reload the page
                location.reload();
            } else {
                // Otherwise, display an error message
                response.json().then((errorData) => {
                    console.error('Error: ', errorData);
                });
            }
        });
}
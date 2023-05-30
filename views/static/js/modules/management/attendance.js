export function displaySections(button) {
    let courseID = button.dataset.courseid;  

    let instructorCourses = document.getElementById("instructor-courses");
    instructorCourses.classList.remove("show-courses");

    let sectionTable = document.getElementById("instructor-sections-table");
    sectionTable.style.display = "block";

    let sectionsTableBody = document.getElementById("sections-table-body");

    fetch(`/api/attendance/${courseID.trim()}`)
        .then(response => response.json())
        .then(data => {
            sectionsTableBody.innerHTML = "";

            // Display the sections in a table
            data.sections.forEach(element => {
                let row = `
                <tr onclick="displayStudentRecords(${element.AttendanceID}, event)" class="">
                    <td>${element.Date}</td>
                    <td>${element.CourseNumber}</td>
                    <td>${element.CourseTitle}</td>
                    <td>${element.SectionName}</td>
                </tr>
              `
                sectionsTableBody.innerHTML += row;
            });
        })
}


export function displayStudentRecords(attendanceID, e) {
    e.preventDefault();

    let studentRecordsList = document.getElementById("student-records-list-body");

    let studentRecords = document.getElementById("student-records");
    studentRecords.style.display = "block";

    let sectionTable = document.getElementById("instructor-sections-table");
    sectionTable.style.display = "none";

    fetch(`/api/attendance/students/${attendanceID}`)
        .then(response => response.json())
        .then(data => {
            studentRecordsList.innerHTML = "";

            // Display the sections in a table
            data.students.forEach(element => {

                let listItem = `
                    <li class="list-group-item d-flex justify-content-between align-items-center">
                        <div class="d-flex align-items-center">
                            <div class="ms-3">
                                <p class="fw-bold mb-1">${element.FirstName} ${element.LastName}</p>
                            </div>
                        </div>
                        ${element.Status == 'present' ? '<span class="badge rounded-pill badge-success">Present</span>'
                        :
                        '<span class="badge rounded-pill badge-danger">Absent</span>'}
                    </li>
              `
                studentRecordsList.innerHTML += listItem;
            });
        })
}

export function showInstructorCourses(e) {
    e.preventDefault();

    let instructorCourses = document.getElementById("instructor-courses");
    instructorCourses.classList.add("show-courses");

    let sectionTable = document.getElementById("instructor-sections-table");
    sectionTable.style.display = "none";
}

export function showSectionTable(e) {
    e.preventDefault();

    let studentRecords = document.getElementById("student-records");
    studentRecords.style.display = "none";

    let sectionTable = document.getElementById("instructor-sections-table");
    sectionTable.style.display = "block";
}
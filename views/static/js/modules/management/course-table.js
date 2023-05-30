
// Important global variable
let courseRowsToDisplay = 10;

// Update the rowSelect to include the current page
export function onCourseRowSelectChange() {
    const rowSelect = document.getElementById("rowSelect");
    const courseRowsToDisplay = parseInt(rowSelect.value, 10);
    localStorage.setItem("courseRowsToDisplay", courseRowsToDisplay); // Save the selected value to localStorage
    loadCourseManagementTable(courseRowsToDisplay, 1);
}

// Update the DOMContentLoaded event listener to include the current page
document.addEventListener("DOMContentLoaded", function () {
    const currentPage = window.location.pathname;

    if (currentPage === "/instructor") {
        const rowSelect = document.getElementById("rowSelect");
        // Restore the selected value from localStorage
        const savedRowsToDisplay = localStorage.getItem("courseRowsToDisplay");

        if (savedRowsToDisplay === null) {
            rowSelect.value = 10;
        }

        if (savedRowsToDisplay) {
            rowSelect.value = savedRowsToDisplay;
        }

        const courseRowsToDisplay = parseInt(rowSelect.value, 10);
        loadCourseManagementTable(courseRowsToDisplay, 1);
    }
});

// loadCourseManagementTable fetchs all courses then display them in the table that has an id of "tableBody"
export function loadCourseManagementTable(courseRowsToDisplay, currentPage = 1, courses = null) {
    if (!courses) {
        fetch("/api/course")
            .then((response) => response.json())
            .then((data) => {
                displayCourseTable(data.courseSections, courseRowsToDisplay, currentPage);
            });
    } else {
        displayCourseTable(courses, courseRowsToDisplay, currentPage);
    }
}

// Create a new function to handle the table rendering
export function displayCourseTable(courses, courseRowsToDisplay, currentPage) {
    const tableBody = document.getElementById("courseTableBody");

    // Clear the table body before appending new rows
    tableBody.innerHTML = "";

    if (courses?.length) {
        // Create the pagination
        createCoursesPagination(courses.length, courseRowsToDisplay, currentPage);

        const startIndex = (currentPage - 1) * courseRowsToDisplay;
        const endIndex = Math.min(currentPage * courseRowsToDisplay, courses.length);

        for (let i = startIndex; i < endIndex; i++) {

            const course = courses[i];
            let row = document.createElement("tr");
            row.innerHTML = `
                            <td> ${course.number}</td>
                            <td> ${course.title}</td>
                            <td> ${course.semester} ${course.year}</td>
                            <td> ${course.startDate} - ${course.endDate}</td>
                            <td>Section ${course.name}</td>
                            <td> ${course.schedule}</td>
                            <td> ${course.numStudents}</td>
                            <td>
                            <button 
                                class="editBtn table-btn"
                                onclick="openEditCourseModal(this)"
                                data-mdb-toggle="modal" 
                                data-mdb-target="#edit-course-modal"
                                data-course-id="${course.courseID}" data-section-id="${course.sectionID}"
                                data-course-number="${course.number}" data-course-title="${course.title}"
                                data-course-semester="${course.semester}" data-course-year="${course.year}"
                                data-section-number="${course.name}"
                                data-course-start="${course.startDate}" data-course-end="${course.endDate}"
                                data-schedule="${course.schedule}" 
                            >
                                <img src="/static/images/icon-edit.svg" alt="">
                            </button>
                            <button
                                class="deleteBtn table-btn" 
                                onclick="openDeleteSectionModal(this)"
                                data-mdb-target="#delete-course-modal"
                                data-course-id="${course.courseID}"
                                data-section-number="${course.name} "
                                data-mdb-toggle="modal" 
                                data-course-title="${course.title}"
                            >
                            <img src="/static/images/icon-trash.svg" alt="">
                            </button>
                            </td> `;
            tableBody.appendChild(row);
        };
    }
}

// Create the pagination
export function createCoursesPagination(totalCourses, courseRowsToDisplay, currentPage) {
    const pagination = document.getElementById("pagination");
    const totalPages = Math.ceil(totalCourses / courseRowsToDisplay);

    // Clear the existing pagination
    pagination.innerHTML = "";

    // Add the previous button
    pagination.innerHTML += `<li id="prevPageBtn" class="page-item"><a class="page-link" href="#">Previous</a></li>`;

    // Add the page buttons
    for (let i = 1; i <= totalPages; i++) {
        pagination.innerHTML += `<li class="page-item"><a class="page-link" href="#">${i}</a></li>`;
    }

    // Add the next button
    pagination.innerHTML += `<li id="nextPageBtn" class="page-item"><a class="page-link" href="#">Next</a></li>`;

    // Add event listeners for page buttons
    const pageButtons = document.querySelectorAll("#pagination .page-item");
    pageButtons.forEach((button, index) => {
        button.addEventListener("click", (event) => {
            event.preventDefault();
            if (index === 0) {
                // Previous button
                if (currentPage > 1) {
                    currentPage--;
                }
            } else if (index === pageButtons.length - 1) {
                // Next button
                if (currentPage < totalPages) {
                    currentPage++;
                }
            } else {
                // Page number button
                currentPage = parseInt(button.textContent, 10); // Get the page number from the button's textContent
            }
            loadCourseManagementTable(courseRowsToDisplay, currentPage);
        });
    });
}

// SearchUser searches for courses by name or email
export function searchCourse() {
    var input, filter;
    input = document.getElementById("searchCourseInput");
    filter = input.value.toUpperCase();

    fetch("/api/course")
        .then((response) => response.json())
        .then((data) => {
            const filteredData = data.courseSections.filter((course) => {
                const courseNumber = course.number;
                const courseTitle = course.title;

                return (
                    courseNumber.toUpperCase().indexOf(filter) > -1 ||
                    courseTitle.toUpperCase().indexOf(filter) > -1
                );
            });

            // Remove existing pagination
            const pagination = document.getElementById("pagination");
            pagination.innerHTML = "";

            // Display filtered data
            const rowSelect = document.getElementById("rowSelect");
            const courseRowsToDisplay = parseInt(rowSelect.value, 10);
            loadCourseManagementTable(courseRowsToDisplay, 1, filteredData);
        });
}
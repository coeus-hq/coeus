// Important global variable
let rowsToDisplay = 10;

// Update the rowSelect to include the current page
export function onRowSelectChange() {
    const rowSelect = document.getElementById("rowSelect");
    const rowsToDisplay = parseInt(rowSelect.value, 10);
    localStorage.setItem("rowsToDisplay", rowsToDisplay); // Save the selected value to localStorage
    loadUserManagementTable(rowsToDisplay, 1);
}

// Update the DOMContentLoaded event listener to include the current page
document.addEventListener("DOMContentLoaded", function () {
    const currentPage = window.location.pathname;

    if (currentPage === "/admin") {
        const rowSelect = document.getElementById("rowSelect");
        // Restore the selected value from localStorage
        const savedRowsToDisplay = localStorage.getItem("rowsToDisplay");

        if (savedRowsToDisplay === null) {
            rowSelect.value = 5;
        }

        if (savedRowsToDisplay) {
            rowSelect.value = savedRowsToDisplay;
        }

        const rowsToDisplay = parseInt(rowSelect.value, 10);
        loadUserManagementTable(rowsToDisplay, 1);
    }
});

// loadUserManagementTable fetchs all users then display them in the table that has an id of "tableBody"
export function loadUserManagementTable(rowsToDisplay, currentPage = 1, users = null) {
    if (!users) {
        fetch("/api/user")
            .then((response) => response.json())
            .then((data) => {
                displayUserTable(data.users, rowsToDisplay, currentPage);
            });
    } else {
        displayUserTable(users, rowsToDisplay, currentPage);
    }
}

// displayUserTable displays the users in the table that has an id of "tableBody"
export function displayUserTable(users, rowsToDisplay, currentPage) {
    const tableBody = document.getElementById("tableBody");

    // Clear the table body before appending new rows
    if (tableBody) {
        tableBody.innerHTML = "";

        // Create the pagination
        createPagination(users.length, rowsToDisplay, currentPage);

        const start = (currentPage - 1) * rowsToDisplay;
        const end = currentPage * rowsToDisplay;

        users.slice(start, end).forEach((user) => {
            let row = document.createElement("tr");
            row.innerHTML = `
        <td>${user.FirstName} ${user.LastName}</td>
        <td>${user.Email}</td>
        <td>${user.HighestModeratorType}</td>
        <td>
        <button class="editBtn table-btn" data-mdb-toggle="modal" data-mdb-target="#edit-user-modal"
        data-user-id="${user.ID}" data-user-firstname="${user.FirstName}"
        data-user-lastname="${user.LastName}" data-user-email="${user.Email}"
        onclick="openEditUserModal(this)">
        <img src="/static/images/icon-edit.svg" alt="">
        </button>
        <button class="deleteBtn table-btn" data-mdb-target="#delete-user-modal" data-mdb-toggle="modal"
        data-user-id="${user.ID}" data-user-firstname="${user.FirstName}"
        data-user-lastname="${user.LastName}" data-user-email="${user.Email}"
        onclick="openDeleteUserModal(this)">
        <img src="/static/images/icon-trash.svg" alt="">
        </button>       
        </td>
      `;

            tableBody.appendChild(row);

        });
    }
}

// createPagination creates the pagination buttons
export function createPagination(totalUsers, rowsToDisplay, currentPage) {
    const pagination = document.getElementById("pagination");
    const totalPages = Math.ceil(totalUsers / rowsToDisplay);

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
            loadUserManagementTable(rowsToDisplay, currentPage);
        });
    });
}

// SearchUser searches for users by name or email
export function searchUser() {
    // Declare variables
    var input, filter;
    input = document.getElementById("searchInput");
    filter = input.value.toUpperCase();

    fetch("/api/user")
        .then((response) => response.json())
        .then((data) => {
            const filteredData = data.users.filter((user) => {
                const fullName = `${user.FirstName} ${user.LastName}`;
                const email = user.Email;
                return (
                    fullName.toUpperCase().indexOf(filter) > -1 ||
                    email.toUpperCase().indexOf(filter) > -1
                );
            });

            // Remove existing pagination
            const pagination = document.getElementById("pagination");
            pagination.innerHTML = "";

            // Display filtered data
            const rowSelect = document.getElementById("rowSelect");
            const rowsToDisplay = parseInt(rowSelect.value, 10);
            loadUserManagementTable(rowsToDisplay, 1, filteredData);
        });
}
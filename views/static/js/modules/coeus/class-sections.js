// Get elements
const addSectionButton = document.getElementById("add-section-button");
const enrolledSections = document.querySelectorAll(".enrolled-section");
const sectionRow = document.querySelectorAll(".section-row");

// Handle add section button
enrolledSections.forEach((section) => {
  if (section.classList.contains("enrolled-section")) {
    addSectionButton.textContent = "Switch Section";
  }
});

// Show hidden enrolled section div if user is already enrolled in a section
if (sectionRow.length === 1 && sectionRow[0].classList.contains("enrolled-section")) {
  addSectionButton.classList.add("d-none");
  document.getElementById("hidden-enrolled-section-div").classList.remove("d-none");
}

// Validate section form
export function validateSectionForm() {
  const checked = document.querySelector('input[name="sectionSelect"]:checked');
  if (!checked) {
    document.getElementById("select-section-alert").style.display = "block";
    setTimeout(function () {
      document.getElementById("select-section-alert").style.display = "none";
    }, 3000);
    return;
  }
  document.querySelector('form').submit();
}
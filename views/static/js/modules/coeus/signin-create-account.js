// validateCreateAccountForm validates the create account form and fetches the create account route if the form is valid
export function validateCreateAccountForm(e) {

  e.preventDefault();

  let firstName = document.forms["createAccountForm"]["firstName"].value;
  let lastName = document.forms["createAccountForm"]["lastName"].value;
  let email = document.forms["createAccountForm"]["email"].value;
  let password = document.forms["createAccountForm"]["password"].value;
  let confirmPassword = document.forms["createAccountForm"]["confirmPassword"].value;
  let alert = document.getElementById("create-account-danger-alert");

  if (firstName == "" || lastName == "" || email == "" || password == "" || confirmPassword == "") {
    alert.innerHTML = "All fields must be filled out";
    alert.style.display = "block";
    setTimeout(function () {
      alert.style.display = "none";
    }, 3000);

    return false;
  }

  if (password != confirmPassword) {
    alert.innerHTML = "Passwords do not match";
    alert.style.display = "block";
    setTimeout(function () {
      alert.style.display = "none";
    }, 3000);
    return false;
  }

  let formData = new FormData(document.getElementById("createAccountForm"));

  fetch('/create-account', {
    method: 'POST',
    body: formData
  })
    .then(response => {

      if (response.status === 200) {
        window.location.href = "/";
      } else {
        alert.innerHTML = "Error creating account.";
        alert.style.display = "block";
        setTimeout(function () {
          alert.style.display = "none";
        }, 3000);
      }
    })
    .catch((error) => {
      console.error('Error:', error);
    });
}

// deleteCookies is called when the sign-in button is clicked if any existing cookies are found, delete them
export function deleteCookies() {
  var cookies = document.cookie.split(";");
  for (var i = 0; i < cookies.length; i++) {
    var cookie = cookies[i];
    var eqPos = cookie.indexOf("=");
    var name = eqPos > -1 ? cookie.substr(0, eqPos) : cookie;
    document.cookie = name + "=;expires=Thu, 01 Jan 1970 00:00:00 GMT";
  }
}
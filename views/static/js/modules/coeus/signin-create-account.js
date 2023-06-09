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
        return response.json();
      }
    })
    .then(data => {
      if (data) {
        alert.innerHTML = data.message || "Error creating account.";
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

export function signInForm(e) {
  e.preventDefault();

  let username = document.getElementById("emailAddress").value;
  let password = document.getElementById("password").value;

  if (username == "" || password == "") {
    let alert = document.getElementById("sign-in-fail-alert");
    alert.innerHTML = "All fields must be filled out";
    alert.classList.remove("hidden");

    setTimeout(function () {
      alert.classList.add("hidden");
    }, 3000);

    return false;
  }

  let formData = new FormData();
  formData.append('username', username);
  formData.append('password', password);

  fetch('/sign-in', {
    method: 'POST',
    body: formData,
})
    .then(response => response.json())
    .then(data => {

        if (data.success) {
            // Redirect based on role
            switch (data.role) {
                case 'admin':
                    window.location.href = "/admin";
                    break;
                case 'instructor':
                    window.location.href = "/";
                    break;
                default:
                    window.location.href = "/"; 
            }
        } else {
            let alert = document.getElementById("sign-in-fail-alert");
            alert.innerHTML = data.content;
            alert.classList.remove("hidden");

            setTimeout(function () {
                alert.classList.add("hidden");
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
  
  return true;
}

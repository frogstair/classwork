var storage; // Two global variables for the storage variables
var code; // and the code if a password is unset

// A shorthand for window.onload in jquery
$(() => {
    storage = window.sessionStorage; // set the storage variable
    code = storage.getItem("_cde") // get the code from storage
    if(code !== "_") { // check if the code already exists, an underscore means the user already has a password
        $("#password-label").text("Create a new password for your account")
    } else { // No underscore means the account doesnt have a password
        code = undefined
    }
})

function checkEmail(el) {
    data = { // Check the email to get a code or confirm that a password is already set
        email: $("#email-input").val()
    }

    // Make a request to check
    axios
    .post("/api/login/pass", data)
    .then((resp) => { // Set all the appropriate storage slots if success
        var code = resp.data.data
        storage.setItem("_eml", data.email)
        storage.setItem("_cde", code)
        window.location.replace("/login/pass")
    })
    .catch((err) => { // Check for error
        var status = err.response.status
        if(status == 410) { // The server returns 410 GONE if the password is already set
            storage.setItem("_eml", data.email)
            storage.setItem("_cde", "_") // Set the code as _
            window.location.replace("/login/pass")
        } else { // If some other error occurred then notify user
            $(el).addClass("border-danger")
            $("#errors").val(err.response.data.error)
        }
    })
}

function login() {
    // Disable button and make a spinner
    $("#login-button").addClass("disabled")
    $("#login-button").prepend($("<span id=\"waiter\" class=\"spinner-grow spinner-grow-sm\" role=\"status\" aria-hidden=\"true\"></span>"))

    // all the data the user sends
    var data = {
        email: storage.getItem("_eml"),
        password: $("#pass-input").val(),
        code: code // If undefined it wont be included
    }

    // Make a request
    axios
    .post("/api/login", data)
    .then(() => { // On success remove the code from the user for security reasons
        storage.removeItem("_eml")
        storage.removeItem("_cde")
        window.location.replace("/dashboard") // Move to dashboard
    })
    .catch((err) => { // Notify of error
        $("#pass-input").addClass("border-danger")
        $("#errors").val(err.response.data.error)
        $("#waiter").remove()
        $("#login-button").removeClass("disabled")
    })
}
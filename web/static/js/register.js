function validateName(element) {
    el = $(element) // Perform the same check as the server
    var name = el.val() // but in real time
    setValid(name.length > 3, "Name must be at least 4 characters long", el)
}

function validateEmail(element) {
    el = $(element) // Get the value of the input field

    var params = { // Prepare the data to be sent to the server
        email: el.val()
    }

    axios // Use axios to form a GET request
    .get("/api/register/email", {
        params: params // Set the parameters of the request to the above value
    })
    .then((res) => { // After a 2xx response was received
        var valid = res.data.data
        setValid(valid, "Email is taken", el)
    })
    .catch((err) => { // If a non-2xx response was received 
        var error = err.response.data.error
        setValid(false, error, el)
    })
}

function setValid(valid, text, ...elements) { // Enable or disable the button
    if(!valid) { // Code to disable the button
        elements.forEach((el) => { // Mark all the elements with an error
            el.addClass("border-danger")
        })
        $("#errors").val(text)
        $("#reg-button").addClass("disabled")
    } else { // Code to enable the button
        elements.forEach((el) => { // Unark all the elements with an error
            el.removeClass("border-danger")
        })
        $("#reg-button").removeClass("disabled")
        $("#errors").val("")
    }
    $("#waiter").remove()
}

function validatePassword() { // Check if the two passwords match
    var pass = $("#pass-input")
    var rpass = $("#pass-repeat")
    setValid(pass.val() === rpass.val(), "Passwords do not match", pass, rpass)
}

function register() {
    // Create a spinner inside the button element and disable it while waiting
    $("#reg-button").addClass("disabled")
    $("#reg-button").prepend($("<span id=\"waiter\" class=\"spinner-grow spinner-grow-sm\" role=\"status\" aria-hidden=\"true\"></span>"))

    // The data being sent to the server
    var data = {
        first_name: $("#fname").val(),
        last_name: $("#lname").val(),
        email: $("#email-input").val(),
        password: $("#pass-input").val(),
    }

    // Make a request
    axios
    .post("/api/register", data) // Set above data as the data for the post request
    .then((res) => { // If no error then redirect to the login page
        window.location.replace("/login")
    })
    .catch((err) => { // If an error occured then set all error fields
        setValid(false, err.response.data.error)
    })
}
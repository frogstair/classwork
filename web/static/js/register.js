function validateName(element) {
    el = $(element)
    var name = el.val()
    setValid(name.length > 3, "Name must be at least 4 characters long", el)
}

function validateEmail(element) {
    el = $(element)
    axios
    .get(`/api/register/email?email=${el.val()}`)
    .then((res) => {
        var valid = res.data.data
        setValid(valid, "Email is taken", el)
    })
    .catch((err) => {
        var error = err.response.data.error
        setValid(false, error, el)
    })
}

function setValid(valid, text, ...elements) {
    if(!valid) {
        elements.forEach((el) => {
            el.addClass("border-danger")
        })
        $("#errors").val(text)
        $("#reg-button").addClass("disabled")
    } else {
        elements.forEach((el) => {
            el.removeClass("border-danger")
        })
        $("#reg-button").removeClass("disabled")
        $("#errors").val("")
    }
    $("#waiter").remove()
}

function validatePassword() {
    var pass = $("#pass-input")
    var rpass = $("#pass-repeat")
    setValid(pass.val() === rpass.val(), "Passwords do not match", pass, rpass)
}

function register() {
    $("#reg-button").addClass("disabled")
    $("#reg-button").prepend($("<span id=\"waiter\" class=\"spinner-grow spinner-grow-sm\" role=\"status\" aria-hidden=\"true\"></span>"))

    var data = {
        first_name: $("#fname").val(),
        last_name: $("#lname").val(),
        email: $("#email-input").val(),
        password: $("#pass-input").val(),
    }

    axios
    .post("/api/register", data)
    .then((res) => {
        window.location.replace("/login")
    })
    .catch((err) => {
        setValid(false, err.response.data.error)
    })
}
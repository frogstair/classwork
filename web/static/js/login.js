var storage;
var code;

$(() => {
    storage = window.sessionStorage;
    code = storage.getItem("_cde")
    if(code !== "_") {
        $("#password-label").text("Create a new password for your account")
    } else {
        code = undefined
    }
})

function checkEmail(el) {
    data = {
        email: $("#email-input").val()
    }

    axios
    .post("/api/login/pass", data)
    .then((resp) => {
        var code = resp.data.data
        storage.setItem("_eml", data.email)
        storage.setItem("_cde", code)
        window.location.replace("/login/pass")
    })
    .catch((err) => {
        var status = err.response.status
        if(status == 410) {
            storage.setItem("_eml", data.email)
            storage.setItem("_cde", "_")
            window.location.replace("/login/pass")
        } else {
            el.addClass("border-danger")
            $("#errors").val(err.response.data.error)
        }
    })
}

function login() {
    $("#login-button").addClass("disabled")
    $("#login-button").prepend($("<span id=\"waiter\" class=\"spinner-grow spinner-grow-sm\" role=\"status\" aria-hidden=\"true\"></span>"))

    var data = {
        email: storage.getItem("_eml"),
        password: $("#pass-input").val(),
        code: code
    }

    axios
    .post("/api/login", data)
    .then((res) => {
        storage.removeItem("_eml")
        storage.removeItem("_cde")
        window.location.replace("/dashboard")
    })
    .catch((err) => {
        $("#pass-input").addClass("border-danger")
        $("#errors").val(err.response.data.error)
        $("#waiter").remove()
        $("#login-button").removeClass("disabled")
    })
}
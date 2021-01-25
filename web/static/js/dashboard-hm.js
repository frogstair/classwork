var dashboard;

const PERMS = {
  Headmaster: 1,
  Teacher: 2,
  Student: 3,
};

var selected = 0;

$(() => {
  $("#waiter").show();
  $("#content").hide();

  $("#content").on("complete", loadWorkspace);

  axios
    .get("/api/dashboard")
    .then((res) => {
      dashboard = res.data.data;
      try {
        $("#content").trigger("complete");
      } catch (err) {
        console.error(err);
      }
    })
    .catch((err) => {
      console.error(err);
      window.location.replace("/login");
    });
});

function loadWorkspace() {
  $("#waiter").remove();
  $("#content").show();
  selected = Cookies.get("_lsel");

  $("#name").text(dashboard.user.first_name);

  if (dashboard.headmaster) {
    var template = `<li class="nav-item">
            <b class="nav-link active">Headmaster</b>
        </li>`;
    $("#navbar").prepend($(template));
  }
  if (dashboard.teacher) {
    var template = `<li class="nav-item">
            <b class="nav-link active">Teacher</b>
        </li>`;
    $("#navbar").prepend($(template));
  }
  if (dashboard.student) {
    var template = `<li class="nav-item">
            <b class="nav-link active">Student</b>
        </li>`;
    $("#navbar").prepend($(template));
  }

  updateSelection();
}

function updateSelection() {
  if (!selected || selected == 0) {
    if (dashboard.headmaster) selected = PERMS.Headmaster;
    else if (dashboard.teacher) selected = PERMS.Teacher;
    else if (dashboard.student) selected = PERMS.Student;
    Cookies.set("_lsel", selected, { path: "/dashboard" });
  }

  if (selected == PERMS.Headmaster) {
    dashboard.headmaster.schools.forEach((school) => {
      template = schoolTemplate(school.id, school.name);
      $("#data").append($(template));
    });

    $("#data").append(
      $(`<div class="row">
        <div class="col-10">
          <input
            type="text"
            id="school-name"
            class="form-control mb-1"
            placeholder="School name"
          />

          <input
          id="errors"
          class="form-control"
          type="text"
          readonly
          />
        </div>
        <div class="col-2 d-grid">
          <button id="schooladder" onclick="addSchool()" class="btn btn-primary">
            <i id="plus" class="fas fa-plus"></i>
          </button>
        </div>
      </div>`)
    );
  }
}

function addSchool() {
  data = {
    name: $("#school-name").val(),
  };

  $("#schooladder").addClass("disabled");
  $("#schooladder").prepend(
    $(
      '<span id="waiter" class="spinner-grow spinner-grow-sm" role="status" aria-hidden="true"></span>'
    )
  );
  $("#plus").remove();

  axios
    .post("/api/school", data)
    .then((res) => {
      var school = res.data.data;
      template = schoolTemplate(school.id, school.name);
      $("#data").prepend($(template));
      $("#errors").val("");
      $("#school-name").val("");
      $("#schooladder").removeClass("disabled");
      $("#school-name").removeClass("border-danger");
      $("#schooladder").append(`<i id="plus" class="fas fa-plus"></i>`);
      $("#waiter").remove();
    })
    .catch((err) => {
      $("#school-name").addClass("border-danger");
      $("#errors").val(err.response.data.error);
      $("#schooladder").removeClass("disabled");
      $("#schooladder").append(`<i id="plus" class="fas fa-plus"></i>`);
      $("#waiter").remove();
    });
}

function logout() {
  axios
    .post("/api/logout")
    .then()
    .then(() => {
      window.location.replace("/login");
    });
}

function schoolTemplate(id, name) {
  return `<div id="${id}" class="card mb-3">
  <div class="row g-0">
    <div class="col-sm-8">
      <div class="card-body">
        <h5 class="card-title mt-4">${name}</h5>
      </div>
    </div>
    <div class="col-sm-4 mt-4 mb-4">
        <div class="row mb-2">
          <div class="col">
            <button onclick="manageSchool('${id}')" class="btn btn-primary">Manage</button>
          </div>
        </div>
        <div class="row">
          <div class="col">
              <button onclick="deleteSchool('${id}')" class="btn btn-primary">Delete</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>`;
}

function deleteSchool(id) {
  axios
    .delete("/api/school?id=" + encodeURI(id))
    .then(() => {
      $("#" + id).remove();
    })
    .catch((err) => {
      $("#errors").val(err.response.data.error);
    });
}

function manageSchool(id) {
  window.localStorage.setItem("_sch", id);
  window.location.href = "/school";
}

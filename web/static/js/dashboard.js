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

  window.localStorage.clear();

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
            <b onclick="selectRole(1)" class="nav-link clickable active">Headmaster</b>
        </li>`;
    $("#navbar").prepend($(template));
  }
  if (dashboard.teacher) {
    var template = `<li class="nav-item">
            <b onclick="selectRole(2)" class="nav-link clickable active">Teacher</b>
        </li>`;
    $("#navbar").prepend($(template));
  }
  if (dashboard.student) {
    var template = `<li class="nav-item">
            <b onclick="selectRole(3)" class="nav-link clickable active">Student</b>
        </li>`;
    $("#navbar").prepend($(template));
  }

  updateSelection();
}

function updateSelection() {
  if (!selected || selected == 0 || selected > 3) {
    if (dashboard.headmaster) selected = PERMS.Headmaster;
    else if (dashboard.teacher) selected = PERMS.Teacher;
    else if (dashboard.student) selected = PERMS.Student;
    Cookies.set("_lsel", selected, { path: "/dashboard" });
  }

  $("#data").empty();

  if (selected == PERMS.Headmaster) {
    $("#data").append($("<h1 id='title'>Headmaster</h1>"));

    dashboard.headmaster.schools.forEach((school) => {
      template = schoolTemplate(school.id, school.name);
      $("#data").append($(template));
      1;
    });

    $("#data").append(
      $(`<div class="row">
        <div class="col-10">
          <input
            type="text"
            class="form-control"
            id="school-name"
            placeholder="School name"
          />

          <input
          id="errors"
          type="text"
          class="form-control"
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
  if (selected == PERMS.Teacher) {
    $("#data").append($("<h1 id='title'>Teacher</h1>"));

    dashboard.teacher.subjects.forEach((subject) => {
      template = subjectTemplate(subject.id, subject.name);
      $("#data").append($(template));
      1;
    });

    $("#data").append(
      $(`<div class="row">
        <div class="col-10">
          <input
            type="text"
            id="subject-name"
            class="form-control mb-1"
            placeholder="Subject name"
          />

          <input
          id="errors"
          class="form-control"
          type="text"
          readonly
          />
        </div>
        <div class="col-2 d-grid">
          <button id="subjectadder" onclick="addSubject()" class="btn btn-primary">
            <i id="plus" class="fas fa-plus"></i>
          </button>
        </div>
      </div>`)
    );
  }
  if (selected == PERMS.Student) {
    $("#data").append($("<h1 id='title'>Student</h1>"));
    dashboard.student.subjects.forEach((subject) => {
      $("#data").append($(subjTemplate(subject.id, subject.name)));
      subject.assignments.forEach((assgn) => {
        var remaining = 0;
        assgn.requests.forEach((req) => {
          if (!req.complete) remaining++;
        });
        $("#" + subject.id).append(
          assgnTemplate(assgn.id, assgn.name, assgn.text, subject.id, remaining)
        );
      });
    });
  }
}

function logout() {
  axios
    .post("/api/logout")
    .then()
    .then(() => {
      window.location.replace("/login");
    });
}

function selectRole(role) {
  selected = role;
  Cookies.set("_lsel", role, { path: "/dashboard" });
  updateSelection();
}

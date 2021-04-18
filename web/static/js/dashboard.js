var dashboard;

const PERMS = { // A map of permissions and their corresponding number
  Headmaster: 1,
  Teacher: 2,
  Student: 3,
};

// Currently selected tab
var selected = 0;

// On window load
$(() => {
  $("#waiter").show(); // Show a waiter until all the information is retrieved from the server
  $("#content").hide();

  // Create an event when the content is loaded
  // When the event is triggered the loadWorkspace is run
  $("#content").on("complete", loadWorkspace);

  // Clear everything in the local storage
  window.localStorage.clear();

  // Get the information for the dashboard
  axios
    .get("/api/dashboard")
    .then((res) => {
      dashboard = res.data.data;
      try {
        // Show the content
        $("#content").trigger("complete");
      } catch (err) {
        // If an error occurs in the loadWorkspace function
        // it is caught here
        console.error(err);
      }
    })
    // If an error occurred then send the user to the login page
    .catch((err) => {
      console.error(err);
      window.location.replace("/login");
    });
});

function loadWorkspace() {
  // Remove all waiters and show the content
  $("#waiter").remove();
  $("#content").show();
  // Get the last selected tab from a cookie
  selected = Cookies.get("_lsel");

  // Get the user's name
  $("#name").text(dashboard.user.first_name);

  // Set the currently selected role and set it as active
  // in the navbar
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

  // Update selection will load in all the necessary info for the selected role
  updateSelection();
}

function updateSelection() {

  // Set the cookie for the last selected role
  if (!selected || selected == 0 || selected > 3) {
    if (dashboard.headmaster) selected = PERMS.Headmaster;
    else if (dashboard.teacher) selected = PERMS.Teacher;
    else if (dashboard.student) selected = PERMS.Student;
    Cookies.set("_lsel", selected, { path: "/dashboard" });
  }

  // Clear the screen
  $("#data").empty();


  // Insert the correct data depending on selected role
  if (selected == PERMS.Headmaster) {
    // Title
    $("#data").append($("<h1 id='title'>Headmaster</h1>"));

    // insert all the schools the headmaster owns
    dashboard.headmaster.schools.forEach((school) => {
      template = schoolTemplate(school.id, school.name);
      $("#data").append($(template));
      1;
    });

    // Append the form to add a new school
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

    // Insert all the subjects
    dashboard.teacher.subjects.forEach((subject) => {
      template = subjectTemplate(subject.id, subject.name);
      $("#data").append($(template));
      1;
    });

    // Insert the form to add a new subject
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
    // Title
    $("#data").append($("<h1 id='title'>Student</h1>"));

    // For each subject
    dashboard.student.subjects.forEach((subject) => {
      $("#data").append($(subjTemplate(subject.id, subject.name)));
      // For each assignment in the subject
      subject.assignments.forEach((assgn) => {
        // Insert all the requests and count the ones that are missing
        var remaining = 0;
        assgn.requests.forEach((req) => {
          if (!req.complete) remaining++;
        });
        $("#" + subject.id).append(
          // Append them all to the subject
          assgnTemplate(assgn.id, assgn.name, assgn.text, subject.id, remaining)
        );
      });
    });
  }
}

function logout() {
  // send a request to the server
  axios
    .post("/api/logout")
    .then() // empty then because it will only run in case of a success
    .then(() => { // Second then always runs
      window.location.replace("/login");
    });
}

function selectRole(role) {
  selected = role; // set the role
  Cookies.set("_lsel", role, { path: "/dashboard" }); // update the cookie
  updateSelection(); // update the selection
}

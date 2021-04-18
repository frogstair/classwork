var subject; // Global variables
var students;
var subject_id = window.localStorage.getItem("_sbj");
var school_id = window.localStorage.getItem("_sch");
var modal;

var upload_req = { // Data for an upload request
  name: "",
  text: "",
  subject_id: "",
  time_due: "",
  files: [],
  uploads: [],
};

$(() => {
  // Set the modal variable (the window that allows the user to create a new assignment)
  modal = new bootstrap.Modal(document.getElementById("assignmentModal"));
  // Show all the waiters
  $("#waiter").show();
  $("#content").hide();

  // Create an event
  $("#content").on("complete", loadWorkspace);

  window.localStorage.removeItem("_asn");

  // Params to get all assignments
  var params = {
    id: subject_id,
    sid: school_id,
  };

  // Make the request
  axios
    .get("/api/school/subject/", {
      params: params,
    })
    .then((res) => {
      subject = res.data.data.subject;
      students = res.data.data.students;

      try {
        $("#content").trigger("complete");
      } catch (err) {
        console.error(err);
      }
    })
    .catch((err) => {
      console.error(err);
      console.error(err.response.data.error);
      window.location.replace("/dashboard");
    });
});

// Template for the assignment
function asnTemplate(id, name, date, needsUpload) {
  var text = needsUpload
    ? `<p>Not completed by</p><div class="mb-3" id="${id}_nc"></div>`
    : "";

  return `<div class="card mb-3" id="${id}">
  <div class="card-body">
    <h5 class="card-title">${name}</h5>
    <p class="card-text"><small class="text-muted">Assigned ${date} </small></p>

    ${text}

    <a onclick="viewAssignment('${id}')" class="btn btn-primary">View</a>
  </div>
</div>`;
}

function loadWorkspace() {
  $("#waiter").remove(); // Show the content
  $("#content").show();
  $("#name").text(subject.name); // Set the name of the subject

  // Place in a try because random errors occur sometimes
  try {
    // Add each student in the subject to a list
    subject.students.forEach((student) => {
      $("#students").append(
        `<li value="${student.id}" class="list-group-item">${
          student.first_name + " " + student.last_name
        }</li>`
      );
    });
  } catch {}

  // Set the locale for the time
  var locale = window.navigator.userLanguage || window.navigator.language;
  moment.locale(locale);

  // If the subject has any assignments
  if (subject.assignments) {
    // Display each assignment with a template
    subject.assignments.forEach((assgn) => {
      // Get the time when each was assigned
      const date = moment(assgn.time_assigned).format("lll");
      // use a template
      var t = asnTemplate(assgn.id, assgn.name, date, assgn.not_completed_by);
      $("#assignments").append($(t));

      // If there are students who didnt complete an assignment
      if (assgn.not_completed_by != null) {
        assgn.not_completed_by.forEach((ncstd) => {
          $("#" + assgn.id + "_nc").append(
            `<span class="badge bg-danger">${
              ncstd.first_name + " " + ncstd.last_name
            }</span>`
          );
        });
      }
    });
  }

  // Add all students in the school to a selection list
  // to add them to the subject
  students.forEach((student) => {
    $("#students_select").append(
      $(
        `<option value="${student.id}">${
          student.first_name + " " + student.last_name
        }</option>`
      )
    );
  });
}

// Add a student to a subject
function addStudent(el) {
  $("#errors").val("");
  if (el.value == "0") return;

  data = { // All the necessary data
    user_id: el.value,
    subject_id: subject_id,
  };

  // Make the request
  axios
    .post("/api/school/subject/students", data)
    .then((res) => {
      student = res.data.data;
      // Add the student to the list
      $("#students").append(
        `<li value="${student.id}" class="list-group-item">${
          student.first_name + " " + student.last_name
        }</li>`
      );
    })
    .catch((err) => {
      $("#errors").val(err.response.data.error);
    });
}

// Add a new upload request
function addUploadRequest() {
  // Get the name of the request
  var name = $("#requestName").val();
  // Name cannot be empty
  if (name.length == 0) {
    return;
  }

  // Set all form elements to 0
  $("#requestName").val("");
  
  // Create a new list of elements
  var el = $("<li></li>");
  el.text(name);
  el.addClass("list-group-item");
  el.addClass("d-flex");
  el.addClass("justify-content-between");
  el.addClass("align-items-center");
  el.attr("id", upload_req.uploads.length);

  // Add the name of the request to the request list
  upload_req.uploads.push(name);

  // Add a delete button
  el.append(
    $(
      `<span onclick="removeUploadReq(this)" class="badge bg-danger rounded-pill">Delete</span>`
    )
  );
  // Add the request to the screen
  $("#uplRequests").append(el);
}

function removeUploadReq(el) {
  // Delete the request from the list
  var parent = $(el).parent();
  upload_req.uploads.splice(parent.attr("id"), 1);
  parent.remove();
}

function addAssignment() {
  // First upload the file
  var formData = new FormData();
  var files = $("#files").prop("files");
  $.each(files, (_, file) => {
    formData.append("files", file);
  });

  // Create a POST request to the files route
  axios
    .post("/files", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
    })
    .then((res) => {
      // After file was uploaded then set all necessary values
      // Path to files on server
      upload_req.files = res.data.files;
      // Set the needed values in the upload_req value
      upload_req.name = $("#asn_name").val();
      upload_req.subject_id = subject_id;
      upload_req.text = $("#text").val();

      // Set the time due into a valid format
      var time_due = moment.tz($("#timedue").val(), moment.tz.guess());

      upload_req.time_due = time_due;
      uploadAssgn(upload_req);
    })
    .catch((err) => {
      // Notify user of error
      console.error(err);
      $("#errors_asn").val(err.response.data.error);
    });
}


function uploadAssgn(data) {
  // Create a post request with all information from before
  axios
    .post("/api/school/subject/assignment", data)
    .then((res) => {
      var assgn = res.data.data;
      var time_created = moment.tz($("#timedue").val(), moment());
      $("#assignments").append(asnTemplate(assgn.id, assgn.name, time_created));
      modal.hide();
    })
    .catch((err) => {
      console.error(err);
      $("#errors_asn").val(err.response.data.error);
    });
}

// View the assignment
function viewAssignment(id) {
  window.localStorage.setItem("_asn", id);
  // Set role as teacher
  window.localStorage.setItem("_rol", "nqrOzz0jmxA=");
  // Redirect
  window.location.href = "/assignment";
}

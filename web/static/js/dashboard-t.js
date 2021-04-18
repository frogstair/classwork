function addSubject() {
  data = { // Get the data for the subject
    name: $("#subject-name").val(),
    school_id: dashboard.teacher.school_id, // Get the school ID from the teacher
  };

  // Disable the button
  $("#subjectadder").addClass("disabled");
  // Add a waiter
  $("#subjectadder").prepend(
    $(
      '<span id="waiter" class="spinner-grow spinner-grow-sm" role="status" aria-hidden="true"></span>'
    )
  );
  // Remove the plus icon
  $("#plus").remove();

  // Send a request to create a new subject
  axios
    .post("/api/school/subject", data)
    .then((res) => {
      var subject = res.data.data;
      // Add the subject at the top
      template = subjectTemplate(subject.id, subject.name);
      $("#title").after($(template));

      // Clear all errors and borders
      $("#errors").val("");
      $("#subject-name").val("");
      $("#subjectadder").removeClass("disabled");
      $("#subject-name").removeClass("border-danger");
    })
    .catch((err) => {
      // Set errors
      $("#subject-name").addClass("border-danger");
      $("#errors").val(err.response.data.error);
      $("#subjectadder").removeClass("disabled");
    })
    .then(() => {
      // remove the waiter and add the plus icon
      $("#subjectadder").append(`<i id="plus" class="fas fa-plus"></i>`);
      $("#waiter").remove();
    });
}

// A template for a subject card
function subjectTemplate(id, name) {
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
            <button onclick="manageSubject('${id}')" class="btn btn-primary">Manage</button>
          </div>
        </div>
        <div class="row">
          <div class="col">
              <button onclick="deleteSubject('${id}')" class="btn btn-primary">Delete</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>`;
}

function deleteSubject(id) {
  // Send a request to delete a subject
  axios
    .delete("/api/school/subject?id=" + encodeURI(id))
    .then(() => {
      $("#" + id).remove();
    })
    .catch((err) => {
      $("#errors").val(err.response.data.error);
    });
}

function manageSubject(id) {
  // Set the subject and school ID for future reference
  window.localStorage.setItem("_sbj", id)
  window.localStorage.setItem("_sch", dashboard.teacher.school_id)
  // Redirect user
  window.location.href = "/subject";
}
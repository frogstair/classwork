function addSchool() {
  data = { // Get the info to send to the server
    name: $("#school-name").val(),
  };

  // Disable the button while the request is running
  $("#schooladder").addClass("disabled");
  // Create a waiter
  $("#schooladder").prepend(
    $(
      '<span id="waiter" class="spinner-grow spinner-grow-sm" role="status" aria-hidden="true"></span>'
    )
  );
  // Remove the plus icon
  $("#plus").remove();

  // Send the request
  axios
    .post("/api/school", data)
    .then((res) => {
      var school = res.data.data;
      // Add the school at the top of the page
      template = schoolTemplate(school.id, school.name);
      $("#title").after($(template));
      // Remove all errors and borders
      $("#errors").val("");
      $("#school-name").val("");
      $("#schooladder").removeClass("disabled");
      $("#school-name").removeClass("border-danger");
    })
    .catch((err) => { // In case of error
      // Add a danger border and show error
      $("#school-name").addClass("border-danger");
      $("#errors").val(err.response.data.error);
      $("#schooladder").removeClass("disabled");
    })
    .then(() => { // Add the plus and remove the waiter
      $("#schooladder").append(`<i id="plus" class="fas fa-plus"></i>`);
      $("#waiter").remove();
    });
}

// School html template
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
  // Send a request to delete the school
  axios
    .delete("/api/school/subject?id=" + encodeURI(id))
    .then(() => {
      $("#" + id).remove();
    })
    .catch((err) => {
      $("#errors").val(err.response.data.error);
    });
}

function manageSchool(id) {
  // Set the storage item to the school ID for reference
  window.localStorage.setItem("_sch", id);
  // Move to the /school page
  window.location.href = "/school";
}

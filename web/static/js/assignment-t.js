var assgn_id = window.localStorage.getItem("_asn"); // Get all the global variables
var assignment;
var selected = 0;

$(() => {
  // Get the role of the person
  var role = window.localStorage.getItem("_rol");
  // If the role is not a teacher
  if (role != "nqrOzz0jmxA=") {
    // Remove the script from the page and let the other script run instead
    $("#t").remove();
    return;
  }

  // Show a waiter
  $("#waiter").show();
  $("#content").hide();

  $("#content").on("complete", loadWorkspace);

  var params = {
    id: assgn_id,
  };

  // Get all the info for the assignment
  axios
    .get("/api/school/subject/assignment", {
      params: params,
    })
    .then((res) => {
      assignment = res.data.data;
      try {
        $("#content").trigger("complete");
      } catch (e) {
        console.error(e);
      }
    })
    .catch((err) => {
      console.error(err);
    });
});

function loadWorkspace() {
  // Set the name of the assignment and the text
  $("#name").text(assignment.name);
  $("#text").text(assignment.text);

  // Parse when it was assigned and when its due
  const time_assigned = moment(assignment.time_assigned).format("lll");
  const time_due = moment(assignment.time_due).format("lll");

  // Add the text to the page
  $("#assigned").text("Assigned " + time_assigned);
  $("#due").text("Due " + time_due);

  // If the assignment has any files attached show that they can be downloaded
  if (assignment.files) {
    $("#files").append("<hr/><h4>Files</h4>")
    assignment.files.forEach((file, index) => {
      // Set the name of the download to
      // assignment_name_filename
      $("#files").append(`
      <a
        style="max-width: 18em"
        class="btn btn-primary"
        href="${file.path}"
        download="${assignment.name.replace(/\s/g, "_")}_${file.name}"
        >File ${index + 1}</a
      >`)
    })
  }

  // If there are any requests show them too
  if (assignment.requests) {
    var index = -1;
    assignment.requests.forEach((request) => {
      index++;
      var subText =
        request.uploads
          ? request.uploads.length + " people submitted"
          : "Nobody submitted";
      $("#requests").append(
        $(`
            <div onclick="select(${index})" class="card clickable">
              <div class="card-body">
                <h5 class="card-title">${request.name}</h5>
                <p>${subText}</p>
              </div>
            </div>`)
      );
    });
  }

  updateSelection();

  $("#waiter").remove();
  $("#content").show();
}

function updateSelection() {
  // Delete everything from the screen
  $("#uploads").empty();

  // Get the selected request
  var request = assignment.requests[selected];
  if (request.uploads) {
    // For each upload get the file to download
    request.uploads.forEach((upload) => {
      var downloadFile = `${
        upload.user.first_name + "_" + upload.user.last_name + "_" + upload.name
      }`;

      // template for each upload
      $("#uploads").append(
        $(`<div class="card">
      <div class="card-body">
        <h5 id="req-name" class="card-title"></h5>
        <div class="card mb-3">
          <div class="card-body">
            <div class="row g-0">
              <div class="col-lg-4">
                <h5 class="card-title">${
                  upload.user.first_name + " " + upload.user.last_name
                }</h5>
              </div>
              <div class="col-lg-8">
                <a class="btn btn-primary" href="${
                  upload.path
                }" download="${downloadFile}">Download user file</a>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>`)
      );
    });
  } else {
    // If nobody uploaded say "no submissions"
    $("#uploads").append(
      $(`<div class="card">
    <div class="card-body">
      <h5 id="req-name" class="card-title"></h5>
      <div class="card mb-3">
        <div class="card-body">
        <h2>No submissions</h2>
        </div>
      </div>
    </div>
  </div>`)
    );
  }
  $("#req-name").text(request.name);
}

// Select each request
function select(id) {
  selected = id;
  updateSelection();
}

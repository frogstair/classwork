var assgn_id = window.localStorage.getItem("_asn");
var assignment;
var selected = 0;

$(() => {
  var role = window.localStorage.getItem("_rol");
  if (role != "nqrOzz0jmxA=") {
    $("#t").remove();
    return;
  }

  $("#waiter").show();
  $("#content").hide();

  $("#content").on("complete", loadWorkspace);

  var params = {
    id: assgn_id,
  };

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
  $("#name").text(assignment.name);
  $("#text").text(assignment.text);

  const time_assigned = moment(assignment.time_assigned).format("lll");
  const time_due = moment(assignment.time_due).format("lll");

  $("#assigned").text("Assigned " + time_assigned);
  $("#due").text("Due " + time_due);

  if (assignment.files.length != 0) {
    $("#files").append("<hr/><h4>Files</h4>")
    assignment.files.forEach((file, index) => {
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

  if (assignment.requests) {
    var index = -1;
    assignment.requests.forEach((request) => {
      index++;
      var subText =
        request.uploads.length > 0
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
  $("#uploads").empty();

  var request = assignment.requests[selected];
  if (request.uploads.length != 0) {
    request.uploads.forEach((upload) => {
      var downloadFile = `${
        upload.user.first_name + upload.user.last_name + "_" + upload.filename
      }`;

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
                  request.file
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

function select(id) {
  selected = id;
  updateSelection();
}

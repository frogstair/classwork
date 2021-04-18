var assgn_id = window.localStorage.getItem("_asn");
var assignment;
var selected = 0;

$(() => {
  var role = window.localStorage.getItem("_rol");
  // If not a student is viewing the page then let the teacher script run
  if (role != "cNFDBScBBJ8=") {
    $("#st").remove();
    return;
  }

  // Show a waiter
  $("#waiter").show();
  $("#content").hide();

  $("#content").on("complete", loadWorkspaceSt);

  var params = {
    id: assgn_id,
  };

  // Get the information for the assignment
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

function loadWorkspaceSt() {
  // Show the name and text of assignment
  $("#name").text(assignment.name);
  $("#text").text(assignment.text);

  // Get when it was assigned and due
  const time_assigned = moment(assignment.time_assigned).format("lll");
  const time_due = moment(assignment.time_due).format("lll");

  // Add the text to the screen
  $("#assigned").text("Assigned " + time_assigned);
  $("#due").text("Due " + time_due);

  // For each file added to the assignment
  // Add a download button
  if (assignment.files) {
    $("#files").append("<hr/><h4>Files</h4>");
    assignment.files.forEach((file, index) => {
      $("#files").append(`
      <a
        style="max-width: 18em"
        class="btn btn-primary"
        href="${file.path}"
        download="${assignment.name.replace(/\s/g, "_")}_${file.name}"
        >File ${index + 1}</a
      >`);
    });
  }

  // If there are any upload requests
  if (assignment.requests) {
    // For each request show a button
    assignment.requests.forEach((request) => {
      if (request.complete) return;
      $("#uploads").append(
        $(`
            <div id="${request.id}" class="card mb-3">
              <div class="card-body">
                <div class="row">
                  <div class="col-lg-3">
                    <h5 class="card-title">${request.name}</h5>
                  </div>
                  <div class="col-lg-9">
                    <input onchange="completeReq('${request.id}', this)" id="file" class="form-control" type="file" id="uploadFile"/>
                  </div>
                </div>
              </div>
            </div>`)
      );
    });
  }

  $("#waiter").remove();
  $("#content").show();
}

function completeReq(id, el) {
  // Send the uploaded file to the server
  var formData = new FormData();
  var files = $(el).prop("files");
  $.each(files, (_, file) => {
    formData.append("files", file);
  });

  // Make the POST request
  axios
    .post("/files", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
    })
    .then((res) => {
      markComplete(id, res.data.files[0]);
    })
    .catch((err) => {
      console.error(err);
      $("#errors_asn").val(err.response.data.error);
    });
}

function markComplete(id, file) {
  // Mark assignment as complete
  data = {
    request_id: id,
    filepath: file,
  };

  // Create a POST request to complete an assignment
  axios.post("/api/school/subject/assignment/complete", data)
  .then((res) => {
    if (res.data.data) {
      // If uploaded successfuly then delete the upload request
      $("#" + id).remove()
    }
  })
  .catch((err) => {
    alert(err.response.data.error);
  });
}

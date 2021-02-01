var assgn_id = window.localStorage.getItem("_asn");
var assignment;
var selected = 0;

$(() => {
  var role = window.localStorage.getItem("_rol");
  if (role != "cNFDBScBBJ8=") {
    $("#st").remove();
    return;
  }

  $("#waiter").show();
  $("#content").hide();

  $("#content").on("complete", loadWorkspaceSt);

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

function loadWorkspaceSt() {
  $("#name").text(assignment.name);
  $("#text").text(assignment.text);

  const time_assigned = moment(assignment.time_assigned).format("lll");
  const time_due = moment(assignment.time_due).format("lll");

  $("#assigned").text("Assigned " + time_assigned);
  $("#due").text("Due " + time_due);

  if (assignment.files.length != 0) {
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

  if (assignment.requests) {
    var index = -1;
    assignment.requests.forEach((request) => {
      if (request.complete) return;
      index++;
      $("#uploads").append(
        $(`
            <div class="card mb-3">
              <div class="card-body">
                <div class="row">
                  <div class="col-lg-3">
                    <h5 class="card-title">${request.name}</h5>
                  </div>
                  <div class="col-lg-9">
                    <input onchange="completeReq('${request.id}')" id="file" class="form-control" type="file" id="uploadFile"/>
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

function completeReq(id) {
  var formData = new FormData();
  var files = $("#file").prop("files");
  $.each(files, (_, file) => {
    formData.append("files", file);
  });

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
  data = {
    request_id: id,
    filepath: file,
  };

  axios.post("/api/school/subject/assignment/complete", data).catch((err) => {
    alert(err.response.data.error);
  });
}

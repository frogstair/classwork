var subject;
var students;
var subject_id = window.localStorage.getItem("_sbj");
var school_id = window.localStorage.getItem("_sch");

var upload_req = {
  name: "",
  text: "",
  subject_id: "",
  time_due: "",
  files: [],
  uploads: [],
};

$(() => {
  $("#waiter").show();
  $("#content").hide();

  $("#content").on("complete", loadWorkspace);

  axios
    .get(
      "/api/school/subject/?id=" +
        encodeURIComponent(subject_id) +
        "&sid=" +
        encodeURIComponent(school_id)
    )
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

function loadWorkspace() {
  $("#name").text(subject.name);

  try {
    subject.students.forEach((student) => {
      $("#students").append(
        `<li value="${student.id}" class="list-group-item">${
          student.first_name + " " + student.last_name
        }</li>`
      );
    });
  } catch {}

  var locale = window.navigator.userLanguage || window.navigator.language;
  moment.locale(locale);

  try {
    subject.assignments.forEach((assgn) => {

      const date = moment(assgn.time_assigned).format("lll");

      $("#assignments").append(
        `<div class="card mb-3" id="${assgn.id}">
        <div class="card-body">
          <h5 class="card-title">${assgn.name}</h5>
          <p class="card-text"><small class="text-muted">Assigned ${date} </small></p>

          <p>Not completed by</p>
          <div class="mb-3" id="${assgn.id}_nc">
          </div>

          <a href="#" class="btn btn-primary">View</a>
        </div>
      </div>`
      );

      try {
        assgn.not_completed_by.forEach((ncstd) => {
          $("#" + assgn.id + "_nc").append(`<span class="badge bg-danger">${ncstd.first_name + " " + ncstd.last_name}</span>`)
        })
      } catch {}

    });
  } catch {}

  students.forEach((student) => {
    $("#students_select").append(
      $(
        `<option value="${student.id}">${
          student.first_name + " " + student.last_name
        }</option>`
      )
    );
  });

  $("#waiter").remove();
  $("#content").show();
}

function addStudent(el) {
  $("#errors").val("");
  if (el.value == "0") return;

  data = {
    user_id: el.value,
    subject_id: subject_id,
  };

  axios.post("/api/school/subject/students", data)
  .then((res) => {
    student = res.data.data
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

function addUploadRequest() {
  var name = $("#requestName").val();
  if (name.length == 0) {
    return;
  }
  $("#requestName").val("");
  var el = $("<li></li>");
  el.text(name);
  el.addClass("list-group-item");
  el.addClass("d-flex");
  el.addClass("justify-content-between");
  el.addClass("align-items-center");
  el.attr("id", upload_req.uploads.length);

  upload_req.uploads.push(name);

  el.append(
    $(
      `<span onclick="removeUploadReq(this)" class="badge bg-danger rounded-pill">Delete</span>`
    )
  );
  $("#uplRequests").append(el);
}

function removeUploadReq(el) {
  var parent = $(el).parent();
  upload_req.uploads.splice(parent.attr("id"), 1);
  parent.remove();
}

function addAssignment() {
  var formData = new FormData();
  var files = $("#files").prop("files");
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
      upload_req.files = res.data.files;
      upload_req.name = $("#asg_name").val();
      upload_req.subject_id = subject_id;
      upload_req.text = $("#text").val();

      moment.tz.load({
        zones: [],
        links: [],
        version: "2014e",
      });

      var time_due = moment.tz($("#timedue").val(), moment.tz.guess());

      upload_req.time_due = time_due;
      uploadAssgn(upload_req);
    })
    .catch((err) => {
      console.error(err);
      $("#errors_asg").val(err.response.data.error);
    });
}

function uploadAssgn(data) {
  axios.post("/api/school/subject/assignment", data).then((res) => {
    console.log(res);
  })
  .catch((res) => {

  });
}

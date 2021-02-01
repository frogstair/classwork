var subject;
var students;
var subject_id = window.localStorage.getItem("_sbj");
var school_id = window.localStorage.getItem("_sch");
var modal;

var upload_req = {
  name: "",
  text: "",
  subject_id: "",
  time_due: "",
  files: [],
  uploads: [],
};

$(() => {
  modal = new bootstrap.Modal(document.getElementById("assignmentModal"));
  $("#waiter").show();
  $("#content").hide();

  $("#content").on("complete", loadWorkspace);

  window.localStorage.removeItem("_asn");

  var params = {
    id: subject_id,
    sid: school_id,
  };

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
  $("#waiter").remove();
  $("#content").show();
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

  if (subject.assignments) {
    subject.assignments.forEach((assgn) => {
      const date = moment(assgn.time_assigned).format("lll");
      var t = asnTemplate(assgn.id, assgn.name, date, assgn.not_completed_by);
      $("#assignments").append($(t));

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

function addStudent(el) {
  $("#errors").val("");
  if (el.value == "0") return;

  data = {
    user_id: el.value,
    subject_id: subject_id,
  };

  axios
    .post("/api/school/subject/students", data)
    .then((res) => {
      student = res.data.data;
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
      upload_req.name = $("#asn_name").val();
      upload_req.subject_id = subject_id;
      upload_req.text = $("#text").val();

      var time_due = moment.tz($("#timedue").val(), moment.tz.guess());

      upload_req.time_due = time_due;
      uploadAssgn(upload_req);
    })
    .catch((err) => {
      console.error(err);
      $("#errors_asn").val(err.response.data.error);
    });
}

function uploadAssgn(data) {
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

function viewAssignment(id) {
  window.localStorage.setItem("_asn", id);
  window.localStorage.setItem("_rol", "nqrOzz0jmxA=");
  window.location.href = "/assignment";
}

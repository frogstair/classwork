function addSubject() {
  data = {
    name: $("#subject-name").val(),
    school_id: dashboard.teacher.school_id,
  };

  $("#subjectadder").addClass("disabled");
  $("#subjectadder").prepend(
    $(
      '<span id="waiter" class="spinner-grow spinner-grow-sm" role="status" aria-hidden="true"></span>'
    )
  );
  $("#plus").remove();

  axios
    .post("/api/school/subject", data)
    .then((res) => {
      var subject = res.data.data;
      template = subjectTemplate(subject.id, subject.name);
      $("#title").after($(template));
      $("#errors").val("");
      $("#subject-name").val("");
      $("#subjectadder").removeClass("disabled");
      $("#subject-name").removeClass("border-danger");
      $("#subjectadder").append(`<i id="plus" class="fas fa-plus"></i>`);
      $("#waiter").remove();
    })
    .catch((err) => {
      $("#subject-name").addClass("border-danger");
      $("#errors").val(err.response.data.error);
      $("#subjectadder").removeClass("disabled");
      $("#subjectadder").append(`<i id="plus" class="fas fa-plus"></i>`);
      $("#waiter").remove();
    });
}

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
            <button onclick="manageSchool('${id}')" class="btn btn-primary">Manage</button>
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
  axios
    .delete("/api/school/subject?id=" + encodeURI(id))
    .then(() => {
      $("#" + id).remove();
    })
    .catch((err) => {
      $("#errors").val(err.response.data.error);
    });
}
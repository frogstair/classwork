function addSchool() {
  data = {
    name: $("#school-name").val(),
  };

  $("#schooladder").addClass("disabled");
  $("#schooladder").prepend(
    $(
      '<span id="waiter" class="spinner-grow spinner-grow-sm" role="status" aria-hidden="true"></span>'
    )
  );
  $("#plus").remove();

  axios
    .post("/api/school", data)
    .then((res) => {
      var school = res.data.data;
      template = schoolTemplate(school.id, school.name);
      $("#title").after($(template));
      $("#errors").val("");
      $("#school-name").val("");
      $("#schooladder").removeClass("disabled");
      $("#school-name").removeClass("border-danger");
      $("#schooladder").append(`<i id="plus" class="fas fa-plus"></i>`);
      $("#waiter").remove();
    })
    .catch((err) => {
      $("#school-name").addClass("border-danger");
      $("#errors").val(err.response.data.error);
      $("#schooladder").removeClass("disabled");
      $("#schooladder").append(`<i id="plus" class="fas fa-plus"></i>`);
      $("#waiter").remove();
    });
}

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
  window.localStorage.setItem("_sch", id);
  window.location.href = "/school";
}
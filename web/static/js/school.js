var school_id = window.localStorage.getItem("_sch");

var data;

$(() => {
  if (!school_id) {
    window.location.href = "/dashboard";
  }

  $("#waiter").show();
  $("#content").hide();

  $("#content").on("complete", loadWorkspace);

  var data = {
    id: school_id
  }

  axios
    .get("/api/school", {
      params: data
    })
    .then((res) => {
      data = res.data.data;
      try {
        $("#content").trigger("complete");
      } catch (err) {
        console.error(err);
      }
    })
    .catch((err) => {
      console.error(err);
      window.localStorage.removeItem("_sch");
      window.location.replace("/dashboard");
    });
});

function loadWorkspace() {
  $("#school_name").text(data.name);

  if (!data.teachers) {
    $("#teacher-form").before($("<h5 id='teacher-heading'>No teachers!</h5>"));
  } else {
    $("#teacher-form").before(
      $(`<ul id="teacher-list" class="list-group mb-2"></ul>`)
    );
    data.teachers.forEach((t) => {
      $("#teacher-list").append(
        $(template("teacher", t.id, t.first_name + " " + t.last_name))
      );
    });
  }

  if (!data.students) {
    $("#student-form").before($("<h5 id='student-heading'>No students!</h5>"));
  } else {
    $("#student-form").before(
      $(`<ul id="student-list" class="list-group mb-2"></ul>`)
    );
    data.students.forEach((s) => {
      $("#student-list").append(
        $(template("student", s.id, s.first_name + " " + s.last_name))
      );
    });
  }

  $("#waiter").remove();
  $("#content").show();
}

function template(who, id, name) {
  t = `
    <li
        class="list-group-item d-flex justify-content-between align-items-center"
        id="${who + '_' + id}"
    >
        ${name}
        <span onclick="remove('${who}', '${id}', this)" class="badge bg-danger clickable rounded-pill">Delete</span>
    </li>`;
  return t;
}

function add(who) {
  $("#" + who + "-errors").val("");
  $("#" + who + "-add").addClass("disabled")
  data = {
    first_name: $("#" + who + "-fname").val(),
    last_name: $("#" + who + "-lname").val(),
    email: $("#" + who + "-email").val(),
    school_id: school_id,
  };

  if ($("#" + who + "-list").length == 0) {
    $("#" + who + "-heading").remove();
    $("#" + who + "-form").before(
      $(`<ul id="${who}-list" class="list-group mb-2"></ul>`)
    );
  }

  axios
    .post("/api/school/" + who, data)
    .then((res) => {
      t = template(
        who,
        res.data.data.id,
        res.data.data.first_name + " " + res.data.data.last_name
      );
      $("#" + who + "-list").append($(t));
      $("#" + who + "-fname").val("");
      $("#" + who + "-lname").val("");
      $("#" + who + "-email").val("");
      $("#" + who + "-add").removeClass("disabled")
    })
    .catch((err) => {
      $("#" + who + "-errors").val(err.response.data.error);
      $("#" + who + "-add").removeClass("disabled")
    });
}

function remove(who, id, button) {
  $("#" + who + "-errors").val("");
  el = $(button)
  el.removeClass("bg-danger")
  el.addClass("bg-secondary")
  axios
    .delete("/api/school/" + who + "?sid=" + school_id + "&uid=" + id)
    .then(() => {
      $("#" + who + "_" + id).remove();
      let amount = $("#" + who + "-list > li").length;
      if (amount == 0) {
        $("#" + who + "-form").before($(`<h5 id='${who}-heading'>No ${who}s!<h5>`));
      }
    })
    .catch((err) => {
      el.addClass("bg-danger")
      el.removeClass("bg-secondary")
      $("#" + who + "-errors").val(err.response.data.error);
    });
}

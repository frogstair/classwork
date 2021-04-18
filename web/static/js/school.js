var school_id = window.localStorage.getItem("_sch"); // Get the school ID

var schoolData;

// When the page loads
$(() => {
  // If the school ID doesnt exist then return to dashboard
  if (!school_id) {
    window.location.href = "/dashboard";
  }

  // Show a waiter
  $("#waiter").show();
  $("#content").hide();

  // Add an event when loading is complete
  $("#content").on("complete", loadWorkspace);

  // The data required for a GET request
  var data = {
    id: school_id
  }

  // Make the network request
  axios
    .get("/api/school", {
      params: data
    })
    .then((res) => {
      schoolData = res.data.data;
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

// Load the workspace
function loadWorkspace() {
  $("#school_name").text(schoolData.name);

  // If there are no teachers then show a message
  if (!schoolData.teachers) {
    $("#teacher-form").before($("<h5 id='teacher-heading'>No teachers!</h5>"));
  } else {
    // Otherwise list all the teachers
    $("#teacher-form").before(
      $(`<ul id="teacher-list" class="list-group mb-2"></ul>`)
    );
    schoolData.teachers.forEach((t) => {
      $("#teacher-list").append(
        $(template("teacher", t.id, t.first_name + " " + t.last_name))
      );
    });
  }

  // Same for students
  if (!schoolData.students) {
    $("#student-form").before($("<h5 id='student-heading'>No students!</h5>"));
  } else {
    $("#student-form").before(
      $(`<ul id="student-list" class="list-group mb-2"></ul>`)
    );
    schoolData.students.forEach((s) => {
      $("#student-list").append(
        $(template("student", s.id, s.first_name + " " + s.last_name))
      );
    });
  }

  // Show all the content
  $("#waiter").remove();
  $("#content").show();
}

function template(who, id, name) {
  // Template for one the list elements of students and teachers
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
  // who acts like selector of whether you want to add
  // a student or a teacher
  $("#" + who + "-errors").val("");
  $("#" + who + "-add").addClass("disabled")
  schoolData = {
    first_name: $("#" + who + "-fname").val(),
    last_name: $("#" + who + "-lname").val(),
    email: $("#" + who + "-email").val(),
    school_id: school_id,
  };

  // If there were no teachers in the beginning, then remove the header
  // that says there are no teachers
  if ($("#" + who + "-list").length == 0) {
    $("#" + who + "-heading").remove();
    $("#" + who + "-form").before(
      $(`<ul id="${who}-list" class="list-group mb-2"></ul>`)
    );
  }

  // Make the request
  axios
    .post("/api/school/" + who, schoolData)
    .then((res) => {
      // Add the new teacher to the list
      t = template(
        who,
        res.data.data.id,
        res.data.data.first_name + " " + res.data.data.last_name
      );
      $("#" + who + "-list").append($(t));
      // Clear the form and enable the button
      $("#" + who + "-fname").val("");
      $("#" + who + "-lname").val("");
      $("#" + who + "-email").val("");
      $("#" + who + "-add").removeClass("disabled")
    })
    .catch((err) => {
      // Show the error to the user
      $("#" + who + "-errors").val(err.response.data.error);
      $("#" + who + "-add").removeClass("disabled")
    });
}

function remove(who, id, button) {
  // Remove the student or the teacher
  $("#" + who + "-errors").val("");
  el = $(button)
  el.removeClass("bg-danger")
  el.addClass("bg-secondary")
  // Make a delete request
  axios
    .delete("/api/school/" + who + "?sid=" + school_id + "&uid=" + id)
    .then(() => {
      // remove the list element with the teacher
      $("#" + who + "_" + id).remove();
      let amount = $("#" + who + "-list > li").length;
      // Show if there is 0 elements of the list left
      if (amount == 0) {
        $("#" + who + "-form").before($(`<h5 id='${who}-heading'>No ${who}s!<h5>`));
      }
    })
    .catch((err) => {
      // Show all errors
      el.addClass("bg-danger")
      el.removeClass("bg-secondary")
      $("#" + who + "-errors").val(err.response.data.error);
    });
}

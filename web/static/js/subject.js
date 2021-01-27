var subject;
var students;
var subject_id = window.localStorage.getItem("_sbj");
var school_id = window.localStorage.getItem("_sch");

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
        `<li value="${student.id}" class="list-group-item">${student.first_name + " " + student.last_name}</li>`
      );
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
  if(el.value == "0") return;

  data = {
    user_id: el.value,
    subject_id: subject_id
  }

  axios.post("/api/school/subject/students", data)
  .then((res) => {
    console.log(res.data.data)
  })
}
// Template for a subject
function subjTemplate(subjid, subjname) {
  return `<div id="${subjid}">
    <h3 class="mt-3">${subjname}</h3>
    
  </div>`;
}

// Template for an assignment
function assgnTemplate(id, name, text, subjid, remaining) {
  return `<div class="card mb-3">
    <div class="card-body">
      <h4>${name}</h4>
      <p>${text}</p>
      <small class="text-muted">${remaining} parts remaining</small
      ><br /><br />
      <button onclick="viewAssgn('${id}', '${subjid}')" class="btn btn-primary">
        View
      </button>
    </div>
  </div>`;
}

// View assignment
function viewAssgn(id, sid) {
  // Set all appropritate fields for future reference
  window.localStorage.setItem("_sub", sid);
  window.localStorage.setItem("_asn", id);
  // Indicate that a student is viewing the assignment
  window.localStorage.setItem("_rol", "cNFDBScBBJ8=");
  // Redirect
  window.location.href = "/assignment";
}
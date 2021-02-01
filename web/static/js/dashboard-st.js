function subjTemplate(subjid, subjname) {
  return `<div id="${subjid}">
    <h3 class="mt-3">${subjname}</h3>
    
  </div>`;
}

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

function viewAssgn(id, sid) {
  window.localStorage.setItem("_sub", sid);
  window.localStorage.setItem("_asn", id);
  window.location.href = "/assignment";
}

//M88ObwDDNNY05MdyglPAcjYlxLXTb9_1
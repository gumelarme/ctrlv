function attachPostSelectorEvents() {
  const radios = document.querySelectorAll("#post-list > [type=radio]");
  for (const r of radios) {
    r.addEventListener("change", function() {
      const postId = r.getAttribute("value");
      window.location.href = "/p/" + postId;
    })
  }
}

// Clear form for new post
function clearForm() {
  const title = document.querySelector("#editor input[name=Title]");
  const id = document.querySelector("#editor input[name=Id]");
  const alias = document.querySelector("#editor input[name=Alias]");
  const category = document.querySelector("#editor input[name=Category]");
  const content = document.querySelector("#editor textarea");
  const visibilityOptions = document.querySelector("#editor select[name=Visibility]").options;

  id.value = "";
  alias.value = "";
  content.value = "";
  title.value = "";
  category.value = "note";
  for (const opt of visibilityOptions) {
    opt.selected = opt.value == "private";
  }
}


attachPostSelectorEvents()
const btnNew = document.querySelector("#btn-new");
btnNew.addEventListener("click", function() {
  clearForm()
})

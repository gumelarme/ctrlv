function attachPostSelectorEvents() {
  const radios = document.querySelectorAll("#post-list > [type=radio]");
  for (const r of radios) {
    r.addEventListener("change", function() {
      const postId = r.getAttribute("value");
      window.location.href = "/p/" + postId;
    })
  }
}

attachPostSelectorEvents()

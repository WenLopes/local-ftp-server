export async function fetchFiles() {
  const response = await fetch(`/files/?t=${new Date().getTime()}`);
  const text = await response.text();
  const parser = new DOMParser();
  const doc = parser.parseFromString(text, "text/html");

  return Array.from(doc.querySelectorAll("a"))
    .map(link => link.getAttribute("href"))
    .filter(href => href !== "../");
}

export async function uploadFile(file) {
  const formData = new FormData();
  formData.append("file", file);

  const response = await fetch("/upload", {
    method: "POST",
    body: formData
  });

  if (!response.ok) {
    throw new Error("Falha no upload");
  }
}

export async function deleteFile(filename) {
  const response = await fetch(`/delete/${filename}`, {
    method: "DELETE"
  });

  if (!response.ok) {
    throw new Error("Erro ao excluir arquivo");
  }
}

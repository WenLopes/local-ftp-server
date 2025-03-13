import { useEffect, useState } from "react";
import "./styles.css";
import { uploadFile, fetchFiles, deleteFile } from "./api";

function App() {
  const [files, setFiles] = useState([]);
  const [selectedFile, setSelectedFile] = useState(null);
  const [message, setMessage] = useState({ text: "", type: "" });

  const loadFiles = async () => {
    const fileList = await fetchFiles();
    setFiles([...fileList]);
  };

  useEffect(() => {
    loadFiles();
  }, []);

  const handleUpload = async (e) => {
    e.preventDefault();
    if (!selectedFile) {
      setMessage({ text: "Selecione um arquivo para enviar.", type: "error" });
      return;
    }

    try {
      await uploadFile(selectedFile);
      setMessage({ text: "Arquivo enviado com sucesso!", type: "success" });
      setSelectedFile(null);
      loadFiles();
    } catch (error) {
      setMessage({ text: "Erro ao enviar arquivo.", type: "error" });
    }
  };

  const handleDelete = async (filename) => {
    if (!window.confirm(`Tem certeza que deseja excluir "${filename}"?`)) return;

    try {
      await deleteFile(filename);
      setMessage({ text: "Arquivo excluído com sucesso!", type: "success" });
      loadFiles();
    } catch (error) {
      setMessage({ text: "Erro ao excluir arquivo.", type: "error" });
    }
  };

  return (
    <div className="container">
      <h2>Servidor de Arquivos</h2>

      {message.text && <div className={`message ${message.type}`}>{message.text}</div>}

      <div className="upload-box">
        <h3>Upload de Arquivo</h3>
        <form onSubmit={handleUpload}>
          <input type="file" onChange={(e) => setSelectedFile(e.target.files[0])} />
          <button type="submit">Enviar</button>
        </form>
      </div>

      <div className="file-list">
        <h3>Arquivos Disponíveis</h3>
        <ul>
          {files.length === 0 ? <li>Nenhum arquivo disponível.</li> : files.map((file, index) => (
            <li key={index}>
              <a href={`/files/${file}`} target="_blank" rel="noopener noreferrer">{file}</a>
              <button className="delete-btn" onClick={() => handleDelete(file)}>Excluir</button>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}

export default App;

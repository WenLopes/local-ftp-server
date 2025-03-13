package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

const (
	sharedDir  = "shared"
	websiteDir = "website/dist"
	port       = "8080"
)

func main() {
	if _, err := os.Stat(sharedDir); os.IsNotExist(err) {
		err := os.Mkdir(sharedDir, os.ModePerm)
		if err != nil {
			log.Fatalf("Erro ao criar diretório: %v", err)
		}
		fmt.Println("Diretório 'shared/' criado.")
	}

	r := mux.NewRouter()

	r.PathPrefix("/files/").Handler(http.StripPrefix("/files/", http.FileServer(http.Dir(sharedDir))))

	r.HandleFunc("/upload", uploadHandler).Methods("POST")

	r.HandleFunc("/delete/{filename}", deleteFileHandler).Methods("DELETE")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(websiteDir)))

	ip, err := getLocalIP()
	if err != nil {
		log.Fatalf("Erro ao obter IP local: %v", err)
	}

	fmt.Println("Servidor rodando!")
	fmt.Printf("Acesse no PC: http://localhost:%s/\n", port)
	fmt.Printf("Acesse na rede WiFi: http://%s:%s/\n", ip, port)

	log.Fatal(http.ListenAndServe(":"+port, r))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Arquivo muito grande", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Erro ao processar arquivo", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	filePath := filepath.Join(sharedDir, handler.Filename)
	outFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Erro ao salvar arquivo", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	written, err := io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Erro ao gravar arquivo", http.StatusInternalServerError)
		return
	}

	if written == 0 {
		http.Error(w, "Arquivo salvo com 0 bytes", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Arquivo deletado com sucesso"))
}

func deleteFileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	filePath := filepath.Join(sharedDir, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Arquivo não encontrado", http.StatusNotFound)
		return
	}

	err := os.Remove(filePath)
	if err != nil {
		http.Error(w, "Erro ao deletar arquivo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Arquivo deletado com sucesso"))
}

func getLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

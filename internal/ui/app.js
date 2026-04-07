const uploadInput = document.getElementById("fileInput");
const uploadMessage = document.getElementById("uploadMessage");
const fileList = document.getElementById("fileList");

// загрузка файла
async function uploadFile() {
    const file = uploadInput.files[0];
    if (!file) {
        uploadMessage.textContent = "Please select a file";
        return;
    }

    const formData = new FormData();
    formData.append("file", file);

    try {
        const res = await fetch("/upload", {
            method: "POST",
            body: formData
        });

        const text = await res.text();
        uploadMessage.textContent = text;
        listFiles(); // обновляем список после загрузки
    } catch (err) {
        uploadMessage.textContent = "Upload failed: " + err;
    }
}

// получение списка файлов
async function listFiles() {
    try {
        const res = await fetch("/list");
        const text = await res.text();
        fileList.textContent = text;
    } catch (err) {
        fileList.textContent = "Failed to fetch file list: " + err;
    }
}
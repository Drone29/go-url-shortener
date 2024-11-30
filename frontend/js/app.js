const backUrl = "/shorten"
const saveBtn = document.getElementById("saveBtn")
const searchBtn = document.getElementById("searchBtn")
const urlInput = document.getElementById("urlInput")
const responseMsg = document.getElementById("responseMsg")

// handle save button
saveBtn.addEventListener("click", async() => {
    const url = urlInput.value;
    if (!url) {
        alert("Please enter a valid URL");
        return;
    }
    fetch(`${backUrl}`, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({url})
    })
    .then(response => {
        if (response.ok) {
            // if 200..299
            const data = response.json();
            console.log(`Response ${response.status} ${response.statusText}\n${data}`);
            responseMsg.innerText = `Saved as ${data.shortCode}`
            urlInput.value = '' // clear input field
        } else {
            console.error(`Error: ${response.status} ${response.statusText}`);
            throw new Error("Failed to POST data");
        }
    })
    .catch(error => {
        console.error("POST error: ", error);
        responseMsg.innerText = "Failed to save URL";
    });
});

// handle search & redirect button
searchBtn.addEventListener("click", async() => {
    const key = urlInput.value;
    if (!key) {
        alert("Please enter a valid key");
        return;
    }
    fetch(`${backUrl}/${key}`)
    .then(response => {
        if (response.ok) {
            const data = response.json();
            console.log(`Response ${response.status} ${response.statusText}\n${data}`);
            // redirect to url
            if (data.url) {
                // open in a new window
                window.open(data.url, '_blank');
                responseMsg.innerText = `${data.url} opened in a new window`;
            } else {
                responseMsg.innerText = "URL not found";
            }
        } else {
            console.error(`Error: ${response.status} ${response.statusText}`);
            throw new Error("Failed to GET data");
        }
    })
    .catch(error => {
        console.error("GET error: ", error);
        responseMsg.innerText = "Failed to fetch URL";
    });
});
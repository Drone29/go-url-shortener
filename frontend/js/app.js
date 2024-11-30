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

    try {

         // check if entered url is valid
        try {
            new URL(url);
        } catch(error) {
            throw error
        }
        const response = await fetch(`${backUrl}`, {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify({url})
        });
        if (response.ok) {
            // if 200..299
            const data = await response.json();
            console.log(`Response ${response.status} ${response.statusText}\n` + JSON.stringify(data));
            responseMsg.innerText = `Saved as ${data.shortCode}`
            urlInput.value = '' // clear input field
        } else {
            throw new Error(await response.text());
        }
    } catch(error) {
        console.error("POST error: ", error);
        responseMsg.innerText = `${error}`;
    }
});

// handle search & redirect button
searchBtn.addEventListener("click", async() => {
    const key = urlInput.value;
    if (!key) {
        alert("Please enter a valid key");
        return;
    }
    try {
        const response = await fetch(`${backUrl}/${key}`);
        if (response.ok) {
            const data = await response.json();
            console.log(`Response ${response.status} ${response.statusText}\n` + JSON.stringify(data));
            // redirect to url
            if (data.url) {
                // open in a new window
                window.open(data.url, '_blank');
                responseMsg.innerText = `${data.url} opened in a new window`;
            } else {
                throw new Error("Invalid response from backend");
            }
            urlInput.value = '' // clear input field
        } else {
            throw new Error(await response.text());
        }
    } catch(error) {
        console.error("GET error: ", error);
        responseMsg.innerText = `${error}`;
    }
});
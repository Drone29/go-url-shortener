const backUrl = "/shorten";
const saveBtn = document.getElementById("saveBtn");
const searchBtn = document.getElementById("searchBtn");
const listBtn = document.getElementById("listBtn");
const urlInput = document.getElementById("urlInput");
const responseMsg = document.getElementById("responseMsg");

async function genericRequest(url, method, body = null) {
    const options = {
        method: method,
        headers: {"Content-Type": "application/json"},
    };
    if (body) {
        options.body = body;
    }
    const response = await fetch(url, options);
    if (!response.ok) {
        throw new Error(await response.text());
    }
    const data = await response.json();
    console.log(`Response ${response.status} ${response.statusText}\n` + JSON.stringify(data));
    return data
}

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
            throw error;
        }
        const data = await genericRequest(`${backUrl}`, "POST", JSON.stringify({url}));
        responseMsg.innerText = `Saved as ${data.shortCode}`
        urlInput.value = '' // clear input field
    }catch(error){
        responseMsg.innerText = `${error}`;
        console.error("Error: ", error);
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
        const data = await genericRequest(`${backUrl}/${key}`, "GET");
        // redirect to url
        if (data.url) {
            // open in a new window
            window.open(data.url, '_blank');
            responseMsg.innerText = `${data.url} opened in a new window`;
        } else {
            throw new Error("Invalid response from backend");
        }
        urlInput.value = '' // clear input field
    }catch(error){
        responseMsg.innerText = `${error}`;
        console.error("Error: ", error);
    }
});

// handle list
listBtn.addEventListener("click", async() => {
    try {
        const data = await genericRequest(`${backUrl}/list`, "GET");
        let result = "List:\n";
        for (let i = 0; i < data.length; i++) {
            const element = data[i];
            result += `${element.shortCode}: ${element.url}\n`;
        }
        responseMsg.innerText = result;
    }catch(error) {
        responseMsg.innerText = `${error}`;
        console.error("Error: ", error);
    }
});
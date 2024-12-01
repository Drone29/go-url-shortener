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

async function errorHandler(func) {
    try {
        await func()
    }catch(error){
        responseMsg.innerText = `${error}`;
        console.error("Error: ", error);
    }
}

// handle save button
saveBtn.addEventListener("click", () => 
    errorHandler(async() => {
        const url = urlInput.value;
        if (!url) {
            alert("Please enter a valid URL");
            return;
        }
        // check if entered url is valid
        new URL(url);

        const data = await genericRequest(`${backUrl}`, "POST", JSON.stringify({url}));
        responseMsg.innerText = `Saved as ${data.shortCode}`
        urlInput.value = '' // clear input field
}));

// handle search & redirect button
searchBtn.addEventListener("click", () =>
    errorHandler(async() => {
        const key = urlInput.value;
        if (!key) {
            alert("Please enter a valid key");
            return;
        }
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
}));

// handle list
listBtn.addEventListener("click", () =>
    errorHandler(async() => {
        const data = await genericRequest(`${backUrl}/list`, "GET");
        let result = "List:\n";
        for (let i = 0; i < data.length; i++) {
            const element = data[i];
            result += `${element.shortCode}: ${element.url}\n`;
        }
        responseMsg.innerText = result;
}));
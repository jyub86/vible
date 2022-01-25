const search = document.querySelector("#search")
const version = document.querySelector("#version")
const version2 = document.querySelector("#version2")
const clip = document.querySelector("#clip")
const title = document.querySelector("#title");
const title2 = document.querySelector("#title2");
const result = document.querySelector("#result");
const result2 = document.querySelector("#result2");
const input = document.querySelector("input");
const prev = document.querySelector("#prev");
const next = document.querySelector("#next");

function searchRequest(option) {
    value = {
        "word" : input.value,
        "version" : version.options[version.selectedIndex].value,
        "version2" : version2.options[version2.selectedIndex].value,
        "clip" : clip.checked,
        "option" : option,
    }
    fetch("/search", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(value)
    })
    .then(response => response.json())
    .catch(error => console.error('Error:', error))
    .then(function(data) {
        if ("error" in data) {
            result.innerHTML = data["error"]
            return
        }
        input.value = data["data"]["title"]
        title.innerHTML = "[ " + data["data"]["title"] + " ]" + " - " + version.options[version.selectedIndex].text
        result.innerHTML = data["data"]["result"]
        if (version2.options[version2.selectedIndex].value !== "") {
            title2.innerHTML = version2.options[version2.selectedIndex].text
            result2.innerHTML = data["data"]["result2"]
        }
        if (data["data"]["clip"]) {
            copyToClipboard(data["data"]["result"])
        }
    })
};

function copyToClipboard(value) {
    const tempElem = document.createElement('textarea');
    tempElem.value = value;  
    document.body.appendChild(tempElem);
    
    tempElem.select();
    document.execCommand("copy");
    document.body.removeChild(tempElem);
}

search.addEventListener('submit', function(event) {
    event.preventDefault();
    searchRequest("")
});

version.addEventListener("change", function(){
    searchRequest("")
})

prev.addEventListener('click', function(event) {
    searchRequest("prev")
});

next.addEventListener('click', function(event) {
    searchRequest("next")
});

version2.addEventListener('change', function(event) {
    searchRequest("")
});
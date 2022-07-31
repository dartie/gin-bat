/* AJAX definition */
function AJAX(url, data, method='GET', async=true){
    return new Promise(function(resolve, reject) {
        var formData = new FormData();

        for (var key in data) {
            // check if the property/key is defined in the object itself, not in parent
            if (data.hasOwnProperty(key)) {
                formData.append(key, data[key]);
            }
        }

        var xhr = new XMLHttpRequest();

        // Add any event handlers here...
        xhr.onload = function() {
            resolve(this.responseText);
        };
        xhr.onerror = reject;

        xhr.open(method, url, async);
        xhr.send(formData);
        //return false; // To avoid actual submission of the form
    });
}


function show_alert(message, style, parentElement){
    // "button" element
    let button_element = document.createElement("button");
    button_element.className = "close btn-close btn-close-dark ms-auto me-2";
    button_element.setAttribute("type", "button");
    button_element.setAttribute("data-dismiss", "alert");
    button_element.setAttribute("data-bs-dismiss", "alert");
    button_element.setAttribute("aria-label", "Close");

    // "div" element
    let div_element = document.createElement("div");
    div_element.className = "alert alert-" + style + " alert-dismissible fade show";
    div_element.setAttribute("role", "alert");
    div_element.textContent = message;

    // link elements
    div_element.appendChild(button_element);
    parentElement.appendChild(div_element);
}


function escapeDoubleQuotes(str) {
	return str.replace(/\\([\s\S])|(")/g,"\\$1$2");
}

function escapeHtml(unsafe) {
    return unsafe
         .replace(/&/g, "&amp;")
         .replace(/</g, "&lt;")
         .replace(/>/g, "&gt;")
         .replace(/"/g, "&quot;")
         .replace(/'/g, "&#039;");
 }


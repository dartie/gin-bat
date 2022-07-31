let fileN = 0;
let uploadedFiles_dict = {};
let unique_uploadedFiles = [];


function uniq(a) {
    return Array.from(new Set(a));
}


function FileListItems (files) {
    var b = new ClipboardEvent("").clipboardData || new DataTransfer()
    for (var i = 0, len = files.length; i<len; i++) b.items.add(files[i])
    return b.files
}


function display_elements (files, inputElement) {
    var parentElement = document.getElementById('filelist');

    /* clear the element list displayed */
    parentElement.innerHTML = '';

    /* create elements to display */
    for (var i = 0; i < files.length; i++) {
        fileN++;
        // get item
        // file = files.item(i);
        //or
        file = files[i];

        /*
        <div class="row input-group mt-1 ml-1" data-filename="Compliance Report-N020E13A8N02344308.xhtml">
            <input type="text" class="form-control col-4" readonly="" value="Compliance ">
            <button style="border-top-right-radius: 15px; border-bottom-right-radius: 15px;" class="btn btn-danger col-1"><span class="glyphicon glyphicon-plus">X</span></button>
            <div class="col-7">&nbsp;</div>
        </div>
        */

        // "row" element
        let row_element = document.createElement("div");
        row_element.className = "row mt-1 input-group ml-1";
        row_element.setAttribute('data-filename', file.name);

        // "remove button" element
        let remove_button_element = document.createElement("button");
        remove_button_element.className = "btn btn-danger col-1";
        remove_button_element.style.borderTopRightRadius = "15px";
        remove_button_element.style.borderBottomRightRadius = "15px";
        remove_button_element.style.maxWidth = "60px";
        //remove_button_element.onclick = "return remove_file(" + fileN + ")";
        remove_button_element.onclick = function () {
            let row_element = this.closest('.row');
            row_element.parentNode.removeChild(row_element);

            fileToRemove = row_element.getAttribute('data-filename');
            delete uploadedFiles_dict[fileToRemove];
            uploadedFiles_updated = [];
            for (var i=0; i < unique_uploadedFiles.length ; i++) {
                var file = unique_uploadedFiles[i];
                var filename = file.name;

                if (filename != fileToRemove) {
                    uploadedFiles_updated.push(file);
                }
            }
            unique_uploadedFiles = uploadedFiles_updated;

            delete uploadedFiles_dict[fileToRemove];

            // update the inputElement which has the contect actually submitted
            inputElement.files = new FileListItems(unique_uploadedFiles);

            return false;
        }

        // "span" element for remove button icon
        var icon_element = document.createElement("span");
        icon_element.className = "glyphicon glyphicon-plus";
        icon_element.textContent = 'X';

        // "input" element
        var new_input = document.createElement("input");
        new_input.type = "text";
        new_input.className = "form-control col-6";    // set the CSS class
        new_input.value = file.name;
        new_input.readOnly = true;

        // "div" column offset element
        let col_offset = document.createElement("div");
        col_offset.className = "col-5";
        col_offset.textContent = ""; // "&nbsp;"

        // Link elements
        row_element.appendChild(new_input);
        row_element.appendChild(remove_button_element);
        remove_button_element.appendChild(icon_element);
        row_element.appendChild(col_offset);

        parentElement.appendChild(row_element);     // put it into the DOM
    }
}


function add_elements (inputElement) {
    var files = inputElement.files;

    for (var i = 0; i < files.length; i++) {
        fileN++;
        // get item
        // file = files.item(i);
        //or
        file = files[i];

        var filename = file.name;

        // add file only in case it has not been added already (the check is performed on the )
        if (filename in uploadedFiles_dict === false) {
            unique_uploadedFiles.push(file);
            uploadedFiles_dict[filename] = file;
        }
    }

    inputElement.files = new FileListItems(unique_uploadedFiles);

    display_elements(unique_uploadedFiles, inputElement);

}

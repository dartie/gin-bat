function set_info(id_element, content){
    let html_element = document.getElementById(id_element);

    if (html_element) {
        html_element.innerHTML = content;
    }
    else
    {
        ;
    }
}


function get_local_date_to_display(lastupdate){
    let date = new Date(lastupdate);
    let lastupdate_local = date.toString(); // "Wed Jun 29 2011 09:52:48 GMT-0700 (PDT)"
    let lastupdate_local_to_display = lastupdate_local.split(" ", 5).join(" ");

    if (lastupdate_local_to_display == "Invalid Date"){
        lastupdate_local_to_display = "";
    }

    return lastupdate_local_to_display
}


function adjust_timezone_date(){
    var datetime_strings = document.getElementsByClassName("lastupdate");

    for (var i=0; i < datetime_strings.length; i++){
        var python_date = datetime_strings[i].textContent;
        datetime_strings[i].textContent = get_local_date_to_display(python_date);
    }
}


function update(){
    //let settings_json = JSON.stringify(settings_dict);

    AJAX("/update", "GET")
    .then(function(result) {
        // Code depending on result
        //console.log(result);

        result = JSON.parse(result);

        var widgets_to_not_update_elements = document.getElementsByClassName("write");
        var widgets_to_not_update = [];

        for(let i=0; i < widgets_to_not_update_elements.length; i++){
            widgets_to_not_update.push(widgets_to_not_update_elements[i].id);
        }

        for(let k in result){
            let v = result[k];
            let html_data = v[0];
            let lastupdate = v[1];

            if (lastupdate !== "") {
                var lastupdate_local_display = get_local_date_to_display(lastupdate);
            }
            else{
                var lastupdate_local_display = "";
            }

            var id_widget_to_write = k + "_cardtext";
            var name_widget_to_write = k;
            if (widgets_to_not_update.includes(id_widget_to_write)){
                //continue;

                var widget_info_json = {};
                widget_info_json["id"] = name_widget_to_write;

                var innerHTML_text = document.getElementById(id_widget_to_write).innerHTML;
                widget_info_json["content"] = escapeHtml(innerHTML_text.replace(/[\n]/g, "<<CR>>"));
                widget_info_json = JSON.stringify(widget_info_json);

                AJAX("/write_widget/" + widget_info_json, "POST")
                .then(function(result) {
                    // Code depending on result
                    console.log(result);

                })
                .catch(function() {
                    // An error occurred
                    console.log("error");
                });
            }

            // update element, getting it by Id which has the module name
            set_info(k + "_cardtext", html_data);
            set_info(k + "_lastupdate", lastupdate_local_display);
        }

    })
    .catch(function() {
        // An error occurred
        console.log("error");
    });

    setTimeout(update, 15000);
}

adjust_timezone_date();
setTimeout(update, 15000);

function getConfirmation() {
    var retVal = confirm("Would you like to drilldown on this data?");
    if (retVal == true) {
        $('#msg').val("how many of those are prime/subprime");
        $('#submitBtn').click();
        return true;
    } else {
        return false;
    }
}


$(document).on("click", ".clickr", function(e) {
    e.preventDefault();
    $('#msg').val($(this).html());
    $('#submitBtn').click();
});

$(document).on("click", "#submitBtn", function(e) {
    e.preventDefault();
    console.log("Clicked submit");
    $("#form").submit();
});

$(document).on("click", ".dickr", function(e) {
    e.preventDefault();
    $('#msg').val($(this).html());
});

$(document).on("click", ".expand", function(e) {
    e.preventDefault();
    $(this).siblings(".bar").toggle();
});

$(document).on("click", ".filterdropdown", function(e) {
    e.preventDefault();
    console.log($(this).html());
});

function setupAutoComplete(userName) {
    $.post("/autocomplete", {
            username: userName
        },
        function(data) {
            $("#msg").autocomplete({
                source: function(req, response) {
                    var results = $.ui.autocomplete.filter(data, req.term);
                    response(results.slice(0, 5)); //for getting 5 results
                },
                open: function(event, ui) {
                    var $input = $(event.target),
                        $results = $input.autocomplete("widget"),
                        top = $results.position().top,
                        height = $results.height(),
                        inputHeight = $input.height(),
                        newTop = top - height - inputHeight;

                    $results.css("top", newTop + "px");
                }
            });
        }
    );
}

function sendFeedback(id, feedback) {
    $.post("/feedback", {
            reqId: id,
            msg: feedback
        },
        function(data) {
            swal(data.message);
        }
    );
}

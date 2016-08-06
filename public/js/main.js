(function () {
    // Update the live data every second.
    setInterval(function () {
        $.getJSON('/livedata', function(data) {
            $.each(data, function(key, value) {
                $("#" + key).text(value);
            });
        });
    }, 200);   
}());

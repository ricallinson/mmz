(function () {

    if (!document.vis) {
        return;
    }

    var graphs = {};

    function makeGraph(nodeId, min, max) {
        var container = document.getElementById(nodeId);
        var dataset = new vis.DataSet();
        var options = {
            start: vis.moment().add(-30, 'seconds'),
            end: vis.moment(),
            graphHeight: '100px',
            showMajorLabels: false,
            showMinorLabels: false,
            dataAxis: {
                left: {
                    range: {
                        min: min, max: max
                    }
                }
            },
            drawPoints: {
                style: 'circle' // square, circle
            },
            shaded: {
                orientation: 'bottom' // top, bottom
            }
        };
        graphs[nodeId] = {
            graph: new vis.Graph2d(container, dataset, options),
            dataset: dataset
        };
    }

    function updateGraph(nodeId, Timestamp, MotorKilowatts) {
        var graph = graphs[nodeId].graph;
        var dataset = graphs[nodeId].dataset;
        // Add a new data point to the dataset.
        dataset.add({
            x: new Date(1000*Timestamp),
            y: MotorKilowatts
        });
        // Remove all data points which are no longer visible.
        var range = graph.getWindow();
        var interval = range.end - range.start;
        var oldIds = dataset.getIds({
            filter: function (item) {
                return item.x < range.start - interval;
            }
        });
        dataset.remove(oldIds);
    }

    function renderGraphs() {
        Object.keys(graphs).forEach(function (nodeId) {
            if (graphs[nodeId]) {
                var graph = graphs[nodeId].graph;
                var now = vis.moment();
                var range = graph.getWindow();
                var interval = range.end - range.start;
                graph.setWindow(now - interval, now, {animation: false});
            }
                
        });
        requestAnimationFrame(renderGraphs);
    }

    // Create a graphs.
    makeGraph('MotorKilowattsVis', 0, 1000);
    makeGraph('MotorVoltageVis', 0, 300);
    makeGraph('BatteryVoltageVis', 200, 300);
    makeGraph('AverageCurrentOnMotorVis', 0, 2000);
    makeGraph('AvailableCurrentFromControllerVis', 0, 2000);
    makeGraph('ControllerTempVis', 0, 200);

    // Update the graphs data every second.
    setInterval(function () {
        $.getJSON('/livedata', function(data) {
            $.each(data, function(nodeId, value) {
                if (graphs[nodeId + 'Vis']) {
                    updateGraph(nodeId + 'Vis', data.Timestamp, value);
                }
                $('#' + nodeId).text(value);
            });
        });
    }, 1000);

    renderGraphs();
}());

(function () {
    $('input').change(function(e) {
        $(this).addClass('dirty');
        var id = this.id
        $.getJSON('/set/' + id + '?value=' + this.value, function (data) {
            console.log(data);
            if (data.status) {
                $('#' + id).removeClass('dirty');
            }
        });
    });
}());

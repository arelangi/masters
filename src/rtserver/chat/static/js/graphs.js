
function columnGraph(ObjectId, Categories, Series, TitleText) {
    Highcharts.chart(ObjectId, {
        chart: {
            type: 'column'
        },
        title: {
            text: TitleText
        },
        xAxis: {
            categories: Categories,
            crosshair: true
        },
        yAxis: {
            min: 0,
            title: {
                text: 'Value'
            }
        },
        tooltip: {
            headerFormat: '<span style="font-size:10px">{point.key}</span><table>',
            pointFormat: '<tr><td style="color:{series.color};padding:0">{series.name}: </td>' +
                '<td style="padding:0"><b>{point.y:.1f}</b></td></tr>',
            footerFormat: '</table>',
            shared: true,
            useHTML: true
        },
        plotOptions: {
            column: {
                pointPadding: 0.2,
                borderWidth: 0
            },
            series: {
                cursor: 'pointer',
                point: {
                    events: {
                        click: function() {
                            console.log(this);
                            getConfirmation();
                        }
                    }
                }
            }
        },
        series: Series
    });
}

function lineGraph(ObjectId, Categories, Series, TitleText) {
    Highcharts.chart(ObjectId, {
        chart: {
            type: 'line'
        },
        title: {
            text: TitleText
        },
        xAxis: {
            categories: Categories
        },
        yAxis: {
            title: {
                text: 'Volume'
            }
        },
        legend: {
            layout: 'vertical',
            align: 'right',
            verticalAlign: 'middle'
        },
        plotOptions: {
            line: {
                dataLabels: {
                    enabled: false
                },
                enableMouseTracking: true
            },
            series: {
                cursor: 'pointer',
                point: {
                    events: {
                        click: function() {
                            console.log(this);
                        }
                    }
                }
            }
        },
        series: Series
    });
}

function combinationChart(ObjectId, Categories, Series, TitleText) {
    Highcharts.chart(ObjectId, {
        chart: {
            zoomType: 'xy'
        },
        title: {
            text: TitleText
        },
        xAxis: [{
            categories: Categories,
            crosshair: true
        }],
        yAxis: [{ // Primary yAxis
            labels: {
                format: '{value}',
                style: {
                    color: Highcharts.getOptions().colors[1]
                }
            },
            title: {
                text: 'Others',
                style: {
                    color: Highcharts.getOptions().colors[1]
                }
            }
        }, { // Secondary yAxis
            title: {
                text: 'CapOne',
                style: {
                    color: Highcharts.getOptions().colors[0]
                }
            },
            labels: {
                format: '{value}',
                style: {
                    color: Highcharts.getOptions().colors[0]
                }
            },
            opposite: true
        }],
        tooltip: {
            shared: true
        },
        legend: {
            layout: 'vertical',
            align: 'left',
            x: 120,
            verticalAlign: 'top',
            y: 100,
            floating: true,
            backgroundColor: (Highcharts.theme && Highcharts.theme.legendBackgroundColor) || '#FFFFFF'
        },
        plotOptions: {
            series: {
                cursor: 'pointer',
                point: {
                    events: {
                        click: function() {
                            console.log(this);
                        }
                    }
                }
            }
        },
        series: [{
            name: Series[0].name,
            type: Series[0].type,
            yAxis: 1,
            data: Series[0].data,
            tooltip: {
                valueSuffix: ' loans'
            }
        }, {
            name: Series[1].name,
            type: Series[1].type,
            data: Series[1].data,
            tooltip: {
                valueSuffix: ' loans'
            }
        }]
    });
}
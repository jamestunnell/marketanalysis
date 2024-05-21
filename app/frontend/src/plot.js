import van from 'vanjs-core'
import  ndjsonParser from 'ndjson-parse'
import { parseISO } from 'date-fns'
import Highcharts from 'highcharts'
import ApexCharts from 'apexcharts'

import { ButtonIcon } from './buttons.js'
import { IconClose } from './icons.js'
import { ModalBackground, ModalForeground } from './modal.js'

const {div} = van.tags

const COLORS = ["royalblue", "seagreen", "plum", "tomato", "goldenrod", "sienna"]

function plotColor(idx) {
    return COLORS[idx % COLORS.length]
}

function tsToUnix(ts) {
    const date = parseISO(ts)
    
    return Math.floor(date.getTime());
}

const makePlot = (series) => {
    const plotArea = div()
    const chart = Highcharts.chart(plotArea, {
        chart: {type: 'line', zooming: {type: 'x'}},
        title: {enabled: false, text: ""},
        xAxis: {type: 'datetime'},
        yAxis: {    },
        plotOptions: {
            series:{
                tooltip: {shared: true},
                states: {
                    hover: {enabled: false}
                },
            },
        },
        legend: {
            enabled: true
        },
        series: series,
        navigation: {
            buttonOptions: {
                enabled: true
            }
        }
    });

    return plotArea
}

function plotNDJSON(text) {
    const records = ndjsonParser(text)
    const names = Object.keys(records[0].values)
    const series = names.map((name, i) => {
        const valuePairs = records.map(r => [tsToUnix(r.timestamp), r.values[name]])
        
        return {name, data: valuePairs, color: plotColor(i)}
    })

    return makePlot(series)
}

const PlotModal = ({text, format}) => {
    if (format !== "ndjson") {
        console.log(`plotting format ${format} is not supported`)

        return
    }

    const closed = van.state(false);
    const closeBtn = ButtonIcon({icon: IconClose(), onclick: () => closed.val = true})
    const plot = plotNDJSON(text)

    const modal = ModalBackground(
        div(
            {class: "block p-16 rounded-lg bg-white min-w-[65%] max-w-[65%]"},
            div(
                {class: "flex flex-col"},
                div({class:"flex self-end"}, closeBtn),
                plot
            )
        ),
    )

    // plot results in a follow-up modal
    van.add(document.body, () => closed.val ? null : modal)
}

export { PlotModal }
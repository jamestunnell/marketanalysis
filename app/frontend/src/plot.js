import van from 'vanjs-core'
import  ndjsonParser from 'ndjson-parse'
import { parseISO } from 'date-fns'
import Highcharts from 'highcharts'

import { ButtonIcon } from './buttons.js'
import { IconClose } from './icons.js'
import { ModalBackground, ModalForeground } from './modal.js'

const {div, p} = van.tags

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
        yAxis: [
            {
                title: {
                    text: 'Price',
                },
            },
            {
                title: {
                    text: 'Oscillator',
                },
                opposite: true,
            },
        ],
        tooltip: {
            shared: true
        },
        legend: {
            enabled: true
        },
        series: series,
    });

    return plotArea
}

function guessYAxis(values) {
    const average = values.reduce((partialSum, a) => partialSum + a, 0.0) / Number(values.length)

    return (Math.abs(average) < 5.0) ? 1 : 0
}

function plotNDJSON(text) {
    const records = ndjsonParser(text)
    const names = Object.keys(records[0].values)
    const series = names.map((name, i) => {
        const values = records.map(r => r.values[name])
        const yAxis = guessYAxis(values)
        const valuePairs = records.map((r,i) => [tsToUnix(r.timestamp), values[i]])
        
        console.log(`series ${name} is using y-axis ${yAxis}`)
        return {name, data: valuePairs, color: plotColor(i), yAxis}
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

    const modal = ModalBackground(
        div(
            {class: "block p-16 rounded-lg bg-white min-w-[65%] max-w-[65%]"},
            div(
                {class: "flex flex-col"},
                div(
                    p({class: "text-xl"}, "Run Results"),
                    div({class: "float-end"}, closeBtn),
                ),
                plotNDJSON(text)
            )
        ),
    )

    // plot results in a follow-up modal
    van.add(document.body, () => closed.val ? null : modal)
}

export { PlotModal }
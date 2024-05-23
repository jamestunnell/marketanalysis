import van from 'vanjs-core'
import  ndjsonParser from 'ndjson-parse'
import { parseISO } from 'date-fns'
import Highcharts from 'highcharts'

import { ButtonIcon } from './buttons.js'
import { IconClose } from './icons.js'
import { kMeansAdaptive } from './clustering/cluster.js'
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
        yAxis: {},
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

// function guessYAxis(values) {
//     const average = values.reduce((partialSum, a) => partialSum + a, 0.0) / Number(values.length)

//     return (Math.abs(average) < 5.0) ? 1 : 0
// }

function plotNDJSON(text) {
    const records = ndjsonParser(text)
    const names = Object.keys(records[0].values)
    const datasets = Object.fromEntries(names.map(name => {
        const values = records.map(r => {
            return r.values[name] ? r.values[name] : null
        })

        return [name, values]
    }))

    console.log(`made result datasets`, datasets)
    
    const clusters = kMeansAdaptive(datasets)

    const plots = clusters.map(members => {
        console.log("making plot with cluster members:", ...members)

        const series = members.map((name, i) => {
            const values = datasets[name]
            const valuePairs = records.map((r,i) => [tsToUnix(r.timestamp), values[i]])
            
            return {name, data: valuePairs, color: plotColor(i)}
        })
    
        return makePlot(series)    
    })

    return plots
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
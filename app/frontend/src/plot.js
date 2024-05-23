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

function plotRecording(recording) {
    const qs = {}

    recording.quantities.forEach(q => {
        if (q.records.length === 0) {
            console.log(`ignoring quantity ${q.name} with no records`)

            return
        }

        console.log(`keeping quantities ${q.name}`)

        qs[q.name] = q
    })

    const valuesByName = Object.fromEntries(Object.entries(qs).map(([name,q]) => {
        const values = q.records.map(r => r.v)
        
        return [name, values]
    }))
    const clusters = kMeansAdaptive(valuesByName)
    const plots = clusters.map(members => {
        console.log(`plotting quantities: ${members}`)

        const series = members.map((name, i) => {
            const valuePairs = qs[name].records.map(r => [tsToUnix(r.t), r.v])
            
            return {name, data: valuePairs, color: plotColor(i)}
        })
    
        return makePlot(series)    
    })

    return plots
}

const PlotRecordingModal = (recording) => {
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
                plotRecording(recording)
            )
        ),
    )

    // plot results in a follow-up modal
    van.add(document.body, () => closed.val ? null : modal)
}

export { PlotRecordingModal }
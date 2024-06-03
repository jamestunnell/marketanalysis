import van from 'vanjs-core'
import { parseISO } from 'date-fns'
import Highcharts from 'highcharts'

import { ButtonIcon } from './buttons.js'
import { IconClose } from './icons.js'
import { kMeansAdaptive } from './clustering/cluster.js'
import { ModalBackground, ModalForeground } from './modal.js'
import userTimeZone from './timezone.js'

const {div, p} = van.tags

const COLORS = ["royalblue", "seagreen", "plum", "tomato", "goldenrod", "sienna"]

function plotColor(idx) {
    return COLORS[idx % COLORS.length]
}

function tsToUnix(ts) {
    const date = parseISO(ts)
    
    return Math.floor(date.getTime());
}

const makePlot = ({series, height}) => {
    const plotArea = div()
    const chart = Highcharts.chart(plotArea, {
        chart: {
            type: 'line',
            zooming: {type: 'x'},
            height,
        },
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
        time: {
            timezone: userTimeZone(),
        },
    });

    return plotArea
}

function plotRecording({recording, totalHeight}) {
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
    const height = totalHeight / Math.min(3, clusters.length)
    const plots = clusters.map(members => {
        console.log(`plotting quantities: ${members}`)

        const series = members.map((name, i) => {
            const valuePairs = qs[name].records.map(r => [tsToUnix(r.t), r.v])
            
            return {name, data: valuePairs, color: plotColor(i)}
        })
    
        return makePlot({series, height})
    })

    return plots
}

const PlotRecordingModal = (recording) => {
    const closed = van.state(false);
    const closeBtn = ButtonIcon({
        icon: IconClose(),
        text: "Close",
        onclick: () => closed.val = true},
    )
    const totalHeight = window.screen.availHeight * 0.9 * 0.75

    const modal = ModalBackground(
        div(
            {class: "block p-8 rounded-lg bg-white min-w-[70%] max-w-[90%] min-h-[70%] max-h-[90%] overflow-y-auto"},
            div(
                {class: "flex flex-col"},
                div(
                    p({class: "text-xl"}, "Run Results"),
                    div({class: "float-end"}, closeBtn),
                ),
                plotRecording({recording, totalHeight})
            )
        ),
    )

    // plot results in a follow-up modal
    van.add(document.body, () => closed.val ? null : modal)
}

export { PlotRecordingModal }
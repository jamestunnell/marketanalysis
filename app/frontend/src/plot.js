import van from 'vanjs-core'
import { parseISO } from 'date-fns'
import Highcharts from 'highcharts'
import mouseWheelZoom from "highcharts/modules/mouse-wheel-zoom";

import { ButtonIcon } from './buttons.js'
import { IconClose } from './icons.js'
import { kMeansAdaptive } from './clustering/cluster.js'
import { ModalBackground, ModalForeground } from './modal.js'
import userTimeZone from './timezone.js'

const {div, p} = van.tags

const COLORS = ["royalblue", "seagreen", "mediumvioletred", "tomato", "goldenrod", "sienna"]

mouseWheelZoom(Highcharts)

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
            panning: true,
            panKey: 'shift',
            animation: false,
            type: 'line',
            zooming: {
                enabled: true,
                mouseWheel: true,
                type: "x",
              },
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

function PlotRecording({recording, totalHeight}) {
    const qs = {}

    recording.quantities.forEach(q => {
        if (q.records.length === 0) {
            console.log(`ignoring quantity ${q.name} with no records`)

            return
        }

        // console.log(`keeping quantities ${q.name}`)

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
            
            return {name, data: valuePairs, color: plotColor(i), animation: false}
        })
    
        return makePlot({series, height})
    })

    return plots
}

export { PlotRecording }
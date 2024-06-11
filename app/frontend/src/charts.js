import van from 'vanjs-core'
import { parseISO } from 'date-fns'
import Highcharts from 'highcharts'
import mouseWheelZoom from "highcharts/modules/mouse-wheel-zoom";

// import { kMeansAdaptive } from './clustering/cluster.js'
import userTimeZone from './timezone.js'

const {div, p} = van.tags

const COLORS = ["royalblue", "seagreen", "mediumvioletred", "tomato", "goldenrod", "sienna"]

mouseWheelZoom(Highcharts)

function plotColor(idx) {
    return COLORS[idx % COLORS.length]
}

function randomPlotColor() {
    const randomIdx = Math.floor(Math.random() * COLORS.length);

    return COLORS[randomIdx]
}

function tsToUnix(ts) {
    const date = parseISO(ts)
    
    return Math.floor(date.getTime());
}

const MakeEmptyChart = ({container, height}) => {
    return Highcharts.chart(container, {
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
        series: [],
        time: {
            timezone: userTimeZone(),
        },
    });
}

function UpdateCharts({charts, recording}) {
    const qAttribs = Object.fromEntries(recording.quantities.map(q => [q.name, q.attributes]))
    const ts = new Set()
    recording.quantities.forEach(q => q.records.forEach(r => ts.add(tsToUnix(r.t))))

    const minTS = Math.min(...ts.values())
    const maxTS = Math.max(...ts.values())

    console.log(`updating charts`, {qAttribs, minTS, maxTS})
    
    const updated = []
    // Update/remove existing series
    charts.forEach((c, cIdx) => {
        const toRemove = []
        c.series.forEach(s => {
            const q = recording.quantities.find(q => q.name === s.name)
            
            // remove series that don't correspond to a quantity or 
            // whose quantity has no data or assigned (chart) cluster
            if (!q || q.records.length === 0 || q.attributes.cluster === null || q.attributes.cluster !== cIdx) {
                console.log(`removing series ${s.name} with no quantity or no records or no assigned cluster or cluster changed`)

                toRemove.push(s.name)

                return
            }

            const valuePairs = q.records.map(r => [tsToUnix(r.t), r.v])
            
            s.setData(valuePairs, false)

            if (q.attributes.cluster >= charts.length) {
                console.log(`updating series ${s.name}`)
            }

            updated.push(q.name)
        })

        toRemove.forEach(name => {
            // remove but don't redraw chart
            c.series.find(s => s.name === name).remove(false)
        })
    })

    // Add new series
    recording.quantities.forEach(q => {
        if (updated.indexOf(q.name) >= 0) {
            return
        }

        if (q.records.length === 0) {
            console.log(`ignoring quantity ${q.name} with no records`)

            return
        }

        if (q.attributes.cluster === null) {
            console.log(`ignoring quantity ${q.name} with no assigned cluster`)

            return
        }

        if (q.attributes.cluster >= charts.length) {
            console.log(`ignoring quantity ${q.name} with cluster ${q.attributes.cluster} > max ${charts.length-1}`)

            return
        }

        const valuePairs = q.records.map(r => [tsToUnix(r.t), r.v])
        const c = charts[q.attributes.cluster]

        console.log(`adding series for ${q.name} with cluster ${q.attributes.cluster}`)

        c.addSeries({
            name: q.name,
            data: valuePairs,
            color: randomPlotColor(),
            animation: false,
        })
    })

    charts.forEach(c => {
        c.xAxis[0].setExtremes(minTS, maxTS, true)
    })
}

export { MakeEmptyChart, UpdateCharts }
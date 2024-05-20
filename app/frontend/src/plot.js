import van from 'vanjs-core'
import  ndjsonParser from 'ndjson-parse'
import { parseISO } from 'date-fns'
import uPlot from 'uplot'

import { ButtonIcon } from './buttons.js'
import { IconClose } from './icons.js'
import { ModalBackground, ModalForeground } from './modal.js'

const {div} = van.tags

const COLORS = ["royalblue", "seagreen", "plum", "tomato", "goldenrod", "sienna"]

function color(idx) {
    return COLORS[idx % COLORS.length]
}

function tsToUnix(ts) {
    const date = parseISO(ts)
    
    return Math.floor(date.getTime() / 1000);
}

const PlotNDJSON = (text) => {
    const records = ndjsonParser(text)
    const names = Object.keys(records[0].values)
    const data = [records.map(r => tsToUnix(r.timestamp))]
    
    names.forEach(name => data.push(records.map(r => r.values[name])))

    const opts = {
        title: "Run Results",
        id: "results-chart",
        // class: "my-chart",
        width: 800,
        height: 600,
        axes: [
            {},
            {
                scale: "price",
                values: (u, vals, space) => vals.map(v => "$"+v.toFixed(2)),
            }
        ],
        series: [{}].concat(names.map((name, i) => {
            return {
                show: true,
                spanGaps: false, // don't connect missing points
                label: name,
                stroke: color(i),
                scale: "price",
                value: (u,v) => v === null ? null : ("$" + v.toFixed(2)),
            }
        })),
    }

    const chartArea = div()

    let u = uPlot(opts, data, chartArea)

    return chartArea
}

const PlotModal = ({text, format}) => {
    if (format !== "ndjson") {
        console.log(`plotting format ${format} is not supported`)

        return
    }

    const closed = van.state(false);
    const closeBtn = ButtonIcon({icon: IconClose(), onclick: () => closed.val = true})
    const plot = PlotNDJSON(text)

    closeBtn.classList.add("self-end")

    const modal = ModalBackground(
        ModalForeground(
            {},
            div(
                {class: "flex flex-col"},
                closeBtn,
                plot,
            )
        ),
    )

    // plot results in a follow-up modal
    van.add(document.body, () => closed.val ? null : modal)
}

export { PlotModal }
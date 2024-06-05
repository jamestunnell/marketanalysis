import van from 'vanjs-core'

import { GraphActionModal } from './graphaction.js'
import { INPUT_CLASS } from './input.js'
import { PostJSON } from './backend.js'
import userTimeZone from './timezone.js'

const {input, label, option, select} = van.tags

const backtestGraph = ({id, date, symbol, showWarmup, predictor}) => {
    return new Promise((resolve, reject) => {
        const route = `/graphs/${id}/backtest`
        const timeZone = userTimeZone()
        const object = {date, symbol, timeZone, showWarmup, predictor}
        const options = {accept: 'application/json'}

        console.log("backtesting graph", object)

        PostJSON({route, object, options}).then(resp => {
            if (resp.status != 200) {
                resp.json().then(appErr => {
                    console.log("failed to backtest graph", appErr);
    
                    reject(appErr);    
                })
            }

            resp.json().then(obj => resolve(obj))
        }).catch(err => {
            console.log("failed to send backtest graph request", err)
            
            reject({
                title: "Action Failed",
                message: "failed to send backtest graph request",
                details: [err.message],
            })
        });
    });
}

const BacktestGraph = ({graph, settings, infoByType}) => {
    const predictor = van.state("")
    const threshold = van.state(0.25)

    const doAction = ({date, symbol, showWarmup}) => {
        return backtestGraph({id: graph.id, symbol, date, showWarmup, predictor: predictor.val})
    }
    const runDisabled = van.derive(() => predictor.val.length === 0)

    const blockOuts = []
    graph.blocks.forEach(blk => {
        infoByType[blk.type].outputs.forEach(out => {
            blockOuts.push(blk.name + "." + out.name)
        })
    })
    
    console.log(`made ${blockOuts.length} block outputs`, blockOuts)

    const predBlockOutOpts = [ option({value:"", selected: true}, "") ].concat(
        blockOuts.map(blkOut => option({value: blkOut}, blkOut))
    )

    const inputElems = [
        label({for: "predictor"}, "Predictor"),
        select({
            id: "predictor",
            class: INPUT_CLASS,
            oninput: e => predictor.val = e.target.value,
            required: true,
        }, predBlockOutOpts),
        label({for: "threshold"}, "Threshold"),
        input({
            id: "threshold",
            type: "number",
            class: INPUT_CLASS,
            value: threshold.val,
            min: -1.0,
            max: 1.0,
            step: 0.01,
            onchange: e => threshold.val = parseFloat(e.target.value)
        }),
    ]

    GraphActionModal({ actionName: "backtest", graph, settings, inputElems, runDisabled, doAction })
}

export { BacktestGraph };
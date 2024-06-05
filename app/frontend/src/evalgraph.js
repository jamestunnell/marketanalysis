import van from 'vanjs-core'

import { PostJSON } from './backend.js'
import { GraphActionModal } from './graphaction.js'
import { INPUT_CLASS } from './input.js'
import userTimeZone from './timezone.js'

const {input, label, option, select} = van.tags

const evalGraph = ({id, symbol, date, showWarmup, source, predictor, horizon}) => {
    return new Promise((resolve, reject) => {
        const route = `/graphs/${id}/eval`
        const timeZone = userTimeZone()
        const object = {type: "slope", symbol, date, timeZone, showWarmup, source, predictor, horizon}
        const options = {accept: 'application/json'}

        console.log("evaluating graph", object)

        PostJSON({route, object, options}).then(resp => {
            if (resp.status != 200) {
                resp.json().then(appErr => {
                    console.log("failed to evaluate graph", appErr);
    
                    reject(appErr);    
                })
            }

            resp.json().then(obj => resolve(obj))
        }).catch(err => {
            console.log("failed to send eval graph request", err)
            
            reject({
                title: "Action Failed",
                message: "failed to send eval graph request",
                details: [err.message],
            })
        });
    });
}

const HORIZON_MIN = 3
const HORIZON_MAX = 100

const EvalGraph = ({graph, selectedSymbol, infoByType}) => {
    const horizon = van.state(HORIZON_MIN)
    const source = van.state("")
    const predictor = van.state("")

    const doAction = ({date, symbol, showWarmup}) => {
        return evalGraph({
            id: graph.id,
            symbol,
            date,
            showWarmup,
            horizon: horizon.val,
            source: source.val,
            predictor: predictor.val,
        })
    }
    const runDisabled = van.derive(() => {
        return (
            (horizon.val < HORIZON_MIN) || 
            (horizon.val > HORIZON_MAX) ||
            (source.val.length === 0) || 
            (predictor.val.length === 0)
        )
    })

    const blockOuts = []
    graph.blocks.forEach(blk => {
        infoByType[blk.type].outputs.forEach(out => {
            blockOuts.push(blk.name + "." + out.name)
        })
    })
    
    console.log(`made ${blockOuts.length} block outputs`, blockOuts)

    const sourceBlockOutOpts = [ option({value:"", selected: true}, "") ].concat(
        blockOuts.map(blkOut => option({value: blkOut}, blkOut))
    )
    const predBlockOutOpts = [ option({value:"", selected: true}, "") ].concat(
        blockOuts.map(blkOut => option({value: blkOut}, blkOut))
    )

    const inputElems = [
        label({for: "horizon"}, "Horizon"),
        input({
            id: "horizon",
            type: "number",
            class: INPUT_CLASS,
            value: horizon.val,
            min: HORIZON_MIN,
            max: HORIZON_MAX,
            step: 1,
            onchange: e => horizon.val = Number(e.target.value),
            required: true,
        }),
        label({for: "source"}, "Source"),
        select({
            id: "source",
            class: INPUT_CLASS,
            oninput: e => source.val = e.target.value,
            required: true,
        }, sourceBlockOutOpts),
        label({for: "predictor"}, "Predictor"),
        select({
            id: "predictor",
            class: INPUT_CLASS,
            oninput: e => predictor.val = e.target.value,
            required: true,
        }, predBlockOutOpts),
    ]
  
    GraphActionModal({ actionName: "evaluate", graph, selectedSymbol, inputElems, runDisabled, doAction })
}

export { EvalGraph };
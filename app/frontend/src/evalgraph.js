import van from 'vanjs-core'

import { PostJSON } from './backend.js'
import { GraphActionModal, INPUT_CLASS } from './graphaction.js'

const {input, label, option, select} = van.tags

const evalGraph = ({id, date, symbol, source, predictor, window}) => {
    return new Promise((resolve, reject) => {
        const route = `/graphs/${id}/eval`
        const object = {type: "slope", date, symbol, source, predictor, window}
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

const WINDOW_MIN = 3
const WINDOW_MAX = 100

const EvalGraph = (graph, infoByType) => {
    const window = van.state(WINDOW_MIN)
    const source = van.state("")
    const predictor = van.state("")

    const doAction = ({date, symbol}) => {
        return evalGraph({
            id: graph.id,
            symbol,
            date,
            window: window.val,
            source: source.val,
            predictor: predictor.val,
        })
    }
    const runDisabled = van.derive(() => {
        return (
            (window.val < WINDOW_MIN) || 
            (window.val > WINDOW_MAX) ||
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
        label({for: "window"}, "Window Size"),
        input({
            id: "window",
            type: "number",
            class: INPUT_CLASS,
            value: window.val,
            min: WINDOW_MIN,
            max: WINDOW_MAX,
            step: 1,
            onchange: e => window.val = Number(e.target.value),
        }),
        label({for: "source"}, "Source"),
        select({
            id: "source",
            class: INPUT_CLASS,
            oninput: e => source.val = e.target.value,
        }, sourceBlockOutOpts),
        label({for: "predictor"}, "Predictor"),
        select({
            id: "predictor",
            class: INPUT_CLASS,
            oninput: e => predictor.val = e.target.value,
        }, predBlockOutOpts),
    ]
  
    GraphActionModal({ actionName: "evaluate", graph, inputElems, runDisabled, doAction })
}

export { EvalGraph };
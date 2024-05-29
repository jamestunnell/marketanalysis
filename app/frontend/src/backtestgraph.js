import van from 'vanjs-core'

import { GraphActionModal, INPUT_CLASS } from './graphaction.js'
import { PostJSON } from './backend.js'

const {label, option, select} = van.tags

const backtestGraph = ({id, date, symbol, source, predictor, window}) => {
    return new Promise((resolve, reject) => {
        const route = `/graphs/${id}/backtest`
        const object = {date, symbol, predictor}
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

const BacktestGraph = (graph, infoByType) => {
    const predictor = van.state("")

    const doAction = ({date, symbol}) => {
        return backtestGraph({id: graph.id, symbol, date, predictor: predictor.val})
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
    ]

    GraphActionModal({ actionName: "backtest", graph, inputElems, runDisabled, doAction })
}

export { BacktestGraph };
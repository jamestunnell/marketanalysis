import van from 'vanjs-core'

import { GraphActionModal } from './graphaction.js'
import { PostJSON } from './backend.js'
import userTimeZone from './timezone.js';

const runGraph = ({id, symbol, date, showWarmup}) => {
    return new Promise((resolve, reject) => {
        const route = `/graphs/${id}/run`
        const timeZone = userTimeZone()
        const object = {type: "day", symbol, date, timeZone, showWarmup}
        const options = {accept: 'application/json'}

        console.log("running graph", object)

        PostJSON({route, object, options}).then(resp => {
            if (resp.status != 200) {
                resp.json().then(appErr => {
                    console.log("failed to run graph", appErr);
    
                    reject(appErr);    
                })
            }

            resp.json().then(obj => resolve(obj))
        }).catch(err => {
            console.log("failed to make run graph request", err)
            
            reject({
                title: "Action Failed",
                message: "failed to make run graph request",
                details: [err.message],
            })
        });
    });
}

const RunGraph = (graph) => {
    const runDisabled = van.state(false)
    const doAction = ({date, symbol, showWarmup}) => {
        return runGraph({id: graph.id, symbol, date, showWarmup})
    }
    const inputElems = []

    GraphActionModal({ actionName: "run", graph, inputElems, runDisabled, doAction })
}

export {RunGraph};
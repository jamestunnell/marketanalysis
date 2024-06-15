import { PostJSON } from '../backend.js'

const runGraph = ({runType, graph, symbol, date, numCharts}) => {
    return new Promise((resolve, reject) => {
        const route = `/graphs/run`
        const object = {runType, graph, symbol, date, numCharts}
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

export {runGraph};
import van from "vanjs-core"

import {Get} from './backend.js'

const {div, h2, p} = van.tags

const getGraph = async (id) => {
    console.log(`getting graph ${id}`);

    const resp = await Get(`/graphs/${id}`);

    if (resp.status != 200) {
        console.log("failed to get graph", await resp.json());

        return []
    }

    const d = await resp.json();

    console.log("received graph", d);

    return d;
}

const Graph = (id) => {
    const graphArea = div(
        {class: ""},
    )

    console.log(`opening graph ${id}`)

    getGraph(id).then(graph => {
        if (!graph) {
            return
        }

        let blockCount = 0;
        let connCount = 0;

        if (graph.blocks) {
            blockCount = Object.keys(graph.blocks).length;
        }

        if (graph.connections) {
            connCount = graph.connections.length;
        }
        
        van.add(graphArea,
            h2(graph.name),
            p(`ID: ${graph.id}`),
            p(`Block count: ${blockCount}`),
            p(`Connection count: ${connCount}`),
        );
    });
    
    return graphArea
}

export default Graph;
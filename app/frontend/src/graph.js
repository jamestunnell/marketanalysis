import van from "vanjs-core"

import {Get} from './backend.js'
import {ButtonAct} from './buttons.js'
import {Table, TableRow} from './table.js'

const {div, p, tbody, td, tr} = van.tags

const getGraph = async (id) => {
    console.log(`getting graph ${id}`);

    const resp = await Get(`/graphs/${id}`);

    if (resp.status != 200) {
        console.log("failed to get graph", await resp.json());

        return null;
    }

    const d = await resp.json();

    console.log("received graph", d);

    return d;
}

const BlockTableRow = ({name, type, paramVals, onDelete}) => {
    const deleted = van.state(false)

    // const viewBtn = ButtonAct({
    //     text: "",
    //     onclick: () => routeTo('graphs', [id]),
    // });
    const deleteBtn = ButtonAct({
        text: "",
        onclick: () => {
            onDelete()
            
            deleted.val = true
        },
    });

    // viewBtn.classList.add("fa-regular");
    // viewBtn.classList.add("fa-eye");

    deleteBtn.classList.add("fa-solid");
    deleteBtn.classList.add("fa-trash");

    return () => deleted.val ? null : tr(
        {class: "border border-solid"},
        td({class: "px-6 py-4"}, name),
        td({class: "px-6 py-4"}, type),
        td({class: "px-6 py-4"}, paramVals ? paramVals.stringify(): ""),
        td(
            {class: "px-6 py-4"},
            div({class:"flex flex-row"}, deleteBtn)
        ),
    )
}


const Graph = (id) => {
    const name = van.state("");
    const blockTableBody = tbody({class:"table-auto"});
    const connTableBody = tbody({class:"table-auto"});

    const graphArea = div(
        {class: "p-6 w-full flex flex-col"},
        p({class: "text-2xl font-medium font-bold mb-4"}, name),
        p({class: "text-lg font-medium mb-2"}, "Blocks"),
        Table({columnNames: ["Name", "Type", "Param Vals", ""], tableBody: blockTableBody}),
        p({class: "text-lg font-medium mb-2"}, "Connections"),
        Table({columnNames: ["Source", "Target", ""], tableBody: connTableBody}),
    )

    console.log(`viewing graph ${id}`)

    getGraph(id).then(graph => {
        if (!graph) {
            return
        }

        name.val = graph.name;

        if (graph.blocks) {
            for (var blkName in graph.blocks) {
                const blk = graph.blocks[blkName];
                const row = BlockTableRow({
                    name: blkName,
                    type: blk.type,
                    paramVals: blk.paramVals,
                    onDelete: () => {

                    }
                });

                van.add(blockTableBody, row);
            }
        }
        
        for (var conn in graph.connections) {
            const row = ConnTableRow({
                source: conn.source,
                target: conn.target,
                onDelete: () => {

                }
            });

            van.add(connTableBody, row);
        }

        // if (graph.connections) {
        //     connCount = graph.connections.length;
        // }
    });
    
    return graphArea
}

export default Graph;
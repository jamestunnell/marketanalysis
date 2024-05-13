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
    const deleted = van.state(false);
    // const viewBtn = ButtonAct({
    //     text: "",
    //     onclick: () => routeTo('graphs', [id]),
    // });
    const deleteBtn = ButtonAct({
        text: "",
        onclick: onDelete,
    });

    // viewBtn.classList.add("fa-regular");
    // viewBtn.classList.add("fa-eye");

    deleteBtn.classList.add("fa-solid");
    deleteBtn.classList.add("fa-trash");

    return deleted.val ? null : tr(
        {class: "border border-solid"},
        td({class: "px-6 py-4"}, name),
        td({class: "px-6 py-4"}, type),
        td({class: "px-6 py-4"}, van.derive(() => JSON.stringify(paramVals.val))),
        td(
            {class: "px-6 py-4"},
            div({class:"flex flex-row"}, deleteBtn)
        ),
    )
}


const ConnTableRow = ({source, target, onDelete}) => {
    // const viewBtn = ButtonAct({
    //     text: "",
    //     onclick: () => routeTo('graphs', [id]),
    // });
    const deleteBtn = ButtonAct({
        text: "",
        onclick: onDelete,
    });

    // viewBtn.classList.add("fa-regular");
    // viewBtn.classList.add("fa-eye");

    deleteBtn.classList.add("fa-solid");
    deleteBtn.classList.add("fa-trash");

    return tr(
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
    const graph = {
        id: van.state(id),
        name: van.state(""),
        connections: {},
        blocks: {},
    };
    const blockTableBody = tbody({class:"table-auto"});
    const connTableBody = tbody({class:"table-auto"});

    const addBlockBtn = ButtonAct({
        text: "Add Block",
        onclick: () => {

        },
    })
    
    const addConnBtn = ButtonAct({
        text: "Add Connection",
        onclick: () => {

        },
    })

    const graphArea = div(
        {class: "p-6 w-full flex flex-col"},
        p({class: "text-2xl font-medium font-bold mb-4"}, name),
        div(
            {class: "flex flex-row-reverse p-4"},
            addBlockBtn,
            addConnBtn,
        ),
        p({class: "text-lg font-medium"}, "Blocks"),
        Table({columnNames: ["Name", "Type", "Param Vals", ""], tableBody: blockTableBody}),
        p({class: "text-lg font-medium"}, "Connections"),
        Table({columnNames: ["Source", "Target", ""], tableBody: connTableBody}),
    )

    console.log(`viewing graph ${id}`)

    getGraph(id).then(g => {
        if (!g) {
            return
        }

        graph.name.val = g.name;
        graph.id.val = g.id;

        for (var blkName in g.blocks) {
            const blk = van.state(Object.assign({}, g.blocks[blkName], {name: blkName}));
            const name = van.derive(() => blk.val.name);
            const row = BlockTableRow({
                name: name,
                type: van.derive(() => blk.val.type),
                paramVals: van.derive(() => blk.val.paramVals || {}),
                onDelete: () => {
                    delete graph.blocks[name.val]
                }
            })

            graph.blocks[blkName] = blk
    
            van.add(blockTableBody, row);
        }
        
        for (var conn in g.connections) {
            // conn.key = uuidv4();

            // console.log("adding conn to conns state", conn);
            
            // graph.connections[key] = van.state(conn);
        }
    });
    
    return graphArea
}

export default Graph;
import van from "vanjs-core"
import { v4 as uuidv4 } from 'uuid';

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

const BlockTableRow = ({name, type, onDelete}) => {
    const deleted = van.state(false);
    // const viewBtn = ButtonAct({
    //     text: "",
    //     onclick: () => routeTo('graphs', [id]),
    // });
    const deleteBtn = ButtonAct({
        text: "",
        onclick: () => {
            deleted.val = true;

            onDelete();
        },
    });

    // viewBtn.classList.add("fa-regular");
    // viewBtn.classList.add("fa-eye");

    deleteBtn.classList.add("fa-solid");
    deleteBtn.classList.add("fa-trash");

    const buttons = div({class:"flex flex-row"}, deleteBtn);
    const rowItems = [name, type, buttons]

    return () => deleted.val ? null : TableRow(rowItems);
}

const ConnTableRow = ({source, target, onDelete}) => {
    const deleted = van.state(false);
    // const viewBtn = ButtonAct({
    //     text: "",
    //     onclick: () => routeTo('graphs', [id]),
    // });
    const deleteBtn = ButtonAct({
        text: "",
        onclick: () => {
            deleted.val = true;

            onDelete();
        },
    });

    // viewBtn.classList.add("fa-regular");
    // viewBtn.classList.add("fa-eye");

    deleteBtn.classList.add("fa-solid");
    deleteBtn.classList.add("fa-trash");
    
    const buttons = div({class:"flex flex-row"}, deleteBtn);
    const rowItems = [source, target, buttons];

    return () => deleted.val ? null : TableRow(rowItems);
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
        Table({columnNames: ["Name", "Type", ""], tableBody: blockTableBody}),
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

        const blockRows = Object.keys(g.blocks).map(blkName => {
            const blk = van.state(Object.assign({}, g.blocks[blkName], {name: blkName}));
            const name = van.derive(() => blk.val.name);
            
            graph.blocks[blkName] = blk
            
            return BlockTableRow({
                name: name,
                type: van.derive(() => blk.val.type),
                onDelete: () => delete graph.blocks[name.val],
            });
        });
        
        const connRows = g.connections.map(c => {
            console.log("found graph connection", c);

            const id = uuidv4();
            const conn = van.state(c);
            
            graph.connections[id] = conn;

            return ConnTableRow({
                source: van.derive(() => conn.val.source),
                target: van.derive(() => conn.val.target),
                onDelete: () => delete graph.connections[id],
            })
        });

        console.log("found %d graph connections", connRows.length);

        van.add(blockTableBody, blockRows);
        van.add(connTableBody, connRows);
    });
    
    return graphArea
}

export default Graph;
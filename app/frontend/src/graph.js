import van from "vanjs-core"
import { v4 as uuidv4 } from 'uuid';
import hash from 'object-hash';

import { Get, Put } from './backend.js';
import { DoAppErrorModal} from './apperror.js';
import { DoBlockModal } from "./block.js";
import { ButtonAdd, ButtonDelete, ButtonSave, ButtonView } from './buttons.js';
import { DoConnectionModal } from "./connection.js";
import { Table, TableRow } from './table.js';

const {div, p, tbody} = van.tags

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

const putGraph = async (graph) => {
    console.log("saving graph", graph);

    const resp = await Put({route:`/graphs/${graph.id}`, content: graph});

    if (resp.status != 204) {
        const appErr = await resp.json()
        
        console.log("failed to save graph", appErr);

        return appErr;
    }

    // Avoid Fetch failed loading
    await resp.text();
    
    console.log("saved graph", graph);

    return null;
}

const BlockTableRow = ({name, type, onView, onDelete}) => {
    const deleted = van.state(false);
    const viewBtn = ButtonView(onView);
    const deleteBtn = ButtonDelete(() => {
        deleted.val = true;

        onDelete();
    });

    const buttons = div({class:"flex flex-row"}, viewBtn, deleteBtn);
    const rowItems = [name, type, buttons]

    return () => deleted.val ? null : TableRow(rowItems);
}

const ConnTableRow = ({source, target, onView, onDelete}) => {
    const deleted = van.state(false);
    const viewBtn = ButtonView(onView);
    const deleteBtn = ButtonDelete(() => {
        deleted.val = true;

        onDelete();
    });

    viewBtn.classList.add("fa-regular","fa-eye");
    deleteBtn.classList.add("fa-solid","fa-trash");
    
    const buttons = div({class:"flex flex-row"}, viewBtn, deleteBtn);
    const rowItems = [source, target, buttons];

    return () => deleted.val ? null : TableRow(rowItems);
}

const Graph = (id) => {
    const graph = {
        name: van.state(""),
        connections: {},
        blocks: {},
    };
    const blockTableBody = tbody({class:"table-auto"});
    const connTableBody = tbody({class:"table-auto"});

    const makeBlockRow = (block) => {
        const id = uuidv4();
        const blk = van.state(block);

        graph.blocks[id] = blk;

        return BlockTableRow({
            name: van.derive(() => blk.val.name),
            type: van.derive(() => blk.val.type),
            onDelete: () => delete graph.blocks[id],
            onView: () => {
                DoBlockModal({
                    block: blk.val,
                    handleResult: (block2) => {
                        if (hash(blk.val) === hash(block2)) {
                            console.log('no block change detected');

                            return
                        }

                        console.log('updating block', block2);

                        blk.val = block2;
                    },
                });
            },
        });
    }
    const makeConnRow = (connection) => {
        const id = uuidv4();
        const conn = van.state(connection);

        graph.connections[id] = conn;

        return ConnTableRow({
            source: van.derive(() => conn.val.source),
            target: van.derive(() => conn.val.target),
            onDelete: () => delete graph.connections[id],
            onView: () => {
                DoConnectionModal({
                    connection: conn.val,
                    handleResult: (connection2) => {
                        if (hash(conn.val) === hash(connection2)) {
                            console.log('no connection change detected');

                            return
                        }

                        console.log('updating connection', connection2);

                        conn.val = connection2;
                    },
                });
            },
        });
    }
    
    const addBlockBtn = ButtonAdd(() => {
        DoBlockModal({
            block: {name: "", type: "", paramVals: {}, recording: []},
            handleResult: (b) => {
                van.add(blockTableBody, makeBlockRow(b));
            },
        });
    });
    const addConnBtn = ButtonAdd(() => {
        DoConnectionModal({
            connection: {source: "", target: ""},
            handleResult: (c) => {
                console.log('adding new connection', c);

                van.add(connTableBody, makeConnRow(c));
            },
        });
    });

    const saveBtn = ButtonSave(() => {
        const blocks = Object.keys(graph.blocks).map(id => graph.blocks[id].val)
        const conns = Object.keys(graph.connections).map(id => graph.connections[id].val);
        const g = {
            id: id,
            name: graph.name.val,
            blocks: blocks,
            connections: conns,
        }

        putGraph(g).then(appErr => {
            if (appErr) {
                DoAppErrorModal(appErr);
            }
        });
    });

    const graphArea = div(
        {class: "p-6 w-full flex flex-col"},
        div(
            {class: "flex flex-row p-4"},
            p({class: "text-2xl font-medium font-bold mb-4"}, graph.name),
            saveBtn,
            // exportBtn,
            // importBtn,
        ),
        div(
            {class: "flex flex-row p-4"},
            p({class: "text-lg font-medium"}, "Blocks"),
            addBlockBtn,
        ),
        Table({columnNames: ["Name", "Type", ""], tableBody: blockTableBody}),
        div(
            {class: "flex flex-row p-4"},
            p({class: "text-lg font-medium"}, "Connections"),
            addConnBtn,
        ),
        Table({columnNames: ["Source", "Target", ""], tableBody: connTableBody}),
    )

    console.log(`viewing graph ${id}`)

    getGraph(id).then(g => {
        if (!g) {
            return
        }

        graph.name.val = g.name;

        van.add(blockTableBody, g.blocks.map(b => makeBlockRow(b)));
        van.add(connTableBody, g.connections.map(c => makeConnRow(c)));
    });
    
    return graphArea
}

export default Graph;
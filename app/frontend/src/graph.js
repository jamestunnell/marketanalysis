import van from "vanjs-core"
import {Modal} from "vanjs-ui"
import { v4 as uuidv4 } from 'uuid';
import hash from 'object-hash';

import {Get, Put} from './backend.js';
import {DoAppErrorModal} from './apperror.js';
import {ButtonAct, ButtonCancel} from './buttons.js';
import {Table, TableRow} from './table.js';

const {div, input, label, p, tbody} = van.tags

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
    const viewBtn = ButtonAct({
        text: "",
        onclick: onView,
    });
    const deleteBtn = ButtonAct({
        text: "",
        onclick: () => {
            deleted.val = true;

            onDelete();
        },
    });

    viewBtn.classList.add("fa-regular", "fa-eye");
    deleteBtn.classList.add("fa-solid", "fa-trash");

    const buttons = div({class:"flex flex-row"}, viewBtn, deleteBtn);
    const rowItems = [name, type, buttons]

    return () => deleted.val ? null : TableRow(rowItems);
}

const ConnTableRow = ({source, target, onView, onDelete}) => {
    const deleted = van.state(false);
    const viewBtn = ButtonAct({
        text: "",
        onclick: onView,
    });
    const deleteBtn = ButtonAct({
        text: "",
        onclick: () => {
            deleted.val = true;

            onDelete();
        },
    });

    viewBtn.classList.add("fa-regular","fa-eye");
    deleteBtn.classList.add("fa-solid","fa-trash");
    
    const buttons = div({class:"flex flex-row"}, viewBtn, deleteBtn);
    const rowItems = [source, target, buttons];

    return () => deleted.val ? null : TableRow(rowItems);
}

const BlockForm = ({name, type, onOK, onCancel}) => {
    const inputClass = "block px-5 py-5 mt-2 border border-gray-200 rounded-md focus:border-indigo-500 focus:outline-none focus:ring";

    return div(
        {class: "flex flex-col drop-shadow hover:drop-shadow-lg w-300 rounded-md"},
        p({class: "text-lg font-medium font-bold text-center"}, "Graph Block"),
        div(
            div(
                label({for: "name"}, "Name"),
                input({
                    id: "name",
                    class: inputClass,
                    type: "text",
                    value: name,
                    placeholder: "Non-empty, unique",
                    oninput: e => name.val = e.target.value,
                }),
                label({for: "type"}, "Type"),
                input({
                    id: "type",
                    class: inputClass,
                    type: "text",
                    value: type,
                    placeholder: "Valid block type",
                    oninput: e => type.val = e.target.value,
                }),
            ),
        ),
        div(
            {class:"mt-4 flex justify-center"},
            ButtonCancel({text: "Cancel", onclick: onCancel}),
            ButtonAct({text: "OK", onclick: onOK}),
        ),
    )
}

const ConnectionForm = ({source, target, onOK, onCancel}) => {
    const inputClass = "block px-5 py-5 mt-2 border border-gray-200 rounded-md focus:border-indigo-500 focus:outline-none focus:ring";

    return div(
        {class: "flex flex-col drop-shadow hover:drop-shadow-lg w-300 rounded-md"},
        p({class: "text-lg font-medium font-bold text-center"}, "Graph Connection"),
        div(
            div(
                label({for: "source"}, "Source Address"),
                input({
                    id: "source",
                    class: inputClass,
                    type: "text",
                    value: source,
                    placeholder: "<block.output>",
                    oninput: e => source.val = e.target.value,
                }),
                label({for: "target"}, "Target Address"),
                input({
                    id: "target",
                    class: inputClass,
                    type: "text",
                    value: target,
                    placeholder: "<block.input>",
                    oninput: e => target.val = e.target.value,
                }),
            ),
        ),
        div(
            {class:"mt-4 flex justify-center"},
            ButtonCancel({text: "Cancel", onclick: onCancel}),
            ButtonAct({text: "OK", onclick: onOK}),
        ),
    )
}

const DoBlockModal = ({block, handleResult}) => {
    const closed = van.state(false)

    const name = van.state(block.name);
    const type = van.state(block.type);
    const paramVals = van.state(block.paramVals);
    const recording = van.state(block.recording);

    van.add(
        document.body,
        Modal({closed},
            BlockForm({
                name: name,
                type: type,
                // paramVals: paramVals,
                // recording: recording,
                onOK: ()=> {
                    handleResult({
                        name: name.val,
                        type: type.val,
                        paramVals: paramVals.val,
                        recording: recording.val,
                    });

                    closed.val = true;
                },
                onCancel: () => {
                    closed.val = true;
                }
            }),
        ),
    );
}

const DoConnectionModal = ({connection, handleResult}) => {
    const closed = van.state(false);
    const source = van.state(connection.source);
    const target = van.state(connection.target);

    van.add(
        document.body,
        Modal({closed},
            ConnectionForm({
                source: source,
                target: target,
                onOK: ()=> {
                    console.log("pressed OK")

                    handleResult({source: source.val, target: target.val});

                    console.log("closing modal")

                    closed.val = true;
                },
                onCancel: () => {
                    closed.val = true;
                }
            }),
        ),
    );
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
    
    const addBlockBtn = ButtonAct({
        text: "Add Block",
        onclick: () => {
            DoBlockModal({
                block: {name: "", type: "", paramVals: {}, recording: []},
                handleResult: (b) => {
                    van.add(blockTableBody, makeBlockRow(b));
                },
            });
        },
    });
    const addConnBtn = ButtonAct({
        text: "Add Connection",
        onclick: () => {
            DoConnectionModal({
                connection: {source: "", target: ""},
                handleResult: (c) => {
                    console.log('adding new connection', c);

                    van.add(connTableBody, makeConnRow(c));
                },
            });
        },
    });

    const saveBtn = ButtonAct({
        text: "",
        onclick: () => {
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
        },
    });

    saveBtn.classList.add("fa-solid", "fa-floppy-disk");
    
    const graphArea = div(
        {class: "p-6 w-full flex flex-col"},
        p({class: "text-2xl font-medium font-bold mb-4"}, graph.name),
        div(
            {class: "flex flex-row-reverse p-4"},
            saveBtn,
            addConnBtn,
            addBlockBtn,
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

        van.add(blockTableBody, g.blocks.map(b => makeBlockRow(b)));
        van.add(connTableBody, g.connections.map(c => makeConnRow(c)));
    });
    
    return graphArea
}

export default Graph;
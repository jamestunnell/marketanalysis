import van from "vanjs-core"
import { v4 as uuidv4 } from 'uuid';
import hash from 'object-hash';

import { DoAppErrorModal} from './apperror.js';
import { Get, Put } from './backend.js';
import { DoBlockModal } from "./block.js";
import { Button, ButtonDanger } from "./buttons.js";
import { DoConnectionModal } from "./connection.js";
import { DownloadCSV, DownloadJSON } from "./download.js";
import { IconAdd, IconDelete, IconExport, IconImport, IconPlay, IconSave, IconView } from "./icons.js";
import { RunGraph } from './rungraph.js'
import { Table, TableRow } from './table.js';
import { UploadJSON } from "./upload.js";
import { Modal } from "vanjs-ui";

const {div, input, p, tbody} = van.tags

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

const putGraph = async ({graph, onSuccess, onErr}) => {
    console.log("saving graph", graph);

    try {
        const resp = await Put({route:`/graphs/${graph.id}`, object: graph});

        if (resp.status != 204) {
            const appErr = await resp.json()
            
            console.log("failed to save graph", appErr);

            onErr(appErr);

            return;
        }

        // Avoid Fetch failed loading
        await resp.text();
        
        console.log("saved graph", graph);

        onSuccess();
    } catch(err) {
        console.log("failed to complete fetch", err)
    }
}

const BlockTableRow = ({name, type, onView, onDelete}) => {
    const deleted = van.state(false);
    const viewBtn = Button({child: IconView(), onclick: onView});
    const deleteBtn = ButtonDanger({
        child: IconDelete(),
        onclick: () => {
            deleted.val = true;

            onDelete();
        },
    });

    const buttons = div({class:"flex flex-row"}, viewBtn, deleteBtn);
    const rowItems = [name, type, buttons]

    return () => deleted.val ? null : TableRow(rowItems);
}

const ConnTableRow = ({source, target, onView, onDelete}) => {
    const deleted = van.state(false);
    const viewBtn = Button({child: IconView(), onclick: onView});
    const deleteBtn = ButtonDanger({
        child: IconDelete(),
        onclick: () => {
            deleted.val = true;

            onDelete();
        },
    });
    
    const buttons = div({class:"flex flex-row"}, viewBtn, deleteBtn);
    const rowItems = [source, target, buttons];

    return () => deleted.val ? null : TableRow(rowItems);
}

const Graph = (id) => {
    const graph = {
        name: van.state(""),
        connections: {},
        blocks: {},
        changed: van.state(false),
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
            onDelete: () => {
                delete graph.blocks[id]

                graph.changed.val = true;
            },
            onView: () => {
                DoBlockModal({
                    block: blk.val,
                    handleResult: (block2) => {
                        if (hash(blk.val) === hash(block2)) {
                            return
                        }

                        graph.changed.val = true

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
            onDelete: () => {
                delete graph.connections[id]

                graph.changed.val = true;
            },
            onView: () => {
                DoConnectionModal({
                    connection: conn.val,
                    handleResult: (connection2) => {
                        if (hash(conn.val) === hash(connection2)) {
                            return
                        }

                        graph.changed.val = true

                        conn.val = connection2;
                    },
                });
            },
        });
    }
    const makeGraph = () => {
        const blocks = Object.keys(graph.blocks).map(id => graph.blocks[id].val)
        const conns = Object.keys(graph.connections).map(id => graph.connections[id].val);
        
        return {
            id: id,
            name: graph.name.val,
            blocks: blocks,
            connections: conns,
        }
    }
    const loadGraph = (g) => {
        if (hash(g) === hash(makeGraph())) {
            console.log('loaded graph is identical to existing, ignoring')

            return
        }

        graph.name.val = g.name;

        // clear existing block & connection rows
        while (blockTableBody.firstChild) {
            blockTableBody.removeChild(blockTableBody.firstChild)
        }

        while (connTableBody.firstChild) {
            connTableBody.removeChild(connTableBody.firstChild)
        }

        graph.blocks = {}
        graph.connections = {}

        van.add(blockTableBody, g.blocks.map(b => makeBlockRow(b)));
        van.add(connTableBody, g.connections.map(c => makeConnRow(c)));
        
        graph.changed.val = true
    }

    const addBlockBtn = Button({
        child: IconAdd(),
        onclick: () => {
            DoBlockModal({
                block: {name: "", type: "", paramVals: {}, recording: []},
                handleResult: (b) => {
                    graph.changed.val = true

                    van.add(blockTableBody, makeBlockRow(b));
                },
            });
        },
    });
    const addConnBtn = Button({
        child: IconAdd(),
        onclick: () => {
            DoConnectionModal({
                connection: {source: "", target: ""},
                handleResult: (c) => {
                    graph.changed.val = true

                    van.add(connTableBody, makeConnRow(c));
                },
            });        
        },
    });
    const runBtn = Button({
        child: IconPlay(),
        disabled: van.derive(() => graph.changed.val),
        onclick: () => RunGraph(makeGraph()),
    });
    const saveBtn = Button({
        child: IconSave(),
        disabled: van.derive(() => !graph.changed.val),
        onclick: () => {
            putGraph({
                graph: makeGraph(),
                onErr: (appErr) => DoAppErrorModal(appErr),
                onSuccess: () => graph.changed.val = false,
            });
        },
    });
    const exportBtn = Button({
        child: IconExport(),
        onclick: () => DownloadJSON({obj: makeGraph(), name: graph.name.val}),
    });
    const importBtn = Button({
        child: IconImport(),
        onclick: () => {
            UploadJSON({
                onSuccess: (g) => {
                    console.log("file imported successfully")

                    if (g.id === id) {
                        loadGraph(g)

                        return
                    }
                        
                    DoAppErrorModal({
                        title: "Invalid Input",
                        message: "Graph IDs do not match",
                        details: ["Go to the Graphs page to import as a different graph."],
                    })
                },
                onErr: (appErr) => DoAppErrorModal(appErr),
            })
        }
    });

    const graphArea = div(
        {class: "container p-6 w-full flex flex-col"},
        div(
            {class: "grid grid-cols-2"},
            div(
                {class: "flex flex-row p-4"},
                p({class: "text-2xl font-medium font-bold"}, graph.name),
            ),
            div(
                {class: "flex flex-row-reverse p-4"},
                exportBtn,
                runBtn,
                saveBtn,
                importBtn,
            ),
        ),
        div(
            {class: "flex flex-row p-4"},
            p({class: "text-xl font-medium"}, "Blocks"),
            addBlockBtn,
        ),
        Table({columnNames: ["Name", "Type", ""], tableBody: blockTableBody}),
        div(
            {class: "flex flex-row p-4"},
            p({class: "text-xl font-medium"}, "Connections"),
            addConnBtn,
        ),
        Table({columnNames: ["Source", "Target", ""], tableBody: connTableBody}),
    )

    console.log(`viewing graph ${id}`)

    getGraph(id).then(g => {
        if (!g) {
            return
        }

        loadGraph(g);

        graph.changed.val = false
    });
    
    return graphArea
}

export default Graph;
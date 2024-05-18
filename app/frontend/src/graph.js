import van from "vanjs-core"
import { v4 as uuidv4 } from 'uuid';
import hash from 'object-hash';
import { nanoid } from 'nanoid'
import { Tooltip } from "vanjs-ui"

import { AppErrorAlert} from './apperror.js';
import { Get, Put } from './backend.js';
import { SelectBlockTypeModal, ConfigureBlockModal } from "./block.js";
import { Button, ButtonDanger } from "./buttons.js";
import { DoConnectionModal } from "./connection.js";
import { DownloadJSON } from "./download.js";
import { IconAdd, IconCheck, IconDelete, IconError, IconExport, IconImport, IconPlay, IconSave, IconView } from "./icons.js";
import { RunGraph } from './rungraph.js'
import { Table, TableRow } from './table.js';
import { UploadJSON } from "./upload.js";

const {div, input, p, tbody} = van.tags

const getAllBlockInfo = () => {
    return new Promise((resolve, reject) => {
        console.log("getting block infos");

        Get('/blocks').
            then(resp => {
                if (resp.status != 200) {
                    resp.json().then(appErr => {
                        console.log("failed to get block infos", appErr);

                        reject(appErr)
                    })
                    
                    return
                }

                resp.json().then(data => {
                    console.log(`received block infos`, data.blocks)

                    resolve(data.blocks)
                })
            }).
            catch(err => {
                reject({
                    title: "Action Failed",
                    message: "failed to make get blocks request",
                    details: [err.message],
                })
            })
    })
}

const getGraph = (id) => {
    return new Promise((resolve, reject) => {
        console.log(`getting graph ${id}`);

        Get(`/graphs/${id}`).then(resp => {
            if (resp.status != 200) {
                resp.json().then(appErr => {
                    console.log("failed to get graph", appErr);

                    reject(appErr)
                })
        
                return;
            }

            resp.json().then(data => {
                console.log("received graph", data);

                resolve(data);
            })    
        }).catch(err => {
            reject({
                title: "Action Failed",
                message: "failed to make get graph requst",
                details: [err.message],
            })
        })
    })
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

const BlockTableRow = ({name, type, info, onView, onDelete}) => {
    const infoMissing = van.derive(() => !info)
    const deleted = van.state(false);
    
    const viewBtn = Button({
        child: [
            IconView(),
            Tooltip({text: `Disabled because ${type} block info is missing`, show: infoMissing}),
        ],
        onclick: onView,
        disabled: infoMissing,
    });
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

const ConnTableRow = ({source, target, validateErr, onView, onDelete}) => {
    const deleted = van.state(false);
    const viewBtn = Button({child: IconView(), onclick: onView});
    const deleteBtn = ButtonDanger({
        child: IconDelete(),
        onclick: () => {
            deleted.val = true;

            onDelete();
        },
    });

    const showTooltip = van.state(false)
    const statusBtn = div(
        {
            class: "rounded-md p-3 m-1 relative",
            onmouseenter: () => showTooltip.val = true,
            onmouseleave: () => showTooltip.val = false,
        },
        () => validateErr.val ? IconError() : IconCheck(),
        Tooltip({
            text: van.derive(() => validateErr.val ? `Connection is invalid: ${validateErr.val.message}` : "Connection is valid"),
            show: showTooltip,
        }),
    );

    const buttons = div({class:"flex flex-row"}, viewBtn, deleteBtn, statusBtn);
    const rowItems = [source, target, buttons];

    return () => deleted.val ? null : TableRow(rowItems);
}

const pageContent = ({graph, infoByType}) => {
    const graphID = graph.id
    const name = van.state(graph.name)
    const changed = van.state(false)
    
    const blockTableBody = tbody({class:"table-auto"});
    const connTableBody = tbody({class:"table-auto"});

    const initialBlocksWithIDs = graph.blocks.map(b => Object.assign(b, {id: nanoid()}))
    const initialConnsWithIDs = graph.connections.map(c => Object.assign(c, {id: nanoid()}))

    const blocks = van.state([].concat(initialBlocksWithIDs))
    const connections = van.state([].concat(initialConnsWithIDs))

    const makeBlockRow = (block) => {
        const id = block.id
        const name = van.state(block.name)
        const digest = van.state(hash(block))
        const info = infoByType[block.type]

        return BlockTableRow({
            name: name,
            type: block.type,
            info: info,
            onDelete: () => {
                blocks.val = blocks.val.filter(b => b.id !== id)
                changed.val = true;
            },
            onView: () => {
                ConfigureBlockModal({
                    info,
                    block: blocks.val.find(b => b.id === id),
                    otherNames: blocks.val.map(b => b.name).filter(n => n != name.val),
                    handleResult: (b) => {
                        const resultDigest = hash(b)
                        if (resultDigest === digest.val) {
                            return
                        }

                        blocks.val = blocks.val.map(blk => blk.id === block.id ? b : blk)
                        name.val = b.name
                        digest.val = resultDigest
                        changed.val = true
                    },
                })
            },
        });
    }
    const makeConnRow = (connection) => {
        const id = connection.id
        const source = van.state(connection.source)
        const target = van.state(connection.target)
        const digest = van.state(hash(connection))

        return ConnTableRow({
            source: source,
            target: target,
            validateErr: van.derive(() => {
                const srcParts = source.val.split(".")
                const tgtParts = target.val.split(".")

                if (srcParts.length !== 2) {
                    return new Error(`source ${source.val} not formatted as <A>.<B>`)
                }

                if (tgtParts.length !== 2) {
                    return new Error(`target ${target.val} not formatted as <A>.<B>`)
                }

                const blks = Object.values(blocks.val)

                console.log("working with blocks", blks)

                const src = blks.find(blk => blk.name === srcParts[0])
                const tgt = blks.find(blk => blk.name === tgtParts[0])

                if (!src) {
                    return new Error(`source block ${srcParts[0]} not found`)
                }
                
                const srcInfo = infoByType[src.type]

                console.log(`found source`, src, srcInfo)

                if (!tgt) {
                    return new Error(`target block ${tgtParts[0]} not found`)
                }

                const tgtInfo = infoByType[tgt.type]

                console.log(`found target`, tgt, tgtInfo)
                
                if (!srcInfo.outputs.find(o => o.name === srcParts[1])) {
                    return new Error(`source block ${srcParts[0]} does not have output ${srcParts[1]}`)
                }
                
                if (!tgtInfo.inputs.find(i => i.name === tgtParts[1])) {
                    return new Error(`target block ${tgtParts[0]} does not have input ${tgtParts[1]}`)
                }

                return null
            }),
            onDelete: () => {
                connections.val = connections.val.filter(c => c.id !== id)
                changed.val = true;
            },
            onView: () => {
                DoConnectionModal({
                    connection: {source: source.val, target: target.val},
                    handleResult: (c) => {
                        const resultDigest = hash(c)
                        if (resultDigest === digest.val) {
                            return
                        }

                        connections.val = connections.val.map(conn => conn.id === id ? c : conn)
                        source.val = c.source
                        target.val = c.target
                        digest.val = resultDigest
                        changed.val = true
                    },
                });
            },
        });
    }
    const makeGraph = () => {
        const g = {
            id: graphID,
            name: name.val,
            blocks: blocks.val.map(b => {
                const { id: _, ...blkWithoutID } = b

                return blkWithoutID
            }),
            connections: connections.val.map(c => {
                const { id: _, ...connWithoutID } = c

                return connWithoutID
            }),
        }

        console.log("made graph", g)

        return g
    }
    const importGraph = (g) => {
        if (hash(g) === hash(makeGraph())) {
            console.log('graph is identical to existing, skipping loading')

            return
        }

        name.val = g.name;

        // clear existing block & connection rows
        while (blockTableBody.firstChild) {
            blockTableBody.removeChild(blockTableBody.firstChild)
        }

        while (connTableBody.firstChild) {
            connTableBody.removeChild(connTableBody.firstChild)
        }

        const blocksWithIDs = graph.blocks.map(b => Object.assign(b, {id: nanoid()}))
        const connsWithIDs = graph.connections.map(c => Object.assign(c, {id: nanoid()}))
        
        van.add(blockTableBody, blocksWithIDs.map(makeBlockRow))
        van.add(connTableBody, connsWithIDs.map(makeConnRow))
        
        blocks.val = blocksWithIDs
        connections.val = connsWithIDs
       
        changed.val = true
    }
    const addBlockBtn = Button({
        child: IconAdd(),
        onclick: () => {
            SelectBlockTypeModal({
                types: Object.keys(infoByType),
                handleResult: (selectedType) => {
                    const info = infoByType[selectedType]

                    console.log(`adding ${selectedType} block`, info)

                    ConfigureBlockModal({
                        info: info,
                        block: {type: selectedType, name: nanoid(), paramVals: {}, recording: []},
                        otherNames: Object.values(blocks.val).map(b => b.name),
                        handleResult: (block) => {
                            block.id = nanoid()

                            van.add(blockTableBody, makeBlockRow(block))

                            blocks.val = [].concat(blocks.val, [block])
                            changed.val = true
                        },
                    })
                },
            })
        },
    });
    const addConnBtn = Button({
        child: IconAdd(),
        onclick: () => {
            DoConnectionModal({
                connection: {source: "", target: ""},
                handleResult: (c) => {
                    c.id = nanoid()

                    van.add(connTableBody, makeConnRow(c))
                    
                    connections.val = [].concat(connections.val, [c])
                    changed.val = true
                },
            });        
        },
    });
    const runBtn = Button({
        child: IconPlay(),
        disabled: van.derive(() => changed.val),
        onclick: () => RunGraph(makeGraph()),
    });
    const saveBtn = Button({
        child: IconSave(),
        disabled: van.derive(() => !changed.val),
        onclick: () => {
            putGraph({
                graph: makeGraph(),
                onErr: (appErr) => AppErrorAlert(appErr),
                onSuccess: () => changed.val = false,
            });
        },
    });
    const exportBtn = Button({
        child: IconExport(),
        onclick: () => DownloadJSON({obj: makeGraph(), name: name.val}),
    });
    const importBtn = Button({
        child: IconImport(),
        onclick: () => {
            UploadJSON({
                onSuccess: (g) => {
                    if (g.id === graphID) {
                        importGraph(g)

                        console.log("imported graph file")

                        return
                    }
                    
                    AppErrorAlert({
                        title: "Invalid Input",
                        message: "Graph IDs do not match",
                        details: ["Go to the Graphs page to import as a different graph."],
                    })
                },
                onErr: (appErr) => AppErrorAlert(appErr),
            })
        },
        disabled: van.derive(() => Object.keys(infoByType).length == 0),
    });
    const graphArea = div(
        {class: "container p-6 w-full flex flex-col"},
        div(
            {class: "grid grid-cols-2"},
            div(
                {class: "flex flex-row p-4"},
                p({class: "text-2xl font-medium font-bold"}, name),
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

    van.add(blockTableBody, initialBlocksWithIDs.map(makeBlockRow))
    van.add(connTableBody, initialConnsWithIDs.map(makeConnRow))

    console.log(`viewing graph ${graphID}`)

    return graphArea
}

const GraphPage = (id) => {
    const page = div()
    
    Promise.all([
        getAllBlockInfo(),
        getGraph(id),
    ]).then(values => {
        const allBlockInfo = values[0]
        const graph = values[1]
        const infoByType = Object.fromEntries(allBlockInfo.map(info => [info.type, info]));

        console.log("got all block info", infoByType)
        
        van.add(page, pageContent({id, graph, infoByType}))
    }).catch(appErr => {
        console.log("failed to resolve all promises", appErr)

        AppErrorAlert(appErr)
    })

    return page
}

export default GraphPage;
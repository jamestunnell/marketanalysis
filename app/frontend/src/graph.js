import van from "vanjs-core"
import hash from 'object-hash';
import { nanoid } from 'nanoid'

import { AppErrorAlert} from './apperror.js';
import { Get, Put } from './backend.js';
import { BacktestGraph } from './backtestgraph.js'
import { SelectBlockTypeModal, BlockRow } from "./block.js";
import { ButtonGroup } from './buttongroup.js'
import { ButtonIcon, ButtonIconDisableable } from "./buttons.js";
import { ConnectionRow } from "./connection.js";
import { DownloadJSON } from "./download.js";
import { EvalGraph } from './evalgraph.js'
import { IconAdd, IconExport, IconImport, IconMagnifyDollar, IconPlay, IconSave, IconStethoscope } from "./icons.js";
import { RunGraph } from './rungraph.js'
import { Table } from './table.js';
import truncateString from "./truncatestring.js";
import { UploadJSON } from "./upload.js";

const {div, p, tbody} = van.tags

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

class PageContent {
    constructor({graph, infoByType}) {
        const digest = hash(graph)

        console.log(`initial graph digest`, truncateString(digest, 10))

        this.id = graph.id
        this.digest = van.state(digest)
        this.digestSaved = van.state(digest)
        this.name = van.state(graph.name)
        this.infoByType = infoByType

        this.blockTableBody = tbody({class:"table-auto"});
        this.connTableBody = tbody({class:"table-auto"});

        this.redoBlockRows(graph.blocks)
        this.redoConnRows(graph.connections)
    }

    redoBlockRows(blocks) {
        this.blockRowsByID = Object.fromEntries(blocks.map(b => {
            const id = nanoid()
            const row = new BlockRow({
                id,
                block: b,
                info: this.infoByType[b.type],
                parent: this,
            })

            return [id, row]
        }))

        van.add(this.blockTableBody, Object.values(this.blockRowsByID).map(row => row.render()))
    }

    redoConnRows(connections) {
        this.connRowsByID = Object.fromEntries(connections.map(c => {
            const id = nanoid()
            const row = new ConnectionRow({id, connection: c, parent: this})

            return [id, row]
        }))


        van.add(this.connTableBody, Object.values(this.connRowsByID).map(row => row.render()))
    }

    deleteBlockRow(id) {
        delete this.blockRowsByID[id]

        this.updateDigest()
    }

    deleteConnectionRow(id) {
        delete this.connRowsByID[id]

        this.updateDigest()
    }

    blockRowsWithoutID(tgtID) {
        return Object.entries(this.blockRowsByID).filter(([id,row]) => id !== tgtID).map(([id, row]) => row)
    }

    blockNames() {
        const names = Object.values(this.blockRowsByID).map(row => row.getName())
        
        console.log("getting block names", names)
        
        return names
    }

    onBlockNameChange() {
        Object.values(this.connRowsByID).forEach(row => row.onBlockNameChange())
    }

    findBlockInfo(name) {
        const row = Object.values(this.blockRowsByID).find(row => name === row.getName())
        
        return row ? row.info : null
    }

    render() {
        console.log("rendering graph page content")
        
        const needsSaved = van.derive(() => this.digest.val !== this.digestSaved.val)
        const addIcon1 = IconAdd()
        const addIcon2 = IconAdd()
        
        addIcon1.classList.add("text-xl")
        addIcon2.classList.add("text-xl")

        const addBlockBtn = ButtonIcon({
            icon: addIcon1,
            onclick: () => {
                SelectBlockTypeModal({
                    types: Object.keys(this.infoByType),
                    handleResult: (selectedType) => {
                        const id  = nanoid()
                        const info = this.infoByType[selectedType]
                        const block = {type: selectedType, name: "", paramVals: {}, recording: []}
                        const row = new BlockRow({id, block, info, parent: this})
                        
                        console.log(`adding ${selectedType} block`, info)
                        
                        van.add(this.blockTableBody, row.render())

                        this.blockRowsByID[id] = row
                        this.updateDigest()
                    },
                })
            },
        });
        const addConnBtn = ButtonIcon({
            icon: addIcon2,
            onclick: () => {
                const connection = {source: "", target: ""}
                const id = nanoid()
                const row = new ConnectionRow({id, connection, parent: this})

                van.add(this.connTableBody, row.render())

                this.connRowsByID[id] = row
                this.updateDigest()
            },
        });
        const runBtn = ButtonIconDisableable({
            icon: IconPlay(),
            disabled: needsSaved,
            onclick: () => RunGraph(this.makeGraph()),
        });
        const evalBtn = ButtonIconDisableable({
            icon: IconStethoscope(),
            disabled: needsSaved,
            onclick: () => EvalGraph(this.makeGraph(), this.infoByType),
        });
        const backtestBtn = ButtonIconDisableable({
            icon: IconMagnifyDollar(),
            disabled: needsSaved,
            onclick: () => BacktestGraph(this.makeGraph(), this.infoByType),
        });
        const saveBtn = ButtonIconDisableable({
            icon: IconSave(),
            disabled: van.derive(() => !needsSaved.val),
            onclick: () => {
                const graph = this.makeGraph()

                putGraph({
                    graph: graph,
                    onErr: (appErr) => AppErrorAlert(appErr),
                    onSuccess: () => {
                        const digest = hash(graph)
                        
                        this.digest.val = digest
                        this.digestSaved.val = digest
                    }
                });
            },
        });
        const exportBtn = ButtonIcon({
            icon: IconExport(),
            onclick: () => DownloadJSON({obj: this.makeGraph(), basename: this.name.val}),
        });
        const importBtn = ButtonIcon({
            icon: IconImport(),
            onclick: () => {
                UploadJSON({
                    onSuccess: (g) => {
                        if (g.id === this.id) {
                            this.importGraph(g)
    
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
            disabled: van.derive(() => Object.keys(this.infoByType).length == 0),
        });

        return div(
            {class: "container p-4 w-full flex flex-col divide-y divide-gray-400"},
            div(
                {class: "grid grid-cols-2"},
                div(
                    {class: "flex flex-row p-3"},
                    p({class: "text-2xl font-medium font-bold"}, this.name),
                ),
                ButtonGroup({
                    buttons: [importBtn, saveBtn, runBtn, evalBtn, backtestBtn, exportBtn],
                    moreClass: "place-self-end",
                })
            ),
            div(
                {class: "flex flex-col mt-4"},
                div(
                    {class: "grid grid-cols-2"},
                    div(
                        {class: "flex flex-row p-2"},
                        p({class: "p-3 m-1 text-xl font-medium"}, "Blocks"),
                    ),
                    div(
                        {class: "flex flex-row-reverse p-2"},
                        addBlockBtn,
                    )
                ),
                Table({columnNames: ["Name", "Type", "Parameters", "Recording", "", ""], tableBody: this.blockTableBody}),
            ),
            div(
                {class: "flex flex-col mt-4"},
                div(
                    {class: "grid grid-cols-2"},
                    div(
                        {class: "flex flex-row p-2"},
                        p({class: "p-3 m-1 text-xl font-medium"}, "Connections"),
                    ),
                    div(
                        {class: "flex flex-row-reverse p-2"},
                        addConnBtn,
                    )
                ),
                Table({columnNames: ["Block", "Output", "Block", "Input", "", ""], tableBody: this.connTableBody}),
            )
        )
    }

    updateDigest() {
        const digest = hash(this.makeGraph())

        console.log(`updated graph digest`, truncateString(digest, 10))

        this.digest.val = digest
    }

    makeGraph() {
        const graph = {
            id: this.id,
            name: this.name.val,
            blocks: Object.values(this.blockRowsByID).map(row => row.makeBlock()),
            connections: Object.values(this.connRowsByID).map(row => row.makeConnection()),
        }

        console.log("made graph", graph)

        return graph
    }

    importGraph(graph) {
        if (hash(graph) === hash(this.makeGraph())) {
            console.log('graph is identical to existing, skipping loading')

            return
        }

        Object.entries(this.blockRowsByID).forEach(([row, id]) => {
            delete blockRowsByID[id]

            row.delete()
        })

        Object.entries(this.connRowsByID).forEach(([row, id]) => {
            delete connRowsByID[id]
            
            row.delete()
        })

        this.name.val = graph.name;

        this.redoBlockRows(graph.blocks)
        this.redoConnRows(graph.connections)
        this.updateDigest()
    }
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

        van.add(page, new PageContent({graph, infoByType}).render())
    }).catch(appErr => {
        console.log("failed to resolve all promises", appErr)

        AppErrorAlert(appErr)
    })

    return page
}

export default GraphPage;
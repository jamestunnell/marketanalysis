import van from "vanjs-core"
import { v4 as uuidv4 } from 'uuid';
import hash from 'object-hash';
import { nanoid } from 'nanoid'

import { AppErrorAlert} from './apperror.js';
import { Get, Put } from './backend.js';
import { SelectBlockTypeModal, ConfigureBlockModal, BlockRow } from "./block.js";
import { ButtonIcon, ButtonIconDisableable, ButtonToggle } from "./buttons.js";
import { ConnectionRow, DoConnectionModal, validateConnection } from "./connection.js";
import { DownloadJSON } from "./download.js";
import { IconAdd, IconCheck, IconDelete, IconError, IconExport, IconImport, IconPlay, IconSave, IconView } from "./icons.js";
import { RunGraph } from './rungraph.js'
import { Table } from './table.js';
import { UploadJSON } from "./upload.js";

const {canvas, div, input, p, tbody} = van.tags

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
        this.id = graph.id
        this.digest = hash(graph)
        this.name = van.state(graph.name)
        this.changed = van.state(false)
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

        this.changed.val = true
    }

    deleteConnectionRow(id) {
        delete this.connRowsByID[id]

        this.changed.val = true
    }

    blockRowsWithoutID(tgtID) {
        return Object.entries(this.blockRowsByID).filter(([id,row]) => id !== tgtID).map(([id, row]) => row)
    }

    blockNames() {
        console.log("getting block names", this.blockRowsByID)

        return Object.values(this.blockRowsByID).map(row => row.getName())
    }

    findBlockInfo(name) {
        const row = Object.values(this.blockRowsByID).find(row => name === row.getName())
        
        return row ? row.info : null
    }

    render() {
        console.log("rendering graph page content")
        
        const needsSaved = van.derive(() => { 
            if (!this.changed.val) {
                return false
            }

            return this.digest !== hash(this.makeGraph())
        })

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
                        this.changed.val = true
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
                this.changed.val = true
            },
        });
        const toggleVisualBtn = ButtonToggle({
            setVal: false,
            onSet: () => console.log("switching to visual editing"),
            onClear: () => console.log("switching away from visual editing"),
        })
        const runBtn = ButtonIconDisableable({
            icon: IconPlay(),
            disabled: needsSaved,
            onclick: () => RunGraph(this.makeGraph()),
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
                        this.changed.val = false
                        this.digest = hash(graph)
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
                div(
                    {class: "flex flex-row-reverse p-2"},
                    exportBtn,
                    runBtn,
                    saveBtn,
                    importBtn,
                ),
                div(
                    {class: "flex flex-row"},
                ),
                div(
                    {class: "flex flex-row-reverse"},
                    toggleVisualBtn,
                    p({class: "p-2"}, "Visual Editing")
                ),
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
                Table({columnNames: ["","Name", "Type", "Parameters", "Recording", "", ""], tableBody: this.blockTableBody}),
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
                Table({columnNames: ["Source", "Target", "", ""], tableBody: this.connTableBody}),
            )
        )
    }

    markChanged() {
        this.changed.val = true
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
       
        this.changed.val = true
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
import van from "vanjs-core"
import hash from 'object-hash';
import { nanoid } from 'nanoid'

import { AppErrorAlert} from './apperror.js';
import { Get, PutJSON } from './backend.js';
import { BacktestGraph } from './backtestgraph.js'
import { AddBlockModal, BlockItem } from "./block.js";
import { ButtonGroup } from './buttongroup.js'
import { ButtonIcon, ButtonIconDisableable } from "./buttons.js";
import { ConnectionRow } from "./connection.js";
import { DownloadJSON } from "./download.js";
import { EvalGraph } from './evalgraph.js'
import GraphSettings from "./graphsettings.js"
import { IconAdd, IconExport, IconImport, IconMagnifyDollar, IconPlay, IconSave, IconStethoscope } from "./icons.js";
import { PlotRecording } from './plot.js'
import { runDay } from './rungraph.js'
import { Table } from './table.js';
import { truncateStringAddElipses } from "./truncatestring.js";
import { UploadJSON } from "./upload.js";

const {div, label, p, nav} = van.tags

const GRAPH_NAV_ID = "graphNav"

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
        const resp = await PutJSON({route:`/graphs/${graph.id}`, object: graph});

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

        console.log(`initial graph digest`, truncateStringAddElipses(digest, 10))

        this.settings = new GraphSettings({containerID: `#${GRAPH_NAV_ID}`})
        this.id = graph.id
        this.digest = van.state(digest)
        this.digestSaved = van.state(digest)
        this.name = van.state(graph.name)
        this.infoByType = infoByType

        this.blockItems = graph.blocks.map(block => {
            return new BlockItem({
                id: nanoid(),
                block,
                info: this.infoByType[block.type],
                parent: this,
            })
        })

        this.settings.load()

        this.settingsArea = nav(
            {id: GRAPH_NAV_ID, class: "nav bg-gray-400 text-white"},
            div(
                {class: "grid grid-cols-5 items-center"},
                div(
                    {class:"flex flex-row"},
                    p({class: "pl-4 text-lg font-semibold"}, this.name),
                    this.renderExportButton(),
                ),
                div(
                    {class: "col-span-4 flex flex-row-reverse pt-2 pb-2 space-x-2 items-center"},
                    div({class:"pr-4"}),
                    this.settings.dateInput,
                    label({for: "date", class: "pr-2 font-semibold"}, "Date"),
                    this.settings.symbolInput,
                    label({for: "symbol", class: "pr-2 font-semibold"}, "Symbol"),
                ),
            )
        ),
        this.blockButtonsArea = div({class: "flex flex-row bg-gray-400"})
        this.plotArea = div({class: "flex flex-col ml-4 mr-8 max-w-max"})
        this.page = div(
            {class: "flex flex-col"},
            this.settingsArea,
            this.blockButtonsArea,
            this.plotArea,
        )
        this.recording = null
        this.connections = graph.connections

        // Create a new ResizeObserver instance
        const pageResizeObserver = new ResizeObserver(entries => {
            // Loop over the entries
            for (let entry of entries) {
                // Get the element's content rect
                const {width, height} = entry.contentRect;
                // Log its width and height
                console.log(`Element's width: ${width}, height: ${height}`);
            }
        });
        
        // Start observing the div element
        pageResizeObserver.observe(this.page);
        
        this.rebuildBlockButtonsArea()

        van.derive(() => {
            const date = this.settings.date.val
            const symbol = this.settings.symbol.val
            const digest = this.digest.val
            
            if (date.length === 0 || symbol.length === 0) {
                return
            }

            console.log(`running graph with digest ${digest}`)

            this.runDayAndRebuildPlotArea(this.makeGraph())
        })
    }

    deleteBlock(id) {
        this.blockItems = this.blockItems.filter(item => item.id !== id)
        
        this.rebuildBlockButtonsArea()
        this.updateDigest()
    }

    // deleteConnectionRow(id) {
    //     delete this.connRowsByID[id]

    //     this.updateDigest()
    // }

    blockNames() {
        return this.blockItems.map(item => item.getName())
    }

    // onBlockNameChange() {
    //     Object.values(this.connRowsByID).forEach(row => row.onBlockNameChange())
    // }

    // findBlockInfo(name) {
    //     const row = Object.values(this.blockRowsByID).find(row => name === row.getName())
        
    //     return row ? row.info : null
    // }

    renderAddBlockButton() {
        const addIcon = IconAdd()
        
        addIcon.classList.add("text-white", "text-xl")

        return ButtonIcon({
            icon: addIcon,
            // text: "Add",
            onclick: () => {
                AddBlockModal({
                    infoByType: this.infoByType,
                    blockNames: this.blockNames(),
                    handleResult: ({block, info}) => {
                        const item = new BlockItem({id: nanoid(), block, info, parent: this})
                        
                        this.blockItems.push(item)

                        console.log(`adding ${info.type} block`, info)
                        
                        this.rebuildBlockButtonsArea()
                        this.updateDigest()
                    },
                })
            },
        });
    }

    renderViewBlockButton(name) {
        const addIcon = IconAdd()
        
        addIcon.classList.add("text-white", "text-xl")

        return ButtonIcon({
            icon: addIcon,
            // text: "Add",
            onclick: () => {
                AddBlockModal({
                    infoByType: this.infoByType,
                    blockNames: this.blockNames(),
                    handleResult: ({block, info}) => {
                        const item = new BlockItem({id: nanoid(), block, info, parent: this})
                        
                        this.blockItems.push(item)

                        console.log(`adding ${info.type} block`, info)
                        
                        this.rebuildBlockButtonsArea()
                        this.updateDigest()
                    },
                })
            },
        });
    }

    renderExportButton() {
        const icon = IconExport()

        icon.classList.add("text-white")

        return ButtonIcon({
            icon,
            // text: "Export",
            onclick: () => DownloadJSON({
                object: this.makeGraph(),
                filename: this.name.val + ".json",
            }),
        });
    }

    rebuildBlockButtonsArea() {
        while (this.blockButtonsArea.firstChild) {
            this.blockButtonsArea.removeChild(this.blockButtonsArea.firstChild)
        }

        const btns = this.blockItems.map(item => item.renderButton())

        btns.push(this.renderAddBlockButton())

        van.add(this.blockButtonsArea, btns)
    }

    runDayAndRebuildPlotArea(graph) {
        while (this.plotArea.firstChild) {
            this.plotArea.removeChild(this.plotArea.firstChild)
        }

        runDay({
            graph,
            date: this.settings.date.val,
            symbol: this.settings.symbol.val,
        }).then(obj => {
            console.log(`run day succeeded`)
            
            const plots = PlotRecording({
                recording: obj,
                totalHeight: 0.7 * (window.innerHeight || document.documentElement.clientHeight || document.body.clientHeight),
            })

            van.add(this.plotArea, plots.map(plot => {
                return div({class:"mt-8"}, plot)
            }))
        }).catch(appErr => {
            AppErrorAlert(appErr)
        })
    }

    render() { return this.page }

    render2() {
        const needsSaved = van.derive(() => this.digest.val !== this.digestSaved.val)
        const disableRun = van.derive(() => {
            return needsSaved.val || this.settings.symbol.val === '' || this.settings.date.val === ''
        })
        const addIcon1 = IconAdd()
        const addIcon2 = IconAdd()
        
        addIcon1.classList.add("text-xl")
        addIcon2.classList.add("text-xl")

        // const addBlockBtn = ButtonIcon({
        //     icon: addIcon1,
        //     // text: "Add",
        //     onclick: () => {
        //         AddBlockModal({
        //             infoByType: this.infoByType,
        //             blockNames: this.blockNames(),
        //             handleResult: ({block, info}) => {
        //                 const id  = nanoid()
        //                 const row = new BlockRow({id, block, info, parent: this})
                        
        //                 console.log(`adding ${info.type} block`, info)
                        
        //                 van.add(this.blockTableBody, row.render())

        //                 this.blockRowsByID[id] = row
        //                 this.updateDigest()
        //             },
        //         })
        //     },
        // });
        const addConnBtn = ButtonIcon({
            icon: addIcon2,
            // text: "Add",
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
            // text: "Run",
            disabled: disableRun,
            onclick: () => {
                RunGraph({
                    graph: this.makeGraph(),
                    settings: this.settings,
                })
            },
        });
        const evalBtn = ButtonIconDisableable({
            icon: IconStethoscope(),
            // text: "Eval",
            disabled: disableRun,
            onclick: () => EvalGraph({
                graph: this.makeGraph(),
                infoByType: this.infoByType,
                settings: this.settings,
            }),
        });
        const backtestBtn = ButtonIconDisableable({
            icon: IconMagnifyDollar(),
            // text: "$$$",
            disabled: disableRun,
            onclick: () => BacktestGraph({
                graph: this.makeGraph(),
                infoByType: this.infoByType,
                settings: this.settings,
            }),
        });
        const saveBtn = ButtonIconDisableable({
            icon: IconSave(),
            // text: "Save",
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

        // const importBtn = ButtonIcon({
        //     icon: IconImport(),
        //     // text: "Import",
        //     onclick: () => {
        //         UploadJSON({
        //             onSuccess: (g) => {
        //                 if (g.id === this.id) {
        //                     this.importGraph(g)
    
        //                     console.log("imported graph file")
    
        //                     return
        //                 }
                        
        //                 AppErrorAlert({
        //                     title: "Invalid Input",
        //                     message: "Graph IDs do not match",
        //                     details: ["Go to the Graphs page to import as a different graph."],
        //                 })
        //             },
        //             onErr: (appErr) => AppErrorAlert(appErr),
        //         })
        //     },
        //     disabled: van.derive(() => Object.keys(this.infoByType).length == 0),
        // });

        return div(
            {class: "container p-4 w-full flex flex-col divide-y divide-gray-400"},
            div(
                {class: "grid grid-cols-2"},
                div(
                    {class: "flex flex-row p-3"},
                    p({class: "text-2xl font-medium font-bold"}, this.name),
                ),
                ButtonGroup({
                    buttons: [saveBtn, runBtn, evalBtn, backtestBtn, exportBtn],
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
        const graph = this.makeGraph()
        const digest = hash(graph)

        console.log(`updated graph digest`, truncateStringAddElipses(digest, 10))

        if (this.digest.val !== digest) {
            this.digest.val = digest
        }
    }

    makeGraph() {
        const graph = {
            id: this.id,
            name: this.name.val,
            blocks: Object.values(this.blockItems).map(item => item.makeBlock()),
            connections: this.connections,//Object.values(this.connRowsByID).map(row => row.makeConnection()),
        }

        console.log("made graph", graph)

        return graph
    }

    // importGraph(graph) {
    //     if (hash(graph) === hash(this.makeGraph())) {
    //         console.log('graph is identical to existing, skipping loading')

    //         return
    //     }

    //     Object.entries(this.connRowsByID).forEach(([row, id]) => {
    //         delete connRowsByID[id]
            
    //         row.delete()
    //     })

    //     this.name.val = graph.name;

    //     this.redoBlockRows(graph.blocks)
    //     this.redoConnRows(graph.connections)
    //     this.updateDigest()
    // }
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
import van from "vanjs-core"
import hash from 'object-hash';
import { nanoid } from 'nanoid'

import { AppErrorAlert} from '../apperror.js';
import { Get, PutJSON } from '../backend.js';
import { BlockItem } from "./block.js";
import { AddBlockModal } from "./addblock.js";
import { Button, ButtonIcon, ButtonIconDisableable } from "../buttons.js";
import { DownloadJSON } from "../download.js";
import GraphSettings from "./graphsettings.js"
import { IconAdd, IconExport, IconSave } from "../icons.js";
import { PlotRecording } from '../plot.js'
import { runDay } from './rungraph.js'
import { truncateStringAddElipses } from "../truncatestring.js";

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
        this.digestsAgree = van.derive(() => this.digest.val === this.digestSaved.val)
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
                    {class:"flex flex-row space-x-3"},
                    p({class: "pl-4 text-lg font-semibold"}, this.name),
                    this.renderExportButton(),
                    this.renderSaveButton(),
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

    getPossibleSources(targetBlockName) {
        const sources = []

        this.blockItems.forEach(item => {
            const blockName = item.getName()
            if (blockName === targetBlockName) { return }

            item.info.outputs.forEach(output => sources.push(blockName +"." + output.name))
        })

        return sources
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

    renderSaveButton() {
        return Button({
            child: IconSave(),
            disabled: this.digestsAgree,
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
    }

    renderExportButton() {
        return Button({
            child: IconExport(),
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
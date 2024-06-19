import van from "vanjs-core"
import hash from 'object-hash';
import { nanoid } from 'nanoid'

import { AppErrorModal} from '../apperror.js'
import { Get, PutJSON } from '../backend.js'
import { BlockItem } from './block.js'
import { AddBlockModal } from './addblock.js'
import { Button, ButtonIcon } from '../elements/buttons.js'
import { DownloadJSON } from "../download.js"
import GraphSettings from './graphsettings.js'
import { IconAdd, IconExport, IconSave } from '../elements/icons.js';
import { MakeEmptyChart, UpdateCharts } from '../charts.js'
import { runGraph } from './rungraph.js'
import Select from '../elements/select.js'
import { truncateStringAddElipses } from '../truncatestring.js'

const {div, label, nav, option, p, select} = van.tags

const GRAPH_NAV_ID = "graphNav"
const RUN_SINGLE_DAY = "singleDay"
const RUN_MULTI_DAY = "multiDay"
const RUN_MULTI_DAY_SUMMARY = "multiDaySummary"

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
        this.runType = van.state(RUN_SINGLE_DAY)

        this.blockItems = graph.blocks.map(block => {
            return new BlockItem({
                id: nanoid(),
                config: block,
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
                    this.renderSelectRunType(),
                    label({for: "runType", class: "pr-2 font-semibold"}, "Run Type"),
                    this.settings.numChartsInput,
                    label({for: "numCharts", class: "pr-2 font-semibold"}, "Charts"),
                    this.settings.dateInput,
                    label({for: "date", class: "pr-2 font-semibold"}, "Date"),
                    this.settings.symbolInput,
                    label({for: "symbol", class: "pr-2 font-semibold"}, "Symbol"),
                ),
            )
        ),
        this.blockButtonsArea = div({class: "flex flex-row bg-gray-400"})
        this.charts = []
        this.chartsArea = div({class: "grid grid-cols-2flex flex-col ml-4 mr-8 max-w-max"})
        this.page = div(
            {class: "flex flex-col"},
            this.settingsArea,
            this.blockButtonsArea,
            this.chartsArea,
        )

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
            const numCharts = this.settings.numCharts.val
            const runType = this.runType.val

            if (date.length === 0 || symbol.length === 0) {
                return
            }

            if (numCharts != this.charts.length) {
                this.rebuildChartsArea()
            }

            this.runAndUpdateCharts(runType)
        })
    }

    rebuildChartsArea() {
        const numCharts = this.settings.numCharts.val

        console.log(`building charts area`, {numCharts})

        while (this.chartsArea.firstChild) {
            this.chartsArea.removeChild(this.chartsArea.firstChild)
        }

        const totalHeight = 0.7 * (window.innerHeight || document.documentElement.clientHeight || document.body.clientHeight)
        const height = totalHeight / Math.min(3, numCharts)

        this.charts = new Array(numCharts)
        const chartDivs = new Array(numCharts)

        for (let i = 0; i < numCharts; i++) {
            chartDivs[i] = div({class:"mt-8"})
            this.charts[i] = MakeEmptyChart({container: chartDivs[i], height})
        }

        van.add(this.chartsArea, chartDivs)
    }
    
    runAndUpdateCharts(runType) {
        // console.log(`running graph with digest ${this.digest.val}`)

        const opts = {
            runType,
            graph: this.makeGraph(),
            date: this.settings.date.val,
            symbol: this.settings.symbol.val,
            numCharts: this.settings.numCharts.val,
        }

        runGraph(opts).then(obj => {
            console.log(`run ${runType} succeeded`)
            
            UpdateCharts({charts: this.charts, recording: obj})
        }).catch(appErr => {
            AppErrorModal(appErr)
        })
    }

    deleteBlock(id) {
        this.blockItems = this.blockItems.filter(item => item.id !== id)
        
        this.rebuildBlockButtonsArea()
        this.updateDigest()
    }

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

        return sources.sort()
    }

    renderSelectRunType() {
        const options = [
            option({value: RUN_SINGLE_DAY, selected: true}, RUN_SINGLE_DAY),
            option({value: RUN_MULTI_DAY}, RUN_MULTI_DAY),
            option({value: RUN_MULTI_DAY_SUMMARY}, RUN_MULTI_DAY_SUMMARY),
        ]

        return Select({
            options, 
            id: "runType",
            onchange: (e) => this.runType.val = e.target.value,
        })
    }

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
                    handleResult: ({config, info}) => {
                        const item = new BlockItem({id: nanoid(), config, info, parent: this})
                        
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
                    onErr: (appErr) => AppErrorModal(appErr),
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
            blocks: this.blockItems.map(item => item.makeConfig()),
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

        AppErrorModal(appErr)
    })

    return page
}

export default GraphPage;
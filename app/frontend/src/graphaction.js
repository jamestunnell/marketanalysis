import van from 'vanjs-core'
import Datepicker from 'flowbite-datepicker/Datepicker'

import { AppErrorAlert} from './apperror.js'
import { Button, ButtonIcon } from "./buttons.js"
import { ButtonGroupHideable } from './buttongroup.js'
import capitalize from './capitalize.js'
import { DownloadJSON } from "./download.js"
import { IconDownload, IconClose, IconPlay, IconPlot } from './icons.js'
import { INPUT_CLASS } from './input.js'
import { ModalBackground } from './modal.js'
import { PlotRecordingModal } from './plot.js'

const {div, input, label, p} = van.tags

const GraphActionModal = ({actionName, graph, settings, inputElems, runDisabled, doAction}) => {
    const closed = van.state(false)
    const completed = van.state(false)
    const showWarmup = van.state(false)
    const recording = van.state({})

    const closeBtn = ButtonIcon({
        icon: IconClose(),
        // text: "Close",
        onclick: ()=> closed.val = true},
    )

    closeBtn.classList.add("self-end")
    
    const runBtn = Button({
        disabled: van.derive(() => {
            return runDisabled.val || (settings.date.val.length === 0) || (settings.symbol.val.length === 0)
        }),
        child: [IconPlay, ` ${capitalize(actionName)}`],
        onclick: () => {
            completed.val = false
    
            doAction({
                date: settings.date.val,
                symbol: settings.symbol.val,
                showWarmup: showWarmup.val,
            }).then(obj => {
                console.log(`graph action ${actionName} succeeded`)
                
                recording.val = obj
                completed.val = true
            }).catch(appErr => {
                AppErrorAlert(appErr)
            })
        },
    })
    const plotBtn = ButtonIcon({
        icon: IconPlot(),
        // text: "Plot",
        onclick: ()=> PlotRecordingModal(recording.val),
    })
    const downloadBtn = ButtonIcon({
        icon: IconDownload(),
        // text: "Download",
        onclick: ()=> {
            DownloadJSON({
                filename: `${graph.name}_${settings.symbol.val}_${settings.date.val}_${actionName}.json`,
                object: recording.val,
            })
        },
    })
    const buttons = div(
        { class: "grid grid-cols-2" },
        runBtn,
        ButtonGroupHideable({
            buttons: [plotBtn, downloadBtn],
            hidden: van.derive(() => !completed.val),
        })
    )
    const modal = ModalBackground(
        div(
            {id: "foreground", class: "flex flex-col space-y-3 block p-16 rounded-lg bg-white min-w-[25%] max-w-[50%]"},
            closeBtn,
            p({class: "text-lg font-medium font-bold text-center"}, `${capitalize(actionName)} Graph`),
            div(
                {class: "flex flew-row space-x-3"},
                input({
                    id: "showWarmup",
                    class: INPUT_CLASS,
                    type: "checkbox",
                    checked: false,
                    onchange: e => showWarmup.val = e.target.checked,
                }),    
                label({for: "showWarmup"}, "Show Warmup"),
            ),
            ...inputElems,
            buttons,
        )
    )
  
    van.add(document.body, () => closed.val ? null : modal);
}

export {GraphActionModal};
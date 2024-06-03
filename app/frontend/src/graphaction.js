import van from 'vanjs-core'
import Datepicker from 'flowbite-datepicker/Datepicker'

import { AppErrorAlert} from './apperror.js'
import { Button, ButtonIcon } from "./buttons.js"
import { ButtonGroupHideable } from './buttongroup.js'
import capitalize from './capitalize.js'
import { DownloadJSON } from "./download.js"
import { IconDownload, IconClose, IconPlay, IconPlot } from './icons.js'
import { ModalBackground } from './modal.js'
import { PlotRecordingModal } from './plot.js'

const {div, input, label, p} = van.tags

const INPUT_CLASS = "block px-3 py-3 border border-gray-200 rounded-md focus:border-indigo-500 focus:outline-none focus:ring"

const GraphActionModal = ({actionName, graph, inputElems, runDisabled, doAction}) => {
    const closed = van.state(false)
    const completed = van.state(false)
    const showWarmup = van.state(false)
    const recording = van.state({})
    const symbol = van.state("")
    const date = van.state("")

    const dateInput = input({
        id: "actionDate",
        class: INPUT_CLASS,
        type: "text",
        placeholder: 'Select date',
    });

    dateInput.addEventListener('changeDate', (e) => {
        date.val = dateInput.value
    })

    const closeBtn = ButtonIcon({
        icon: IconClose(),
        text: "Close",
        onclick: ()=> closed.val = true},
    )
    const runBtn = Button({
        disabled: van.derive(() => {
            return runDisabled.val || (date.val.length === 0) || (symbol.val.length === 0)
        }),
        child: [IconPlay, ` ${capitalize(actionName)}`],
        onclick: () => {
            completed.val = false
    
            doAction({date: date.val, symbol: symbol.val, showWarmup: showWarmup.val}).then(obj => {
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
        text: "Plot",
        onclick: ()=> PlotRecordingModal(recording.val),
    })
    const downloadBtn = ButtonIcon({
        icon: IconDownload(),
        text: "Download",
        onclick: ()=> {
            DownloadJSON({
                filename: `${graph.name}_${symbol.val}_${date.val}_${actionName}.json`,
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

    closeBtn.classList.add("self-end")

    const modal = ModalBackground(
        div(
            {id: "foreground", class: "flex flex-col space-y-3 block p-16 rounded-lg bg-white min-w-[25%] max-w-[50%]"},
            closeBtn,
            p({class: "text-lg font-medium font-bold text-center"}, `${capitalize(actionName)} Graph`),
            label({for: "actionSymbol"}, "Symbol"),
            input({
                id: "actionSymbol",
                class: INPUT_CLASS,
                type: "text",
                placeholder: 'Symbol (SPY, QQQ, etc.)',
                onchange: e => symbol.val = e.target.value,
            }),
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
            label({for: "actionDate"}, "Date"),
            dateInput,
            ...inputElems,
            buttons,
        )
    )

    const datePickerOpts = {
        autohide: true,
        container: "#foreground",
        daysOfWeekDisabled: [0, 6], // disable saturday and sunday
        format: "yyyy-mm-dd",
        maxDate: new Date(),
        todayHighlight: true,
    }
    const datepicker = new Datepicker(dateInput, datePickerOpts)
  
    van.add(document.body, () => closed.val ? null : modal);
}

export {GraphActionModal, INPUT_CLASS};
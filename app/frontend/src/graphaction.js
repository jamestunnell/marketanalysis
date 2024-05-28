import van from 'vanjs-core'
import Datepicker from 'flowbite-datepicker/Datepicker'

import { AppErrorAlert} from './apperror.js'
import { Button, ButtonIcon } from "./buttons.js"
import capitalize from './capitalize.js'
import { DownloadJSON } from "./download.js"
import { IconDownload, IconClose, IconPlay, IconPlot } from './icons.js'
import { ModalBackground } from './modal.js'
import { PlotRecordingModal } from './plot.js'

const {div, input, label} = van.tags

const INPUT_CLASS = "block px-3 py-3 border border-gray-200 rounded-md focus:border-indigo-500 focus:outline-none focus:ring"

const GraphActionModal = ({actionName, graph, inputElems, runDisabled, doAction}) => {
    const closed = van.state(false)
    const completed = van.state(false)
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

    const closeBtn = ButtonIcon({icon: IconClose(), onclick: ()=> closed.val = true})
    const runBtn = Button({
        disabled: van.derive(() => {
            return runDisabled.val || (date.val.length === 0) || (symbol.val.length === 0)
        }),
        child: [IconPlay, ` ${capitalize(actionName)}`],
        onclick: () => {
            completed.val = false
    
            doAction({date: date.val, symbol: symbol.val}).then(obj => {
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
        onclick: ()=> PlotRecordingModal(recording.val),
    })
    const downloadBtn = ButtonIcon({
        icon: IconDownload(),
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
        div(
            { class: van.derive(() => `flex flex-row ${!completed.val ? "hidden" : ""}`) },
            plotBtn,
            downloadBtn,    
        )
    )

    closeBtn.classList.add("self-end")
    runBtn.classList.add("self-center")

    const modal = ModalBackground(
        div(
            {id: "foreground", class: "flex flex-col block p-16 rounded-lg bg-white min-w-[25%] max-w-[50%]"},
            closeBtn,
            label({for: "actionSymbol"}, "Symbol"),
            input({
                id: "actionSymbol",
                class: INPUT_CLASS,
                type: "text",
                placeholder: 'Symbol (SPY, QQQ, etc.)',
                onchange: e => symbol.val = e.target.value,
            }),  
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
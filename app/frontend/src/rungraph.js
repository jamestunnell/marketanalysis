import van from 'vanjs-core'
import Datepicker from 'flowbite-datepicker/Datepicker'

import { AppErrorAlert} from './apperror.js'
import { PostJSON } from './backend.js'
import { Button, ButtonCancel, ButtonIcon } from "./buttons.js"
import { Download } from "./download.js"
import { IconDownload, IconClose, IconPlay, IconPlot } from './icons.js'
import { ModalBackground, ModalForeground } from './modal.js'
import { PlotModal } from './plot.js'

const {div, input, p, label} = van.tags

const runGraph = ({id, date, symbol}) => {
    return new Promise((resolve, reject) => {
        const route = `/graphs/${id}/run-day`
        const object = {date, symbol, format: "ndjson"}
        const options = {accept: 'application/x-ndjson'}

        console.log("running graph", object)

        PostJSON({route, object, options}).then(resp => {
            if (resp.status != 200) {
                resp.json().then(appErr => {
                    console.log("failed to run graph", appErr);
    
                    reject(appErr);    
                })
            }

            resp.text().then(text => resolve(text))
        }).catch(err => {
            console.log("failed to make run graph request", err)
            
            reject({
                title: "Action Failed",
                message: "failed to make run graph request",
                details: [err.message],
            })
        });
    });
}

const RunGraph = (graph) => {
    const closed = van.state(false)
    const completed = van.state(false)
    const resultText = van.state("")
    const inputClass = "block px-3 py-3 border border-gray-200 rounded-md focus:border-indigo-500 focus:outline-none focus:ring"
    const today = new Date()
    const date = van.state("")
    const symbol = van.state("")
    const dateInput = input({
        id: "runDate",
        class: inputClass,
        type: "text",
        placeholder: 'Select date',
    });

    dateInput.addEventListener('changeDate', (e) => {
        date.val = dateInput.value
    })

    const onRun = () => {
        if (date.val === "") {
            AppErrorAlert({
                title: "Invalid Input",
                message: "date is empty",
                details: [],
            })

            return
        }

        runGraph({
            id: graph.id,
            date: date.val,
            symbol: symbol.val,
        }).then(text => {
            console.log("run graph succeeded")
            
            resultText.val = text
            completed.val = true
        }).catch(appErr => {
            AppErrorAlert(appErr)
        })
    }

    const closeBtn = ButtonIcon({icon: IconClose(), onclick: ()=> closed.val = true})
    const runBtn = Button({
        disabled: van.derive(() => date.val.length == 0),
        child: [IconPlay, " Run"],
        onclick: onRun,
    })
    const plotBtn = ButtonIcon({
        icon: IconPlot(),
        onclick: ()=> PlotModal({text: resultText.val, format: "ndjson"}),
    })
    const downloadBtn = ButtonIcon({
        icon: IconDownload(),
        onclick: ()=> {
            Download({
                filename: `${graph.name}_${date.val}.ndjson`,
                blob: new Blob([resultText.val]),
            })
        },
    })
    const resultsButtons = div(
        {
            class: van.derive(() => `flex flex-row justify-center ${!completed.val ? "hidden" : ""}`),
        },
        plotBtn,
        downloadBtn,
    )

    closeBtn.classList.add("self-end")
    runBtn.classList.add("self-center")

    const modal = ModalBackground(
        ModalForeground(
            {id: "foreground"},
            div(
                {id: "modalContent", class: "flex flex-col rounded-md space-y-4"},
                closeBtn,
                p({class: "text-lg font-medium font-bold text-center"}, "Run Graph"),
                div(
                    {class: "flex flex-col"},
                    label({for: "runDate"}, "Date"),
                    dateInput,
                ),
                div(
                    {class: "flex flex-col"},
                    label({for: "runSymbol"}, "Symbol"),
                    input({
                        id: "runSymbol",
                        class: inputClass,
                        type: "text",
                        placeholder: 'Security symbol (SPY, QQQ, etc.)',
                        oninput: e => symbol.val = e.target.value,
                    }),                        
                ),
                runBtn,
                resultsButtons,
            ),
        ),
    )
  
    van.add(document.body, () => closed.val ? null : modal);

    const datePickerOpts = {
        autohide: true,
        container: "#foreground",
        daysOfWeekDisabled: [0, 6], // disable saturday and sunday
        format: "yyyy-mm-dd",
        maxDate: today,
        todayHighlight: true,
    }
    const datepicker = new Datepicker(dateInput, datePickerOpts)
}

export {RunGraph};